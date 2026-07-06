package safety

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Zone classifies path safety for scan, display, and delete operations.
type Zone string

const (
	ZoneForbidden   Zone = "forbidden"
	ZoneCriticalOS  Zone = "critical_os"
	ZoneDiagnostic  Zone = "diagnostic"
	ZoneCaution     Zone = "caution"
	ZoneReview      Zone = "review"
	ZoneMaintenance Zone = "maintenance"
	ZoneNormal      Zone = "normal"
)

// ClassifyPath returns the safety zone for an absolute or relative path.
func ClassifyPath(path string) Zone {
	if path == "" {
		return ZoneNormal
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return ZoneNormal
	}
	abs = filepath.Clean(abs)
	lower := strings.ToLower(abs)
	base := strings.ToLower(filepath.Base(abs))

	if runtime.GOOS == "windows" {
		return classifyWindows(abs, lower, base)
	}
	return classifyUnix(abs, lower, base)
}

func classifyWindows(abs, lower, base string) Zone {
	for _, name := range []string{"pagefile.sys", "hiberfil.sys", "swapfile.sys", "memory.dmp"} {
		if base == name {
			return ZoneForbidden
		}
	}
	for _, seg := range []string{"system32", "syswow64", "winsxs", "servicing"} {
		if pathHasSegment(lower, seg) {
			return ZoneForbidden
		}
	}
	for _, seg := range []string{"minidump", "livekernelreports"} {
		if pathHasSegment(lower, seg) {
			return ZoneDiagnostic
		}
	}
	if strings.Contains(lower, "crashdumps") {
		return ZoneDiagnostic
	}
	if pathHasSegment(lower, "program files") || pathHasSegment(lower, "program files (x86)") {
		return ZoneCriticalOS
	}
	if isWindowsDir(lower, "windows") {
		if pathHasSegment(lower, "temp") && !pathHasSegment(lower, "system32") {
			return ZoneMaintenance
		}
		return ZoneCriticalOS
	}
	if isTempPath(lower) {
		return ZoneMaintenance
	}
	return ZoneNormal
}

func classifyUnix(abs, lower, base string) Zone {
	for _, prefix := range []string{"/proc", "/sys", "/dev"} {
		if lower == prefix || strings.HasPrefix(lower, prefix+"/") {
			return ZoneForbidden
		}
	}
	for _, prefix := range []string{"/boot", "/usr", "/bin", "/sbin", "/etc"} {
		if lower == prefix || strings.HasPrefix(lower, prefix+"/") {
			return ZoneCriticalOS
		}
	}
	if strings.HasPrefix(lower, "/lib") {
		return ZoneCriticalOS
	}
	for _, seg := range []string{"/var/crash", "/var/lib/systemd/coredump"} {
		if lower == seg || strings.HasPrefix(lower, seg+"/") {
			return ZoneDiagnostic
		}
	}
	if base == "vmcore" {
		return ZoneDiagnostic
	}
	if lower == "/tmp" || strings.HasPrefix(lower, "/tmp/") {
		return ZoneMaintenance
	}
	return ZoneNormal
}

func pathHasSegment(lowerPath, segment string) bool {
	segment = strings.ToLower(segment)
	if runtime.GOOS == "windows" {
		return strings.Contains(lowerPath, "\\"+segment+"\\") ||
			strings.HasSuffix(lowerPath, "\\"+segment) ||
			strings.Contains(lowerPath, "/"+segment+"/") ||
			strings.HasSuffix(lowerPath, "/"+segment)
	}
	return strings.Contains(lowerPath, "/"+segment+"/") || strings.HasSuffix(lowerPath, "/"+segment)
}

func isWindowsDir(lowerPath, name string) bool {
	return pathHasSegment(lowerPath, name) || strings.HasSuffix(lowerPath, `\`+name) || strings.HasSuffix(lowerPath, `/`+name)
}

func isTempPath(lower string) bool {
	if strings.Contains(lower, `\temp\`) || strings.HasSuffix(lower, `\temp`) {
		return true
	}
	temp := os.TempDir()
	if temp != "" {
		tl := strings.ToLower(filepath.Clean(temp))
		return lower == tl || strings.HasPrefix(lower, tl+string(os.PathSeparator))
	}
	return false
}

// ShouldSkipScan returns true when the scanner must not descend into path.
func ShouldSkipScan(path string) bool {
	return ClassifyPath(path) == ZoneForbidden
}

// IsDeletable returns whether disk-tool may delete path in this zone.
func IsDeletable(zone Zone) bool {
	switch zone {
	case ZoneReview, ZoneMaintenance, ZoneNormal, ZoneCaution:
		return true
	default:
		return false
	}
}

// IsDriveRoot reports whether path is a volume or filesystem root.
func IsDriveRoot(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	abs = filepath.Clean(abs)
	if runtime.GOOS == "windows" {
		if len(abs) == 3 && abs[1] == ':' && (abs[2] == '\\' || abs[2] == '/') {
			return true
		}
		return false
	}
	return abs == string(os.PathSeparator)
}

// ZoneLabel returns a short human label for UI.
func ZoneLabel(zone Zone) string {
	switch zone {
	case ZoneForbidden:
		return "Forbidden"
	case ZoneCriticalOS:
		return "Critical OS"
	case ZoneDiagnostic:
		return "Diagnostic"
	case ZoneCaution:
		return "Caution"
	case ZoneReview:
		return "Review"
	case ZoneMaintenance:
		return "Maintenance"
	default:
		return "Normal"
	}
}

// ApplyCandidateZone sets zone and deletable on a cleanup candidate from its path and risk.
func ApplyCandidateZone(path, risk string) (Zone, bool) {
	zone := ClassifyPath(path)
	if zone == ZoneNormal {
		if risk == "caution" {
			zone = ZoneCaution
		} else if risk == "review" {
			zone = ZoneReview
		}
	}
	return zone, IsDeletable(zone)
}

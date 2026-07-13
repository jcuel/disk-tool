package safety

import (
	"path/filepath"
	"strings"
)

// Virtual disk / disk-image extensions that must never be deleted via disk-tool.
// Deleting these can destroy VMs, Hyper-V disks, WSL images, or system recovery media.
var diskImageExt = map[string]bool{
	".vhd":   true,
	".vhdx":  true,
	".avhd":  true,
	".avhdx": true,
	".vmdk":  true,
	".vdi":   true,
	".qcow":  true,
	".qcow2": true,
	".wim":   true,
	".esd":   true,
	".vfd":   true,
}

// IsDiskImagePath reports whether path is a virtual disk or disk image file.
func IsDiskImagePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return diskImageExt[ext]
}

// CanDeletePath reports whether disk-tool may delete path (zone + disk-image rules).
func CanDeletePath(path string) (ok bool, reason string) {
	if IsDiskImagePath(path) {
		return false, "virtual disk / disk image — deletion disabled"
	}
	zone := ClassifyPath(path)
	if !IsDeletable(zone) {
		return false, "protected zone: " + string(zone)
	}
	return true, ""
}

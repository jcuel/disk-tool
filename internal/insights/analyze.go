package insights

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/safety"
)

// Options controls age-based and analysis thresholds.
type Options struct {
	AgeThresholdDays int
	MinSizeBytes     int64
}

const defaultAgeDays = 90
const defaultMinStale = 50 * 1024 * 1024

var dirRules = []struct {
	name     string
	category model.CleanupCategory
	hint     string
	risk     string
}{
	{"node_modules", model.CategoryNodeModules, "Dependencies from npm/yarn/pnpm — safe to delete; reinstall with npm install", "review"},
	{".venv", model.CategoryPythonVenv, "Python virtual environment — recreate with python -m venv", "review"},
	{"venv", model.CategoryPythonVenv, "Python virtual environment — recreate if still needed", "review"},
	{"__pycache__", model.CategoryBuildOutput, "Python bytecode cache — regenerated automatically", "review"},
	{"target", model.CategoryBuildOutput, "Rust/Java build output — rebuild with your build tool", "review"},
	{"dist", model.CategoryBuildOutput, "Build output folder — rebuild from source", "review"},
	{"build", model.CategoryBuildOutput, "Build artifacts — rebuild from source", "review"},
	{".gradle", model.CategoryPackageCache, "Gradle cache — may re-download on next build", "caution"},
	{".m2", model.CategoryPackageCache, "Maven local repository — shared across Java projects", "caution"},
	{".npm", model.CategoryPackageCache, "npm cache — can often be cleared with npm cache clean", "caution"},
	{".cache", model.CategoryPackageCache, "Application cache — review before deleting", "caution"},
	{"vendor", model.CategoryBuildOutput, "Go/PHP vendor tree — may be regenerable depending on project", "review"},
}

var installerExt = map[string]bool{
	".exe": true, ".msi": true, ".zip": true, ".dmg": true,
	".iso": true, ".tar": true, ".gz": true, ".7z": true,
	".rar": true, ".deb": true, ".rpm": true, ".pkg": true,
}

func Analyze(job *model.ScanJob) *model.InsightsReport {
	return AnalyzeWithOptions(job, jobInsightsOptions(job))
}

func jobInsightsOptions(job *model.ScanJob) Options {
	opts := Options{AgeThresholdDays: defaultAgeDays, MinSizeBytes: defaultMinStale}
	if job == nil {
		return opts
	}
	if job.InsightsConfig.AgeThresholdDays > 0 {
		opts.AgeThresholdDays = job.InsightsConfig.AgeThresholdDays
	}
	if job.InsightsConfig.MinSizeBytes > 0 {
		opts.MinSizeBytes = job.InsightsConfig.MinSizeBytes
	}
	return opts
}

func AnalyzeWithOptions(job *model.ScanJob, opts Options) *model.InsightsReport {
	if job == nil || job.Tree == nil {
		return &model.InsightsReport{Summary: "No scan data yet."}
	}

	report := &model.InsightsReport{
		TopConsumers:      []model.TopConsumer{},
		CleanupCandidates: []model.CleanupCandidate{},
	}
	rootSize := job.Tree.Size
	if rootSize <= 0 {
		rootSize = 1
	}

	for _, c := range job.Tree.Children {
		report.TopConsumers = append(report.TopConsumers, model.TopConsumer{
			Name: c.Name,
			Path: c.Path,
			Size: c.Size,
			Pct:  float64(c.Size) * 100 / float64(rootSize),
		})
	}

	seen := map[string]bool{}
	var walk func(*model.ScanNode)
	walk = func(n *model.ScanNode) {
		if n == nil {
			return
		}
		base := strings.ToLower(n.Name)
		for _, rule := range dirRules {
			if base == rule.name && n.Size > 0 && !seen[n.Path] {
				seen[n.Path] = true
				addCandidate(report, model.CleanupCandidate{
					Category: rule.category,
					Path:     n.Path,
					Size:     n.Size,
					Hint:     rule.hint,
					Risk:     rule.risk,
				})
			}
		}
		if isDownloadsDir(n.Path) {
			for _, f := range job.LargestFiles {
				if strings.HasPrefix(strings.ToLower(f.Path), strings.ToLower(n.Path)) {
					addDownloadCandidate(report, seen, f)
				}
			}
		}
		for _, ch := range n.Children {
			walk(ch)
		}
	}
	walk(job.Tree)

	for _, f := range job.LargestFiles {
		if isDownloadsPath(f.Path) {
			addDownloadCandidate(report, seen, f)
		}
		addStaleCandidate(report, seen, f, opts)
	}

	tagCandidates(report)
	report.SafetyGrid = buildSafetyGrid(job.Root, report.CleanupCandidates)
	report.Summary = buildSummary(job, report)
	report.TicketText = buildTicketText(job, report)
	return report
}

func addCandidate(report *model.InsightsReport, c model.CleanupCandidate) {
	zone, deletable := safety.ApplyCandidateZone(c.Path, c.Risk)
	if !deletable {
		c.Hint = c.Hint + " — protected zone; deletion disabled"
		report.CleanupCandidates = append(report.CleanupCandidates, tagCandidate(c, zone, deletable))
		return
	}
	report.CleanupCandidates = append(report.CleanupCandidates, tagCandidate(c, zone, deletable))
	report.TotalReclaimable += c.Size
}

func tagCandidate(c model.CleanupCandidate, zone safety.Zone, deletable bool) model.CleanupCandidate {
	c.Zone = string(zone)
	c.Deletable = deletable
	return c
}

func tagCandidates(report *model.InsightsReport) {
	var reclaim int64
	for i, c := range report.CleanupCandidates {
		if c.Zone == "" {
			z, d := safety.ApplyCandidateZone(c.Path, c.Risk)
			c.Zone = string(z)
			c.Deletable = d
		}
		if c.Deletable {
			reclaim += c.Size
		}
		report.CleanupCandidates[i] = c
	}
	report.TotalReclaimable = reclaim
}

func buildSafetyGrid(root string, candidates []model.CleanupCandidate) *model.SafetyGrid {
	grid := &model.SafetyGrid{
		Zones:     make(map[string]model.SafetyZoneStats),
		DriveRoot: safety.IsDriveRoot(root),
	}
	for _, c := range candidates {
		key := c.Zone
		if key == "" {
			key = string(safety.ZoneNormal)
		}
		st := grid.Zones[key]
		st.Count++
		st.Bytes += c.Size
		grid.Zones[key] = st
		if !c.Deletable {
			grid.ProtectedBytes += c.Size
		}
	}
	return grid
}

func addStaleCandidate(report *model.InsightsReport, seen map[string]bool, f model.FileEntry, opts Options) {
	if safety.IsDiskImagePath(f.Path) || safety.IsDiskImagePath(f.Name) {
		return
	}
	if seen[f.Path] {
		return
	}
	minSize := opts.MinSizeBytes
	if minSize <= 0 {
		minSize = defaultMinStale
	}
	if f.Size < minSize {
		return
	}
	ageDays := opts.AgeThresholdDays
	if ageDays <= 0 {
		ageDays = defaultAgeDays
	}
	info, err := os.Stat(f.Path)
	if err != nil {
		return
	}
	age := time.Since(info.ModTime())
	if age < time.Duration(ageDays)*24*time.Hour {
		return
	}
	if !isStaleLocation(f.Path) {
		return
	}
	seen[f.Path] = true
	hint := fmt.Sprintf("Not modified in %d+ days — review before deleting", ageDays)
	addCandidate(report, model.CleanupCandidate{
		Category: model.CategoryStaleLarge,
		Path:     f.Path,
		Size:     f.Size,
		Hint:     hint,
		Risk:     "review",
	})
}

func isStaleLocation(path string) bool {
	if isDownloadsPath(path) || isTempPath(path) {
		return true
	}
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.Contains(lower, "/temp/") || strings.Contains(lower, "/tmp/")
}

func isTempPath(path string) bool {
	z := safety.ClassifyPath(path)
	return z == safety.ZoneMaintenance
}

func addDownloadCandidate(report *model.InsightsReport, seen map[string]bool, f model.FileEntry) {
	if safety.IsDiskImagePath(f.Path) || safety.IsDiskImagePath(f.Name) {
		return
	}
	ext := strings.ToLower(filepath.Ext(f.Name))
	if !installerExt[ext] && f.Size < 50*1024*1024 {
		return
	}
	if seen[f.Path] {
		return
	}
	seen[f.Path] = true
	hint := "Installer or archive in Downloads — remove if no longer needed"
	if f.Size >= 100*1024*1024 {
		hint = "Large installer/archive — likely safe to remove if already installed"
	}
	addCandidate(report, model.CleanupCandidate{
		Category: model.CategoryDownload,
		Path:     f.Path,
		Size:     f.Size,
		Hint:     hint,
		Risk:     "review",
	})
}

func isDownloadsDir(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return base == "downloads" || base == "descargas"
}

func isDownloadsPath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	return strings.Contains(lower, "/downloads/") || strings.Contains(lower, "/descargas/") ||
		strings.Contains(lower, "\\downloads\\") || strings.Contains(lower, "\\descargas\\")
}

func buildSummary(job *model.ScanJob, r *model.InsightsReport) string {
	if len(r.TopConsumers) == 0 {
		return fmt.Sprintf("Scanned %s — drill into folders to find where space is used.", job.Root)
	}
	top := r.TopConsumers[0]
	msg := fmt.Sprintf("%s uses %.1f%% of scanned space (%s). Found %d cleanup candidate(s) worth ~%s if reviewed.",
		top.Name, top.Pct, formatBytes(top.Size), len(r.CleanupCandidates), formatBytes(r.TotalReclaimable))
	if r.SafetyGrid != nil && r.SafetyGrid.DriveRoot {
		msg += " Scanning a full drive includes protected OS zones — prefer your user profile for cleanup."
	}
	return msg
}

func buildTicketText(job *model.ScanJob, r *model.InsightsReport) string {
	var b strings.Builder
	b.WriteString("Disk usage report\n")
	b.WriteString("=================\n")
	b.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("Scan root: %s\n", job.Root))
	b.WriteString(fmt.Sprintf("Total scanned: %s (%d files indexed)\n\n", formatBytes(job.Tree.Size), job.Tree.FileCount))
	b.WriteString("Top space consumers\n")
	b.WriteString("-------------------\n")
	for i, c := range r.TopConsumers {
		if i >= 10 {
			break
		}
		b.WriteString(fmt.Sprintf("  %.1f%%  %s  (%s)\n", c.Pct, c.Name, formatBytes(c.Size)))
	}
	if len(r.CleanupCandidates) > 0 {
		b.WriteString("\nSuggested cleanup review\n")
		b.WriteString("------------------------\n")
		for i, c := range r.CleanupCandidates {
			if i >= 20 {
				b.WriteString(fmt.Sprintf("  ... and %d more\n", len(r.CleanupCandidates)-20))
				break
			}
			b.WriteString(fmt.Sprintf("  [%s/%s] %s — %s\n", c.Category, c.Zone, formatBytes(c.Size), c.Path))
			b.WriteString(fmt.Sprintf("         %s\n", c.Hint))
		}
		b.WriteString(fmt.Sprintf("\nPotential reclaimable (deletable only): ~%s\n", formatBytes(r.TotalReclaimable)))
	}
	b.WriteString("\nNote: Sizes are based on scanned portions of the tree. Protected OS and diagnostic zones cannot be deleted via disk-tool.\n")
	return b.String()
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

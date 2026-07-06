package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	CleanupStatusDeleted           = "deleted"
	CleanupStatusSkippedLocked     = "skipped_locked"
	CleanupStatusSkippedMissing    = "skipped_missing"
	CleanupStatusSkippedOutside    = "skipped_outside_root"
	CleanupStatusSkippedScanRoot   = "skipped_scan_root"
	CleanupStatusSkippedProtected  = "skipped_protected_zone"
	CleanupStatusFailed            = "failed"
	CleanupStatusWouldDelete       = "would_delete"
)

type CleanupItemResult struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status"`
	Reason   string `json:"reason,omitempty"`
}

type CleanupReport struct {
	DryRun         bool                `json:"dryRun"`
	StartedAt      time.Time           `json:"startedAt"`
	FinishedAt     time.Time           `json:"finishedAt"`
	TotalRequested int                 `json:"totalRequested"`
	BytesReclaimed int64               `json:"bytesReclaimed"`
	Results        []CleanupItemResult `json:"results"`
	ReportText     string              `json:"reportText"`
}

type CleanupRequest struct {
	Paths         []string `json:"paths"`
	DryRun        bool     `json:"dryRun"`
	Confirm       bool     `json:"confirm"`
	ConfirmPhrase string   `json:"confirmPhrase"`
}

func BuildCleanupReportText(report *CleanupReport) string {
	if report == nil {
		return ""
	}
	var b strings.Builder
	mode := "EXECUTED"
	if report.DryRun {
		mode = "DRY RUN"
	}
	fmt.Fprintf(&b, "disk-tool cleanup report (%s)\n", mode)
	fmt.Fprintf(&b, "Started: %s\n", report.StartedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Finished: %s\n", report.FinishedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "Requested: %d | Reclaimed: %s\n\n", report.TotalRequested, formatCleanupBytes(report.BytesReclaimed))
	for _, r := range report.Results {
		line := fmt.Sprintf("[%s] %s (%s)", r.Status, r.Path, formatCleanupBytes(r.Size))
		if r.Category != "" {
			line += " " + r.Category
		}
		if r.Reason != "" {
			line += " — " + r.Reason
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func formatCleanupBytes(n int64) string {
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

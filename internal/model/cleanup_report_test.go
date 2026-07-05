package model

import (
	"strings"
	"testing"
	"time"
)

func TestBuildCleanupReportText_includesResults(t *testing.T) {
	report := &CleanupReport{
		DryRun:         false,
		StartedAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		FinishedAt:     time.Date(2026, 1, 1, 0, 1, 0, 0, time.UTC),
		TotalRequested: 1,
		BytesReclaimed: 1024,
		Results: []CleanupItemResult{
			{Path: "/x", Size: 1024, Status: CleanupStatusDeleted, Category: "node_modules"},
		},
	}
	text := BuildCleanupReportText(report)
	if !strings.Contains(text, "EXECUTED") {
		t.Fatal("missing executed marker")
	}
	if !strings.Contains(text, "/x") {
		t.Fatal("missing path")
	}
}

package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcuel/disk-tool/internal/model"
)

func TestRunCleanup_dryRun(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "node_modules")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "x"), []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}

	job := &model.ScanJob{
		Root: dir,
		Insights: &model.InsightsReport{
			CleanupCandidates: []model.CleanupCandidate{
				{Category: model.CategoryNodeModules, Path: target, Size: 3},
			},
		},
	}

	report, err := RunCleanup(job, model.CleanupRequest{
		Paths:  []string{target},
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !report.DryRun {
		t.Fatal("expected dry run")
	}
	if len(report.Results) != 1 || report.Results[0].Status != model.CleanupStatusWouldDelete {
		t.Fatalf("unexpected results: %+v", report.Results)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatal("dry run should not delete")
	}
}

func TestRunCleanup_executeRequiresConfirm(t *testing.T) {
	dir := t.TempDir()
	job := &model.ScanJob{Root: dir}
	_, err := RunCleanup(job, model.CleanupRequest{
		Paths:  []string{filepath.Join(dir, "missing")},
		DryRun: false,
	})
	if err != errConfirmRequired {
		t.Fatalf("expected confirm required, got %v", err)
	}
}

func TestRunCleanup_executeRequiresPhrase(t *testing.T) {
	dir := t.TempDir()
	job := &model.ScanJob{Root: dir}
	_, err := RunCleanup(job, model.CleanupRequest{
		Paths:   []string{filepath.Join(dir, "missing")},
		DryRun:  false,
		Confirm: true,
	})
	if err != errCleanupConfirmPhrase {
		t.Fatalf("expected phrase error, got %v", err)
	}
}

func TestRunCleanup_executeDeletes(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "cache")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "a"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	job := &model.ScanJob{
		Root: dir,
		Insights: &model.InsightsReport{
			CleanupCandidates: []model.CleanupCandidate{
				{Category: model.CategoryPackageCache, Path: target, Size: 4},
			},
			TotalReclaimable: 4,
		},
	}

	report, err := RunCleanup(job, model.CleanupRequest{
		Paths:         []string{target},
		DryRun:        false,
		Confirm:       true,
		ConfirmPhrase: "DELETE",
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.BytesReclaimed == 0 {
		t.Fatal("expected reclaimed bytes")
	}
	if report.Results[0].Status != model.CleanupStatusDeleted {
		t.Fatalf("unexpected status: %s", report.Results[0].Status)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatal("target should be deleted")
	}
	pruneDeletedCandidates(job, report)
	if len(job.Insights.CleanupCandidates) != 0 {
		t.Fatal("expected candidates pruned")
	}
}

func TestRunCleanup_skipsScanRoot(t *testing.T) {
	dir := t.TempDir()
	job := &model.ScanJob{Root: dir}
	report, err := RunCleanup(job, model.CleanupRequest{
		Paths:  []string{dir},
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.Results[0].Status != model.CleanupStatusSkippedScanRoot {
		t.Fatalf("unexpected status: %s", report.Results[0].Status)
	}
}

func TestRunCleanup_skipsOutsideRoot(t *testing.T) {
	dir := t.TempDir()
	job := &model.ScanJob{Root: dir}
	report, err := RunCleanup(job, model.CleanupRequest{
		Paths:  []string{filepath.Join(dir, "..", "outside")},
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.Results[0].Status != model.CleanupStatusSkippedOutside {
		t.Fatalf("unexpected status: %s", report.Results[0].Status)
	}
}

func TestBuildCleanupReportText(t *testing.T) {
	report := &model.CleanupReport{
		DryRun:         true,
		TotalRequested: 1,
		Results: []model.CleanupItemResult{
			{Path: "/tmp/x", Size: 10, Status: model.CleanupStatusWouldDelete},
		},
	}
	text := model.BuildCleanupReportText(report)
	if text == "" || !strings.Contains(text, "DRY RUN") {
		t.Fatalf("unexpected report text: %q", text)
	}
}

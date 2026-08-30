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
	if len(job.Insights.CleanupCandidates) != 0 {
		t.Fatal("expected candidates pruned")
	}
}

func TestApplyPostDeleteUpdates_prunesLargestFilesAndTree(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	target := filepath.Join(sub, "big.bin")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("0123456789"), 0o644); err != nil {
		t.Fatal(err)
	}

	child := &model.ScanNode{
		Name: "big.bin", Path: target, Size: 10, FileCount: 1, IsDir: false, Scanned: true,
	}
	subNode := &model.ScanNode{
		Name: "subdir", Path: sub, Size: 10, FileCount: 1, IsDir: true, Scanned: true,
		Children: []*model.ScanNode{child},
	}
	root := &model.ScanNode{
		Name: filepath.Base(dir), Path: dir, Size: 10, FileCount: 1, IsDir: true, Scanned: true,
		Children: []*model.ScanNode{subNode},
	}

	job := &model.ScanJob{
		Root: dir,
		Tree: root,
		LargestFiles: []model.FileEntry{
			{Path: target, Name: "big.bin", Size: 10},
			{Path: filepath.Join(dir, "other.txt"), Name: "other.txt", Size: 1},
		},
		Insights: &model.InsightsReport{
			CleanupCandidates: []model.CleanupCandidate{
				{Category: model.CategoryPackageCache, Path: target, Size: 10},
			},
			TotalReclaimable: 10,
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
	if report.Results[0].Status != model.CleanupStatusDeleted {
		t.Fatalf("expected deleted, got %s", report.Results[0].Status)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatal("file should be gone on disk (Bucket B check)")
	}
	if len(job.LargestFiles) != 1 || job.LargestFiles[0].Name != "other.txt" {
		t.Fatalf("largest files not pruned: %+v", job.LargestFiles)
	}
	if len(subNode.Children) != 0 {
		t.Fatal("tree child should be removed")
	}
	if subNode.Size != 0 {
		t.Fatalf("parent size should recompute to 0, got %d", subNode.Size)
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

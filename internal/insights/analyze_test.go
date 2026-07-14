package insights

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcuel/disk-tool/internal/model"
)

func TestAnalyzeNodeModules(t *testing.T) {
	root := t.TempDir()
	nm := filepath.Join(root, "proj", "node_modules")
	if err := os.MkdirAll(nm, 0o755); err != nil {
		t.Fatal(err)
	}

	job := &model.ScanJob{
		Root: root,
		Tree: &model.ScanNode{
			Name: "root", Path: root, Size: 500_000_000, IsDir: true, Scanned: true,
			Children: []*model.ScanNode{
				{
					Name: "proj", Path: filepath.Join(root, "proj"), Size: 400_000_000, IsDir: true, Scanned: true,
					Children: []*model.ScanNode{
						{Name: "node_modules", Path: nm, Size: 380_000_000, IsDir: true, Scanned: true},
					},
				},
			},
		},
	}
	r := AnalyzeWithOptions(job, Options{AgeThresholdDays: 90, MinSizeBytes: 50 * 1024 * 1024, SkipDocker: true})
	if len(r.CleanupCandidates) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(r.CleanupCandidates))
	}
	if r.CleanupCandidates[0].Category != model.CategoryNodeModules {
		t.Fatalf("wrong category %s", r.CleanupCandidates[0].Category)
	}
	if r.TicketText == "" {
		t.Fatal("expected ticket text")
	}
}

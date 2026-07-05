package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestScanOverview(t *testing.T) {
	root := t.TempDir()
	mkfile(t, filepath.Join(root, "a.txt"), 100)
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	mkfile(t, filepath.Join(sub, "b.txt"), 250)

	sc := New()
	tree, _, err := sc.ScanOverview(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if tree.Size != 350 {
		t.Fatalf("expected 350 bytes, got %d", tree.Size)
	}
	if len(tree.Children) != 1 {
		t.Fatalf("expected 1 top-level child, got %d", len(tree.Children))
	}
	if tree.Children[0].Scanned {
		t.Fatal("overview child should not be scanned")
	}
	if tree.Children[0].Size != 250 {
		t.Fatalf("expected subdir aggregate 250, got %d", tree.Children[0].Size)
	}
}

func TestScanBranchDepth(t *testing.T) {
	root := t.TempDir()
	// root/a/leaf/file.txt
	a := filepath.Join(root, "a")
	leaf := filepath.Join(a, "leaf")
	if err := os.MkdirAll(leaf, 0o755); err != nil {
		t.Fatal(err)
	}
	mkfile(t, filepath.Join(leaf, "f.txt"), 500)

	sc := New()
	branch, _, err := sc.ScanBranch(context.Background(), a, root, 2, Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if !branch.Scanned {
		t.Fatal("branch root should be scanned")
	}
	if len(branch.Children) != 1 {
		t.Fatalf("expected leaf dir, got %d children", len(branch.Children))
	}
	leafNode := branch.Children[0]
	if leafNode.Expandable {
		t.Fatal("leaf at depth limit with file inside should not need expand on dir if scanned")
	}
}

func TestScanFull(t *testing.T) {
	root := t.TempDir()
	mkfile(t, filepath.Join(root, "a.txt"), 100)
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	mkfile(t, filepath.Join(sub, "b.txt"), 250)

	sc := New()
	tree, largest, err := sc.Scan(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if tree.Size != 350 {
		t.Fatalf("expected 350 bytes, got %d", tree.Size)
	}
	if !tree.Children[0].Scanned {
		t.Fatal("full scan should mark children scanned")
	}
	if len(largest) != 2 {
		t.Fatalf("expected 2 largest files, got %d", len(largest))
	}
}

func TestScanCancel(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 20; i++ {
		dir := filepath.Join(root, "dir"+string(rune('0'+i%10)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		mkfile(t, filepath.Join(dir, "f.txt"), 10)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sc := New()
	_, _, err := sc.ScanOverview(ctx, Options{Root: root})
	if err == nil {
		t.Fatal("expected cancel error")
	}
}

func mkfile(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if size > 0 {
		if _, err := f.Write(make([]byte, size)); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()
}

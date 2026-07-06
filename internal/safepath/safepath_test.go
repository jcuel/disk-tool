package safepath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRoot(t *testing.T) {
	_, err := ResolveRoot("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	root, err := ResolveRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if root == "" {
		t.Fatal("expected abs path")
	}
	f := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = ResolveRoot(f)
	if err == nil {
		t.Fatal("expected error for file path")
	}
}

func TestUnderRoot(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "sub")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := UnderRoot(root, child)
	if err != nil {
		t.Fatal(err)
	}
	if p != child {
		t.Fatalf("got %s", p)
	}
	_, err = UnderRoot(root, filepath.Join(root, "..", "outside"))
	if err == nil {
		t.Fatal("expected outside root error")
	}
}

func TestRel(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "sub", "file.txt")
	if err := os.MkdirAll(filepath.Dir(child), 0o755); err != nil {
		t.Fatal(err)
	}
	rel, err := Rel(root, child)
	if err != nil {
		t.Fatal(err)
	}
	if rel != filepath.Join("sub", "file.txt") {
		t.Fatalf("got %q", rel)
	}
	_, err = Rel(root, filepath.Join(root, "..", "escape"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpenRoot(t *testing.T) {
	root := t.TempDir()
	r, abs, err := OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	if abs != root {
		t.Fatalf("abs %q", abs)
	}
	entries, err := r.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	defer entries.Close()
	if _, err := entries.ReadDir(-1); err != nil {
		t.Fatal(err)
	}
}

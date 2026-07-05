package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRoot(t *testing.T) {
	_, err := ValidateRoot("")
	if err == nil {
		t.Fatal("expected error for empty root")
	}
	root, err := ValidateRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if root == "" {
		t.Fatal("expected abs path")
	}
}

func TestPathWithinRoot(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "sub")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := PathWithinRoot(root, child)
	if err != nil {
		t.Fatal(err)
	}
	if p != child {
		t.Fatalf("got %s", p)
	}
	_, err = PathWithinRoot(root, filepath.Join(root, "..", "outside"))
	if err == nil {
		t.Fatal("expected outside root error")
	}
}

package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenInFileManager_missingPath(t *testing.T) {
	err := OpenInFileManager(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestOpenInFileManager_existingDir(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping OS file manager launch in CI")
	}
	dir := t.TempDir()
	if err := OpenInFileManager(dir); err != nil {
		t.Fatal(err)
	}
}

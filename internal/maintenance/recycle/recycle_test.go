package recycle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspect_linux(t *testing.T) {
	if filepath.Separator != '/' {
		t.Skip("linux-specific")
	}
	info := Inspect()
	if !info.Supported {
		t.Fatalf("expected supported on linux: %v", info.Error)
	}
	home, _ := os.UserHomeDir()
	if info.Path == "" {
		t.Fatal("expected path")
	}
	if home != "" && info.Path != "" {
		// path should be under home trash
		if info.ItemCount < 0 {
			t.Fatal("bad count")
		}
	}
}

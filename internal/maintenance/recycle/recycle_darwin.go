//go:build darwin

package recycle

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func trashDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".Trash")
}

func inspectPlatform() *BinInfo {
	dir := trashDir()
	if dir == "" {
		return &BinInfo{Supported: false, Error: "home directory not found"}
	}
	count, total, err := walkTrash(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return &BinInfo{Path: dir, Supported: true}
		}
		return &BinInfo{Path: dir, Supported: true, Error: err.Error()}
	}
	return &BinInfo{
		Path:       dir,
		ItemCount:  count,
		TotalBytes: total,
		Supported:  true,
	}
}

func walkTrash(dir string) (count int, total int64, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, err
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		count++
		if info.IsDir() {
			_, sub, _ := walkTrash(filepath.Join(dir, e.Name()))
			total += sub
		} else {
			total += info.Size()
		}
	}
	return count, total, nil
}

func emptyPlatform(ctx context.Context) (*EmptyReport, error) {
	script := `tell application "Finder" to empty trash`
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return &EmptyReport{Error: msg}, err
	}
	after := inspectPlatform()
	return &EmptyReport{
		ItemCount:  after.ItemCount,
		TotalBytes: after.TotalBytes,
	}, nil
}

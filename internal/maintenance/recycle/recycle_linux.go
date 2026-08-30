//go:build linux

package recycle

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func trashDirs() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".local", "share", "Trash", "files"),
		filepath.Join(home, ".Trash", "files"),
	}
}

func inspectPlatform() *BinInfo {
	dirs := trashDirs()
	if len(dirs) == 0 {
		return &BinInfo{Supported: false, Error: "home directory not found"}
	}
	var total int64
	var count int
	var usedPath string
	for _, d := range dirs {
		c, b, err := walkTrash(d)
		if err != nil {
			continue
		}
		if c > 0 || usedPath == "" {
			usedPath = d
		}
		count += c
		total += b
	}
	if usedPath == "" {
		usedPath = dirs[0]
	}
	return &BinInfo{
		Path:       usedPath,
		ItemCount:  count,
		TotalBytes: total,
		Supported:  true,
	}
}

func walkTrash(dir string) (count int, total int64, err error) {
	if _, err = os.Stat(dir); err != nil {
		return 0, 0, err
	}
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		count++
		total += info.Size()
		return nil
	})
	return count, total, nil
}

func emptyPlatform(ctx context.Context) (*EmptyReport, error) {
	if _, err := exec.LookPath("gio"); err == nil {
		cmd := exec.CommandContext(ctx, "gio", "trash", "--empty")
		if out, err := cmd.CombinedOutput(); err != nil {
			return &EmptyReport{Error: strings.TrimSpace(string(out))}, err
		}
		after := inspectPlatform()
		return &EmptyReport{
			ItemCount:  after.ItemCount,
			TotalBytes: after.TotalBytes,
		}, nil
	}
	for _, d := range trashDirs() {
		filesDir := d
		infoDir := strings.TrimSuffix(d, "/files") + "/info"
		_ = os.RemoveAll(filesDir)
		_ = os.RemoveAll(infoDir)
		_ = os.MkdirAll(filesDir, 0o755)
		_ = os.MkdirAll(infoDir, 0o755)
	}
	after := inspectPlatform()
	return &EmptyReport{
		ItemCount:  after.ItemCount,
		TotalBytes: after.TotalBytes,
	}, nil
}

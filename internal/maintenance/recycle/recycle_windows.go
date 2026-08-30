//go:build windows

package recycle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func inspectPlatform() *BinInfo {
	root, err := recycleRoot()
	if err != nil {
		return &BinInfo{Supported: false, Error: err.Error()}
	}
	count, total, err := walkTrash(root)
	if err != nil {
		return &BinInfo{Path: root, Supported: true, Error: err.Error()}
	}
	return &BinInfo{
		Path:       root,
		ItemCount:  count,
		TotalBytes: total,
		Supported:  true,
	}
}

func recycleRoot() (string, error) {
	drive := os.Getenv("SystemDrive")
	if drive == "" {
		drive = "C:"
	}
	return filepath.Join(drive+string(os.PathSeparator), "$Recycle.Bin"), nil
}

func walkTrash(dir string) (count int, total int64, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, err
	}
	for _, sid := range entries {
		if !sid.IsDir() {
			continue
		}
		c, b, _ := walkFiles(filepath.Join(dir, sid.Name()))
		count += c
		total += b
	}
	return count, total, nil
}

func walkFiles(dir string) (count int, total int64, err error) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil || info.IsDir() {
			return nil
		}
		count++
		total += info.Size()
		return nil
	})
	return count, total, nil
}

func emptyPlatform(ctx context.Context) (*EmptyReport, error) {
	ps := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command",
		"Clear-RecycleBin -Force -ErrorAction Stop")
	if out, err := ps.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return &EmptyReport{Error: msg}, fmt.Errorf("%s", msg)
	}
	after := inspectPlatform()
	return &EmptyReport{
		ItemCount:  after.ItemCount,
		TotalBytes: after.TotalBytes,
	}, nil
}

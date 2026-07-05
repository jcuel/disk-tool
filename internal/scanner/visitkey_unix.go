//go:build !windows

package scanner

import (
	"os"
	"syscall"
)

type visitKey struct {
	dev uint64
	ino uint64
}

func visitKeyFromInfo(path string, info os.FileInfo) visitKey {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return visitKey{dev: uint64(stat.Dev), ino: stat.Ino}
	}
	return visitKey{dev: 0, ino: uint64(info.ModTime().UnixNano())}
}

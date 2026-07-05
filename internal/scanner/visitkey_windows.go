//go:build windows

package scanner

import "os"

type visitKey struct {
	dev uint64
	ino uint64
}

func visitKeyFromInfo(path string, info os.FileInfo) visitKey {
	h := uint64(len(path))
	for i := 0; i < len(path); i++ {
		h = h*31 + uint64(path[i])
	}
	return visitKey{dev: 0, ino: h ^ uint64(info.ModTime().UnixNano())}
}

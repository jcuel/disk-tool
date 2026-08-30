//go:build !windows

package wsl

import "context"

// ListDisks returns empty on non-Windows platforms.
func ListDisks() ([]Disk, error) {
	return nil, nil
}

// Compact is unsupported off Windows.
func Compact(ctx context.Context, req CompactRequest) (*CompactReport, error) {
	return &CompactReport{Supported: false, Error: errUnsupported.Error()}, errUnsupported
}

// Supported reports whether WSL compact is available on this OS.
func Supported() bool {
	return false
}

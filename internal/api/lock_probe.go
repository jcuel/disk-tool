package api

import (
	"os"
	"strings"
)

func isLockError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "being used") ||
		strings.Contains(msg, "used by another process") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "sharing violation")
}

// PathInUse returns true when a file appears locked. Directories return false (checked on delete).
func PathInUse(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return isLockError(err)
	}
	_ = f.Close()
	return false
}

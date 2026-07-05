package api

import (
	"errors"
	"os"
)

var errConfirmRequired = errors.New("confirm must be true")

// DeletePath removes a file or directory. Caller must validate path is within scan root.
func DeletePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

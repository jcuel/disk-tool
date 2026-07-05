package api

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenInFileManager opens path in the OS file manager (folder or file).
func OpenInFileManager(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			cmd = exec.Command("cmd", "/c", "explorer", "/select,", path)
		} else {
			cmd = exec.Command("cmd", "/c", "explorer", path)
		}
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

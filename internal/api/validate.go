package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ValidateRoot(root string) (string, error) {
	if root == "" {
		return "", errors.New("root path is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("root must be a directory")
	}
	return abs, nil
}

func PathWithinRoot(root, target string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	rootAbs = filepath.Clean(rootAbs)
	targetAbs = filepath.Clean(targetAbs)

	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", errors.New("path outside scan root")
	}
	return targetAbs, nil
}

func CommonRoots() []string {
	var roots []string
	if home, err := os.UserHomeDir(); err == nil {
		roots = append(roots, home)
	}
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if vol := os.Getenv("SystemDrive"); vol != "" {
		roots = append(roots, vol+string(os.PathSeparator))
	} else {
		roots = append(roots, string(os.PathSeparator))
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(roots))
	for _, r := range roots {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

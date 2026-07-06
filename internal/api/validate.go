package api

import (
	"errors"
	"os"

	"github.com/jcuel/disk-tool/internal/safepath"
)

func ValidateRoot(root string) (string, error) {
	abs, err := safepath.ResolveRoot(root)
	if err != nil {
		if errors.Is(err, safepath.ErrEmptyPath) {
			return "", errors.New("root path is required")
		}
		if errors.Is(err, safepath.ErrNotDir) {
			return "", errors.New("root must be a directory")
		}
		return "", err
	}
	return abs, nil
}

func PathWithinRoot(root, target string) (string, error) {
	abs, err := safepath.UnderRoot(root, target)
	if err != nil {
		if errors.Is(err, safepath.ErrOutsideRoot) {
			return "", errors.New("path outside scan root")
		}
		return "", err
	}
	return abs, nil
}

func CommonRoots() []string {
	var roots []string
	if vol := os.Getenv("SystemDrive"); vol != "" {
		roots = append(roots, vol+string(os.PathSeparator))
	}
	if home, err := os.UserHomeDir(); err == nil {
		roots = append(roots, home)
	}
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if len(roots) == 0 {
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

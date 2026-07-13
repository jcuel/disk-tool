package docker

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DataRoot describes a Docker data directory (report only; never FS-delete).
type DataRoot struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hint string `json:"hint"`
}

// DataRoots returns known Docker data directories that exist on this machine.
func DataRoots() []DataRoot {
	var roots []DataRoot
	for _, p := range candidateDataPaths() {
		info, err := os.Stat(p)
		if err != nil || !info.IsDir() {
			continue
		}
		size := dirSizeLimited(p, 2)
		roots = append(roots, DataRoot{
			Path: p,
			Size: size,
			Hint: "Docker data root — reclaim via docker system prune or Docker Desktop; do not delete this folder",
		})
	}
	return roots
}

// IsDataRoot reports whether path is a Docker synthetic reclaim path or under a known data root.
func IsDataRoot(path string) bool {
	if path == "" {
		return false
	}
	if strings.HasPrefix(path, "docker://") {
		return true
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	abs = filepath.Clean(abs)
	lower := strings.ToLower(abs)
	for _, root := range candidateDataPaths() {
		r, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		r = filepath.Clean(r)
		rl := strings.ToLower(r)
		if lower == rl || strings.HasPrefix(lower, rl+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func candidateDataPaths() []string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		local := os.Getenv("LOCALAPPDATA")
		var out []string
		if local != "" {
			out = append(out, filepath.Join(local, "Docker"))
		}
		if home != "" {
			out = append(out, filepath.Join(home, "AppData", "Local", "Docker"))
		}
		return out
	case "darwin":
		if home == "" {
			return nil
		}
		return []string{
			filepath.Join(home, "Library", "Containers", "com.docker.docker", "Data"),
		}
	default:
		var out []string
		out = append(out, "/var/lib/docker")
		if home != "" {
			out = append(out, filepath.Join(home, ".local", "share", "docker"))
		}
		return out
	}
}

// dirSizeLimited walks at most maxDepth levels to avoid scanning huge trees.
func dirSizeLimited(root string, maxDepth int) int64 {
	var total int64
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if rel != "." {
			depth := strings.Count(rel, string(os.PathSeparator)) + 1
			if d.IsDir() && depth > maxDepth {
				return filepath.SkipDir
			}
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total
}

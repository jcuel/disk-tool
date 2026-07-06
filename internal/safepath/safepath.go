package safepath

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyPath   = errors.New("path is required")
	ErrOutsideRoot = errors.New("path outside scan root")
	ErrNotDir      = errors.New("root must be a directory")
)

func ResolveRoot(userPath string) (string, error) {
	root, abs, err := OpenRoot(userPath)
	if err != nil {
		return "", err
	}
	_ = root.Close()
	return abs, nil
}

func OpenRoot(userPath string) (*os.Root, string, error) {
	if userPath == "" {
		return nil, "", ErrEmptyPath
	}
	if strings.Contains(userPath, "\x00") {
		return nil, "", ErrEmptyPath
	}
	abs, err := filepath.Abs(userPath)
	if err != nil {
		return nil, "", err
	}
	abs = filepath.Clean(abs)
	root, err := os.OpenRoot(abs)
	if err != nil {
		if errors.Is(err, os.ErrInvalid) {
			return nil, "", ErrNotDir
		}
		return nil, "", err
	}
	return root, abs, nil
}

func UnderRoot(root, target string) (string, error) {
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
		return "", ErrOutsideRoot
	}
	return targetAbs, nil
}

func Rel(rootAbs, targetAbs string) (string, error) {
	rootAbs = filepath.Clean(rootAbs)
	targetAbs = filepath.Clean(targetAbs)
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", ErrOutsideRoot
	}
	return rel, nil
}

func AbsUnder(rootAbs, rel string) (string, error) {
	if rel == "." {
		return filepath.Clean(rootAbs), nil
	}
	targetAbs := filepath.Join(rootAbs, rel)
	return UnderRoot(rootAbs, targetAbs)
}

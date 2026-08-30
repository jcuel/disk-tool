//go:build windows

package wsl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	errConfirmRequired = errors.New("confirm required")
	errConfirmPhrase   = errors.New("confirmPhrase must be DELETE")
	errPathRequired    = errors.New("path required")
)

// Supported reports whether WSL compact is available on this OS.
func Supported() bool {
	return true
}

// ListDisks finds Docker Desktop and WSL2 VHDX files.
func ListDisks() ([]Disk, error) {
	var out []Disk
	seen := map[string]struct{}{}
	add := func(path, kind string) {
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			return
		}
		info, err := os.Stat(path)
		if err != nil {
			return
		}
		seen[path] = struct{}{}
		out = append(out, Disk{Path: path, SizeBytes: info.Size(), Kind: kind})
	}

	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		dockerWSL := filepath.Join(local, "Docker", "wsl")
		_ = filepath.Walk(dockerWSL, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			if strings.EqualFold(filepath.Ext(path), ".vhdx") {
				add(path, "docker-wsl")
			}
			return nil
		})
	}

	localPackages := filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages")
	_ = filepath.Walk(localPackages, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Base(path), "ext4.vhdx") {
			add(path, "wsl-distro")
		}
		return nil
	})

	return out, nil
}

// Compact shuts down WSL and compacts the given VHDX.
func Compact(ctx context.Context, req CompactRequest) (*CompactReport, error) {
	if req.Path == "" {
		return nil, errPathRequired
	}
	abs, err := filepath.Abs(req.Path)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(filepath.Ext(abs), ".vhdx") {
		return &CompactReport{Supported: true, Error: "not a .vhdx file"}, errors.New("not a .vhdx file")
	}
	before, err := fileSize(abs)
	if err != nil {
		return &CompactReport{Supported: true, Path: abs, Error: err.Error()}, err
	}
	if req.DryRun {
		return &CompactReport{
			DryRun:      true,
			Path:        abs,
			BytesBefore: before,
			BytesAfter:  before,
			FreedBytes:  0,
			Supported:   true,
		}, nil
	}
	if !req.Confirm {
		return nil, errConfirmRequired
	}
	if req.ConfirmPhrase != ConfirmPhrase {
		return nil, errConfirmPhrase
	}

	var log strings.Builder
	shutdown := exec.CommandContext(ctx, "wsl", "--shutdown")
	if out, err := shutdown.CombinedOutput(); err != nil {
		log.WriteString("wsl --shutdown: " + strings.TrimSpace(string(out)) + "\n")
	}
	time.Sleep(2 * time.Second)

	compactOut, compactErr := compactVHDX(ctx, abs)
	log.WriteString(compactOut)
	if compactErr != nil {
		return &CompactReport{
			Path:        abs,
			BytesBefore: before,
			Output:      log.String(),
			Error:       compactErr.Error(),
			Supported:   true,
		}, compactErr
	}

	after, _ := fileSize(abs)
	freed := before - after
	if freed < 0 {
		freed = 0
	}
	return &CompactReport{
		Path:        abs,
		BytesBefore: before,
		BytesAfter:  after,
		FreedBytes:  freed,
		Output:      log.String(),
		Supported:   true,
	}, nil
}

func fileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func compactVHDX(ctx context.Context, path string) (string, error) {
	psScript := fmt.Sprintf(`Optimize-VHD -Path '%s' -Mode Full`, strings.ReplaceAll(path, "'", "''"))
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err == nil {
		return text, nil
	}
	// diskpart fallback
	script := fmt.Sprintf(`select vdisk file="%s"
attach vdisk readonly
compact vdisk
detach vdisk
exit
`, strings.ReplaceAll(path, `"`, `\"`))
	dp := exec.CommandContext(ctx, "diskpart")
	dp.Stdin = strings.NewReader(script)
	dpOut, dpErr := dp.CombinedOutput()
	combined := text + "\n" + strings.TrimSpace(string(dpOut))
	if dpErr != nil {
		return combined, fmt.Errorf("compact failed (Optimize-VHD and diskpart): %v", dpErr)
	}
	return combined, nil
}

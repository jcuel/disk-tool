package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SyntheticPath is the insights path for Docker CLI reclaim (not a filesystem path).
const SyntheticPath = "docker://reclaimable"

// DiskUsage summarizes docker system df.
type DiskUsage struct {
	Available         bool   `json:"available"`
	DaemonOK          bool   `json:"daemonOk"`
	Error             string `json:"error,omitempty"`
	ImagesSize        int64  `json:"imagesSize"`
	ImagesReclaim     int64  `json:"imagesReclaimable"`
	ContainersSize    int64  `json:"containersSize"`
	ContainersReclaim int64  `json:"containersReclaimable"`
	VolumesSize       int64  `json:"volumesSize"`
	VolumesReclaim    int64  `json:"volumesReclaimable"`
	BuildCacheSize    int64  `json:"buildCacheSize"`
	BuildCacheReclaim int64  `json:"buildCacheReclaimable"`
	// Reclaimable is images+containers+build cache (volumes excluded from prune).
	Reclaimable int64  `json:"reclaimable"`
	RawDF       string `json:"rawDf,omitempty"`
}

// PruneReport is the result of a prune dry-run or execute.
type PruneReport struct {
	DryRun      bool   `json:"dryRun"`
	Reclaimable int64  `json:"reclaimable"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Available reports whether the docker binary is on PATH.
func Available() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// Detect returns Docker disk usage. If CLI missing, Available=false.
func Detect(ctx context.Context) *DiskUsage {
	if ctx == nil {
		ctx = context.Background()
	}
	u := &DiskUsage{}
	if !Available() {
		u.Error = "docker CLI not found on PATH"
		return u
	}
	u.Available = true
	out, err := runDocker(ctx, "system", "df")
	if err != nil {
		u.Error = err.Error()
		return u
	}
	u.DaemonOK = true
	u.RawDF = out
	parseSystemDF(out, u)
	u.Reclaimable = u.ImagesReclaim + u.ContainersReclaim + u.BuildCacheReclaim
	return u
}

func runDocker(ctx context.Context, args ...string) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cctx, "docker", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("docker %s: %s", strings.Join(args, " "), msg)
	}
	return stdout.String(), nil
}

// sizeRe matches docker df size cells like "12.3GB" or "1.2MB (60%)".
var (
	typeLine = regexp.MustCompile(`(?i)^(Images|Containers|Local Volumes|Build Cache)\s+`)
	sizeTok  = regexp.MustCompile(`(?i)([0-9]+(?:\.[0-9]+)?)\s*(B|KB|MB|GB|TB)\b`)
)

func parseSystemDF(out string, u *DiskUsage) {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToLower(line), "type") {
			continue
		}
		m := typeLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		kind := strings.ToLower(m[1])
		sizes := sizeTok.FindAllStringSubmatch(line, -1)
		if len(sizes) == 0 {
			continue
		}
		total := parseSizeToken(sizes[0])
		reclaim := int64(0)
		if len(sizes) >= 2 {
			reclaim = parseSizeToken(sizes[len(sizes)-1])
		}
		switch {
		case strings.HasPrefix(kind, "image"):
			u.ImagesSize, u.ImagesReclaim = total, reclaim
		case strings.HasPrefix(kind, "container"):
			u.ContainersSize, u.ContainersReclaim = total, reclaim
		case strings.Contains(kind, "volume"):
			u.VolumesSize, u.VolumesReclaim = total, reclaim
		case strings.Contains(kind, "build"):
			u.BuildCacheSize, u.BuildCacheReclaim = total, reclaim
		}
	}
}

func parseSizeToken(m []string) int64 {
	if len(m) < 3 {
		return 0
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0
	}
	switch strings.ToUpper(m[2]) {
	case "B":
		return int64(v)
	case "KB":
		return int64(v * 1024)
	case "MB":
		return int64(v * 1024 * 1024)
	case "GB":
		return int64(v * 1024 * 1024 * 1024)
	case "TB":
		return int64(v * 1024 * 1024 * 1024 * 1024)
	default:
		return 0
	}
}

// PruneDryRun reports reclaimable bytes without deleting (uses system df).
func PruneDryRun(ctx context.Context) (*PruneReport, error) {
	u := Detect(ctx)
	if !u.Available {
		return &PruneReport{DryRun: true, Error: u.Error}, errors.New(u.Error)
	}
	if !u.DaemonOK {
		return &PruneReport{DryRun: true, Error: u.Error}, errors.New(u.Error)
	}
	return &PruneReport{
		DryRun:      true,
		Reclaimable: u.Reclaimable,
		Output:      u.RawDF,
	}, nil
}

// Prune runs docker system prune -af without --volumes.
func Prune(ctx context.Context) (*PruneReport, error) {
	if !Available() {
		return &PruneReport{Error: "docker CLI not found on PATH"}, errors.New("docker CLI not found on PATH")
	}
	before := Detect(ctx)
	out, err := runDocker(ctx, "system", "prune", "-af")
	if err != nil {
		return &PruneReport{Output: out, Error: err.Error()}, err
	}
	reclaimed := before.Reclaimable
	after := Detect(ctx)
	if after.DaemonOK && before.DaemonOK {
		diff := before.Reclaimable - after.Reclaimable
		if diff > 0 {
			reclaimed = diff
		}
	}
	return &PruneReport{
		DryRun:      false,
		Reclaimable: reclaimed,
		Output:      out,
	}, nil
}

package docker

import "testing"

func TestParseSystemDF(t *testing.T) {
	sample := `TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE
Images          12        3         8.5GB     6.2GB
Containers      5         1         120.5MB   80.1MB
Local Volumes   4         2         1.1GB     400MB
Build Cache     20        0         2.3GB     2.3GB
`
	u := &DiskUsage{}
	parseSystemDF(sample, u)
	if u.ImagesSize == 0 || u.ImagesReclaim == 0 {
		t.Fatalf("images: size=%d reclaim=%d", u.ImagesSize, u.ImagesReclaim)
	}
	if u.BuildCacheReclaim == 0 {
		t.Fatal("expected build cache reclaim")
	}
	if u.VolumesReclaim == 0 {
		t.Fatal("expected volumes reclaim parsed (even if not pruned)")
	}
}

func TestParseSizeToken(t *testing.T) {
	if parseSizeToken([]string{"", "1.5", "GB"}) != int64(1.5*1024*1024*1024) {
		t.Fatal("1.5GB parse")
	}
	if parseSizeToken([]string{"", "512", "MB"}) != 512*1024*1024 {
		t.Fatal("512MB parse")
	}
}

func TestIsDataRoot_synthetic(t *testing.T) {
	if !IsDataRoot(SyntheticPath) {
		t.Fatal("expected synthetic docker:// path to be protected")
	}
}

func TestCandidateDataPaths_nonEmpty(t *testing.T) {
	paths := candidateDataPaths()
	if len(paths) == 0 {
		t.Fatal("expected at least one candidate path for this OS")
	}
}

package diskspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVolumeRoot_windowsStyle(t *testing.T) {
	if filepath.VolumeName(`C:\Users\foo`) == "" {
		t.Skip("no drive letters on this platform")
	}
	got := VolumeRoot(`C:\Users\foo`)
	if got != `C:\` && got != `C:`+string(os.PathSeparator) {
		t.Fatalf("unexpected root: %q", got)
	}
}

func TestForPath(t *testing.T) {
	root := VolumeRoot(os.TempDir())
	info, err := ForPath(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if info.Path != root {
		t.Fatalf("path %q want %q", info.Path, root)
	}
	if info.Total <= 0 || info.Free < 0 || info.Used < 0 {
		t.Fatalf("invalid sizes: %+v", info)
	}
	if info.Used+info.Free != info.Total {
		t.Fatalf("used+free != total: %+v", info)
	}
}

package safety

import (
	"runtime"
	"testing"
)

func TestClassifyPath_forbiddenWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows paths")
	}
	cases := map[string]Zone{
		`C:\Windows\System32`:         ZoneForbidden,
		`C:\Windows\SysWOW64\drivers`: ZoneForbidden,
		`C:\pagefile.sys`:             ZoneForbidden,
	}
	for path, want := range cases {
		if got := ClassifyPath(path); got != want {
			t.Fatalf("%s: got %q want %q", path, got, want)
		}
	}
}

func TestClassifyPath_unixForbidden(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix paths")
	}
	cases := map[string]Zone{
		"/proc/1/status": ZoneForbidden,
		"/sys/class":     ZoneForbidden,
		"/dev/sda":       ZoneForbidden,
		"/usr/bin":       ZoneCriticalOS,
		"/var/crash/pkg": ZoneDiagnostic,
		"/tmp/cache":     ZoneMaintenance,
	}
	for path, want := range cases {
		if got := ClassifyPath(path); got != want {
			t.Fatalf("%s: got %q want %q", path, got, want)
		}
	}
}

func TestIsDiskImagePath(t *testing.T) {
	cases := map[string]bool{
		`C:\Users\x\Hyper-V\VM.vhdx`: true,
		`/home/u/disk.vhd`:           true,
		`C:\Users\x\vm.vmdk`:         true,
		`C:\Users\x\box.vdi`:         true,
		`/var/lib/libvirt/img.qcow2`: true,
		`C:\Users\x\install.iso`:     false,
		`C:\Users\x\big.zip`:         false,
	}
	for path, want := range cases {
		if got := IsDiskImagePath(path); got != want {
			t.Fatalf("%s: got %v want %v", path, got, want)
		}
	}
}

func TestCanDeletePath_diskImage(t *testing.T) {
	ok, reason := CanDeletePath(`C:\Users\x\Disks\dev.vhdx`)
	if ok {
		t.Fatal("vhdx must not be deletable")
	}
	if reason == "" {
		t.Fatal("expected reason")
	}
}

func TestShouldSkipScan(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix virtual fs")
	}
	if !ShouldSkipScan("/proc/cpuinfo") {
		t.Fatal("expected skip for /proc")
	}
	if ShouldSkipScan("/home/user/projects") {
		t.Fatal("user path should not skip")
	}
}

func TestApplyCandidateZone(t *testing.T) {
	if runtime.GOOS == "windows" {
		zone, ok := ApplyCandidateZone(`C:\Users\x\AppData\Local\Temp\old`, "review")
		if zone != ZoneMaintenance || !ok {
			t.Fatalf("temp: zone=%s deletable=%v", zone, ok)
		}
		return
	}
	zone, ok := ApplyCandidateZone("/tmp/old", "review")
	if zone != ZoneMaintenance || !ok {
		t.Fatalf("tmp: zone=%s deletable=%v", zone, ok)
	}
	zone, ok = ApplyCandidateZone("/var/crash/core", "review")
	if zone != ZoneDiagnostic || ok {
		t.Fatalf("crash: zone=%s deletable=%v", zone, ok)
	}
}

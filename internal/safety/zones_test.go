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

func TestIsDeletable(t *testing.T) {
	if !IsDeletable(ZoneReview) {
		t.Fatal("review should be deletable")
	}
	if IsDeletable(ZoneForbidden) || IsDeletable(ZoneCriticalOS) || IsDeletable(ZoneDiagnostic) {
		t.Fatal("protected zones must not be deletable")
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

package docker

import "testing"

func TestComputeReclaimed_positiveDiff(t *testing.T) {
	got, noChange := computeReclaimed(10*1024*1024, 2*1024*1024)
	if noChange {
		t.Fatal("expected change")
	}
	if got != 8*1024*1024 {
		t.Fatalf("got %d", got)
	}
}

func TestComputeReclaimed_noChange(t *testing.T) {
	got, noChange := computeReclaimed(5*1024*1024, 5*1024*1024)
	if !noChange {
		t.Fatal("expected noChange")
	}
	if got != 0 {
		t.Fatalf("got %d want 0", got)
	}
}

func TestComputeReclaimed_afterLarger(t *testing.T) {
	got, noChange := computeReclaimed(1*1024*1024, 2*1024*1024)
	if !noChange {
		t.Fatal("expected noChange when after >= before")
	}
	if got != 0 {
		t.Fatalf("got %d want 0", got)
	}
}

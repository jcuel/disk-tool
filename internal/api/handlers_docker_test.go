package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jcuel/disk-tool/internal/model"
)

func putCompletedJob(t *testing.T, store *Store, root string) *model.ScanJob {
	t.Helper()
	id := newScanID()
	job := &model.ScanJob{
		ID:     id,
		Root:   root,
		Status: model.ScanStatusCompleted,
	}
	store.mu.Lock()
	store.jobs[id] = job
	store.subscribers[id] = make(map[subscriber]struct{})
	store.mu.Unlock()
	return job
}

func TestDockerPrune_requiresConfirm(t *testing.T) {
	store := NewStore()
	job := putCompletedJob(t, store, t.TempDir())
	s := NewServer(store, nil)
	body := `{"dryRun":false,"confirm":false,"confirmPhrase":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/scans/"+job.ID+"/docker/prune", strings.NewReader(body))
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
}

func TestDockerPrune_dryRunOK(t *testing.T) {
	store := NewStore()
	job := putCompletedJob(t, store, t.TempDir())
	s := NewServer(store, nil)
	body := `{"dryRun":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/scans/"+job.ID+"/docker/prune", strings.NewReader(body))
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var report map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&report); err != nil {
		t.Fatal(err)
	}
	if report["dryRun"] != true {
		t.Fatalf("expected dryRun true: %#v", report)
	}
}

func TestDockerStatus_notFound(t *testing.T) {
	s := NewServer(NewStore(), nil)
	req := httptest.NewRequest(http.MethodGet, "/api/scans/missing/docker", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status %d", rr.Code)
	}
}

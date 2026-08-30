package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRecycleInspect(t *testing.T) {
	s := NewServer(NewStore(), nil)
	req := httptest.NewRequest(http.MethodGet, "/api/maintenance/recycle", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var info struct {
		Supported bool `json:"supported"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &info); err != nil {
		t.Fatal(err)
	}
	if !info.Supported {
		t.Fatal("expected supported on linux")
	}
}

func TestHandleWSLDisks(t *testing.T) {
	s := NewServer(NewStore(), nil)
	req := httptest.NewRequest(http.MethodGet, "/api/wsl/disks", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Supported bool `json:"supported"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
}

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jcuel/disk-tool/internal/maintenance/recycle"
	"github.com/jcuel/disk-tool/internal/maintenance/wsl"
)

func (s *Server) registerMaintenanceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/maintenance/recycle", s.handleRecycleInspect)
	mux.HandleFunc("POST /api/maintenance/recycle/empty", s.handleRecycleEmpty)
	mux.HandleFunc("GET /api/wsl/disks", s.handleWSLDisks)
	mux.HandleFunc("POST /api/wsl/compact", s.handleWSLCompact)
}

func (s *Server) handleRecycleInspect(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, recycle.Inspect())
}

func (s *Server) handleRecycleEmpty(w http.ResponseWriter, r *http.Request) {
	var req recycle.EmptyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	report, err := recycle.Empty(ctx, req)
	if err != nil {
		if report != nil && report.Error != "" {
			writeJSON(w, http.StatusBadRequest, report)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleWSLDisks(w http.ResponseWriter, _ *http.Request) {
	disks, err := wsl.ListDisks()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"supported": wsl.Supported(),
		"disks":     disks,
	})
}

func (s *Server) handleWSLCompact(w http.ResponseWriter, r *http.Request) {
	var req wsl.CompactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	report, err := wsl.Compact(ctx, req)
	if err != nil {
		if report != nil && report.Error != "" {
			writeJSON(w, http.StatusBadRequest, report)
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

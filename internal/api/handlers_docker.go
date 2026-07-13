package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jcuel/disk-tool/internal/docker"
	"github.com/jcuel/disk-tool/internal/model"
)

func (s *Server) registerDockerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/scans/{id}/docker", s.handleDockerStatus)
	mux.HandleFunc("POST /api/scans/{id}/docker/prune", s.handleDockerPrune)
}

func (s *Server) handleDockerStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, ok := s.store.Get(id); !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	u := docker.Detect(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{
		"usage":     u,
		"dataRoots": docker.DataRoots(),
	})
}

type dockerPruneRequest struct {
	DryRun        bool   `json:"dryRun"`
	Confirm       bool   `json:"confirm"`
	ConfirmPhrase string `json:"confirmPhrase"`
}

func (s *Server) handleDockerPrune(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	if job.Status != model.ScanStatusCompleted {
		writeError(w, http.StatusBadRequest, "scan not completed")
		return
	}
	var req dockerPruneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if req.DryRun {
		report, err := docker.PruneDryRun(ctx)
		if report == nil {
			report = &docker.PruneReport{DryRun: true}
		}
		if err != nil && report.Error == "" {
			report.Error = err.Error()
		}
		writeJSON(w, http.StatusOK, report)
		return
	}
	if !req.Confirm {
		writeError(w, http.StatusBadRequest, errConfirmRequired.Error())
		return
	}
	if req.ConfirmPhrase != cleanupConfirmPhrase {
		writeError(w, http.StatusBadRequest, errCleanupConfirmPhrase.Error())
		return
	}
	report, err := docker.Prune(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

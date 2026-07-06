package api

import (
	"encoding/json"
	"net/http"

	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/safety"
)

func (s *Server) registerSafetyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/scans/{id}/maintenance-presets", s.handleMaintenancePresets)
	mux.HandleFunc("POST /api/scans/{id}/reanalyze", s.handleReanalyze)
	mux.HandleFunc("POST /api/scans/{id}/duplicates", s.handleFindDuplicates)
}

func (s *Server) handleMaintenancePresets(w http.ResponseWriter, r *http.Request) {
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
	candidates := []model.CleanupCandidate{}
	if job.Insights != nil {
		candidates = job.Insights.CleanupCandidates
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"presets": safety.AllPresets(),
		"matches": safety.MatchPresets(candidates),
	})
}

func (s *Server) handleReanalyze(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cfg model.InsightsConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := s.store.Reanalyze(id, cfg); err != nil {
		switch err {
		case errNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		case errNotReady:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	job, _ := s.store.Get(id)
	writeJSON(w, http.StatusOK, job.Insights)
}

func (s *Server) handleFindDuplicates(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.FindDuplicates(id); err != nil {
		switch err {
		case errNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		case errNotReady:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	job, _ := s.store.Get(id)
	writeJSON(w, http.StatusOK, map[string]any{"duplicateGroups": job.DuplicateGroups})
}

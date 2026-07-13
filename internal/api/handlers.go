package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jcuel/disk-tool/internal/diskspace"
	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/safety"
	"github.com/jcuel/disk-tool/internal/scanner"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	store   *Store
	static  http.Handler
}

func NewServer(store *Store, static http.Handler) *Server {
	return &Server{store: store, static: static}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/roots", s.handleRoots)
	mux.HandleFunc("GET /api/disk", s.handleDisk)
	mux.HandleFunc("POST /api/scans", s.handleStartScan)
	mux.HandleFunc("GET /api/scans/{id}", s.handleGetScan)
	mux.HandleFunc("DELETE /api/scans/{id}", s.handleCancelScan)
	mux.HandleFunc("POST /api/scans/{id}/expand", s.handleExpandScan)
	mux.HandleFunc("GET /api/scans/{id}/events", s.handleScanEvents)
	mux.HandleFunc("GET /api/scans/{id}/export", s.handleExport)
	mux.HandleFunc("POST /api/scans/{id}/open", s.handleOpenPath)
	mux.HandleFunc("POST /api/scans/{id}/delete", s.handleDeletePath)
	mux.HandleFunc("POST /api/scans/{id}/cleanup", s.handleCleanup)
	s.registerSafetyRoutes(mux)
	if s.static != nil {
		mux.Handle("/", s.static)
	}
	return mux
}

func (s *Server) handleRoots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"roots": CommonRoots()})
}

func (s *Server) handleDisk(w http.ResponseWriter, r *http.Request) {
	path := SanitizePath(r.URL.Query().Get("path"))
	if path == "" {
		path = string(os.PathSeparator)
		if vol := os.Getenv("SystemDrive"); vol != "" {
			path = vol + string(os.PathSeparator)
		}
	}
	info, err := diskspace.ForPath(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (s *Server) handleStartScan(w http.ResponseWriter, r *http.Request) {
	var req model.StartScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	job, err := s.store.Start(req.Root, insightsConfigFromRequest(req))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, model.ScanResponse{ScanID: job.ID})
}

func (s *Server) handleGetScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.store.GetMutable(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	annotateLargestFiles(job)
	writeJSON(w, http.StatusOK, job)
}

func annotateLargestFiles(job *model.ScanJob) {
	if job == nil {
		return
	}
	for i := range job.LargestFiles {
		ok, _ := safety.CanDeletePath(job.LargestFiles[i].Path)
		d := ok
		job.LargestFiles[i].Deletable = &d
	}
}

func (s *Server) handleExpandScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req model.ExpandScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}
	if err := s.store.Expand(id, req.Path, req.Depth); err != nil {
		switch err {
		case errNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		case errNotReady, errExpandBusy:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "expanding", "path": req.Path})
}

func (s *Server) handleCancelScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !s.store.Cancel(id) {
		writeError(w, http.StatusNotFound, "scan not found or not running")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelling"})
}

func (s *Server) handleScanEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ch, ok := s.store.Subscribe(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	defer s.store.Unsubscribe(id, ch)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	job, _ := s.store.Get(id)
	if job != nil {
		_ = conn.WriteJSON(model.ProgressEvent{
			Type:         "snapshot",
			ScanID:       id,
			Status:       string(job.Status),
			Root:         job.Root,
			DirsScanned:  job.DirsScanned,
			FilesScanned: job.FilesScanned,
			BytesScanned: job.BytesScanned,
			CurrentPath:  job.CurrentPath,
		})
	}

	for {
		select {
		case ev, open := <-ch:
			if !open {
				return
			}
			if err := conn.WriteJSON(ev); err != nil {
				return
			}
			if ev.Type == "completed" || ev.Type == "cancelled" || ev.Type == "error" {
				// keep connection open for expand events
				if ev.Type != "completed" {
					return
				}
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) handleOpenPath(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Path = SanitizePath(req.Path)
	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}
	abs, err := PathWithinRoot(job.Root, req.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := OpenInFileManager(abs); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "opened", "path": abs})
}

func (s *Server) handleDeletePath(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	var req struct {
		Path          string `json:"path"`
		Confirm       bool   `json:"confirm"`
		ConfirmPhrase string `json:"confirmPhrase"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
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
	req.Path = SanitizePath(req.Path)
	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}
	abs, err := PathWithinRoot(job.Root, req.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if abs == job.Root {
		writeError(w, http.StatusBadRequest, "cannot delete scan root")
		return
	}
	if ok, reason := safety.CanDeletePath(abs); !ok {
		writeError(w, http.StatusBadRequest, reason)
		return
	}
	if err := DeletePath(abs); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "path": abs})
}

func insightsConfigFromRequest(req model.StartScanRequest) model.InsightsConfig {
	if req.InsightsConfig != nil {
		return *req.InsightsConfig
	}
	return model.InsightsConfig{}
}

func (s *Server) handleCleanup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.store.GetMutable(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	if job.Status != model.ScanStatusCompleted {
		writeError(w, http.StatusBadRequest, "scan not completed")
		return
	}
	var req model.CleanupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	report, err := RunCleanup(job, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	job.LastCleanupReport = report
	if !req.DryRun {
		pruneDeletedCandidates(job, report)
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	job, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "scan not found")
		return
	}
	if job.Status != model.ScanStatusCompleted {
		writeError(w, http.StatusBadRequest, "scan not completed")
		return
	}
	switch format {
	case "json":
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"scan-%s.json\"", id))
		writeJSON(w, http.StatusOK, job)
	case "html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"scan-%s.html\"", id))
		if err := renderHTMLReport(w, job); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
	case "ticket":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"disk-report-%s.txt\"", id))
		if job.Insights == nil {
			writeError(w, http.StatusBadRequest, "no insights available")
			return
		}
		_, _ = w.Write([]byte(job.Insights.TicketText))
	case "cleanup-json":
		if job.LastCleanupReport == nil {
			writeError(w, http.StatusNotFound, "no cleanup report available")
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"cleanup-%s.json\"", id))
		writeJSON(w, http.StatusOK, job.LastCleanupReport)
	case "cleanup-html":
		if job.LastCleanupReport == nil {
			writeError(w, http.StatusNotFound, "no cleanup report available")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"cleanup-%s.html\"", id))
		if err := renderCleanupHTMLReport(w, job.LastCleanupReport, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
	case "cleanup-ticket":
		if job.LastCleanupReport == nil {
			writeError(w, http.StatusNotFound, "no cleanup report available")
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"cleanup-%s.txt\"", id))
		_, _ = w.Write([]byte(job.LastCleanupReport.ReportText))
	default:
		writeError(w, http.StatusBadRequest, "format must be json, html, ticket, cleanup-json, cleanup-html, or cleanup-ticket")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func renderHTMLReport(w http.ResponseWriter, job *model.ScanJob) error {
	const tpl = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>disk-tool scan {{.ID}}</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#0f1419;color:#e6edf3}
table{border-collapse:collapse;width:100%}th,td{border:1px solid #30363d;padding:.5rem;text-align:left}
th{background:#161b22}.bar{height:8px;background:#238636;display:inline-block;min-width:2px}
</style></head><body>
<h1>Disk usage report</h1>
<p>{{.Summary}}</p>
<p>Root: {{.Root}} | Indexed files: {{.FilesScanned}} | Bytes scanned: {{.BytesScanned}}</p>
<h2>Top space consumers</h2>
<table><tr><th>Name</th><th>Size</th><th>%</th></tr>
{{range .Rows}}<tr><td>{{.Name}}</td><td>{{.SizeHuman}}</td><td>{{.Pct}}%</td></tr>
{{end}}</table>
{{if .Cleanup}}
<h2>Cleanup candidates</h2>
<table><tr><th>Category</th><th>Path</th><th>Size</th><th>Hint</th></tr>
{{range .Cleanup}}<tr><td>{{.Category}}</td><td>{{.Path}}</td><td>{{.SizeHuman}}</td><td>{{.Hint}}</td></tr>
{{end}}
<p>Potential reclaimable if reviewed: {{.Reclaimable}}</p>
{{end}}
<h2>Largest files</h2>
<table><tr><th>Path</th><th>Size</th></tr>
{{range .Largest}}<tr><td>{{.Path}}</td><td>{{.SizeHuman}}</td></tr>{{end}}
</table></body></html>`

	type row struct {
		Name       string
		SizeHuman  string
		FileCount  int64
		Pct        string
	}
	type largestRow struct {
		Path      string
		SizeHuman string
	}
	type cleanupRow struct {
		Category  string
		Path      string
		SizeHuman string
		Hint      string
	}
	data := struct {
		ID           string
		Root         string
		Summary      string
		FilesScanned int64
		BytesScanned int64
		Rows         []row
		Largest      []largestRow
		Cleanup      []cleanupRow
		Reclaimable  string
	}{
		ID:           job.ID,
		Root:         job.Root,
		FilesScanned: job.FilesScanned,
		BytesScanned: job.BytesScanned,
	}
	if job.Insights != nil {
		data.Summary = job.Insights.Summary
		data.Reclaimable = formatBytes(job.Insights.TotalReclaimable)
		for _, c := range job.Insights.CleanupCandidates {
			data.Cleanup = append(data.Cleanup, cleanupRow{
				Category:  string(c.Category),
				Path:      c.Path,
				SizeHuman: formatBytes(c.Size),
				Hint:      c.Hint,
			})
		}
	}
	rootSize := int64(1)
	if job.Tree != nil && job.Tree.Size > 0 {
		rootSize = job.Tree.Size
	}
	if job.Tree != nil {
		for _, c := range job.Tree.Children {
			data.Rows = append(data.Rows, row{
				Name:      c.Name,
				SizeHuman: formatBytes(c.Size),
				FileCount: c.FileCount,
				Pct:       fmt.Sprintf("%.1f", float64(c.Size)*100/float64(rootSize)),
			})
		}
	}
	if job.Insights != nil && len(job.Insights.TopConsumers) > 0 {
		data.Rows = nil
		for _, c := range job.Insights.TopConsumers {
			data.Rows = append(data.Rows, row{
				Name:      c.Name,
				SizeHuman: formatBytes(c.Size),
				Pct:       fmt.Sprintf("%.1f", c.Pct),
			})
		}
	}
	for _, f := range job.LargestFiles {
		data.Largest = append(data.Largest, largestRow{Path: f.Path, SizeHuman: formatBytes(f.Size)})
	}
	t, err := template.New("report").Parse(tpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func renderCleanupHTMLReport(w http.ResponseWriter, report *model.CleanupReport, scanID string) error {
	const tpl = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>disk-tool cleanup {{.ScanID}}</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#0f1419;color:#e6edf3}
table{border-collapse:collapse;width:100%}th,td{border:1px solid #30363d;padding:.5rem;text-align:left}
th{background:#161b22}.deleted{color:#3fb950}.skipped{color:#d29922}.failed{color:#f85149}
</style></head><body>
<h1>Cleanup report</h1>
<p>Scan: {{.ScanID}} | Mode: {{.Mode}} | Reclaimed: {{.Reclaimed}}</p>
<p>Started: {{.Started}} | Finished: {{.Finished}}</p>
<table><tr><th>Status</th><th>Path</th><th>Size</th><th>Category</th><th>Reason</th></tr>
{{range .Rows}}<tr class="{{.Class}}"><td>{{.Status}}</td><td>{{.Path}}</td><td>{{.Size}}</td><td>{{.Category}}</td><td>{{.Reason}}</td></tr>
{{end}}</table></body></html>`

	type row struct {
		Status   string
		Path     string
		Size     string
		Category string
		Reason   string
		Class    string
	}
	mode := "executed"
	if report.DryRun {
		mode = "dry run"
	}
	data := struct {
		ScanID    string
		Mode      string
		Reclaimed string
		Started   string
		Finished  string
		Rows      []row
	}{
		ScanID:    scanID,
		Mode:      mode,
		Reclaimed: formatBytes(report.BytesReclaimed),
		Started:   report.StartedAt.Format(time.RFC3339),
		Finished:  report.FinishedAt.Format(time.RFC3339),
	}
	for _, r := range report.Results {
		class := "skipped"
		switch r.Status {
		case model.CleanupStatusDeleted, model.CleanupStatusWouldDelete:
			class = "deleted"
		case model.CleanupStatusFailed:
			class = "failed"
		}
		data.Rows = append(data.Rows, row{
			Status:   r.Status,
			Path:     r.Path,
			Size:     formatBytes(r.Size),
			Category: r.Category,
			Reason:   r.Reason,
			Class:    class,
		})
	}
	t, err := template.New("cleanup").Parse(tpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func FlattenTree(node *model.ScanNode, path string) *model.ScanNode {
	if node == nil {
		return nil
	}
	if path == "" || node.Path == path {
		return node
	}
	return scanner.FindNode(node, path)
}

func NodeChildrenJSON(node *model.ScanNode) []map[string]any {
	if node == nil {
		return nil
	}
	out := make([]map[string]any, 0, len(node.Children))
	parentSize := node.Size
	if parentSize == 0 {
		parentSize = 1
	}
	for _, c := range node.Children {
		out = append(out, map[string]any{
			"name":      c.Name,
			"path":      c.Path,
			"size":      c.Size,
			"fileCount": c.FileCount,
			"pct":       float64(c.Size) * 100 / float64(parentSize),
		})
	}
	return out
}

func SanitizePath(p string) string {
	return strings.TrimSpace(p)
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jcuel/disk-tool/internal/model"
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
	mux.HandleFunc("POST /api/scans", s.handleStartScan)
	mux.HandleFunc("GET /api/scans/{id}", s.handleGetScan)
	mux.HandleFunc("DELETE /api/scans/{id}", s.handleCancelScan)
	mux.HandleFunc("POST /api/scans/{id}/expand", s.handleExpandScan)
	mux.HandleFunc("GET /api/scans/{id}/events", s.handleScanEvents)
	mux.HandleFunc("GET /api/scans/{id}/export", s.handleExport)
	if s.static != nil {
		mux.Handle("/", s.static)
	}
	return mux
}

func (s *Server) handleRoots(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"roots": CommonRoots()})
}

func (s *Server) handleStartScan(w http.ResponseWriter, r *http.Request) {
	var req model.StartScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	job, err := s.store.Start(req.Root)
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
	writeJSON(w, http.StatusOK, job)
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
	default:
		writeError(w, http.StatusBadRequest, "format must be json, html, or ticket")
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

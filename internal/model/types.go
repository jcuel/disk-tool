package model

import "time"

const DefaultDrillDepth = 5

type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusCancelled ScanStatus = "cancelled"
	ScanStatusFailed    ScanStatus = "failed"
)

type FileEntry struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ScanNode struct {
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Size       int64       `json:"size"`
	FileCount  int64       `json:"fileCount"`
	IsDir      bool        `json:"isDir"`
	Scanned    bool        `json:"scanned"`
	Expandable bool        `json:"expandable"`
	Children   []*ScanNode `json:"children,omitempty"`
}

type ProgressEvent struct {
	Type         string `json:"type"`
	ScanID       string `json:"scanId"`
	Status       string `json:"status,omitempty"`
	Root         string `json:"root,omitempty"`
	TargetPath   string `json:"targetPath,omitempty"`
	DirsScanned  int64  `json:"dirsScanned,omitempty"`
	FilesScanned int64  `json:"filesScanned,omitempty"`
	BytesScanned int64  `json:"bytesScanned,omitempty"`
	CurrentPath  string `json:"currentPath,omitempty"`
	Error        string `json:"error,omitempty"`
	ElapsedMs    int64  `json:"elapsedMs,omitempty"`
}

type ScanJob struct {
	ID           string      `json:"id"`
	Root         string      `json:"root"`
	Status       ScanStatus  `json:"status"`
	StartedAt    time.Time   `json:"startedAt"`
	CompletedAt  *time.Time  `json:"completedAt,omitempty"`
	Error        string      `json:"error,omitempty"`
	Tree         *ScanNode   `json:"tree,omitempty"`
	LargestFiles []FileEntry `json:"largestFiles,omitempty"`
	DirsScanned  int64       `json:"dirsScanned"`
	FilesScanned int64       `json:"filesScanned"`
	BytesScanned int64       `json:"bytesScanned"`
	CurrentPath       string          `json:"currentPath"`
	Insights          *InsightsReport `json:"insights,omitempty"`
	LastCleanupReport *CleanupReport  `json:"lastCleanupReport,omitempty"`
}

type ScanResponse struct {
	ScanID string `json:"scanId"`
}

type StartScanRequest struct {
	Root string `json:"root"`
}

type ExpandScanRequest struct {
	Path  string `json:"path"`
	Depth int    `json:"depth"`
}

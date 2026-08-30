package wsl

import "errors"

var errUnsupported = errors.New("WSL/VHDX compact is only supported on Windows")

// Disk describes a compactable virtual disk on the host.
type Disk struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"sizeBytes"`
	Kind      string `json:"kind"`
}

// CompactRequest controls VHDX compact dry-run or execute.
type CompactRequest struct {
	Path          string `json:"path"`
	DryRun        bool   `json:"dryRun"`
	Confirm       bool   `json:"confirm"`
	ConfirmPhrase string `json:"confirmPhrase"`
}

// CompactReport is the result of compact dry-run or execute.
type CompactReport struct {
	DryRun      bool   `json:"dryRun"`
	Path        string `json:"path"`
	BytesBefore int64  `json:"bytesBefore"`
	BytesAfter  int64  `json:"bytesAfter"`
	FreedBytes  int64  `json:"freedBytes"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
	Supported   bool   `json:"supported"`
}

const ConfirmPhrase = "DELETE"

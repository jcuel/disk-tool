package recycle

import (
	"context"
	"errors"
)

const ConfirmPhrase = "DELETE"

var (
	errConfirmRequired = errors.New("confirm required")
	errConfirmPhrase   = errors.New("confirmPhrase must be DELETE")
)

// BinInfo describes the platform recycle/trash bin.
type BinInfo struct {
	Path       string `json:"path"`
	ItemCount  int    `json:"itemCount"`
	TotalBytes int64  `json:"totalBytes"`
	Supported  bool   `json:"supported"`
	Error      string `json:"error,omitempty"`
}

// EmptyRequest controls recycle bin empty dry-run or execute.
type EmptyRequest struct {
	DryRun        bool   `json:"dryRun"`
	Confirm       bool   `json:"confirm"`
	ConfirmPhrase string `json:"confirmPhrase"`
}

// EmptyReport is the result of inspect-only or empty execute.
type EmptyReport struct {
	DryRun     bool   `json:"dryRun"`
	ItemCount  int    `json:"itemCount"`
	TotalBytes int64  `json:"totalBytes"`
	Error      string `json:"error,omitempty"`
}

// Inspect returns recycle bin size without deleting.
func Inspect() *BinInfo {
	return inspectPlatform()
}

// Empty dry-runs or empties the recycle bin / trash.
func Empty(ctx context.Context, req EmptyRequest) (*EmptyReport, error) {
	info := Inspect()
	if !info.Supported {
		return &EmptyReport{DryRun: req.DryRun, Error: info.Error}, errors.New(info.Error)
	}
	if req.DryRun {
		return &EmptyReport{
			DryRun:     true,
			ItemCount:  info.ItemCount,
			TotalBytes: info.TotalBytes,
		}, nil
	}
	if !req.Confirm {
		return nil, errConfirmRequired
	}
	if req.ConfirmPhrase != ConfirmPhrase {
		return nil, errConfirmPhrase
	}
	return emptyPlatform(ctx)
}

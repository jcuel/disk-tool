package api

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/safety"
)

const cleanupConfirmPhrase = "DELETE"

var errCleanupPathsRequired = errors.New("paths required")
var errCleanupConfirmPhrase = errors.New("confirmPhrase must be DELETE")

type pathItem struct {
	path string
	size int64
}

func RunCleanup(job *model.ScanJob, req model.CleanupRequest) (*model.CleanupReport, error) {
	if len(req.Paths) == 0 {
		return nil, errCleanupPathsRequired
	}
	if !req.DryRun {
		if !req.Confirm {
			return nil, errConfirmRequired
		}
		if req.ConfirmPhrase != cleanupConfirmPhrase {
			return nil, errCleanupConfirmPhrase
		}
	}

	knownSizes := candidateSizes(job)
	categories := candidateCategories(job)
	unique := dedupePaths(req.Paths)
	items := make([]pathItem, 0, len(unique))
	for _, p := range unique {
		p = SanitizePath(p)
		if p == "" {
			continue
		}
		items = append(items, pathItem{path: p, size: lookupCandidateSize(p, knownSizes)})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].size > items[j].size
	})

	report := &model.CleanupReport{
		DryRun:         req.DryRun,
		StartedAt:      time.Now().UTC(),
		TotalRequested: len(items),
		Results:        make([]model.CleanupItemResult, 0, len(items)),
	}

	rootAbs, _ := filepath.Abs(job.Root)

	for _, item := range items {
		result := model.CleanupItemResult{
			Path:     item.path,
			Size:     item.size,
			Category: categories[item.path],
		}

		abs, err := PathWithinRoot(job.Root, item.path)
		if err != nil {
			result.Status = model.CleanupStatusSkippedOutside
			result.Reason = err.Error()
			report.Results = append(report.Results, result)
			continue
		}
		if abs == rootAbs {
			result.Status = model.CleanupStatusSkippedScanRoot
			result.Reason = "cannot delete scan root"
			report.Results = append(report.Results, result)
			continue
		}

		zone := safety.ClassifyPath(abs)
		if !safety.IsDeletable(zone) {
			result.Status = model.CleanupStatusSkippedProtected
			result.Reason = "protected zone: " + string(zone)
			report.Results = append(report.Results, result)
			continue
		}

		info, err := os.Stat(abs)
		if err != nil {
			if os.IsNotExist(err) {
				result.Status = model.CleanupStatusSkippedMissing
				result.Reason = "path not found"
			} else {
				result.Status = model.CleanupStatusFailed
				result.Reason = err.Error()
			}
			report.Results = append(report.Results, result)
			continue
		}

		result.Size = resolveItemSize(abs, knownSizes, req.DryRun, info)

		if PathInUse(abs) {
			result.Status = model.CleanupStatusSkippedLocked
			result.Reason = "file in use"
			report.Results = append(report.Results, result)
			continue
		}

		if req.DryRun {
			result.Status = model.CleanupStatusWouldDelete
			report.Results = append(report.Results, result)
			continue
		}

		if err := DeletePath(abs); err != nil {
			if isLockError(err) {
				result.Status = model.CleanupStatusSkippedLocked
				result.Reason = err.Error()
			} else {
				result.Status = model.CleanupStatusFailed
				result.Reason = err.Error()
			}
			report.Results = append(report.Results, result)
			continue
		}

		result.Status = model.CleanupStatusDeleted
		report.BytesReclaimed += result.Size
		report.Results = append(report.Results, result)
	}

	report.FinishedAt = time.Now().UTC()
	report.ReportText = model.BuildCleanupReportText(report)
	return report, nil
}

func lookupCandidateSize(path string, known map[string]int64) int64 {
	if s, ok := known[filepath.Clean(path)]; ok {
		return s
	}
	return 0
}

func resolveItemSize(abs string, known map[string]int64, dryRun bool, info os.FileInfo) int64 {
	if s, ok := known[filepath.Clean(abs)]; ok && s > 0 {
		return s
	}
	if !info.IsDir() {
		return info.Size()
	}
	if dryRun {
		return 0
	}
	return dirSize(abs)
}

func candidateSizes(job *model.ScanJob) map[string]int64 {
	out := make(map[string]int64)
	if job == nil || job.Insights == nil {
		return out
	}
	for _, c := range job.Insights.CleanupCandidates {
		out[filepath.Clean(c.Path)] = c.Size
	}
	return out
}

func dedupePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		p = SanitizePath(p)
		if p == "" {
			continue
		}
		key := filepath.Clean(p)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}
	return out
}

func candidateCategories(job *model.ScanJob) map[string]string {
	out := make(map[string]string)
	if job == nil || job.Insights == nil {
		return out
	}
	for _, c := range job.Insights.CleanupCandidates {
		out[c.Path] = string(c.Category)
	}
	return out
}

func dirSize(path string) int64 {
	var total int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total
}

func pruneDeletedCandidates(job *model.ScanJob, report *model.CleanupReport) {
	if job == nil || job.Insights == nil || report == nil {
		return
	}
	deleted := make(map[string]struct{})
	for _, r := range report.Results {
		if r.Status == model.CleanupStatusDeleted {
			deleted[filepath.Clean(r.Path)] = struct{}{}
		}
	}
	if len(deleted) == 0 {
		return
	}
	kept := make([]model.CleanupCandidate, 0, len(job.Insights.CleanupCandidates))
	var reclaim int64
	for _, c := range job.Insights.CleanupCandidates {
		if _, ok := deleted[c.Path]; ok {
			continue
		}
		kept = append(kept, c)
		reclaim += c.Size
	}
	job.Insights.CleanupCandidates = kept
	job.Insights.TotalReclaimable = reclaim
}

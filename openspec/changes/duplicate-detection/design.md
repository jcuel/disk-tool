# Design: Duplicate file detection

**Change:** duplicate-detection

## Approach

1. `internal/dedup` package: hash files with SHA-256, bucket by hash.
2. Skip symlinks; respect scan root via existing `PathWithinRoot`.
3. Store `DuplicateGroup[]` on `ScanJob` after explicit `POST /api/scans/{id}/duplicates` or lazy during full branch scan (v1: on-demand endpoint).
4. UI tab or section under Insights listing groups sorted by wasted bytes.

## Performance

- Parallel workers bounded by `GOMAXPROCS`.
- Skip files over 500MB unless user confirms (future).

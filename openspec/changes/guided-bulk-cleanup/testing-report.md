# Testing report: guided-bulk-cleanup

**Change:** guided-bulk-cleanup
**Date:** 2026-07-05

## Automated

- [x] `go test ./...` — pass
- [x] `npm run build` (web) — pass
- [x] `scripts/smoke-api.ps1` — pass (includes cleanup dry-run)
- [x] `scripts/smoke-api.sh` — updated with cleanup dry-run check

## Unit coverage

- `internal/api/cleanup_test.go` — dry-run, confirm gates, delete, skip root/outside
- `internal/model/cleanup_report_test.go` — report text builder

## Manual (recommended)

- [ ] Select safe (review) candidates, dry-run, confirm with DELETE, verify report exports
- [ ] Verify locked file is skipped with reason in report

## Notes

- Tree sizes may be stale after bulk delete until re-scan or drill-down
- Docker cleanup deferred to phase 2

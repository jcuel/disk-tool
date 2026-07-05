# Proposal: Guided bulk cleanup

**Change:** guided-bulk-cleanup
**Status:** in-progress

## Summary

Batch delete scan-root cleanup candidates with dry-run, double confirmation, lock preflight, and downloadable reports.

## Scope

- `POST /api/scans/{id}/cleanup` (dry-run + execute)
- Cleanup report model and export formats
- UI multi-select, review/confirm modals, report panel

## Out of Scope

- Docker/container resources (`docker-cleanup` follow-up)

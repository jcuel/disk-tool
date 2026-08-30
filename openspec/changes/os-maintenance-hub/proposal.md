# Proposal: os-maintenance-hub

**Change:** `os-maintenance-hub`
**Status:** applying
**Issue:** Recycle bin + guided OS cleanup

## Summary

Replace passive protected-zone messaging with actionable maintenance cards and a cross-platform empty-recycle-bin API under `internal/maintenance/recycle/`.

## Approach

- New package `internal/maintenance/recycle/` with OS-specific implementations
- Handlers in `internal/api/handlers_maintenance.go`
- Maintenance hub section in Insights UI

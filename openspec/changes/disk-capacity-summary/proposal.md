# Proposal: Disk capacity summary

**Change:** disk-capacity-summary
**Status:** in-progress

## Summary

Show volume-level disk usage (pie chart + stats) before drill-down scan completes. Auto-start overview on system drive. Fix largest-files grid overlap and table layout.

## Scope

- `internal/diskspace` + `GET /api/disk`
- Disk summary panel + auto-scan in UI
- Grid layout fix for files panel
- Largest files: name, path, size, open, delete

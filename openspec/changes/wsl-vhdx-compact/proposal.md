# Proposal: wsl-vhdx-compact

**Change:** `wsl-vhdx-compact`
**Status:** applying
**Issue:** Compact Docker/WSL VHDX on Windows

## Summary

Windows-only maintenance workflow to list and compact VHDX files after Docker/WSL usage, reporting actual on-disk size change.

## Approach

- `internal/maintenance/wsl/` with `wsl_windows.go` + non-Windows stub
- Linked from Docker maintenance card copy

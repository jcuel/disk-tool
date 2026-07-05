# Proposal: Open folder in explorer

**Change:** cleanup-insights-explorer
**Status:** in-progress

## Summary

Add Open action on tree, cleanup, and largest-files rows to launch the OS file manager for paths under the scan root.

## Scope

- POST `/api/scans/{id}/open` with path validation
- UI Open buttons

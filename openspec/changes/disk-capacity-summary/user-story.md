# User Story: Disk capacity summary

**Change:** disk-capacity-summary
**Status:** refined

## Story

As a user, I want to see volume capacity (used vs free) at a glance when disk-tool opens, with an automatic overview scan of the system drive.

## Acceptance Criteria

- [x] GET /api/disk returns capacity, used, free for scan volume
- [x] Pie chart and text stats at top of UI
- [x] Auto-scan system drive (C:\) on load
- [x] Largest files panel visible with size, open, delete
- [x] Path column readable (filename + path, no character-by-character wrap)

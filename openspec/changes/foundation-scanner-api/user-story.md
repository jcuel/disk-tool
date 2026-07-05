# User Story: Local web disk space analyzer

**Change:** `foundation-scanner-api`
**Status:** refined

## Story

As a user, I want to open a local web page, pick a folder or drive, and see how disk space is distributed across subfolders with a treemap and sortable tree, so I can find what is consuming space.

## Context

TreeSize-inspired cross-platform tool. Local-only (127.0.0.1). Go backend + Vite/TypeScript frontend with ECharts. Delivered in four phases: scanner API, web UI, charts, export/hardening.

## Acceptance Criteria

- [ ] User can start `disk-tool serve` and open the UI in a browser on localhost
- [ ] User can enter or pick a root path and start a scan
- [ ] Scan progress updates live via WebSocket
- [ ] User can cancel an in-progress scan
- [ ] Folder tree shows name, size, percent of parent, and file count sorted by size
- [ ] Treemap and bar chart visualize children of the selected folder
- [ ] Largest files panel lists top 100 files globally for a completed scan
- [ ] User can drill into subfolders from tree or chart clicks
- [ ] User can export scan results as JSON or HTML
- [ ] Headless CLI `disk-tool scan /path --json` works without the UI
- [ ] Docker smoke test builds and runs a sample scan

## Out of Scope

- Duplicate file detection
- Cloud/NAS/SSH remote scanning
- Delete or move files from the UI
- PDF/Excel export

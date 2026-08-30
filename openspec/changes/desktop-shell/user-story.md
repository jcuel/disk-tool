# User Story: Desktop app shell

**Change:** `desktop-shell`
**Status:** refined

## Story

As a disk-tool user, I want a native desktop app I can install and launch without opening a browser or terminal, so that disk analysis feels like a normal desktop application.

## Context

disk-tool v1 ships as a Go binary that serves the embedded Vite UI on localhost. Tauri was deferred in `foundation-scanner-api` for localhost-web distribution. This change adds a Tauri 2 shell that bundles the Go binary as a sidecar and loads the existing UI in an OS webview.

## Acceptance Criteria

- [ ] Double-clicking the desktop app opens a native window with the disk-tool UI
- [ ] Scan, drill-down, delete, cleanup, and Docker flows work in the desktop window
- [ ] Quitting the app stops the Go sidecar (no orphan processes)
- [ ] Installers build for Windows, macOS, and Linux on release tags
- [ ] CLI `disk-tool serve` and `disk-tool scan` remain available as standalone binaries

## Out of Scope

- Rewriting the Go backend in Rust
- React migration or frontend framework changes
- Code signing (optional follow-up)
- System tray and single-instance lock (v1.1)

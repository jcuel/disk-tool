# Proposal: Tauri desktop shell

**Change:** desktop-shell
**Status:** in-progress
**Domain:** disk-tool

## Summary

Add a Tauri 2 desktop application that spawns the existing Go binary as a sidecar, navigates the webview to the sidecar URL, and packages cross-platform installers alongside the current CLI releases.

## Motivation

End users should install disk-tool like a normal desktop app—no terminal, no browser tab, no port conflicts. The Go scanner/API stack stays unchanged; Tauri provides window chrome, taskbar presence, and installer packaging (mirroring the multi-git release workflow).

## Scope

### In scope

- Tauri 2 project under `desktop/` with Go sidecar lifecycle
- Go `serve` hardening: dynamic port, readiness stdout, graceful shutdown
- Build scripts and release CI for desktop installers
- README desktop install section

### Out of scope

- Rust rewrite of scanner/API
- Vite HMR dev proxy (optional follow-up)
- Code signing and auto-update

## Risks

| Risk | Mitigation |
|------|------------|
| Orphan sidecar processes | Kill sidecar on Tauri exit; smoke test shutdown |
| Port conflicts | `--port 0` with dynamic assignment |
| Duplicate UI bundles | Go embeds UI only; Tauri is a thin shell |
| CI complexity (Rust + WebKit) | Sidecar smoke in CI; full Tauri build on release tags |

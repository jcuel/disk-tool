# Testing report: desktop-shell

**Date:** 2026-08-29

| Check | Result | Notes |
|-------|--------|-------|
| Unit — parseReadyPort | pass | `cmd/disk-tool/serve_test.go` |
| Unit — serve dynamic port + ready stdout | pass | helper subprocess test |
| Sidecar smoke | pass | `scripts/smoke-sidecar.sh` |
| Tauri Rust compile | pass | `cargo build --release` in `desktop/src-tauri` |
| Full Tauri bundle (Linux) | skipped locally | requires WebKit/GTK; CI `desktop-check` job |

## Manual verification

- [ ] Desktop window opens against sidecar URL on Windows
- [ ] Desktop window opens against sidecar URL on macOS
- [ ] Quitting desktop app stops sidecar process
- [ ] Release installers attach to GitHub Release on tag push

## Notes

Closes #99. Desktop v1 uses Go sidecar only (no duplicate frontend bundle in Tauri).

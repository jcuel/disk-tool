# Testing report: desktop-shell

**Date:** 2026-08-29

| Check | Result | Notes |
|-------|--------|-------|
| Unit — parseReadyPort | pass | `cmd/disk-tool/serve_test.go` |
| Unit — serve dynamic port + ready stdout | pass | helper subprocess test |
| Sidecar smoke | pass | `scripts/smoke-sidecar.sh` |
| Tauri Rust compile | pass | `cargo check` in `desktop/src-tauri` |
| Full Tauri bundle (Linux) | skipped locally | requires WebKit/GTK; CI `desktop-check` job |

## PR #100 review fixes (2026-08-29)

| Fix | Result | Notes |
|-----|--------|-------|
| Release publish consolidated | done | `publish-release` waits for CLI + desktop artifacts |
| Sidecar kill on timeout/error | done | all error paths in `start_sidecar` |
| Capabilities trimmed | done | `default.json` = `core:default`; `sidecar.json` for Rust spawn |
| Localhost navigation guard | done | `WebviewWindowBuilder` + `on_navigation` |

## Manual verification

- [ ] Desktop window opens against sidecar URL on Windows
- [ ] Desktop window opens against sidecar URL on macOS
- [ ] Quitting desktop app stops sidecar process
- [ ] Release installers attach to GitHub Release on tag push (single `publish-release` job)

## Notes

Closes #99. Desktop v1 uses Go sidecar only (no duplicate frontend bundle in Tauri).

# Design: Tauri desktop shell

**Change:** `desktop-shell`

## Approach

Tauri 2 hosts an OS webview. On startup, Rust spawns the Go binary (`disk-tool serve --no-open --port 0 --ready-stdout`) as an `externalBin` sidecar, parses `disk-tool-ready port=N` from stdout, and navigates the main window to `http://127.0.0.1:N`. On app exit, the sidecar is killed.

The Go binary continues to embed `web/dist` and serve REST/WebSocket APIs unchanged.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Tauri + Go sidecar | Reuses all Go code; multi-git-style installers | Two processes | Selected |
| Wails (Go webview) | Single process | Different release pipeline; refactor serve loop | Rejected for v1 |
| Tauri + Rust backend | Native invoke API | Full backend rewrite | Rejected |
| Browser-only | Already works | Poor end-user UX | Superseded for distribution |

## Components

```
desktop/
  package.json
  src-tauri/
    tauri.conf.json       bundle.externalBin, window config
    src/lib.rs              sidecar spawn, URL navigation, shutdown
    src/main.rs
    capabilities/default.json
    binaries/               platform Go binary (build output)

cmd/disk-tool/main.go       --port 0, --ready-stdout, signal shutdown
scripts/build-desktop.sh    web → embed → go → tauri build
build.ps1 / Makefile          desktop targets
.github/workflows/            release tauri-action + sidecar smoke
```

## Data & API Touchpoints

| Component | Input | Output |
|-----------|-------|--------|
| Go sidecar | `serve --no-open --port 0 --ready-stdout` | HTTP UI + `/api/*` on dynamic localhost port |
| Tauri shell | Sidecar stdout line | Webview URL `http://127.0.0.1:{port}` |
| Release CI | Go cross-compile + Rust tauri build | NSIS/MSI, DMG, AppImage/deb + CLI binaries |

No new REST routes. Existing `POST /api/scans/{id}/open` handles file-manager launch from Go.

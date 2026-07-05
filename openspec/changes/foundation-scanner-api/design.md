# Design: disk-tool v1

**Change:** `foundation-scanner-api`

## Approach

Single Go binary serves REST/WebSocket API and embedded static frontend. Scanner uses bounded goroutine pool per directory walk with bottom-up size aggregation.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Go + embed Vite | Fast scan, one binary | Two build steps | Selected |
| Python FastAPI | Rapid prototype | Slower scans | Rejected |
| Tauri desktop | Native feel | Not local-web as planned | Rejected |

## Components

```
cmd/disk-tool/main.go       CLI entry
internal/model/             ScanNode, ScanJob, events
internal/scanner/           Parallel walk, cancel, largest files
internal/api/               HTTP handlers, store, path validation
web/                        Vite + TS + ECharts UI
```

## Data & API Touchpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/roots` | Common scan roots per OS |
| POST | `/api/scans` | `{ "root": "..." }` → `{ "scanId" }` |
| GET | `/api/scans/{id}` | Status + tree |
| DELETE | `/api/scans/{id}` | Cancel |
| WS | `/api/scans/{id}/events` | Progress stream |
| GET | `/api/scans/{id}/export?format=json\|html` | Download report |

## Cross-platform

- `filepath` for paths; skip permission errors; detect symlink cycles via `sys` stat device+inode where available.

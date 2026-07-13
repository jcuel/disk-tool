# Testing Report: Protect disk images and confirm single deletes

**Change:** `protect-disk-images`
**Date:** 2026-07-13

## Automated

| Check | Result |
|-------|--------|
| `go test ./internal/safety/` | pass (includes `IsDiskImagePath`, `CanDeletePath`) |
| `go test ./internal/insights/` | pass |
| `go test ./internal/api/` | pass |
| Vite production build | pass |
| Local `disk-tool serve` on :8080 | pass (HTTP 200 after embed) |

## Manual

| Check | Result |
|-------|--------|
| Largest Files Delete opens Review delete modal (not browser prompt) | pass (hard-refresh required) |
| Dry-run + checkbox + type DELETE required | pass (implemented) |
| VHD/VHDX Delete disabled / API blocked | pass (server + unit tests) |

## Security notes

- Disk-image guard is extension-based on the path basename; rename bypass remains a residual risk
- Confirmation phrase still required server-side for cleanup execute

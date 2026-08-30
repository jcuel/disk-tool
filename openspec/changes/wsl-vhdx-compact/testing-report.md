# Testing report: wsl-vhdx-compact

**Date:** 2026-08-30

| Check | Result | Notes |
|-------|--------|-------|
| API — wsl disks (Linux) | pass | returns `supported: false` via stub |
| Web build | pass | `npm run build` |
| Go tests | pass | `go test ./...` |

## Manual verification

- [ ] Windows: list VHDX paths and sizes
- [ ] Compact after Docker prune shows NTFS size decrease
- [ ] `wsl --shutdown` warning shown before execute

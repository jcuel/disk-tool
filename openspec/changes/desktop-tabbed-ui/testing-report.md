# Testing report: desktop-tabbed-ui

**Date:** 2026-08-30

| Check | Result | Notes |
|-------|--------|-------|
| Web build | pass | `npm run build` |
| Browser regression | pass | no `desktop-app` without Tauri |
| Go tests | pass | `go test ./...` |

## Manual verification

- [ ] Tauri 1280×900: no page scroll; tabs switch panes
- [ ] Chart click drill-down on Browse tab
- [ ] Browser / GitHub Pages layout unchanged

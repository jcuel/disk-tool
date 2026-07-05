# Testing report: e2e-browser-tests

**Date:** 2026-07-05

- [x] `go test ./...`
- [x] `scripts/smoke-api.ps1`
- [ ] `scripts/e2e-run.sh` — pending CI `e2e-linux` job on PR
- [x] Manual review: Cypress specs for overview, drill-down, layout

**Security:** Cypress is a dev-only dependency installed at E2E runtime; no new production endpoints.

**Notes:** E2E uses `testdata/e2e-root/` via `?root=` query param. Chart canvas clicks remain API-smoke / manual only.

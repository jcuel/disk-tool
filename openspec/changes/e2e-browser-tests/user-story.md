# User Story: Cypress E2E browser tests

**Change:** e2e-browser-tests
**Status:** refined

## Story

As a maintainer, I want automated browser tests for critical UI flows so that regressions in scan, drill-down, and layout are caught before merge.

## Acceptance Criteria

- [x] Cypress runs headless against a local `disk-tool serve` instance in CI (Linux)
- [x] Tests use `testdata/e2e-root/` fixture, not system drives
- [x] `?root=` and `?noAutoScan=1` query params support deterministic test entry
- [x] MVP specs: overview scan, tree drill-down, layout smoke
- [x] `scripts/e2e-run.sh` documented in README

## Out of Scope

- Cypress Cloud, Windows E2E job, ECharts canvas clicks, visual regression

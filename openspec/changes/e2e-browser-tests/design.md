# Design: e2e-browser-tests

**Change:** e2e-browser-tests

## Runner

Cypress 13+ in [`web/`](web/package.json) devDependencies. Headless `cypress run` only — no Cypress Cloud.

## Server lifecycle

[`scripts/e2e-run.sh`](scripts/e2e-run.sh) mirrors [`scripts/smoke-api.sh`](scripts/smoke-api.sh): build embed static, `go build`, start `serve --port 18081`, run Cypress, trap EXIT to kill server.

## Test entry

`CYPRESS_SCAN_ROOT` env → `Cypress.env('scanRoot')`. Visit `/?root=<fixture>` to scan fixture on load. `?noAutoScan=1` defers scan for manual `#start-btn` tests.

## CI

Self-contained `e2e-linux` job on `ubuntu-latest` after unit smoke pattern (no artifact coupling).

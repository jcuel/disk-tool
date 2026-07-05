# Proposal: Cypress E2E browser tests

**Change:** e2e-browser-tests
**Status:** archived

## Summary

Add free open-source Cypress tests for overview scan, tree drill-down, and layout smoke. Wire a Linux CI job and minimal URL hooks so tests avoid auto-scanning `C:\`.

## Scope

- `web/cypress/` + `cypress.config.ts`
- `testdata/e2e-root/` fixture
- `?root=` / `?noAutoScan=1` in `web/src/main.ts`
- `scripts/e2e-run.sh`, CI job `e2e-linux`

# Design: GitHub Pages site

## Architecture

Two static bundles deployed as one Pages artifact:

| Path | Source | Purpose |
|------|--------|---------|
| `/disk-tool/` | `site/dist/` | Landing page |
| `/disk-tool/demo/` | `web/` demo build | UI + mock API |

## Demo mode

`VITE_DEMO_MODE=true` routes [`web/src/api.ts`](../../../web/src/api.ts) exports to [`web/src/demo/mock-api.ts`](../../../web/src/demo/mock-api.ts), which serves frozen JSON from [`web/src/demo/fixtures.json`](../../../web/src/demo/fixtures.json).

Fixtures are exported from a real scan of `testdata/e2e-root/` via [`scripts/export-demo-fixtures.sh`](../../../scripts/export-demo-fixtures.sh).

## Deployment

GitHub Actions `pages.yml` builds both packages on `master` push and deploys via `actions/deploy-pages@v4`. PRs run build-only smoke in CI.

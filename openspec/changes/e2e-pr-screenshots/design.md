# Design: e2e-pr-screenshots

**Change:** e2e-pr-screenshots

## Capture

Custom `cy.captureStep(name)` in [`web/cypress/support/commands.ts`](web/cypress/support/commands.ts). Viewport `1280×800` in [`web/cypress.config.ts`](web/cypress.config.ts). Four snapshots per run:

| Name | Spec | Capture |
|------|------|---------|
| `01-overview-ready` | overview | viewport |
| `02-manual-scan-start` | overview (noAutoScan) | viewport |
| `03-drill-down-big-dir` | drill-down | viewport |
| `04-layout-full` | layout | fullPage |

Runner [`scripts/e2e-run.sh`](scripts/e2e-run.sh) unchanged; Cypress writes to `web/cypress/screenshots/`.

## CI artifact

After `bash scripts/e2e-run.sh`, `actions/upload-artifact@v4` uploads `web/cypress/screenshots/` with 14-day retention. Runs on all E2E executions (`if: always()`).

## Inline PR comment

[`scripts/e2e-pr-screenshots-comment.sh`](scripts/e2e-pr-screenshots-comment.sh) pushes PNGs to branch `e2e-screenshots/pr-<number>` and writes markdown with `raw.githubusercontent.com` image URLs. [`peter-evans/create-or-update-comment@v4`](https://github.com/peter-evans/create-or-update-comment) posts/updates one gallery comment (`comment-tag: e2e-screenshots-gallery`).

PR-only: `if: github.event_name == 'pull_request' && always()`.

## Fallback comment

When publish fails or no PNGs exist, peter-evans posts/updates a comment linking to the workflow run and artifact name. Marker: `comment-tag: e2e-screenshots-fallback`.

## Fork PR caveat

Workflows triggered from fork PRs receive a read-only `GITHUB_TOKEN`; branch push and inline comment may be skipped. Artifact upload still works from the upstream repo workflow when the PR is from a branch in the same repo.

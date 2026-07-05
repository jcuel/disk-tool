# Proposal: Cypress PR screenshot gallery

**Change:** e2e-pr-screenshots
**Status:** archived

## Summary

Add intentional Cypress screenshots at key UI milestones and publish them on pull requests via a sticky inline comment plus CI artifact backup.

## Scope

- `captureStep` custom command and 4 capture points in existing specs
- `e2e-linux` job: artifact upload + PR comment (opengisch/comment-pr-with-images + peter-evans fallback)
- README note

## Out of scope

Cypress Cloud, visual regression baselines, Windows E2E, ECharts canvas tests.

## Risks

| Risk | Mitigation |
|------|------------|
| Screenshot branch clutter | Branch scoped per PR (`e2e-screenshots/pr-N`) |
| Fork PR read-only token | Artifact + fallback comment with run link |
| Flaky chart renders | Capture only after `waitForOverviewReady` assertions |

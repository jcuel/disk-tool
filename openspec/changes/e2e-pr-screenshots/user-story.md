# User Story: Cypress PR screenshot gallery

**Change:** e2e-pr-screenshots
**Status:** refined

## Story

As a reviewer, I want Cypress E2E to capture UI snapshots and show them on pull requests so that I can visually confirm scan, drill-down, and layout without running the app locally.

## Acceptance Criteria

- [ ] `captureStep` Cypress command captures 4 named snapshots per run (overview, manual scan, drill-down, layout)
- [ ] Fixed viewport (1280×800) for stable screenshots across runs
- [ ] CI uploads `web/cypress/screenshots/` as artifact on every E2E run
- [ ] PRs get a sticky comment with inline screenshot gallery (updated on re-push)
- [ ] Fallback comment with workflow/artifact link when inline upload is unavailable (e.g. fork PRs)
- [ ] README documents PR screenshot behavior

## Out of Scope

- Cypress Cloud, Percy/Chromatic baselines, Windows E2E job, ECharts canvas pixel tests, visual regression diffs

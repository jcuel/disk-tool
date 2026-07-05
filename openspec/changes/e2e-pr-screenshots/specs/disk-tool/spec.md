## MODIFIED Requirements

### Requirement: Browser E2E tests

The project SHALL provide Cypress browser tests that run headless against a local server using a deterministic scan fixture under `testdata/e2e-root/`. On pull requests, E2E SHALL capture named UI snapshots at key flow milestones, upload screenshots as CI artifacts, and post or update a sticky PR comment with inline images (with artifact link fallback when inline upload is unavailable).

#### Scenario: PR screenshot gallery

- **WHEN** Cypress E2E runs on a pull request
- **THEN** key UI states are captured as named PNGs
- **AND** a sticky PR comment shows inline snapshots
- **AND** screenshots are uploaded as a downloadable CI artifact

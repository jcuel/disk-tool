## ADDED Requirements

### Requirement: Browser E2E tests

The project SHALL provide Cypress browser tests that run headless against a local server using a deterministic scan fixture under `testdata/e2e-root/`.

#### Scenario: Overview scan in browser

- **WHEN** Cypress visits the UI with `?root=` pointing at the E2E fixture
- **THEN** the overview scan completes and the folder tree shows child rows

#### Scenario: Tree drill-down in browser

- **WHEN** the user clicks a folder row in the tree after overview completes
- **THEN** the breadcrumb updates to that folder path

#### Scenario: CI E2E gate

- **WHEN** a pull request is opened
- **THEN** the Linux Cypress job runs and must pass before merge

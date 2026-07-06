# disk-tool — System Specification

## Purpose

Cross-platform local web disk usage analyzer. Users scan a folder or drive on localhost and explore usage via folder tree, treemap, bar chart, and largest-files views.

## Requirements

### Requirement: Project documentation

The repository SHALL include a root README describing purpose and development workflow.

#### Scenario: Developer onboarding

- **WHEN** a developer opens the repository
- **THEN** README.md explains disk-tool purpose and SPECBOOT slash commands

### Requirement: Local web disk scan

The system SHALL provide a localhost web interface to scan a user-selected directory and display disk usage distribution.

#### Scenario: Start scan from UI

- **WHEN** the user submits a valid root path and clicks Start
- **THEN** a scan job starts and progress events stream to the browser

#### Scenario: Cancel scan

- **WHEN** the user clicks Cancel during an active scan
- **THEN** the scan stops and status becomes cancelled

### Requirement: Folder tree view

The system SHALL display a sortable folder tree with name, logical size, percent of parent, and file count.

#### Scenario: Tree sorted by size

- **WHEN** a scan completes
- **THEN** folders are listed largest-first under each parent

### Requirement: Visual charts

The system SHALL render a treemap and bar chart for the selected folder's immediate children.

#### Scenario: Drill-down from chart

- **WHEN** the user clicks a folder segment in the treemap or bar chart
- **THEN** that folder becomes selected and expandable children trigger a deeper scan

#### Scenario: Drill-down from tree

- **WHEN** the user clicks a folder in the tree or chart
- **THEN** charts update to show that folder's children

### Requirement: Layout balance

The web UI SHALL use a two-column grid for tree vs charts/files with equal column width, SHALL span Insights across the full width below, and SHALL show scan-in-progress placeholders when overview data is not yet available.

### Requirement: Largest files

The system SHALL list the top 100 largest files found during a completed scan.

### Requirement: Export

The system SHALL export completed scan results as JSON or self-contained HTML.

### Requirement: Headless CLI

The system SHALL support `disk-tool scan <path> --json` without starting the web server.

### Requirement: Browser E2E tests

The project SHALL provide Cypress browser tests that run headless against a local server using a deterministic scan fixture under `testdata/e2e-root/`. On pull requests, E2E SHALL capture named UI snapshots at key flow milestones, upload screenshots as CI artifacts, and post or update a sticky PR comment with inline images (with artifact link fallback when inline upload is unavailable).

#### Scenario: Overview scan in browser

- **WHEN** Cypress visits the UI with `?root=` pointing at the E2E fixture
- **THEN** the overview scan completes and the folder tree shows child rows

#### Scenario: Tree drill-down in browser

- **WHEN** the user clicks a folder row in the tree after overview completes
- **THEN** the breadcrumb updates to that folder path

#### Scenario: CI E2E gate

- **WHEN** a pull request is opened
- **THEN** the Linux Cypress job runs and must pass before merge

#### Scenario: PR screenshot gallery

- **WHEN** Cypress E2E runs on a pull request
- **THEN** key UI states are captured as named PNGs
- **AND** a sticky PR comment shows inline snapshots
- **AND** screenshots are uploaded as a downloadable CI artifact

### Requirement: Localhost binding

The server SHALL bind to 127.0.0.1 only and reject paths outside the scan root.

### Requirement: Cleanup insights

The system SHALL detect pattern-based cleanup candidates (dev artifacts, caches, downloads) with risk tiers and export them in insights, HTML, JSON, and support ticket formats.

### Requirement: Safe delete

The system SHALL expose `POST /api/scans/{id}/delete` requiring `confirm: true` and `confirmPhrase: "DELETE"`. Paths outside the scan root, the scan root itself, and protected safety zones SHALL be rejected.

### Requirement: Open in file manager

The system SHALL expose `POST /api/scans/{id}/open` to launch the OS file manager for paths under the scan root.

### Requirement: Guided bulk cleanup

The system SHALL expose `POST /api/scans/{id}/cleanup` with dry-run and execute modes. Execute SHALL require typed confirmation and skip locked, missing, outside-root, and protected-zone paths.

### Requirement: OS safety zones

The system SHALL classify paths into safety zones (forbidden, critical_os, diagnostic, caution, review, maintenance). Forbidden paths SHALL be skipped during scan. Delete and bulk cleanup SHALL block critical_os, diagnostic, and forbidden zones.

### Requirement: Safety grid

The insights report SHALL include a safety grid summarizing candidate counts and bytes per zone. The web UI SHALL display the grid with zone and risk badges on cleanup candidates.

### Requirement: Maintenance presets

The system SHALL expose `GET /api/scans/{id}/maintenance-presets` returning preset definitions and matched deletable paths for one-click maintenance flows.

### Requirement: Age-based cleanup

The system SHALL flag stale large files under Downloads and temp locations via configurable age and size thresholds. The UI SHALL allow threshold adjustment and re-analyze via `POST /api/scans/{id}/reanalyze`.

### Requirement: Duplicate detection

The system SHALL expose `POST /api/scans/{id}/duplicates` to find duplicate file groups under the scan root, excluding protected zones.

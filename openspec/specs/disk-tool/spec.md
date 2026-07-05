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

### Requirement: Localhost binding

The server SHALL bind to 127.0.0.1 only and reject paths outside the scan root.

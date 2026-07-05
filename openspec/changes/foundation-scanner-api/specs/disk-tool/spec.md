# Delta Spec: disk-tool v1

**Change:** `foundation-scanner-api`
**Domain:** `disk-tool`

## ADDED Requirements

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

#### Scenario: Drill-down

- **WHEN** the user clicks a folder in the tree or chart
- **THEN** charts update to show that folder's children

### Requirement: Largest files

The system SHALL list the top 100 largest files found during a completed scan.

### Requirement: Export

The system SHALL export completed scan results as JSON or self-contained HTML.

### Requirement: Headless CLI

The system SHALL support `disk-tool scan <path> --json` without starting the web server.

### Requirement: Localhost binding

The server SHALL bind to 127.0.0.1 only and reject paths outside the scan root.

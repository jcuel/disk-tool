# Delta Spec: OS maintenance hub

**Change:** `os-maintenance-hub`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: Maintenance hub

The UI SHALL present a Maintenance section with actionable cards for temp cleanup, empty recycle bin, Docker reclaim, and a read-only critical OS explainer. Protected-zone bytes SHALL be shown to users with guidance to use maintenance actions instead of filesystem delete on protected paths.

#### Scenario: Protected bytes surfaced

- **WHEN** insights include a safety grid with `protectedBytes > 0`
- **THEN** the UI displays the protected total and points users to Maintenance actions

### Requirement: Empty recycle bin

The system SHALL expose recycle bin / trash inspection and empty APIs independent of scan ID. Execute SHALL require `confirm: true` and `confirmPhrase: "DELETE"`. Dry-run SHALL report item count and total bytes without deleting.

#### Scenario: Recycle dry-run

- **WHEN** `POST /api/maintenance/recycle/empty` with `dryRun: true`
- **THEN** the response reports item count and bytes without modifying the bin

#### Scenario: Recycle execute

- **WHEN** the user confirms with phrase DELETE
- **THEN** the OS recycle bin / trash is emptied on supported platforms

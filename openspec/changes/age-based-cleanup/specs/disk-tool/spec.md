# Delta Spec: Age-based cleanup suggestions

**Change:** age-based-cleanup
**Domain:** disk-tool

## ADDED Requirements

### Requirement: Age-based cleanup candidates

The system SHALL flag large files older than a configurable threshold as cleanup candidates.

#### Scenario: Default thresholds

- **WHEN** a scan completes with default settings (90 days, 50MB)
- **THEN** matching files appear in cleanup candidates with age and hint

#### Scenario: User adjusts thresholds

- **WHEN** the user changes age or size threshold in the UI
- **THEN** cleanup candidates refresh on next scan or re-analyze action

## MODIFIED Requirements

### Requirement: Cleanup insights

The insights engine SHALL include age-based stale file detection alongside pattern-based candidates.

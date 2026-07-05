# Delta Spec: Duplicate file detection

**Change:** duplicate-detection
**Domain:** disk-tool

## ADDED Requirements

### Requirement: Duplicate file detection

The system SHALL identify groups of files with identical content under the scan root.

#### Scenario: On-demand duplicate scan

- **WHEN** the user requests duplicate detection on a completed scan
- **THEN** duplicate groups are computed and returned with paths and wasted bytes

#### Scenario: Export includes duplicates

- **WHEN** the user exports JSON or HTML after duplicate scan
- **THEN** duplicate groups are included in the report

# Delta Spec: Docker maintenance

**Change:** `docker-cleanup`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: Docker maintenance

The system SHALL detect Docker CLI availability and report reclaimable Docker usage (images, containers, build cache) when the daemon is reachable. The system SHALL expose scan-scoped Docker status and prune APIs. Prune execute SHALL require `confirm: true` and `confirmPhrase: "DELETE"`, SHALL run `docker system prune` without volume deletion, and SHALL NOT delete Docker data-root directories or virtual disk image files via the filesystem.

#### Scenario: Docker status when CLI available

- **WHEN** Docker CLI is on PATH and the daemon responds
- **THEN** `GET /api/scans/{id}/docker` returns usage summary including reclaimable bytes

#### Scenario: Docker status when CLI missing

- **WHEN** Docker CLI is not available
- **THEN** the status response indicates Docker is unavailable
- **AND** insights MAY list known data-root paths as non-deletable caution candidates

#### Scenario: Confirmed prune

- **WHEN** the user confirms prune with phrase DELETE after dry-run review
- **THEN** the server runs Docker system prune without `--volumes`
- **AND** filesystem cleanup of Docker data roots remains blocked

### Requirement: Docker reclaim preset

The system SHALL provide a maintenance preset `docker-reclaim` that surfaces Docker reclaim candidates and drives the Docker prune confirmation flow (not bulk filesystem delete of Docker roots).

## MODIFIED Requirements

### Requirement: Cleanup insights

The system SHALL detect pattern-based cleanup candidates (dev artifacts, caches, downloads, Docker reclaimable usage) with risk tiers and export them in insights, HTML, JSON, and support ticket formats.

### Requirement: Maintenance presets

The system SHALL expose `GET /api/scans/{id}/maintenance-presets` returning preset definitions (including `docker-reclaim`) and matched deletable paths for one-click maintenance flows.

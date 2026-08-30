# Delta Spec: Docker prune honesty

**Change:** `fix-docker-prune-reporting`
**Domain:** `disk-tool`

## MODIFIED Requirements

### Requirement: Docker maintenance

The system SHALL detect Docker CLI availability and report reclaimable Docker usage (images, containers, build cache) when the daemon is reachable. The system SHALL expose scan-scoped Docker status and prune APIs. Prune execute SHALL require `confirm: true` and `confirmPhrase: "DELETE"`, SHALL run `docker system prune` without volume deletion, and SHALL NOT delete Docker data-root directories or virtual disk image files via the filesystem. After execute, reclaimed bytes SHALL equal `max(0, beforeReclaimable − afterReclaimable)` from consecutive `docker system df` reads; the system SHALL NOT inflate reclaimed bytes with the pre-prune estimate. The prune response SHALL include `noChange: true` when prune succeeded but reclaimable usage did not decrease.

#### Scenario: Prune with no df change

- **WHEN** the user confirms Docker prune and `docker system df` reclaimable total is unchanged or higher after prune
- **THEN** the response reports `reclaimable: 0` and `noChange: true`
- **AND** the UI warns that host disk may not shrink until WSL/VHDX compact on Windows

#### Scenario: Prune with positive reclaim

- **WHEN** reclaimable usage decreases after prune
- **THEN** `reclaimable` equals the before/after difference
- **AND** insights are re-analyzed for the scan

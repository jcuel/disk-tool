# User Story: Honest Docker prune reporting

**Change:** `fix-docker-prune-reporting`
**Status:** refined

## Story

As a developer who prunes Docker to free disk space, I want disk-tool to report only the bytes Docker actually reclaimed (before/after `docker system df`), so I am not misled when host disk space does not change — especially on Windows where VHDX files may not shrink.

## Context

Live QA (#110) showed delete UI staleness; Docker prune had a related trust issue: when `before.Reclaimable - after.Reclaimable <= 0`, the server could inflate reclaimed bytes with the pre-prune estimate. Users saw success alerts without insights refresh.

## Acceptance Criteria

- [x] `PruneReport` includes `beforeReclaimable`, `afterReclaimable`, `beforeDf`, `afterDf`, and `noChange`
- [x] `reclaimedBytes` is `max(0, before − after)` only — never inflated with pre-prune estimate
- [x] UI shows status banner (not generic alert) and warns when `noChange` or zero reclaimed
- [x] After prune: `reanalyzeInsights()` + `refreshJob()`
- [x] Docker modal copy explains layer prune ≠ VHDX shrink on Windows
- [x] Unit tests for reclaim math; handler returns `noChange` in JSON

## Out of Scope

- WSL/VHDX compact (separate change `wsl-vhdx-compact`)
- `docker volume prune` / `--volumes`

# User Story: OS maintenance hub

**Change:** `os-maintenance-hub`
**Status:** refined

## Story

As a disk-tool user, I want guided maintenance cards (temp cleanup, empty recycle bin, Docker reclaim, critical OS explainer) instead of passive "danger zone" messaging, so I know what safe OS-level actions to take for protected paths.

## Context

Safety zones block filesystem delete on critical paths but only showed badges. `protectedBytes` was computed but not rendered. Recycle bin empty was not implemented.

## Acceptance Criteria

- [x] Maintenance hub UI with action cards in Insights
- [x] Surface `protectedBytes` in safety grid
- [x] Zone labels with one-line safe-alternative guidance
- [x] `GET /api/maintenance/recycle` and `POST /api/maintenance/recycle/empty` (cross-platform)
- [x] Confirm gate: checkbox + type `DELETE` for execute
- [x] Dry-run reports item count and bytes
- [x] Critical OS card is read-only with OS tool hints

## Out of Scope

- WSL compact (separate change)
- Auto-delete without typed DELETE

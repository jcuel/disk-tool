# Proposal: CI project board sync

**Change:** ci-project-board-sync
**Status:** in-progress

## Summary

Run `scripts/sync-project-board.sh` automatically after CI succeeds on a push to `master`. Bash script is the single source of truth; PowerShell delegates to it.

## Scope

- Refactor sync script: env-driven IDs, dynamic OpenSpec proposal detection (no hardcoded issue numbers)
- `scripts/project-board.env` for board field IDs
- `.github/workflows/sync-project-board.yml` triggered by successful CI on merge pushes
- Repository secret `GH_PROJECT_SYNC` (PAT with `project` scope)

## Out of scope

- Creating or editing project views (UI-only)

# Proposal: fix-docker-prune-reporting

**Change:** `fix-docker-prune-reporting`
**Status:** applying
**Issue:** Docker prune reports false reclaim

## Summary

Fix Docker prune reporting to use honest before/after `docker system df` diff, expose structured fields in API/UI, and refresh insights after maintenance.

## Approach

1. Add `computeReclaimed()` in `internal/docker` — never inflate when diff ≤ 0.
2. Extend `PruneReport` JSON with before/after df and `noChange`.
3. Replace Docker success `alert()` with status banner; call `refreshInsightsAfterMaintenance()`.
4. Expand modal hint for Windows VHDX behavior.

## Risks

Low — behavior change is more accurate; users who saw inflated reclaim will now see 0/noChange.

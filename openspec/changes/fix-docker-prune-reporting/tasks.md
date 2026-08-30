# Tasks: fix-docker-prune-reporting

**Change:** `fix-docker-prune-reporting`

## Implementation Checklist

- [x] `computeReclaimed()` + extended `PruneReport` in `internal/docker/detect.go`
- [x] Unit tests in `internal/docker/reclaim_test.go`
- [x] UI status banner + `noChange` messaging in `web/src/main.ts`
- [x] `refreshInsightsAfterMaintenance()` after prune
- [x] Docker modal copy for VHDX / WSL compact pointer
- [x] Delta spec + testing report

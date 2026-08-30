# Testing report: fix-docker-prune-reporting

**Date:** 2026-08-30

| Check | Result | Notes |
|-------|--------|-------|
| Unit — computeReclaimed | pass | `internal/docker/reclaim_test.go` |
| Unit — docker package | pass | `go test ./internal/docker/...` |
| API — prune dry-run | pass | `handlers_docker_test.go` |
| Web build | pass | `npm run build` in `web/` |

## Manual verification

- [ ] Prune when nothing reclaimable → UI shows 0 / noChange banner
- [ ] Prune when df drops → reclaimed matches diff; insights refresh

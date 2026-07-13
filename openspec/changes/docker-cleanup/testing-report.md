# Testing report: docker-cleanup

**Change:** docker-cleanup  
**Date:** 2026-07-13  
**Issue:** [#77](https://github.com/jcuel/disk-tool/issues/77)

## Automated

- [x] `go test ./internal/docker/ ./internal/insights/ ./internal/safety/ ./internal/api/` — pass
- [x] `node web/node_modules/vite/bin/vite.js build` — pass
- [x] Embedded `web/dist` → `cmd/disk-tool/static`
- [x] API: prune rejects without confirm; dry-run returns OK without deleting

## Unit coverage

- `internal/docker/detect_test.go` — `docker system df` parse fixtures, size tokens, synthetic path protection, candidate paths
- `internal/api/handlers_docker_test.go` — confirm gate, dry-run OK, missing scan 404
- Insights tests use `SkipDocker` so host Docker/Desktop dirs do not pollute fixture counts

## Manual (recommended)

- [ ] With Docker daemon running: scan → insights show CategoryDocker → **Docker reclaim** → dry-run → type DELETE → prune
- [ ] Without Docker CLI: caution data-root candidates (`deletable: false`); preset opens install/start hint
- [ ] Confirm filesystem cleanup skips `docker://` and Docker data roots
- [ ] Confirm VHDX / disk images remain non-deletable (protect-disk-images)

## Notes

- Prune uses `docker system prune -af` only — **no `--volumes`**
- Reclaim requires `docker` on PATH and a running daemon
- CI does not require a Docker daemon; unit fixtures cover parsing

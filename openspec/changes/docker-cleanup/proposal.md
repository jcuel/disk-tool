# Proposal: Docker maintenance module

**Change:** docker-cleanup
**Status:** in progress
**Domain:** disk-tool

## Summary

Add CLI-first Docker maintenance: detect reclaimable usage, surface insights + `docker-reclaim` preset, execute `docker system prune -af` after dry-run and typed DELETE. Never delete Docker data roots or VHDX files.

## Deliverables

- `internal/docker` package (detect, df parse, prune, path awareness)
- `CategoryDocker` insights + preset
- `GET/POST /api/scans/{id}/docker` endpoints
- UI maintenance preset with review/confirm modal
- OpenSpec delta + tests

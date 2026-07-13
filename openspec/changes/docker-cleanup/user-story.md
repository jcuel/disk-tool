# User Story: Docker maintenance (CLI prune + insights)

**Change:** `docker-cleanup`
**Status:** refined

## Story

As a developer with a full disk from Docker images and build cache, I want disk-tool to detect reclaimable Docker usage and run a confirmed CLI prune, so I can free space without deleting Docker Desktop / WSL disk images or named volumes.

## Context

Docker cleanup was deferred from guided-bulk-cleanup. VHDX/WSL disks are already non-deletable. Filesystem delete of `/var/lib/docker` or Desktop data roots is unsafe.

## Acceptance Criteria

- [x] `internal/docker` detects Docker CLI and parses reclaimable usage (`docker system df`)
- [x] Insights show `CategoryDocker` candidates; Docker data roots are not filesystem-deletable
- [x] Preset `docker-reclaim` + confirmed prune (no `--volumes`)
- [x] VHDX / Docker Desktop disk images remain non-deletable
- [x] OpenSpec package + tests + UI dry-run/confirm flow
- [x] Docs note: Docker CLI must be on PATH / daemon running for reclaim

## Out of Scope

- `docker volume prune` / `--volumes`
- WSL Optimize-VHD / compact of docker-desktop-data.vhdx
- Remote Docker contexts
- Filesystem delete of `/var/lib/docker` or Desktop data directories

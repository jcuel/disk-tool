# User Story: WSL / Docker VHDX compact

**Change:** `wsl-vhdx-compact`
**Status:** refined

## Story

As a Windows user whose Docker prune did not free host disk space, I want to compact WSL / Docker Desktop VHDX files, so NTFS file size actually shrinks after layer cleanup.

## Context

Docker prune removes layers inside VHDX; `%LOCALAPPDATA%\Docker\wsl\*.vhdx` often stays the same size on disk. Previously out of scope in `docker-cleanup`.

## Acceptance Criteria

- [x] Detect VHDX candidates (Docker wsl + WSL2 distro paths)
- [x] `GET /api/wsl/disks` lists path, size, modified
- [x] `POST /api/wsl/compact` with dry-run and DELETE confirm gate
- [x] Workflow: `wsl --shutdown`, then `Optimize-VHD` or diskpart fallback
- [x] Return `bytesBefore`, `bytesAfter`, `freedOnDisk`
- [x] Never delete VHDX — compact only
- [x] Windows UI maintenance card; stub on other OS

## Out of Scope

- macOS / Linux compact
- Remote Docker contexts

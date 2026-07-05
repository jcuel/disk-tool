# Proposal: Foundation scanner API through full v1 UI

**Change:** `foundation-scanner-api`
**Status:** archived
**Domain:** `disk-tool`

## Summary

Build a cross-platform local web disk analyzer: Go parallel scanner, REST/WebSocket API, Vite UI with tree table, ECharts treemap/bar charts, largest-files panel, JSON/HTML export, and Docker smoke verification.

## Motivation

Users need TreeSize-like visibility into disk usage without proprietary desktop software. A localhost web app works on Windows, Linux, and macOS with one codebase.

## Scope

### In scope

- Go scanner with concurrent walk, cancel, progress events, symlink cycle protection
- REST + WebSocket API on 127.0.0.1
- Vite + TypeScript + ECharts frontend
- CLI: `serve`, `scan --json`, `version`
- JSON and HTML export
- Dockerfile + docker-compose smoke test

### Out of scope

- Duplicate finder, remote mounts policy beyond user-selected path
- File deletion, PDF export, SQLite spill for huge trees

## Risks

| Risk | Mitigation |
|------|------------|
| Memory on huge trees | Top-N largest files capped at 100; full tree in memory for v1 |
| Symlink cycles | Dev/inode visit tracking |
| Path traversal | Validate all paths stay under scan root |

# Proposal: OS safety zones

**Change:** safety-zones
**Status:** in-progress
**Domain:** disk-tool

## Summary

Classify paths into OS-aware safety zones, skip forbidden paths during scan, block delete on protected zones, and surface a Safety Grid in the UI.

## Scope

- `internal/safety/zones.go` — Windows and Linux path catalogs
- Scanner skip for `forbidden`
- Delete/cleanup guards for protected zones
- Safety grid + maintenance presets API and UI

## Out of scope

- macOS zone catalog (stub for v2)
- Auto-delete in diagnostic zones

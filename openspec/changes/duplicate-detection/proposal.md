# Proposal: Duplicate file detection

**Change:** duplicate-detection
**Status:** draft
**Domain:** disk-tool

## Summary

Add content-hash duplicate detection scoped to the active scan root, surfaced in insights and exports.

## Motivation

Users often have identical copies across Downloads, backups, and project folders. Finding them manually is slow.

## Scope

### In scope

- Hash files during or after branch scan (size threshold configurable, default skip files <4KB)
- Duplicate groups in insights panel
- API field on scan job

### Out of scope

- Auto-delete duplicates
- Global cross-volume scan without user-selected root

## Risks

| Risk | Mitigation |
|------|------------|
| Long scan on large trees | Run on demand; cap file count per group display |
| Memory use | Stream hashing; don't load all paths at once |

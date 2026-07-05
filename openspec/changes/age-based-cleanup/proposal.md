# Proposal: Age-based cleanup suggestions

**Change:** age-based-cleanup
**Status:** draft
**Domain:** disk-tool

## Summary

Extend insights to surface large files not modified within a configurable age threshold.

## Motivation

Downloads and temp directories accumulate stale installers and archives users forget to remove.

## Scope

### In scope

- Age + size filters during insights pass
- UI controls for threshold days and min size
- New cleanup candidate category `stale-large-file`

### Out of scope

- macOS/iOS backup folders requiring special permissions
- Deletion (use safe-delete flow separately)

## Risks

| Risk | Mitigation |
|------|------------|
| False positives on intentionally archived files | Conservative defaults; clear hints |

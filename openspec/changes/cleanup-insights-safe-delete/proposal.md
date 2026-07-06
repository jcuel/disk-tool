# Proposal: Safe delete cleanup candidates

**Change:** cleanup-insights-safe-delete
**Status:** archived

## Summary

Allow deleting cleanup candidate paths from the UI after explicit confirmation. Paths must be under the active scan root; scan root itself cannot be deleted.

## Scope

- POST `/api/scans/{id}/delete` with `{ path, confirm: true }`
- Delete button on cleanup candidates table only

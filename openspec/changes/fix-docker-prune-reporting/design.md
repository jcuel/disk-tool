# Design: fix-docker-prune-reporting

**Change:** `fix-docker-prune-reporting`

## Reclaim calculation

```
reclaimed = before.Reclaimable - after.Reclaimable  if before > after else 0
noChange  = reclaimed == 0 && prune command succeeded
```

Never substitute pre-prune estimate when diff ≤ 0.

## API

Existing `POST /api/scans/{id}/docker/prune` response gains optional fields:

- `beforeReclaimable`, `afterReclaimable` (int64)
- `beforeDf`, `afterDf` (raw CLI output)
- `noChange` (bool)

## UI

Mirror delete-fix pattern: `#status-banner` for prune outcome; reanalyze insights on success.

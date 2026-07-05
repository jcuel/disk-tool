# Design: Age-based cleanup suggestions

**Change:** age-based-cleanup

## Approach

1. Extend `internal/insights` to accept `AgeThresholdDays` and `MinSizeBytes`.
2. During scan/largest-files collection, evaluate mtime against threshold.
3. Prefer paths under `Downloads`, `Temp`, `tmp` but include any large stale file.
4. UI: two number inputs in insights panel (persist in sessionStorage).

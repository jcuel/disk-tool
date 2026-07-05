# Design: Guided bulk cleanup

## Flow

1. User selects candidates (default preset: `risk=review`).
2. `POST cleanup` with `dryRun: true` — validate paths, probe locks, return manifest.
3. User confirms via checkbox + typed `DELETE`.
4. `POST cleanup` with `dryRun: false`, `confirm: true`, `confirmPhrase: "DELETE"`.
5. Process largest paths first; skip locked/missing; store `LastCleanupReport`.
6. Prune deleted paths from in-memory insights; client refreshes job.

## Safety

- `PathWithinRoot` on every path; reject scan root.
- Best-effort lock probe on files; directory lock detected on delete failure.
- No process termination in v1.

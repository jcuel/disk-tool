## ADDED Requirements

### Guided bulk cleanup

The system SHALL expose `POST /api/scans/{id}/cleanup` accepting a list of paths under the scan root.

When `dryRun` is true, the system SHALL validate each path and return a cleanup report without deleting files.

When `dryRun` is false, the system SHALL require `confirm: true` and `confirmPhrase: "DELETE"` before deleting.

The system SHALL skip paths that are missing, outside the scan root, locked/in use, or equal to the scan root, recording the reason per item.

The system SHALL store the last cleanup report on the scan job and support export via `format=cleanup-json`, `cleanup-html`, and `cleanup-ticket`.

# User Story: Guided bulk cleanup

**Change:** guided-bulk-cleanup
**Status:** refined

## Story

As a user, I want to select multiple cleanup candidates, review them with a dry-run, confirm twice, and delete them in one session with a report of what was removed or skipped.

## Acceptance Criteria

- [ ] Multi-select cleanup candidates in the UI
- [ ] Dry-run preflight with per-path status (locked, missing, outside root)
- [ ] Execute requires confirm + typed DELETE
- [ ] Cleanup report with bytes reclaimed and export (json/html/ticket)
- [ ] Skip locked/in-use paths without stopping the batch

## Out of Scope

- Docker/container cleanup (phase 2)
- Automatic deletion without confirmation

# User Story: Age-based cleanup suggestions

**Change:** age-based-cleanup
**Status:** refined

## Story

As a user, I want disk-tool to highlight large old files (especially in Downloads and temp folders) so I can remove stale data.

## Acceptance Criteria

- [ ] Track last-modified time during scan (already available via os.Stat)
- [ ] Flag files older than N days (default 90) above size threshold (default 50MB)
- [ ] Show age-based candidates in insights with path, size, age, hint
- [ ] Configurable thresholds in UI (days + min size)
- [ ] Include in support ticket export

## Out of Scope

- Access-time/atime (platform unreliable)
- Automatic deletion

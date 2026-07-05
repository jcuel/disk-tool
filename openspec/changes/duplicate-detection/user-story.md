# User Story: Duplicate file detection

**Change:** duplicate-detection
**Status:** refined

## Story

As a user, I want disk-tool to find duplicate files by content hash so I can reclaim space from redundant copies.

## Acceptance Criteria

- [ ] Optional duplicate scan phase after folder drill-down (or dedicated action)
- [ ] Group duplicates by hash with total reclaimable size
- [ ] Show paths in each duplicate group in UI
- [ ] Limit hash scan to paths under active scan root
- [ ] Export duplicate groups in JSON/HTML report

## Out of Scope

- Automatic deletion of duplicates (manual review only)
- Cross-drive deduplication without explicit scan root

# User Story: Cleanup insights and support reports

**Change:** `cleanup-insights`
**Status:** refined

## Story

As a user, I want disk-tool to highlight where space is consumed and flag common leftovers (node_modules, caches, Downloads installers) so I can clean up safely and attach a report to a support ticket.

## Acceptance Criteria

- [x] Insights summary after overview scan (top consumer + reclaimable estimate)
- [x] Cleanup candidates table with category, path, size, hint
- [x] Copy report / export support ticket (plain text)
- [x] HTML export includes cleanup section
- [ ] One-click open folder in explorer (future)
- [ ] Safe delete with confirmation (future, v2)

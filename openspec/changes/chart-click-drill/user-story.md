# User Story: Chart click drill-down and layout balance

**Change:** chart-click-drill
**Status:** refined

## Story

As a user exploring disk usage, I want to click folders in the Distribution charts and see a balanced layout so that drill-down matches the folder tree and panels use the full screen width.

## Acceptance Criteria

- [x] Clicking a treemap or bar segment selects that folder path
- [x] Expandable folders (`+`) auto-scan when clicked from a chart
- [x] Tree and Distribution columns are equal width; tree is not capped at 70vh
- [x] Insights panel spans full width on the bottom row
- [x] Empty tree/chart panels show scan-in-progress messaging during overview scan

## Out of Scope

- Partial tree updates during overview scan (snapshot streaming)
- Chart click on disk pie (capacity summary)

# Proposal: Chart click drill-down and layout balance

**Change:** chart-click-drill
**Status:** archived

## Summary

Fix Distribution chart clicks so they drill into folders like tree rows. Balance the main grid: equal tree/chart columns, full-width Insights, scan placeholders during overview.

## Scope

- `web/src/charts.ts` — click path resolution, ECharts `triggerEvent`, `notMerge`
- `web/src/main.ts` — scan placeholders, chart `needsExpand` parity with tree
- `web/src/styles.css` — grid column balance, remove tree `70vh` cap, Insights span
- Embedded static assets rebuild

## Cause (chart clicks)

ECharts `setOption` merge dropped custom `path` on series data; bar series lacked reliable click events without `triggerEvent`.

# Chart click drill-down

## Problem
Clicking folders in the Distribution treemap or bar chart did not select the path or trigger expand/scan.

## Cause
ECharts `setOption` merge dropped custom `path` fields on series data; bar series lacked `triggerEvent`, so clicks were not emitted reliably.

## Fix
- Full option replace (`notMerge: true`) on chart updates
- `triggerEvent: true` on treemap and bar series (and bar y-axis)
- `resolveChildPath` helper with fallbacks by `data.path`, label name, and `dataIndex`
- Pointer cursor on chart containers

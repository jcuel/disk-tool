# Design: chart-click-drill

**Change:** chart-click-drill

## Chart click resolution

ECharts merges series data on `setOption` by default, dropping custom `path` fields. Use `setOption(option, { notMerge: true })` and `triggerEvent: true` on treemap, bar series, and bar y-axis.

`pickChildPath` resolves folder path in order: `data.path`, category label / name map, `dataIndex`, then `convertFromPixel` fallback.

Chart callback mirrors tree rows: `selectPath(path, needsExpand(child))`.

## Layout

Grid areas:

```
tree   charts
tree   files
insights (span 2 columns)
```

Remove `grid-column: auto` on `.insights-panel` (broke full-width span). Remove `max-height: 70vh` on `.tree-panel`. Use `minmax(0, 1fr)` columns so both sides stay 50/50.

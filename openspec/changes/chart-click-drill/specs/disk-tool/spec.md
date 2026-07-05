## MODIFIED Requirements

### Requirement: Visual charts

The system SHALL render a treemap and bar chart for the selected folder's immediate children.

#### Scenario: Drill-down from chart

- **WHEN** the user clicks a folder segment in the treemap or bar chart
- **THEN** that folder becomes selected and expandable children trigger a deeper scan

#### Scenario: Drill-down from tree

- **WHEN** the user clicks a folder in the tree or chart
- **THEN** charts update to show that folder's children

## ADDED Requirements

### Requirement: Layout balance

The web UI SHALL use a two-column grid for tree vs charts/files with equal column width, SHALL span Insights across the full width below, and SHALL show scan-in-progress placeholders when overview data is not yet available.

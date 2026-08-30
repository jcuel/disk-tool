# Delta Spec: Desktop tabbed UI

**Change:** `desktop-tabbed-ui`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: Desktop tabbed layout

When running inside the Tauri desktop shell, the UI SHALL present a tab bar with **Browse**, **Files**, and **Insights** panes. Browse SHALL contain the folder tree and distribution charts; Files SHALL contain the largest-files table; Insights SHALL contain maintenance and cleanup panels. Header, controls, and disk summary SHALL remain fixed; only the active tab pane SHALL scroll. Switching to Browse SHALL resize charts to fit the visible viewport. The browser-hosted UI SHALL retain the existing single-scroll layout without tabs.

#### Scenario: Desktop tabs visible

- **WHEN** the UI loads in the Tauri webview
- **THEN** a tab bar is shown and exactly one tab pane is active

#### Scenario: Browser unchanged

- **WHEN** the UI loads in a normal browser without Tauri
- **THEN** no tab bar is shown and the 2×2 scroll layout is preserved

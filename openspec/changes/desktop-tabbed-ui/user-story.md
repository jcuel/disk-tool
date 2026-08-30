# User Story: Desktop tabbed UI

**Change:** `desktop-tabbed-ui`
**Status:** refined

## Story

As a desktop disk-tool user, I want Browse / Files / Insights tabs in the Tauri app, so I can navigate a large scan without endless scrolling while the browser demo keeps its current layout.

## Context

Extends `desktop-shell`. Single long scroll works in browser; desktop window benefits from fixed chrome and tab panes with chart resize on tab switch.

## Acceptance Criteria

- [x] Tauri detected via `__TAURI_INTERNALS__`; `body.desktop-app` class applied
- [x] Tabs: Browse (tree + charts), Files (largest files), Insights (maintenance + cleanup)
- [x] Fixed header/controls/disk summary; only tab pane scrolls
- [x] Charts resize when switching to Browse tab
- [x] Browser layout unchanged (no tab bar without Tauri)
- [x] Default Tauri window height 900px

## Out of Scope

- Browser tabbed layout
- System tray / single-instance

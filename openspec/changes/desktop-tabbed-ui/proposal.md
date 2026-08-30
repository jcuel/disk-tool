# Proposal: desktop-tabbed-ui

**Change:** `desktop-tabbed-ui`
**Status:** applying
**Issue:** Tabbed layout for Tauri shell

## Summary

Add desktop-only tab navigation preserving tree↔chart coupling on Browse tab. Browser users see unchanged scroll layout.

## Approach

- `web/src/desktop-tabs.ts` — detection, tab bar, pane visibility
- `data-desktop-tab` on panels; CSS flex viewport under `body.desktop-app`
- Export `resizeCharts()` from charts module

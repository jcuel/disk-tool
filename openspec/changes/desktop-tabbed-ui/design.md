# Design: desktop-tabbed-ui

**Change:** `desktop-tabbed-ui`

## Tab model (Option A)

| Tab | Panels |
|-----|--------|
| Browse | tree + distribution charts |
| Files | largest files table |
| Insights | maintenance hub + cleanup + duplicates |

Tab bar inserted before `main.layout`. Panels tagged with `data-desktop-tab`; CSS shows only active pane in desktop mode.

## Detection

`window.__TAURI_INTERNALS__` — no extra npm dep required for v1.

## Layout CSS

```css
body.desktop-app { height: 100vh; overflow: hidden; }
body.desktop-app .tab-pane { overflow: auto; }
```

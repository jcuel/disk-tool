# Design: os-maintenance-hub

**Change:** `os-maintenance-hub`

## Recycle bin adapters

| OS | Mechanism |
|----|-----------|
| Windows | PowerShell `Clear-RecycleBin -Force` |
| Linux | `gio trash --empty` or walk `~/.local/share/Trash/files` |
| macOS | AppleScript empty trash via Finder |

## API

- `GET /api/maintenance/recycle` — `{ supported, itemCount, totalBytes }`
- `POST /api/maintenance/recycle/empty` — `{ dryRun, confirm, confirmPhrase }`

Not tied to scan ID.

## UI cards

Temp cleanup, empty recycle bin, Docker reclaim, critical OS explainer; WSL card added by `wsl-vhdx-compact` on Windows.

# Design: wsl-vhdx-compact

**Change:** `wsl-vhdx-compact`

## VHDX discovery

- `%LOCALAPPDATA%\Docker\wsl\*\*.vhdx`
- `%LOCALAPPDATA%\Packages\*\LocalState\ext4.vhdx`

## Compact sequence

1. User confirms (warn: stops all WSL distros)
2. `wsl --shutdown`
3. `Optimize-VHD -Mode Full` if Hyper-V module available, else diskpart compact script
4. Re-stat file; return before/after bytes

## Safety

VHDX paths remain in disk-image protection set; compact never deletes.

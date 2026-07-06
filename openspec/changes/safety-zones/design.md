# Design: OS safety zones

**Change:** safety-zones

## Zone tiers

| Zone | Scan | Delete |
|------|------|--------|
| forbidden | skip | never |
| critical_os | show | never |
| diagnostic | show | never |
| maintenance, review, caution | show | review flow |

## Packages

- `internal/safety/zones.go` — path classification
- `internal/safety/presets.go` — maintenance preset matching
- Hooks in scanner, insights, cleanup, delete handlers

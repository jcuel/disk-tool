## ADDED Requirements

### Disk capacity endpoint

The system SHALL expose `GET /api/disk?path=...` returning volume total, used, and free bytes.

### Disk summary UI

The web UI SHALL display volume capacity stats and a used/free chart at the top, and SHALL auto-start an overview scan on the default system drive when loaded.

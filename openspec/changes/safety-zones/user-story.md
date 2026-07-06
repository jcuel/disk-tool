# User story: Safety zones

As a user scanning my full drive, I want disk-tool to protect OS and crash-dump paths so I cannot accidentally delete files that would break my system.

## Acceptance criteria

- [x] Scanning skips `/proc`, `/sys`, `System32`, etc.
- [x] Delete and bulk cleanup reject protected zones
- [x] Safety grid shows zone breakdown
- [x] Full-drive scan shows guidance banner
- [x] Maintenance presets offer safe reclaim paths only

# Delta Spec: WSL VHDX compact

**Change:** `wsl-vhdx-compact`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: WSL VHDX compact (Windows)

On Windows, the system SHALL detect compactable WSL and Docker Desktop VHDX files, expose list and compact APIs, and report on-disk file size before and after compact. Compact SHALL shut down WSL, run `Optimize-VHD` or an equivalent diskpart compact fallback, and SHALL NOT delete VHDX files. Execute SHALL require `confirm: true` and `confirmPhrase: "DELETE"`.

#### Scenario: List VHDX disks

- **WHEN** `GET /api/wsl/disks` on Windows with Docker/WSL VHDX present
- **THEN** the response lists paths with current file sizes

#### Scenario: Compact with confirm

- **WHEN** the user confirms compact for a listed VHDX path
- **THEN** WSL is shut down, the VHDX is compacted, and `freedOnDisk` reflects NTFS file size change

#### Scenario: Non-Windows stub

- **WHEN** `GET /api/wsl/disks` on Linux or macOS
- **THEN** the response indicates WSL compact is not supported

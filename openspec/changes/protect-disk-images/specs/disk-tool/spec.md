# Delta Spec: Protect disk images and confirm single deletes

**Change:** `protect-disk-images`
**Domain:** `disk-tool`

## ADDED Requirements

### Requirement: Virtual disk image protection

The system SHALL treat virtual disk and disk-image file extensions (including `.vhd`, `.vhdx`, `.avhd`, `.avhdx`, `.vmdk`, `.vdi`, `.qcow`, `.qcow2`, `.wim`, `.esd`, `.vfd`) as non-deletable. Single delete and bulk cleanup SHALL reject those paths. Cleanup insights SHALL NOT list them as download or stale-large candidates. Largest Files MAY still list them and SHALL mark them non-deletable in the UI.

#### Scenario: Delete VHDX blocked

- **WHEN** the client requests delete or cleanup execute for a `.vhdx` path under the scan root
- **THEN** the server rejects the operation with a disk-image protection reason

#### Scenario: Downloads VHD not a cleanup candidate

- **WHEN** insights analyze a Downloads folder containing a large `.vhd` file
- **THEN** that file is not included in cleanup candidates

#### Scenario: Largest Files shows protected disk image

- **WHEN** a scan's largest files include a `.vhdx`
- **THEN** the scan response marks that entry `deletable: false`
- **AND** the UI disables Delete for that row

### Requirement: Single-path delete confirmation workflow

The web UI SHALL NOT use a browser prompt for single-path delete from Largest Files or cleanup candidate rows. Delete SHALL open a modal review step, run a cleanup dry-run preflight, then require an explicit checkbox and typed `DELETE` before execute.

#### Scenario: Largest Files delete confirmation

- **WHEN** the user clicks Delete on a deletable Largest Files row
- **THEN** a Review delete modal is shown
- **AND** Continue runs a dry-run
- **AND** Confirm delete requires review checkbox and phrase `DELETE` before permanent deletion

## MODIFIED Requirements

### Requirement: Safe delete

The system SHALL expose `POST /api/scans/{id}/delete` requiring `confirm: true` and `confirmPhrase: "DELETE"`. Paths outside the scan root, the scan root itself, protected safety zones, and virtual disk / disk-image extensions SHALL be rejected.

### Requirement: Guided bulk cleanup

The system SHALL expose `POST /api/scans/{id}/cleanup` with dry-run and execute modes. Execute SHALL require typed confirmation and skip locked, missing, outside-root, protected-zone, and virtual disk / disk-image paths.

### Requirement: Largest files

The system SHALL list the top 100 largest files found during a completed scan. Each entry MAY include `deletable` indicating whether the UI may offer delete. Virtual disk / disk-image paths SHALL be marked non-deletable.

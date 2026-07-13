# User Story: Protect disk images and confirm single deletes

**Change:** `protect-disk-images`
**Status:** refined

## Story

As a user scanning a full drive, I want virtual disk files (VHD/VHDX and similar) blocked from delete and every single-file delete to go through an explicit review workflow, so I cannot accidentally destroy VMs or wipe a large file with a one-click prompt.

## Context

Largest Files used a browser `prompt()` for confirmation. Cleanup insights treated any Downloads file ≥50 MB as a candidate regardless of extension, so Hyper-V disks appeared as installers. Safe-delete already blocked OS zones but not disk-image extensions.

## Acceptance Criteria

- [x] Virtual disk / disk-image extensions cannot be deleted via single delete or bulk cleanup
- [x] Those extensions are excluded from download and stale cleanup candidates
- [x] Largest Files still lists disk images but Delete is disabled
- [x] Single-file Delete (Largest Files and cleanup row) uses review → dry-run → confirm modal (checkbox + type DELETE)
- [x] Unit tests cover disk-image detection and non-deletable paths

## Out of Scope

- Blocking `.iso` (still allowed as download installer candidate)
- Soft-delete / Recycle Bin
- Product-level plugin system for custom protected extensions

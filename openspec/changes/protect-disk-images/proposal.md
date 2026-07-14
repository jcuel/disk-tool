# Proposal: Protect virtual disk images from delete

**Change:** protect-disk-images
**Status:** archived

## Problem

Largest Files and cleanup insights can offer `.vhd` / `.vhdx` (and similar) for deletion. Any Downloads file ≥50MB becomes a cleanup candidate regardless of extension, so Hyper-V disks appear as "safe to review" installers. Accidental delete can destroy VMs. Single-file delete also used a browser `prompt()`, which is easy to miss.

## Solution

- Block delete of virtual disk / disk-image extensions server-side (`CanDeletePath`)
- Exclude them from download/stale cleanup candidates
- Disable Delete in Largest Files UI (still visible for awareness via `deletable: false`)
- Replace browser prompt with review → dry-run → confirm modal (same gates as bulk cleanup)

## Artifacts

- `user-story.md`, `design.md`, `tasks.md`, `testing-report.md`
- Delta: `specs/disk-tool/spec.md`

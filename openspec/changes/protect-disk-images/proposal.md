# Proposal: Protect virtual disk images from delete

**Change:** protect-disk-images
**Status:** in progress

## Problem

Largest Files and cleanup insights can offer `.vhd` / `.vhdx` (and similar) for deletion. Any Downloads file ≥50MB becomes a cleanup candidate regardless of extension, so Hyper-V disks appear as "safe to review" installers. Accidental delete can destroy VMs.

## Solution

- Block delete of virtual disk / disk-image extensions server-side
- Exclude them from cleanup candidates
- Disable Delete in Largest Files UI (still visible for awareness)

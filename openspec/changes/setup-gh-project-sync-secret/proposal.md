# Proposal: Classic PAT setup for GH_PROJECT_SYNC

**Change:** setup-gh-project-sync-secret
**Status:** in-progress

## Summary

Add cross-platform scripts to rotate `GH_PROJECT_SYNC` with a classic PAT. User-owned project boards require classic `project` scope; fine-grained tokens are unsupported.

## Scope

- `scripts/setup-gh-project-sync-secret.sh` and `.ps1`
- Update `.github/PROJECT.md` with classic PAT requirement and script usage

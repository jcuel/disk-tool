# Proposal: Outreach automation

**Change:** outreach-automation
**Status:** in progress

## Summary

Automate Show HN and Reddit launch posts after release binaries ship, using repo-side scripts, GitHub Actions, and optional Cursor Automation re-runs.

## Deliverables

- `config/outreach.yaml` + templates under `openspec/changes/product-launch/templates/`
- `scripts/outreach/` — render, HN submit, Reddit submit, orchestrator
- `.github/workflows/outreach-on-release.yml`
- `scripts/release-version.sh` — strip Cursor footer from version parsing
- Cursor Automation setup guide

## Risks

- HN has no official write API; session login may be rate-limited or flagged.
- Reddit subs may auto-remove self-promo; karma/age gates apply.
- PolyForm NC license line must appear in every post (enforced via config footer).

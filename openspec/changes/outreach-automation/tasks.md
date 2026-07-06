# Tasks: Outreach automation

- [x] Add `config/outreach.yaml` and post templates
- [x] Add `scripts/outreach/` (render, hn_submit, reddit_post, run)
- [x] Add `.github/workflows/outreach-on-release.yml`
- [x] Fix `scripts/release-version.sh` Cursor footer parsing
- [x] Document in CONTRIBUTING.md and launch-checklist.md
- [x] Add Cursor Automation setup guide
- [ ] Maintainer: set GitHub secrets and `OUTREACH_ENABLED`
- [ ] Maintainer: install hn-mcp and create Cursor Automation
- [ ] Validate with `gh workflow run outreach-on-release.yml -f version=1.3.0 -f dry_run=true`

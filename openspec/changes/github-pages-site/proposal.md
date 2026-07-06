# Proposal: GitHub Pages product site

**Change:** github-pages-site
**Status:** in progress

## Summary

Publish a GitHub Pages site at `https://jcuel.github.io/disk-tool/` with a product landing page and an interactive demo that reuses the disk-tool UI with static fixtures (no Go backend).

## Scope

- `site/` — Vite landing page
- `web/src/demo/` — fixtures + mock API for `VITE_DEMO_MODE`
- `scripts/export-demo-fixtures.sh`
- `.github/workflows/pages.yml`
- CI smoke job for site + demo builds

## Out of scope

- Custom domain, release binaries on Pages, real filesystem scanning on Pages

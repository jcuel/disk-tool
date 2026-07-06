# Design: Outreach automation

## Flow

1. Tag `vX.Y.Z` → `release-assets.yml` uploads binaries.
2. `outreach-on-release.yml` runs on `workflow_run` success when `OUTREACH_ENABLED=true`.
3. `run.py` checks release notes for `outreach-posted:` marker, renders templates, posts HN + Reddit (staggered).
4. Release notes updated with marker and summary.

## Components

| Layer | Role |
|-------|------|
| `config/outreach.yaml` | URLs, channels, stagger, license footer |
| `render.py` | `{{version}}`, `{{demo_url}}`, etc. |
| `hn_submit.py` | Session CSRF login + link submit |
| `reddit_post.py` | OAuth2 script app + `/api/submit` |
| `run.py` | Idempotency + orchestration |

## Credentials

- **GitHub Actions secrets** — primary path for release-triggered posts.
- **Cursor MCP env** — hn-mcp for manual/tag-triggered agent re-runs (see `cursor-automation-setup.md`).

## Guards

- `OUTREACH_ENABLED` repo variable (default off).
- `workflow_dispatch` with `dry_run=true` for testing.
- Optional `OUTREACH_DRY_RUN=true` variable for tag-triggered dry runs.

# Cursor Automation setup (outreach re-runs)

Use this when you want a **tag-triggered or manual** Cursor agent to re-run outreach without waiting for GitHub Actions, or to post via hn-mcp from your machine.

## 1. Install hn-mcp

Clone [hn-mcp](https://github.com/booklib-ai/hn-mcp) and add to Cursor MCP settings (`~/.cursor/mcp.json` or project config):

```json
{
  "mcpServers": {
    "hackernews": {
      "command": "uv",
      "args": ["run", "--directory", "/path/to/hn-mcp", "python3", "server.py"],
      "env": {
        "HN_USERNAME": "your_hn_username",
        "HN_PASSWORD": "your_hn_password"
      }
    }
  }
}
```

Authenticate hn-mcp in the **Agents Window** before saving a Cursor Automation.

## 2. Reddit credentials

Set the same Reddit OAuth env vars used by GitHub Actions (locally or in Cursor cloud agent secrets):

- `REDDIT_CLIENT_ID`, `REDDIT_CLIENT_SECRET`, `REDDIT_USERNAME`, `REDDIT_PASSWORD`
- `REDDIT_USER_AGENT` (optional)

## 3. Create the automation

In Cursor → Automations → New:

| Field | Value |
|-------|--------|
| **Name** | disk-tool release outreach |
| **Trigger** | Git push tag matching `v*.*.*` on `jcuel/disk-tool` (or manual) |
| **Tools** | MCP (hackernews / hn-mcp), shell |
| **Instructions** | See below |

### Agent instructions (paste)

```
When triggered for tag vX.Y.Z on jcuel/disk-tool:

1. Read config/outreach.yaml and openspec/changes/product-launch/templates/.
2. Run: python scripts/outreach/run.py --dry-run --version X.Y.Z
   Review rendered titles/bodies; ensure PolyForm NC license footer is present.
3. If OUTREACH_ENABLED is not true and this is not an explicit manual run, stop.
4. HN: use hn_submit MCP tool with title/url from config (Show HN link to demo URL).
5. Reddit: run python scripts/outreach/reddit_post.py per subreddit from config,
   or python scripts/outreach/run.py --version X.Y.Z --channel reddit
6. Skip channels already marked in GitHub Release notes (outreach-posted: hn,reddit).
7. Log post URLs in the automation output.
```

## 4. GitHub Actions (primary)

For unattended release posts, prefer the workflow:

1. Set secrets: `HN_*`, `REDDIT_*`
2. Set variable `OUTREACH_ENABLED=true`
3. Dry-run: `gh workflow run outreach-on-release.yml -f version=1.3.0 -f dry_run=true`

## 5. First live test

Cut a patch release (`v1.3.1`) with `OUTREACH_ENABLED=true` only after a successful dry-run.

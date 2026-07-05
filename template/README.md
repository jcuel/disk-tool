# SPECBOOT + OpenSpec bootstrap

Copy this folder into a new repository root to start with SPECBOOT slash commands, OpenSpec layout, GitHub issue templates, and workflow rules.

## Quick start

```bash
# From your new repo root (after copying template/ contents):
mkdir -p openspec/changes
git init && git add .
```

1. Edit `openspec/config.yaml` — set `project` and `domain`.
2. Edit `openspec/specs/app/spec.md` — rename `app` to your domain folder if needed.
3. Edit `.cursor/rules/github-workflow.mdc` — set your GitHub repo URL and project board link.
4. Run `/enrich-us <change-name>` with your first user story.

## What is included

| Path | Purpose |
|------|---------|
| `.cursor/commands/` | SPECBOOT slash commands (`/enrich-us` … `/commit`) |
| `.cursor/rules/` | Engineering guardrails, workflow state machine, GitHub tracking |
| `.cursor/templates/openspec/` | Artifact templates for changes |
| `openspec/config.yaml` | Project metadata |
| `openspec/specs/app/spec.md` | Source-of-truth spec stub |
| `openspec/changes/` | Empty; one folder per change |
| `.github/ISSUE_TEMPLATE/` | OpenSpec change + bug templates |

## Workflow

```
/enrich-us → /propose → /apply → /verify → /code-review → /archive → /commit
```

Track work in GitHub Issues + Projects (see `.cursor/rules/github-workflow.mdc`).

## Origin

Scaffold from [disk-tool](https://github.com/jcuel/disk-tool). Customize before use.

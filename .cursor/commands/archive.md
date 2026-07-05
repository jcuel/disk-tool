---
description: Merge delta spec into source of truth and mark change archived
argument-hint: <change-name>
---

# Archive Change

The user invoked this command with: $ARGUMENTS

## Goal

Merge approved delta specs into `openspec/specs/` and mark the change as archived.

## Prerequisites

- `/code-review` gate met: review complete, no critical blockers, summary `ready`
- If not, stop and tell the user to run `/code-review <change-name>` first

## Inputs

- Change name from `$ARGUMENTS`
- Delta spec at `openspec/changes/<change-name>/specs/disk-tool/spec.md`
- Source spec at `openspec/specs/disk-tool/spec.md`

## Steps

1. Read the delta spec ADDED/MODIFIED/REMOVED blocks.
2. Merge changes into `openspec/specs/disk-tool/spec.md` (apply additions, update modifications, remove deletions).
3. Set `proposal.md` status to `archived`.
4. Keep the change folder in place as historical record.
5. Summarize what was merged.
6. **GitHub:** after PR merge, move issue to **Done**; post final comment with merge commit; issue closes via `Closes #N` in PR body.

## Outputs

- Updated `openspec/specs/disk-tool/spec.md`
- `proposal.md` status = `archived`

## Done Criteria

- Source spec reflects merged deltas
- Change folder retained with archived status
- Prompt the user: "Run `/commit` to create the git commit."

## Rules

- Follow `engineering-guardrails.mdc`, `specboot-workflow.mdc`, and `github-workflow.mdc`
- Do not commit, push, or delete the change folder

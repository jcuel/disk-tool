---
description: Draft and create a git commit with user approval
argument-hint: [optional commit scope hint]
---

# Commit

The user invoked this command with: $ARGUMENTS

## Goal

Create a git commit after explicit user approval. No push unless explicitly requested.

## Prerequisites

- Changes are staged or ready to stage
- Prefer running after `/archive` when a feature is complete

## Steps

1. Run in parallel: `git status`, `git diff`, `git log -1 --format='%an %ae'`.
2. Analyze staged and unstaged changes; draft a conventional commit message focused on why.
3. Present the proposed message to the user for approval.
4. Stage relevant files (exclude secrets: `.env`, credentials).
5. Commit only when the user explicitly approves.
6. Run `git status` after commit to confirm success.
7. If pre-commit hook fails, fix and create a new commit (do not amend unless allowed).

## Outputs

- Git commit on the current branch

## Done Criteria

- Commit created with user approval
- Remind user: push and run cloudrabbit tests when CI is configured (only push if asked)

## Rules

- Follow user git safety rules: no force push, no config changes, no push unless asked
- Never commit `.env` or credential files

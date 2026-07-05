---
description: Implement a change — branch, code, tests, docs, testing report
argument-hint: <change-name>
---

# Apply Change

The user invoked this command with: $ARGUMENTS

## Goal

Execute the implementation loop for `openspec/changes/<change-name>/`.

## Prerequisites

- `proposal.md`, `design.md`, `tasks.md`, and delta `specs/app/spec.md` must exist
- If missing, stop and tell the user to run `/propose <change-name>` first

## Inputs

- Change name from `$ARGUMENTS`
- Artifacts in `openspec/changes/<change-name>/`

## Steps

1. Create feature branch `feat/<change-name>` (or switch if it already exists).
2. Confirm a **GitHub Issue** exists for this change; move project Status to **In progress**; comment with branch name (see `github-workflow.mdc`).
3. Read `tasks.md` and implement the next unchecked item.
3. Write or update code with minimal, focused diffs.
4. Add or update tests for the change.
5. Update relevant documentation.
6. Create or update `testing-report.md` (initial test results).
7. If scope shifted during implementation, update `proposal.md` (set status to `in-progress`).
8. Check off completed items in `tasks.md`.
9. Ask whether to continue the implement loop or proceed to `/verify <change-name>`.

## Implement Loop

Repeat steps 2–8 until all tasks are checked or the user stops.

## Outputs

- Feature branch with code and tests
- Updated `tasks.md`, `testing-report.md`, and docs
- Updated `proposal.md` if scope changed

## Done Criteria

- Branch exists; all tasks checked
- `testing-report.md` exists with at least unit test results
- Tell the user to run `/verify <change-name>` next

## Rules

- Follow `engineering-guardrails.mdc`, `specboot-workflow.mdc`, and `github-workflow.mdc`
- Prefer containerized verification when Dockerfile exists; otherwise document the gap

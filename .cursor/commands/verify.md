---
description: Verify tests, container smoke, and security scan for a change
argument-hint: <change-name>
---

# Verify Change

The user invoked this command with: $ARGUMENTS

## Goal

Run validation and record results in `testing-report.md`.

## Prerequisites

- `/apply` gate met: branch, code, tests, docs, `testing-report.md` exist
- If not, stop and tell the user to run `/apply <change-name>` first

## Inputs

- Change name from `$ARGUMENTS`
- `user-story.md` acceptance criteria
- `testing-report.md` template sections

## Steps

1. Run the project test suite (or document that no runner exists yet).
2. If `Dockerfile` or `docker-compose.yml` exists, run container smoke test.
   Otherwise note "container smoke: unavailable — no Dockerfile" in the report.
3. Note security considerations (dependency scan, secrets, input validation) in the report.
4. Verify each acceptance criterion from `user-story.md`; record pass/fail.
5. Update `testing-report.md` verify checklist and summary.
6. Fix failures if straightforward; otherwise report blockers.

## Outputs

- Updated `testing-report.md` with verify checklist complete

## Done Criteria

- All tests pass (or failures documented as blockers)
- Verify checklist filled
- Summary is `ready` or `blocked`
- If ready, tell the user to run `/code-review <change-name>` next

## Rules

- Follow `engineering-guardrails.mdc` and `specboot-workflow.mdc`

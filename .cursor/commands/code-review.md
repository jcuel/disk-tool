---
description: Structured code review against delta spec and guardrails
argument-hint: <change-name>
---

# Code Review

Alias: `/code_review`

The user invoked this command with: $ARGUMENTS

## Goal

Review the implementation against the delta spec, user story, and engineering guardrails.

## Prerequisites

- `/verify` gate met: tests green, `testing-report.md` verify checklist complete
- If not, stop and tell the user to run `/verify <change-name>` first

## Inputs

- Change name from `$ARGUMENTS`
- Git diff on the feature branch
- Delta spec at `openspec/changes/<change-name>/specs/disk-tool/spec.md`
- `user-story.md` acceptance criteria

## Steps

1. Review code for correctness, security, and minimal scope.
2. Check implementation matches delta spec requirements.
3. Confirm tests cover acceptance criteria.
4. Record findings in `testing-report.md` code review section:
   - Critical (must fix), Suggestion, Nice-to-have
5. Set summary to `ready` only if no critical blockers remain.

## Outputs

- Updated `testing-report.md` with review notes and final summary

## Done Criteria

- Review table populated in `testing-report.md`
- No critical blockers, or blockers listed with required fixes
- If ready, tell the user to run `/archive <change-name>` next

## Rules

- Follow `engineering-guardrails.mdc` and `specboot-workflow.mdc`
- Do not commit or push

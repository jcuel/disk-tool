---
description: Scaffold proposal artifacts and delta spec for a change
argument-hint: <change-name>
---

# Propose Change

The user invoked this command with: $ARGUMENTS

## Goal

Create proposal artifacts under `openspec/changes/<change-name>/` from the refined user story.

## Prerequisites

- `openspec/changes/<change-name>/user-story.md` must exist
- If missing, stop and tell the user to run `/enrich-us` first

## Inputs

- Change name from `$ARGUMENTS`
- `user-story.md` in the change folder
- Templates in `.cursor/templates/openspec/`

## Steps

1. Read `user-story.md` and derive scope, design approach, and tasks.
2. Create or update:
   - `proposal.md` (status: `draft`)
   - `design.md`
   - `tasks.md` (checkbox list aligned with implementation)
   - `specs/disk-tool/spec.md` (delta spec using ADDED/MODIFIED/REMOVED blocks)
3. Map each acceptance criterion to at least one delta requirement or task.

## Outputs

- `proposal.md`, `design.md`, `tasks.md`
- `specs/disk-tool/spec.md` (delta)

## Done Criteria

- All four artifacts exist and reference the same change name
- Tell the user to run `/apply <change-name>` next

## Rules

- Follow `engineering-guardrails.mdc` and `specboot-workflow.mdc`
- Keep the delta spec focused on behavior changes, not implementation details

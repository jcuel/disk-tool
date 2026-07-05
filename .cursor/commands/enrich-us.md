---
description: Enrich a raw user story into a refined OpenSpec artifact
argument-hint: [change-name or story text]
---

# Enrich User Story

Alias: `/enrich_us`

The user invoked this command with: $ARGUMENTS

## Goal

Expand a raw user story into `openspec/changes/<change-name>/user-story.md` with acceptance criteria.

## Inputs

- Raw story from `$ARGUMENTS` or the current conversation
- Template: `.cursor/templates/openspec/user-story.md`

## Steps

1. Derive `<change-name>` as a kebab-case slug from the story (or use it directly if provided).
2. Create `openspec/changes/<change-name>/` if it does not exist.
3. Write `user-story.md` using the template: actors, goal, benefit, acceptance criteria, out-of-scope.
4. Confirm the gate: acceptance criteria are testable and unambiguous.

## Outputs

- `openspec/changes/<change-name>/user-story.md`

## Done Criteria

- File exists with at least two acceptance criteria
- Tell the user to run `/propose <change-name>` next

## Rules

- Follow `engineering-guardrails.mdc` and `specboot-workflow.mdc`
- Do not create proposal artifacts yet — that is `/propose`

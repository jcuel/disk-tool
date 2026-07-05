# disk-tool — GitHub Project

Board: https://github.com/users/jcuel/projects/3

## Status (project field)

Maps to SPECBOOT stages:

| Status | When to use |
|--------|-------------|
| **Backlog** | Idea or `/enrich-us` not started; future work |
| **Ready** | Scoped issue; `/propose` done or not needed; pick up next |
| **In progress** | Branch open; `/apply` active |
| **In review** | PR open; `/verify` + `/code-review` |
| **Done** | Merged to `master`; `/archive` complete |

Move status on the project board when stage changes. Keep issue **comments** updated at each transition.

## Milestones

| Milestone | Scope |
|-----------|--------|
| **v0.1 — Foundation** | Shipped (scanner, API, UI, insights v1, CI) |
| **v0.2 — Fixes & tooling** | #3 CLI JSON bug, #6 govulncheck, #7 template restore |
| **v1.0 — Cleanup insights v2** | #1 open in explorer, #2 safe delete |
| **Future** | #4 duplicates, #5 age-based cleanup (needs `/propose`) |

Set milestone when creating an issue. Filter the **By milestone** view to plan releases.

## Comments convention

Post a short comment when status changes:

```
**Status:** Ready → In progress
**Branch:** fix/cli-json-windows
**Notes:** Starting /apply; repro confirmed on Windows.
```

On PR open, link the PR. On merge, note `/archive` and close.

## Recommended views (create in UI)

GitHub does not expose view creation in the stable API — add these tabs manually (**New view** on the project):

### 1. Board (default)

- Layout: **Board**
- Group by: **Status**
- Fields: Title, Assignees, Labels, Milestone, Linked pull requests

### 2. Ready queue

- Layout: **Table**
- Filter: `Status = Ready`
- Sort: **Priority** (label `priority:high` first) or Milestone due date

### 3. By milestone

- Layout: **Table**
- Group by: **Milestone**
- Sort: Status

### 4. In review

- Layout: **Table**
- Filter: `Status = In review`
- Fields: Reviewers, Linked pull requests

Rename the default **View 1** tab to **All issues** (table, no filter).

## Quick links

- [Issues](https://github.com/jcuel/disk-tool/issues)
- [Milestones](https://github.com/jcuel/disk-tool/milestones)
- [CI](https://github.com/jcuel/disk-tool/actions)

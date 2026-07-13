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
| **Done** | Merged to `dev` (or `master` for release); `/archive` complete |

Move status on the project board when stage changes. Keep issue **comments** updated at each transition.

## Milestones

| Milestone | Scope |
|-----------|--------|
| **v0.1 — Foundation** | Shipped (scanner, API, UI, insights v1, CI) |
| **v0.2 — Fixes & tooling** | #3 CLI JSON bug, #6 govulncheck |
| **v1.0 — Cleanup insights v2** | #1 open in explorer, #2 safe delete |
| **v1.1 — Safety & maintenance** | OS zones, safety grid, maintenance presets, age-based, duplicates |
| **Future** | macOS zones, [Docker cleanup](https://github.com/jcuel/disk-tool/issues/77) |

Set milestone when creating an issue. Filter the **By milestone** view to plan releases.

## Comments convention

Post a short comment when status changes:

```
**Status:** Ready → In progress
**Branch:** fix/cli-json-windows
**Notes:** Starting /apply; repro confirmed on Windows.
```

On PR open, link the PR. On merge, note `/archive` and close.

## Sync board with issues

Project **Status** can drift when issues close via merged PRs.

### Automatic (CI)

After a successful **CI** run on a push to `master`, [`.github/workflows/sync-project-board.yml`](../workflows/sync-project-board.yml) runs `scripts/sync-project-board.sh`.

After a push to **`master`** (release merge), [`.github/workflows/release-version.yml`](../workflows/release-version.yml) bumps semver in `openspec/config.yaml`, the CLI, and web package, then tags `vX.Y.Z`. See [CONTRIBUTING.md](../CONTRIBUTING.md#release-to-master-version-bump).

If `master` drifts ahead of `dev` (common after release), [`.github/workflows/sync-dev-from-master.yml`](../workflows/sync-dev-from-master.yml) opens a **`master` → `dev`** sync PR. Monitor with [branch-drift-check.yml](../workflows/branch-drift-check.yml).

Requires repository secret **`GH_PROJECT_SYNC`**: a **classic** PAT (`ghp_…`) with **`project`** and **`repo`** scopes. Fine-grained PATs cannot write user-owned project boards ([GitHub docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)). The default `GITHUB_TOKEN` cannot write user-owned project boards either.

Create or rotate the secret:

```powershell
./scripts/setup-gh-project-sync-secret.ps1
```

```bash
bash scripts/setup-gh-project-sync-secret.sh
```

Pre-filled token page: [classic PAT — project + repo](https://github.com/settings/tokens/new?scopes=project,repo&description=disk-tool-GH_PROJECT_SYNC). Pick an expiration (e.g. 90 days) and generate.

If the secret is missing, the workflow skips sync with a log message (CI still passes).

### Manual

```powershell
./scripts/sync-project-board.ps1
```

```bash
bash scripts/sync-project-board.sh
```

Board IDs live in [`scripts/project-board.env`](../scripts/project-board.env) (override via env vars if needed).

- Closed issues on the board → **Done**
- Open issues with `openspec/changes/<change>/proposal.md` still in **Backlog** → **Ready** (matches `/propose <change>` or change name in issue body)

Requires `gh` with `project` scope. Run from repo root.

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

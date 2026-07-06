# disk-tool

Cross-platform local web disk usage analyzer (TreeSize-inspired). **Overview first, drill down on demand** — no full-tree scan upfront.

## Product goal

Guide users to **where disk space is consumed**, surface **actionable cleanup candidates**, and produce **support-ready reports**.

| User need | disk-tool response |
|-----------|-------------------|
| "What's eating my disk?" | Overview % breakdown → drill into heavy folders |
| "Leftover dev junk?" | Detects `node_modules`, `.venv`, `target/`, caches |
| "Old installers in Downloads?" | Flags large `.exe`, `.zip`, `.msi`, etc. |
| "Need a ticket for IT/support?" | **Copy report** or export **Support ticket** (plain text) |

Insights improve as you drill — scan `Users` or `Projects` to uncover nested `node_modules` and caches.

## Safe cleanup guide

| Zone | Meaning | Delete via disk-tool? |
|------|---------|----------------------|
| **review** | Regenerable dev artifacts (`node_modules`, `target/`) | Yes, after review |
| **maintenance** | User temp folders | Yes, with age filter |
| **caution** | Shared caches (`.npm`, `.gradle`) | Review carefully |
| **diagnostic** | Crash dumps, minidumps | No — view only |
| **critical_os** | Windows, `/usr`, Program Files | No |
| **forbidden** | System32, `/proc` | Hidden from scan |

Use **Maintenance presets** (Dev reclaim, Temp cleanup) for quick safe wins. Scan your user profile rather than `C:\` or `/` when cleaning up.

## How scanning works

1. **Overview** — lists top-level folders with accurate total sizes (% of root). Each top-level folder is sized in parallel (no full-tree walk).
2. **Drill-down** — click a folder or **Scan folder deeper** to scan **5 levels** inside that branch only.
3. **Repeat** — folders at the depth limit show a **+** badge; drill again to go deeper.

This avoids loading huge trees into memory and prevents scan deadlocks on large disks.

## Live demo

Try a **sample scan in your browser** (no install): [jcuel.github.io/disk-tool/demo/](https://jcuel.github.io/disk-tool/demo/)

Product overview: [jcuel.github.io/disk-tool/](https://jcuel.github.io/disk-tool/)

## Download

Pre-built binaries for Windows, Linux, and macOS are on the [Releases](https://github.com/jcuel/disk-tool/releases/latest) page. Download, extract, and run `disk-tool serve` (or `disk-tool.exe serve` on Windows).

## Quick start

```powershell
# Build (requires Node.js + Go)
cd web; npm ci; npm run build; cd ..
Copy-Item -Recurse -Force web\dist\* cmd\disk-tool\static\
go build -o bin/disk-tool ./cmd/disk-tool

# Web UI (opens browser)
./bin/disk-tool serve

# Headless overview (top-level only)
./bin/disk-tool scan C:\Users --json

# Full tree (legacy, use for small paths only)
./bin/disk-tool scan . --full --json
```

On Linux/macOS use `make build` if Make is available.

## Commands

| Command | Description |
|---------|-------------|
| `disk-tool serve [--port 8080]` | Local web UI on 127.0.0.1 |
| `disk-tool scan <path> [--json]` | Scan without UI |
| `disk-tool version` | Print version |

## API (localhost)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/roots` | Common scan roots |
| GET | `/api/disk?path=...` | Volume capacity, used, and free bytes |
| POST | `/api/scans` | Start scan `{ "root": "..." }` |
| GET | `/api/scans/{id}` | Scan status and tree |
| DELETE | `/api/scans/{id}` | Cancel overview scan |
| POST | `/api/scans/{id}/expand` | Drill into folder `{ "path": "...", "depth": 5 }` |
| WS | `/api/scans/{id}/events` | Progress + expand events |
| GET | `/api/scans/{id}/export?format=json\|html\|ticket\|cleanup-json\|cleanup-html\|cleanup-ticket` | Export scan or cleanup report |
| POST | `/api/scans/{id}/open` | Open path in OS file manager `{ "path": "..." }` |
| POST | `/api/scans/{id}/delete` | Delete path under scan root `{ "path": "...", "confirm": true }` |
| POST | `/api/scans/{id}/cleanup` | Bulk cleanup dry-run or execute `{ "paths": [...], "dryRun": true, "confirm": false }` |

## CI and smoke tests

Pipeline: [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

| Job | Environment | Checks |
|-----|-------------|--------|
| `test-linux` | ubuntu-latest | `go test`, govulncheck, build, [`scripts/smoke-api.sh`](scripts/smoke-api.sh) |
| `test-windows` | windows-latest | `build.ps1`, govulncheck, [`scripts/smoke-api.ps1`](scripts/smoke-api.ps1) |
| `docker-smoke` | ubuntu + Docker | CLI scan + in-container API smoke |
| `e2e-linux` | ubuntu-latest | [`scripts/e2e-run.sh`](scripts/e2e-run.sh) — Cypress browser tests |
| `security` | Trivy | Filesystem + container image (CRITICAL/HIGH); PR comment with scan tables |
| `sync` | ubuntu-latest | Project board sync after merge ([`sync-project-board.yml`](.github/workflows/sync-project-board.yml)) |

**Local smoke**

```powershell
# Windows
.\build.ps1
.\scripts\smoke-api.ps1

# Linux / WSL / macOS
make build
make smoke-api
bash scripts/smoke-docker.sh   # requires Docker

# Cypress E2E (Linux/WSL — builds server + runs headless browser tests)
bash scripts/e2e-run.sh

On pull requests, CI uploads E2E screenshots as artifacts and posts a sticky PR comment with inline UI snapshots (artifact link fallback when inline upload is unavailable).

# Windows E2E (after ./build.ps1)
.\scripts\e2e-run.ps1
```

## Docker

```bash
docker compose build
docker compose run --rm smoke          # CLI JSON scan
docker compose run --rm smoke-api      # API smoke in Alpine
```

## Development workflow (SPECBOOT)

| Command | Stage |
|---------|-------|
| `/enrich-us` | Refine user story |
| `/propose` | Create proposal artifacts |
| `/apply` | Implement on a feature branch |
| `/verify` | Run tests and validation |
| `/code-review` | Review against spec |
| `/archive` | Merge spec deltas |
| `/commit` | Git commit (with approval) |

OpenSpec layout: `openspec/`. Slash commands: `.cursor/commands/`.

**Tracking:** [Issues](https://github.com/jcuel/disk-tool/issues) · [Project board](https://github.com/users/jcuel/projects/3) · [Milestones](https://github.com/jcuel/disk-tool/milestones). See [`.github/PROJECT.md`](.github/PROJECT.md) for status, comments, and views.

## Community

| Document | Purpose |
|----------|---------|
| [CONTRIBUTING.md](CONTRIBUTING.md) | Fork → `dev` workflow, tests, PR expectations |
| [Issue templates](.github/ISSUE_TEMPLATE/) | Bug, OpenSpec change, chore, question |
| [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) | Community standards |
| [SECURITY.md](SECURITY.md) | Vulnerability reporting |
| [LICENSE](LICENSE) | PolyForm Noncommercial — collaborate freely; commercial use requires permission |

Contributors open PRs against **`dev`**. Only the maintainer merges **`dev` → `master`** for releases.

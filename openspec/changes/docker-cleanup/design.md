# Design: Docker maintenance module

**Change:** `docker-cleanup`

## Approach

CLI-first reclaim via `docker system prune -af` (no volumes). Path awareness reports known Desktop/Linux data roots as non-deletable caution candidates when useful.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Delete `/var/lib/docker` via FS | Fast reclaim | Breaks Docker/daemon | Rejected |
| Prune including volumes | More space | Risk of DB data loss | Deferred |
| CLI prune + insights | Safe, matches product UX | Needs Docker on PATH | **Chosen** |

## Components

| Area | Role |
|------|------|
| `internal/docker` | Detect, SystemDF, PruneDryRun, Prune, DataRoots |
| `internal/insights` | Emit CategoryDocker candidates |
| `internal/safety/presets.go` | `docker-reclaim` preset |
| `internal/api` | GET docker status, POST docker prune |
| `web/src/main.ts` | Preset → review → confirm modal |

## Prune command

```
docker system prune -af
```

No `--volumes`. Dry-run uses `docker system df` reclaimable totals (prune has limited dry-run; we report df reclaimable before execute).

## Data roots (report only, never FS-delete)

- Windows: `%LOCALAPPDATA%\Docker`
- macOS: `~/Library/Containers/com.docker.docker/Data`
- Linux: `/var/lib/docker`, `~/.local/share/docker`

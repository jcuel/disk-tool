# Contributing to disk-tool

Thank you for your interest in contributing. This project welcomes collaboration under the terms of the [LICENSE](LICENSE) (noncommercial use).

## Branch model

| Branch | Purpose | Who merges here |
|--------|---------|-----------------|
| **`dev`** | Integration branch for all contributor work | Maintainers (via reviewed PRs) |
| **`master`** | Stable / release-ready | **Maintainer only** (`dev` → `master`) |

**Contributors:** fork the repo, branch from `dev`, and open pull requests **into `dev`**.

**Do not** open PRs targeting `master` unless you are the maintainer performing a release integration.

```
fork ──► feat/your-change ──► PR ──► dev ──► (maintainer) ──► master
```

## Getting started

1. Fork [jcuel/disk-tool](https://github.com/jcuel/disk-tool) on GitHub.
2. Clone your fork and add upstream:
   ```bash
   git clone https://github.com/YOUR_USER/disk-tool.git
   cd disk-tool
   git remote add upstream https://github.com/jcuel/disk-tool.git
   ```
3. Create a feature branch from **`dev`**:
   ```bash
   git fetch upstream
   git checkout -b feat/my-change upstream/dev
   ```
4. Make changes, run tests locally (see [README.md](README.md)).
5. Push to your fork and open a PR **against `dev`**.

## Pull requests

- Use the [pull request template](.github/pull_request_template.md).
- Link a GitHub issue when the work is non-trivial (`Closes #N` in the description).
- Keep PRs focused; one logical change per PR.
- CI must pass before merge.
- For larger features, follow the OpenSpec / SPECBOOT flow documented in [README.md](README.md) and `.cursor/rules/`.

## Development checks

```bash
# Linux / WSL / macOS
go test ./...
bash scripts/smoke-api.sh ./bin/disk-tool
bash scripts/e2e-run.sh          # optional; needs Node + Go

# Windows
.\build.ps1
.\scripts\smoke-api.ps1
```

## Code style

- Match existing patterns in the file you edit.
- Go: `go test ./...` must pass; run `gofmt` on changed files.
- Web: TypeScript in `web/`; build with `npm run build` in `web/`.
- No unrelated drive-by refactors in the same PR.

## Security

Report vulnerabilities privately — see [SECURITY.md](SECURITY.md). Do not open public issues for security bugs.

## Code of conduct

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md). Participants are expected to uphold it.

## Questions

Open a [GitHub Discussion](https://github.com/jcuel/disk-tool/discussions) or an issue labeled `question` if Discussions are not enabled.

## Summary

<!-- What changed and why? One short paragraph. -->

## Target branch

- [ ] This PR targets **`dev`** (required for contributor changes)
- [ ] Maintainer release: **`dev` → `master`** (maintainer only)

For maintainer releases, add **`Release-Version: X.Y.Z`** below (e.g. `Release-Version: 1.1.0`) so CI tags the correct semver. If omitted, the release workflow bumps the **minor** version automatically.

## Linked issue

Closes #<!-- issue number, or "N/A — chore/docs" -->

## Test plan

- [ ] `go test ./...`
- [ ] API smoke (`scripts/smoke-api.sh` / `.ps1`)
- [ ] Other: <!-- e.g. e2e-run.sh, manual UI check -->

## OpenSpec (if applicable)

Change folder: `openspec/changes/<name>/`  
Testing report: <!-- link or "N/A" -->

## Checklist

- [ ] CI green (or explain expected failures)
- [ ] No secrets or local paths committed
- [ ] [CONTRIBUTING.md](../CONTRIBUTING.md) branch model followed

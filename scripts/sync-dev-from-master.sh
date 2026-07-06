#!/usr/bin/env bash
# Open or update a master → dev PR when master is ahead of dev (post-release drift).
set -euo pipefail

REPO="${1:-$(gh repo view --json nameWithOwner -q .nameWithOwner)}"
REMOTE="${REMOTE:-origin}"
AUTO_MERGE="${AUTO_MERGE_BRANCH_SYNC:-false}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if ! command -v gh >/dev/null 2>&1; then
  echo "gh CLI is required" >&2
  exit 1
fi

git fetch "$REMOTE" master dev

set +e
report="$(bash "$ROOT/scripts/check-branch-drift.sh" --quiet 2>&1)"
code=$?
set -e
echo "$report" >&2
eval "$(echo "$report" | grep -E '^[a-z_]+=')"

if [[ "$branch_drift_status" == "aligned" || "$branch_drift_status" == "dev_ahead" ]]; then
  echo "No sync required ($branch_drift_status)."
  exit 0
fi

if [[ "$branch_drift_status" == "diverged" && "$tree_match" != "true" ]]; then
  echo "::warning::Branches diverged with different trees. Open a manual master → dev PR and resolve conflicts."
fi

existing="$(gh pr list --repo "$REPO" --base dev --head master --state open --json number --jq '.[0].number' 2>/dev/null || true)"
if [[ -n "$existing" && "$existing" != "null" ]]; then
  echo "Sync PR already open: #$existing"
  pr_num="$existing"
else
  body="$(cat <<EOF
## Summary
Automated sync after release drift: \`master\` is **${master_ahead}** commit(s) ahead of \`dev\`.

- **Status:** \`${branch_drift_status}\`
- **Tree match:** \`${tree_match}\` (history-only drift when \`true\`)

## Maintainer
Merge this PR to realign \`dev\` with \`master\` so both branches share release history.

See [CONTRIBUTING.md](../CONTRIBUTING.md#sync-dev-from-master-branch-drift).
EOF
)"
  pr_url="$(gh pr create --repo "$REPO" \
    --base dev \
    --head master \
    --title "chore: sync dev from master (branch drift)" \
    --body "$body")"
  pr_num="$(echo "$pr_url" | sed -n 's|.*/pull/\([0-9]*\)$|\1|p')"
  echo "Created sync PR #$pr_num"
fi

if [[ "$AUTO_MERGE" == "true" && "$tree_match" == "true" ]]; then
  echo "Auto-merge enabled for history-only drift."
  gh pr merge "$pr_num" --repo "$REPO" --merge --admin || \
    echo "Auto-merge failed; merge PR #$pr_num manually."
fi

echo "Sync PR: https://github.com/$REPO/pull/$pr_num"

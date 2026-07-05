#!/usr/bin/env bash
# Apply repository rulesets from .github/rulesets/*.json
# Primary path: .github/workflows/sync-branch-rulesets.yml on push to dev.
# Local use: requires gh CLI authenticated with admin access on the repo.
set -euo pipefail

REPO="${1:-$(gh repo view --json nameWithOwner -q .nameWithOwner)}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RULESET_DIR="$ROOT/.github/rulesets"

if ! command -v gh >/dev/null 2>&1; then
  echo "gh CLI is required" >&2
  exit 1
fi

echo "Repository: $REPO"

existing="$(gh api "repos/$REPO/rulesets" --jq '.[].name' 2>/dev/null || true)"

for file in "$RULESET_DIR"/*.json; do
  name="$(jq -r .name "$file")"
  echo "Applying ruleset: $name"

  if echo "$existing" | grep -Fxq "$name"; then
    id="$(gh api "repos/$REPO/rulesets" --jq ".[] | select(.name==\"$name\") | .id")"
    gh api "repos/$REPO/rulesets/$id" --method PUT --input "$file" >/dev/null
    echo "  updated ruleset id=$id"
  else
    id="$(gh api "repos/$REPO/rulesets" --method POST --input "$file" --jq .id)"
    echo "  created ruleset id=$id"
  fi
done

echo "Done. Verify at https://github.com/$REPO/settings/rules"

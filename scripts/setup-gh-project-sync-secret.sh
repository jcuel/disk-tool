#!/usr/bin/env bash
# Create/replace GH_PROJECT_SYNC with a classic PAT (required for user-owned project boards).
set -euo pipefail

REPO="${1:-jcuel/disk-tool}"
PAT_URL="https://github.com/settings/tokens/new?scopes=project,repo&description=disk-tool-GH_PROJECT_SYNC"

echo "GitHub fine-grained PATs cannot access user-owned Projects (board #3)."
echo "Use a classic PAT (ghp_...) with scopes: project, repo"
echo ""
echo "Opening token creation page..."
if command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$PAT_URL" >/dev/null 2>&1 || true
elif command -v open >/dev/null 2>&1; then
  open "$PAT_URL" || true
else
  echo "$PAT_URL"
fi

echo ""
read -r -s -p "Paste classic PAT (ghp_...), then Enter: " PAT
echo ""
if [[ -z "$PAT" ]]; then
  echo "No token provided." >&2
  exit 1
fi
if [[ ! "$PAT" =~ ^ghp_ ]]; then
  echo "Expected a classic PAT starting with ghp_." >&2
  exit 1
fi

echo "Setting GH_PROJECT_SYNC on $REPO..."
gh secret set GH_PROJECT_SYNC --repo "$REPO" --body "$PAT"
echo "Secret updated."

echo "Verifying token can list project items..."
export GH_TOKEN="$PAT"
if gh project item-list 3 --owner jcuel --format json --limit 1 >/dev/null 2>&1; then
  echo "Project API check passed."
else
  echo "Warning: project API check failed — confirm project + repo scopes." >&2
  exit 1
fi

echo "Done. Next merge to master will auto-sync the board."

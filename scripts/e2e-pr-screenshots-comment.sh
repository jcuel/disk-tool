#!/usr/bin/env bash
# Push Cypress PNGs to e2e-screenshots/pr-N branch and write PR comment markdown.
set -euo pipefail

PR_NUMBER="${PR_NUMBER:?PR_NUMBER required}"
RUN_ID="${RUN_ID:?RUN_ID required}"
SCREENSHOT_ROOT="${SCREENSHOT_ROOT:-web/cypress/screenshots}"
COMMENT_FILE="${COMMENT_FILE:-.e2e-screenshots-comment.md}"
BRANCH="e2e-screenshots/pr-${PR_NUMBER}"
DEST="pr-${PR_NUMBER}"
OWNER="${GITHUB_REPOSITORY%%/*}"
REPO="${GITHUB_REPOSITORY##*/}"
ROOT="${GITHUB_WORKSPACE:?GITHUB_WORKSPACE required}"

mapfile -t PNGS < <(find "$ROOT/$SCREENSHOT_ROOT" -name '*.png' 2>/dev/null | sort)
if [[ ${#PNGS[@]} -eq 0 ]]; then
  echo "No PNG screenshots under ${SCREENSHOT_ROOT}"
  exit 1
fi

WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

git -C "$WORKDIR" init -q
git -C "$WORKDIR" remote add origin "https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
git -C "$WORKDIR" checkout -B "$BRANCH"
mkdir -p "$WORKDIR/$DEST"
rm -f "$WORKDIR/$DEST"/*.png

i=1
for png in "${PNGS[@]}"; do
  base="$(basename "$png" .png)"
  out="$WORKDIR/$DEST/$(printf '%02d-%s.png' "$i" "$base")"
  cp "$png" "$out"
  i=$((i + 1))
done

git -C "$WORKDIR" config user.name "github-actions[bot]"
git -C "$WORKDIR" config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git -C "$WORKDIR" add "$DEST"
git -C "$WORKDIR" commit -m "E2E screenshots for PR #${PR_NUMBER} (run ${RUN_ID})"
git -C "$WORKDIR" push -f origin "$BRANCH"

{
  echo "## E2E UI snapshots"
  echo ""
  echo "Run [\`${RUN_ID}\`](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${RUN_ID}) · Artifact: \`cypress-screenshots-pr-${PR_NUMBER}\`"
  echo ""
  for f in "$WORKDIR/$DEST"/*.png; do
    name="$(basename "$f")"
    label="${name%.png}"
    url="https://raw.githubusercontent.com/${OWNER}/${REPO}/${BRANCH}/${DEST}/${name}"
    echo "<details open><summary>${label}</summary>"
    echo ""
    echo "![${label}](${url})"
    echo ""
    echo "</details>"
    echo ""
  done
} > "$ROOT/$COMMENT_FILE"

echo "Published ${#PNGS[@]} screenshot(s) to ${BRANCH}; comment body at ${COMMENT_FILE}"

#!/usr/bin/env bash
# Build markdown PR comment from Trivy filesystem + image scan outputs.
set -euo pipefail

COMMENT_FILE="${COMMENT_FILE:-.trivy-pr-comment.md}"
RUN_ID="${RUN_ID:?RUN_ID required}"
FS_OUTCOME="${FS_OUTCOME:-unknown}"
IMAGE_OUTCOME="${IMAGE_OUTCOME:-unknown}"
ROOT="${GITHUB_WORKSPACE:?GITHUB_WORKSPACE required}"

status_label() {
  case "$1" in
    success) echo "PASS" ;;
    failure) echo "FAIL" ;;
    skipped) echo "SKIPPED" ;;
    *) echo "$1" ;;
  esac
}

report_block() {
  local file="$1"
  local max=12000
  if [[ ! -f "$file" ]]; then
    echo "_No report file generated._"
    return
  fi
  if [[ ! -s "$file" ]]; then
    echo "_No CRITICAL/HIGH findings._"
    return
  fi
  local content
  content="$(cat "$file")"
  if ((${#content} > max)); then
    content="${content:0:max}"$'\n\n...(truncated)...'
  fi
  printf '```\n%s\n```' "$content"
}

{
  echo "## Trivy security scan"
  echo ""
  echo "Run [\`${RUN_ID}\`](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${RUN_ID}) · Severity filter: CRITICAL, HIGH"
  echo ""
  echo "### Filesystem — $(status_label "$FS_OUTCOME")"
  echo ""
  echo "<details><summary>Filesystem scan (\`trivy-fs.txt\`)</summary>"
  echo ""
  report_block "$ROOT/trivy-fs.txt"
  echo ""
  echo "</details>"
  echo ""
  echo "### Container image — $(status_label "$IMAGE_OUTCOME")"
  echo ""
  echo "<details><summary>Image scan (\`trivy-image.txt\`)</summary>"
  echo ""
  report_block "$ROOT/trivy-image.txt"
  echo ""
  echo "</details>"
  if [[ "$FS_OUTCOME" == "failure" || "$IMAGE_OUTCOME" == "failure" ]]; then
    echo ""
    echo "> CI will fail until CRITICAL/HIGH findings are resolved or accepted in \`.trivyignore\`."
  fi
} > "$ROOT/$COMMENT_FILE"

echo "Wrote Trivy PR comment to ${COMMENT_FILE}"

#!/usr/bin/env bash
# Sync GitHub Project board Status with issue state.
# Closed issues -> Done; open issues with OpenSpec proposal -> Ready.
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
env_file="${script_dir}/project-board.env"
if [[ -f "$env_file" ]]; then
  set -a
  # shellcheck source=/dev/null
  source "$env_file"
  set +a
fi

OWNER="${PROJECT_BOARD_OWNER:-jcuel}"
PROJECT_NUMBER="${PROJECT_BOARD_NUMBER:-3}"
REPO="${PROJECT_BOARD_REPO:-jcuel/disk-tool}"
PROJECT_ID="${PROJECT_BOARD_ID:-PVT_kwHOA0E_1c4Bcf_W}"
FIELD_ID="${PROJECT_BOARD_STATUS_FIELD:-PVTSSF_lAHOA0E_1c4Bcf_WzhXIUKk}"
STATUS_DONE="${PROJECT_BOARD_STATUS_DONE:-0d582c4b}"
STATUS_READY="${PROJECT_BOARD_STATUS_READY:-7eda7d29}"

if ! command -v gh >/dev/null 2>&1; then
  echo "gh CLI is required" >&2
  exit 1
fi
if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

find_change_for_issue() {
  local body="$1"
  local change=""
  if [[ "$body" =~ /propose[[:space:]]+([a-z0-9-]+) ]]; then
    change="${BASH_REMATCH[1]}"
    if [[ -f "openspec/changes/${change}/proposal.md" ]]; then
      echo "$change"
      return
    fi
  fi
  for dir in openspec/changes/*/; do
    [[ -d "$dir" ]] || continue
    change="$(basename "$dir")"
    [[ -f "$dir/proposal.md" ]] || continue
    if [[ "$body" == *"$change"* ]]; then
      echo "$change"
      return
    fi
  done
}

items_json="$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json --limit 100)"
count="$(echo "$items_json" | jq '.items | length')"
updated=0

for ((i = 0; i < count; i++)); do
  num="$(echo "$items_json" | jq -r ".items[$i].content.number // empty")"
  item_id="$(echo "$items_json" | jq -r ".items[$i].id")"
  status="$(echo "$items_json" | jq -r ".items[$i].status // empty")"
  [[ -n "$num" ]] || continue

  state="$(gh issue view "$num" --repo "$REPO" --json state --jq '.state' 2>/dev/null || echo "")"

  if [[ "$state" == "CLOSED" && "$status" != "Done" ]]; then
    echo "Issue #$num: $status -> Done"
    gh project item-edit --id "$item_id" --project-id "$PROJECT_ID" --field-id "$FIELD_ID" \
      --single-select-option-id "$STATUS_DONE"
    updated=$((updated + 1))
    continue
  fi

  if [[ "$state" == "OPEN" && "$status" == "Backlog" ]]; then
    body="$(gh issue view "$num" --repo "$REPO" --json body --jq '.body')"
    change="$(find_change_for_issue "$body")"
    if [[ -n "$change" ]]; then
      echo "Issue #$num: Backlog -> Ready (openspec/changes/${change})"
      gh project item-edit --id "$item_id" --project-id "$PROJECT_ID" --field-id "$FIELD_ID" \
        --single-select-option-id "$STATUS_READY"
      updated=$((updated + 1))
    fi
  fi
done

echo "Project board sync complete ($updated item(s) updated)."

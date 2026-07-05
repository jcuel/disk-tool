#!/usr/bin/env bash
# Sync GitHub Project board Status with issue state (closed -> Done).
set -euo pipefail

OWNER="${1:-jcuel}"
PROJECT_NUMBER="${2:-3}"
REPO="${3:-jcuel/disk-tool}"

PROJECT_ID="PVT_kwHOA0E_1c4Bcf_W"
FIELD_ID="PVTSSF_lAHOA0E_1c4Bcf_WzhXIUKk"
STATUS_DONE="0d582c4b"
STATUS_READY="7eda7d29"

mapfile -t CLOSED < <(gh issue list --repo "$REPO" --state closed --json number --jq '.[].number')

items_json="$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json --limit 100)"
count="$(echo "$items_json" | jq '.items | length')"

for ((i = 0; i < count; i++)); do
  num="$(echo "$items_json" | jq -r ".items[$i].content.number // empty")"
  item_id="$(echo "$items_json" | jq -r ".items[$i].id")"
  status="$(echo "$items_json" | jq -r ".items[$i].status // empty")"
  state="$(gh issue view "$num" --repo "$REPO" --json state --jq '.state' 2>/dev/null || echo "")"

  if [[ "$state" == "CLOSED" && "$status" != "Done" ]]; then
    echo "Issue #$num: $status -> Done"
    gh project item-edit --id "$item_id" --project-id "$PROJECT_ID" --field-id "$FIELD_ID" \
      --single-select-option-id "$STATUS_DONE"
  elif [[ "$state" == "OPEN" && "$status" == "Backlog" ]]; then
    change=""
    title="$(gh issue view "$num" --repo "$REPO" --json title --jq '.title')"
    body="$(gh issue view "$num" --repo "$REPO" --json body --jq '.body')"
    if [[ "$body" =~ /propose[[:space:]]+([a-z0-9-]+) ]]; then
      change="${BASH_REMATCH[1]}"
    fi
    case "$num" in
      4) change="duplicate-detection" ;;
      5) change="age-based-cleanup" ;;
    esac
    if [[ -n "$change" && -f "openspec/changes/${change}/proposal.md" ]]; then
      echo "Issue #$num: Backlog -> Ready (proposal exists)"
      gh project item-edit --id "$item_id" --project-id "$PROJECT_ID" --field-id "$FIELD_ID" \
        --single-select-option-id "$STATUS_READY"
    fi
  fi
done

echo "Project board sync complete."

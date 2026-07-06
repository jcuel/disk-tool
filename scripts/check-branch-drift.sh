#!/usr/bin/env bash
# Report drift between origin/master and origin/dev.
# Exit 0 when aligned or dev is ahead (normal pre-release).
# Exit 1 when master is ahead (sync recommended). Exit 2 when diverged with different trees.
set -euo pipefail

QUIET=false
for arg in "$@"; do
  case "$arg" in
    --quiet) QUIET=true ;;
  esac
done

log() {
  if [[ "$QUIET" != "true" ]]; then
    echo "$@" >&2
  fi
}

REMOTE="${REMOTE:-origin}"
MASTER_REF="${MASTER_REF:-$REMOTE/master}"
DEV_REF="${DEV_REF:-$REMOTE/dev}"

git fetch "$REMOTE" master dev >/dev/null 2>&1 || git fetch "$REMOTE" master dev

master_ahead="$(git rev-list --count "$DEV_REF..$MASTER_REF" 2>/dev/null || echo 0)"
dev_ahead="$(git rev-list --count "$MASTER_REF..$DEV_REF" 2>/dev/null || echo 0)"

tree_match="false"
if git diff --quiet "$DEV_REF" "$MASTER_REF" 2>/dev/null; then
  tree_match="true"
fi

status="aligned"
if [[ "$master_ahead" -gt 0 && "$dev_ahead" -gt 0 ]]; then
  status="diverged"
elif [[ "$master_ahead" -gt 0 ]]; then
  status="master_ahead"
elif [[ "$dev_ahead" -gt 0 ]]; then
  status="dev_ahead"
fi

echo "branch_drift_status=$status"
echo "master_ahead=$master_ahead"
echo "dev_ahead=$dev_ahead"
echo "tree_match=$tree_match"
echo "master_ref=$(git rev-parse "$MASTER_REF")"
echo "dev_ref=$(git rev-parse "$DEV_REF")"

case "$status" in
  aligned)
    log "Branches aligned."
    exit 0
    ;;
  dev_ahead)
    log "dev is $dev_ahead commit(s) ahead of master (normal pre-release)."
    exit 0
    ;;
  master_ahead)
    if [[ "$tree_match" == "true" ]]; then
      log "History drift only (same tree). Run scripts/sync-dev-from-master.sh."
    else
      log "master is $master_ahead commit(s) ahead with file changes."
    fi
    exit 1
    ;;
  diverged)
    log "Diverged: master +$master_ahead, dev +$dev_ahead."
    if [[ "$tree_match" == "true" ]]; then
      exit 1
    fi
    exit 2
    ;;
esac

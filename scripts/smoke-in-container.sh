#!/bin/sh
# Runs inside Alpine container — API smoke without host curl dependency on CI runner
set -eu

PORT="${SMOKE_PORT:-18080}"
ROOT="${SMOKE_ROOT:-/data/src}"
BASE="http://127.0.0.1:$PORT"

disk-tool serve --port "$PORT" --no-open &
SERVER_PID=$!
trap 'kill $SERVER_PID 2>/dev/null || true' EXIT

i=0
while [ "$i" -lt 30 ]; do
  if wget -q -O /dev/null "$BASE/api/roots" 2>/dev/null; then
    break
  fi
  i=$((i + 1))
  sleep 0.2
done

wget -q -O - "$BASE/api/roots" | grep -q roots
echo "OK /api/roots"

SCAN_JSON=$(wget -q -O - --header='Content-Type: application/json' \
  --post-data="{\"root\":\"$ROOT\"}" "$BASE/api/scans")
echo "$SCAN_JSON" | grep -q scanId
SCAN_ID=$(echo "$SCAN_JSON" | sed -n 's/.*"scanId"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
echo "OK POST /api/scans ($SCAN_ID)"

i=0
STATUS=""
while [ "$i" -lt 60 ]; do
  STATUS=$(wget -q -O - "$BASE/api/scans/$SCAN_ID" | sed -n 's/.*"status"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
  if [ "$STATUS" = "completed" ]; then break; fi
  i=$((i + 1))
  sleep 0.5
done
[ "$STATUS" = "completed" ] || { echo "scan not completed: $STATUS" >&2; exit 1; }
echo "OK scan completed"

wget -q -O - "$BASE/api/scans/$SCAN_ID/export?format=ticket" | grep -q 'Disk usage report'
echo "OK export ticket"

echo "smoke-in-container: all checks passed"

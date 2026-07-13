#!/usr/bin/env bash
# API smoke test — Linux, macOS, WSL
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="${1:-$ROOT/bin/disk-tool}"
PORT="${SMOKE_PORT:-18080}"
BASE="http://127.0.0.1:$PORT"

if [[ ! -x "$BIN" && ! -f "$BIN" ]]; then
  echo "Binary not found: $BIN" >&2
  exit 1
fi

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

"$BIN" serve --port "$PORT" --no-open &
SERVER_PID=$!

for i in $(seq 1 30); do
  if curl -sf "$BASE/api/roots" >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done

curl -sf "$BASE/api/roots" | grep -q 'roots'
echo "OK /api/roots"

curl -sf "$BASE/api/disk?path=$ROOT" | grep -q '"total"'
echo "OK GET /api/disk"

SCAN_JSON=$(curl -sf -X POST "$BASE/api/scans" \
  -H 'Content-Type: application/json' \
  -d "{\"root\":\"$ROOT\"}")
echo "$SCAN_JSON" | grep -q 'scanId'
SCAN_ID=$(echo "$SCAN_JSON" | sed -n 's/.*"scanId"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
[[ -n "$SCAN_ID" ]] || { echo "missing scanId" >&2; exit 1; }
echo "OK POST /api/scans ($SCAN_ID)"

for i in $(seq 1 60); do
  STATUS=$(curl -sf "$BASE/api/scans/$SCAN_ID" | sed -n 's/.*"status"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
  [[ "$STATUS" == "completed" ]] && break
  sleep 0.5
done
[[ "$STATUS" == "completed" ]] || { echo "scan not completed: $STATUS" >&2; exit 1; }
echo "OK GET /api/scans/{id} completed"

JOB=$(curl -sf "$BASE/api/scans/$SCAN_ID")
echo "$JOB" | grep -q 'insights'
echo "$JOB" | grep -q 'tree'
echo "OK insights + tree"

FIRST_PATH=""
if command -v python3 >/dev/null 2>&1; then
  FIRST_PATH=$(JOB_JSON="$JOB" python3 -c 'import json,os; j=json.loads(os.environ["JOB_JSON"]); ch=(j.get("tree") or {}).get("children") or []; print(ch[0]["path"] if ch else "", end="")')
fi
if [[ -n "$FIRST_PATH" && "$FIRST_PATH" != "$ROOT" ]]; then
  curl -sf -X POST "$BASE/api/scans/$SCAN_ID/expand" \
    -H 'Content-Type: application/json' \
    -d "{\"path\":\"$FIRST_PATH\",\"depth\":5}" | grep -q expanding
  echo "OK POST /api/scans/{id}/expand"
fi

# Docker status endpoint (CLI may be absent; must still return 200)
DOCKER_JSON=$(curl -sf "$BASE/api/scans/$SCAN_ID/docker")
echo "$DOCKER_JSON" | grep -q 'usage'
echo "OK GET /api/scans/{id}/docker"

DOCKER_DRY=$(curl -sf -X POST "$BASE/api/scans/$SCAN_ID/docker/prune" \
  -H 'Content-Type: application/json' \
  -d '{"dryRun":true,"confirm":false,"confirmPhrase":""}')
echo "$DOCKER_DRY" | grep -q '"dryRun":true\|"dryRun": true'
echo "OK POST /api/scans/{id}/docker/prune dry-run"

curl -sf "$BASE/api/scans/$SCAN_ID/export?format=ticket" | grep -q 'Disk usage report'
echo "OK export ticket"

SMOKE_DIR="$ROOT/.smoke-cleanup-$$"
mkdir -p "$SMOKE_DIR/nested"
echo test > "$SMOKE_DIR/nested/file.txt"
CLEANUP_JSON=$(curl -sf -X POST "$BASE/api/scans/$SCAN_ID/cleanup" \
  -H 'Content-Type: application/json' \
  -d "{\"paths\":[\"$SMOKE_DIR/nested\"],\"dryRun\":true,\"confirm\":false,\"confirmPhrase\":\"\"}")
echo "$CLEANUP_JSON" | grep -q 'would_delete'
rm -rf "$SMOKE_DIR"
echo "OK POST /api/scans/{id}/cleanup dry-run"

CODE=$(curl -sf -o /dev/null -w '%{http_code}' "$BASE/")
[[ "$CODE" == "200" ]] || { echo "UI status $CODE" >&2; exit 1; }
echo "OK UI /"

echo "smoke-api: all checks passed"

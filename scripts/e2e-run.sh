#!/usr/bin/env bash
# Build disk-tool, start server, run Cypress E2E, tear down.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$ROOT/bin/disk-tool"
PORT="${E2E_PORT:-18081}"
BASE="http://127.0.0.1:$PORT"
FIXTURE="$ROOT/testdata/e2e-root"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

cd "$ROOT/web"
npm ci
npm install --no-save cypress@15.18.1
npm run build
cd "$ROOT"
rm -rf cmd/disk-tool/static/*
cp -r web/dist/* cmd/disk-tool/static/

go build -o "$BIN" ./cmd/disk-tool

"$BIN" serve --port "$PORT" --no-open &
SERVER_PID=$!

for i in $(seq 1 30); do
  if curl -sf "$BASE/api/roots" >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done
curl -sf "$BASE/api/roots" >/dev/null

cd "$ROOT/web"
export CYPRESS_BASE_URL="$BASE"
# Cypress 15 uses tsx (not ts-node) so Node 20/22 + TypeScript 7 configs work.
npx cypress@15.18.1 run --env "scanRoot=$FIXTURE"

SHOT_COUNT="$(find cypress/screenshots -name '*.png' 2>/dev/null | wc -l | tr -d ' ')"
echo "e2e-run: captured ${SHOT_COUNT} screenshot(s)"

echo "e2e-run: all Cypress specs passed"

#!/usr/bin/env bash
# Smoke test Go sidecar readiness (--port 0 --ready-stdout).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="${1:-$ROOT/bin/disk-tool}"

if [[ ! -x "$BIN" ]]; then
  echo "binary not found: $BIN" >&2
  exit 1
fi

FIFO="$(mktemp -u)"
mkfifo "$FIFO"
"$BIN" serve --no-open --port 0 --ready-stdout >"$FIFO" 2>/dev/null &
PID=$!
trap 'kill -INT "$PID" 2>/dev/null || true; wait "$PID" 2>/dev/null || true; rm -f "$FIFO"' EXIT

if ! read -r -t 10 line <"$FIFO"; then
  echo "no readiness line from sidecar" >&2
  exit 1
fi

if [[ "$line" != disk-tool-ready\ port=* ]]; then
  echo "unexpected readiness line: $line" >&2
  exit 1
fi

PORT="${line#disk-tool-ready port=}"
curl -sf "http://127.0.0.1:${PORT}/api/roots" >/dev/null
echo "sidecar smoke ok (port=${PORT})"

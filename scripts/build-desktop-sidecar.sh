#!/usr/bin/env bash
# Copy the Go sidecar binary into Tauri's externalBin location with the target triple suffix.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="${1:-$ROOT/bin/disk-tool}"
DEST_DIR="$ROOT/desktop/src-tauri/binaries"

if [[ ! -x "$SRC" && ! -f "$SRC" ]]; then
  echo "sidecar source not found: $SRC" >&2
  exit 1
fi

TARGET="$(rustc -vV | sed -n 's/^host: //p')"
EXT=""
if [[ "$TARGET" == *windows* ]]; then
  EXT=".exe"
fi

mkdir -p "$DEST_DIR"
DEST="$DEST_DIR/disk-tool-${TARGET}${EXT}"
cp "$SRC" "$DEST"
chmod +x "$DEST"
echo "Installed sidecar: $DEST"

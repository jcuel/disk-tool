#!/usr/bin/env bash
# Set project version across openspec, CLI, and web package files.
# Usage: bash scripts/bump-version.sh 1.1.0
set -euo pipefail

NEW="${1:?usage: bump-version.sh X.Y.Z}"
if ! [[ "$NEW" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "invalid semver: $NEW" >&2
  exit 1
fi

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG="$ROOT/openspec/config.yaml"
MAIN="$ROOT/cmd/disk-tool/main.go"
PKG="$ROOT/web/package.json"
LOCK="$ROOT/web/package-lock.json"
DESKTOP_PKG="$ROOT/desktop/package.json"
DESKTOP_LOCK="$ROOT/desktop/package-lock.json"
TAURI_CONF="$ROOT/desktop/src-tauri/tauri.conf.json"
CARGO_TOML="$ROOT/desktop/src-tauri/Cargo.toml"

if [[ ! -f "$CONFIG" ]]; then
  echo "missing $CONFIG" >&2
  exit 1
fi

CURRENT="$(grep -E '^version:' "$CONFIG" | sed 's/version:[[:space:]]*//')"
echo "version: $CURRENT -> $NEW"

python3 - "$CONFIG" "$NEW" <<'PY'
import re, sys
path, ver = sys.argv[1], sys.argv[2]
text = open(path, encoding="utf-8").read()
text = re.sub(r"^version: .*$", f"version: {ver}", text, count=1, flags=re.M)
open(path, "w", encoding="utf-8", newline="\n").write(text)
PY

python3 - "$MAIN" "$NEW" <<'PY'
import re, sys
path, ver = sys.argv[1], sys.argv[2]
text = open(path, encoding="utf-8").read()
text = re.sub(
    r'fmt\.Println\("disk-tool [0-9]+\.[0-9]+\.[0-9]+"\)',
    f'fmt.Println("disk-tool {ver}")',
    text,
    count=1,
)
open(path, "w", encoding="utf-8", newline="\n").write(text)
PY

if command -v jq >/dev/null 2>&1; then
  tmp="$(mktemp)"
  jq --arg v "$NEW" '.version = $v' "$PKG" > "$tmp" && mv "$tmp" "$PKG"
  if [[ -f "$LOCK" ]]; then
    tmp="$(mktemp)"
    jq --arg v "$NEW" '.version = $v | .packages[""].version = $v' "$LOCK" > "$tmp" && mv "$tmp" "$LOCK"
  fi
else
  echo "jq required to update package.json" >&2
  exit 1
fi

for f in "$DESKTOP_PKG" "$DESKTOP_LOCK"; do
  if [[ -f "$f" ]]; then
    tmp="$(mktemp)"
    jq --arg v "$NEW" '.version = $v' "$f" > "$tmp" && mv "$tmp" "$f"
  fi
done

if [[ -f "$TAURI_CONF" ]]; then
  python3 - "$TAURI_CONF" "$NEW" <<'PY'
import re, sys
path, ver = sys.argv[1], sys.argv[2]
text = open(path, encoding="utf-8").read()
text = re.sub(r'"version": "[0-9]+\.[0-9]+\.[0-9]+"', f'"version": "{ver}"', text, count=1)
open(path, "w", encoding="utf-8", newline="\n").write(text)
PY
fi

if [[ -f "$CARGO_TOML" ]]; then
  python3 - "$CARGO_TOML" "$NEW" <<'PY'
import re, sys
path, ver = sys.argv[1], sys.argv[2]
text = open(path, encoding="utf-8").read()
text = re.sub(r'^version = "[0-9]+\.[0-9]+\.[0-9]+"', f'version = "{ver}"', text, count=1, flags=re.M)
open(path, "w", encoding="utf-8", newline="\n").write(text)
PY
fi

echo "bumped to $NEW"

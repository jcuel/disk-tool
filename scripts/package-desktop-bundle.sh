#!/usr/bin/env bash
# Copy only user-facing desktop installers from Tauri bundle output.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DEST="${1:-$ROOT/dist/desktop}"
PLATFORM="${2:-linux}"

mkdir -p "$DEST"

if [ "$PLATFORM" = "macos" ]; then
  SEARCH_ROOT="$ROOT/desktop/src-tauri/target/aarch64-apple-darwin/release/bundle"
else
  SEARCH_ROOT="$ROOT/desktop/src-tauri/target/release/bundle"
fi

if [ ! -d "$SEARCH_ROOT" ]; then
  echo "Bundle directory not found: $SEARCH_ROOT" >&2
  find "$ROOT/desktop/src-tauri/target" -type d -name bundle 2>/dev/null || true
  exit 1
fi

mapfile -d '' installers < <(
  find "$SEARCH_ROOT" -type f \( \
    -name '*.AppImage' -o \
    -name '*.deb' -o \
    -name '*.rpm' -o \
    -name '*.dmg' -o \
    -name '*-setup.exe' -o \
    -name '*.msi' \
  \) -print0
)

if [ ${#installers[@]} -eq 0 ]; then
  echo "No desktop installers found under $SEARCH_ROOT" >&2
  find "$SEARCH_ROOT" -type f 2>/dev/null | head -20
  exit 1
fi

for f in "${installers[@]}"; do
  cp "$f" "$DEST/"
  echo "Packaged: $(basename "$f")"
done

ls -la "$DEST"

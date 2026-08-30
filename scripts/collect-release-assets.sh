#!/usr/bin/env bash
# Filter downloaded CI artifacts down to release-quality files only.
set -euo pipefail

SRC="${1:?source dir}"
DEST="${2:?dest dir}"

mkdir -p "$DEST"

# CLI binaries (exact names from build job)
while IFS= read -r -d '' f; do
  cp "$f" "$DEST/"
  echo "CLI: $(basename "$f")"
done < <(
  find "$SRC" -type f \( \
    -name 'disk-tool-linux-amd64' -o \
    -name 'disk-tool-windows-amd64.exe' -o \
    -name 'disk-tool-darwin-amd64' -o \
    -name 'disk-tool-darwin-arm64' \
  \) -print0
)

# Desktop installers only — never AppImage innards (.so, gschema, etc.)
while IFS= read -r -d '' f; do
  cp "$f" "$DEST/"
  echo "Desktop: $(basename "$f")"
done < <(
  find "$SRC" -type f \( \
    -name '*.AppImage' -o \
    -name '*.deb' -o \
    -name '*.rpm' -o \
    -name '*.dmg' -o \
    -name '*-setup.exe' -o \
    -name '*.msi' \
  \) -print0
)

count="$(find "$DEST" -type f | wc -l | tr -d ' ')"
if [ "$count" -eq 0 ]; then
  echo "No release assets collected from $SRC" >&2
  exit 1
fi

echo "Collected $count release file(s):"
ls -la "$DEST"

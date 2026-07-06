#!/usr/bin/env bash
# Compute release version and apply bump-version.sh.
# Used by .github/workflows/release-version.yml on push to master.
#
# Priority:
#   1. RELEASE_VERSION env (explicit X.Y.Z)
#   2. Release-Version: X.Y.Z in commit message / merge body
#   3. Semver bump of openspec/config.yaml (default: minor)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG="$ROOT/openspec/config.yaml"
BUMP_KIND="${BUMP_KIND:-minor}"

current_version() {
  grep -E '^version:' "$CONFIG" | sed 's/version:[[:space:]]*//' | tr -d '\r'
}

semver_bump() {
  local ver="$1" kind="$2"
  local major minor patch
  IFS=. read -r major minor patch <<< "$ver"
  case "$kind" in
    major) echo "$((major + 1)).0.0" ;;
    minor) echo "${major}.$((minor + 1)).0" ;;
    patch) echo "${major}.${minor}.$((patch + 1))" ;;
    *) echo "unknown bump kind: $kind" >&2; exit 1 ;;
  esac
}

resolve_target() {
  if [[ -n "${RELEASE_VERSION:-}" ]]; then
    echo "$RELEASE_VERSION"
    return
  fi
  local msg="${COMMIT_MESSAGE:-}"
  if [[ "$msg" =~ [Rr]elease-[Vv]ersion:[[:space:]]*([0-9]+\.[0-9]+\.[0-9]+) ]]; then
    echo "${BASH_REMATCH[1]}"
    return
  fi
  local cur
  cur="$(current_version)"
  semver_bump "$cur" "$BUMP_KIND"
}

cur="$(current_version)"
target="$(resolve_target)"

if [[ "$target" == "$cur" ]]; then
  echo "version unchanged ($cur); nothing to bump"
  exit 0
fi

echo "release version: $cur -> $target"
bash "$ROOT/scripts/bump-version.sh" "$target"
echo "target_version=$target"

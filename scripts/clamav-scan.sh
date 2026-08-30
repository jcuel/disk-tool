#!/usr/bin/env bash
# Scan release/CI binaries with ClamAV.
# Usage: clamav-scan.sh [dir...]
# Exits 1 when clamscan reports an infection.
set -euo pipefail

REPORT="${CLAMAV_REPORT:-clamav-report.txt}"
SCAN_ROOTS=("$@")
if [ ${#SCAN_ROOTS[@]} -eq 0 ]; then
  SCAN_ROOTS=(".")
fi

if ! command -v clamscan >/dev/null 2>&1; then
  echo "clamscan not found; install clamav package" >&2
  exit 2
fi

: >"$REPORT"
infected=0
scanned=0

scan_one() {
  local file="$1"
  local output rc
  scanned=$((scanned + 1))
  echo "=== $file ===" | tee -a "$REPORT"
  set +e
  output="$(clamscan --no-summary "$file" 2>&1)"
  rc=$?
  set -e
  echo "$output" | tee -a "$REPORT"
  echo >>"$REPORT"
  case "$rc" in
    0) ;;
    1)
      infected=$((infected + 1))
      echo "INFECTED: $file" | tee -a "$REPORT"
      ;;
    *)
      echo "clamscan error on $file (exit $rc)" | tee -a "$REPORT"
      exit "$rc"
      ;;
  esac
}

find_scannable() {
  local root="$1"
  find "$root" -type f \( \
    -name '*.exe' -o \
    -name '*.msi' -o \
    -name '*.dll' -o \
    -name '*.AppImage' -o \
    -name '*.deb' -o \
    -name '*.dmg' -o \
    -name 'disk-tool' -o \
    -name 'disk-tool-*' \
  \) -print 2>/dev/null || true
}

for root in "${SCAN_ROOTS[@]}"; do
  if [ ! -e "$root" ]; then
    echo "Skip missing path: $root" | tee -a "$REPORT"
    continue
  fi
  mapfile -t files < <(find_scannable "$root")
  if [ ${#files[@]} -eq 0 ]; then
    echo "No scannable binaries under $root" | tee -a "$REPORT"
    continue
  fi
  for file in "${files[@]}"; do
    scan_one "$file"
  done
done

{
  echo "Summary: scanned=$scanned infected=$infected"
} | tee -a "$REPORT"

if [ "$scanned" -eq 0 ]; then
  echo "No binaries scanned" >&2
  exit 2
fi

if [ "$infected" -gt 0 ]; then
  exit 1
fi

#!/usr/bin/env bash
# Build disk-tool desktop app: web UI -> Go embed -> sidecar copy -> Tauri bundle.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Building web UI"
(cd web && npm ci && npm run build)

echo "==> Embedding static assets"
rm -rf cmd/disk-tool/static/*
cp -r web/dist/* cmd/disk-tool/static/

echo "==> Building Go sidecar"
go build -ldflags="-s -w" -o bin/disk-tool ./cmd/disk-tool

echo "==> Installing sidecar for Tauri"
bash scripts/build-desktop-sidecar.sh bin/disk-tool

echo "==> Generating desktop icons (ico/icns) if needed"
if [[ ! -f desktop/src-tauri/icons/icon.ico || ! -f desktop/src-tauri/icons/icon.icns ]]; then
  (cd desktop && npm ci && npx tauri icon src-tauri/icons/icon.png)
fi

echo "==> Building Tauri desktop app"
(cd desktop && npm ci && npm run build)

echo "Desktop build complete. Bundles under desktop/src-tauri/target/release/bundle/"

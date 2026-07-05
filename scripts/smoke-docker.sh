#!/usr/bin/env bash
# Docker smoke — Linux CI and WSL (requires Docker)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "== docker compose build =="
docker compose build

echo "== CLI scan smoke =="
docker compose run --rm smoke

echo "== API smoke in container =="
docker compose run --rm smoke-api

echo "smoke-docker: all checks passed"

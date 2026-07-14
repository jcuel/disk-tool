#!/usr/bin/env bash
# Clone hn-mcp for Cursor MCP outreach (optional).
# Usage: bash scripts/outreach/install-hn-mcp.sh [target-dir]
set -euo pipefail

TARGET="${1:-$HOME/tools/hn-mcp}"
if [[ -d "$TARGET/.git" ]]; then
  echo "hn-mcp already at $TARGET"
  exit 0
fi
mkdir -p "$(dirname "$TARGET")"
git clone https://github.com/booklib-ai/hn-mcp.git "$TARGET"
echo "Installed hn-mcp at $TARGET"
echo "Add MCP config — see openspec/changes/outreach-automation/cursor-automation-setup.md"

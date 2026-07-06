#!/usr/bin/env bash
# Export demo fixtures from a real scan of testdata/e2e-root/.
# Run when testdata/e2e-root/ changes: bash scripts/export-demo-fixtures.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$ROOT/bin/disk-tool"
PORT="${EXPORT_PORT:-18082}"
BASE="http://127.0.0.1:$PORT"
FIXTURE="$ROOT/testdata/e2e-root"
OUT="$ROOT/web/src/demo/fixtures.json"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

cd "$ROOT/web"
npm ci --silent
npm run build --silent
cd "$ROOT"
rm -rf cmd/disk-tool/static/*
cp -r web/dist/* cmd/disk-tool/static/
go build -o "$BIN" ./cmd/disk-tool

"$BIN" serve --port "$PORT" --no-open &
SERVER_PID=$!

for i in $(seq 1 50); do
  if curl -sf "$BASE/api/roots" >/dev/null 2>&1; then
    break
  fi
  sleep 0.2
done
curl -sf "$BASE/api/roots" >/dev/null

# Normalize fixture path for cross-platform demo display
DEMO_ROOT="/demo/projects"
scan_id="$(curl -sf -X POST "$BASE/api/scans" \
  -H "Content-Type: application/json" \
  -d "$(python -c "import json,sys; print(json.dumps({'root': sys.argv[1]}))" "$FIXTURE")" \
  | python -c "import json,sys; print(json.load(sys.stdin)['scanId'])")"

for i in $(seq 1 120); do
  status="$(curl -sf "$BASE/api/scans/$scan_id" | python -c "import json,sys; print(json.load(sys.stdin).get('status',''))")"
  if [[ "$status" == "completed" || "$status" == "failed" ]]; then
    break
  fi
  sleep 0.25
done

big_path="$(python -c "
import json, urllib.request
j = json.load(urllib.request.urlopen('$BASE/api/scans/$scan_id'))
for c in j.get('tree', {}).get('children', []):
    if c.get('name') == 'big-dir':
        print(c['path'])
        break
")"
if [[ -n "$big_path" ]]; then
  curl -sf -X POST "$BASE/api/scans/$scan_id/expand" \
    -H "Content-Type: application/json" \
    -d "$(python -c "import json; print(json.dumps({'path': '$big_path', 'depth': 5}))")" >/dev/null
  for i in $(seq 1 60); do
    sleep 0.25
    done="$(curl -sf "$BASE/api/scans/$scan_id" | python -c "
import json,sys
j=json.load(sys.stdin)
def find(n,p):
    if n.get('path')==p: return n
    for c in n.get('children') or []:
        r=find(c,p)
        if r: return r
for c in j.get('tree',{}).get('children',[]):
    if c.get('name')=='big-dir':
        print('yes' if c.get('scanned') else 'no')
        break
")"
    [[ "$done" == "yes" ]] && break
  done
fi

mkdir -p "$(dirname "$OUT")"
python - "$BASE" "$scan_id" "$DEMO_ROOT" "$OUT" <<'PY'
import json, sys, urllib.request

base, scan_id, demo_root, out = sys.argv[1:5]
job = json.load(urllib.request.urlopen(f"{base}/api/scans/{scan_id}"))
disk = json.load(urllib.request.urlopen(
    f"{base}/api/disk?path={urllib.parse.quote(job['root'])}"
))

def remap_path(p: str, old_root: str, new_root: str) -> str:
    old = old_root.replace("\\", "/").rstrip("/")
    new = new_root.rstrip("/")
    p = p.replace("\\", "/")
    if p == old:
        return new
    if p.startswith(old + "/"):
        return new + p[len(old):]
    return p

def remap_node(n, old_root, new_root):
    if not n:
        return n
    out = dict(n)
    out["path"] = remap_path(n["path"], old_root, new_root)
    if "children" in n and n["children"]:
        out["children"] = [remap_node(c, old_root, new_root) for c in n["children"]]
    return out

old_root = job["root"]
job["root"] = demo_root
job["id"] = "demo-scan-1"
if job.get("tree"):
    job["tree"] = remap_node(job["tree"], old_root, demo_root)
    job["tree"]["path"] = demo_root

payload = {
    "demoRoot": demo_root,
    "roots": [demo_root],
    "disk": {
        "path": demo_root,
        "total": disk["total"],
        "used": disk["used"],
        "free": disk["free"],
    },
    "scan": job,
    "maintenancePresets": {
        "presets": [
            {"id": "dev-reclaim", "name": "Dev reclaim", "description": "node_modules, target, build caches", "autoSelect": True},
            {"id": "temp-cleanup", "name": "Temp cleanup", "description": "User temp folders", "autoSelect": False},
        ],
        "matches": [],
    },
    "duplicateGroups": [],
}

with open(out, "w", encoding="utf-8", newline="\n") as f:
    json.dump(payload, f, indent=2)
    f.write("\n")
print(f"Wrote {out}")
PY

echo "export-demo-fixtures: done"

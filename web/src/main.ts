import {
  cancelScan,
  connectEvents,
  expandScan,
  fetchRoots,
  findNode,
  formatBytes,
  getScan,
  needsExpand,
  openPath,
  startScan,
  type ScanJob,
  type ScanNode,
} from "./api";
import { initCharts, pct, renderCharts } from "./charts";
import "./styles.css";

const app = document.getElementById("app")!;

app.innerHTML = `
<header>
  <h1>disk-tool</h1>
  <p class="subtitle">Find where disk space goes — review cleanup candidates and export reports for support</p>
</header>
<section class="controls">
  <label>Scan path</label>
  <div class="path-row">
    <input id="path-input" type="text" placeholder="C:\\ or /home/user" />
    <select id="roots-select"><option value="">Quick pick…</option></select>
    <button id="start-btn">Scan overview</button>
    <button id="cancel-btn" disabled>Cancel</button>
    <button id="expand-btn" disabled>Scan folder deeper</button>
    <button id="export-json" disabled>JSON</button>
    <button id="export-html" disabled>HTML</button>
    <button id="export-ticket" disabled>Support ticket</button>
    <button id="copy-ticket" disabled>Copy report</button>
  </div>
  <div id="progress" class="progress hidden">
    <div class="progress-bar"><div id="progress-fill"></div></div>
    <span id="progress-text"></span>
  </div>
</section>
<main class="layout">
  <section class="panel tree-panel">
    <h2>Folder tree</h2>
    <div id="breadcrumb"></div>
    <table id="tree-table">
      <thead><tr><th>Name</th><th>Size</th><th>%</th><th>Files</th><th></th><th></th></tr></thead>
      <tbody></tbody>
    </table>
  </section>
  <section class="panel charts-panel">
    <h2>Distribution</h2>
    <div id="treemap" class="chart"></div>
    <div id="barchart" class="chart"></div>
  </section>
  <section class="panel insights-panel">
    <h2>Insights</h2>
    <p id="insights-summary" class="hint">Run an overview scan to see where space is used.</p>
    <h3>Cleanup candidates</h3>
    <p class="hint">Leftover project deps, caches, and downloads — review before deleting</p>
    <table id="cleanup-table">
      <thead><tr><th>Type</th><th>Path</th><th>Size</th><th>Hint</th><th></th></tr></thead>
      <tbody></tbody>
    </table>
  </section>
  <section class="panel files-panel">
    <h2>Largest files</h2>
    <p class="hint">Updated as you drill into folders</p>
    <table id="files-table">
      <thead><tr><th>Path</th><th>Size</th><th></th></tr></thead>
      <tbody></tbody>
    </table>
  </section>
</main>
`;

const pathInput = document.getElementById("path-input") as HTMLInputElement;
const rootsSelect = document.getElementById("roots-select") as HTMLSelectElement;
const startBtn = document.getElementById("start-btn") as HTMLButtonElement;
const cancelBtn = document.getElementById("cancel-btn") as HTMLButtonElement;
const expandBtn = document.getElementById("expand-btn") as HTMLButtonElement;
const exportJson = document.getElementById("export-json") as HTMLButtonElement;
const exportHtml = document.getElementById("export-html") as HTMLButtonElement;
const exportTicket = document.getElementById("export-ticket") as HTMLButtonElement;
const copyTicket = document.getElementById("copy-ticket") as HTMLButtonElement;
const progressEl = document.getElementById("progress")!;
const progressFill = document.getElementById("progress-fill")!;
const progressText = document.getElementById("progress-text")!;
const treeBody = document.querySelector("#tree-table tbody")!;
const filesBody = document.querySelector("#files-table tbody")!;
const cleanupBody = document.querySelector("#cleanup-table tbody")!;
const insightsSummary = document.getElementById("insights-summary")!;
const breadcrumb = document.getElementById("breadcrumb")!;

let scanId: string | null = null;
let job: ScanJob | null = null;
let selectedPath: string | null = null;
let ws: WebSocket | null = null;
let expanding = false;

initCharts(
  document.getElementById("treemap")!,
  document.getElementById("barchart")!
);

fetchRoots().then((roots) => {
  for (const r of roots) {
    const opt = document.createElement("option");
    opt.value = r;
    opt.textContent = r;
    rootsSelect.appendChild(opt);
  }
});

rootsSelect.onchange = () => {
  if (rootsSelect.value) pathInput.value = rootsSelect.value;
};

function setProgress(
  label: string,
  ev: { dirsScanned?: number; filesScanned?: number; bytesScanned?: number; currentPath?: string; targetPath?: string }
) {
  progressEl.classList.remove("hidden");
  const pctVal = Math.min(95, (ev.dirsScanned || 0) % 100);
  progressFill.style.width = `${pctVal}%`;
  const where = ev.targetPath || ev.currentPath || "";
  progressText.textContent = `${label}: ${ev.filesScanned || 0} files · ${formatBytes(ev.bytesScanned || 0)} · ${where}`;
}

async function refreshJob() {
  if (!scanId) return;
  job = await getScan(scanId);
  renderUI();
}

function makeOpenBtn(path: string): HTMLButtonElement {
  const btn = document.createElement("button");
  btn.textContent = "Open";
  btn.className = "link-btn";
  btn.type = "button";
  btn.onclick = async (e) => {
    e.stopPropagation();
    if (!scanId) return;
    try {
      await openPath(scanId, path);
    } catch (err) {
      alert(String(err));
    }
  };
  return btn;
}

function renderUI() {
  if (!job?.tree) return;
  const node = selectedPath ? findNode(job.tree, selectedPath) : job.tree;
  if (!node) return;

  expandBtn.disabled = !scanId || expanding || !needsExpand(node);

  breadcrumb.innerHTML = "";
  const up = document.createElement("button");
  up.textContent = "↑ Up";
  up.className = "link-btn";
  up.disabled = node.path === job.tree.path;
  up.onclick = () => {
    const parent = parentPath(node.path, job!.tree!.path);
    selectPath(parent, false);
  };
  breadcrumb.appendChild(up);
  const label = document.createElement("span");
  label.textContent = node.path;
  breadcrumb.appendChild(label);

  treeBody.innerHTML = "";
  const parentSize = node.size || 1;
  const rows = node.scanned && node.children?.length ? node.children : node.children || [];

  if (!node.scanned && node.path !== job.tree.path && rows.length === 0) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td colspan="5" class="muted">Folder not expanded — click <strong>Scan folder deeper</strong></td>`;
    treeBody.appendChild(tr);
  }

  for (const c of rows) {
    const tr = document.createElement("tr");
    tr.className = "clickable";
    const badge = needsExpand(c) ? '<span class="badge">+</span>' : "";
    tr.innerHTML = `<td>${escapeHtml(c.name)}</td><td>${formatBytes(c.size)}</td><td>${pct(c.size, parentSize)}%</td><td>${c.fileCount}</td><td>${badge}</td>`;
    const openTd = document.createElement("td");
    openTd.appendChild(makeOpenBtn(c.path));
    tr.appendChild(openTd);
    tr.onclick = () => selectPath(c.path, needsExpand(c));
    treeBody.appendChild(tr);
  }

  renderCharts(node, (p) => selectPath(p, true));

  renderInsights();

  filesBody.innerHTML = "";
  for (const f of job.largestFiles || []) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${escapeHtml(f.path)}</td><td>${formatBytes(f.size)}</td>`;
    const openTd = document.createElement("td");
    openTd.appendChild(makeOpenBtn(f.path));
    tr.appendChild(openTd);
    filesBody.appendChild(tr);
  }
}

function renderInsights() {
  const ins = job?.insights;
  if (!ins) {
    insightsSummary.textContent = "Run an overview scan to see where space is used.";
    cleanupBody.innerHTML = "";
    return;
  }
  insightsSummary.textContent = ins.summary;
  cleanupBody.innerHTML = "";
  if (ins.cleanupCandidates.length === 0) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td colspan="4" class="muted">No known cleanup patterns yet — drill into Users, Projects, or Downloads</td>`;
    cleanupBody.appendChild(tr);
    return;
  }
  for (const c of ins.cleanupCandidates) {
    const tr = document.createElement("tr");
    tr.className = "clickable";
    tr.innerHTML = `<td><span class="badge-cat">${escapeHtml(c.category)}</span></td><td>${escapeHtml(c.path)}</td><td>${formatBytes(c.size)}</td><td>${escapeHtml(c.hint)}</td>`;
    const openTd = document.createElement("td");
    openTd.appendChild(makeOpenBtn(c.path));
    tr.appendChild(openTd);
    tr.onclick = () => selectPath(c.path, true);
    cleanupBody.appendChild(tr);
  }
}

async function selectPath(path: string, autoExpand: boolean) {
  selectedPath = path;
  renderUI();
  if (autoExpand && scanId && job?.tree) {
    const node = findNode(job.tree, path);
    if (node && needsExpand(node)) {
      await doExpand(path);
    }
  }
}

async function doExpand(path: string) {
  if (!scanId || expanding) return;
  expanding = true;
  expandBtn.disabled = true;
  try {
    await expandScan(scanId, path);
  } catch (e) {
    alert(String(e));
    expanding = false;
    renderUI();
  }
}

function parentPath(current: string, root: string): string {
  if (current === root) return root;
  const sep = current.includes("\\") ? "\\" : "/";
  const parts = current.split(sep);
  parts.pop();
  const p = parts.join(sep) || root;
  return p.length < root.length ? root : p;
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;");
}

startBtn.onclick = async () => {
  const root = pathInput.value.trim();
  if (!root) return;
  startBtn.disabled = true;
  cancelBtn.disabled = false;
  exportJson.disabled = true;
  exportHtml.disabled = true;
  exportTicket.disabled = true;
  copyTicket.disabled = true;
  job = null;
  selectedPath = null;
  try {
    scanId = await startScan(root);
    ws?.close();
    ws = connectEvents(scanId, async (ev) => {
      if (ev.type === "progress" || ev.type === "snapshot") {
        setProgress("Overview", ev);
      }
      if (ev.type === "expand-progress" || ev.type === "expand-started") {
        setProgress("Drill-down", ev);
      }
      if (ev.type === "expand-completed") {
        expanding = false;
        await refreshJob();
      }
      if (ev.type === "expand-error") {
        expanding = false;
        alert(ev.error || "expand failed");
        renderUI();
      }
      if (ev.type === "completed" || ev.type === "cancelled" || ev.type === "error") {
        startBtn.disabled = false;
        cancelBtn.disabled = true;
        if (ev.type === "completed") {
          progressFill.style.width = "100%";
          progressText.textContent = "Overview ready — click a folder to drill down";
        }
        if (scanId && ev.type === "completed") {
          job = await getScan(scanId);
          selectedPath = job.tree?.path || null;
          exportJson.disabled = false;
          exportHtml.disabled = false;
          exportTicket.disabled = false;
          copyTicket.disabled = false;
          renderUI();
        }
      }
    });
  } catch (e) {
    alert(String(e));
    startBtn.disabled = false;
    cancelBtn.disabled = true;
  }
};

cancelBtn.onclick = async () => {
  if (scanId) await cancelScan(scanId);
};

expandBtn.onclick = async () => {
  if (selectedPath) await doExpand(selectedPath);
};

exportJson.onclick = () => {
  if (scanId) window.open(`/api/scans/${scanId}/export?format=json`, "_blank");
};
exportHtml.onclick = () => {
  if (scanId) window.open(`/api/scans/${scanId}/export?format=html`, "_blank");
};
exportTicket.onclick = () => {
  if (scanId) window.open(`/api/scans/${scanId}/export?format=ticket`, "_blank");
};
copyTicket.onclick = async () => {
  const text = job?.insights?.ticketText;
  if (!text) return;
  await navigator.clipboard.writeText(text);
  copyTicket.textContent = "Copied!";
  setTimeout(() => { copyTicket.textContent = "Copy report"; }, 2000);
};

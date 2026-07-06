import {
  cancelScan,
  connectEvents,
  deletePath,
  expandScan,
  fetchDisk,
  fetchMaintenancePresets,
  fetchRoots,
  findDuplicates,
  findNode,
  formatBytes,
  getScan,
  needsExpand,
  openPath,
  reanalyzeInsights,
  runCleanup,
  startScan,
  type CleanupReport,
  type MaintenancePresetMatch,
  type ScanJob,
  type ScanNode,
} from "./api";
import { initCharts, pct, renderCharts, renderDiskPie } from "./charts";
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
<section class="disk-summary panel" id="disk-summary">
  <div class="disk-summary-text">
    <h2>Disk capacity</h2>
    <dl id="disk-stats" class="disk-stats">
      <dt>Volume</dt><dd id="disk-volume">—</dd>
      <dt>Capacity</dt><dd id="disk-capacity">—</dd>
      <dt>Used</dt><dd id="disk-used">—</dd>
      <dt>Free</dt><dd id="disk-free">—</dd>
    </dl>
  </div>
  <div id="disk-pie" class="chart disk-pie-chart"></div>
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
    <p id="charts-hint" class="hint charts-hint hidden">Treemap and bar chart appear after the overview scan completes. Click a segment to drill into that folder.</p>
    <div id="treemap" class="chart"></div>
    <div id="barchart" class="chart"></div>
  </section>
  <section class="panel insights-panel">
    <h2>Insights</h2>
    <p id="insights-summary" class="hint">Run an overview scan to see where space is used.</p>
    <h3>Cleanup candidates</h3>
    <p class="hint">Leftover project deps, caches, and downloads — review before deleting</p>
    <div id="drive-root-banner" class="warning-banner hidden"></div>
    <div id="safety-grid" class="safety-grid hidden"></div>
    <div id="maintenance-presets" class="maintenance-presets hidden"></div>
    <div id="age-controls" class="age-controls hidden">
      <label>Stale age (days) <input type="number" id="age-days" min="1" value="90" /></label>
      <label>Min size (MB) <input type="number" id="min-size-mb" min="1" value="50" /></label>
      <button type="button" id="reanalyze-btn" class="secondary-btn">Refresh insights</button>
      <button type="button" id="duplicates-btn" class="secondary-btn">Find duplicates</button>
    </div>
    <div id="duplicates-panel" class="duplicates-panel hidden"></div>
    <div id="cleanup-toolbar" class="cleanup-toolbar hidden">
      <button type="button" id="select-review-btn" class="secondary-btn">Select safe (review)</button>
      <button type="button" id="select-all-btn" class="secondary-btn">Select all</button>
      <button type="button" id="clear-select-btn" class="secondary-btn">Clear</button>
      <button type="button" id="review-cleanup-btn">Review cleanup…</button>
    </div>
    <div id="cleanup-report-panel" class="cleanup-report hidden"></div>
    <table id="cleanup-table">
      <thead><tr><th></th><th>Type</th><th>Zone</th><th>Risk</th><th>Path</th><th>Size</th><th>Hint</th><th></th></tr></thead>
      <tbody></tbody>
    </table>
  </section>
  <section class="panel files-panel">
    <h2>Largest files</h2>
    <p class="hint">Updated as you drill into folders</p>
    <table id="files-table">
      <thead><tr><th>Path</th><th>Size</th><th>Actions</th></tr></thead>
      <tbody></tbody>
    </table>
  </section>
</main>
<div id="modal-backdrop" class="modal-backdrop hidden" aria-hidden="true">
  <div class="modal" role="dialog" aria-modal="true">
    <h3 id="modal-title"></h3>
    <div id="modal-body"></div>
    <div id="modal-actions" class="modal-actions"></div>
  </div>
</div>
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
const diskVolume = document.getElementById("disk-volume")!;
const diskCapacity = document.getElementById("disk-capacity")!;
const diskUsed = document.getElementById("disk-used")!;
const diskFree = document.getElementById("disk-free")!;
const treeBody = document.querySelector("#tree-table tbody")!;
const filesBody = document.querySelector("#files-table tbody")!;
const cleanupBody = document.querySelector("#cleanup-table tbody")!;
const insightsSummary = document.getElementById("insights-summary")!;
const breadcrumb = document.getElementById("breadcrumb")!;
const cleanupToolbar = document.getElementById("cleanup-toolbar")!;
const cleanupReportPanel = document.getElementById("cleanup-report-panel")!;
const selectReviewBtn = document.getElementById("select-review-btn") as HTMLButtonElement;
const selectAllBtn = document.getElementById("select-all-btn") as HTMLButtonElement;
const clearSelectBtn = document.getElementById("clear-select-btn") as HTMLButtonElement;
const reviewCleanupBtn = document.getElementById("review-cleanup-btn") as HTMLButtonElement;
const modalBackdrop = document.getElementById("modal-backdrop")!;
const modalTitle = document.getElementById("modal-title")!;
const modalBody = document.getElementById("modal-body")!;
const modalActions = document.getElementById("modal-actions")!;
const chartsHint = document.getElementById("charts-hint")!;
const driveRootBanner = document.getElementById("drive-root-banner")!;
const safetyGridEl = document.getElementById("safety-grid")!;
const maintenancePresetsEl = document.getElementById("maintenance-presets")!;
const ageControls = document.getElementById("age-controls")!;
const ageDaysInput = document.getElementById("age-days") as HTMLInputElement;
const minSizeMbInput = document.getElementById("min-size-mb") as HTMLInputElement;
const reanalyzeBtn = document.getElementById("reanalyze-btn") as HTMLButtonElement;
const duplicatesBtn = document.getElementById("duplicates-btn") as HTMLButtonElement;
const duplicatesPanel = document.getElementById("duplicates-panel")!;

let modalBusy = false;

const modalDialog = modalBackdrop.querySelector(".modal") as HTMLDivElement;
modalDialog.addEventListener("click", (e) => e.stopPropagation());

document.addEventListener("keydown", (e) => {
  if (e.key === "Escape" && !modalBackdrop.classList.contains("hidden") && !modalBusy) {
    closeModal();
  }
});
let scanId: string | null = null;
let job: ScanJob | null = null;
let selectedPath: string | null = null;
let ws: WebSocket | null = null;
let expanding = false;
let scanning = false;
const selectedCleanup = new Set<string>();
let pendingDryRun: CleanupReport | null = null;

initCharts(
  document.getElementById("treemap")!,
  document.getElementById("barchart")!,
  document.getElementById("disk-pie")!
);

async function loadDiskSummary(path: string) {
  const root = path.trim();
  if (!root) return;
  try {
    const info = await fetchDisk(root);
    diskVolume.textContent = info.path;
    diskCapacity.textContent = formatBytes(info.total);
    diskUsed.textContent = `${formatBytes(info.used)} (${pct(info.used, info.total)}%)`;
    diskFree.textContent = `${formatBytes(info.free)} (${pct(info.free, info.total)}%)`;
    renderDiskPie(info.free, info.used);
  } catch {
    diskVolume.textContent = root;
    diskCapacity.textContent = "—";
    diskUsed.textContent = "—";
    diskFree.textContent = "—";
  }
}

function pickDefaultRoot(roots: string[]): string {
  const drive = roots.find((r) => /^[A-Za-z]:[\\/]?$/.test(r) || /^[A-Za-z]:\\/.test(r));
  return drive || roots[0] || "";
}

function queryParams(): { root: string | null; noAutoScan: boolean } {
  const params = new URLSearchParams(location.search);
  return {
    root: params.get("root"),
    noAutoScan: params.get("noAutoScan") === "1",
  };
}

fetchRoots().then(async (roots) => {
  for (const r of roots) {
    const opt = document.createElement("option");
    opt.value = r;
    opt.textContent = r;
    rootsSelect.appendChild(opt);
  }
  const { root: queryRoot, noAutoScan } = queryParams();
  const scanRoot = queryRoot || pickDefaultRoot(roots);
  if (scanRoot) {
    pathInput.value = scanRoot;
    if (roots.includes(scanRoot)) {
      rootsSelect.value = scanRoot;
    }
    await loadDiskSummary(scanRoot);
    if (!noAutoScan) {
      void beginScan(scanRoot);
    }
  }
});

rootsSelect.onchange = () => {
  if (rootsSelect.value) {
    pathInput.value = rootsSelect.value;
    void loadDiskSummary(rootsSelect.value);
  }
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

function makeDeleteBtn(path: string, label: string, deletable = true): HTMLButtonElement {
  const btn = document.createElement("button");
  btn.textContent = "Delete";
  btn.className = "link-btn danger";
  btn.type = "button";
  if (!deletable) {
    btn.disabled = true;
    btn.title = "Protected zone — deletion disabled";
    return btn;
  }
  btn.onclick = async (e) => {
    e.stopPropagation();
    if (!scanId) return;
    const phrase = prompt(`Type DELETE to permanently remove:\n\n${path}`);
    if (phrase !== "DELETE") return;
    try {
      await deletePath(scanId, path, "DELETE");
      await refreshJob();
    } catch (err) {
      alert(String(err));
    }
  };
  return btn;
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

function renderLargestFiles() {
  filesBody.innerHTML = "";
  const files = job?.largestFiles || [];
  if (files.length === 0) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td colspan="3" class="muted">No large files indexed yet — drill into folders</td>`;
    filesBody.appendChild(tr);
    return;
  }
  for (const f of files) {
    const tr = document.createElement("tr");
    const pathTd = document.createElement("td");
    pathTd.className = "path-cell";
    pathTd.title = f.path;
    const nameEl = document.createElement("div");
    nameEl.className = "file-name";
    nameEl.textContent = f.name;
    const pathEl = document.createElement("div");
    pathEl.className = "file-path";
    pathEl.textContent = f.path;
    pathTd.appendChild(nameEl);
    pathTd.appendChild(pathEl);
    tr.appendChild(pathTd);
    const sizeTd = document.createElement("td");
    sizeTd.textContent = formatBytes(f.size);
    tr.appendChild(sizeTd);
    const actionsTd = document.createElement("td");
    actionsTd.className = "actions-cell";
    actionsTd.appendChild(makeOpenBtn(f.path));
    actionsTd.appendChild(document.createTextNode(" "));
    actionsTd.appendChild(makeDeleteBtn(f.path, `${f.name}: ${formatBytes(f.size)}`));
    tr.appendChild(actionsTd);
    filesBody.appendChild(tr);
  }
}

function renderEmptyPanels() {
  chartsHint.classList.toggle("hidden", !scanning);
  if (scanning) {
    treeBody.innerHTML =
      `<tr><td colspan="6" class="muted">Overview scan in progress — see progress bar above. Charts fill in when the scan finishes.</td></tr>`;
    breadcrumb.innerHTML = "";
    return;
  }
  chartsHint.classList.add("hidden");
  treeBody.innerHTML =
    `<tr><td colspan="6" class="muted">Click <strong>Scan overview</strong> to index this path</td></tr>`;
  breadcrumb.innerHTML = "";
}

function renderUI() {
  renderLargestFiles();
  if (!job?.tree) {
    renderEmptyPanels();
    return;
  }
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

  renderCharts(node, (path) => {
    const child = (node?.children || []).find((c) => c.path === path);
    selectPath(path, child ? needsExpand(child) : true);
  });

  chartsHint.classList.add("hidden");
  renderInsights();
}

function renderInsights() {
  const ins = job?.insights;
  if (!ins) {
    insightsSummary.textContent = "Run an overview scan to see where space is used.";
    cleanupBody.innerHTML = "";
    cleanupToolbar.classList.add("hidden");
    cleanupReportPanel.classList.add("hidden");
    safetyGridEl.classList.add("hidden");
    maintenancePresetsEl.classList.add("hidden");
    ageControls.classList.add("hidden");
    driveRootBanner.classList.add("hidden");
    return;
  }
  insightsSummary.textContent = ins.summary;
  ageControls.classList.remove("hidden");
  if (job?.insightsConfig?.ageThresholdDays) {
    ageDaysInput.value = String(job.insightsConfig.ageThresholdDays);
  }
  if (job?.insightsConfig?.minSizeBytes) {
    minSizeMbInput.value = String(Math.round(job.insightsConfig.minSizeBytes / (1024 * 1024)));
  }

  if (ins.safetyGrid?.driveRoot) {
    driveRootBanner.classList.remove("hidden");
    driveRootBanner.textContent =
      "Scanning a full drive includes protected OS areas. Prefer your user profile (/home or Users) for cleanup. Protected zones cannot be deleted.";
  } else {
    driveRootBanner.classList.add("hidden");
  }

  renderSafetyGrid(ins);
  void renderMaintenancePresets();
  renderDuplicatesPanel();

  cleanupBody.innerHTML = "";
  const candidates = ins.cleanupCandidates || [];
  cleanupToolbar.classList.toggle("hidden", candidates.length === 0);
  if (candidates.length === 0) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td colspan="8" class="muted">No known cleanup patterns yet — drill into Users, Projects, or Downloads</td>`;
    cleanupBody.appendChild(tr);
    return;
  }
  for (const c of candidates) {
    const tr = document.createElement("tr");
    tr.className = "clickable";
    const checkTd = document.createElement("td");
    const cb = document.createElement("input");
    cb.type = "checkbox";
    cb.checked = selectedCleanup.has(c.path);
    cb.disabled = c.deletable === false;
    cb.onclick = (e) => e.stopPropagation();
    cb.onchange = () => {
      if (cb.checked) selectedCleanup.add(c.path);
      else selectedCleanup.delete(c.path);
    };
    checkTd.appendChild(cb);
    tr.appendChild(checkTd);
    const typeTd = document.createElement("td");
    typeTd.innerHTML = `<span class="badge-cat">${escapeHtml(c.category)}</span>`;
    tr.appendChild(typeTd);
    const zoneTd = document.createElement("td");
    zoneTd.innerHTML = `<span class="badge-zone zone-${escapeHtml(c.zone || "normal")}">${escapeHtml(c.zone || "normal")}</span>`;
    tr.appendChild(zoneTd);
    const riskTd = document.createElement("td");
    riskTd.innerHTML = `<span class="badge-risk risk-${escapeHtml(c.risk)}">${escapeHtml(c.risk)}</span>`;
    tr.appendChild(riskTd);
    const pathTd = document.createElement("td");
    pathTd.textContent = c.path;
    tr.appendChild(pathTd);
    const sizeTd = document.createElement("td");
    sizeTd.textContent = formatBytes(c.size);
    tr.appendChild(sizeTd);
    const hintTd = document.createElement("td");
    hintTd.textContent = c.hint;
    tr.appendChild(hintTd);
    const openTd = document.createElement("td");
    openTd.appendChild(makeOpenBtn(c.path));
    openTd.appendChild(document.createTextNode(" "));
    openTd.appendChild(makeDeleteBtn(c.path, `${c.category}: ${formatBytes(c.size)}`, c.deletable !== false));
    tr.appendChild(openTd);
    tr.onclick = () => selectPath(c.path, true);
    cleanupBody.appendChild(tr);
  }
  renderCleanupReport();
}

function renderSafetyGrid(ins: NonNullable<ScanJob["insights"]>) {
  const grid = ins.safetyGrid;
  if (!grid || Object.keys(grid.zones).length === 0) {
    safetyGridEl.classList.add("hidden");
    return;
  }
  safetyGridEl.classList.remove("hidden");
  const cells = Object.entries(grid.zones)
    .map(
      ([zone, st]) =>
        `<button type="button" class="safety-cell zone-${escapeHtml(zone)}" data-zone="${escapeHtml(zone)}">
          <strong>${escapeHtml(zone)}</strong>
          <span>${st.count} item(s) · ${formatBytes(st.bytes)}</span>
        </button>`
    )
    .join("");
  safetyGridEl.innerHTML = `<h3>Safety grid</h3><p class="hint">Click a zone to filter cleanup candidates</p><div class="safety-cells">${cells}</div>`;
  safetyGridEl.querySelectorAll(".safety-cell").forEach((btn) => {
    btn.addEventListener("click", () => {
      const zone = (btn as HTMLElement).dataset.zone;
      selectedCleanup.clear();
      for (const c of ins.cleanupCandidates || []) {
        if (c.zone === zone && c.deletable !== false) selectedCleanup.add(c.path);
      }
      renderInsights();
    });
  });
}

async function renderMaintenancePresets() {
  if (!scanId || !job?.insights) {
    maintenancePresetsEl.classList.add("hidden");
    return;
  }
  try {
    const data = await fetchMaintenancePresets(scanId);
    maintenancePresetsEl.classList.remove("hidden");
    maintenancePresetsEl.innerHTML = `<h3>Maintenance presets</h3><div class="preset-buttons"></div>`;
    const wrap = maintenancePresetsEl.querySelector(".preset-buttons")!;
    for (const m of data.matches) {
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className = "secondary-btn preset-btn";
      btn.title = m.description;
      btn.textContent = `${m.name} (${m.matchCount})`;
      btn.disabled = m.matchCount === 0;
      btn.onclick = () => applyPreset(m);
      wrap.appendChild(btn);
    }
  } catch {
    maintenancePresetsEl.classList.add("hidden");
  }
}

function applyPreset(match: MaintenancePresetMatch) {
  selectedCleanup.clear();
  for (const p of match.paths) selectedCleanup.add(p);
  renderInsights();
  if (match.id === "cache-review") {
    openReviewModal((job?.insights?.cleanupCandidates || []).filter((c) => match.paths.includes(c.path)));
  } else if (match.matchCount > 0) {
    void startReviewCleanup();
  }
}

function renderDuplicatesPanel() {
  const groups = job?.duplicateGroups || [];
  if (groups.length === 0) {
    duplicatesPanel.classList.add("hidden");
    return;
  }
  duplicatesPanel.classList.remove("hidden");
  const rows = groups
    .slice(0, 10)
    .map(
      (g) =>
        `<tr><td>${escapeHtml(g.hash)}</td><td>${g.files.length}</td><td>${formatBytes(g.wasted)}</td><td>${escapeHtml(g.files.map((f) => f.path).join("; "))}</td></tr>`
    )
    .join("");
  duplicatesPanel.innerHTML = `<h3>Duplicate files</h3><table class="modal-table"><thead><tr><th>Hash</th><th>Copies</th><th>Wasted</th><th>Paths</th></tr></thead><tbody>${rows}</tbody></table>`;
}

reanalyzeBtn.onclick = async () => {
  if (!scanId) return;
  const cfg = {
    ageThresholdDays: parseInt(ageDaysInput.value, 10) || 90,
    minSizeBytes: (parseInt(minSizeMbInput.value, 10) || 50) * 1024 * 1024,
  };
  try {
    reanalyzeBtn.disabled = true;
    await reanalyzeInsights(scanId, cfg);
    await refreshJob();
  } catch (e) {
    alert(String(e));
  } finally {
    reanalyzeBtn.disabled = false;
  }
};

duplicatesBtn.onclick = async () => {
  if (!scanId) return;
  try {
    duplicatesBtn.disabled = true;
    await findDuplicates(scanId);
    await refreshJob();
  } catch (e) {
    alert(String(e));
  } finally {
    duplicatesBtn.disabled = false;
  }
};

function renderCleanupReport() {
  const report = job?.lastCleanupReport;
  if (!report || !scanId) {
    cleanupReportPanel.classList.add("hidden");
    return;
  }
  cleanupReportPanel.classList.remove("hidden");
  const mode = report.dryRun ? "Dry run" : "Executed";
  cleanupReportPanel.innerHTML = `
    <p><strong>Last cleanup (${mode}):</strong> ${report.results.length} item(s), reclaimed ${formatBytes(report.bytesReclaimed)}</p>
    <div class="cleanup-report-actions">
      <button type="button" class="secondary-btn" data-export="cleanup-json">JSON</button>
      <button type="button" class="secondary-btn" data-export="cleanup-html">HTML</button>
      <button type="button" class="secondary-btn" data-export="cleanup-ticket">Ticket</button>
      <button type="button" class="secondary-btn" id="copy-cleanup-report">Copy report</button>
    </div>`;
  cleanupReportPanel.querySelectorAll("[data-export]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const fmt = (btn as HTMLElement).dataset.export;
      window.open(`/api/scans/${scanId}/export?format=${fmt}`, "_blank");
    });
  });
  document.getElementById("copy-cleanup-report")!.onclick = async () => {
    await navigator.clipboard.writeText(report.reportText);
  };
}

function closeModal() {
  modalBackdrop.classList.add("hidden");
  modalBackdrop.setAttribute("aria-hidden", "true");
  modalBody.innerHTML = "";
  modalActions.innerHTML = "";
}

function setModalLoading(message: string) {
  modalTitle.textContent = "Please wait";
  modalBody.innerHTML = `<p class="modal-loading">${escapeHtml(message)}</p>`;
  modalActions.innerHTML = "";
  modalBackdrop.classList.remove("hidden");
  modalBackdrop.setAttribute("aria-hidden", "false");
}

function openModal(title: string, bodyHtml: string, actions: { label: string; primary?: boolean; danger?: boolean; onClick: () => void | Promise<void> }[]) {
  modalTitle.textContent = title;
  modalBody.innerHTML = bodyHtml;
  modalActions.innerHTML = "";
  for (const a of actions) {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.textContent = a.label;
    if (a.primary) btn.className = "primary-btn";
    else if (a.danger) btn.className = "danger-btn";
    else btn.className = "secondary-btn";
    btn.onclick = async (e) => {
      e.stopPropagation();
      if (modalBusy) return;
      modalBusy = true;
      try {
        await a.onClick();
      } finally {
        modalBusy = false;
      }
    };
    modalActions.appendChild(btn);
  }
  modalBackdrop.classList.remove("hidden");
  modalBackdrop.setAttribute("aria-hidden", "false");
}

function selectedCandidates() {
  return (job?.insights?.cleanupCandidates || []).filter((c) => selectedCleanup.has(c.path));
}

function totalSelectedBytes(): number {
  return selectedCandidates().reduce((n, c) => n + c.size, 0);
}

async function startReviewCleanup() {
  if (!scanId) return;
  const items = selectedCandidates();
  if (items.length === 0) {
    alert("Select at least one cleanup candidate.");
    return;
  }
  openReviewModal(items);
}

function openReviewModal(items: ReturnType<typeof selectedCandidates>) {
  const rows = items.map((c) =>
    `<tr><td>${escapeHtml(c.category)}</td><td>${escapeHtml(c.zone || "normal")}</td><td>${escapeHtml(c.path)}</td><td>${formatBytes(c.size)}</td><td>${escapeHtml(c.risk)}</td></tr>`
  ).join("");
  const body = `
    <p>${items.length} item(s), total ${formatBytes(totalSelectedBytes())}</p>
    <div class="modal-table-wrap">
    <table class="modal-table"><thead><tr><th>Type</th><th>Zone</th><th>Path</th><th>Size</th><th>Risk</th></tr></thead><tbody>${rows}</tbody></table>
    </div>`;
  openModal("Review cleanup", body, [
    { label: "Cancel", onClick: closeModal },
    {
      label: "Continue",
      primary: true,
      onClick: async () => {
        if (!scanId) return;
        setModalLoading("Running dry-run preflight…");
        try {
          pendingDryRun = await runCleanup(scanId, {
            paths: items.map((c) => c.path),
            dryRun: true,
            confirm: false,
            confirmPhrase: "",
          });
          showConfirmModal(items);
        } catch (e) {
          alert(String(e));
          openReviewModal(items);
        }
      },
    },
  ]);
}

function showConfirmModal(items: { path: string; size: number; category: string }[]) {
  const dry = pendingDryRun;
  const dryRows = (dry?.results || []).map((r) =>
    `<tr><td>${escapeHtml(r.status)}</td><td>${escapeHtml(r.path)}</td><td>${formatBytes(r.size)}</td><td>${escapeHtml(r.reason || "")}</td></tr>`
  ).join("");
  const body = `
    <p class="warning">This permanently deletes the selected paths. This cannot be undone.</p>
    ${dry ? `<p>Dry-run: ${dry.results.filter((r) => r.status === "would_delete").length} ready, ${dry.results.filter((r) => r.status.startsWith("skipped")).length} skipped</p>` : ""}
    ${dryRows ? `<div class="modal-table-wrap"><table class="modal-table"><thead><tr><th>Status</th><th>Path</th><th>Size</th><th>Reason</th></tr></thead><tbody>${dryRows}</tbody></table></div>` : ""}
    <label class="confirm-check"><input type="checkbox" id="cleanup-reviewed" /> I reviewed these paths</label>
    <label>Type <strong>DELETE</strong> to confirm<br/><input type="text" id="cleanup-phrase" class="modal-input" autocomplete="off" /></label>`;
  openModal("Confirm cleanup", body, [
    { label: "Back", onClick: () => openReviewModal(items as ReturnType<typeof selectedCandidates>) },
    {
      label: "Run cleanup",
      danger: true,
      onClick: async () => {
        const reviewed = (document.getElementById("cleanup-reviewed") as HTMLInputElement).checked;
        const phrase = (document.getElementById("cleanup-phrase") as HTMLInputElement).value.trim();
        if (!reviewed) {
          alert("Check the review confirmation box.");
          return;
        }
        if (phrase !== "DELETE") {
          alert('Type DELETE to confirm.');
          return;
        }
        if (!scanId) return;
        try {
          reviewCleanupBtn.disabled = true;
          setModalLoading("Deleting selected paths…");
          const report = await runCleanup(scanId, {
            paths: items.map((c) => c.path),
            dryRun: false,
            confirm: true,
            confirmPhrase: "DELETE",
          });
          pendingDryRun = null;
          for (const r of report.results) {
            if (r.status === "deleted") selectedCleanup.delete(r.path);
          }
          closeModal();
          await refreshJob();
          alert(`Cleanup complete. Reclaimed ${formatBytes(report.bytesReclaimed)}.`);
        } catch (e) {
          alert(String(e));
          showConfirmModal(items);
        } finally {
          reviewCleanupBtn.disabled = false;
        }
      },
    },
  ]);
}

selectReviewBtn.onclick = () => {
  selectedCleanup.clear();
  for (const c of job?.insights?.cleanupCandidates || []) {
    if (c.risk === "review" && c.deletable !== false) selectedCleanup.add(c.path);
  }
  renderInsights();
};

selectAllBtn.onclick = () => {
  selectedCleanup.clear();
  for (const c of job?.insights?.cleanupCandidates || []) {
    if (c.deletable !== false) selectedCleanup.add(c.path);
  }
  renderInsights();
};

clearSelectBtn.onclick = () => {
  selectedCleanup.clear();
  renderInsights();
};

reviewCleanupBtn.onclick = () => startReviewCleanup();

modalBackdrop.onclick = (e) => {
  if (e.target === modalBackdrop && !modalBusy) closeModal();
};

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
  await beginScan(root);
};

async function applyOverviewCompleted() {
  scanning = false;
  startBtn.disabled = false;
  cancelBtn.disabled = true;
  progressFill.style.width = "100%";
  progressText.textContent = "Overview ready — click a folder to drill down";
  if (!scanId) return;
  job = await getScan(scanId);
  selectedPath = job.tree?.path || null;
  exportJson.disabled = false;
  exportHtml.disabled = false;
  exportTicket.disabled = false;
  copyTicket.disabled = false;
  renderUI();
}

async function beginScan(root: string) {
  pathInput.value = root;
  void loadDiskSummary(root);
  scanning = true;
  startBtn.disabled = true;
  cancelBtn.disabled = false;
  exportJson.disabled = true;
  exportHtml.disabled = true;
  exportTicket.disabled = true;
  copyTicket.disabled = true;
  job = null;
  selectedPath = null;
  renderUI();
  try {
    scanId = await startScan(root);
    ws?.close();
    ws = connectEvents(scanId, async (ev) => {
      if (ev.type === "progress" || ev.type === "snapshot") {
        setProgress("Overview", ev);
        if (ev.type === "snapshot" && ev.status === "completed") {
          await applyOverviewCompleted();
        }
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
        scanning = false;
        startBtn.disabled = false;
        cancelBtn.disabled = true;
        if (ev.type === "completed") {
          await applyOverviewCompleted();
        }
      }
    });
  } catch (e) {
    scanning = false;
    alert(String(e));
    startBtn.disabled = false;
    cancelBtn.disabled = true;
    renderUI();
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

import type { CleanupReport, ScanJob } from "./api";
import { formatBytes } from "./api";

function downloadBlob(content: string, mime: string, filename: string, openInTab = false) {
  const blob = new Blob([content], { type: mime });
  const url = URL.createObjectURL(blob);
  if (openInTab) {
    window.open(url, "_blank");
    setTimeout(() => URL.revokeObjectURL(url), 60_000);
    return;
  }
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function renderHtmlReport(job: ScanJob): string {
  const rootSize = job.tree?.size || 1;
  const rows = (job.tree?.children || [])
    .map((c) => {
      const pct = ((c.size / rootSize) * 100).toFixed(1);
      return `<tr><td>${escapeHtml(c.name)}</td><td>${formatBytes(c.size)}</td><td>${pct}%</td></tr>`;
    })
    .join("");
  const cleanup = (job.insights?.cleanupCandidates || [])
    .map(
      (c) =>
        `<tr><td>${escapeHtml(c.category)}</td><td>${escapeHtml(c.path)}</td><td>${formatBytes(c.size)}</td><td>${escapeHtml(c.hint)}</td></tr>`
    )
    .join("");
  const largest = (job.largestFiles || [])
    .map((f) => `<tr><td>${escapeHtml(f.path)}</td><td>${formatBytes(f.size)}</td></tr>`)
    .join("");
  const summary = job.insights?.summary ? `<p>${escapeHtml(job.insights.summary)}</p>` : "";
  const cleanupSection = cleanup
    ? `<h2>Cleanup candidates</h2><table><tr><th>Category</th><th>Path</th><th>Size</th><th>Hint</th></tr>${cleanup}</table>`
    : "";

  return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>disk-tool scan ${escapeHtml(job.id)}</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#0f1419;color:#e6edf3}
table{border-collapse:collapse;width:100%}th,td{border:1px solid #30363d;padding:.5rem;text-align:left}
th{background:#161b22}
</style></head><body>
<h1>Disk usage report</h1>
${summary}
<p>Root: ${escapeHtml(job.root)} | Indexed files: ${job.filesScanned} | Bytes scanned: ${formatBytes(job.bytesScanned)}</p>
<h2>Top space consumers</h2>
<table><tr><th>Name</th><th>Size</th><th>%</th></tr>${rows}</table>
${cleanupSection}
<h2>Largest files</h2>
<table><tr><th>Path</th><th>Size</th></tr>${largest}</table>
</body></html>`;
}

function renderCleanupHtml(report: CleanupReport, scanId: string): string {
  const rows = report.results
    .map(
      (r) =>
        `<tr><td>${escapeHtml(r.status)}</td><td>${escapeHtml(r.path)}</td><td>${formatBytes(r.size)}</td><td>${escapeHtml(r.reason || "")}</td></tr>`
    )
    .join("");
  return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>disk-tool cleanup ${escapeHtml(scanId)}</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#0f1419;color:#e6edf3}
table{border-collapse:collapse;width:100%}th,td{border:1px solid #30363d;padding:.5rem;text-align:left}
th{background:#161b22}
</style></head><body>
<h1>Cleanup report</h1>
<p>${report.dryRun ? "Dry run" : "Executed"} — reclaimed ${formatBytes(report.bytesReclaimed)}</p>
<table><tr><th>Status</th><th>Path</th><th>Size</th><th>Reason</th></tr>${rows}</table>
</body></html>`;
}

export function exportScanClient(job: ScanJob, format: string): void {
  const id = job.id;
  switch (format) {
    case "json":
      downloadBlob(JSON.stringify(job, null, 2), "application/json", `scan-${id}.json`);
      break;
    case "html":
      downloadBlob(renderHtmlReport(job), "text/html;charset=utf-8", `scan-${id}.html`, true);
      break;
    case "ticket":
      if (!job.insights?.ticketText) {
        alert("No insights available");
        return;
      }
      downloadBlob(job.insights.ticketText, "text/plain;charset=utf-8", `disk-report-${id}.txt`);
      break;
    case "cleanup-json":
      if (!job.lastCleanupReport) {
        alert("No cleanup report available");
        return;
      }
      downloadBlob(JSON.stringify(job.lastCleanupReport, null, 2), "application/json", `cleanup-${id}.json`);
      break;
    case "cleanup-html":
      if (!job.lastCleanupReport) {
        alert("No cleanup report available");
        return;
      }
      downloadBlob(
        renderCleanupHtml(job.lastCleanupReport, id),
        "text/html;charset=utf-8",
        `cleanup-${id}.html`,
        true
      );
      break;
    case "cleanup-ticket":
      if (!job.lastCleanupReport) {
        alert("No cleanup report available");
        return;
      }
      downloadBlob(job.lastCleanupReport.reportText, "text/plain;charset=utf-8", `cleanup-${id}.txt`);
      break;
    default:
      alert(`Unknown export format: ${format}`);
  }
}

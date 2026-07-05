export interface ScanNode {
  name: string;
  path: string;
  size: number;
  fileCount: number;
  isDir?: boolean;
  scanned?: boolean;
  expandable?: boolean;
  children?: ScanNode[];
}

export interface FileEntry {
  path: string;
  name: string;
  size: number;
}

export interface CleanupCandidate {
  category: string;
  path: string;
  size: number;
  hint: string;
  risk: string;
}

export interface TopConsumer {
  name: string;
  path: string;
  size: number;
  pct: number;
}

export interface InsightsReport {
  summary: string;
  topConsumers: TopConsumer[];
  cleanupCandidates: CleanupCandidate[];
  totalReclaimable: number;
  ticketText: string;
}

export interface ScanJob {
  id: string;
  root: string;
  status: string;
  tree?: ScanNode;
  largestFiles?: FileEntry[];
  insights?: InsightsReport;
  dirsScanned: number;
  filesScanned: number;
  bytesScanned: number;
  currentPath: string;
  error?: string;
}

export interface ProgressEvent {
  type: string;
  scanId?: string;
  status?: string;
  targetPath?: string;
  dirsScanned?: number;
  filesScanned?: number;
  bytesScanned?: number;
  currentPath?: string;
  error?: string;
}

export const DEFAULT_DRILL_DEPTH = 5;

export function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  const units = ["KB", "MB", "GB", "TB"];
  let v = n / 1024;
  let i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  return `${v.toFixed(1)} ${units[i]}`;
}

export async function fetchRoots(): Promise<string[]> {
  const r = await fetch("/api/roots");
  const j = await r.json();
  return j.roots as string[];
}

export async function startScan(root: string): Promise<string> {
  const r = await fetch("/api/scans", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ root }),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "scan failed");
  }
  const j = await r.json();
  return j.scanId as string;
}

export async function expandScan(
  id: string,
  path: string,
  depth = DEFAULT_DRILL_DEPTH
): Promise<void> {
  const r = await fetch(`/api/scans/${id}/expand`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path, depth }),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "expand failed");
  }
}

export async function getScan(id: string): Promise<ScanJob> {
  const r = await fetch(`/api/scans/${id}`);
  if (!r.ok) throw new Error("scan not found");
  return r.json();
}

export async function cancelScan(id: string): Promise<void> {
  await fetch(`/api/scans/${id}`, { method: "DELETE" });
}

export function connectEvents(
  id: string,
  onEvent: (ev: ProgressEvent) => void
): WebSocket {
  const proto = location.protocol === "https:" ? "wss" : "ws";
  const ws = new WebSocket(`${proto}://${location.host}/api/scans/${id}/events`);
  ws.onmessage = (msg) => onEvent(JSON.parse(msg.data));
  return ws;
}

export function findNode(root: ScanNode | undefined, path: string): ScanNode | undefined {
  if (!root) return undefined;
  if (root.path === path) return root;
  for (const c of root.children || []) {
    const found = findNode(c, path);
    if (found) return found;
  }
  return undefined;
}

export function needsExpand(node: ScanNode | undefined): boolean {
  if (!node || !node.isDir) return false;
  return node.expandable === true || (node.scanned === false && (node.fileCount > 0 || node.size > 0));
}

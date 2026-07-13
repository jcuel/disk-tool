import type {
  CleanupReport,
  CleanupRequest,
  DiskInfo,
  DockerDiskUsage,
  DockerPruneReport,
  DockerPruneRequest,
  DuplicateGroup,
  InsightsReport,
  MaintenancePresetMatch,
  ScanEvent,
  ScanJob,
} from "./api";

const DEFAULT_DRILL_DEPTH = 5;

export async function fetchRoots(): Promise<string[]> {
  const r = await fetch("/api/roots");
  const j = await r.json();
  return j.roots as string[];
}

export async function fetchDisk(path: string): Promise<DiskInfo> {
  const r = await fetch(`/api/disk?path=${encodeURIComponent(path)}`);
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "disk info failed");
  }
  return r.json();
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

export async function deletePath(id: string, path: string, confirmPhrase = "DELETE"): Promise<void> {
  const r = await fetch(`/api/scans/${id}/delete`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path, confirm: true, confirmPhrase }),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "delete failed");
  }
}

export async function openPath(id: string, path: string): Promise<void> {
  const r = await fetch(`/api/scans/${id}/open`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path }),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "open failed");
  }
}

export async function runCleanup(id: string, req: CleanupRequest): Promise<CleanupReport> {
  const r = await fetch(`/api/scans/${id}/cleanup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "cleanup failed");
  }
  return r.json();
}

export async function cancelScan(id: string): Promise<void> {
  await fetch(`/api/scans/${id}`, { method: "DELETE" });
}

export async function fetchMaintenancePresets(id: string): Promise<{
  presets: { id: string; name: string; description: string; autoSelect: boolean }[];
  matches: MaintenancePresetMatch[];
}> {
  const r = await fetch(`/api/scans/${id}/maintenance-presets`);
  if (!r.ok) throw new Error("presets failed");
  return r.json();
}

export async function fetchDockerStatus(id: string): Promise<{
  usage: DockerDiskUsage;
  dataRoots: { path: string; size: number; hint: string }[];
}> {
  const r = await fetch(`/api/scans/${id}/docker`);
  if (!r.ok) {
    const e = await r.json().catch(() => ({}));
    throw new Error((e as { error?: string }).error || "docker status failed");
  }
  return r.json();
}

export async function dockerPrune(id: string, req: DockerPruneRequest): Promise<DockerPruneReport> {
  const r = await fetch(`/api/scans/${id}/docker/prune`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!r.ok) {
    const e = await r.json().catch(() => ({}));
    throw new Error((e as { error?: string }).error || "docker prune failed");
  }
  return r.json();
}

export async function reanalyzeInsights(
  id: string,
  cfg: { ageThresholdDays: number; minSizeBytes: number }
): Promise<InsightsReport> {
  const r = await fetch(`/api/scans/${id}/reanalyze`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(cfg),
  });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "reanalyze failed");
  }
  return r.json();
}

export async function findDuplicates(id: string): Promise<DuplicateGroup[]> {
  const r = await fetch(`/api/scans/${id}/duplicates`, { method: "POST" });
  if (!r.ok) {
    const e = await r.json();
    throw new Error(e.error || "duplicate scan failed");
  }
  const j = await r.json();
  return j.duplicateGroups as DuplicateGroup[];
}

export function connectEvents(
  id: string,
  onEvent: (ev: ScanEvent) => void
): WebSocket {
  const proto = location.protocol === "https:" ? "wss" : "ws";
  const ws = new WebSocket(`${proto}://${location.host}/api/scans/${id}/events`);
  ws.onmessage = (msg) => onEvent(JSON.parse(msg.data));
  return ws;
}

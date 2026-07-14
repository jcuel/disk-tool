import type {
  CleanupReport,
  DockerDiskUsage,
  DuplicateGroup,
  InsightsReport,
  MaintenancePresetMatch,
  ScanEvent,
  ScanJob,
  ScanNode,
} from "../api";
import raw from "./fixtures.json";

const DEMO_ERROR = "Demo mode — install disk-tool locally to perform this action.";

interface DemoFixtures {
  demoRoot: string;
  roots: string[];
  disk: { path: string; total: number; used: number; free: number };
  scan: ScanJob;
  maintenancePresets: {
    presets: { id: string; name: string; description: string; autoSelect: boolean }[];
    matches: MaintenancePresetMatch[];
  };
  duplicateGroups: DuplicateGroup[];
}

const fixtures = raw as DemoFixtures;
const listeners = new Set<(ev: ScanEvent) => void>();

function cloneScan(): ScanJob {
  return JSON.parse(JSON.stringify(fixtures.scan)) as ScanJob;
}

let job = cloneScan();
let scanId = job.id;
let pendingOverview: (() => void) | null = null;

function delay(ms: number): Promise<void> {
  return new Promise((r) => setTimeout(r, ms));
}

function emit(ev: ScanEvent) {
  for (const fn of listeners) fn(ev);
}

function findNodeInTree(root: ScanNode | undefined, path: string): ScanNode | undefined {
  if (!root) return undefined;
  if (root.path === path) return root;
  for (const c of root.children || []) {
    const found = findNodeInTree(c, path);
    if (found) return found;
  }
  return undefined;
}

class FakeWebSocket {
  static readonly CONNECTING = 0;
  static readonly OPEN = 1;
  static readonly CLOSING = 2;
  static readonly CLOSED = 3;
  readonly CONNECTING = 0;
  readonly OPEN = 1;
  readonly CLOSING = 2;
  readonly CLOSED = 3;
  readyState = FakeWebSocket.OPEN;
  onmessage: ((ev: MessageEvent) => void) | null = null;
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;

  constructor(_url: string) {
    queueMicrotask(() => this.onopen?.());
  }

  close() {
    this.readyState = FakeWebSocket.CLOSED;
    this.onclose?.();
  }

  send() {
    /* demo: no-op */
  }
}

export async function fetchRoots(): Promise<string[]> {
  return [...fixtures.roots];
}

export async function fetchDisk(_path: string) {
  return { ...fixtures.disk };
}

export async function startScan(root: string): Promise<string> {
  if (root !== fixtures.demoRoot && !fixtures.roots.includes(root)) {
    throw new Error(`Demo only supports ${fixtures.demoRoot}`);
  }
  job = cloneScan();
  job.status = "running";
  scanId = job.id;
  await delay(50);
  job.status = "completed";
  pendingOverview = () => {
    emit({ type: "progress", scanId, dirsScanned: 1, filesScanned: 2, bytesScanned: 100, currentPath: fixtures.demoRoot });
    emit({
      type: "snapshot",
      scanId,
      status: "completed",
      dirsScanned: job.dirsScanned,
      filesScanned: job.filesScanned,
      bytesScanned: job.bytesScanned,
    });
    emit({ type: "completed", scanId, status: "completed" });
  };
  return scanId;
}

export async function expandScan(_id: string, path: string, _depth = 5): Promise<void> {
  const node = findNodeInTree(job.tree, path);
  if (!node) throw new Error("path not found");
  emit({ type: "expand-started", scanId, targetPath: path });
  await delay(150);
  if (path === "/demo/projects/small-dir") {
    node.children = [
      {
        name: "tiny.txt",
        path: "/demo/projects/small-dir/tiny.txt",
        size: 6,
        fileCount: 0,
        isDir: false,
        scanned: true,
        expandable: false,
      },
    ];
  }
  node.scanned = true;
  node.expandable = false;
  emit({ type: "expand-completed", scanId, targetPath: path });
}

export async function getScan(id: string): Promise<ScanJob> {
  if (id !== scanId) throw new Error("scan not found");
  return JSON.parse(JSON.stringify(job)) as ScanJob;
}

export async function deletePath(): Promise<void> {
  throw new Error(DEMO_ERROR);
}

export async function openPath(): Promise<void> {
  throw new Error(DEMO_ERROR);
}

export async function runCleanup(): Promise<CleanupReport> {
  throw new Error(DEMO_ERROR);
}

export async function cancelScan(): Promise<void> {
  job.status = "cancelled";
  emit({ type: "cancelled", scanId });
}

export async function fetchMaintenancePresets(_id: string) {
  return JSON.parse(JSON.stringify(fixtures.maintenancePresets));
}

export async function fetchDockerStatus(_id: string): Promise<{
  usage: DockerDiskUsage;
  dataRoots: { path: string; size: number; hint: string }[];
}> {
  const usage: DockerDiskUsage = {
    available: true,
    daemonOk: true,
    imagesSize: 13.08 * 1024 * 1024 * 1024,
    imagesReclaimable: 10.26 * 1024 * 1024 * 1024,
    containersSize: 3.665 * 1024 * 1024,
    containersReclaimable: 3.66 * 1024 * 1024,
    volumesSize: 9.039 * 1024 * 1024 * 1024,
    volumesReclaimable: 939.3 * 1024 * 1024,
    buildCacheSize: 2.35 * 1024 * 1024 * 1024,
    buildCacheReclaimable: 2.35 * 1024 * 1024 * 1024,
    reclaimable: 13543722188,
    rawDf:
      "TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE\n" +
      "Images          34        7         13.08GB   10.26GB (78%)\n" +
      "Containers      7         2         3.665MB   3.66MB (99%)\n" +
      "Local Volumes   11        5         9.039GB   939.3MB (10%)\n" +
      "Build Cache     116       0         2.35GB    2.35GB\n",
  };
  return {
    usage,
    dataRoots: [
      {
        path: "/Users/demo/Library/Containers/com.docker.docker/Data",
        size: 47231042854,
        hint: "Docker data root — reclaim via docker system prune or Docker Desktop; do not delete this folder",
      },
    ],
  };
}

export async function dockerPrune(
  _id: string,
  req: { dryRun: boolean; confirm: boolean; confirmPhrase: string }
): Promise<{ dryRun: boolean; reclaimable: number; output?: string; error?: string }> {
  if (req.dryRun) {
    return {
      dryRun: true,
      reclaimable: 13543722188,
      output:
        "Demo dry-run only — install disk-tool locally to run docker system prune -af (volumes kept).",
    };
  }
  throw new Error(DEMO_ERROR);
}

export async function reanalyzeInsights(_id: string): Promise<InsightsReport> {
  if (!job.insights) throw new Error("no insights");
  return JSON.parse(JSON.stringify(job.insights)) as InsightsReport;
}

export async function findDuplicates(_id: string): Promise<DuplicateGroup[]> {
  return JSON.parse(JSON.stringify(fixtures.duplicateGroups)) as DuplicateGroup[];
}

export function connectEvents(id: string, onEvent: (ev: ScanEvent) => void): WebSocket {
  const handler = (ev: ScanEvent) => {
    if (!ev.scanId || ev.scanId === id) onEvent(ev);
  };
  listeners.add(handler);
  if (pendingOverview && id === scanId) {
    const play = pendingOverview;
    pendingOverview = null;
    queueMicrotask(() => play());
  }
  const ws = new FakeWebSocket("");
  ws.onmessage = (msg) => onEvent(JSON.parse(msg.data as string));
  const origClose = ws.close.bind(ws);
  ws.close = () => {
    listeners.delete(handler);
    origClose();
  };
  return ws as unknown as WebSocket;
}

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
  zone?: string;
  deletable?: boolean;
}

export interface SafetyZoneStats {
  count: number;
  bytes: number;
}

export interface SafetyGrid {
  zones: Record<string, SafetyZoneStats>;
  driveRoot: boolean;
  protectedBytes: number;
}

export interface DuplicateGroup {
  hash: string;
  files: FileEntry[];
  wasted: number;
}

export interface MaintenancePresetMatch {
  id: string;
  name: string;
  description: string;
  matchCount: number;
  matchBytes: number;
  paths: string[];
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
  safetyGrid?: SafetyGrid;
}

export interface ScanJob {
  id: string;
  root: string;
  status: string;
  tree?: ScanNode;
  largestFiles?: FileEntry[];
  insights?: InsightsReport;
  lastCleanupReport?: CleanupReport;
  duplicateGroups?: DuplicateGroup[];
  insightsConfig?: { ageThresholdDays: number; minSizeBytes: number };
  dirsScanned: number;
  filesScanned: number;
  bytesScanned: number;
  currentPath: string;
  error?: string;
}

export interface CleanupItemResult {
  path: string;
  size: number;
  category?: string;
  status: string;
  reason?: string;
}

export interface CleanupReport {
  dryRun: boolean;
  startedAt: string;
  finishedAt: string;
  totalRequested: number;
  bytesReclaimed: number;
  results: CleanupItemResult[];
  reportText: string;
}

export interface CleanupRequest {
  paths: string[];
  dryRun: boolean;
  confirm: boolean;
  confirmPhrase: string;
}

export interface DockerDiskUsage {
  available: boolean;
  daemonOk: boolean;
  error?: string;
  imagesSize: number;
  imagesReclaimable: number;
  containersSize: number;
  containersReclaimable: number;
  volumesSize: number;
  volumesReclaimable: number;
  buildCacheSize: number;
  buildCacheReclaimable: number;
  reclaimable: number;
  rawDf?: string;
}

export interface DockerPruneReport {
  dryRun: boolean;
  reclaimable: number;
  output?: string;
  error?: string;
}

export interface DockerPruneRequest {
  dryRun: boolean;
  confirm: boolean;
  confirmPhrase: string;
}

export interface ScanEvent {
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

export interface DiskInfo {
  path: string;
  total: number;
  free: number;
  used: number;
}

export const DEFAULT_DRILL_DEPTH = 5;

export const isDemoMode = import.meta.env.VITE_DEMO_MODE === "true";

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

import * as live from "./api-live";
import * as mock from "./demo/mock-api";

const impl = isDemoMode ? mock : live;

export const fetchRoots = impl.fetchRoots;
export const fetchDisk = impl.fetchDisk;
export const startScan = impl.startScan;
export const expandScan = impl.expandScan;
export const getScan = impl.getScan;
export const deletePath = impl.deletePath;
export const openPath = impl.openPath;
export const runCleanup = impl.runCleanup;
export const cancelScan = impl.cancelScan;
export const fetchMaintenancePresets = impl.fetchMaintenancePresets;
export const fetchDockerStatus = impl.fetchDockerStatus;
export const dockerPrune = impl.dockerPrune;
export const reanalyzeInsights = impl.reanalyzeInsights;
export const findDuplicates = impl.findDuplicates;
export const connectEvents = impl.connectEvents;

import * as echarts from "echarts";
import type { ScanNode } from "./api";
import { formatBytes } from "./api";

let treemapChart: echarts.ECharts | null = null;
let barChart: echarts.ECharts | null = null;
let diskPieChart: echarts.ECharts | null = null;

type ChartChild = { name: string; path: string; size: number };

let lastTreemapChildren: ChartChild[] = [];
let lastBarChildren: ChartChild[] = [];

export function initCharts(treemapEl: HTMLElement, barEl: HTMLElement, diskPieEl?: HTMLElement) {
  treemapChart = echarts.init(treemapEl);
  barChart = echarts.init(barEl);
  if (diskPieEl) {
    diskPieChart = echarts.init(diskPieEl);
  }
  window.addEventListener("resize", () => {
    treemapChart?.resize();
    barChart?.resize();
    diskPieChart?.resize();
  });
}

export function renderDiskPie(free: number, used: number) {
  if (!diskPieChart) return;
  diskPieChart.setOption({
    tooltip: {
      trigger: "item",
      formatter: (p: { name: string; value: number; percent: number }) =>
        `${p.name}<br/>${formatBytes(p.value)} (${p.percent.toFixed(1)}%)`,
    },
    legend: { bottom: 0, textStyle: { color: "#8b949e" } },
    series: [
      {
        type: "pie",
        radius: ["45%", "70%"],
        center: ["50%", "45%"],
        avoidLabelOverlap: true,
        label: { show: false },
        data: [
          { name: "Free", value: free, itemStyle: { color: "#238636" } },
          { name: "Used", value: used, itemStyle: { color: "#388bfd" } },
        ],
      },
    ],
  });
}

function resolveChildPath(
  params: {
    data?: unknown;
    name?: string;
    value?: unknown;
    dataIndex?: number;
  },
  byName: Map<string, string>,
  byIndex?: ChartChild[]
): string | null {
  const data = params.data as { path?: string; name?: string } | undefined;
  if (data?.path) return data.path;
  const label =
    params.name ??
    data?.name ??
    (typeof params.value === "string" ? params.value : undefined);
  if (label && byName.has(label)) return byName.get(label)!;
  if (params.dataIndex != null && byIndex?.[params.dataIndex]) {
    return byIndex[params.dataIndex].path;
  }
  return null;
}

function pickChildPath(
  chart: echarts.ECharts,
  params: {
    componentType?: string;
    data?: unknown;
    name?: string;
    value?: unknown;
    dataIndex?: number;
    event?: { offsetX?: number; offsetY?: number };
  },
  byName: Map<string, string>,
  byIndex: ChartChild[]
): string | null {
  if (params.componentType === "yAxis") {
    const label = typeof params.value === "string" ? params.value : String(params.value ?? "");
    if (label && byName.has(label)) return byName.get(label)!;
  }

  const fromParams = resolveChildPath(params, byName, byIndex);
  if (fromParams) return fromParams;

  const idx = params.dataIndex;
  if (idx != null && byIndex[idx]) return byIndex[idx].path;

  const ox = params.event?.offsetX;
  const oy = params.event?.offsetY;
  if (ox == null || oy == null) return null;
  try {
    const raw = chart.convertFromPixel({ seriesIndex: 0 }, [ox, oy]);
    const dataIdx = Array.isArray(raw) ? Number(raw[raw.length - 1]) : Number(raw);
    if (!Number.isNaN(dataIdx) && byIndex[dataIdx]) return byIndex[dataIdx].path;
  } catch {
    /* chart not ready */
  }
  return null;
}

function bindDrillClick(
  chart: echarts.ECharts,
  byName: Map<string, string>,
  byIndex: ChartChild[],
  onSelect: (path: string) => void
) {
  chart.off("click");
  chart.on("click", (p) => {
    const path = pickChildPath(chart, p, byName, byIndex);
    if (path) onSelect(path);
  });
}

export function renderCharts(
  node: ScanNode | undefined,
  onSelect: (path: string) => void
) {
  if (!treemapChart || !barChart) return;
  const children: ChartChild[] = (node?.children || []).map((c) => ({
    name: c.name,
    path: c.path,
    size: c.size,
  }));
  const pathByName = new Map(children.map((c) => [c.name, c.path]));
  const top = [...children].sort((a, b) => b.size - a.size).slice(0, 15);
  lastTreemapChildren = children;
  lastBarChildren = top;

  treemapChart.setOption(
    {
      tooltip: {
        formatter: (p: { name: string; value: number }) =>
          `${p.name}<br/>${formatBytes(p.value)}`,
      },
      series: [
        {
          type: "treemap",
          roam: false,
          nodeClick: false,
          triggerEvent: true,
          data: children.map((c) => ({ name: c.name, value: c.size, path: c.path })),
          label: { show: true, formatter: "{b}" },
          levels: [{ itemStyle: { borderWidth: 1, gapWidth: 2 } }],
          emphasis: { itemStyle: { borderColor: "#58a6ff" } },
        },
      ],
    },
    true
  );

  barChart.setOption(
    {
      tooltip: {
        trigger: "item",
        formatter: (p: { name: string; value: number }) =>
          `${p.name}<br/>${formatBytes(p.value)}`,
      },
      grid: { left: 120, right: 20, top: 10, bottom: 30 },
      xAxis: { type: "value", axisLabel: { formatter: (v: number) => formatBytes(v) } },
      yAxis: {
        type: "category",
        triggerEvent: true,
        data: top.map((c) => c.name),
        inverse: true,
        axisLabel: { width: 110, overflow: "truncate" },
      },
      series: [
        {
          type: "bar",
          triggerEvent: true,
          data: top.map((c) => ({
            value: c.size,
            path: c.path,
            name: c.name,
          })),
        },
      ],
    },
    true
  );

  bindDrillClick(treemapChart, pathByName, lastTreemapChildren, onSelect);
  bindDrillClick(barChart, pathByName, lastBarChildren, onSelect);
}

export function pct(size: number, parent: number): string {
  if (parent <= 0) return "0";
  return ((size / parent) * 100).toFixed(1);
}

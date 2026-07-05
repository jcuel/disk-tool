import * as echarts from "echarts";
import type { ScanNode } from "./api";
import { formatBytes } from "./api";

let treemapChart: echarts.ECharts | null = null;
let barChart: echarts.ECharts | null = null;

export function initCharts(treemapEl: HTMLElement, barEl: HTMLElement) {
  treemapChart = echarts.init(treemapEl);
  barChart = echarts.init(barEl);
  window.addEventListener("resize", () => {
    treemapChart?.resize();
    barChart?.resize();
  });
}

export function renderCharts(
  node: ScanNode | undefined,
  onSelect: (path: string) => void
) {
  if (!treemapChart || !barChart) return;
  const children = node?.children || [];
  const parentSize = node?.size || 1;

  const data = children.map((c) => ({
    name: c.name,
    value: c.size,
    path: c.path,
  }));

  treemapChart.setOption({
    tooltip: {
      formatter: (p: { name: string; value: number }) =>
        `${p.name}<br/>${formatBytes(p.value)}`,
    },
    series: [
      {
        type: "treemap",
        roam: false,
        nodeClick: false,
        data,
        label: { show: true, formatter: "{b}" },
        levels: [{ itemStyle: { borderWidth: 1, gapWidth: 2 } }],
      },
    ],
  });

  const top = [...children].sort((a, b) => b.size - a.size).slice(0, 15);
  barChart.setOption({
    tooltip: {
      trigger: "axis",
      formatter: (params: { name: string; value: number }[]) => {
        const p = params[0];
        return `${p.name}<br/>${formatBytes(p.value as number)}`;
      },
    },
    grid: { left: 120, right: 20, top: 10, bottom: 30 },
    xAxis: { type: "value", axisLabel: { formatter: (v: number) => formatBytes(v) } },
    yAxis: {
      type: "category",
      data: top.map((c) => c.name),
      inverse: true,
    },
    series: [
      {
        type: "bar",
        data: top.map((c) => ({
          value: c.size,
          itemStyle: {},
        })),
      },
    ],
  });

  treemapChart.off("click");
  barChart.off("click");
  treemapChart.on("click", (p) => {
    const d = p.data as { path?: string } | undefined;
    if (d?.path) onSelect(d.path);
  });
  barChart.on("click", (p) => {
    const idx = p.dataIndex ?? -1;
    const item = top[idx];
    if (item?.path) onSelect(item.path);
  });
}

export function pct(size: number, parent: number): string {
  if (parent <= 0) return "0";
  return ((size / parent) * 100).toFixed(1);
}

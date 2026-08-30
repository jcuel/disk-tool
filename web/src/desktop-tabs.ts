import { resizeCharts } from "./charts";

export type DesktopTab = "browse" | "files" | "insights";

export function isDesktopApp(): boolean {
  return !!(window as unknown as { __TAURI_INTERNALS__?: unknown }).__TAURI_INTERNALS__;
}

let activeTab: DesktopTab = "browse";

export function getActiveDesktopTab(): DesktopTab {
  return activeTab;
}

export function initDesktopTabs(): void {
  if (!isDesktopApp()) return;

  document.body.classList.add("desktop-app");
  document.body.dataset.activeTab = activeTab;

  const layout = document.querySelector("main.layout");
  if (!layout) return;

  const tabBar = document.createElement("nav");
  tabBar.id = "desktop-tab-bar";
  tabBar.className = "tab-bar";
  tabBar.setAttribute("role", "tablist");
  tabBar.innerHTML = `
    <button type="button" class="tab-btn active" data-tab="browse" role="tab" aria-selected="true">Browse</button>
    <button type="button" class="tab-btn" data-tab="files" role="tab" aria-selected="false">Files</button>
    <button type="button" class="tab-btn" data-tab="insights" role="tab" aria-selected="false">Insights</button>
  `;
  layout.parentElement?.insertBefore(tabBar, layout);

  layout.querySelectorAll("[data-desktop-tab]").forEach((el) => {
    el.classList.add("desktop-tab-panel");
  });

  tabBar.querySelectorAll(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const tab = btn.getAttribute("data-tab") as DesktopTab;
      switchDesktopTab(tab);
    });
  });
  applyDesktopTab();
}

export function switchDesktopTab(tab: DesktopTab): void {
  if (!isDesktopApp()) return;
  activeTab = tab;
  document.body.dataset.activeTab = tab;
  applyDesktopTab();
}

function applyDesktopTab(): void {
  const tab = activeTab;
  document.querySelectorAll(".tab-btn").forEach((btn) => {
    const on = btn.getAttribute("data-tab") === tab;
    btn.classList.toggle("active", on);
    btn.setAttribute("aria-selected", on ? "true" : "false");
  });
  document.querySelectorAll(".desktop-tab-panel").forEach((el) => {
    const panelTab = el.getAttribute("data-desktop-tab");
    el.classList.toggle("desktop-tab-active", panelTab === tab);
  });
  if (tab === "browse") {
    requestAnimationFrame(() => resizeCharts());
  }
}

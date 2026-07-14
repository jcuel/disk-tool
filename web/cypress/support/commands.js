export function scanRoot() {
  const root = Cypress.env("scanRoot");
  if (!root) {
    throw new Error("Cypress env scanRoot is required (set by scripts/e2e-run.sh)");
  }
  return root;
}

function visitWithScan(root, noAutoScan = false) {
  const q = new URLSearchParams({ root });
  if (noAutoScan) q.set("noAutoScan", "1");
  cy.visit(`/?${q.toString()}`);
}

function waitForOverviewReady() {
  cy.get("#tree-table tbody tr.clickable", { timeout: 60000 }).should("have.length.at.least", 1);
  cy.get("#insights-summary").should("not.contain", "Run an overview scan");
}

function sanitizeScreenshotName(name) {
  return name.replace(/[^\w-]+/g, "-");
}

Cypress.Commands.add("visitWithScan", (root, noAutoScan = false) => {
  visitWithScan(root, noAutoScan);
});

Cypress.Commands.add("waitForOverviewReady", () => {
  waitForOverviewReady();
});

Cypress.Commands.add("captureStep", (name, options = {}) => {
  const safe = sanitizeScreenshotName(name);
  cy.screenshot(safe, { capture: "viewport", overwrite: true, ...options });
});

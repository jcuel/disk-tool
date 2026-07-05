/// <reference types="cypress" />

export function scanRoot(): string {
  const root = Cypress.env("scanRoot") as string | undefined;
  if (!root) {
    throw new Error("Cypress env scanRoot is required (set by scripts/e2e-run.sh)");
  }
  return root;
}

function visitWithScan(root: string, noAutoScan = false): void {
  const q = new URLSearchParams({ root });
  if (noAutoScan) q.set("noAutoScan", "1");
  cy.visit(`/?${q.toString()}`);
}

function waitForOverviewReady(): void {
  cy.get("#progress-text", { timeout: 60000 }).should("contain", "Overview ready");
  cy.get("#tree-table tbody tr.clickable", { timeout: 10000 }).should("have.length.at.least", 1);
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitWithScan(root: string, noAutoScan?: boolean): Chainable<void>;
      waitForOverviewReady(): Chainable<void>;
    }
  }
}

Cypress.Commands.add("visitWithScan", (root: string, noAutoScan = false) => {
  visitWithScan(root, noAutoScan);
});

Cypress.Commands.add("waitForOverviewReady", () => {
  waitForOverviewReady();
});

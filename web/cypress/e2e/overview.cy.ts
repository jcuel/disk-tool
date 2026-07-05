import { scanRoot } from "../support/commands";

describe("overview scan", () => {
  it("completes fixture scan and shows tree rows", () => {
    const root = scanRoot();
    cy.visitWithScan(root);
    cy.waitForOverviewReady();
    cy.get("#insights-summary").should("not.contain", "Run an overview scan");
    cy.get("#tree-table tbody tr.clickable").first().should("contain", "big-dir");
    cy.captureStep("01-overview-ready");
  });

  it("starts scan manually when noAutoScan is set", () => {
    const root = scanRoot();
    cy.visitWithScan(root, true);
    cy.get("#path-input").should("have.value", root);
    cy.get("#start-btn").should("not.be.disabled").click();
    cy.waitForOverviewReady();
    cy.captureStep("02-manual-scan-start");
  });
});

import { scanRoot } from "../support/commands.js";

describe("tree drill-down", () => {
  beforeEach(() => {
    cy.visitWithScan(scanRoot());
    cy.waitForOverviewReady();
  });

  it("updates breadcrumb when a folder row is clicked", () => {
    cy.get("#tree-table tbody tr.clickable").contains("td", "big-dir").parent().click();
    cy.get("#breadcrumb span").invoke("text").should("match", /big-dir/);
    cy.captureStep("03-drill-down-big-dir");
  });
});

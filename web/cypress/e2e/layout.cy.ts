import { scanRoot } from "../support/commands";

describe("layout", () => {
  beforeEach(() => {
    cy.visitWithScan(scanRoot());
    cy.waitForOverviewReady();
  });

  it("shows main panels and insights spans layout width", () => {
    cy.contains("h2", "Folder tree").should("be.visible");
    cy.contains("h2", "Distribution").should("be.visible");
    cy.contains("h2", "Largest files").should("be.visible");
    cy.contains("h2", "Insights").should("be.visible");

    cy.get(".layout").then(($layout) => {
      const layoutW = $layout[0].getBoundingClientRect().width;
      cy.get(".insights-panel").then(($ins) => {
        const insW = $ins[0].getBoundingClientRect().width;
        expect(insW).to.be.greaterThan(layoutW * 0.9);
      });
    });
    cy.captureStep("04-layout-full", { capture: "fullPage" });
  });
});

import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: process.env.CYPRESS_BASE_URL || "http://127.0.0.1:18081",
    // Plain JS specs: Cypress webpack/ts-loader is incompatible with TypeScript 7 (fileExists).
    supportFile: "cypress/support/e2e.js",
    specPattern: "cypress/e2e/**/*.cy.js",
    viewportWidth: 1280,
    viewportHeight: 800,
    screenshotsFolder: "cypress/screenshots",
    video: false,
    screenshotOnRunFailure: true,
    defaultCommandTimeout: 15000,
    requestTimeout: 15000,
    setupNodeEvents() {},
  },
});

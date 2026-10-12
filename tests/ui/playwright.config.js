import { defineConfig } from "@playwright/test";

const demoUrl = process.env.DEMO_URL || "http://localhost:8377";

// Lightpanda stacks the DOM without layout, so the demo sits at large document
// coordinates. Playwright refuses to click an element outside the viewport,
// and a 100,000 pixel viewport still left demo elements outside it.
const lightpandaClickableViewport = { width: 1280, height: 10_000_000 };

const realEngineDiagnostics = {
  screenshot: "only-on-failure",
  trace: "retain-on-failure",
};

export default defineConfig({
  testDir: "./specs",
  fullyParallel: true,
  workers: 3,
  reporter: [["list"]],
  use: {
    baseURL: demoUrl,
  },
  projects: [
    {
      name: "lightpanda",
      grepInvert: /@browser/,
      use: {
        viewport: lightpandaClickableViewport,
        screenshot: "off",
        trace: "off",
      },
    },
    {
      name: "chromium",
      use: { browserName: "chromium", ...realEngineDiagnostics },
    },
    {
      name: "firefox",
      use: { browserName: "firefox", ...realEngineDiagnostics },
    },
    {
      name: "webkit",
      use: { browserName: "webkit", ...realEngineDiagnostics },
    },
  ],
});

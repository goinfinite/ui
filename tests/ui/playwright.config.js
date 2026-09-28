import { defineConfig } from "@playwright/test";

const demoUrl = process.env.DEMO_URL || "http://localhost:8377";

export default defineConfig({
  testDir: "./specs",
  fullyParallel: true,
  workers: 3,
  reporter: [["list"]],
  use: {
    baseURL: demoUrl,
  },
  projects: [
    { name: "chromium", use: { browserName: "chromium" } },
    { name: "firefox", use: { browserName: "firefox" } },
    { name: "webkit", use: { browserName: "webkit" } },
  ],
});

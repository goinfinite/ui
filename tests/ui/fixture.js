import { test as base, chromium, expect } from "@playwright/test";

const lightpandaEndpoint = process.env.UI_TEST_LIGHTPANDA_CDP_URL;

// Lightpanda speaks CDP, not the Playwright wire protocol, so the config
// connectOptions cannot reach it. The worker browser connects over CDP instead.
const test = lightpandaEndpoint
  ? base.extend({
      browser: [
        // biome-ignore lint/correctness/noEmptyPattern: Playwright fixture API
        async ({}, use) => {
          const browser = await chromium.connectOverCDP(lightpandaEndpoint);
          await use(browser);
          await browser.close();
        },
        { scope: "worker" },
      ],
    })
  : base;

export { expect, test };

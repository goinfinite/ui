import { writeFileSync } from "node:fs";
import { expect, test } from "@playwright/test";

const resultsPath = new URL(
  "../../../performance-results.json",
  import.meta.url,
);
const warmupRunCount = 1;
const measuredRunCount = 3;
const expandedHeightPx = 432;

function median(samples) {
  const sorted = [...samples].sort((left, right) => left - right);
  return Math.round(sorted[Math.floor(sorted.length / 2)]);
}

async function measureSamples(measureOnce) {
  const samples = [];
  for (let run = 0; run < warmupRunCount + measuredRunCount; run++) {
    const elapsed = await measureOnce();
    if (run >= warmupRunCount) samples.push(elapsed);
  }
  return median(samples);
}

async function measureExpandLatency(page) {
  return page.evaluate(async (targetHeightPx) => {
    const textarea = document.querySelector("#text-area-demo textarea");
    const icons = textarea.closest("fieldset").querySelector("div.absolute");
    icons.style.setProperty("display", "flex", "important");
    const toggle = icons.querySelector("button");
    const start = performance.now();
    toggle.click();
    await new Promise((resolve) => {
      const poll = () => {
        if (
          Math.abs(textarea.getBoundingClientRect().height - targetHeightPx) < 1
        )
          resolve();
        else requestAnimationFrame(poll);
      };
      poll();
    });
    const elapsed = performance.now() - start;
    toggle.click();
    await new Promise((resolve) => setTimeout(resolve, 400));
    icons.style.removeProperty("display");
    return elapsed;
  }, expandedHeightPx);
}

async function measureDropdownOpenLatency(page) {
  return page.evaluate(async () => {
    const trigger = document.querySelector("#select-input-demo .group.flex");
    const list = document.querySelector("#select-input-demo ul");
    const start = performance.now();
    trigger.click();
    await new Promise((resolve) => {
      const poll = () => {
        if (list.offsetParent !== null) resolve();
        else requestAnimationFrame(poll);
      };
      poll();
    });
    const elapsed = performance.now() - start;
    trigger.click();
    await new Promise((resolve) => setTimeout(resolve, 400));
    return elapsed;
  });
}

test("@perf interaction latencies are measured and reported", async ({
  page,
}) => {
  await page.addInitScript(() => {
    window.__lcpMs = 0;
    new PerformanceObserver((records) => {
      const entries = records.getEntries();
      window.__lcpMs = entries[entries.length - 1].startTime;
    }).observe({ type: "largest-contentful-paint", buffered: true });
  });

  await page.goto("/index.html");
  await page.waitForFunction(() => document.fonts.ready);
  await page.waitForFunction(() => {
    const fieldset = document.querySelector("#text-area-demo fieldset");
    return window.Alpine && fieldset.classList.contains("mt-0");
  });

  const lcpMs = Math.round(await page.evaluate(() => window.__lcpMs));
  const expandMs = await measureSamples(() => measureExpandLatency(page));
  const dropdownMs = await measureSamples(() =>
    measureDropdownOpenLatency(page),
  );

  expect(lcpMs).toBeGreaterThan(0);
  expect(expandMs).toBeGreaterThan(0);
  expect(dropdownMs).toBeGreaterThan(0);

  writeFileSync(
    resultsPath,
    `${JSON.stringify(
      {
        "demo.navigation-lcp": lcpMs,
        "form.textarea-expand": expandMs,
        "form.select-dropdown-open": dropdownMs,
      },
      null,
      2,
    )}\n`,
  );
});

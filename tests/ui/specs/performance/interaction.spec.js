import { writeFileSync } from "node:fs";
import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

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

async function measureSearchableFilterLatency(page) {
  return page.evaluate(async () => {
    const combobox = document.querySelector(
      "#searchable-select-input-demo input[role=combobox]",
    );
    combobox.focus();
    combobox.click();
    const start = performance.now();
    combobox.value = "can";
    combobox.dispatchEvent(new Event("input", { bubbles: true }));
    await new Promise((resolve) => {
      const poll = () => {
        const optionRows = [
          ...document.querySelectorAll(
            "#searchable-select-input-demo ul[role=listbox] li",
          ),
        ];
        const canadaRow = optionRows.find(
          (optionRow) => optionRow.textContent.trim() === "Canada",
        );
        if (canadaRow && canadaRow.offsetParent !== null) resolve();
        else requestAnimationFrame(poll);
      };
      poll();
    });
    const elapsed = performance.now() - start;
    combobox.value = "Brazil";
    combobox.dispatchEvent(new Event("input", { bubbles: true }));
    document.body.click();
    return elapsed;
  });
}

async function measureButtonActionLatency(page) {
  return page.evaluate(async () => {
    const button = document.querySelector("#button-demo-sync");
    const start = performance.now();
    button.click();
    await new Promise((resolve) => {
      const poll = () => {
        if (button.getAttribute("aria-busy") === "true") resolve();
        else requestAnimationFrame(poll);
      };
      poll();
    });
    return performance.now() - start;
  });
}

async function measureSliderDragLatency(page) {
  return page.evaluate(async () => {
    const rangeInput = document.querySelector(
      "#number-range-input-field-demo input[type=range]",
    );
    const numberInput = document.querySelector(
      "#number-range-input-field-demo input[type=number]",
    );
    const targetValue = rangeInput.value === "250" ? "150" : "250";
    const start = performance.now();
    rangeInput.value = targetValue;
    rangeInput.dispatchEvent(new Event("input", { bubbles: true }));
    await new Promise((resolve) => {
      const poll = () => {
        if (numberInput.value === targetValue) resolve();
        else requestAnimationFrame(poll);
      };
      poll();
    });
    return performance.now() - start;
  });
}

test("@browser @perf interaction latencies are measured and reported", async ({
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

  await openExamplePanel(page, "#button-demo", "Action States");

  const lcpMs = Math.round(await page.evaluate(() => window.__lcpMs));
  const expandMs = await measureSamples(() => measureExpandLatency(page));
  const dropdownMs = await measureSamples(() =>
    measureDropdownOpenLatency(page),
  );
  const searchableFilterMs = await measureSamples(() =>
    measureSearchableFilterLatency(page),
  );
  const sliderDragMs = await measureSamples(() =>
    measureSliderDragLatency(page),
  );
  const buttonActionMs = await measureSamples(async () => {
    const elapsed = await measureButtonActionLatency(page);
    await page.waitForFunction(() => {
      const button = document.querySelector("#button-demo-sync");
      return (
        button !== null &&
        !button.disabled &&
        button.textContent.trim() === "Sync"
      );
    });
    return elapsed;
  });

  expect(lcpMs).toBeGreaterThan(0);
  expect(expandMs).toBeGreaterThan(0);
  expect(dropdownMs).toBeGreaterThan(0);
  expect(searchableFilterMs).toBeGreaterThan(0);
  expect(sliderDragMs).toBeGreaterThan(0);
  expect(buttonActionMs).toBeGreaterThan(0);

  writeFileSync(
    resultsPath,
    `${JSON.stringify(
      {
        "demo.navigation-lcp": lcpMs,
        "form.textarea-expand": expandMs,
        "form.select-dropdown-open": dropdownMs,
        "form.searchable-select-filter": searchableFilterMs,
        "form.number-range-slider-drag": sliderDragMs,
        "control.button-action-state": buttonActionMs,
      },
      null,
      2,
    )}\n`,
  );
});

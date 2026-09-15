import { readFileSync, writeFileSync } from "node:fs";
import { test, expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";

const baselinePath = new URL("../../a11y-baseline.json", import.meta.url);

function countViolations(results) {
  const counts = {};
  for (const violation of results.violations) {
    counts[violation.id] = (counts[violation.id] ?? 0) + violation.nodes.length;
  }
  return counts;
}

function findRegressions(current, baseline) {
  const regressions = [];
  for (const [rule, count] of Object.entries(current)) {
    if (!(rule in baseline)) {
      regressions.push(`${rule}: new violation (${count} nodes)`);
      continue;
    }
    if (count > baseline[rule]) {
      regressions.push(`${rule}: ${baseline[rule]} -> ${count} nodes`);
    }
  }
  return regressions;
}

test("@a11y axe violations do not grow past the recorded baseline", async ({ page }) => {
  await page.goto("/index.html");
  await page.waitForFunction(() => document.fonts.ready);
  await page.waitForFunction(() =>
    [...document.querySelectorAll("div")].every(
      (d) => !(d.style.position === "fixed" && d.style.display === "flex"),
    ),
  );

  const current = countViolations(await new AxeBuilder({ page }).analyze());

  if (process.env.UPDATE_A11Y_BASELINE === "1") {
    writeFileSync(baselinePath, JSON.stringify(current, null, 2) + "\n");
    return;
  }

  const baseline = JSON.parse(readFileSync(baselinePath, "utf8"));
  expect(findRegressions(current, baseline)).toEqual([]);
});

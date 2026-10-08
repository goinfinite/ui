import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const filterBarSection = "#filter-bar-demo-main";
const stateReadout = "#filter-bar-demo-state span";

const chipOf = (page, section, filterLabel) =>
  page.locator(
    `${section} span[x-show]:has(button[aria-label="Remove ${filterLabel} filter"])`,
  );

test.describe("FilterBar @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#filter-bar-demo");
    await expect(page.locator(filterBarSection)).toBeVisible();
  });

  test("@smoke active value renders as a removable chip", async ({ page }) => {
    const chip = chipOf(page, filterBarSection, "Status");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("Status");
    await expect(chip).toContainText("running");
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":"running"',
    );

    await page
      .locator(`${filterBarSection} button[aria-label="Remove Status filter"]`)
      .click();
    await expect(chip).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });

  test("@smoke text filter updates the value and the chip", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} input[name=name]`).fill("alpha");

    const chip = chipOf(page, filterBarSection, "Name");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("alpha");
    await expect(page.locator(stateReadout).first()).toContainText(
      '"name":"alpha"',
    );
  });

  test("@smoke number range fills the chip with the bounds", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} input[name=cpuMin]`).fill("2");
    const chip = chipOf(page, filterBarSection, "CPU");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("≥ 2");

    await page.locator(`${filterBarSection} input[name=cpuMax]`).fill("8");
    await expect(chip).toContainText("2–8");
  });

  test("@smoke number range chip shows the maximum alone", async ({ page }) => {
    await page.locator(`${filterBarSection} input[name=cpuMax]`).fill("8");
    const chip = chipOf(page, filterBarSection, "CPU");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("≤ 8");
  });

  test("@smoke enum select updates the chip", async ({ page }) => {
    await page.locator(`${filterBarSection} [role=button]`).click();
    await page
      .locator(`${filterBarSection} ul li label`)
      .filter({ hasText: "stopped" })
      .click();

    const chip = chipOf(page, filterBarSection, "Status");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("stopped");
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":"stopped"',
    );
  });

  test("@smoke clearing the enum select empties the value and hides the chip", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} .ph-x-circle`).first().click();

    await expect(chipOf(page, filterBarSection, "Status")).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });

  test("@smoke clear filters resets every value", async ({ page }) => {
    await page.locator(`${filterBarSection} input[name=name]`).fill("alpha");
    await page
      .locator(`${filterBarSection} button`)
      .filter({ hasText: "Clear filters" })
      .click();

    await expect(chipOf(page, filterBarSection, "Name")).toBeHidden();
    await expect(chipOf(page, filterBarSection, "Status")).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText('"name":""');
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });

  test("@smoke multi-value filter holds an array and joins the chip", async ({
    page,
  }) => {
    await openExamplePanel(page, "#filter-bar-demo", "Multi-Value Filter");
    const multiSection = "#filter-bar-demo-multi";
    await page
      .locator(`${multiSection} button[aria-label="Remove Status filter"]`)
      .click();
    await page.locator(`${multiSection} [role=button]`).click();
    await page
      .locator(`${multiSection} ul li label`)
      .filter({ hasText: "running" })
      .click();
    await page
      .locator(`${multiSection} ul li label`)
      .filter({ hasText: "stopped" })
      .click();
    await page.keyboard.press("Escape");

    const chip = chipOf(page, multiSection, "Status");
    await expect(chip).toBeVisible();
    await expect(chip).toContainText("running, stopped");
    await expect(page.locator("#filter-bar-demo-multi-state")).toContainText(
      '"status":["running","stopped"]',
    );

    await page
      .locator(`${multiSection} button[aria-label="Remove Status filter"]`)
      .click();
    await expect(chip).toBeHidden();
    await expect(page.locator("#filter-bar-demo-multi-state")).toContainText(
      '"status":[]',
    );
  });
});

import { test, expect } from "@playwright/test";

const inputFieldSection = "#input-field-demo";
const selectSection = "#select-input-demo";
const multiSelectSection = "#multi-select-input-demo";
const textAreaSection = "#text-area-demo";

async function tooltipFor(trigger) {
  return trigger.page().locator(`#${await trigger.getAttribute("aria-describedby")}`);
}

test.describe("InputHint", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    await page.waitForFunction(() => window.Alpine !== undefined);
  });

  test("@smoke input field tooltip reveals when reached by keyboard", async ({ page }) => {
    const trigger = page.locator(`${inputFieldSection} [role=note]`).first();
    const input = trigger.locator("xpath=ancestor::fieldset").locator("input");
    const tooltip = await tooltipFor(trigger);

    await expect(tooltip).toBeHidden();
    await input.click();
    await page.keyboard.press("Tab");

    await expect(trigger).toBeFocused();
    await expect(tooltip).toBeVisible();
    await expect(tooltip).toContainText("This is a helpful hint displayed as a tooltip.");
  });

  test("@smoke select tooltip opens on tap without opening the dropdown", async ({ page }) => {
    const trigger = page.locator(`${selectSection} [role=note]`).first();
    const dropdown = page.locator(`${selectSection} ul`).first();
    const tooltip = await tooltipFor(trigger);

    await trigger.click();

    await expect(tooltip).toBeVisible();
    await expect(dropdown).toBeHidden();
  });

  test("@smoke multi-select tooltip reveals on focus", async ({ page }) => {
    const trigger = page.locator(`${multiSelectSection} [role=note]`).first();
    const tooltip = await tooltipFor(trigger);

    await trigger.focus();

    await expect(tooltip).toBeVisible();
    await expect(tooltip).toContainText("This is a helpful hint displayed as a tooltip.");
  });

  test("@smoke multi-select tooltip trigger does not toggle the dropdown from the keyboard", async ({ page }) => {
    const trigger = page.locator(`${multiSelectSection} [role=note]`).first();
    const dropdown = page.locator(`${multiSelectSection} ul`).first();

    await trigger.focus();
    await page.keyboard.press("Enter");
    await expect(dropdown).toBeHidden();

    await page.keyboard.press("Space");
    await expect(dropdown).toBeHidden();
  });

  test("@smoke textarea description hint renders below the field", async ({ page }) => {
    const hint = page
      .locator(textAreaSection)
      .getByText("This is a helpful hint displayed as a description below the textarea.");

    await expect(hint).toBeVisible();
  });

  test("@smoke textarea tooltip hint reveals on focus", async ({ page }) => {
    const trigger = page.locator(`${textAreaSection} [role=note]`).first();
    const tooltip = await tooltipFor(trigger);

    await trigger.focus();

    await expect(tooltip).toBeVisible();
    await expect(tooltip).toContainText("This is a helpful hint displayed as a tooltip.");
  });

  test("textarea hint text bound to a state path updates at runtime", async ({ page }) => {
    const trigger = page.locator(`${textAreaSection} [role=note]`).nth(1);
    const tooltip = await tooltipFor(trigger);
    const changeButton = page
      .locator(`${textAreaSection} button`)
      .filter({ hasText: "Change Hint Text" });

    await trigger.focus();
    await expect(tooltip).toHaveText("First hint text.");

    await changeButton.click();
    await trigger.focus();
    await expect(tooltip).toHaveText("Second hint text.");
  });
});

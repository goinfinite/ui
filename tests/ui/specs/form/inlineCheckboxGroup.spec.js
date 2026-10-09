import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const section = "#inline-checkbox-group-demo";

test.describe("InlineCheckboxGroup", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${section.slice(1)}`);
    await expect(page.locator(`${section} p.font-mono`).first()).toContainText(
      "banana",
    );
  });

  test("@smoke checking a box adds its value to the bound array", async ({
    page,
  }) => {
    const stateText = page.locator(`${section} p.font-mono`).first();
    const mainGroup = page.locator(`${section} fieldset`).first();

    await mainGroup.getByLabel("Apple", { exact: true }).check();
    await expect(stateText).toContainText("apple");

    await mainGroup.getByLabel("Banana", { exact: true }).uncheck();
    await expect(stateText).not.toContainText("banana");
  });

  test("vertical orientation stacks the checkboxes", async ({ page }) => {
    const panel = await openExamplePanel(page, section, "Vertical Orientation");
    const content = panel.locator("fieldset > div").first();
    const flexDirection = await content.evaluate(
      (element) => getComputedStyle(element).flexDirection,
    );
    expect(flexDirection).toBe("column");
  });
});

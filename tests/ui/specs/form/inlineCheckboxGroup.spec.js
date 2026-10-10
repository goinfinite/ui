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

  test("@smoke select all checks every option and fills the bound array", async ({
    page,
  }) => {
    const fieldset = page.locator(`${section} fieldset`).first();
    const stateText = page.locator(`${section} p.font-mono`).first();

    await fieldset.hover();
    await fieldset
      .getByRole("button", { name: "Select all options", exact: true })
      .click();

    await expect(stateText).toContainText('["apple","banana","orange"]');
    await expect(fieldset.getByLabel("Apple", { exact: true })).toBeChecked();
    await expect(fieldset.getByLabel("Banana", { exact: true })).toBeChecked();
    await expect(fieldset.getByLabel("Orange", { exact: true })).toBeChecked();
  });

  test("@smoke deselect all unchecks every option and empties the bound array", async ({
    page,
  }) => {
    const fieldset = page.locator(`${section} fieldset`).first();
    const stateText = page.locator(`${section} p.font-mono`).first();

    await fieldset.hover();
    await fieldset
      .getByRole("button", { name: "Deselect all options", exact: true })
      .click();

    await expect(stateText).toContainText("[]");
    await expect(fieldset.getByLabel("Apple", { exact: true })).not.toBeChecked();
    await expect(
      fieldset.getByLabel("Banana", { exact: true }),
    ).not.toBeChecked();
    await expect(
      fieldset.getByLabel("Orange", { exact: true }),
    ).not.toBeChecked();
  });

  test("actions stay hidden until hover or focus", async ({ page }) => {
    const fieldset = page.locator(`${section} fieldset`).first();
    const selectAllButton = fieldset.getByRole("button", {
      name: "Select all options",
      exact: true,
    });

    await expect(selectAllButton).toBeHidden();
    await fieldset.hover();
    await expect(selectAllButton).toBeVisible();

    await page.mouse.move(0, 0);
    await expect(selectAllButton).toBeHidden();
    await fieldset.getByLabel("Apple", { exact: true }).focus();
    await expect(selectAllButton).toBeVisible();
  });

  test("clicking the empty area of a vertical option row toggles it", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Vertical Orientation");
    const fieldset = panel.locator("fieldset").first();
    const appleBox = fieldset.getByLabel("Apple", { exact: true });
    const row = fieldset.locator(":scope > div > div").first();

    await expect(appleBox).not.toBeChecked();
    const rowBox = await row.boundingBox();
    await page.mouse.click(
      rowBox.x + rowBox.width / 2,
      rowBox.y + rowBox.height - 10,
    );
    await expect(appleBox).toBeChecked();
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

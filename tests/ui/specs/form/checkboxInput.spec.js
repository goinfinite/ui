import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const checkboxSection = "#checkbox-input-demo";

function checkboxInput(page, label) {
  return page
    .locator(`${checkboxSection} label`)
    .filter({ hasText: label })
    .first()
    .locator("input[type=checkbox]");
}

test.describe("CheckboxInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${checkboxSection.slice(1)}`);
  });

  test("@smoke clicking the checkbox updates the bound boolean", async ({
    page,
  }) => {
    const checkbox = checkboxInput(page, "Accept the terms");
    const boundValue = page
      .locator(`${checkboxSection} p`)
      .filter({ hasText: "Value:" })
      .locator("span");
    await expect(checkbox).not.toBeChecked();
    await expect(boundValue).toHaveText("false");

    await checkbox.click();
    await expect(checkbox).toBeChecked();
    await expect(boundValue).toHaveText("true");
  });

  test("@smoke select all tracks partial and full selection", async ({
    page,
  }) => {
    const selectAll = checkboxInput(page, "Select all");
    await expect(selectAll).toHaveJSProperty("indeterminate", true);

    await checkboxInput(page, "banana").click();
    await expect(selectAll).toHaveJSProperty("indeterminate", true);

    await checkboxInput(page, "orange").click();
    await expect(selectAll).toHaveJSProperty("indeterminate", false);
    await expect(selectAll).toBeChecked();

    await selectAll.click();
    await expect(checkboxInput(page, "apple")).not.toBeChecked();
  });

  test("@smoke static checked and disabled states render", async ({ page }) => {
    await openExamplePanel(page, checkboxSection, "States");
    await expect(checkboxInput(page, "Checked")).toBeChecked();
    await expect(checkboxInput(page, "Disabled")).toBeDisabled();
    await expect(checkboxInput(page, "Checked and disabled")).toBeChecked();
    await expect(checkboxInput(page, "Checked and disabled")).toBeDisabled();
  });
});

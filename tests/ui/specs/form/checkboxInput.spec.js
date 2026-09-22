import { expect, test } from "@playwright/test";

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
    await expect(checkbox).not.toBeChecked();

    await checkbox.click();
    await expect(checkbox).toBeChecked();
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
    await expect(checkboxInput(page, "Checked")).toBeChecked();
    await expect(checkboxInput(page, "Disabled")).toBeDisabled();
    await expect(checkboxInput(page, "Checked and disabled")).toBeChecked();
    await expect(checkboxInput(page, "Checked and disabled")).toBeDisabled();
  });

  test("label position left keeps the box next to the label", async ({
    page,
  }) => {
    const label = page
      .locator(`${checkboxSection} label`)
      .filter({ hasText: "Label on the left" })
      .first();
    const box = label.locator("span").first();

    const labelBox = await label.boundingBox();
    const boxBox = await box.boundingBox();
    const parentBox = await label.locator("xpath=..").boundingBox();

    expect(labelBox.width).toBeLessThan(parentBox.width / 2);
    expect(
      boxBox.x + boxBox.width - (labelBox.x + labelBox.width),
    ).toBeLessThanOrEqual(2);
  });
});

import { expect, test } from "../../fixture.js";

const numberRangeSection = "#number-range-input-field-demo";
const cpuQuotaBox = `${numberRangeSection} input[name=cpuQuota]`;
const cpuQuotaState = `${numberRangeSection} p:has-text("cpuQuota:")`;

async function typeAndDispatchChange(numberBox, typedValue) {
  await numberBox.fill(typedValue);
  await numberBox.dispatchEvent("change");
}

test.describe("NumberRangeInputField", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html${numberRangeSection}`);
    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 100");
  });

  test("@smoke a typed value commits clamped and snapped to the step", async ({
    page,
  }) => {
    await typeAndDispatchChange(page.locator(cpuQuotaBox), "137");

    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 140");
    await expect(page.locator(cpuQuotaBox)).toHaveValue("140");
  });

  test("@browser leaving the box commits the typed value", async ({ page }) => {
    await page.locator(cpuQuotaBox).fill("250");
    await page.locator(cpuQuotaBox).blur();

    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 250");
  });

  test("a partial value does not reach the state before the commit", async ({
    page,
  }) => {
    await page.locator(cpuQuotaBox).fill("5");

    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 100");
  });

  test("a value above the maximum commits as the maximum", async ({ page }) => {
    await typeAndDispatchChange(page.locator(cpuQuotaBox), "9999");

    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 400");
  });

  test("an empty box reverts to the bound value", async ({ page }) => {
    await typeAndDispatchChange(page.locator(cpuQuotaBox), "");

    await expect(page.locator(cpuQuotaBox)).toHaveValue("100");
    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 100");
  });

  test("@browser the slider writes the state and the box follows", async ({
    page,
  }) => {
    const quotaSlider = page
      .locator(`${numberRangeSection} input[type=range]`)
      .first();
    await quotaSlider.focus();
    await quotaSlider.press("ArrowRight");

    await expect(page.locator(cpuQuotaState)).toHaveText("cpuQuota: 110");
    await expect(page.locator(cpuQuotaBox)).toHaveValue("110");
  });
});

import { test, expect } from "@playwright/test";

const selectSection = "#select-input-demo";

test.describe("SelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${selectSection.slice(1)}`);
    const hiddenInput = page.locator(`${selectSection} input[type=hidden]`).first();
    await expect(hiddenInput).toHaveValue("Brazil");
  });

  test("@smoke dropdown opens, selects an option and closes", async ({ page }) => {
    const trigger = page.locator(`${selectSection} .group.flex`).first();
    await trigger.click();

    const options = page.locator(`${selectSection} ul`).first();
    await expect(options).toBeVisible();

    await options.getByText("Argentina", { exact: true }).click();
    await expect(options).toBeHidden();

    const hiddenInput = page.locator(`${selectSection} input[type=hidden]`).first();
    await expect(hiddenInput).toHaveValue("Argentina");
  });

  test("@smoke clear button empties the selection", async ({ page }) => {
    await page.locator(`${selectSection} .ph-x-circle`).first().click();
    const hiddenInput = page.locator(`${selectSection} input[type=hidden]`).first();
    await expect(hiddenInput).toHaveValue("");
  });
});

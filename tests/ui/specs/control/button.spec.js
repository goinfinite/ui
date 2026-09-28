import { expect, test } from "@playwright/test";

const buttonSection = "#button-demo";

test.describe("Button", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${buttonSection.slice(1)}`);
  });

  test("@control every button carries a type attribute", async ({ page }) => {
    const buttons = page.locator(`${buttonSection} button`);
    const buttonCount = await buttons.count();

    expect(buttonCount).toBeGreaterThan(0);
    for (let index = 0; index < buttonCount; index += 1) {
      await expect(buttons.nth(index)).toHaveAttribute(
        "type",
        /^(button|submit)$/,
      );
    }
  });

  test("@control a plain button defaults to type=button", async ({ page }) => {
    const newRecordButton = page
      .locator(`${buttonSection} button`)
      .filter({ hasText: "New Record" })
      .first();

    await expect(newRecordButton).toHaveAttribute("type", "button");
  });

  test("@control IsSubmit renders a submit button", async ({ page }) => {
    await expect(
      page.locator(`${buttonSection} button[type=submit]`),
    ).toHaveCount(1);
  });
});

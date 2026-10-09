import { expect, test } from "@playwright/test";

const toastSection = "#toast-demo";

const toastBackgroundOf = (page) =>
  page
    .locator("#toast")
    .evaluate((element) => getComputedStyle(element).backgroundColor);

test.describe("Toast", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${toastSection.slice(1)}`);
    await expect(page.locator("#toast")).toBeAttached();
  });

  test("@smoke success tone tints the surface green", async ({ page }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Saved!", "success"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgba(34, 197, 94, 0.15)");
  });

  test("@smoke partial success tone tints the surface yellow", async ({
    page,
  }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Two of three.", "partialSuccess"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgba(234, 179, 8, 0.15)");
  });

  test("@smoke danger tone tints the surface red", async ({ page }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Failed.", "danger"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgba(239, 68, 68, 0.15)");
  });
});
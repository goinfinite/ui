import { expect, test } from "@playwright/test";

const toastSection = "#toast-demo";

const toastBackgroundOf = (page) =>
  page
    .locator("#toast")
    .evaluate((element) => getComputedStyle(element).backgroundColor);

const toastShadowOf = (page) =>
  page
    .locator("#toast")
    .evaluate((element) => getComputedStyle(element).boxShadow);

test.describe("Toast", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${toastSection.slice(1)}`);
    await expect(page.locator("#toast")).toBeAttached();
  });

  test("@smoke success tone keeps the neutral surface and rings green", async ({
    page,
  }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Saved!", "success"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgb(38, 38, 38)");
    expect(await toastShadowOf(page)).toContain("rgb(34, 197, 94)");
  });

  test("@smoke partial success tone keeps the neutral surface and rings yellow", async ({
    page,
  }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Two of three.", "partialSuccess"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgb(38, 38, 38)");
    expect(await toastShadowOf(page)).toContain("rgb(234, 179, 8)");
  });

  test("@smoke danger tone keeps the neutral surface and rings red", async ({
    page,
  }) => {
    await page.evaluate(() =>
      Alpine.store("toast").displayToast("Failed.", "danger"),
    );
    await expect(page.locator("#toast")).toBeVisible();
    expect(await toastBackgroundOf(page)).toBe("rgb(38, 38, 38)");
    expect(await toastShadowOf(page)).toContain("rgb(239, 68, 68)");
  });
});

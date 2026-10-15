import { expect, test } from "../../fixture.js";

const toastSection = "#toast-demo";

const toastShadowOf = (page) =>
  page
    .locator("#toast")
    .evaluate((element) => getComputedStyle(element).boxShadow);

const toneCases = [
  {
    name: "success tone rings green",
    message: "Saved!",
    type: "success",
    ringColor: "rgb(34, 197, 94)",
  },
  {
    name: "partial success tone rings yellow",
    message: "Two of three.",
    type: "partialSuccess",
    ringColor: "rgb(234, 179, 8)",
  },
  {
    name: "danger tone rings red",
    message: "Failed.",
    type: "danger",
    ringColor: "rgb(239, 68, 68)",
  },
];

test.describe("Toast", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${toastSection.slice(1)}`);
    await expect(page.locator("#toast")).toBeAttached();
  });

  for (const toneCase of toneCases) {
    test(`@browser @smoke ${toneCase.name}`, async ({ page }) => {
      await page.evaluate(
        ({ message, type }) =>
          Alpine.store("toast").displayToast(message, type),
        { message: toneCase.message, type: toneCase.type },
      );
      await expect(page.locator("#toast")).toBeVisible();
      expect(await toastShadowOf(page)).toContain(toneCase.ringColor);
    });
  }
});

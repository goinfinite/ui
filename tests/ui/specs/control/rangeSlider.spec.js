import { test, expect } from "@playwright/test";

const sliderSection = "#range-slider-demo";

function sliderByLabel(page, label) {
  return page.locator(`${sliderSection} input[type=range][aria-label="${label}"]`);
}

function trackOf(slider) {
  return slider.locator("xpath=..");
}

function thumbOf(slider) {
  return trackOf(slider).locator("div.cursor-pointer").first();
}

function tickMarksOf(slider) {
  return trackOf(slider).locator('div[aria-hidden="true"] > span');
}

test.describe("RangeSlider", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${sliderSection.slice(1)}`);
  });

  test("@control arrow keys change the bound value", async ({ page }) => {
    const slider = sliderByLabel(page, "Value").first();
    await slider.focus();
    await page.keyboard.press("ArrowRight");
    await expect(page.locator(`${sliderSection} p`, { hasText: "Value:" }).locator("span")).toHaveText(
      "51",
    );
  });

  test("@control single slider exposes the default accessible name", async ({ page }) => {
    await expect(page.getByRole("slider", { name: "Value", exact: true }).first()).toBeVisible();
  });

  test("@control dual sliders expose distinct accessible names", async ({ page }) => {
    await expect(page.getByRole("slider", { name: "Minimum price" })).toHaveCount(1);
    await expect(page.getByRole("slider", { name: "Maximum price" })).toHaveCount(1);
    await expect(page.getByRole("slider", { name: "Lower value" })).toHaveCount(1);
    await expect(page.getByRole("slider", { name: "Upper value" })).toHaveCount(1);
  });

  test("@control focusing the input reveals the focus ring", async ({ page }) => {
    const slider = sliderByLabel(page, "Value").first();
    const thumb = thumbOf(slider);
    const restingShadow = await thumb.evaluate((element) => getComputedStyle(element).boxShadow);
    await slider.focus();
    await expect
      .poll(() => thumb.evaluate((element) => getComputedStyle(element).boxShadow))
      .not.toBe(restingShadow);
  });

  test("@control tick marks render at the configured tick step", async ({ page }) => {
    const slider = sliderByLabel(page, "Ticks every 25");
    const ticks = tickMarksOf(slider);
    await expect(ticks).toHaveCount(5);
    const trackBox = await trackOf(slider).boundingBox();
    const firstTickBox = await ticks.first().boundingBox();
    const lastTickBox = await ticks.last().boundingBox();
    const firstTickCenter = firstTickBox.x + firstTickBox.width / 2;
    const lastTickCenter = lastTickBox.x + lastTickBox.width / 2;
    expect(Math.abs(firstTickCenter - trackBox.x)).toBeLessThanOrEqual(2);
    expect(Math.abs(lastTickCenter - (trackBox.x + trackBox.width))).toBeLessThanOrEqual(2);
  });

  test("@control tick marks fall back to the slider step", async ({ page }) => {
    await expect(tickMarksOf(sliderByLabel(page, "Ticks every step"))).toHaveCount(6);
  });
});

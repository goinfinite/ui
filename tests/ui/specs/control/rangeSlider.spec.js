import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const sliderSection = "#range-slider-demo";
const trackRightInsetPx = 2;

function sliderByLabel(page, label) {
  return page.locator(
    `${sliderSection} input[type=range][aria-label="${label}"]`,
  );
}

function trackOf(slider) {
  return slider.locator("xpath=..");
}

function thumbOf(slider, index = 0) {
  return trackOf(slider).locator("div.cursor-pointer").nth(index);
}

async function dragThumb(page, slider, index, deltaX) {
  await page
    .locator('body > div[style*="z-index: 9999"]')
    .waitFor({ state: "hidden" });
  await slider.scrollIntoViewIfNeeded();
  const thumbBox = await thumbOf(slider, index).boundingBox();
  const startX = thumbBox.x + thumbBox.width / 2;
  const startY = thumbBox.y + thumbBox.height / 2;
  await page.mouse.move(startX, startY);
  await page.mouse.down();
  await page.mouse.move(startX + deltaX, startY, { steps: 5 });
  await page.mouse.up();
}

function tickMarksOf(slider) {
  return trackOf(slider).locator('div[aria-hidden="true"] > span');
}

async function writeExternalState(page, values) {
  await page.evaluate((state) => {
    const dualModeRoot = document.querySelector(
      '#range-slider-demo div[x-data^="{lowerValue"]',
    );
    const alpineData = window.Alpine.$data(dualModeRoot);
    for (const [property, value] of Object.entries(state)) {
      alpineData[property] = value;
    }
  }, values);
}

async function openTrackTooltipSlider(page, label) {
  await openExamplePanel(page, sliderSection, "Track Tooltip");
  const slider = sliderByLabel(page, label);
  const track = trackOf(slider);
  const tooltip = track.locator("[data-ui-range-slider-track-tooltip]");
  await slider.scrollIntoViewIfNeeded();
  const trackBox = await track.boundingBox();
  const hover = (horizontalFraction) => {
    const hoverX = Math.min(
      trackBox.width * horizontalFraction,
      trackBox.width - trackRightInsetPx,
    );
    return page.mouse.move(
      trackBox.x + hoverX,
      trackBox.y + trackBox.height / 2,
    );
  };
  return { tooltip, hover, trackBox };
}

async function expectTrackTooltipText(
  page,
  label,
  horizontalFraction,
  expectedText,
) {
  const { tooltip, hover } = await openTrackTooltipSlider(page, label);
  await hover(horizontalFraction);
  await expect(tooltip).toBeVisible();
  await expect(tooltip).toHaveText(expectedText);
}

test.describe("RangeSlider", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${sliderSection.slice(1)}`);
  });

  test("@control arrow keys change the bound value", async ({ page }) => {
    const slider = sliderByLabel(page, "Value").first();
    await slider.focus();
    await page.keyboard.press("ArrowRight");
    await expect(
      page.locator(`${sliderSection} p`, { hasText: "Value:" }).locator("span"),
    ).toHaveText("51");
  });

  test("@control single slider exposes the default accessible name", async ({
    page,
  }) => {
    await expect(
      page.getByRole("slider", { name: "Value", exact: true }).first(),
    ).toBeVisible();
  });

  test("@control dual sliders expose distinct accessible names", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Dual-Thumb Mode");
    await expect(
      page.getByRole("slider", { name: "Minimum price" }),
    ).toHaveCount(1);
    await expect(
      page.getByRole("slider", { name: "Maximum price" }),
    ).toHaveCount(1);
    await expect(page.getByRole("slider", { name: "Lower value" })).toHaveCount(
      1,
    );
    await expect(page.getByRole("slider", { name: "Upper value" })).toHaveCount(
      1,
    );
  });

  test("@control dragging the thumb changes the bound value", async ({
    page,
  }) => {
    const slider = sliderByLabel(page, "Value").first();
    await expect(slider).toHaveValue("50");
    await dragThumb(page, slider, 0, 60);
    expect(Number(await slider.inputValue())).toBeGreaterThan(50);
  });

  test("@control dragging each dual thumb changes its own bound value", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Dual-Thumb Mode");
    const lowerSlider = sliderByLabel(page, "Minimum price");
    const upperSlider = sliderByLabel(page, "Maximum price");
    await expect(lowerSlider).toHaveValue("25");
    await expect(upperSlider).toHaveValue("75");
    await dragThumb(page, lowerSlider, 0, 40);
    expect(Number(await lowerSlider.inputValue())).toBeGreaterThan(25);
    await expect(upperSlider).toHaveValue("75");
    await dragThumb(page, upperSlider, 1, -40);
    expect(Number(await upperSlider.inputValue())).toBeLessThan(75);
  });

  test("@control lower thumb crossing the upper clamps and syncs both bound values", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Dual-Thumb Mode");
    const lowerSlider = sliderByLabel(page, "Minimum price");
    const upperSlider = sliderByLabel(page, "Maximum price");

    await lowerSlider.focus();
    await page.keyboard.press("End");

    await expect(lowerSlider).toHaveValue("100");
    await expect(upperSlider).toHaveValue("100");
    await expect(
      page.locator(`${sliderSection} p`, { hasText: "Price Range" }),
    ).toContainText("$100 - $100");
  });

  test("@control external state writes keep both thumbs ordered", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Dual-Thumb Mode");
    const lowerSlider = sliderByLabel(page, "Minimum price");
    const upperSlider = sliderByLabel(page, "Maximum price");
    const rangeDisplay = page.locator(`${sliderSection} p`, {
      hasText: "Price Range",
    });

    await page.waitForFunction(() => window.Alpine !== undefined);

    await writeExternalState(page, { lowerValue: 90 });
    await expect(lowerSlider).toHaveValue("90");
    await expect(upperSlider).toHaveValue("91");
    await expect(rangeDisplay).toContainText("$90 - $91");

    await writeExternalState(page, { upperValue: 10 });
    await expect(upperSlider).toHaveValue("10");
    await expect(lowerSlider).toHaveValue("9");
    await expect(rangeDisplay).toContainText("$9 - $10");

    await writeExternalState(page, { lowerValue: 150 });
    await expect(lowerSlider).toHaveValue("100");
    await expect(upperSlider).toHaveValue("100");
    await expect(rangeDisplay).toContainText("$100 - $100");
  });

  test("@control initial crossed and out-of-bounds state normalizes when the setting is on", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Initial State Normalization");
    const initialSection = page.locator("#range-slider-initial-state-demo");
    const sliders = initialSection.locator("input[type=range]");

    await expect(sliders.nth(0)).toHaveValue("80");
    await expect(sliders.nth(1)).toHaveValue("81");
    await expect(sliders.nth(2)).toHaveValue("100");
    await expect(sliders.nth(3)).toHaveValue("100");
    await expect(sliders.nth(4)).toHaveValue("100");
    await expect(initialSection.locator("p")).toContainText("$80 - $81");
    await expect(initialSection.locator("p")).toContainText("$100 - $100");
    await expect(initialSection.locator("p")).toContainText("Slider: 100");
  });

  test("@control initial crossed state stays untouched when the setting is off", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Initial State Normalization");
    const unnormalizedSection = page.locator("#range-slider-unnormalized-demo");
    const sliders = unnormalizedSection.locator("input[type=range]");

    await expect(sliders.nth(0)).toHaveValue("80");
    await expect(sliders.nth(1)).toHaveValue("20");
    await expect(unnormalizedSection.locator("p")).toContainText("$80 - $20");
  });

  test("@control focusing the input reveals the focus ring", async ({
    page,
  }) => {
    const slider = sliderByLabel(page, "Value").first();
    const thumb = thumbOf(slider);
    const restingShadow = await thumb.evaluate(
      (element) => getComputedStyle(element).boxShadow,
    );
    await slider.focus();
    await expect
      .poll(() =>
        thumb.evaluate((element) => getComputedStyle(element).boxShadow),
      )
      .not.toBe(restingShadow);
  });

  test("@control tick marks render at the configured tick step", async ({
    page,
  }) => {
    await openExamplePanel(page, sliderSection, "Steps & Ticks");
    const slider = sliderByLabel(page, "Ticks every 25");
    const ticks = tickMarksOf(slider);
    await expect(ticks).toHaveCount(5);
    const trackBox = await trackOf(slider).boundingBox();
    const firstTickBox = await ticks.first().boundingBox();
    const lastTickBox = await ticks.last().boundingBox();
    const firstTickCenter = firstTickBox.x + firstTickBox.width / 2;
    const lastTickCenter = lastTickBox.x + lastTickBox.width / 2;
    expect(Math.abs(firstTickCenter - trackBox.x)).toBeLessThanOrEqual(2);
    expect(
      Math.abs(lastTickCenter - (trackBox.x + trackBox.width)),
    ).toBeLessThanOrEqual(2);
  });

  test("@control tick marks fall back to the slider step", async ({ page }) => {
    await openExamplePanel(page, sliderSection, "Steps & Ticks");
    await expect(
      tickMarksOf(sliderByLabel(page, "Ticks every step")),
    ).toHaveCount(6);
  });

  test("@control hovering the track previews the value under the pointer", async ({
    page,
  }) => {
    const { tooltip, hover, trackBox } = await openTrackTooltipSlider(
      page,
      "Tooltip top",
    );

    await expect(tooltip).toBeHidden();

    await hover(0.3);
    await expect(tooltip).toBeVisible();
    await expect(tooltip).toHaveText("30");

    await hover(1);
    await expect(tooltip).toHaveText("100");

    await page.mouse.move(trackBox.x, trackBox.y - 40);
    await expect(tooltip).toBeHidden();
  });

  test("@control the track tooltip snaps to a fractional step", async ({
    page,
  }) => {
    await expectTrackTooltipText(page, "Tooltip bottom", 0.3, "0.3");
  });

  test("@control the track tooltip keeps a sub-micro step preview nonzero", async ({
    page,
  }) => {
    await expectTrackTooltipText(page, "Tooltip tiny step", 0.3, "3e-7");
  });

  test("@control the track tooltip caps an unaligned endpoint to the step grid", async ({
    page,
  }) => {
    await expectTrackTooltipText(page, "Tooltip unaligned", 1, "60");
  });

  test("@control the track tooltip keeps an aligned decimal endpoint on the grid", async ({
    page,
  }) => {
    await expectTrackTooltipText(page, "Tooltip aligned decimal", 1, "0.3");
  });

  test("@control the value bubble display mode gates visibility", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, sliderSection, "Value Bubble");
    const alwaysBubble = trackOf(
      panel.locator("input[type=range]").nth(0),
    ).locator("[data-ui-range-slider-thumb-bubble]");
    await expect(alwaysBubble).toBeVisible();
    await expect(alwaysBubble).toHaveText("50");

    const hoverSlider = sliderByLabel(page, "Hover bubble");
    const hoverTrack = trackOf(hoverSlider);
    const hoverBubble = hoverTrack.locator(
      "[data-ui-range-slider-thumb-bubble]",
    );
    await expect(hoverBubble).toBeHidden();

    await hoverSlider.scrollIntoViewIfNeeded();
    const trackBox = await hoverTrack.boundingBox();
    await page.mouse.move(
      trackBox.x + trackBox.width / 2,
      trackBox.y + trackBox.height / 2,
    );
    await expect(hoverBubble).toBeVisible();
    await expect(hoverBubble).toHaveText("50");

    await page.mouse.move(trackBox.x, trackBox.y - 80);
    await expect(hoverBubble).toBeHidden();

    await hoverSlider.focus();
    await expect(hoverBubble).toBeVisible();
  });
});

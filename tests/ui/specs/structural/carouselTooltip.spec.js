import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const carouselSection = "#carousel-demo";
const inlineCarouselRoot = "#carousel-demo-tooltips-carousel";
const modalCarouselRoot = "#carousel-demo-tooltips-modal-carousel";

function countTeleportedTooltips(page) {
  return page.evaluate(
    () => document.querySelectorAll('body > div[x-ref="tooltip"]').length,
  );
}

test.describe("Carousel item tooltips @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${carouselSection.slice(1)}`);
    await openExamplePanel(page, carouselSection, "Item Tooltips");
  });

  test("hovering an item shows one body tooltip wired to the item", async ({
    page,
  }) => {
    const firstTrigger = page
      .locator(
        `${inlineCarouselRoot} [data-ui-carousel-track] [x-ref="trigger"]`,
      )
      .first();
    await firstTrigger.hover();

    const tooltip = page.locator('body > div[x-ref="tooltip"]:visible');
    await expect(tooltip).toHaveCount(1);

    const tooltipId = await tooltip.getAttribute("id");
    expect(await firstTrigger.getAttribute("aria-describedby")).toBe(tooltipId);
    await expect(tooltip).toContainText("Handles the public API traffic");

    const viewportSize = page.viewportSize();
    const tooltipBox = await tooltip.boundingBox();
    expect(tooltipBox.x).toBeGreaterThanOrEqual(0);
    expect(tooltipBox.y).toBeGreaterThanOrEqual(0);
    expect(tooltipBox.x + tooltipBox.width).toBeLessThanOrEqual(
      viewportSize.width,
    );
    expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(
      viewportSize.height,
    );

    await page.mouse.move(0, 0);
    await expect(tooltip).toHaveCount(0);
  });

  test("focusing an item shows its tooltip for keyboard users", async ({
    page,
  }) => {
    await page.mouse.move(0, 0);
    const firstTrigger = page
      .locator(
        `${inlineCarouselRoot} [data-ui-carousel-track] [x-ref="trigger"]`,
      )
      .first();

    await firstTrigger.focus();

    const tooltip = page.locator('body > div[x-ref="tooltip"]:visible');
    await expect(tooltip).toHaveCount(1);
    await expect(tooltip).toContainText("Handles the public API traffic");

    await firstTrigger.blur();
    await expect(tooltip).toHaveCount(0);
  });

  test("the tooltip survives inside a modal", async ({ page }) => {
    await page
      .getByRole("button", { name: "Open the carousel in a modal" })
      .click();

    const modalCarousel = page.locator(modalCarouselRoot);
    await expect(modalCarousel).toBeVisible();

    const firstTrigger = page
      .locator(
        `${modalCarouselRoot} [data-ui-carousel-track] [x-ref="trigger"]`,
      )
      .first();
    await firstTrigger.hover();

    const tooltip = page.locator('body > div[x-ref="tooltip"]:visible');
    await expect(tooltip).toHaveCount(1);

    const isTeleportedToBody = await tooltip.evaluate(
      (element) => element.parentElement === document.body,
    );
    expect(isTeleportedToBody).toBe(true);

    const viewportSize = page.viewportSize();
    const tooltipBox = await tooltip.boundingBox();
    expect(tooltipBox.x).toBeGreaterThanOrEqual(0);
    expect(tooltipBox.y).toBeGreaterThanOrEqual(0);
    expect(tooltipBox.x + tooltipBox.width).toBeLessThanOrEqual(
      viewportSize.width,
    );
    expect(tooltipBox.y + tooltipBox.height).toBeLessThanOrEqual(
      viewportSize.height,
    );

    const backdrop = page
      .locator("div.fixed.inset-0.z-100")
      .filter({ has: modalCarousel });
    const paintsAboveBackdrop = await page.evaluate(
      ([tooltipElement, backdropElement]) => {
        const tooltipZ = Number(getComputedStyle(tooltipElement).zIndex);
        const backdropZ = Number(getComputedStyle(backdropElement).zIndex);
        const tooltipIsLaterInBody = Boolean(
          backdropElement.compareDocumentPosition(tooltipElement) &
            Node.DOCUMENT_POSITION_FOLLOWING,
        );
        return (
          tooltipZ > backdropZ ||
          (tooltipZ === backdropZ && tooltipIsLaterInBody)
        );
      },
      [await tooltip.elementHandle(), await backdrop.elementHandle()],
    );
    expect(paintsAboveBackdrop).toBe(true);
  });

  test("a refresh does not leak teleported tooltips", async ({ page }) => {
    const initialTooltipCount = await countTeleportedTooltips(page);

    await page
      .locator(`${inlineCarouselRoot} input[name=search]`)
      .pressSequentially("bravo", { delay: 40 });

    await expect
      .poll(() => countTeleportedTooltips(page))
      .toBe(initialTooltipCount);
  });
});

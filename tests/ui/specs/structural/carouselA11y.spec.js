import { expect, test } from "@playwright/test";

const carouselRoot = "#carousel-demo-main-carousel";

test.describe("Carousel accessibility @structural @a11y", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#carousel-demo");
    await expect(page.locator(`${carouselRoot} nav p`)).toHaveText("1–6 of 6");
  });

  test("prev and next controls carry accessible names", async ({ page }) => {
    await expect(
      page.locator(`${carouselRoot} button[aria-label="Previous items"]`),
    ).toBeVisible();
    await expect(
      page.locator(`${carouselRoot} button[aria-label="Next items"]`),
    ).toBeVisible();
  });

  test("the controls work with the keyboard", async ({ page }) => {
    const nextButton = page.locator(
      `${carouselRoot} button[aria-label="Next items"]`,
    );
    await nextButton.focus();
    await page.keyboard.press("Enter");

    await expect
      .poll(() =>
        page.evaluate(
          (selector) =>
            Alpine.$data(document.querySelector(selector)).windowStart,
          carouselRoot,
        ),
      )
      .toBe(1);
  });

  test("each dot carries an accessible name and marks the current window", async ({
    page,
  }) => {
    const firstDot = page.locator(
      `${carouselRoot} [aria-label="Go to item 1"]`,
    );
    await expect(firstDot).toBeVisible();
    await expect(firstDot).toHaveAttribute("aria-current", "true");

    await page.locator(`${carouselRoot} [aria-label="Go to item 3"]`).click();
    await expect(
      page.locator(`${carouselRoot} [aria-label="Go to item 3"]`),
    ).toHaveAttribute("aria-current", "true");
  });
  test("pagination is a named navigation landmark with named controls", async ({
    page,
  }) => {
    const pagination = page.locator(
      `${carouselRoot} nav[aria-label="Carousel pagination"]`,
    );
    await expect(pagination).toBeVisible();
    for (const buttonLabel of [
      "First page",
      "Previous page",
      "Next page",
      "Last page",
    ]) {
      await expect(
        pagination.locator(`button[aria-label="${buttonLabel}"]`),
      ).toBeVisible();
    }
  });

  test("off-window items are removed from the tab order", async ({ page }) => {
    await expect
      .poll(() =>
        page.evaluate((selector) => {
          const track = document.querySelector(
            `${selector} [data-ui-carousel-track]`,
          );
          return Array.from(track.children).map((item) => item.inert);
        }, carouselRoot),
      )
      .toEqual([false, false, false, true, true, true]);
  });
});

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const carouselSection = "#carousel-demo";
const carouselRoot = "#carousel-demo-main-carousel";
const refreshFragmentPattern = /carouselDemoRefresh/;

function trackRefreshUrls(page) {
  const refreshUrls = [];
  page.on("request", (request) => {
    if (request.url().includes("carouselDemoRefresh")) {
      refreshUrls.push(request.url());
    }
  });
  return refreshUrls;
}

function carouselData(page, rootSelector) {
  return page.evaluate((selector) => {
    const root = document.querySelector(selector);
    const data = Alpine.$data(root);
    return {
      canSlideNext: data.canSlideNext,
      canSlidePrevious: data.canSlidePrevious,
      dotCount: data.dotIndexes.length,
      itemsCount: data.itemsCount,
      itemsPerView: data.itemsPerView,
      pageNumber: data.pageNumber,
      windowStart: data.windowStart,
    };
  }, rootSelector);
}

async function visibleItemCount(page, rootSelector) {
  return page.evaluate((selector) => {
    const root = document.querySelector(selector);
    const viewport = root.querySelector("[x-ref=viewport]");
    const viewportBounds = viewport.getBoundingClientRect();
    const items = root.querySelectorAll("[data-ui-carousel-track] > div");
    let visibleCount = 0;
    for (const item of items) {
      const itemBounds = item.getBoundingClientRect();
      if (
        itemBounds.left >= viewportBounds.left - 2 &&
        itemBounds.right <= viewportBounds.right + 2
      ) {
        visibleCount++;
      }
    }
    return visibleCount;
  }, rootSelector);
}

test.describe("Carousel @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${carouselSection.slice(1)}`);
    await expect(page.locator(`${carouselRoot} nav p`)).toHaveText("1–6 of 6");
  });

  test("@smoke renders three items per view on a wide window", async ({
    page,
  }) => {
    await expect(
      page.locator(`${carouselRoot} [data-ui-carousel-track] > div`),
    ).toHaveCount(6);
    expect((await carouselData(page, carouselRoot)).itemsPerView).toBe(3);
    expect(await visibleItemCount(page, carouselRoot)).toBe(3);
    await expect(page.locator(`${carouselRoot} nav p`)).toHaveText("1–6 of 6");
  });

  test("@smoke next and previous move the window one item", async ({
    page,
  }) => {
    await expect(
      page.locator(`${carouselRoot} button[aria-label="Previous items"]`),
    ).toBeDisabled();

    await page
      .locator(`${carouselRoot} button[aria-label="Next items"]`)
      .click();
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(1);

    await page
      .locator(`${carouselRoot} button[aria-label="Previous items"]`)
      .click();
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(0);
  });

  test("@smoke the last window disables next", async ({ page }) => {
    const carouselState = await carouselData(page, carouselRoot);
    const lastWindowStart =
      carouselState.itemsCount - carouselState.itemsPerView;

    await page.evaluate(
      ({ selector, start }) => {
        Alpine.$data(document.querySelector(selector)).goToWindow(start);
      },
      { selector: carouselRoot, start: lastWindowStart },
    );

    await expect(
      page.locator(`${carouselRoot} button[aria-label="Next items"]`),
    ).toBeDisabled();
    await expect(
      page.locator(`${carouselRoot} button[aria-label="Previous items"]`),
    ).toBeEnabled();
  });

  test("@smoke the dots jump to a window and mark the current one", async ({
    page,
  }) => {
    expect((await carouselData(page, carouselRoot)).dotCount).toBe(4);

    await page.locator(`${carouselRoot} [aria-label="Go to item 3"]`).click();
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(2);
    await expect(
      page.locator(`${carouselRoot} [aria-label="Go to item 3"]`),
    ).toHaveAttribute("aria-current", "true");
  });

  test("@smoke a swipe moves the window", async ({ page }) => {
    const viewport = page.locator(`${carouselRoot} [x-ref=viewport]`);
    await viewport.scrollIntoViewIfNeeded();

    const swipe = async (fromRatio, toRatio) => {
      const bounds = await viewport.boundingBox();
      const centerY = bounds.y + bounds.height / 2;
      await page.mouse.move(bounds.x + bounds.width * fromRatio, centerY);
      await page.mouse.down();
      await page.mouse.move(bounds.x + bounds.width * toRatio, centerY, {
        steps: 5,
      });
      await page.mouse.up();
      await page.waitForTimeout(400);
    };

    await swipe(0.85, 0.15);
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(1);

    await swipe(0.15, 0.85);
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(0);
  });

  test("@smoke a touch pointer swipe moves the window", async ({ page }) => {
    const viewport = page.locator(`${carouselRoot} [x-ref=viewport]`);
    await viewport.scrollIntoViewIfNeeded();

    const touchSwipe = (fromRatio, toRatio) =>
      page.evaluate(
        ({ selector, from, to }) => {
          const viewportElement = document.querySelector(
            `${selector} [x-ref=viewport]`,
          );
          const bounds = viewportElement.getBoundingClientRect();
          const centerY = bounds.y + bounds.height / 2;
          const pointerEvent = (eventType, ratio) =>
            new PointerEvent(eventType, {
              bubbles: true,
              clientX: bounds.x + bounds.width * ratio,
              clientY: centerY,
              pointerId: 1,
              pointerType: "touch",
            });
          viewportElement.dispatchEvent(pointerEvent("pointerdown", from));
          viewportElement.dispatchEvent(pointerEvent("pointerup", to));
        },
        { selector: carouselRoot, from: fromRatio, to: toRatio },
      );

    await touchSwipe(0.85, 0.15);
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(1);

    await touchSwipe(0.15, 0.85);
    expect((await carouselData(page, carouselRoot)).windowStart).toBe(0);
  });

  test("@smoke the named breakpoints change the visible item count", async ({
    page,
  }) => {
    await page.setViewportSize({ width: 500, height: 900 });
    await expect
      .poll(() =>
        carouselData(page, carouselRoot).then((data) => data.itemsPerView),
      )
      .toBe(1);
    await expect.poll(() => visibleItemCount(page, carouselRoot)).toBe(1);

    await page.setViewportSize({ width: 700, height: 900 });
    await expect
      .poll(() =>
        carouselData(page, carouselRoot).then((data) => data.itemsPerView),
      )
      .toBe(2);
    await expect.poll(() => visibleItemCount(page, carouselRoot)).toBe(2);

    await page.setViewportSize({ width: 1280, height: 900 });
    await expect
      .poll(() =>
        carouselData(page, carouselRoot).then((data) => data.itemsPerView),
      )
      .toBe(3);
    await expect.poll(() => visibleItemCount(page, carouselRoot)).toBe(3);
  });

  test("@smoke autoplay advances the window and pauses on hover", async ({
    page,
  }) => {
    await openExamplePanel(page, carouselSection, "Autoplay");
    const autoplayRoot = "#carousel-demo-autoplay-carousel";

    await expect
      .poll(
        () => carouselData(page, autoplayRoot).then((data) => data.windowStart),
        { timeout: 6000 },
      )
      .toBeGreaterThan(0);

    const pausedStart = (await carouselData(page, autoplayRoot)).windowStart;
    await page.locator(autoplayRoot).hover();
    await page.waitForTimeout(3500);
    expect((await carouselData(page, autoplayRoot)).windowStart).toBe(
      pausedStart,
    );
  });

  test("@smoke a page change refreshes the chunk from the server", async ({
    page,
  }) => {
    await openExamplePanel(page, carouselSection, "Server Pagination");
    const serverRoot = "#carousel-demo-server-carousel";
    await expect(page.locator(`${serverRoot} nav p`)).toHaveText("1–6 of 25");

    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${serverRoot} button[aria-label="Next page"]`).click();

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("page=1");
    await expect(page.locator(`${serverRoot} nav p`)).toHaveText("7–12 of 25");
    await expect(page.locator(serverRoot)).toContainText("golf");
  });

  test("refreshing one carousel leaves another carousel's pagination alone", async ({
    page,
  }) => {
    await openExamplePanel(page, carouselSection, "Server Pagination");
    const serverRoot = "#carousel-demo-server-carousel";
    const mainRoot = "#carousel-demo-main-carousel";

    await expect(
      page.locator(`${mainRoot} [data-ui-carousel-pagination]`),
    ).toContainText("1–6 of 6");

    await page.locator(`${serverRoot} button[aria-label="Next page"]`).click();

    await expect(page.locator(`${serverRoot} nav p`)).toHaveText("7–12 of 25");
    await expect(
      page.locator(`${mainRoot} [data-ui-carousel-pagination]`),
    ).toContainText("1–6 of 6");
  });

  test("typing in the search box debounces into one request", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page
      .locator(`${carouselRoot} input[name=search]`)
      .pressSequentially("bravo", { delay: 40 });

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("search=bravo");
  });

  test("a filter change refreshes with the filter param", async ({ page }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${carouselRoot} input[name=name]`).fill("alpha");

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("name=alpha");
  });

  test("a failed refresh shows the error state and retry recovers", async ({
    page,
  }) => {
    await openExamplePanel(page, carouselSection, "Server Pagination");
    const serverRoot = "#carousel-demo-server-carousel";

    let shouldFail = true;
    await page.route(refreshFragmentPattern, (route) => {
      if (shouldFail) {
        shouldFail = false;
        return route.fulfill({ status: 500, body: "server error" });
      }
      return route.continue();
    });

    await page.locator(`${serverRoot} button[aria-label="Next page"]`).click();
    await expect(
      page.locator(serverRoot).getByText("Could not refresh the carousel."),
    ).toBeVisible();

    await page
      .locator(serverRoot)
      .getByRole("button", { name: "Retry" })
      .click();
    await expect(
      page.locator(serverRoot).getByText("Could not refresh the carousel."),
    ).toBeHidden();
    await expect(
      page.locator(`${serverRoot} [data-ui-carousel-track] > div`),
    ).toHaveCount(6);
  });

  test("items per page dropdown stays inside the carousel surface", async ({
    page,
  }) => {
    const selectTrigger = page.locator(`${carouselRoot} nav [role=button]`);
    await selectTrigger.scrollIntoViewIfNeeded();
    await selectTrigger.click();

    const dropdown = page.locator(`${carouselRoot} nav ul`);
    await expect(dropdown).toBeVisible();

    const surface = page.locator(carouselRoot);
    await expect
      .poll(async () => {
        const dropdownBox = await dropdown.boundingBox();
        const surfaceBox = await surface.boundingBox();
        return {
          belowSurfaceTop: dropdownBox.y >= surfaceBox.y,
          aboveSurfaceBottom:
            dropdownBox.y + dropdownBox.height <=
            surfaceBox.y + surfaceBox.height,
        };
      })
      .toEqual({ belowSurfaceTop: true, aboveSurfaceBottom: true });
  });
});

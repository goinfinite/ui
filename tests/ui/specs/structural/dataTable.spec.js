import { expect, test } from "@playwright/test";

const tableSection = "#data-table-demo";
const tableRoot = "#data-table-demo-table";
const refreshFragmentPattern = /dataTableDemoRefresh/;

function trackRefreshUrls(page) {
  const refreshUrls = [];
  page.on("request", (request) => {
    if (request.url().includes("dataTableDemoRefresh")) {
      refreshUrls.push(request.url());
    }
  });
  return refreshUrls;
}

function rowNames(page) {
  return page.locator(`${tableRoot} tbody tr td:nth-child(2)`);
}

test.describe("DataTable @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${tableSection.slice(1)}`);
    await expect(rowNames(page).first()).toHaveText("alpha");
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("1–5 of 25");
  });

  test("@smoke renders headers, rows, and the pagination readout", async ({
    page,
  }) => {
    await expect(page.locator(`${tableRoot} thead th`)).toHaveCount(5);
    await expect(page.locator(`${tableRoot} thead`)).toContainText("Name");
    await expect(page.locator(`${tableRoot} thead`)).toContainText("CPU cores");
    await expect(rowNames(page)).toHaveCount(5);
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("1–5 of 25");
  });

  test("@smoke selecting a row shows the bulk action and announces the count", async ({
    page,
  }) => {
    await page
      .locator(`${tableRoot} tbody tr`)
      .first()
      .locator("input[type=checkbox]")
      .click();

    await expect(
      page
        .locator(`${tableRoot} button`)
        .filter({ hasText: "Clear selection" }),
    ).toBeVisible();
    await expect(page.locator(`${tableRoot} p[aria-live=polite]`)).toHaveText(
      "1 row selected",
    );
    await expect(
      page.locator(`${tableRoot} span`).filter({ hasText: "1 selected" }),
    ).toBeVisible();
  });

  test("@smoke select-all toggles the current page", async ({ page }) => {
    const headerCheckbox = page.locator(
      `${tableRoot} thead input[type=checkbox]`,
    );
    await headerCheckbox.click();

    await expect(
      page.locator(`${tableRoot} tbody input[type=checkbox]:checked`),
    ).toHaveCount(5);
    await expect(headerCheckbox).toBeChecked();

    await headerCheckbox.click();
    await expect(
      page.locator(`${tableRoot} tbody input[type=checkbox]:checked`),
    ).toHaveCount(0);
  });

  test("@smoke next page requests page two and swaps the rows", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("page=2");
    await expect(rowNames(page).first()).toHaveText("foxtrot");
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("6–10 of 25");
  });

  test("@smoke clicking a page number refreshes to that page", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await expect(
      page.locator(`${tableRoot} nav [aria-current=page]`),
    ).toHaveText("1");

    await page
      .locator(`${tableRoot} nav button`)
      .filter({ hasText: /^2$/ })
      .click();

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("page=2");
    await expect(rowNames(page).first()).toHaveText("foxtrot");
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("6–10 of 25");
    await expect(
      page.locator(`${tableRoot} nav [aria-current=page]`),
    ).toHaveText("2");
  });

  test("sorting cycles ascending, descending, and unsorted", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    const createdAtHeader = page
      .locator(`${tableRoot} thead th`)
      .filter({ hasText: "Created at" });
    const createdAtButton = createdAtHeader.locator("button");

    await createdAtButton.click();
    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("sort=createdAt");
    expect(refreshUrls[0]).toContain("direction=asc");
    expect(refreshUrls[0]).toContain("page=1");
    await expect(createdAtHeader).toHaveAttribute("aria-sort", "ascending");

    await createdAtButton.click();
    await expect.poll(() => refreshUrls.length).toBe(2);
    expect(refreshUrls[1]).toContain("direction=desc");
    await expect(createdAtHeader).toHaveAttribute("aria-sort", "descending");

    await createdAtButton.click();
    await expect.poll(() => refreshUrls.length).toBe(3);
    expect(refreshUrls[2]).not.toContain("sort=createdAt");
    await expect(createdAtHeader).toHaveAttribute("aria-sort", "none");
  });

  test("typing a filter debounces into one request", async ({ page }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page
      .locator(`${tableRoot} input[name=name]`)
      .pressSequentially("alpha", { delay: 40 });

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("name=alpha");
    await expect(
      page.locator(`${tableRoot} span`).filter({ hasText: "Name: alpha" }),
    ).toBeVisible();
  });

  test("search debounces into one request and resets the page", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();
    await expect.poll(() => refreshUrls.length).toBe(1);

    await page
      .locator(`${tableRoot} input[name=search]`)
      .pressSequentially("bravo", { delay: 40 });
    await expect.poll(() => refreshUrls.length).toBe(2);
    expect(refreshUrls[1]).toContain("search=bravo");
    expect(refreshUrls[1]).toContain("page=1");
  });

  test("selection survives a refresh", async ({ page }) => {
    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();
    await expect(rowNames(page).first()).toHaveText("foxtrot");

    await page
      .locator(`${tableRoot} tbody tr`)
      .first()
      .locator("input[type=checkbox]")
      .click();
    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Refresh" })
      .click();

    await expect(
      page.locator(`${tableRoot} tbody input[type=checkbox]:checked`),
    ).toHaveCount(1);
    await expect(page.locator(`${tableRoot} p[aria-live=polite]`)).toHaveText(
      "1 row selected",
    );
  });

  test("clearing the enum filter refreshes without the filter param", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${tableRoot} .ph-x-circle`).first().click();

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).not.toContain("status=");
    expect(refreshUrls[0]).toContain("page=1");
  });

  test("removing a filter chip refreshes without the filter param", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page
      .locator(`${tableRoot} button[aria-label="Remove Status filter"]`)
      .click();

    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).not.toContain("status=");
    expect(refreshUrls[0]).toContain("page=1");
    await expect(
      page.locator(`${tableRoot} span`).filter({ hasText: "Status: running" }),
    ).toBeHidden();
  });

  test("refresh keeps the current page's rows", async ({ page }) => {
    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Refresh" })
      .click();
    await expect(rowNames(page).first()).toHaveText("alpha");
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("1–5 of 25");

    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();
    await expect(rowNames(page).first()).toHaveText("foxtrot");

    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Refresh" })
      .click();
    await expect(rowNames(page).first()).toHaveText("foxtrot");
  });

  test("the larger page size and the last page resolve to their own fragments", async ({
    page,
  }) => {
    await page.locator(`${tableRoot} nav [role=button]`).click();
    await page
      .locator(`${tableRoot} nav ul`)
      .getByText("25", { exact: true })
      .click();
    await expect(page.locator(`${tableRoot} tbody tr`)).toHaveCount(25);
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("1–25 of 25");

    await page.locator(`${tableRoot} nav [role=button]`).click();
    await page
      .locator(`${tableRoot} nav ul`)
      .getByText("5", { exact: true })
      .click();
    await expect(rowNames(page).first()).toHaveText("alpha");

    await page.locator(`${tableRoot} button[aria-label="Last page"]`).click();
    await expect(rowNames(page).first()).toHaveText("uniform");

    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Refresh" })
      .click();
    await expect(rowNames(page).first()).toHaveText("uniform");
  });

  test("clear filters resets every filter and refreshes", async ({ page }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${tableRoot} input[name=name]`).fill("alpha");
    await expect.poll(() => refreshUrls.length).toBe(1);

    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Clear filters" })
      .click();
    await expect.poll(() => refreshUrls.length).toBe(2);
    expect(refreshUrls[1]).not.toContain("status=");
    expect(refreshUrls[1]).not.toContain("name=");
    expect(refreshUrls[1]).toContain("page=1");
    await expect(rowNames(page).first()).toHaveText("alpha");
  });

  test("a number range filter refreshes with min and max params", async ({
    page,
  }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.locator(`${tableRoot} input[name=cpuMin]`).fill("2");
    await expect.poll(() => refreshUrls.length).toBe(1);
    expect(refreshUrls[0]).toContain("cpuMin=2");

    await page.locator(`${tableRoot} input[name=cpuMax]`).fill("8");
    await expect.poll(() => refreshUrls.length).toBe(2);
    expect(refreshUrls[1]).toContain("cpuMin=2");
    expect(refreshUrls[1]).toContain("cpuMax=8");
  });

  test("controls disable while a refresh is in flight", async ({ page }) => {
    await page.route(refreshFragmentPattern, async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 800));
      return route.continue();
    });

    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();
    await expect(
      page.locator(`${tableRoot} nav [role=button]`),
    ).toHaveAttribute("aria-disabled", "true");
    await expect(
      page.locator(`${tableRoot} button[aria-label="Next page"]`),
    ).toBeDisabled();

    await expect(rowNames(page).first()).toHaveText("foxtrot");
    await expect(
      page.locator(`${tableRoot} nav [role=button]`),
    ).not.toHaveAttribute("aria-disabled", "true");
  });

  test("a failed refresh shows the error row and retry recovers", async ({
    page,
  }) => {
    let shouldFail = true;
    await page.route(refreshFragmentPattern, (route) => {
      if (shouldFail) {
        shouldFail = false;
        return route.fulfill({ status: 500, body: "server error" });
      }
      return route.continue();
    });

    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Refresh" })
      .click();
    await expect(
      page.locator(tableRoot).getByText("Could not refresh the table."),
    ).toBeVisible();
    await expect(page.locator(`${tableRoot} table`)).toBeHidden();

    await page
      .locator(`${tableRoot} button`)
      .filter({ hasText: "Retry" })
      .click();
    await expect(page.locator(`${tableRoot} table`)).toBeVisible();
    await expect(
      page.locator(tableRoot).getByText("Could not refresh the table."),
    ).toBeHidden();
  });

  test("a stale failed refresh does not show the error row", async ({
    page,
  }) => {
    let requestCount = 0;
    await page.route(refreshFragmentPattern, async (route) => {
      requestCount++;
      if (requestCount === 1) {
        return route.abort("failed");
      }
      await new Promise((resolve) => setTimeout(resolve, 600));
      return route.continue();
    });

    await page.evaluate(() => {
      const dataTable = Alpine.$data(
        document.getElementById("data-table-demo-table"),
      );
      dataTable.refresh();
      dataTable.refresh();
    });

    await expect.poll(() => requestCount).toBe(2);
    await expect(rowNames(page).first()).toHaveText("alpha");
    await expect(page.locator(`${tableRoot} p[aria-live=polite]`)).toHaveText(
      "Table refreshed",
    );
    await expect(
      page.locator(tableRoot).getByText("Could not refresh the table."),
    ).toBeHidden();
  });

  test("refreshes through the fetch fallback when htmx is absent", async ({
    browser,
  }) => {
    const context = await browser.newContext({ baseURL: process.env.DEMO_URL });
    await context.route("**/htmx.min.js*", (route) => route.abort());
    const page = await context.newPage();

    await page.goto(`/index.html#${tableSection.slice(1)}`);
    await expect(
      page.locator(`${tableRoot} tbody tr td:nth-child(2)`).first(),
    ).toHaveText("alpha");
    expect(await page.evaluate(() => typeof window.htmx)).toBe("undefined");

    await page.locator(`${tableRoot} button[aria-label="Next page"]`).click();
    await expect(
      page.locator(`${tableRoot} tbody tr td:nth-child(2)`).first(),
    ).toHaveText("foxtrot");
    await expect(page.locator(`${tableRoot} nav p`)).toHaveText("6–10 of 25");

    await context.close();
  });

  test("a global refresh event reloads the table", async ({ page }) => {
    const refreshUrls = trackRefreshUrls(page);
    await page.evaluate(() =>
      window.dispatchEvent(new CustomEvent("refresh:data-table-demo")),
    );

    await expect.poll(() => refreshUrls.length).toBe(1);
  });

  test("items per page dropdown opens upward when it would overflow the viewport", async ({
    page,
  }) => {
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
    await page.locator(`${tableRoot} nav [role=button]`).click();

    const dropdown = page.locator(`${tableRoot} nav ul`);
    await expect(dropdown).toBeVisible();
    const dropdownBounds = await dropdown.evaluate((element) => {
      const rect = element.getBoundingClientRect();
      return { bottom: rect.bottom, top: rect.top };
    });
    const viewportHeight = await page.evaluate(() => window.innerHeight);
    expect(dropdownBounds.top).toBeGreaterThanOrEqual(0);
    expect(dropdownBounds.bottom).toBeLessThanOrEqual(viewportHeight);
  });
});

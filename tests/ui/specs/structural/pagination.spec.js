import { expect, test } from "@playwright/test";

const paginationSection = "#pagination-demo";

test.describe("Pagination @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${paginationSection.slice(1)}`);
    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "11–20 of 240",
    );
  });

  test("@smoke readout tracks the bound page number", async ({ page }) => {
    await page
      .locator(`${paginationSection} button[aria-label="Next page"]`)
      .click();
    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "21–30 of 240",
    );

    await page
      .locator(`${paginationSection} button[aria-label="Previous page"]`)
      .click();
    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "11–20 of 240",
    );
  });

  test("@smoke first and previous disable on the first page", async ({
    page,
  }) => {
    await page
      .locator(`${paginationSection} button[aria-label="First page"]`)
      .click();

    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "1–10 of 240",
    );
    await expect(
      page.locator(`${paginationSection} button[aria-label="First page"]`),
    ).toBeDisabled();
    await expect(
      page.locator(`${paginationSection} button[aria-label="Previous page"]`),
    ).toBeDisabled();
    await expect(
      page.locator(`${paginationSection} button[aria-label="Next page"]`),
    ).toBeEnabled();
  });

  test("@smoke last page disables next and last", async ({ page }) => {
    await page
      .locator(`${paginationSection} button[aria-label="Last page"]`)
      .click();

    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "231–240 of 240",
    );
    await expect(
      page.locator(`${paginationSection} button[aria-label="Next page"]`),
    ).toBeDisabled();
    await expect(
      page.locator(`${paginationSection} button[aria-label="Last page"]`),
    ).toBeDisabled();
  });

  test("@smoke items per page resets the page and calls OnChangeFunc", async ({
    page,
  }) => {
    await page.locator(`${paginationSection} nav [role=button]`).click();
    await page
      .locator(`${paginationSection} nav ul li label`)
      .filter({ hasText: "30" })
      .click();

    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "1–30 of 240",
    );
    await expect(page.locator("#pagination-demo-state span").nth(1)).toHaveText(
      "30",
    );
    await expect(page.locator("#pagination-demo-state span").nth(2)).toHaveText(
      "1",
    );
  });

  test("@smoke page strip marks the current page and shows first, last, and ellipses", async ({
    page,
  }) => {
    const strip = page.locator(`${paginationSection} nav [aria-current=page]`);
    await expect(strip).toHaveText("2");

    for (const pageLabel of ["1", "2", "3", "24"]) {
      await expect(
        page
          .locator(`${paginationSection} nav button`)
          .filter({ hasText: new RegExp(`^${pageLabel}$`) }),
      ).toBeVisible();
    }
    await expect(
      page.locator(`${paginationSection} nav`).getByText("…", { exact: true }),
    ).toBeVisible();
  });

  test("@smoke clicking a page number jumps to it and calls OnChangeFunc", async ({
    page,
  }) => {
    await page
      .locator(`${paginationSection} nav button`)
      .filter({ hasText: /^3$/ })
      .click();

    await expect(page.locator(`${paginationSection} nav p`)).toHaveText(
      "21–30 of 240",
    );
    await expect(
      page.locator(`${paginationSection} nav [aria-current=page]`),
    ).toHaveText("3");
    await expect(page.locator("#pagination-demo-state span").nth(2)).toHaveText(
      "1",
    );
  });

  test("@smoke page strip at the edges shows the neighbor pages without stray ellipses", async ({
    page,
  }) => {
    await page
      .locator(`${paginationSection} nav button[aria-label="First page"]`)
      .click();

    await expect(
      page.locator(`${paginationSection} nav [aria-current=page]`),
    ).toHaveText("1");
    await expect(
      page
        .locator(`${paginationSection} nav button`)
        .filter({ hasText: /^2$/ }),
    ).toBeVisible();
    await expect(
      page
        .locator(`${paginationSection} nav button`)
        .filter({ hasText: /^24$/ }),
    ).toBeVisible();
    await expect(
      page.locator(`${paginationSection} nav`).getByText("…", { exact: true }),
    ).toHaveCount(1);

    await page
      .locator(`${paginationSection} nav button`)
      .filter({ hasText: /^2$/ })
      .click();
    await page
      .locator(`${paginationSection} nav button`)
      .filter({ hasText: /^3$/ })
      .click();
    await page
      .locator(`${paginationSection} nav button`)
      .filter({ hasText: /^4$/ })
      .click();

    await expect(
      page.locator(`${paginationSection} nav [aria-current=page]`),
    ).toHaveText("4");
    await expect(
      page
        .locator(`${paginationSection} nav button`)
        .filter({ hasText: /^2$/ }),
    ).toBeVisible();
    await expect(
      page.locator(`${paginationSection} nav`).getByText("…", { exact: true }),
    ).toHaveCount(1);

    await page
      .locator(`${paginationSection} nav button[aria-label="Last page"]`)
      .click();

    await expect(
      page.locator(`${paginationSection} nav [aria-current=page]`),
    ).toHaveText("24");
    await expect(
      page
        .locator(`${paginationSection} nav button`)
        .filter({ hasText: /^23$/ }),
    ).toBeVisible();
    await expect(
      page.locator(`${paginationSection} nav`).getByText("…", { exact: true }),
    ).toHaveCount(1);
  });
});

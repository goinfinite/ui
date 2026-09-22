import { expect, test } from "@playwright/test";

const tableRoot = "#data-table-demo-table";

test.describe("DataTable accessibility @structural @a11y", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#data-table-demo");
    await expect(page.locator(`${tableRoot} tbody tr`).first()).toBeVisible();
  });

  test("sortable headers are named buttons inside th with aria-sort", async ({
    page,
  }) => {
    await expect(page.locator(`${tableRoot} thead th:has(button)`)).toHaveCount(
      3,
    );

    const nameHeader = page
      .locator(`${tableRoot} thead th`)
      .filter({ hasText: "Name" });
    await expect(nameHeader).toHaveAttribute("aria-sort", "ascending");
    await expect(nameHeader.locator("button")).toHaveText(/Name/);
  });

  test("sort buttons work with the keyboard", async ({ page }) => {
    const createdAtHeader = page
      .locator(`${tableRoot} thead th`)
      .filter({ hasText: "Created at" });
    await createdAtHeader.locator("button").focus();
    await page.keyboard.press("Enter");

    await expect(createdAtHeader).toHaveAttribute("aria-sort", "ascending");
  });

  test("selection checkboxes have accessible names", async ({ page }) => {
    await expect(
      page.locator(`${tableRoot} thead input[type=checkbox]`),
    ).toHaveAttribute("aria-label", "Select all rows on this page");
    await expect(
      page
        .locator(`${tableRoot} tbody tr`)
        .first()
        .locator("input[type=checkbox]"),
    ).toHaveAttribute("aria-label", "alpha");
  });

  test("pagination is a named navigation landmark with named controls", async ({
    page,
  }) => {
    const pagination = page.locator(
      `${tableRoot} nav[aria-label="Table pagination"]`,
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

  test("the selection count announces politely", async ({ page }) => {
    const liveRegion = page.locator(`${tableRoot} p[aria-live=polite]`);
    await expect(liveRegion).toHaveAttribute("aria-live", "polite");

    await page
      .locator(`${tableRoot} tbody tr`)
      .first()
      .locator("input[type=checkbox]")
      .click();
    await expect(liveRegion).toHaveText("1 row selected");
  });
});

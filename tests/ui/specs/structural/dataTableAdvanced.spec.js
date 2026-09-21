import { expect, test } from "@playwright/test";

const demoSection = "#data-table-demo";

function examplePanel(page, title) {
  return page.locator(`${demoSection} details`).filter({ hasText: title });
}

async function openExamplePanel(page, title) {
  const target = examplePanel(page, title);
  await target.locator("summary").click();
  return target;
}

test.describe("DataTable advanced examples @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#data-table-demo");
  });

  test("rows configuration renders a named selection column", async ({
    page,
  }) => {
    const target = await openExamplePanel(page, "Rows Configuration");
    const rowCheckbox = target.locator("tbody input[type=checkbox]").first();

    await expect(target.locator("tbody input[type=checkbox]")).toHaveCount(4);
    await expect(rowCheckbox).toHaveAttribute("aria-label", "alpha");

    const box = rowCheckbox.locator("xpath=following-sibling::span[1]");
    await expect(box).toHaveClass(/peer-checked:border-emerald-500/);
    await rowCheckbox.click();
    await expect
      .poll(() =>
        box.evaluate((element) => getComputedStyle(element).borderTopColor),
      )
      .toBe("rgb(16, 185, 129)");
  });

  test("filters panel starts with its chip", async ({ page }) => {
    const target = await openExamplePanel(page, "Filters");

    await expect(
      target.locator("span").filter({ hasText: "Status: running" }),
    ).toBeVisible();
  });

  test("search box panel binds the starting query", async ({ page }) => {
    const target = await openExamplePanel(page, "Search Box");
    const searchInput = target.locator("input[name=search]");

    await expect(searchInput).toHaveValue("alpha");
    await expect(
      searchInput.locator("xpath=ancestor::div[contains(@class, 'flex-1')]"),
    ).toHaveClass(/justify-center/);
  });

  test("header action button runs and the panel counts it", async ({
    page,
  }) => {
    const target = await openExamplePanel(page, "Header Actions");
    const syncButton = target.locator("button").filter({ hasText: "Sync" });

    await syncButton.click();
    await syncButton.click();

    await expect(target.locator('span[x-text="headerActionCount"]')).toHaveText(
      "2",
    );
  });

  test("bulk action appears only after a row is selected", async ({ page }) => {
    const target = await openExamplePanel(page, "Bulk Actions");
    const archiveButton = target
      .locator("button")
      .filter({ hasText: "Archive" });

    await expect(archiveButton).toBeHidden();
    await target.locator("tbody input[type=checkbox]").first().click();
    await expect(archiveButton).toBeVisible();
    await expect(target.locator("p[aria-live=polite]")).toHaveText(
      "1 row selected",
    );
  });

  test("row and column styling applies the header, stripe, cell, and row classes", async ({
    page,
  }) => {
    const target = await openExamplePanel(page, "Row and Column Styling");

    await expect(target.locator("thead tr")).toHaveClass(/bg-neutral-50\/5/);

    const bodyRows = target.locator("tbody tr");
    const oddRow = bodyRows.first();
    await expect(oddRow).toHaveClass(/odd:bg-neutral-50\/5/);
    await expect(oddRow.locator("td").first()).toHaveClass(/font-bold/);
    await expect(bodyRows.nth(2)).toHaveClass(/text-neutral-400/);

    const stripeBackground = await oddRow.evaluate(
      (row) => getComputedStyle(row).backgroundColor,
    );
    await oddRow.hover();
    await expect
      .poll(() =>
        oddRow.evaluate((row) => getComputedStyle(row).backgroundColor),
      )
      .not.toBe(stripeBackground);
  });
});

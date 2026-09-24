import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const demoSection = "#data-table-demo";

test.describe("DataTable advanced examples @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#data-table-demo");
  });

  test("rows configuration renders a named selection column", async ({
    page,
  }) => {
    const target = await openExamplePanel(
      page,
      demoSection,
      "Rows Configuration",
    );
    const rowCheckbox = target.locator("tbody input[type=checkbox]").first();

    await expect(target.locator("tbody input[type=checkbox]")).toHaveCount(4);
    await expect(rowCheckbox).toHaveAttribute("aria-label", "alpha");

    const box = rowCheckbox.locator("xpath=following-sibling::span[1]");
    await rowCheckbox.click();
    await expect
      .poll(() =>
        box.evaluate((element) => getComputedStyle(element).borderTopColor),
      )
      .toBe("rgb(16, 185, 129)");
  });

  test("filters panel starts with its chip", async ({ page }) => {
    const target = await openExamplePanel(page, demoSection, "Filters");

    await expect(
      target.locator("span").filter({ hasText: "Status: running" }),
    ).toBeVisible();
  });

  test("search box panel binds the starting query", async ({ page }) => {
    const target = await openExamplePanel(page, demoSection, "Search Box");
    const searchInput = target.locator("input[name=search]");

    await expect(searchInput).toHaveValue("alpha");
    const alignmentRow = searchInput.locator(
      "xpath=ancestor::div[contains(@class, 'flex-1')][1]",
    );
    const searchBoxFrame = searchInput.locator(
      "xpath=ancestor::div[contains(@class, 'w-72')][1]",
    );
    const [rowBox, frameBox] = await Promise.all([
      alignmentRow.boundingBox(),
      searchBoxFrame.boundingBox(),
    ]);
    expect(
      Math.abs(frameBox.x + frameBox.width / 2 - (rowBox.x + rowBox.width / 2)),
    ).toBeLessThanOrEqual(2);
  });

  test("header action button runs and the panel counts it", async ({
    page,
  }) => {
    const target = await openExamplePanel(page, demoSection, "Header Actions");
    const syncButton = target.locator("button").filter({ hasText: "Sync" });

    await syncButton.click();
    await syncButton.click();

    await expect(
      target.getByText("Header action clicks:").first(),
    ).toContainText("Header action clicks: 2");
  });

  test("bulk action appears only after a row is selected", async ({ page }) => {
    const target = await openExamplePanel(page, demoSection, "Bulk Actions");
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

  test("row and column styling applies the header, stripe, cell, and row styling", async ({
    page,
  }) => {
    const target = await openExamplePanel(
      page,
      demoSection,
      "Row and Column Styling",
    );

    const headerRow = target.locator("thead tr");
    expect(
      await headerRow.evaluate((row) => getComputedStyle(row).backgroundColor),
    ).not.toBe("rgba(0, 0, 0, 0)");

    const bodyRows = target.locator("tbody tr");
    const oddRow = bodyRows.first();
    const evenRow = bodyRows.nth(1);
    const [oddBackground, evenBackground] = await Promise.all([
      oddRow.evaluate((row) => getComputedStyle(row).backgroundColor),
      evenRow.evaluate((row) => getComputedStyle(row).backgroundColor),
    ]);
    expect(oddBackground).not.toBe(evenBackground);
    expect(
      await oddRow
        .locator("td")
        .first()
        .evaluate((cell) => getComputedStyle(cell).fontWeight),
    ).toBe("700");

    const [styledRowColor, plainRowColor] = await Promise.all([
      bodyRows.nth(2).evaluate((row) => getComputedStyle(row).color),
      evenRow.evaluate((row) => getComputedStyle(row).color),
    ]);
    expect(styledRowColor).not.toBe(plainRowColor);

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

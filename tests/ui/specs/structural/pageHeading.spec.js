import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const pageHeadingSection = "#page-heading-demo";

test.describe("PageHeading @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${pageHeadingSection.slice(1)}`);
  });

  test("@smoke page-level heading renders the h1 header block with the action slot", async ({
    page,
  }) => {
    const headerBlock = page
      .locator(pageHeadingSection)
      .locator("h1", { hasText: "Records" })
      .locator("xpath=ancestor::div[contains(@class,'flex-col')][1]");

    await expect(headerBlock.locator("h1")).toHaveText("Records");
    await expect(headerBlock.locator("i.ph-table")).toHaveCount(1);
    await expect(headerBlock.locator("p")).toHaveText(
      "Manage the server records and their lifecycle.",
    );
    await expect(
      headerBlock.getByRole("button", { name: "New Record" }),
    ).toBeVisible();
    await expect(
      headerBlock.getByRole("button", { name: "Export" }),
    ).toBeVisible();
  });

  test("title-only section heading omits the description and actions", async ({
    page,
  }) => {
    await openExamplePanel(page, pageHeadingSection, "Title Only");
    const headerBlock = page
      .locator(pageHeadingSection)
      .locator("h2", { hasText: "Title Only" })
      .locator("xpath=ancestor::div[contains(@class,'flex-col')][1]");

    await expect(headerBlock.locator("h2")).toHaveText("Title Only");
    await expect(headerBlock.locator("p")).toHaveCount(0);
    await expect(headerBlock.getByRole("button")).toHaveCount(0);
  });

  test("@smoke level page renders an h1 and level section renders an h2", async ({
    page,
  }) => {
    await openExamplePanel(page, pageHeadingSection, "Levels & Sizes");
    await expect(
      page.locator(pageHeadingSection).locator("h1", { hasText: "Page Level" }),
    ).toBeVisible();
    await expect(
      page
        .locator(pageHeadingSection)
        .locator("h2", { hasText: "Section Level" }),
    ).toBeVisible();
  });
});

import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const tagSection = "#tag-demo";

test.describe("Tag", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${tagSection.slice(1)}`);
  });

  test("@smoke outer and inner segments render their labels", async ({
    page,
  }) => {
    const tag = page
      .locator(`${tagSection} div.flex.w-fit`)
      .filter({ hasText: "Info" })
      .first();

    await expect(tag).toContainText("Info");
    await expect(tag).toContainText("Warning");
  });

  test("a bound inner value mirrors its state path", async ({ page }) => {
    const boundTag = page
      .locator(`${tagSection} div.flex.w-fit`)
      .filter({ hasText: "count" })
      .first();

    await expect(boundTag).toContainText("4");
  });

  test("@smoke side element clicks run their callback and update the bound value", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      tagSection,
      "Side Element Actions",
    );
    const actionTag = panel
      .locator("div.flex.w-fit")
      .filter({ hasText: "Add" })
      .first();
    const countReadout = panel
      .locator("p", { hasText: "Count:" })
      .locator("span");

    await expect(countReadout).toHaveText("0");

    await actionTag.getByText("Add", { exact: true }).click();
    await actionTag.getByText("Add", { exact: true }).click();
    await expect(countReadout).toHaveText("2");

    await actionTag.getByText("Remove", { exact: true }).click();
    await expect(countReadout).toHaveText("1");
  });

  test("bound side labels follow their state path", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      tagSection,
      "Side Element Actions",
    );
    const boundTag = panel.locator("div.flex.w-fit:has(i.ph-pencil)").first();

    await expect(boundTag).toContainText("Add");
    await expect(boundTag).toContainText("Remove");

    await panel.getByRole("button", { name: "Rename Actions" }).click();

    await expect(boundTag).toContainText("Increase");
    await expect(boundTag).toContainText("Decrease");
  });

  test("@smoke the remove button runs OnRemoveFunc", async ({ page }) => {
    const panel = await openExamplePanel(page, tagSection, "Removable Tags");
    const tag = panel
      .locator("span")
      .filter({ hasText: "Status: running" })
      .first();

    await expect(tag).toBeVisible();

    await tag.getByRole("button", { name: "Remove status filter" }).click();
    await expect(tag).toBeHidden();

    await panel.getByRole("button", { name: "Restore" }).click();
    await expect(tag).toBeVisible();
  });

  test("@browser text case transforms the tag labels", async ({ page }) => {
    const panel = await openExamplePanel(page, tagSection, "Text Case");

    await expect(panel.getByText("Info (lower)")).toHaveCSS(
      "text-transform",
      "lowercase",
    );
    await expect(panel.getByText("Info (upper)")).toHaveCSS(
      "text-transform",
      "uppercase",
    );
    await expect(panel.getByText("Info (capitalize)")).toHaveCSS(
      "text-transform",
      "capitalize",
    );
  });
});

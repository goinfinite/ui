import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const accordionSection = "#accordion-demo";

test.describe("Accordion", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${accordionSection.slice(1)}`);
  });

  test("@browser @smoke single-open keeps one section open at a time", async ({
    page,
  }) => {
    const singleOpenPanel = await openExamplePanel(
      page,
      accordionSection,
      "Single Open",
    );

    const sectionItems = singleOpenPanel.locator("details");
    await sectionItems.nth(0).locator("summary").click();
    await expect(sectionItems.nth(0)).toHaveAttribute("open", "");

    await sectionItems.nth(1).locator("summary").click();
    await expect(sectionItems.nth(0)).not.toHaveAttribute("open", "");
    await expect(sectionItems.nth(1)).toHaveAttribute("open", "");
  });

  test("@smoke first item open renders the first section expanded", async ({
    page,
  }) => {
    const firstItemPanel = await openExamplePanel(
      page,
      accordionSection,
      "First Item Open",
    );

    const sectionItems = firstItemPanel.locator("details");
    await expect(sectionItems.nth(0)).toHaveAttribute("open", "");
    await expect(sectionItems.nth(1)).not.toHaveAttribute("open", "");
    await expect(sectionItems.nth(2)).not.toHaveAttribute("open", "");
  });
});

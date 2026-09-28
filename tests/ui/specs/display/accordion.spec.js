import { expect, test } from "@playwright/test";

const accordionSection = "#accordion-demo";

test.describe("Accordion", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${accordionSection.slice(1)}`);
  });

  test("@smoke open item content reads as a cutout of the parent surface", async ({
    page,
  }) => {
    const firstLiveItem = page.locator(`${accordionSection} details`).first();
    await firstLiveItem.locator("summary").click();
    const cutout = firstLiveItem.locator("summary + div");
    const cutoutStyles = await cutout.evaluate((element) => {
      const computed = getComputedStyle(element);
      return {
        backgroundColor: computed.backgroundColor,
        borderRadius: computed.borderTopLeftRadius,
        marginTop: computed.marginTop,
        padding: computed.padding,
      };
    });

    expect(cutoutStyles.backgroundColor).toBe("rgb(23, 23, 23)");
    expect(cutoutStyles.borderRadius).toBe("8px");
    expect(cutoutStyles.marginTop).toBe("16px");
    expect(cutoutStyles.padding).toBe("16px");
  });

  test("@smoke only the outer items keep rounded corners", async ({ page }) => {
    const liveAccordionItems = page
      .locator(`${accordionSection} details`)
      .filter({ hasText: "This is the content for section" });

    await expect
      .poll(async () => {
        const cornerStyles = await liveAccordionItems.evaluateAll((elements) =>
          elements.slice(0, 3).map((element) => {
            const computed = getComputedStyle(element);
            return [
              computed.borderTopLeftRadius,
              computed.borderBottomLeftRadius,
            ].join("/");
          }),
        );
        return cornerStyles;
      })
      .toEqual(["6px/0px", "0px/0px", "0px/6px"]);
  });

  test("@smoke single-open keeps one section open at a time", async ({
    page,
  }) => {
    const singleOpenPanel = page
      .locator(`${accordionSection} details`)
      .filter({ hasText: "Single Open" })
      .first();
    await singleOpenPanel.locator("summary").first().click();

    const sectionItems = singleOpenPanel.locator("details");
    await sectionItems.nth(0).locator("summary").click();
    await expect(sectionItems.nth(0)).toHaveAttribute("open", "");

    await sectionItems.nth(1).locator("summary").click();
    await expect(sectionItems.nth(0)).not.toHaveAttribute("open", "");
    await expect(sectionItems.nth(1)).toHaveAttribute("open", "");
  });
});

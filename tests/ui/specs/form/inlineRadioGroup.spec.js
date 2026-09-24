import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const radioGroupSection = "#inline-radio-group-demo";

test.describe("InlineRadioGroup", () => {
  test("@smoke wrapper carries no top margin so it aligns with sibling controls", async ({
    page,
  }) => {
    await page.goto(`/index.html#${radioGroupSection.slice(1)}`);
    const group = page.locator(`${radioGroupSection} .border-1`).first();
    await expect(group).toBeVisible();
    await expect
      .poll(async () => group.evaluate((el) => getComputedStyle(el).marginTop))
      .toBe("0px");
  });

  test("@smoke large radios grow the group and keep bottom clearance", async ({
    page,
  }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    await openExamplePanel(page, radioGroupSection, "Sizes");
    const group = page
      .locator(`${radioGroupSection} details`)
      .filter({ hasText: "Pick a size" })
      .locator("div.border-1");
    await expect(group).toBeVisible();

    const geometry = await group.evaluate((el) => {
      const groupRect = el.getBoundingClientRect();
      const rowRect = el.querySelector("div.flex-row").getBoundingClientRect();
      return {
        paddingBottom: getComputedStyle(el).paddingBottom,
        groupBottom: groupRect.bottom,
        rowBottom: rowRect.bottom,
      };
    });

    expect(geometry.paddingBottom).toBe("8px");
    expect(geometry.groupBottom - geometry.rowBottom).toBeGreaterThanOrEqual(8);
  });
});

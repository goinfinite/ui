import { expect, test } from "@playwright/test";

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
});

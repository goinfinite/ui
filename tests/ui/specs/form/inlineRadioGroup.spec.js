import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const radioGroupSection = "#inline-radio-group-demo";

test.describe("InlineRadioGroup", () => {
  test("@smoke label notches into the fieldset border like the other fields", async ({
    page,
  }) => {
    await page.goto(`/index.html#${radioGroupSection.slice(1)}`);
    const fieldset = page.locator(`${radioGroupSection} fieldset`).first();
    await expect(fieldset).toBeVisible();
    const legend = fieldset.locator("legend");
    await expect(legend).toHaveText("Select an option");
    const geometry = await fieldset.evaluate((el) => {
      const legendRect = el.querySelector("legend").getBoundingClientRect();
      const contentRect = el
        .querySelector("div.flex-row")
        .getBoundingClientRect();
      return {
        legendTop: Math.round(legendRect.top),
        fieldsetTop: Math.round(el.getBoundingClientRect().top),
        contentTop: Math.round(contentRect.top),
      };
    });
    expect(geometry.legendTop).toBe(geometry.fieldsetTop);
    expect(geometry.contentTop).toBeGreaterThan(geometry.legendTop);
  });

  test("@smoke large radios grow the group and keep bottom clearance", async ({
    page,
  }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    await openExamplePanel(page, radioGroupSection, "Sizes");
    const fieldset = page
      .locator(`${radioGroupSection} details`)
      .filter({ hasText: "Pick a size" })
      .locator("fieldset");
    await expect(fieldset).toBeVisible();

    const geometry = await fieldset.evaluate((el) => {
      const content = el.querySelector("div.flex-row");
      const contentRect = content.getBoundingClientRect();
      const radioRect = content.querySelector("label").getBoundingClientRect();
      return {
        paddingBottom: getComputedStyle(content).paddingBottom,
        contentBottom: contentRect.bottom,
        radioBottom: radioRect.bottom,
      };
    });

    expect(geometry.paddingBottom).toBe("6px");
    expect(
      geometry.contentBottom - geometry.radioBottom,
    ).toBeGreaterThanOrEqual(6);
  });

  test("vertical orientation stacks the radios", async ({ page }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    const panel = await openExamplePanel(
      page,
      radioGroupSection,
      "Vertical Orientation",
    );
    const content = panel.locator("fieldset > div").first();
    const flexDirection = await content.evaluate(
      (element) => getComputedStyle(element).flexDirection,
    );
    expect(flexDirection).toBe("column");
  });
});

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const inputFieldSection = "#input-field-demo";

test.describe("InputField", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    await page.waitForFunction(() => window.Alpine !== undefined);
    await openExamplePanel(page, inputFieldSection, "Affixes");
  });

  test("@smoke a long suffix stays on one line inside the field border", async ({
    page,
  }) => {
    const fieldset = page
      .locator(inputFieldSection)
      .locator("fieldset")
      .filter({ hasText: "Slug" });
    const suffix = fieldset.locator("div.truncate").last();
    await expect(suffix).toHaveText("ui.demo.goinfinite.net");

    const geometry = await fieldset.evaluate((el) => {
      const fieldRect = el.getBoundingClientRect();
      const suffixRect = el
        .querySelector("div.truncate:last-of-type")
        .getBoundingClientRect();
      return {
        fieldBottom: fieldRect.bottom,
        fieldLeft: fieldRect.left,
        fieldRight: fieldRect.right,
        suffixBottom: suffixRect.bottom,
        suffixLeft: suffixRect.left,
        suffixRight: suffixRect.right,
        suffixHeight: suffixRect.height,
        lineHeight: parseFloat(
          getComputedStyle(el.querySelector("input")).lineHeight,
        ),
      };
    });

    expect(geometry.suffixHeight).toBeLessThan(geometry.lineHeight * 2);
    expect(geometry.suffixBottom).toBeLessThanOrEqual(geometry.fieldBottom);
    expect(geometry.suffixLeft).toBeGreaterThanOrEqual(geometry.fieldLeft);
    expect(geometry.suffixRight).toBeLessThanOrEqual(geometry.fieldRight);
  });

  test("@smoke a long suffix truncates and leaves the input usable", async ({
    page,
  }) => {
    const input = page.locator(`${inputFieldSection} input[name=slug]`);
    await input.fill("hello");

    await expect(input).toHaveValue("hello");
    const suffixScroll = await page
      .locator(inputFieldSection)
      .locator("fieldset")
      .filter({ hasText: "Slug" })
      .locator("div.truncate")
      .last()
      .evaluate((el) => ({
        clientWidth: el.clientWidth,
        scrollWidth: el.scrollWidth,
      }));
    expect(suffixScroll.scrollWidth).toBeGreaterThan(suffixScroll.clientWidth);
  });

  test("@smoke text case transforms the floating label", async ({ page }) => {
    const panel = await openExamplePanel(page, inputFieldSection, "Label Case");
    const fields = panel.locator("fieldset");
    const cases = [
      ["none", "none"],
      ["lower", "lowercase"],
      ["upper", "uppercase"],
      ["capitalize", "capitalize"],
    ];

    for (let index = 0; index < cases.length; index++) {
      const [caseName, transform] = cases[index];
      const field = fields.nth(index);
      const legend = field.locator("legend");
      await expect(legend).toHaveText(`Network Interface (${caseName})`);
      expect(
        await legend.evaluate(
          (element) => getComputedStyle(element).textTransform,
        ),
      ).toBe(transform);
      expect(
        await field
          .locator("input")
          .evaluate(
            (element) =>
              getComputedStyle(element, "::placeholder").textTransform,
          ),
      ).toBe(transform);
    }
  });
});

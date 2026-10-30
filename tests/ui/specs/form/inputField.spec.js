import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const inputFieldSection = "#input-field-demo";

test.describe("InputField", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html", { waitUntil: "domcontentloaded" });
    await page.waitForFunction(() => window.Alpine !== undefined);
    await openExamplePanel(page, inputFieldSection, "Affixes");
  });

  test("@browser @smoke a long suffix truncates and leaves the input usable", async ({
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

  test("@browser @smoke a state path drives the legend label", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, inputFieldSection, "Labels");
    const dynamicField = panel.locator(
      'fieldset:has(input[name="runtimeField"])',
    );
    const legend = dynamicField.locator("legend");
    await expect(legend).toHaveText("Admin Username");

    await panel.getByRole("button", { name: "Use the password label" }).click();
    await expect(legend).toHaveText("Admin Password");
  });

  test("@browser @smoke text case transforms the floating label", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, inputFieldSection, "Labels");
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

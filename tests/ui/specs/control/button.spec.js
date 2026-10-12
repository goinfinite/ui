import { openEveryExamplePanel, openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const buttonSection = "#button-demo";

test.describe("Button", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${buttonSection.slice(1)}`);
  });

  test("@control every button carries a type attribute", async ({ page }) => {
    await openEveryExamplePanel(page, buttonSection);

    const buttons = page.locator(`${buttonSection} button`);
    const buttonCount = await buttons.count();

    expect(buttonCount).toBeGreaterThan(0);
    for (let index = 0; index < buttonCount; index += 1) {
      await expect(buttons.nth(index)).toHaveAttribute(
        "type",
        /^(button|submit)$/,
      );
    }
  });

  test("@control a plain button defaults to type=button", async ({ page }) => {
    const newRecordButton = page
      .locator(`${buttonSection} button`)
      .filter({ hasText: "New Record" })
      .first();

    await expect(newRecordButton).toHaveAttribute("type", "button");
  });

  test("@control IsSubmit renders a submit button", async ({ page }) => {
    const submitPanel = await openExamplePanel(
      page,
      buttonSection,
      "Full-Width & Submit",
    );

    await expect(submitPanel.locator("button[type=submit]")).toHaveCount(1);
  });

  test("@control re-embedding the tooltip state script does not throw", async ({
    page,
  }) => {
    const pageErrors = [];
    page.on("pageerror", (error) => pageErrors.push(error.message));

    await page.evaluate(() => {
      const tooltipScript = Array.from(
        document.querySelectorAll("script"),
      ).find((script) => script.textContent.includes("showTooltip"));
      const reembeddedScript = document.createElement("script");
      reembeddedScript.textContent = tooltipScript.textContent;
      document.body.append(reembeddedScript);
    });

    await page.waitForTimeout(100);
    expect(pageErrors).toEqual([]);
  });
});

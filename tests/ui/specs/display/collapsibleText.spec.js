import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const collapsibleSection = "#collapsible-text-demo";
const textSelector = `${collapsibleSection} [x-ref=textContent]`;
const toggleSelector = `${collapsibleSection} [x-data^=collapsibleText] button`;

test.describe("CollapsibleText", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${collapsibleSection.slice(1)}`);
    await expect(page.locator(toggleSelector).first()).toContainText(
      "show more",
    );
  });

  test("@smoke the text truncates with the affordance inline after the ellipsis", async ({
    page,
  }) => {
    const text = page.locator(textSelector).first();
    const toggle = page.locator(toggleSelector).first();

    await expect(text).toContainText(/…\s*show more/);
    await expect(text).not.toContainText("bytes=642");

    await toggle.dispatchEvent("click");

    await expect(toggle).toContainText("show less");
    await expect(text).toContainText("bytes=642");
    await expect(text).not.toContainText("…");
  });

  test("a second click collapses the text back", async ({ page }) => {
    const text = page.locator(textSelector).first();
    const toggle = page.locator(toggleSelector).first();

    await toggle.dispatchEvent("click");
    await expect(text).not.toContainText("…");

    await toggle.dispatchEvent("click");

    await expect(toggle).toContainText("show more");
    await expect(text).toContainText(/…\s*show more/);
  });

  test("custom labels replace the affordance text", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      collapsibleSection,
      "Custom Labels",
    );
    const toggle = panel.locator("[x-data^=collapsibleText] button").first();

    await expect(toggle).toContainText("Expand the line");

    await toggle.dispatchEvent("click");

    await expect(toggle).toContainText("Collapse the line");
  });

  test("@browser a normal line height still truncates the text", async ({
    page,
  }) => {
    const text = page.locator(textSelector).first();

    await text.evaluate((element) => {
      element.style.lineHeight = "normal";
      element.style.width = "200px";
      element.style.flex = "none";
    });

    await expect(text).toContainText(/…\s*show more/);
    await expect(text).not.toContainText("bytes=642");
  });
});

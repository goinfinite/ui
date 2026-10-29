import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const textAreaSection = "#text-area-demo";
const collapsedHeightPx = 144;
const expandedHeightPx = 432;

async function waitForAlpineHydration(page) {
  await page.waitForFunction((section) => {
    const fieldset = document.querySelector(`${section} fieldset`);
    return window.Alpine && fieldset.classList.contains("mt-0");
  }, textAreaSection);
}

async function heightOf(page, index) {
  return page
    .locator(`${textAreaSection} textarea`)
    .nth(index)
    .evaluate((ta) => ta.getBoundingClientRect().height);
}

async function toggleExpand(page, index) {
  await page.locator(`${textAreaSection} textarea`).nth(index).hover();
  await page
    .locator(`${textAreaSection} fieldset`)
    .nth(index)
    .locator(".ph-arrows-out, .ph-arrows-in")
    .first()
    .click();
}

async function expectStableHeight(page, index, expected) {
  await expect
    .poll(async () => heightOf(page, index), { timeout: 5000 })
    .toBeCloseTo(expected, 1);
}

test.describe("TextArea", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${textAreaSection.slice(1)}`);
    await waitForAlpineHydration(page);
  });

  test("@browser @smoke expand toggle grows the box and collapse restores it", async ({
    page,
  }) => {
    await toggleExpand(page, 0);
    await expectStableHeight(page, 0, expandedHeightPx);
    await toggleExpand(page, 0);
    await expectStableHeight(page, 0, collapsedHeightPx);
  });

  test("@browser read-only instance expands without a state path", async ({
    page,
  }) => {
    await openExamplePanel(page, textAreaSection, "Value & Read-Only");
    await toggleExpand(page, 1);
    await expectStableHeight(page, 1, expandedHeightPx);
  });

  test("@browser action buttons reveal on focus and expand with the keyboard", async ({
    page,
  }) => {
    const fieldset = page.locator(`${textAreaSection} fieldset`).first();
    const toggleHeightButton = fieldset.locator(
      'button[aria-label="Toggle text area height"]',
    );

    await expect(toggleHeightButton).toBeAttached();
    await expect(toggleHeightButton).toBeHidden();
    await fieldset.locator("textarea").focus();
    await expect(toggleHeightButton).toBeVisible();

    await toggleHeightButton.press("Enter");
    await expectStableHeight(page, 0, expandedHeightPx);
  });

  test("@browser label state path updates the legend, placeholder, and aria label", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, textAreaSection, "State Paths");
    const fieldset = panel.locator("fieldset");
    const textarea = fieldset.locator("textarea");

    await expect(fieldset.locator("legend")).toHaveText("Description");
    await panel.getByRole("button", { name: "Change Label" }).click();
    await expect(fieldset.locator("legend")).toHaveText("Message");
    await expect(textarea).toHaveAttribute("placeholder", "Message");
    await expect(textarea).toHaveAttribute("aria-label", "Message");
  });
});

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const toggleSection = "#toggle-switch-demo";
const booleanStateInput = `${toggleSection} input[type=hidden][name=notificationsEnabled]`;

function toggleByLabel(page, label) {
  return page
    .locator(`${toggleSection} label`)
    .filter({ hasText: label })
    .first();
}

async function labelOffsetFromTrack(page, label) {
  const toggle = toggleByLabel(page, label);
  const trackBox = await toggle.locator("div").first().boundingBox();
  const labelBox = await toggle.locator("span").first().boundingBox();
  return labelBox.x - trackBox.x;
}

async function startEdgeOf(page, label) {
  const toggle = toggleByLabel(page, label);
  const trackBox = await toggle.locator("div").first().boundingBox();
  const labelBox = await toggle.locator("span").first().boundingBox();
  return Math.min(trackBox.x, labelBox.x);
}

test.describe("ToggleSwitch", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${toggleSection.slice(1)}`);
    await expect(page.locator(booleanStateInput)).toHaveValue("false");
  });

  test("@smoke clicking the switch updates the bound boolean", async ({
    page,
  }) => {
    await toggleByLabel(page, "Enable notifications")
      .locator("div")
      .first()
      .click();
    await expect(page.locator(booleanStateInput)).toHaveValue("true");
  });

  test("@smoke default label renders right of the track", async ({ page }) => {
    await expect
      .poll(() => labelOffsetFromTrack(page, "Enable notifications"))
      .toBeGreaterThan(0);
  });

  test("label position left renders the label before the track", async ({
    page,
  }) => {
    await openExamplePanel(page, toggleSection, "Label Position");
    await expect
      .poll(() => labelOffsetFromTrack(page, "Label on left"))
      .toBeLessThan(0);
  });

  test("label position left keeps the switch at the column start", async ({
    page,
  }) => {
    await openExamplePanel(page, toggleSection, "Label Position");
    await expect
      .poll(async () => {
        const leftPositionStart = await startEdgeOf(page, "Label on left");
        const rightPositionStart = await startEdgeOf(page, "Label on right");
        return Math.abs(leftPositionStart - rightPositionStart);
      })
      .toBeLessThanOrEqual(2);
  });

  test("label position left keeps the label and track together", async ({
    page,
  }) => {
    await openExamplePanel(page, toggleSection, "Label Position");
    const toggle = toggleByLabel(page, "Label on left");
    const labelBox = await toggle.boundingBox();
    const parentBox = await toggle.locator("xpath=..").boundingBox();
    expect(labelBox.width).toBeLessThan(parentBox.width / 2);
  });
});

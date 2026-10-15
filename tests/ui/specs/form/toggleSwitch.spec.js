import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

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

  test("@browser @smoke default label renders right of the track", async ({
    page,
  }) => {
    await expect
      .poll(() => labelOffsetFromTrack(page, "Enable notifications"))
      .toBeGreaterThan(0);
  });

  test("@browser label position left renders the label before the track", async ({
    page,
  }) => {
    await openExamplePanel(page, toggleSection, "Label Position");
    await expect
      .poll(() => labelOffsetFromTrack(page, "Label on left"))
      .toBeLessThan(0);
  });
});

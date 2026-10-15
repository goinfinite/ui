import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const alertSection = "#alert-demo";

function alertWith(page, title) {
  return page.locator(`${alertSection} [role=alert]`, { hasText: title });
}

test.describe("Alert", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${alertSection.slice(1)}`);
  });

  test("@smoke each variation renders its tone icon", async ({ page }) => {
    const variationIcons = [
      ["Success Alert", "ph-check-circle"],
      ["Warning Alert", "ph-warning"],
      ["Error Alert", "ph-x-circle"],
      ["Info Alert", "ph-info"],
    ];

    for (const [title, iconClass] of variationIcons) {
      await expect(
        alertWith(page, title).locator(`i.${iconClass}`),
      ).toBeVisible();
    }
  });

  test("@smoke the close button dismisses the alert and runs OnCloseFunc", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, alertSection, "Close & Timeout");
    const closableAlert = panel.locator("[role=alert]", {
      hasText: "Closable alert",
    });
    const closedReadout = panel
      .locator("p", { hasText: "Closed:" })
      .locator("span")
      .first();

    await expect(closableAlert).toBeVisible();
    await expect(closedReadout).toHaveText("false");

    await closableAlert.getByRole("button").click();

    await expect(closableAlert).toBeHidden();
    await expect(closedReadout).toHaveText("true");
  });

  test("@smoke the timeout dismisses the alert and runs OnTimeoutFunc", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, alertSection, "Close & Timeout");
    const timedAlert = panel.locator("[role=alert]", {
      hasText: "Self-dismissing alert",
    });
    const timedOutReadout = panel
      .locator("p", { hasText: "Timed out:" })
      .locator("span")
      .last();

    await expect(timedAlert).toBeVisible();
    await expect(timedOutReadout).toHaveText("false");

    await expect(timedAlert).toBeHidden({ timeout: 12000 });
    await expect(timedOutReadout).toHaveText("true");
  });

  test("state-bound title updates at runtime", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      alertSection,
      "State-Bound Content",
    );
    const dynamicAlert = panel.locator("[role=alert]").first();

    await expect(dynamicAlert).toContainText("Dynamic Content Alert");
    await expect(dynamicAlert).toContainText("Dynamic alert content!");

    await panel.getByRole("button", { name: "Change Content" }).click();

    await expect(dynamicAlert).toContainText("Title updated at");
    await expect(dynamicAlert).toContainText("Message updated at");
  });

  test("a state path toggles the close button", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      alertSection,
      "State-Bound Content",
    );
    const conditionalAlert = panel.locator("[role=alert]").nth(1);
    const closeButton = conditionalAlert.getByRole("button");

    await expect(conditionalAlert).toContainText("Conditional Close Button");
    await expect(closeButton).toBeVisible();

    await panel.getByRole("button", { name: "Toggle Close Button" }).click();
    await expect(closeButton).toBeHidden();

    await panel.getByRole("button", { name: "Toggle Close Button" }).click();
    await expect(closeButton).toBeVisible();
  });
});

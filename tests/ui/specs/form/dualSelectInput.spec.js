import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const dualSelectSection = "#dual-select-input-demo";
const triggerSelector = `${dualSelectSection} [x-ref=selectTrigger]`;
const dropdownSelector = `${dualSelectSection} [x-ref=selectDropdown]`;
const valueInputSelector = `${dualSelectSection} input[type=text]`;
const liveStateSelector = `${dualSelectSection} p span`;

test.describe("DualSelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${dualSelectSection.slice(1)}`);
    await expect(page.locator(triggerSelector)).toContainText("port");
    await expect(page.locator(valueInputSelector)).toHaveValue("80");
  });

  test("@smoke picking an option updates the unit half", async ({ page }) => {
    const trigger = page.locator(triggerSelector);
    const dropdown = page.locator(dropdownSelector);

    await trigger.click();
    await expect(dropdown).toBeVisible();

    await dropdown.locator("label").filter({ hasText: "socket" }).click();

    await expect(dropdown).toBeHidden();
    await expect(trigger).toContainText("socket");
    await expect(page.locator(liveStateSelector).first()).toHaveText("socket");
  });

  test("@smoke typing a value updates the value half", async ({ page }) => {
    const valueInput = page.locator(valueInputSelector);

    await valueInput.fill("8443");

    await expect(valueInput).toHaveValue("8443");
    await expect(page.locator(liveStateSelector).nth(1)).toHaveText("8443");
  });

  test("@smoke a typed value that matches no option shows Custom", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      dualSelectSection,
      "Shared Slot",
    );
    const panelTrigger = panel.locator("[x-ref=selectTrigger]");
    const panelDropdown = panel.locator("[x-ref=selectDropdown]");
    const panelInput = panel.locator("input[type=text]");

    await expect(panelTrigger).toContainText("Every 5 minutes");
    await panelInput.fill("3,7");

    await expect(panelTrigger).toContainText("Custom");

    await panelTrigger.click();
    await expect(
      panelDropdown.locator("label").filter({ hasText: "Custom" }),
    ).toBeVisible();
  });

  test("@browser a preset pick writes the shared path on both halves", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      dualSelectSection,
      "Shared Slot",
    );
    const panelTrigger = panel.locator("[x-ref=selectTrigger]");
    const panelInput = panel.locator("input[type=text]");

    await panelTrigger.click();
    await panel
      .locator("[x-ref=selectDropdown] label")
      .filter({ hasText: "Every 15 minutes" })
      .click();

    await expect(panelInput).toHaveValue("*/15");
    await expect(panel.locator("p span")).toHaveText("*/15");
  });

  test("@browser the value input spans the fieldset height inside the border", async ({
    page,
  }) => {
    const valueInput = page.locator(valueInputSelector);
    const geometry = await valueInput.evaluate((element) => {
      const fieldsetRect = element.closest("fieldset").getBoundingClientRect();
      const rowRect = element.parentElement.getBoundingClientRect();
      const inputRect = element.getBoundingClientRect();
      return {
        fieldsetTop: fieldsetRect.top,
        fieldsetBottom: fieldsetRect.bottom,
        rowTop: rowRect.top,
        rowBottom: rowRect.bottom,
        inputTop: inputRect.top,
        inputBottom: inputRect.bottom,
      };
    });

    expect(geometry.inputTop).toBeGreaterThan(geometry.fieldsetTop);
    expect(geometry.inputTop).toBeLessThan(geometry.rowTop);
    expect(
      Math.abs(geometry.inputBottom - geometry.rowBottom),
    ).toBeLessThanOrEqual(1);
    expect(
      Math.abs(geometry.inputBottom - (geometry.fieldsetBottom - 1)),
    ).toBeLessThanOrEqual(1);
  });

  test("@browser arrow down opens with the checked option focused and enter commits", async ({
    page,
  }) => {
    const trigger = page.locator(triggerSelector);
    const dropdown = page.locator(dropdownSelector);

    await trigger.focus();
    await page.keyboard.press("ArrowDown");
    await expect(dropdown).toBeVisible();
    await expect(dropdown.locator("input:checked")).toHaveValue("port");

    await page.keyboard.press("Enter");
    await expect(dropdown).toBeHidden();
    await expect(trigger).toBeFocused();
  });
});

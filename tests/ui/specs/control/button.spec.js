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

  test("@control a tracked submit shows the HTMX request result", async ({
    page,
  }) => {
    const actionPanel = await openExamplePanel(
      page,
      buttonSection,
      "Action States",
    );
    const savingButton = actionPanel
      .locator("form[hx-get='assets/ui-logo.svg'] button")
      .first();
    const failingButton = actionPanel
      .locator("form[hx-get='assets/missing-action-target'] button")
      .first();

    await savingButton.click();
    await failingButton.click();

    await expect(savingButton).toHaveText("Saved");
    await expect(failingButton).toHaveText("Save failed");
    await expect(savingButton).toHaveText("Save", { timeout: 4000 });
  });

  test("@control a tracked submit refuses a second submit while pending", async ({
    page,
  }) => {
    const actionPanel = await openExamplePanel(
      page,
      buttonSection,
      "Action States",
    );
    const savingForm = actionPanel
      .locator("form[hx-get='assets/ui-logo.svg']")
      .first();

    const passedSubmitCount = await savingForm.evaluate((formElement) => {
      let submitCount = 0;
      formElement.addEventListener("submit", () => {
        submitCount += 1;
      });
      formElement.requestSubmit();
      formElement.requestSubmit();
      return submitCount;
    });

    expect(passedSubmitCount).toBe(1);
  });

  test("@control a tracked click follows the promise it returns", async ({
    page,
  }) => {
    await openExamplePanel(page, buttonSection, "Action States");
    const syncButton = page.locator("#button-demo-sync");
    const brokenSyncButton = page.locator("#button-demo-broken-sync");

    await syncButton.click();
    await brokenSyncButton.click();

    await expect(syncButton).toHaveText("Syncing");
    await expect(syncButton).toBeDisabled();
    await expect(syncButton).toHaveAttribute("aria-busy", "true");
    await expect(syncButton).toHaveText("Synced", { timeout: 4000 });
    await expect(brokenSyncButton).toHaveText("Sync failed", { timeout: 4000 });
    await expect(syncButton).toHaveText("Sync", { timeout: 4000 });
    await expect(syncButton).toBeEnabled();
  });

  test("@control the pending path disables every bound button", async ({
    page,
  }) => {
    await openExamplePanel(page, buttonSection, "Action States");
    const startButton = page.locator("#button-demo-start-sync");
    const deleteButton = page.locator("#button-demo-delete-account");

    await startButton.click();

    await expect(deleteButton).toHaveText("Wait for the sync");
    await expect(deleteButton).toBeDisabled();
    await expect(deleteButton).toHaveText("Delete account", { timeout: 4000 });
    await expect(deleteButton).toBeEnabled();
  });
});

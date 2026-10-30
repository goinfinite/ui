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

  test("@control icon top renders above the label and combines with side icons", async ({
    page,
  }) => {
    const iconTopPanel = await openExamplePanel(page, buttonSection, "Icons");
    const iconTopButton = iconTopPanel
      .locator("button")
      .filter({ hasText: "Upload" })
      .first();

    await expect(iconTopButton).toHaveClass(/flex-col/);
    await expect(iconTopButton.locator("i.ph-cloud-arrow-up")).toHaveCount(1);
    await expect(iconTopButton.locator("i.ph-folder")).toHaveCount(1);
    await expect(iconTopButton.locator("i.ph-arrow-right")).toHaveCount(1);

    const iconOnlyButton = iconTopPanel
      .getByRole("button", { name: "Add" })
      .first();
    await expect(iconOnlyButton).toHaveClass(/flex-col/);
    await expect(iconOnlyButton.locator("i.ph-plus")).toHaveCount(1);
    await expect(iconOnlyButton.locator("span")).toHaveCount(0);

    const actionOnlyButton = iconTopPanel.getByRole("button", {
      name: "Favorite",
    });
    await expect(actionOnlyButton).toHaveClass(/flex-col/);
    await expect(actionOnlyButton.locator("span.flex.flex-row")).toBeHidden();
  });

  test("@control countdown disables the button until zero and restores the label", async ({
    page,
  }) => {
    const countdownPanel = await openExamplePanel(
      page,
      buttonSection,
      "Countdown",
    );
    const countdownButton = countdownPanel.locator("#button-demo-countdown");

    await expect(countdownButton).toBeDisabled();
    await expect(countdownButton).toContainText("Delete (");
    await expect(countdownButton).toBeEnabled({ timeout: 8000 });
    await expect(countdownButton).toHaveText("Delete");
  });

  test("@control a trigger path restarts the countdown", async ({ page }) => {
    const countdownPanel = await openExamplePanel(
      page,
      buttonSection,
      "Countdown",
    );
    const triggeredButton = countdownPanel.locator(
      "#button-demo-triggered-countdown",
    );

    await expect(triggeredButton).toBeEnabled({ timeout: 6000 });
    await countdownPanel
      .getByRole("button", { name: "Restart countdown" })
      .click();
    await expect(triggeredButton).toBeDisabled();
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

  test("@control dynamic labels follow the action state", async ({ page }) => {
    await openExamplePanel(page, buttonSection, "Action States");
    const successButton = page.locator("#button-demo-dynamic-labels");
    const failureButton = page.locator("#button-demo-dynamic-failure-labels");
    const idleLabelButton = page.locator("#button-demo-dynamic-idle-label");

    await expect(successButton).toHaveText("Dynamic idle");
    await expect(idleLabelButton).toHaveText("Dynamic idle");

    await page.evaluate(() => {
      const actionScope = document.querySelector(
        "#button-demo [x-data='buttonDemoActions']",
      );
      Alpine.$data(actionScope).dynamicButtonLabel = "Dynamic idle updated";
    });
    await expect(successButton).toHaveText("Dynamic idle updated");
    await expect(idleLabelButton).toHaveText("Dynamic idle updated");

    await successButton.click();
    await expect(successButton).toHaveText("Dynamic pending");
    await expect(successButton).toHaveText("Dynamic succeeded", {
      timeout: 4000,
    });

    await failureButton.click();
    await expect(failureButton).toHaveText("Dynamic pending");
    await expect(failureButton).toHaveText("Dynamic failed", {
      timeout: 4000,
    });
  });

  test("@browser @control a FormId form mounted after init attaches the submit tracking", async ({
    page,
  }) => {
    const trackingResult = await page.evaluate(async () => {
      const container = document.createElement("div");
      container.innerHTML = `
        <div x-data='buttonState({"idle":"","pending":"ph-spinner","succeeded":"ph-check-circle","failed":"ph-warning-circle"}, true, 0)'>
          <button type="submit" form="lazy-tracking-form" id="lazy-tracking-button">Save</button>
        </div>
      `;
      document.body.appendChild(container);
      window.Alpine.initTree(container);
      const state = window.Alpine.$data(container.firstElementChild);

      const form = document.createElement("form");
      form.id = "lazy-tracking-form";
      form.innerHTML = '<input name="value" value="x">';
      document.body.appendChild(form);

      const deadline = performance.now() + 1000;
      while (state.submittedForm === null && performance.now() < deadline) {
        await new Promise((resolve) => setTimeout(resolve, 10));
      }
      if (state.submittedForm === null) {
        return {
          submittedFormId: "",
          isPending: false,
          isSecondSubmitBlocked: false,
        };
      }

      form.dispatchEvent(
        new Event("submit", { bubbles: true, cancelable: true }),
      );
      await new Promise((resolve) => setTimeout(resolve, 10));
      const isPending = state.actionStatus === "pending";

      const secondSubmitEvent = new Event("submit", {
        bubbles: true,
        cancelable: true,
      });
      form.dispatchEvent(secondSubmitEvent);

      return {
        submittedFormId: state.submittedForm.id,
        isPending,
        isSecondSubmitBlocked: secondSubmitEvent.defaultPrevented,
      };
    });

    expect(trackingResult).toEqual({
      submittedFormId: "lazy-tracking-form",
      isPending: true,
      isSecondSubmitBlocked: true,
    });
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

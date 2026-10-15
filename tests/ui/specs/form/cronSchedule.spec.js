import { expect, test } from "../../fixture.js";

const cronSection = "#cron-schedule-demo";
const boundValue = `${cronSection} p span`;
const previewText = `${cronSection} span[x-text="cronPreviewTextBuilder()"]`;
const expressionChip = `${cronSection} fieldset code[x-text]`;
const hiddenInput = `${cronSection} input[type=hidden][name=schedule]`;

function presetTrigger(page) {
  return page.locator(
    `${cronSection} div[x-show="cronScheduleMode === 'predefined'"] [x-ref=selectTrigger]`,
  );
}

function customModeLabel(page) {
  return page
    .locator(`${cronSection} label`)
    .filter({ has: page.locator('input[type=radio][value="custom"]') });
}

function fieldInput(page, fieldLabel) {
  return page.locator(`${cronSection} input[aria-label=${fieldLabel}]`);
}

test.describe("CronSchedule", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${cronSection.slice(1)}`);
    await expect(page.locator(boundValue)).toHaveText("0 0 * * *");
    await expect(page.locator(previewText)).toHaveText("At 00:00, every day");
  });

  test("@smoke picking a preset writes the expression and the preview", async ({
    page,
  }) => {
    await presetTrigger(page).click();
    await page
      .locator(`${cronSection} ul`)
      .getByText("Every 5 days", { exact: true })
      .click();

    await expect(page.locator(boundValue)).toHaveText("0 0 */5 * *");
    await expect(page.locator(hiddenInput)).toHaveValue("0 0 */5 * *");
    await expect(page.locator(previewText)).toHaveText(
      "At 00:00, every 5 days",
    );
  });

  test("@smoke custom mode composes the five fields into the bound path", async ({
    page,
  }) => {
    await customModeLabel(page).click();
    await fieldInput(page, "Minute").fill("30");
    await fieldInput(page, "Hour").fill("9");
    await fieldInput(page, "Weekday").fill("1-5");

    await expect(page.locator(boundValue)).toHaveText("30 9 * * 1-5");
    await expect(page.locator(previewText)).toHaveText("At 09:30, on weekdays");
  });

  test("@smoke an invalid field clears the bound path and reports it", async ({
    page,
  }) => {
    await customModeLabel(page).click();
    await fieldInput(page, "Minute").fill("99");

    await expect(page.locator(previewText)).toHaveText(
      "Invalid schedule: fix the highlighted fields",
    );
    await expect(page.locator(boundValue)).toHaveText("");

    await fieldInput(page, "Minute").fill("15");
    await expect(page.locator(boundValue)).toHaveText("15 0 * * *");
    await expect(page.locator(previewText)).toHaveText("At 00:15, every day");
  });

  test("@smoke a descending weekday range is rejected and weekends is valid", async ({
    page,
  }) => {
    await customModeLabel(page).click();
    await fieldInput(page, "Weekday").fill("6-0");

    await expect(page.locator(previewText)).toHaveText(
      "Invalid schedule: fix the highlighted fields",
    );
    await expect(page.locator(boundValue)).toHaveText("");

    await fieldInput(page, "Weekday").fill("6,0");
    await expect(page.locator(boundValue)).toHaveText("0 0 * * 6,0");
    await expect(page.locator(previewText)).toHaveText("At 00:00, on weekends");
  });

  test("@browser an invalid field marks the border red", async ({ page }) => {
    await customModeLabel(page).click();
    await fieldInput(page, "Minute").fill("99");

    await expect(page.locator(expressionChip)).toHaveText("—");
    await expect
      .poll(async () =>
        fieldInput(page, "Minute").evaluate(
          (input) => getComputedStyle(input.closest("fieldset")).borderColor,
        ),
      )
      .toContain("239");
  });
});

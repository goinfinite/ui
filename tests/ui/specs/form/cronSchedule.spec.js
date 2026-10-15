import { expect, test } from "../../fixture.js";

const cronSection = "#cron-schedule-demo";
const boundValue = `${cronSection} p span`;
const previewText = `${cronSection} .ph-clock + span`;
const expressionChip = `${cronSection} fieldset code`;
const hiddenInput = `${cronSection} input[type=hidden][name=schedule]`;

function presetTrigger(page) {
  return page.locator(`${cronSection} [role=button]`).first();
}

function customModeLabel(page) {
  return page
    .locator(`${cronSection} label`)
    .filter({ has: page.locator('input[type=radio][value="custom"]') });
}

function fieldInput(page, fieldLabel) {
  return page.locator(`${cronSection} input[aria-label=${fieldLabel}]`);
}

async function expectInvalidFieldThenRecovers(
  page,
  fieldLabel,
  invalidValue,
  validValue,
  expectedExpression,
  expectedPreview,
) {
  await fieldInput(page, fieldLabel).fill(invalidValue);

  await expect(page.locator(previewText)).toHaveText(
    "Invalid schedule: fix the highlighted fields",
  );
  await expect(page.locator(boundValue)).toHaveText("");

  await fieldInput(page, fieldLabel).fill(validValue);
  await expect(page.locator(boundValue)).toHaveText(expectedExpression);
  await expect(page.locator(previewText)).toHaveText(expectedPreview);
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

  test("@browser arrow keys move through the unnamed preset radios and write the expression", async ({
    page,
  }) => {
    const trigger = presetTrigger(page);
    const options = page.locator(`${cronSection} ul`).first();

    await trigger.focus();
    await page.keyboard.press("ArrowDown");
    await expect(options).toBeVisible();

    await page.keyboard.press("ArrowDown");
    await expect(options.locator("input:checked")).toHaveValue("0 1 * * *");
    await expect(page.locator(boundValue)).toHaveText("0 1 * * *");

    await page.keyboard.press("Enter");
    await expect(options).toBeHidden();
    await expect(trigger).toBeFocused();
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
    await expectInvalidFieldThenRecovers(
      page,
      "Minute",
      "99",
      "15",
      "15 0 * * *",
      "At 00:15, every day",
    );
  });

  test("@smoke a descending weekday range is rejected and weekends is valid", async ({
    page,
  }) => {
    await customModeLabel(page).click();
    await expectInvalidFieldThenRecovers(
      page,
      "Weekday",
      "6-0",
      "6,0",
      "0 0 * * 6,0",
      "At 00:00, on weekends",
    );
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

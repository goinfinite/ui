import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const pickerSection = "#date-time-picker-demo";
const mainPickerSelector = `${pickerSection} [x-data*="bookingDate"]`;
const triggerSelector = 'button[aria-label="Toggle date picker"]';
const popoverSelector = 'div[x-show="isOpen"]';
const monthLabelSelector = '[x-text="monthYearLabel"]';

function pickDay(picker, dayNumber) {
  return picker
    .locator("button[data-day-index]")
    .filter({ hasText: new RegExp(`^${dayNumber}$`) })
    .first();
}

function selectByLegend(picker, legendText) {
  return picker
    .locator('[x-data^="selectInput"]')
    .filter({ has: picker.page().locator("legend", { hasText: legendText }) })
    .first();
}

test.describe("DateTimePicker", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${pickerSection.slice(1)}`);
    await expect(
      page.locator(`${mainPickerSelector} input[type=text]`),
    ).toHaveValue("10/16/2026");
  });

  test("@smoke picking a day writes the ISO value to its state path", async ({
    page,
  }) => {
    const picker = page.locator(mainPickerSelector).first();

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();
    await expect(picker.locator(monthLabelSelector)).toHaveText("October 2026");

    await pickDay(picker, 20).click();

    await expect(picker.locator("input[type=text]")).toHaveValue("10/20/2026");
    await expect(picker.locator("p span").first()).toHaveText("2026-10-20");
  });

  test("a typed ISO value commits and a bad value reverts", async ({
    page,
  }) => {
    const picker = page.locator(mainPickerSelector).first();
    const input = picker.locator("input[type=text]");

    await input.fill("2027-01-05");
    await input.press("Enter");

    await expect(input).toHaveValue("01/05/2027");

    await input.fill("nonsense");
    await input.press("Enter");

    await expect(input).toHaveValue("01/05/2027");

    await input.fill("2026-02-31");
    await input.press("Enter");

    await expect(input).toHaveValue("01/05/2027");
  });

  test("Escape closes the calendar", async ({ page }) => {
    const picker = page.locator(mainPickerSelector).first();

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    await page.evaluate(() =>
      window.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" })),
    );

    await expect(picker.locator(popoverSelector)).toBeHidden();
  });

  test("@smoke the range picker fills start then end and writes both paths", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, pickerSection, "Range");
    const picker = panel.locator('[x-data*="bookingStart"]').first();
    const readouts = picker.locator("p span");

    await expect(readouts.nth(0)).toHaveText("2026-10-16");
    await expect(readouts.nth(1)).toHaveText("2026-10-21");

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    await pickDay(picker, 5).click();
    await expect(readouts.nth(0)).toHaveText("2026-10-05");

    await pickDay(picker, 25).click();
    await expect(readouts.nth(1)).toHaveText("2026-10-25");
  });

  test("each range input commits its own side", async ({ page }) => {
    const panel = await openExamplePanel(page, pickerSection, "Range");
    const picker = panel.locator('[x-data*="bookingStart"]').first();
    const readouts = picker.locator("p span");

    const startInput = picker.locator(
      'input[aria-label="Booking window start"]',
    );
    const endInput = picker.locator('input[aria-label="Booking window end"]');

    await startInput.fill("2026-10-05");
    await startInput.press("Enter");
    await expect(readouts.nth(0)).toHaveText("2026-10-05");

    await endInput.fill("2026-11-09");
    await endInput.press("Enter");
    await expect(readouts.nth(1)).toHaveText("2026-11-09");
  });

  test("out-of-bounds days render disabled", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      pickerSection,
      "Locale and Bounds",
    );
    const picker = panel.locator("[x-data^=dateTimePicker]").first();

    await expect(
      picker.locator("button[data-day-index][disabled]").first(),
    ).toBeDisabled();
  });

  test("@smoke the time selector commits hours and minutes through SelectInput", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      pickerSection,
      "Date, Time, and Timezone",
    );
    const picker = panel.locator('[x-data*="publishAt"]').first();
    const readout = picker.locator("p span").first();

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    const hourSelect = selectByLegend(picker, "Hour");
    await hourSelect.locator('[role="button"]').click();
    await hourSelect
      .locator("ul label")
      .getByText("14", { exact: true })
      .click();

    await expect(readout).toHaveText(/T14:\d{2}-03:00$/);

    const minuteSelect = selectByLegend(picker, "Minute");
    await minuteSelect.locator('[role="button"]').click();
    await minuteSelect
      .locator("ul label")
      .getByText("45", { exact: true })
      .click();

    await expect(readout).toHaveText(/T14:45-03:00$/);
  });

  test("@browser the timezone select rewrites the value offset and its state path", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      pickerSection,
      "Date, Time, and Timezone",
    );
    const picker = panel.locator('[x-data*="publishAt"]').first();
    const readouts = picker.locator("p span");

    await expect(readouts.nth(0)).toHaveText("2026-10-16T09:30-03:00");
    await expect(readouts.nth(1)).toHaveText("America/Sao_Paulo");

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    const timezoneSelect = selectByLegend(picker, "Timezone");
    await expect(timezoneSelect.locator('[role="button"]')).toContainText(
      "America/Sao Paulo",
    );
    await timezoneSelect.locator('[role="button"]').click();
    await timezoneSelect
      .locator("ul label")
      .getByText("UTC", { exact: true })
      .click();

    await expect(readouts.nth(0)).toHaveText("2026-10-16T09:30+00:00");
    await expect(readouts.nth(1)).toHaveText("UTC");
  });

  test("@browser the timezone offset follows the selected wall-clock time across a DST edge", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      pickerSection,
      "Date, Time, and Timezone",
    );
    const picker = panel.locator('[x-data*="publishAt"]').first();
    const readouts = picker.locator("p span");
    const input = picker.locator("input[type=text]");

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    const timezoneSelect = selectByLegend(picker, "Timezone");
    await timezoneSelect.locator('[role="button"]').click();
    await timezoneSelect
      .locator("ul label")
      .getByText("America/New York", { exact: true })
      .click();

    await page.evaluate(() =>
      window.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" })),
    );

    await input.fill("2026-03-08T09:30");
    await input.press("Enter");

    await expect(readouts.nth(0)).toHaveText("2026-03-08T09:30-04:00");
    await expect(readouts.nth(1)).toHaveText("America/New_York");
  });

  test("the range field captions each side and draws the calendar band", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, pickerSection, "Range");
    const picker = panel.locator('[x-data*="bookingStart"]').first();

    await expect(
      picker
        .locator("span")
        .filter({ hasText: /^Start$/ })
        .first(),
    ).toBeVisible();
    await expect(
      picker.locator("span").filter({ hasText: /^End$/ }).first(),
    ).toBeVisible();
    await expect(picker.locator("span").filter({ hasText: "→" })).toHaveCount(
      0,
    );

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    await expect(pickDay(picker, 18)).toHaveClass(/bg-secondary-500\/20/);
    await expect(pickDay(picker, 16)).toHaveClass(/rounded-l-md/);
    await expect(pickDay(picker, 21)).toHaveClass(/rounded-r-md/);
  });

  test("TextCase and the Done button settings reach the rendered control", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, pickerSection, "Appearance");
    const picker = panel.locator('[x-data*="stayStart"]').first();

    await expect(picker.locator("legend").first()).toHaveClass(/uppercase/);
    await expect(
      picker
        .locator("span")
        .filter({ hasText: /^Check-in$/ })
        .first(),
    ).toHaveClass(/uppercase/);

    await picker.locator(triggerSelector).click();
    await expect(picker.locator(popoverSelector)).toBeVisible();

    const doneButton = picker.getByRole("button", { name: "Apply dates" });
    await expect(doneButton).toBeVisible();
    await expect(doneButton).toHaveClass(/bg-primary-500/);
  });
});

import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const fieldsetSection = "#multi-column-repeatable-fieldset-demo";
const mainWrapperSelector = `${fieldsetSection} [x-data*=teamMemberRows]`;
const readoutSelector = "p span";

const anaRow = '{"name":"Ana Souza","role":"owner","age":"38","active":true}';

test.describe("MultiColumnRepeatableFieldset", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${fieldsetSection.slice(1)}`);
    await expect(
      page.locator(`${mainWrapperSelector} ${readoutSelector}`),
    ).toHaveText(`[${anaRow}]`);
  });

  test("@smoke add, duplicate, move, and remove change the bound rows", async ({
    page,
  }) => {
    const wrapper = page.locator(mainWrapperSelector);
    const readout = wrapper.locator(readoutSelector);
    const rows = wrapper.locator("tbody tr:has(input)");

    await wrapper.getByText("Add row", { exact: true }).click();
    const addedRows = JSON.parse(await readout.textContent());
    expect(addedRows).toHaveLength(2);
    expect(addedRows[0]).toEqual(JSON.parse(anaRow));
    expect(addedRows[1].name).toMatch(/^Member \d+$/);
    const randomRow = JSON.stringify(addedRows[1]);

    await rows.first().getByLabel("Duplicate row").click();
    await expect(readout).toHaveText(`[${anaRow},${anaRow},${randomRow}]`);

    await rows.nth(2).getByLabel("Move row up").click();
    await expect(readout).toHaveText(`[${anaRow},${randomRow},${anaRow}]`);

    await rows.first().getByLabel("Remove row").click();
    await expect(readout).toHaveText(`[${randomRow},${anaRow}]`);
  });

  test("move buttons disable at the list edges", async ({ page }) => {
    const wrapper = page.locator(mainWrapperSelector);
    const rows = wrapper.locator("tbody tr:has(input)");

    await expect(rows.first().getByLabel("Move row up")).toBeDisabled();
    await expect(rows.last().getByLabel("Move row down")).toBeDisabled();
  });

  test("the hidden input carries the rows as JSON for form submission", async ({
    page,
  }) => {
    const hiddenInput = page.locator(
      `${mainWrapperSelector} input[name=teamMembers]`,
    );

    await expect(hiddenInput).toHaveValue(`[${anaRow}]`);
  });

  test("a required cell without a value is marked invalid", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      fieldsetSection,
      "Per-Row Validation",
    );

    await expect(panel.locator("[aria-invalid=true]").first()).toHaveAttribute(
      "aria-invalid",
      "true",
    );
  });

  test("the empty state shows its label and adds the first row", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Empty State");
    const emptyCell = panel.locator("td", { hasText: "No port mappings yet" });

    await expect(emptyCell).toBeVisible();

    await panel.locator("button", { hasText: "Add port mapping" }).click();

    await expect(emptyCell).toBeHidden();
    await expect(panel.locator("input[type=number]")).toHaveCount(1);
  });
});

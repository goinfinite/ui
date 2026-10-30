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

    await wrapper.getByText("add row", { exact: true }).click();
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

  test("@smoke bulk edit converts the rows to text and back", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=envVarRows]");
    const readout = wrapper.locator('p:has-text("envVarRows:") span');

    await expect(readout).toHaveText(
      JSON.stringify([
        { name: "APP_ENV", value: "production" },
        { name: "LOG_LEVEL", value: "info" },
      ]),
    );

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toBeVisible();
    await expect(textarea).toHaveValue("APP_ENV=production\nLOG_LEVEL=info");

    await textarea.fill(
      "NEW_VAR=first\nDATABASE_URL=postgres://db.example/app?a=b\n\nEMPTY_VALUE=",
    );
    await expect(readout).toHaveText(
      JSON.stringify([
        { name: "NEW_VAR", value: "first" },
        { name: "DATABASE_URL", value: "postgres://db.example/app?a=b" },
        { name: "EMPTY_VALUE", value: "" },
      ]),
    );

    await wrapper.getByRole("button", { name: "table view" }).click();
    await expect(textarea).toBeHidden();
    const rows = wrapper.locator("tbody tr:has(input)");
    await expect(rows).toHaveCount(3);
    await expect(rows.nth(1).locator('input[type="text"]').nth(1)).toHaveValue(
      "postgres://db.example/app?a=b",
    );
  });

  test("@smoke toggling bulk edit without an edit keeps the rows unchanged", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=envVarRows]");
    const readout = wrapper.locator('p:has-text("envVarRows:") span');
    const rows = wrapper.locator("tbody tr:has(input)");
    const nameInput = rows.first().locator('input[type="text"]').first();

    await nameInput.fill("  Ana  ");
    await wrapper.getByRole("button", { name: "add row" }).click();

    const expectedRows = [
      { name: "  Ana  ", value: "production" },
      { name: "LOG_LEVEL", value: "info" },
      { name: "", value: "" },
    ];
    await expect(readout).toHaveText(JSON.stringify(expectedRows));

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    await wrapper.getByRole("button", { name: "table view" }).click();

    await expect(readout).toHaveText(JSON.stringify(expectedRows));
  });

  test("@smoke the bulk textarea clear button empties the bound rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=envVarRows]");
    const readout = wrapper.locator('p:has-text("envVarRows:") span');

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toBeVisible();

    await textarea.focus();
    await wrapper
      .getByRole("button", { name: "Clear text area content" })
      .click();

    await expect(readout).toHaveText("[]");
    await expect(textarea).toHaveValue("");
  });

  test("@smoke bulk edit parses typed cells with the default delimiter", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=typedBulkRows]");
    const readout = wrapper.locator('p:has-text("typedBulkRows:") span');
    const hiddenInput = wrapper.locator('input[name="typedBulkRows"]');

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("Ana|owner|true");

    const enabledTokens = [
      ["true", true],
      ["1", true],
      ["on", true],
      ["yes", true],
      ["false", false],
      ["0", false],
    ];
    for (const [token, expectedEnabled] of enabledTokens) {
      await textarea.fill(`Test|member|${token}`);
      await expect(readout).toHaveText(
        JSON.stringify([
          { name: "Test", role: "member", enabled: expectedEnabled },
        ]),
      );
    }

    await textarea.fill("Bob|member|yes\nCara|owner|0");
    const expectedRows = [
      { name: "Bob", role: "member", enabled: true },
      { name: "Cara", role: "owner", enabled: false },
    ];
    await expect(readout).toHaveText(JSON.stringify(expectedRows));
    await expect(hiddenInput).toHaveValue(JSON.stringify(expectedRows));

    await wrapper.getByRole("button", { name: "table view" }).click();
    const rows = wrapper.locator("tbody tr:has(input)");
    await expect(rows).toHaveCount(2);
    await expect(rows.nth(0).locator('input[type="checkbox"]')).toBeChecked();
    await expect(
      rows.nth(1).locator('input[type="checkbox"]'),
    ).not.toBeChecked();
    await expect(rows.nth(0).locator('input[type="text"]')).toHaveValue("Bob");
  });

  test("@smoke bulk edit escapes delimiters and backslashes in a non-final cell", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=typedBulkRows]");
    const readout = wrapper.locator('p:has-text("typedBulkRows:") span');
    const nameInput = wrapper
      .locator("tbody tr:has(input)")
      .first()
      .locator('input[type="text"]')
      .first();

    await nameInput.fill("A\\B|C");
    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("A\\\\B\\|C|owner|true");

    await textarea.fill("A\\\\B\\|C|member|false");
    await expect(readout).toHaveText(
      JSON.stringify([{ name: "A\\B|C", role: "member", enabled: false }]),
    );

    await wrapper.getByRole("button", { name: "table view" }).click();
    await expect(
      wrapper
        .locator("tbody tr:has(input)")
        .first()
        .locator('input[type="text"]')
        .first(),
    ).toHaveValue("A\\B|C");
  });

  test("@smoke bulk edit escapes newlines inside a cell", async ({ page }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=typedBulkRows]");
    const readout = wrapper.locator('p:has-text("typedBulkRows:") span');

    await wrapper.evaluate((wrapperElement) => {
      const wrapperState = window.Alpine.$data(wrapperElement);
      wrapperState.typedBulkRows = [
        { name: "Line\nBreak", role: "owner", enabled: true },
      ];
    });
    await expect(readout).toHaveText(
      JSON.stringify([{ name: "Line\nBreak", role: "owner", enabled: true }]),
    );

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("Line\\nBreak|owner|true");

    await textarea.fill("Line\\nBreak|member|true");
    await expect(readout).toHaveText(
      JSON.stringify([{ name: "Line\nBreak", role: "member", enabled: true }]),
    );

    await wrapper.getByRole("button", { name: "table view" }).click();
    await expect(wrapper.locator("tbody tr:has(input)")).toHaveCount(1);
  });

  test("@smoke bulk edit preserves properties outside the configured columns", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=envVarRows]");
    const readout = wrapper.locator('p:has-text("envVarRows:") span');
    const initialRows = [
      { id: "env-1", name: "APP_ENV", value: "production" },
      { id: "env-2", name: "LOG_LEVEL", value: "info" },
    ];

    await wrapper.evaluate((wrapperElement, rows) => {
      const wrapperState = window.Alpine.$data(wrapperElement);
      wrapperState.envVarRows = rows;
    }, initialRows);
    await expect(readout).toHaveText(JSON.stringify(initialRows));

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await textarea.fill("APP_ENV=staging\nLOG_LEVEL=debug");
    await expect(readout).toHaveText(
      JSON.stringify([
        { id: "env-1", name: "APP_ENV", value: "staging" },
        { id: "env-2", name: "LOG_LEVEL", value: "debug" },
      ]),
    );

    await textarea.fill("LOG_LEVEL=debug");
    await expect(readout).toHaveText(
      JSON.stringify([{ id: "env-2", name: "LOG_LEVEL", value: "debug" }]),
    );
  });

  test("@smoke bulk edit keeps row properties when lines move", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=envVarRows]");
    const readout = wrapper.locator('p:has-text("envVarRows:") span');
    const initialRows = [
      { id: "env-1", name: "APP_ENV", value: "production" },
      { id: "env-2", name: "LOG_LEVEL", value: "info" },
    ];

    await wrapper.evaluate((wrapperElement, rows) => {
      const wrapperState = window.Alpine.$data(wrapperElement);
      wrapperState.envVarRows = rows;
    }, initialRows);
    await expect(readout).toHaveText(JSON.stringify(initialRows));

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await textarea.fill("LOG_LEVEL=debug\nAPP_ENV=production");
    await expect(readout).toHaveText(
      JSON.stringify([
        { id: "env-2", name: "LOG_LEVEL", value: "debug" },
        { id: "env-1", name: "APP_ENV", value: "production" },
      ]),
    );
  });

  test("@smoke bulk edit keeps unchanged cell types and whitespace", async ({
    page,
  }) => {
    const wrapper = page.locator(mainWrapperSelector);
    const readout = wrapper.locator(readoutSelector);
    const initialRows = [
      { name: "  Ana  ", role: "owner", age: 38, active: true },
    ];

    await wrapper.evaluate((wrapperElement, rows) => {
      const wrapperState = window.Alpine.$data(wrapperElement);
      wrapperState.teamMemberRows = rows;
    }, initialRows);
    await expect(readout).toHaveText(JSON.stringify(initialRows));

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("  Ana  |owner|38|true");

    await textarea.fill("  Ana  |member|38|true");
    await expect(readout).toHaveText(
      JSON.stringify([
        { name: "  Ana  ", role: "member", age: 38, active: true },
      ]),
    );
  });

  test("@smoke bulk edit defaults to | for four columns", async ({ page }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=portBindingRows]");
    const readout = wrapper.locator('p:has-text("portBindingRows:") span');
    const hiddenInput = wrapper.locator('input[name="portBindingRows"]');

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("web|8080|80|tcp");

    await textarea.fill("web|8080|80|tcp\ncustom|9090|90|udp");
    const expectedRows = [
      {
        serviceName: "web",
        publicPort: "8080",
        containerPort: "80",
        protocol: "tcp",
      },
      {
        serviceName: "custom",
        publicPort: "9090",
        containerPort: "90",
        protocol: "udp",
      },
    ];
    await expect(readout).toHaveText(JSON.stringify(expectedRows));
    await expect(hiddenInput).toHaveValue(JSON.stringify(expectedRows));

    await wrapper.getByRole("button", { name: "table view" }).click();
    const rows = wrapper.locator("tbody tr:has(input)");
    await expect(rows).toHaveCount(2);
    await expect(rows.nth(1).locator('input[type="text"]').nth(0)).toHaveValue(
      "9090",
    );
    await expect(rows.nth(1).locator('input[type="text"]').nth(1)).toHaveValue(
      "90",
    );
  });

  test("@smoke bulk edit honors an explicit custom delimiter", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Delimiters");
    const wrapper = panel.locator("[x-data*=customDelimiterRows]");
    const readout = wrapper.locator('p:has-text("customDelimiterRows:") span');
    const hiddenInput = wrapper.locator('input[name="customDelimiterRows"]');

    await wrapper.getByRole("button", { name: "bulk edit" }).click();
    const textarea = wrapper.locator("textarea");
    await expect(textarea).toHaveValue("Ana;owner;true");

    await textarea.fill("Bob;member;yes\nCara;owner;0");
    const expectedRows = [
      { name: "Bob", role: "member", enabled: true },
      { name: "Cara", role: "owner", enabled: false },
    ];
    await expect(readout).toHaveText(JSON.stringify(expectedRows));
    await expect(hiddenInput).toHaveValue(JSON.stringify(expectedRows));
  });

  test("@smoke button styles customize the add row and bulk edit controls", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, fieldsetSection, "Appearance");
    const wrapper = panel.locator("[x-data*=endpointRows]");
    const addRowButton = wrapper.getByRole("button", { name: "Add endpoint" });
    const bulkEditButton = wrapper.getByRole("button", {
      name: "Paste endpoints",
    });

    await expect(addRowButton).toHaveClass(/rounded-full/);
    await expect(addRowButton).toHaveClass(/bg-emerald-500/);
    await expect(addRowButton.locator("i")).toHaveClass(/ph-plus-circle/);
    await expect(bulkEditButton).toHaveClass(/rounded-none/);
    await expect(bulkEditButton).toHaveClass(/text-emerald-400/);
    await expect(bulkEditButton.locator("i")).toHaveClass(/ph-clipboard-text/);

    await bulkEditButton.click();
    const tableButton = wrapper.getByRole("button", {
      name: "Show endpoints",
    });
    await expect(tableButton.locator("i")).toHaveClass(/ph-table/);
  });
});

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const section = "#remote-searchable-select-input-demo";

function mainCombobox(page) {
  return page.locator(`${section} input[role=combobox]`).first();
}

function mainListbox(page) {
  return page.locator(`${section} ul[role=listbox]`).first();
}

test.describe("RemoteSearchableSelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${section.slice(1)}`);
    await expect(mainCombobox(page)).toBeVisible();
  });

  test("@smoke typing three characters fetches options and selects one", async ({
    page,
  }) => {
    const combobox = mainCombobox(page);
    await combobox.click();
    await combobox.fill("bra");

    const listbox = mainListbox(page);
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeVisible({
      timeout: 5000,
    });
    await listbox.getByRole("option", { name: "Brazil" }).click();

    await expect(combobox).toHaveValue("Brazil");
    await expect(page.locator(`${section} p.font-mono`).first()).toContainText(
      "Value: Brazil",
    );
  });

  test("@smoke initial options label a preset value before any request", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Preset Value");
    const combobox = panel.locator("input[role=combobox]");
    await expect(combobox).toHaveValue("Brazil");

    await combobox.click();
    const listbox = panel.locator("ul[role=listbox]");
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeVisible();

    await combobox.fill("chi");
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeHidden();
  });

  test("loading indicator shows while the request is in flight", async ({
    page,
  }) => {
    await mainCombobox(page).click();
    await mainCombobox(page).fill("bra");

    await expect(mainListbox(page).locator(".ph-spinner")).toBeVisible({
      timeout: 5000,
    });
  });

  test("minimum query length shows a prompt and loads no options", async ({
    page,
  }) => {
    await mainCombobox(page).click();
    await mainCombobox(page).fill("br");

    const listbox = mainListbox(page);
    await expect(listbox.getByText("Type at least 3 characters")).toBeVisible();
    await expect(listbox.getByRole("option")).toHaveCount(0);
  });

  test("a response is discarded when the query drops below the minimum", async ({
    page,
  }) => {
    await mainCombobox(page).click();
    await mainCombobox(page).fill("bra");
    await page.waitForTimeout(450);

    await mainCombobox(page).fill("b");
    await page.waitForTimeout(700);

    const listbox = mainListbox(page);
    await expect(listbox.getByText("Type at least 3 characters")).toBeVisible();
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeHidden();
  });

  test("multi-select collects options from several queries", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Multi-Select");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");
    const listbox = panel.locator("ul[role=listbox]");

    await combobox.click();
    await combobox.fill("bra");
    await listbox.getByRole("option", { name: "Brazil" }).click();
    await expect(stateText).toHaveText('["Brazil"]');

    await combobox.fill("chi");
    await listbox.getByRole("option", { name: "Chile" }).click();
    await expect(stateText).toHaveText('["Brazil","Chile"]');
  });

  test("a failing request shows the error state", async ({ page }) => {
    const panel = await openExamplePanel(page, section, "Error State");
    const combobox = panel.locator("input[role=combobox]");

    await combobox.click();
    await combobox.fill("bra");

    await expect(panel.getByText("Could not load options.")).toBeVisible({
      timeout: 5000,
    });
  });
});

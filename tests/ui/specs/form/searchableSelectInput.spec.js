import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const section = "#searchable-select-input-demo";

function mainCombobox(page) {
  return page.locator(`${section} input[role=combobox]`).first();
}

function mainListbox(page) {
  return page.locator(`${section} ul[role=listbox]`).first();
}

test.describe("SearchableSelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${section.slice(1)}`);
    await expect(mainCombobox(page)).toHaveValue("Brazil");
  });

  test("@smoke filter narrows the options and selects one", async ({
    page,
  }) => {
    const combobox = mainCombobox(page);
    await combobox.click();
    await combobox.fill("can");

    const listbox = mainListbox(page);
    await expect(listbox.getByRole("option", { name: "Canada" })).toBeVisible();
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeHidden();

    await listbox.getByRole("option", { name: "Canada" }).click();
    await expect(listbox).toBeHidden();
    await expect(combobox).toHaveValue("Canada");
    await expect(page.locator(`${section} p.font-mono`).first()).toContainText(
      "Value: Canada",
    );
  });

  test("@smoke clear button empties the selection", async ({ page }) => {
    await mainCombobox(page).click();
    await page.locator(`${section} .ph-x-circle`).first().click();
    await expect(mainCombobox(page)).toHaveValue("");
  });

  test("no matches row appears for an unmatched query", async ({ page }) => {
    const combobox = mainCombobox(page);
    await combobox.click();
    await combobox.fill("xyz");

    const listbox = mainListbox(page);
    await expect(listbox.getByText("No matches")).toBeVisible();
    await expect(listbox.getByRole("option", { name: "Brazil" })).toBeHidden();
  });

  test("Escape closes the dropdown", async ({ page }) => {
    const combobox = mainCombobox(page);
    await combobox.click();
    await expect(mainListbox(page)).toBeVisible();

    await page.keyboard.press("Escape");
    await expect(mainListbox(page)).toBeHidden();
  });

  test("single select submits one hidden entry", async ({ page }) => {
    const formEntries = await page.evaluate(() => {
      const section = document.querySelector("#searchable-select-input-demo");
      const hiddenInput = section.querySelector(
        "input[type=hidden][name=country]",
      );
      const componentRoot = hiddenInput.closest("[x-data]");
      const form = document.createElement("form");
      componentRoot.parentNode.insertBefore(form, componentRoot);
      form.appendChild(componentRoot);
      const entries = [...new FormData(form).entries()];
      form.parentNode.insertBefore(componentRoot, form);
      form.remove();
      return entries;
    });

    expect(formEntries).toEqual([["country", "Brazil"]]);
  });

  test("multi-select toggles values and clears the array", async ({ page }) => {
    const panel = await openExamplePanel(page, section, "Multi-Select");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");
    await expect(stateText).toHaveText('["Brazil","Canada"]');

    await combobox.click();
    await panel
      .locator("ul[role=listbox]")
      .getByRole("option", { name: "Argentina" })
      .click();
    await expect(stateText).toHaveText('["Brazil","Canada","Argentina"]');

    await panel.locator(".ph-x-circle").click();
    await expect(stateText).toHaveText("[]");
  });

  test("@smoke tags: option click adds a tag and the tag button removes it", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Selected as Tags");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");
    await expect(stateText).toHaveText('["Brazil"]');

    await combobox.click();
    await panel
      .locator("ul[role=listbox]")
      .getByRole("option", { name: "Argentina" })
      .click();
    await expect(stateText).toHaveText('["Brazil","Argentina"]');

    await panel.getByRole("button", { name: "Remove Argentina" }).click();
    await expect(stateText).toHaveText('["Brazil"]');
  });

  test("@smoke tags: a typed value becomes a tag on Enter", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Selected as Tags");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");

    await combobox.fill("Atlantis");
    await page.keyboard.press("Enter");
    await expect(stateText).toHaveText('["Brazil","Atlantis"]');
    await expect(
      panel.getByRole("button", { name: "Remove Atlantis" }),
    ).toBeVisible();
  });

  test("tags: Backspace on an empty input removes the last tag", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Selected as Tags");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");

    await combobox.fill("Atlantis");
    await page.keyboard.press("Enter");
    await expect(stateText).toHaveText('["Brazil","Atlantis"]');

    await combobox.fill("");
    await page.keyboard.press("Backspace");
    await expect(stateText).toHaveText('["Brazil"]');
  });

  test("tags: options-only mode ignores typed values", async ({ page }) => {
    const panel = await openExamplePanel(page, section, "Tags Options Only");
    const combobox = panel.locator("input[role=combobox]");
    const stateText = panel
      .locator("p", { hasText: "Values:" })
      .locator("span");

    await combobox.fill("Atlantis");
    await page.keyboard.press("Enter");
    await expect(stateText).toHaveText('["Brazil"]');
  });
});

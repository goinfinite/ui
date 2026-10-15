import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const radioGroupSection = "#inline-radio-group-demo";

function mainGroup(page) {
  return page.locator(`${radioGroupSection} fieldset`).first();
}

function stateReadout(page) {
  return page
    .locator(`${radioGroupSection} p`, { hasText: "Value:" })
    .locator("span");
}

test.describe("InlineRadioGroup", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${radioGroupSection.slice(1)}`);
    await expect(stateReadout(page)).toHaveText("option2");
  });

  test("@smoke selecting an option updates the bound state", async ({
    page,
  }) => {
    const group = mainGroup(page);
    await group.getByText("Option 1", { exact: true }).click();

    await expect(stateReadout(page)).toHaveText("option1");
    await expect(group.getByLabel("Option 1", { exact: true })).toBeChecked();
    await expect(
      group.getByLabel("Option 2", { exact: true }),
    ).not.toBeChecked();
  });

  test("@smoke the group submits the selected value under its input name", async ({
    page,
  }) => {
    await mainGroup(page).getByText("Option 1", { exact: true }).click();

    const formEntries = await page.evaluate(() => {
      const checkedRadio = document.querySelector(
        "#inline-radio-group-demo input[type=radio]:checked",
      );
      const groupRoot = checkedRadio.closest("[x-data]");
      const form = document.createElement("form");
      groupRoot.parentNode.insertBefore(form, groupRoot);
      form.appendChild(groupRoot);
      const entries = [...new FormData(form).entries()];
      form.parentNode.insertBefore(groupRoot, form);
      form.remove();
      return entries;
    });

    expect(formEntries).toEqual([["group1", "option1"]]);
  });

  test("vertical orientation keeps the group bound to one shared path", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      radioGroupSection,
      "Vertical Orientation",
    );
    const readout = panel.locator("p", { hasText: "Value:" }).locator("span");
    await expect(readout).toHaveText("option2");

    await panel.getByText("Option 3", { exact: true }).click();

    await expect(readout).toHaveText("option3");
  });
});

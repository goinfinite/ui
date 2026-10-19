import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const selectSection = "#select-input-demo";

function checkedRadioOf(page, name) {
  return page.locator(
    `${selectSection} input[type=radio][name=${name}]:checked`,
  );
}

test.describe("SelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${selectSection.slice(1)}`);
    await expect(checkedRadioOf(page, "country")).toHaveValue("Brazil");
  });

  test("@smoke dropdown opens, selects an option and closes", async ({
    page,
  }) => {
    const trigger = page.locator(`${selectSection} .group.flex`).first();
    await trigger.click();

    const options = page.locator(`${selectSection} ul`).first();
    await expect(options).toBeVisible();

    await options.getByText("Argentina", { exact: true }).click();
    await expect(options).toBeHidden();

    await expect(checkedRadioOf(page, "country")).toHaveValue("Argentina");
  });

  test("@smoke clear button empties the selection", async ({ page }) => {
    await page.locator(`${selectSection} .ph-x-circle`).first().click();
    await expect(
      page.locator(
        `${selectSection} input[type=radio][name=country][value=""]`,
      ),
    ).toBeChecked();
  });

  test("@browser @smoke option value with an apostrophe selects and highlights", async ({
    page,
  }) => {
    const trigger = page.locator(`${selectSection} .group.flex`).first();
    await trigger.click();

    const options = page.locator(`${selectSection} ul`).first();
    const apostropheOption = options
      .locator("li")
      .filter({ hasText: "Côte d'Ivoire" });
    await apostropheOption.click();
    await expect(options).toBeHidden();

    await expect(apostropheOption.locator("input[type=radio]")).toBeChecked();
    await expect(trigger).toContainText("Côte d'Ivoire");
    await expect(apostropheOption).toHaveCSS(
      "background-color",
      "rgba(250, 250, 250, 0.1)",
    );
  });

  test("@browser @smoke trigger opens with the keyboard and closes with Escape", async ({
    page,
  }) => {
    const trigger = page.locator(`${selectSection} [role=button]`).first();
    const options = page.locator(`${selectSection} ul`).first();

    await trigger.focus();
    await page.keyboard.press("Enter");
    await expect(options).toBeVisible();

    await page.keyboard.press("Escape");
    await expect(options).toBeHidden();
    await expect(trigger).toBeFocused();
  });

  test("@browser @smoke arrow down opens with the checked option focused and enter commits", async ({
    page,
  }) => {
    const trigger = page.locator(`${selectSection} [role=button]`).first();
    const options = page.locator(`${selectSection} ul`).first();

    await trigger.focus();
    await page.keyboard.press("ArrowDown");
    await expect(options).toBeVisible();
    await expect(checkedRadioOf(page, "country")).toHaveValue("Brazil");
    await expect(checkedRadioOf(page, "country")).toBeFocused();

    await page.keyboard.press("ArrowDown");
    await expect(checkedRadioOf(page, "country")).toHaveValue("Chile");
    await expect(trigger).toContainText("Chile");

    await page.keyboard.press("Enter");
    await expect(options).toBeHidden();
    await expect(trigger).toBeFocused();
  });

  test("@browser the Style example moves the checked row icon with the selection", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, selectSection, "Appearance");
    const trigger = panel.locator("[role=button]").first();
    const options = panel.locator("ul").first();

    await trigger.click();
    const selectedRow = options
      .locator("li")
      .filter({ has: page.locator("i.ph-check-circle") });
    await expect(selectedRow).toHaveCount(1);
    await expect(selectedRow).toContainText("Brazil");

    await options.getByText("Chile", { exact: true }).click();
    await trigger.click();
    const newSelectedRow = options
      .locator("li")
      .filter({ has: page.locator("i.ph-check-circle") });
    await expect(newSelectedRow).toHaveCount(1);
    await expect(newSelectedRow).toContainText("Chile");
  });

  test("@browser @smoke arrow keys move through the named radio group and Escape closes", async ({
    page,
  }) => {
    const trigger = page.locator(`${selectSection} [role=button]`).first();
    const options = page.locator(`${selectSection} ul`).first();
    await trigger.click();
    await expect(options).toBeVisible();

    await page
      .locator(
        `${selectSection} input[type=radio][name=country][value="Brazil"]`,
      )
      .focus();
    await page.keyboard.press("ArrowDown");

    await expect(checkedRadioOf(page, "country")).toHaveValue("Chile");
    await expect(trigger).toContainText("Chile");
    await expect(options).toBeVisible();

    await page.keyboard.press("Escape");
    await expect(options).toBeHidden();
  });

  test("@smoke form submission carries one entry for the selected value", async ({
    page,
  }) => {
    const formEntries = await page.evaluate(() => {
      const checkedRadio = document.querySelector(
        "#select-input-demo input[type=radio][name=country]:checked",
      );
      const selectRoot = checkedRadio.closest("[x-data]");
      const form = document.createElement("form");
      selectRoot.parentNode.insertBefore(form, selectRoot);
      form.appendChild(selectRoot);
      const entries = [...new FormData(form).entries()];
      form.parentNode.insertBefore(selectRoot, form);
      form.remove();
      return entries;
    });

    expect(formEntries).toEqual([["country", "Brazil"]]);
  });

  test("@smoke OnChangeFunc runs on selection and on clear", async ({
    page,
  }) => {
    await openExamplePanel(page, selectSection, "OnChangeFunc");
    const block = page
      .locator(`${selectSection} p`)
      .filter({ hasText: "Change count:" })
      .locator("xpath=..");
    const trigger = block.locator(".group.flex");
    const changeCount = block
      .locator("p", { hasText: "Change count:" })
      .locator("span");

    await trigger.click();
    await block.getByText("Argentina", { exact: true }).click();
    await expect(changeCount).toHaveText("1");

    await trigger.locator(".ph-x-circle").click();
    await expect(changeCount).toHaveText("2");
  });

  test("@browser @smoke multi-select keeps the dropdown open and joins the summary", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, selectSection, "Multi-Select");
    const trigger = panel.locator(".group.flex").first();
    const dropdown = panel.locator("ul").first();

    await trigger.click();
    await expect(dropdown).toBeVisible();

    await dropdown.getByText("Argentina", { exact: true }).click();
    await expect(dropdown).toBeVisible();
    await expect(trigger).toContainText("Brazil, Argentina");

    await dropdown.getByText("Chile", { exact: true }).click();
    await expect(trigger).toContainText("Brazil, Argentina, Chile");

    await page.keyboard.press("Escape");
    await expect(dropdown).toBeHidden();
  });

  test("@smoke multi-select clear button empties the selection", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, selectSection, "Multi-Select");
    const trigger = panel.locator(".group.flex").first();

    await trigger.locator(".ph-x-circle").click();

    await expect(panel.locator("input[type=checkbox]:checked")).toHaveCount(0);
    await expect(trigger).toContainText("Countries");
  });

  test("@smoke multi-select form submission carries one entry per selected value", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, selectSection, "Multi-Select");
    const trigger = panel.locator(".group.flex").first();
    const dropdown = panel.locator("ul").first();

    await trigger.click();
    await dropdown.getByText("Argentina", { exact: true }).click();
    await dropdown.getByText("Chile", { exact: true }).click();

    const formEntries = await page.evaluate(() => {
      const checkedCheckbox = document.querySelector(
        "#select-input-demo input[type=checkbox][name=countries]:checked",
      );
      const selectRoot = checkedCheckbox.closest("[x-data]");
      const form = document.createElement("form");
      selectRoot.parentNode.insertBefore(form, selectRoot);
      form.appendChild(selectRoot);
      const entries = [...new FormData(form).entries()];
      form.parentNode.insertBefore(selectRoot, form);
      form.remove();
      return entries;
    });

    expect(formEntries).toHaveLength(3);
    expect(formEntries.map(([, value]) => value).sort()).toEqual([
      "Argentina",
      "Brazil",
      "Chile",
    ]);
    expect(new Set(formEntries.map(([name]) => name))).toEqual(
      new Set(["countries"]),
    );
  });
});

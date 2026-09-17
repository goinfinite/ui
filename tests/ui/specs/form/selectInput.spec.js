import { test, expect } from "@playwright/test";

const selectSection = "#select-input-demo";

function checkedRadioOf(page, name) {
  return page.locator(`${selectSection} input[type=radio][name=${name}]:checked`);
}

test.describe("SelectInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${selectSection.slice(1)}`);
    await expect(checkedRadioOf(page, "country")).toHaveValue("Brazil");
  });

  test("@smoke dropdown opens, selects an option and closes", async ({ page }) => {
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
    await expect(page.locator(`${selectSection} input[type=radio][name=country][value=""]`)).toBeChecked();
  });

  test("@smoke option value with an apostrophe selects and highlights", async ({ page }) => {
    const trigger = page.locator(`${selectSection} .group.flex`).first();
    await trigger.click();

    const options = page.locator(`${selectSection} ul`).first();
    const apostropheOption = options.locator("li").filter({ hasText: "Côte d'Ivoire" });
    await apostropheOption.click();
    await expect(options).toBeHidden();

    await expect(apostropheOption.locator("input[type=radio]")).toBeChecked();
    await expect(trigger).toContainText("Côte d'Ivoire");
    await expect(apostropheOption).toHaveCSS("background-color", "rgba(250, 250, 250, 0.1)");
  });

  test("@smoke trigger opens with the keyboard and closes with Escape", async ({ page }) => {
    const trigger = page.locator(`${selectSection} [role=button]`).first();
    const options = page.locator(`${selectSection} ul`).first();

    await trigger.focus();
    await page.keyboard.press("Enter");
    await expect(options).toBeVisible();

    await page.keyboard.press("Escape");
    await expect(options).toBeHidden();
  });

  test("@smoke arrow keys move through the named radio group and Escape closes", async ({ page }) => {
    const trigger = page.locator(`${selectSection} [role=button]`).first();
    const options = page.locator(`${selectSection} ul`).first();
    await trigger.click();
    await expect(options).toBeVisible();

    await page.locator(`${selectSection} input[type=radio][name=country][value="Brazil"]`).focus();
    await page.keyboard.press("ArrowDown");

    await expect(checkedRadioOf(page, "country")).toHaveValue("Chile");
    await expect(trigger).toContainText("Chile");
    await expect(options).toBeVisible();

    await page.keyboard.press("Escape");
    await expect(options).toBeHidden();
  });

  test("@smoke form submission carries one entry for the selected value", async ({ page }) => {
    const formEntries = await page.evaluate(() => {
      const checkedRadio = document.querySelector("#select-input-demo input[type=radio][name=country]:checked");
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

  test("@smoke OnChangeFunc runs on selection and on clear", async ({ page }) => {
    const block = page.locator(`${selectSection} .grid`).filter({ hasText: "Change count:" });
    const trigger = block.locator(".group.flex");
    const changeCount = block.locator("p", { hasText: "Change count:" }).locator("span");

    await trigger.click();
    await block.getByText("Argentina", { exact: true }).click();
    await expect(changeCount).toHaveText("1");

    await trigger.locator(".ph-x-circle").click();
    await expect(changeCount).toHaveText("2");
  });
});

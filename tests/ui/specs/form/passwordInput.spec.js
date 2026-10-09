import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const section = "#password-input-demo";

function mainPasswordInput(page) {
  return page
    .locator(`${section} input[type=password], ${section} input[type=text]`)
    .first();
}

test.describe("PasswordInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${section.slice(1)}`);
    await expect(mainPasswordInput(page)).toBeVisible();
  });

  test("@smoke reveal toggle switches the input type", async ({ page }) => {
    await expect(mainPasswordInput(page)).toHaveAttribute("type", "password");

    await page.getByRole("button", { name: "Show password" }).first().click();
    await expect(mainPasswordInput(page)).toHaveAttribute("type", "text");

    await page.getByRole("button", { name: "Hide password" }).first().click();
    await expect(mainPasswordInput(page)).toHaveAttribute("type", "password");
  });

  test("@smoke action buttons expose a tooltip", async ({ page }) => {
    const revealButton = page
      .getByRole("button", { name: "Show password" })
      .first();
    const generateButton = page
      .getByRole("button", { name: "Generate random password" })
      .first();

    const revealTooltip = page.locator(
      `#${await revealButton.getAttribute("aria-describedby")}`,
    );
    const generateTooltip = page.locator(
      `#${await generateButton.getAttribute("aria-describedby")}`,
    );

    await expect(revealTooltip).toHaveAttribute("role", "tooltip");
    await expect(revealTooltip).toHaveText("show/hide the password");
    await expect(generateTooltip).toHaveText("generate random password");
  });

  test("@smoke generate fills the field and completes the meter", async ({
    page,
  }) => {
    await page
      .getByRole("button", { name: "Generate random password" })
      .first()
      .click();

    const generatedValue = await mainPasswordInput(page).inputValue();
    expect(generatedValue.length).toBe(16);
    await expect(page.locator(`${section} p.font-mono`).first()).toContainText(
      generatedValue,
    );

    const meterWidth = await page
      .locator(`${section} .bg-secondary-500`)
      .first()
      .evaluate((element) => element.style.width);
    expect(meterWidth).toBe("100%");
  });

  test("typing updates the criteria checklist and the percentage", async ({
    page,
  }) => {
    await mainPasswordInput(page).fill("abc");

    const criteria = page.locator(`${section} ul li`);
    await expect(
      criteria.filter({ hasText: "lowercase" }).locator(".ph-check").first(),
    ).toBeVisible();
    await expect(
      criteria.filter({ hasText: "number" }).locator(".ph-check").first(),
    ).toBeHidden();

    const meterWidth = await page
      .locator(`${section} .bg-secondary-500`)
      .first()
      .evaluate((element) => element.style.width);
    expect(meterWidth).toBe("20%");
  });

  test("custom rules drive the generator and the criteria list", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Custom Rules");
    const generateButton = panel.getByRole("button", {
      name: "Generate random password",
    });
    await generateButton.click();

    const generatedValue = await panel.locator("input").inputValue();
    expect(generatedValue.length).toBe(16);
    expect(generatedValue).toMatch(/[0-9]/);
    expect(generatedValue).toMatch(/[a-z]/);
    expect(generatedValue).toMatch(/[A-Z]/);

    await expect(
      panel.locator("li").filter({ hasText: "special character" }),
    ).toHaveCount(0);

    const meterWidth = await panel
      .locator(".bg-secondary-500")
      .evaluate((element) => element.style.width);
    expect(meterWidth).toBe("100%");
  });
});

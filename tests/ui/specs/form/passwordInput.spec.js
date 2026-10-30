import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

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

  test("caps lock shows a hint while typing and clears on focusout", async ({
    page,
  }) => {
    const capsLockHint = page.locator(`${section} p[role=status]`).first();
    await expect(capsLockHint).toBeHidden();

    await mainPasswordInput(page).evaluate((inputElement) => {
      const keyboardEvent = new KeyboardEvent("keydown", {
        key: "a",
        bubbles: true,
      });
      Object.defineProperty(keyboardEvent, "getModifierState", {
        value: () => true,
      });
      inputElement.dispatchEvent(keyboardEvent);
    });
    await expect(capsLockHint).toBeVisible();

    await mainPasswordInput(page).evaluate((inputElement) => {
      inputElement.dispatchEvent(new FocusEvent("focusout", { bubbles: true }));
    });
    await expect(capsLockHint).toBeHidden();

    const generateButton = page
      .getByRole("button", { name: "Generate random password" })
      .first();
    await generateButton.focus();
    await generateButton.evaluate((buttonElement) => {
      const keyboardEvent = new KeyboardEvent("keydown", {
        key: "a",
        bubbles: true,
      });
      Object.defineProperty(keyboardEvent, "getModifierState", {
        value: () => true,
      });
      buttonElement.dispatchEvent(keyboardEvent);
    });
    await expect(capsLockHint).toBeHidden();
  });

  test("@browser caps lock warning accepts custom text, color, case, and hiding", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, section, "Hints");
    const customField = panel.locator('input[name="capsLockPassword"]');
    const customHint = panel.locator(
      'div[x-data^="passwordInput"]:has(input[name="capsLockPassword"]) p[role=status]',
    );
    const hiddenHint = panel.locator(
      'div[x-data^="passwordInput"]:has(input[name="hiddenCapsLockPassword"]) p[role=status]',
    );

    await customField.evaluate((inputElement) => {
      const keyboardEvent = new KeyboardEvent("keydown", {
        key: "a",
        bubbles: true,
      });
      Object.defineProperty(keyboardEvent, "getModifierState", {
        value: () => true,
      });
      inputElement.dispatchEvent(keyboardEvent);
    });

    await expect(customHint).toBeVisible();
    await expect(customHint).toHaveText("Caps Lock is active");
    await expect(customHint).toHaveClass(/text-red-400/);
    await expect(customHint).toHaveCSS("text-transform", "uppercase");
    await expect(hiddenHint).toHaveCount(0);
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

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const confirmationSection = "#confirmation-dialog-demo";

function backdropWith(page, text) {
  return page.locator("div.fixed.inset-0.z-100", { hasText: text });
}

test.describe("ConfirmationDialog", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${confirmationSection.slice(1)}`);
  });

  test("@smoke delete dialog keeps confirm disabled when the target name is empty", async ({
    page,
  }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    const modal = backdropWith(page, "Delete");
    await expect(modal).toBeVisible();

    const confirmButton = modal.getByRole("button", { name: "Yes, Delete!" });
    await modal
      .getByRole("textbox", { name: "Type the name to confirm" })
      .fill("anything");
    await expect(confirmButton).toBeDisabled();

    await modal.evaluate((element) => {
      Alpine.$data(element).demoTargetName = "";
    });
    await modal
      .getByRole("textbox", { name: "Type the name to confirm" })
      .fill("");
    await expect(confirmButton).toBeDisabled();
  });

  test("@smoke delete dialog gates confirm behind the typed target name", async ({
    page,
  }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    const modal = backdropWith(page, "Delete");
    await expect(modal).toBeVisible();

    const confirmButton = modal.getByRole("button", { name: "Yes, Delete!" });
    const typedInput = modal.getByRole("textbox", {
      name: "Type the name to confirm",
    });
    await expect(confirmButton).toBeDisabled();

    await typedInput.fill("wrong-name");
    await expect(confirmButton).toBeDisabled();

    await typedInput.fill("primary-database");
    await expect(confirmButton).toBeEnabled();
  });

  test("@smoke delete dialog interpolates the target name and id", async ({
    page,
  }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    const modal = backdropWith(page, "Delete");
    await expect(modal).toBeVisible();

    await expect(modal.locator("strong").first()).toHaveText(
      "primary-database",
    );
    await expect(modal.locator("strong span").first()).toHaveText("42");
  });

  test("@smoke critical dialog gates confirm behind the typed target name", async ({
    page,
  }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Critical$/ })
      .click();
    const modal = backdropWith(page, "Critical");
    await expect(modal).toBeVisible();

    const confirmButton = modal.getByRole("button", { name: "Yes, Proceed!" });
    const typedInput = modal.getByRole("textbox", {
      name: "Type the name to confirm",
    });
    await expect(confirmButton).toBeDisabled();

    await typedInput.fill("primary-database");
    await expect(confirmButton).toBeEnabled();
  });

  test("confirm and warning presets do not gate the confirm action", async ({
    page,
  }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Confirm$/ })
      .click();
    const confirmModal = backdropWith(page, "Confirm Action");
    await expect(confirmModal).toBeVisible();
    await expect(
      confirmModal.getByRole("button", { name: "Yes, Proceed!" }),
    ).toBeEnabled();
    await expect(
      confirmModal.getByRole("textbox", { name: "Type the name to confirm" }),
    ).toHaveCount(0);
    await confirmModal.click({ position: { x: 5, y: 5 } });

    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Warning$/ })
      .click();
    const warningModal = backdropWith(page, "Warning");
    await expect(warningModal).toBeVisible();
    await expect(
      warningModal.getByRole("button", { name: "Yes, Proceed!" }),
    ).toBeEnabled();
  });

  test("confirm runs the confirm function", async ({ page }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    const modal = backdropWith(page, "Delete");
    await expect(modal).toBeVisible();

    await modal
      .getByRole("textbox", { name: "Type the name to confirm" })
      .fill("primary-database");
    await modal.getByRole("button", { name: "Yes, Delete!" }).click();

    await expect(modal).toBeHidden();
    await expect(
      page.locator(confirmationSection).getByText("deleted primary-database"),
    ).toBeVisible();
  });

  test("typed confirmation resets when the modal reopens", async ({ page }) => {
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    const modal = backdropWith(page, "Delete");
    await modal
      .getByRole("textbox", { name: "Type the name to confirm" })
      .fill("primary-database");
    await expect(
      modal.getByRole("button", { name: "Yes, Delete!" }),
    ).toBeEnabled();
    await modal.click({ position: { x: 5, y: 5 } });
    await expect(modal).toBeHidden();

    await page
      .locator(confirmationSection)
      .getByRole("button", { name: /Delete$/ })
      .click();
    await expect(modal).toBeVisible();
    await expect(
      modal.getByRole("button", { name: "Yes, Delete!" }),
    ).toBeDisabled();
  });

  test("@smoke typed prompt renders its field and gates confirm behind the expected phrase", async ({
    page,
  }) => {
    await openExamplePanel(page, confirmationSection, "Type-to-Confirm");
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: "Open Typed Prompt" })
      .click();
    const modal = backdropWith(page, "Typed Prompt");
    await expect(modal).toBeVisible();

    const confirmButton = modal.getByRole("button", { name: "Yes, Proceed!" });
    const typedInput = modal.getByRole("textbox", {
      name: 'Type "proceed" to confirm',
    });
    await expect(typedInput).toBeVisible();
    await expect(confirmButton).toBeDisabled();

    await typedInput.fill("nope");
    await expect(confirmButton).toBeDisabled();

    await typedInput.fill("proceed");
    await expect(confirmButton).toBeEnabled();
    await confirmButton.click();

    await expect(modal).toBeHidden();
    await expect(
      page.locator(confirmationSection).getByText("Confirmed: true"),
    ).toBeVisible();
  });

  test("typed prompt dialog hugs its content height", async ({ page }) => {
    await openExamplePanel(page, confirmationSection, "Type-to-Confirm");
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: "Open Typed Prompt" })
      .click();
    const modal = backdropWith(page, "Typed Prompt");
    await expect(modal).toBeVisible();

    const panel = modal.locator("div.relative.flex.flex-col").first();
    const panelBox = await panel.boundingBox();
    const viewportSize = page.viewportSize();
    expect(panelBox.height / viewportSize.height).toBeLessThan(0.5);
  });

  test("@smoke left icon position places the icon beside the title", async ({
    page,
  }) => {
    await openExamplePanel(page, confirmationSection, "Icon On Left");
    await page
      .locator(confirmationSection)
      .getByRole("button", { name: "Open Left Icon" })
      .click();
    const modal = backdropWith(page, "Icon Beside Title");
    await expect(modal).toBeVisible();

    const title = modal.locator("h2", { hasText: "Icon Beside Title" });
    await expect(title).toBeVisible();
    const icon = modal.locator("i.ph-warning-circle");
    await expect(icon).toBeVisible();
    await expect
      .poll(async () => {
        const iconBox = await icon.boundingBox();
        const titleBox = await title.boundingBox();
        const iconEndsLeftOfTitle = iconBox.x + iconBox.width <= titleBox.x;
        const iconSitsOnTheTitleRow = iconBox.y + iconBox.height > titleBox.y;
        return iconEndsLeftOfTitle && iconSitsOnTheTitleRow;
      })
      .toBe(true);
  });
});

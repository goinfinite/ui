import { expect, test } from "@playwright/test";

const modalSection = "#modal-demo";

function backdropWith(page, text) {
  return page.locator("div.fixed.inset-0.z-100", { hasText: text });
}

test.describe("Modal", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${modalSection.slice(1)}`);
  });

  test("@smoke uncloseable modal stays open on a backdrop click", async ({
    page,
  }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open", exact: true })
      .click();
    const modal = backdropWith(page, "Uncloseable Modal");
    await expect(modal).toBeVisible();
    await expect(modal.locator("i.ph-x")).toHaveCount(0);
    await modal.click({ position: { x: 5, y: 5 } });
    await expect(modal).toBeVisible();
  });

  test("@smoke uncloseable modal closes from an in-modal button", async ({
    page,
  }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open", exact: true })
      .click();
    const modal = backdropWith(page, "Uncloseable Modal");
    await expect(modal).toBeVisible();
    await page.getByRole("button", { name: "Close Modal" }).click();
    await expect(modal).toBeHidden();
  });

  test("closeable modal still closes on a backdrop click", async ({ page }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Blue", exact: true })
      .click();
    const modal = backdropWith(page, "Blue Backdrop Modal");
    await expect(modal).toBeVisible();
    await modal.click({ position: { x: 5, y: 5 } });
    await expect(modal).toBeHidden();
  });

  test("backdrop close runs the close callback", async ({ page }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Show Modal" })
      .click();
    const modal = backdropWith(page, "Interactive Modal");
    await expect(modal).toBeVisible();
    await modal.click({ position: { x: 5, y: 5 } });
    await expect(modal).toBeHidden();
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Show Modal" })
      .click();
    await expect(backdropWith(page, "Closed by backdrop click")).toBeVisible();
  });
});

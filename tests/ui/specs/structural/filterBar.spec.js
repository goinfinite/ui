import { expect, test } from "@playwright/test";

const filterBarSection = "#filter-bar-demo";
const stateReadout = "#filter-bar-demo-state span";

test.describe("FilterBar @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${filterBarSection.slice(1)}`);
    await expect(page.locator(filterBarSection)).toBeVisible();
  });

  test("@smoke active value renders as a removable chip", async ({ page }) => {
    const chip = page.locator(filterBarSection).getByText("Status: running");
    await expect(chip).toBeVisible();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":"running"',
    );

    await page
      .locator(`${filterBarSection} button[aria-label="Remove Status filter"]`)
      .click();
    await expect(chip).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });

  test("@smoke text filter updates the value and the chip", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} input[name=name]`).fill("alpha");

    await expect(
      page.locator(filterBarSection).getByText("Name: alpha"),
    ).toBeVisible();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"name":"alpha"',
    );
  });

  test("@smoke number range fills the chip with the bounds", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} input[name=cpuMin]`).fill("2");
    await expect(
      page.locator(filterBarSection).getByText("CPU: ≥ 2"),
    ).toBeVisible();

    await page.locator(`${filterBarSection} input[name=cpuMax]`).fill("8");
    await expect(
      page.locator(filterBarSection).getByText("CPU: 2–8"),
    ).toBeVisible();
  });

  test("@smoke number range chip shows the maximum alone", async ({ page }) => {
    await page.locator(`${filterBarSection} input[name=cpuMax]`).fill("8");
    await expect(
      page.locator(filterBarSection).getByText("CPU: ≤ 8"),
    ).toBeVisible();
  });

  test("@smoke enum select updates the chip", async ({ page }) => {
    await page.locator(`${filterBarSection} [role=button]`).click();
    await page
      .locator(`${filterBarSection} ul li label`)
      .filter({ hasText: "stopped" })
      .click();

    await expect(
      page.locator(filterBarSection).getByText("Status: stopped"),
    ).toBeVisible();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":"stopped"',
    );
  });

  test("@smoke clearing the enum select empties the value and hides the chip", async ({
    page,
  }) => {
    await page.locator(`${filterBarSection} .ph-x-circle`).first().click();

    await expect(
      page.locator(filterBarSection).getByText("Status: running"),
    ).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });

  test("@smoke clear filters resets every value", async ({ page }) => {
    await page.locator(`${filterBarSection} input[name=name]`).fill("alpha");
    await page
      .locator(`${filterBarSection} button`)
      .filter({ hasText: "Clear filters" })
      .click();

    await expect(
      page.locator(filterBarSection).getByText("Name: alpha"),
    ).toBeHidden();
    await expect(
      page.locator(filterBarSection).getByText("Status: running"),
    ).toBeHidden();
    await expect(page.locator(stateReadout).first()).toContainText('"name":""');
    await expect(page.locator(stateReadout).first()).toContainText(
      '"status":""',
    );
  });
});

import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const modalSection = "#modal-demo";
const closeTransitionDurationMs = 350;
const resizeTransitionDurationMs = 350;

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
    await openExamplePanel(page, modalSection, "Uncloseable");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open", exact: true })
      .click();
    const modal = backdropWith(page, "Uncloseable Modal");
    await expect(modal).toBeVisible();
    await expect(modal.locator("i.ph-x")).toHaveCount(0);
    await modal.dispatchEvent("mousedown");
    await expect(modal).toBeVisible();
  });

  test("@smoke a footer submit button submits the form in the middle content", async ({
    page,
  }) => {
    const footerFormPanel = await openExamplePanel(
      page,
      modalSection,
      "Footer Form Submit",
    );
    await footerFormPanel
      .getByRole("button", { name: "Open form modal" })
      .click();
    const modal = backdropWith(page, "Edit Profile");
    await expect(modal).toBeVisible();

    const saveButton = modal.getByRole("button", { name: "Save" });
    await expect(saveButton).toHaveAttribute("form", "modal-footer-form-demo");
    await expect(modal.locator("form#modal-footer-form-demo")).toBeVisible();
    await saveButton.click();
    await expect(modal).toBeHidden();
  });

  test("@smoke a tracked footer submit follows the form named by FormId", async ({
    page,
  }) => {
    const footerFormPanel = await openExamplePanel(
      page,
      modalSection,
      "Footer Form Submit",
    );
    let releaseSaveRequest = () => {};
    await page.route(/ui-logo\.svg/, async (route) => {
      await new Promise((resolve) => {
        releaseSaveRequest = resolve;
      });
      await route.continue();
    });
    await footerFormPanel
      .getByRole("button", { name: "Open form modal" })
      .click();
    const modal = backdropWith(page, "Edit Profile");
    await expect(modal).toBeVisible();
    const saveButton = page.locator('button[form="modal-footer-form-demo"]');

    await page.waitForFunction(() => {
      const saveButton = document.querySelector(
        "button[form=modal-footer-form-demo]",
      );
      if (!saveButton) {
        return false;
      }
      const actionData = window.Alpine?.$data(saveButton.closest("[x-data]"));
      return actionData?.submittedForm !== null;
    });
    await saveButton.evaluate((button) => button.click());
    await expect(saveButton).toBeDisabled();
    await expect(saveButton).toHaveAttribute("aria-busy", "true");
    releaseSaveRequest();
    await expect(modal).toBeHidden({ timeout: 5000 });
  });

  test("@smoke uncloseable modal closes from an in-modal button", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Uncloseable");
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
    await openExamplePanel(page, modalSection, "Appearance");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Blue", exact: true })
      .click();
    const modal = backdropWith(page, "Blue Backdrop Modal");
    await expect(modal).toBeVisible();
    await modal.dispatchEvent("mousedown");
    await expect(modal).toBeHidden();
  });

  test("@browser dragging from the panel to the backdrop keeps the modal open", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Appearance");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Blue", exact: true })
      .click();
    const modal = backdropWith(page, "Blue Backdrop Modal");
    await expect(modal).toBeVisible();

    const panel = modal.locator("div.relative.flex.flex-col").first();
    const panelBox = await panel.boundingBox();
    const backdropBox = await modal.boundingBox();
    await page.mouse.move(
      panelBox.x + panelBox.width / 2,
      panelBox.y + panelBox.height / 2,
    );
    await page.mouse.down();
    await page.mouse.move(backdropBox.x + 5, backdropBox.y + 5);
    await page.mouse.up();

    await page.waitForTimeout(closeTransitionDurationMs);
    await expect(modal).toBeVisible();
  });

  test("backdrop close runs the close callback", async ({ page }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Show Modal" })
      .click();
    const modal = backdropWith(page, "Interactive Modal");
    await expect(modal).toBeVisible();
    await modal.dispatchEvent("mousedown");
    await expect(modal).toBeHidden();
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Show Modal" })
      .click();
    await expect(backdropWith(page, "Closed by backdrop click")).toBeVisible();
  });

  async function openModalPanel(page, text) {
    const modal = backdropWith(page, text);
    await expect(modal).toBeVisible();
    const panel = modal.locator("div.relative.flex.flex-col").first();
    const viewportSize = page.viewportSize();
    const panelBox = await panel.boundingBox();
    return {
      modal,
      panel,
      viewportSize,
      widthRatio: panelBox.width / viewportSize.width,
      heightRatio: panelBox.height / viewportSize.height,
    };
  }

  test("@browser @smoke near-full terminal modal spans about 90% of the viewport", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Near-Full Size");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open XXL 90%" })
      .click();
    const { modal, widthRatio, heightRatio } = await openModalPanel(
      page,
      "WebTerminal Modal",
    );
    expect(widthRatio).toBeGreaterThan(0.85);
    expect(heightRatio).toBeGreaterThan(0.85);
    await expect(modal.getByText("agent@opencode").first()).toBeVisible();
  });

  test("@browser enlarge and reduce step through every size and hide at the ends", async ({
    page,
  }) => {
    await page
      .locator(modalSection)
      .getByRole("button", { name: "XS", exact: true })
      .click();
    const { modal, panel, viewportSize } = await openModalPanel(
      page,
      "Interactive Modal",
    );
    const enlargeButton = modal.locator("button:has(i.ph-arrows-out)");
    const reduceButton = modal.locator("button:has(i.ph-arrows-in)");

    await expect(reduceButton).toBeHidden();
    await expect(enlargeButton).toBeVisible();

    const expectedWidthRatios = {
      xs: 0.4,
      sm: 0.5,
      md: 0.6,
      lg: 0.7,
      xl: 0.8,
      xxl: 0.9,
      full: 1.0,
    };
    for (const step of ["sm", "md", "lg", "xl", "xxl", "full"]) {
      await enlargeButton.click();
      await page.waitForTimeout(resizeTransitionDurationMs);
      const panelBox = await panel.boundingBox();
      expect(
        Math.abs(
          panelBox.width / viewportSize.width - expectedWidthRatios[step],
        ),
      ).toBeLessThan(0.05);
    }

    await expect(enlargeButton).toBeHidden();
    await expect(reduceButton).toBeVisible();

    for (const step of ["xxl", "xl", "lg", "md", "sm", "xs"]) {
      await reduceButton.click();
      await page.waitForTimeout(resizeTransitionDurationMs);
      const panelBox = await panel.boundingBox();
      expect(
        Math.abs(
          panelBox.width / viewportSize.width - expectedWidthRatios[step],
        ),
      ).toBeLessThan(0.05);
    }

    await expect(reduceButton).toBeHidden();
    await expect(enlargeButton).toBeVisible();
  });

  test("@smoke custom header modal renders a close button and closes", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Header & Footer Slots");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open Slots Modal" })
      .click();
    const modal = backdropWith(page, "Slots Modal");
    await expect(modal).toBeVisible();
    const closeButton = modal.locator("div.flex.items-center.gap-1 i.ph-x");
    await closeButton.click();
    await expect(modal).toBeHidden();
  });

  test("@browser @smoke width and height settings pin the panel dimensions", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Custom Dimensions");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open Pinned Dimensions" })
      .click();
    const { modal, widthRatio, heightRatio } = await openModalPanel(
      page,
      "Pinned Dimensions Modal",
    );
    expect(Math.abs(widthRatio - 0.7)).toBeLessThan(0.05);
    expect(Math.abs(heightRatio - 0.45)).toBeLessThan(0.05);
    await expect(
      modal.locator('button:has(i[class*="ph-arrows"])'),
    ).toHaveCount(0);
  });

  test("@browser @smoke size slice steps only through its entries", async ({
    page,
  }) => {
    await openExamplePanel(page, modalSection, "Size Slice");
    await page
      .locator(modalSection)
      .getByRole("button", { name: "Open Size Slice" })
      .click();
    const { modal, panel, viewportSize } = await openModalPanel(
      page,
      "Size Slice Modal",
    );
    const enlargeButton = modal.locator("button:has(i.ph-arrows-out)");
    const reduceButton = modal.locator("button:has(i.ph-arrows-in)");

    await expect(reduceButton).toBeHidden();

    const expectedWidthRatios = { md: 0.6, lg: 0.7, xl: 0.8 };
    for (const step of ["lg", "xl"]) {
      await enlargeButton.click();
      await page.waitForTimeout(resizeTransitionDurationMs);
      const panelBox = await panel.boundingBox();
      expect(
        Math.abs(
          panelBox.width / viewportSize.width - expectedWidthRatios[step],
        ),
      ).toBeLessThan(0.05);
    }
    await expect(enlargeButton).toBeHidden();

    for (const step of ["lg", "md"]) {
      await reduceButton.click();
      await page.waitForTimeout(resizeTransitionDurationMs);
      const panelBox = await panel.boundingBox();
      expect(
        Math.abs(
          panelBox.width / viewportSize.width - expectedWidthRatios[step],
        ),
      ).toBeLessThan(0.05);
    }
    await expect(reduceButton).toBeHidden();
  });
});

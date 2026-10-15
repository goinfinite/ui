import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const sidebarSection = "#sidebar-demo";
const mainSidebar = `${sidebarSection} aside`;

function sidebarWidth(sidebar) {
  return sidebar.evaluate((element) => element.getBoundingClientRect().width);
}

test.describe("Sidebar @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${sidebarSection.slice(1)}`);
    await expect(page.locator(mainSidebar)).toBeVisible();
  });

  test("@smoke the toggle hides and restores the sidebar", async ({ page }) => {
    const sidebar = page.locator(mainSidebar);

    await expect(sidebar).toBeVisible();

    await page.getByRole("button", { name: "Toggle Visibility" }).click();
    await expect(sidebar).toBeHidden();

    await page.getByRole("button", { name: "Toggle Visibility" }).click();
    await expect(sidebar).toBeVisible();
  });

  test("@browser the collapse toggle narrows the sidebar and hides the labels", async ({
    page,
  }) => {
    const sidebar = page.locator(mainSidebar);
    const expandedWidth = await sidebarWidth(sidebar);
    await expect(sidebar.getByText("Navigation")).toBeVisible();

    await page.getByRole("button", { name: "Toggle Collapsed" }).click();

    await expect.poll(() => sidebarWidth(sidebar)).toBeLessThan(expandedWidth);
    await expect(sidebar.getByText("Navigation")).toBeHidden();
    await expect(sidebar.getByText("N", { exact: true })).toBeVisible();
  });

  test("@browser the collapse handle toggles the sidebar", async ({ page }) => {
    const sidebar = page.locator(mainSidebar);
    const expandedWidth = await sidebarWidth(sidebar);

    await sidebar.locator("i.ph-arrow-square-right").click();

    await expect.poll(() => sidebarWidth(sidebar)).toBeLessThan(expandedWidth);
  });

  test("@browser switching the side flips the sidebar float", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, sidebarSection, "Side");
    const dynamicSidebar = panel.locator("aside").first();
    const floatOf = () =>
      dynamicSidebar.evaluate((element) => getComputedStyle(element).float);

    await expect.poll(floatOf).toBe("left");

    await panel.getByRole("button", { name: "Right" }).click();

    await expect.poll(floatOf).toBe("right");
  });

  test("@browser a fixed attachment pins the sidebar to the viewport", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      sidebarSection,
      "Attachment Modes",
    );
    const sidebar = panel.locator("aside").first().locator("div").first();

    await expect(sidebar).toHaveCSS("position", "relative");

    await panel.getByRole("button", { name: "Fixed" }).click();

    await expect(sidebar).toHaveCSS("position", "fixed");

    await panel.getByRole("button", { name: "Inline" }).click();

    await expect(sidebar).toHaveCSS("position", "relative");
  });

  test("@browser an off-canvas sidebar overlays its container", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, sidebarSection, "Off-Canvas");
    const sidebar = panel.locator("aside").first().locator("div").first();

    await expect(sidebar).toHaveCSS("position", "relative");

    await panel.locator("button").filter({ hasText: "On" }).click();

    await expect(sidebar).toHaveCSS("position", "absolute");
  });

  test("@browser a static collapsed sidebar renders at the narrow width", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      sidebarSection,
      "Width & Background",
    );
    const [wideSidebar, narrowSidebar] = await panel.locator("aside").all();

    const wideWidth = await sidebarWidth(wideSidebar);
    const narrowWidth = await sidebarWidth(narrowSidebar);

    expect(narrowWidth).toBeLessThan(wideWidth);
  });
});

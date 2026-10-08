import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const tabsSection = "#tabs-demo";
const mainTabs = "#tabs-demo-main";

test.describe("Tabs @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#tabs-demo");
    await expect(page.locator(`${mainTabs} [role=tablist]`)).toBeVisible();
  });

  test("@smoke renders a tablist with the first tab selected and its panel visible", async ({
    page,
  }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await expect(tablist.locator("[role=tab]")).toHaveCount(3);
    await expect(tablist.locator("[role=tab][aria-selected=true]")).toHaveCount(
      1,
    );
    await expect(tablist.locator("[role=tab][aria-selected=true]")).toHaveText(
      /General/,
    );

    const panels = page.locator(`${mainTabs} [role=tabpanel]`);
    await expect(panels.filter({ hasText: "General" }).first()).toBeVisible();
    await expect(panels.filter({ hasText: "Security" }).first()).toBeHidden();
  });

  test("@smoke clicking a tab selects it and switches the visible panel", async ({
    page,
  }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await tablist.locator("[role=tab][data-tab-value=security]").click();

    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "security");

    const panels = page.locator(`${mainTabs} [role=tabpanel]`);
    await expect(panels.filter({ hasText: "Security" }).first()).toBeVisible();
    await expect(panels.filter({ hasText: "General" }).first()).toBeHidden();
  });

  test("@smoke badge count renders on the tab", async ({ page }) => {
    const securityTab = page
      .locator(`${mainTabs} [role=tab]`)
      .filter({ hasText: "Security" });
    await expect(securityTab.locator("span").last()).toHaveText("3");
  });

  test("reactive badge count follows its state path and hides at zero", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Badge Counts");

    const notificationsTab = panel
      .locator("[role=tab]")
      .filter({ hasText: "Notifications" });
    const badge = notificationsTab.locator("span").last();
    await expect(badge).toHaveText("5");

    const badgeClasses = await badge.getAttribute("class");
    expect(badgeClasses).toContain("bg-red-500/20");
    expect(badgeClasses).toContain("ring-1.5");
    expect(badgeClasses).toContain("ring-red-500/50");
    expect(badgeClasses).toContain("text-red-50");

    await panel.getByRole("button", { name: "Clear Alerts" }).click();
    await expect(badge).not.toBeVisible();

    await panel.getByRole("button", { name: "Add Alert" }).click();
    await expect(badge).toHaveText("1");
  });

  test("arrow keys move focus and selection to the neighbor tab", async ({
    page,
  }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await tablist.locator("[role=tab][data-tab-value=general]").focus();
    await page.keyboard.press("ArrowRight");

    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "security");
    await expect(
      tablist.locator("[role=tab][data-tab-value=security]"),
    ).toHaveAttribute("tabindex", "0");
    await expect(
      tablist.locator("[role=tab][data-tab-value=general]"),
    ).toHaveAttribute("tabindex", "-1");
  });

  test("Home and End jump to the first and last tab", async ({ page }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await tablist.locator("[role=tab][data-tab-value=security]").focus();
    await page.keyboard.press("End");
    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "advanced");

    await page.keyboard.press("Home");
    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "general");
  });

  test("vertical orientation lays the list out as a column beside the panels", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Vertical");

    const tablist = panel.locator("[role=tablist]").first();
    await expect(tablist.locator("[role=tab]")).toHaveCount(3);

    const firstTab = tablist.locator("[role=tab]").first();
    const tabClasses = await firstTab.getAttribute("class");
    expect(tabClasses).toContain("w-full");

    const scrollsX = await tablist.evaluate(
      (element) => element.scrollWidth > element.clientWidth,
    );
    expect(scrollsX).toBe(false);

    await tablist.locator("[role=tab][data-tab-value=general]").focus();
    await page.keyboard.press("ArrowDown");
    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "security");
  });

  test("right side places the vertical tab list beside the panels on the right", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Right Side");

    const tablist = panel.locator("[role=tablist]").first();
    await expect(tablist.locator("[role=tab]")).toHaveCount(3);
    await expect(tablist).toHaveAttribute("aria-orientation", "vertical");

    const rootClasses = await tablist.locator("xpath=..").getAttribute("class");
    expect(rootClasses).toContain("flex-row-reverse");

    const firstTab = tablist.locator("[role=tab]").first();
    const tabClasses = await firstTab.getAttribute("class");
    expect(tabClasses).toContain("rounded-r-md");
    expect(tabClasses).toContain("border-l-2");

    await firstTab.focus();
    await page.keyboard.press("ArrowDown");
    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "security");
  });

  test("icon position top stacks the icon above the label", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Icon Position");

    const firstTab = panel.locator("[role=tab]").first();
    const tabClasses = await firstTab.getAttribute("class");
    expect(tabClasses).toContain("flex-col");
    expect(tabClasses).toContain("text-center");

    await expect(firstTab.locator("i.ph-duotone")).toBeVisible();
    await expect(firstTab.locator("span").first()).toHaveText("General");
  });

  test("border radius rounds the tab edges", async ({ page }) => {
    const panel = await openExamplePanel(page, tabsSection, "Border Radius");

    const firstTab = panel.locator("[role=tab]").first();
    const tabClasses = await firstTab.getAttribute("class");
    expect(tabClasses).toContain("rounded-t-xl");
  });

  test("vertical alignment centers the tab list with the content", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      tabsSection,
      "Vertical Alignment",
    );

    const tablist = panel.locator("[role=tablist]").first();
    const tablistClasses = await tablist.getAttribute("class");
    expect(tablistClasses).toContain("self-center");

    const listBox = await tablist.boundingBox();
    const rootBox = await tablist.locator("xpath=..").boundingBox();
    const listCenter = listBox.y + listBox.height / 2;
    const rootCenter = rootBox.y + rootBox.height / 2;
    expect(Math.abs(listCenter - rootCenter)).toBeLessThan(2);

    const scrollsX = await tablist.evaluate(
      (element) => element.scrollWidth > element.clientWidth,
    );
    expect(scrollsX).toBe(false);
  });

  test("vertical overflow tab list scrolls vertically", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      tabsSection,
      "Vertical Overflow Scroll",
    );

    const tablist = panel.locator("[role=tablist]").first();
    await expect(tablist.locator("[role=tab]")).toHaveCount(9);
    const tablistClasses = await tablist.getAttribute("class");
    expect(tablistClasses).toContain("max-h-[60%]");
    expect(tablistClasses).toContain("overflow-y-auto");

    const scrolls = await tablist.evaluate(
      (element) => element.scrollHeight > element.clientHeight,
    );
    expect(scrolls).toBe(true);

    const scrollsX = await tablist.evaluate(
      (element) => element.scrollWidth > element.clientWidth,
    );
    expect(scrollsX).toBe(false);
  });

  test("surface customization applies the background, padding, ring, radius, and text color", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      tabsSection,
      "Surface Customization",
    );

    const tablist = panel.locator("[role=tablist]").first();
    const tablistClasses = await tablist.getAttribute("class");
    expect(tablistClasses).toContain("bg-neutral-50/5");
    expect(tablistClasses).toContain("p-3");
    expect(tablistClasses).toContain("ring-2");
    expect(tablistClasses).toContain("ring-neutral-50/10");

    const firstTab = panel.locator("[role=tab]").first();
    const tabClasses = await firstTab.getAttribute("class");
    expect(tabClasses).toContain("rounded-t-none");
    expect(tabClasses).toContain("text-neutral-50/60");
  });

  test("content surface applies the background, padding, ring, radius, shadow, and text color", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Content Surface");

    const panels = panel.locator("[role=tabpanel]").first().locator("xpath=..");
    const panelsClasses = await panels.getAttribute("class");
    expect(panelsClasses).toContain("bg-neutral-50/5");
    expect(panelsClasses).toContain("p-5");
    expect(panelsClasses).toContain("ring-2");
    expect(panelsClasses).toContain("ring-neutral-50/10");
    expect(panelsClasses).toContain("rounded");
    expect(panelsClasses).toContain("shadow-md");
    expect(panelsClasses).toContain("text-neutral-50/70");
    expect(panelsClasses).not.toContain("pt-4");

    await expect(panel.locator("[role=tabpanel]").first()).toBeVisible();
  });

  test("overflow tab list scrolls horizontally", async ({ page }) => {
    const panel = await openExamplePanel(page, tabsSection, "Overflow Scroll");

    const tablist = panel.locator("[role=tablist]").first();
    await expect(tablist.locator("[role=tab]")).toHaveCount(9);
    await expect(tablist).toHaveClass(/max-w-\[80%\]/);

    const scrolls = await tablist.evaluate(
      (element) => element.scrollWidth > element.clientWidth,
    );
    expect(scrolls).toBe(true);
  });

  test("deep-linked tab writes the hash on click and reads it on change", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, tabsSection, "Deep Linking");

    const tablist = panel.locator("[role=tablist]").first();
    await tablist.locator("[role=tab][data-tab-value=security]").click();
    await expect(page).toHaveURL(/#security$/);

    await page.evaluate(() => {
      location.hash = "advanced";
    });
    await expect(
      tablist.locator("[role=tab][aria-selected=true]"),
    ).toHaveAttribute("data-tab-value", "advanced");
  });

  test("vertical overflow caps the list to MaxVisibleTabs tabs", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      tabsSection,
      "Vertical Overflow Scroll",
    );
    const tabList = panel.locator(
      "#tabs-demo-vertical-overflow-visible [data-ui-tabs-list]",
    );
    const geometry = () =>
      tabList.evaluate((element) => {
        const tab = element.querySelector("[role=tab]");
        const styles = getComputedStyle(element);
        const rowGap = Number.parseFloat(styles.rowGap) || 0;
        const verticalPadding =
          (Number.parseFloat(styles.paddingTop) || 0) +
          (Number.parseFloat(styles.paddingBottom) || 0);
        return {
          visibleTabs: Math.round(
            (element.clientHeight - verticalPadding) /
              (tab.getBoundingClientRect().height + rowGap),
          ),
          scrolls: element.scrollHeight > element.clientHeight,
        };
      });

    await expect.poll(geometry).toEqual({ visibleTabs: 4, scrolls: true });
  });
});

test.describe("Tabs accessibility @structural @a11y", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/index.html#tabs-demo");
    await expect(page.locator(`${mainTabs} [role=tablist]`)).toBeVisible();
  });

  test("tabs and panels are wired with aria-controls, aria-labelledby, and aria-selected", async ({
    page,
  }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await expect(tablist).toHaveAttribute("aria-label", "Tabs");

    const firstTab = tablist.locator("[role=tab]").first();
    const controlsId = await firstTab.getAttribute("aria-controls");
    const panel = page.locator(`#${controlsId}`);
    await expect(panel).toHaveAttribute("role", "tabpanel");
    await expect(panel).toHaveAttribute(
      "aria-labelledby",
      await firstTab.getAttribute("id"),
    );
    await expect(firstTab).toHaveAttribute("aria-selected", "true");
  });

  test("only the selected tab is in the tab order", async ({ page }) => {
    const tablist = page.locator(`${mainTabs} [role=tablist]`);
    await expect(tablist.locator('[role=tab][tabindex="0"]')).toHaveCount(1);
    await expect(tablist.locator('[role=tab][tabindex="-1"]')).toHaveCount(2);
  });
});

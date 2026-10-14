import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const chartSection = "#chart-demo";
const liveChartTableRows = `${chartSection} table.sr-only tbody tr`;

test.describe("Chart", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html${chartSection}`);
  });

  test("@smoke the live chart lists its series and trims the streamed rows", async ({
    page,
  }) => {
    const legendLabels = page.locator(`${chartSection} ul li span[x-text]`);
    await expect(legendLabels.first()).toHaveText("CPU");
    await expect(legendLabels.nth(1)).toHaveText("Memory");

    const firstRowTime = page.locator(liveChartTableRows).first().locator("th");
    await expect(page.locator(liveChartTableRows)).toHaveCount(20);
    const initialFirstRowTime = await firstRowTime.textContent();

    await expect(firstRowTime).not.toHaveText(initialFirstRowTime, {
      timeout: 5000,
    });
    await expect(page.locator(liveChartTableRows)).toHaveCount(20);
  });

  test("an empty chart shows the empty label", async ({ page }) => {
    const emptyPanel = await openExamplePanel(
      page,
      chartSection,
      "Empty State",
    );

    await expect(
      emptyPanel.getByText("No requests yet", { exact: true }),
    ).toBeVisible();
  });

  test("@browser a category axis keeps the row order", async ({ page }) => {
    const barPanel = await openExamplePanel(page, chartSection, "Bar");
    const axisLabels = barPanel.locator("svg text");
    await expect(axisLabels.first()).toHaveText("Jan");

    const monthLabels = await axisLabels.evaluateAll((labelElements) =>
      labelElements
        .map((labelElement) => labelElement.textContent)
        .filter((labelText) => /^[A-Z][a-z]{2}$/.test(labelText)),
    );
    expect(monthLabels).toEqual(["Jan", "Feb", "Mar", "Apr", "May", "Jun"]);
  });

  test("@browser an identity color token paints the bars and the legend swatch", async ({
    page,
  }) => {
    const stackedPanel = await openExamplePanel(page, chartSection, "Stacked");
    const receivedSwatch = stackedPanel
      .locator("ul li span[aria-hidden]")
      .first();
    await expect(receivedSwatch).toHaveCSS(
      "background-color",
      "rgb(14, 165, 233)",
    );

    const barFills = await stackedPanel
      .locator("svg path[fill]")
      .evaluateAll((barElements) =>
        barElements.map((barElement) => barElement.getAttribute("fill")),
      );
    expect(barFills).toContain("rgb(14, 165, 233)");
  });

  test("@browser hovering the plot shows every series value", async ({
    page,
  }) => {
    const chartCanvas = page
      .locator(`${chartSection} [x-ref=chartCanvas]`)
      .first();
    await expect(chartCanvas.locator("svg")).toBeVisible();
    await chartCanvas.scrollIntoViewIfNeeded();
    const canvasBox = await chartCanvas.boundingBox();

    await page.mouse.move(
      canvasBox.x + canvasBox.width * 0.6,
      canvasBox.y + canvasBox.height * 0.5,
    );

    const tooltipRows = page.locator(
      `${chartSection} div.fixed div.flex-row span.text-neutral-300`,
    );
    await expect(tooltipRows.first()).toBeVisible();
    await expect(tooltipRows).toHaveText(["CPU", "Memory"]);
  });
});

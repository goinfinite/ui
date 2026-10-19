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
    const legendItems = page.locator(`${chartSection} ul li`);
    await expect(legendItems.first()).toHaveText("CPU");
    await expect(legendItems.nth(1)).toHaveText("Memory");

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
    const panel = await openExamplePanel(page, chartSection, "Line");
    const chartCanvas = panel.locator("[x-ref=chartCanvas]").first();
    await expect(chartCanvas.locator("svg")).toBeVisible();
    await chartCanvas.scrollIntoViewIfNeeded();
    const canvasBox = await chartCanvas.boundingBox();

    await page.mouse.move(
      canvasBox.x + canvasBox.width * 0.6,
      canvasBox.y + canvasBox.height * 0.5,
    );

    const tooltipRows = panel.locator(
      "div.fixed div.flex-row span.text-neutral-300",
    );
    await expect(tooltipRows.first()).toBeVisible();
    await expect(tooltipRows).toHaveText(["Received (GiB)", "Transmitted (GiB)"]);
  });

  test("@browser horizontal bars run the category axis down the left edge", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Horizontal Bar");
    const axisLabels = panel.locator("svg text");
    await expect(axisLabels.first()).toHaveText("Jan");

    const barBoxes = await panel
      .locator("svg path[fill]")
      .evaluateAll((barElements) =>
        barElements.map((barElement) => {
          const barBox = barElement.getBBox();
          return { width: barBox.width, height: barBox.height };
        }),
      );
    expect(barBoxes.some((barBox) => barBox.width > barBox.height)).toBe(true);
  });

  test("@browser a pie draws one solid slice per row and labels each share", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Pie");
    const sliceFills = panel.locator("svg path[fill][fill-opacity='1']");
    await expect(sliceFills).toHaveCount(4);
    await expect(panel.locator("svg text")).toHaveCount(4);
    await expect(panel.locator("svg text").last()).toHaveText(/%$/);
    await expect(panel.locator("ul li").first()).toHaveText("mariadb");
  });

  test("@smoke the radial bar lists its categories and their table rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Radial Bar");
    const legendItems = panel.locator("ul li");

    await expect(legendItems).toHaveCount(4);
    await expect(legendItems.first()).toHaveText("mariadb");
    await expect(legendItems.last()).toHaveText("redis");
    await expect(panel.locator("table.sr-only tbody tr")).toHaveCount(4);
  });

  test("@browser a radial bar grows each category outward from the hole", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Radial Bar");
    const barPaths = panel.locator("svg path[fill][fill-opacity='1']");
    await expect(barPaths).toHaveCount(4);
    await expect(panel.locator("svg text")).toHaveCount(4);

    const barBoxes = await barPaths.evaluateAll((radialBarPaths) =>
      radialBarPaths.map((radialBarPath) => {
        const radialBarBox = radialBarPath.getBBox();
        return { width: radialBarBox.width };
      }),
    );
    expect(barBoxes[0].width).toBeGreaterThan(barBoxes[1].width);
  });

  test("@smoke the waffle lists its categories and their table rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Waffle");
    const legendItems = panel.locator("ul li");

    await expect(legendItems).toHaveCount(4);
    await expect(legendItems.first()).toHaveText("mariadb");
    await expect(legendItems.last()).toHaveText("redis");
    await expect(panel.locator("table.sr-only tbody tr")).toHaveCount(4);
  });

  test("@browser the waffle tiles every cell and keeps the squares equal", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Waffle");
    const cellPaths = panel.locator("svg path[fill]");
    await expect(cellPaths).toHaveCount(100);

    const cellFills = await cellPaths.evaluateAll((cellElements) =>
      cellElements.map((cellElement) => getComputedStyle(cellElement).fill),
    );
    expect(
      cellFills.filter((cellFill) => cellFill === "rgb(14, 165, 233)"),
    ).toHaveLength(60);

    const firstCellBox = await cellPaths.first().evaluate((cellElement) => {
      const cellBox = cellElement.getBBox();
      return { width: cellBox.width, height: cellBox.height };
    });
    expect(firstCellBox.width).toBeCloseTo(firstCellBox.height, 1);
  });

  test("@smoke the sunburst lists its categories and their table rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Sunburst");
    const legendItems = panel.locator("ul li");

    await expect(legendItems).toHaveCount(3);
    await expect(legendItems.first()).toHaveText("storage");
    await expect(legendItems.last()).toHaveText("backup");
    await expect(panel.locator("table.sr-only tbody tr")).toHaveCount(5);
  });

  test("@browser the sunburst draws one ring per hierarchy level", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Sunburst");
    const arcPaths = panel.locator("svg path[fill][fill-opacity]");
    await expect(arcPaths).toHaveCount(7);

    const arcFills = await arcPaths.evaluateAll((arcElements) =>
      arcElements.map((arcElement) => getComputedStyle(arcElement).fill),
    );
    expect(
      arcFills.filter((arcFill) => arcFill === "rgb(14, 165, 233)"),
    ).toHaveLength(3);

    const categoryArcs = panel.locator("svg path[fill-opacity='1']");
    await expect(categoryArcs).toHaveCount(3);
    const leafArcs = panel.locator("svg path[fill-opacity='0.6']");
    await expect(leafArcs).toHaveCount(4);

    const categoryArcFills = await categoryArcs.evaluateAll((arcElements) =>
      arcElements.map((arcElement) => getComputedStyle(arcElement).fill),
    );
    const legendSwatchColors = await panel
      .locator("ul li span:first-child")
      .evaluateAll((swatchElements) =>
        swatchElements.map(
          (swatchElement) => getComputedStyle(swatchElement).backgroundColor,
        ),
      );
    expect(categoryArcFills).toEqual(legendSwatchColors);

    const arcRings = await panel
      .locator("svg path[fill][fill-opacity]")
      .evaluateAll((arcElements) =>
        arcElements.map((arcElement) => {
          const radii = [
            ...arcElement.getAttribute("d").matchAll(/A([\d.]+),/g),
          ].map((radiusMatch) => Number(radiusMatch[1]));
          radii.sort((first, second) => second - first);
          return {
            opacity: arcElement.getAttribute("fill-opacity"),
            outerRadius: radii[0],
            innerRadius: radii[1],
          };
        }),
      );
    const categoryRings = arcRings.filter((arcRing) => arcRing.opacity === "1");
    const leafRings = arcRings.filter((arcRing) => arcRing.opacity === "0.6");
    for (const categoryRing of categoryRings) {
      for (const leafRing of leafRings) {
        expect(leafRing.innerRadius).toBeCloseTo(categoryRing.outerRadius, 0);
        expect(leafRing.outerRadius - leafRing.innerRadius).toBeCloseTo(
          categoryRing.outerRadius - categoryRing.innerRadius,
          0,
        );
      }
    }
  });
  test("@smoke the parliament lists its parties and their table rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Parliament");
    const legendItems = panel.locator("ul li");

    await expect(legendItems).toHaveCount(3);
    await expect(legendItems.first()).toHaveText("Center");
    await expect(legendItems.last()).toHaveText("Greens");
    await expect(panel.locator("table.sr-only tbody tr")).toHaveCount(3);
  });

  test("@browser the parliament seats every party total in semicircle rows", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Parliament");
    const seatPaths = panel.locator("svg path[fill]");
    await expect(seatPaths).toHaveCount(120);

    const seatFills = await seatPaths.evaluateAll((seatElements) =>
      seatElements.map((seatElement) => getComputedStyle(seatElement).fill),
    );
    expect(
      seatFills.filter((seatFill) => seatFill === "rgb(14, 165, 233)"),
    ).toHaveLength(52);
    expect(
      seatFills.filter((seatFill) => seatFill === "rgb(217, 89, 38)"),
    ).toHaveLength(41);
    expect(
      seatFills.filter((seatFill) => seatFill === "rgb(25, 158, 112)"),
    ).toHaveLength(27);

    const seatRadii = await panel
      .locator("svg")
      .first()
      .evaluate((svgElement) => {
        const chartCentreX = svgElement.width.baseVal.value / 2;
        const chartCentreY = svgElement.height.baseVal.value - 4;
        return [...svgElement.querySelectorAll("path[fill]")].map(
          (seatElement) => {
            const seatMatrix = seatElement.transform.baseVal.getItem(0).matrix;
            return {
              radius: Math.round(
                Math.hypot(
                  seatMatrix.e - chartCentreX,
                  seatMatrix.f - chartCentreY,
                ),
              ),
              isAboveCentre: seatMatrix.f < chartCentreY,
            };
          },
        );
      });
    const seatBands = seatRadii
      .map((seat) => seat.radius)
      .sort((first, second) => first - second);
    let seatBandCount = 1;
    for (let seatIndex = 1; seatIndex < seatBands.length; seatIndex++) {
      if (seatBands[seatIndex] - seatBands[seatIndex - 1] > 10) {
        seatBandCount++;
      }
    }
    expect(seatBandCount).toBe(6);
    expect(seatRadii.every((seat) => seat.isAboveCentre)).toBe(true);

    await seatPaths.first().hover();
    const tooltipContainer = panel.locator("div.fixed").first();
    await expect(tooltipContainer).toContainText("Center");
    await expect(tooltipContainer).toContainText("52");
  });
});

test.describe("Chart progress", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html${chartSection}`);
  });

  test("@smoke the progress bar shows its percent and paints the threshold color", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Progress");
    const bar = panel.getByRole("progressbar").first();
    const fill = bar.locator(":scope > div.h-full");

    await expect(bar).toHaveAttribute("aria-valuenow", "65.00");
    await expect(fill).toHaveAttribute("style", /width: 65%/);
    await expect(fill).toHaveClass(/bg-amber-300/);
  });

  test("a striped bar paints the diagonal gradient", async ({ page }) => {
    const panel = await openExamplePanel(page, chartSection, "Progress");
    const stripedFill = panel
      .locator("[role=progressbar]")
      .nth(1)
      .locator(":scope > div.h-full");

    await expect(stripedFill).toHaveCSS(
      "background-image",
      /repeating-linear-gradient/,
    );
  });

  test("label affixes wrap the percent text", async ({ page }) => {
    const panel = await openExamplePanel(page, chartSection, "Progress");
    const affixedBar = panel.locator("[role=progressbar]").nth(2);

    await expect(affixedBar).toContainText("Used 82% of capacity");
  });

  test("@smoke custom thresholds pick the color per percent", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      chartSection,
      "Progress Thresholds",
    );
    const fills = panel.locator("[role=progressbar] > div.h-full");

    await expect(fills.nth(5)).toHaveClass(/bg-rose-500/);
    await expect(fills.nth(6)).toHaveClass(/bg-amber-500/);
    await expect(fills.nth(7)).toHaveClass(/bg-emerald-500/);
  });

  test("min and max rescale the fill", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      chartSection,
      "Progress Thresholds",
    );
    const boundsBar = panel.locator("[role=progressbar]").nth(8);

    await expect(boundsBar).toHaveAttribute("aria-valuenow", "12.00");
    await expect(boundsBar.locator(":scope > div.h-full")).toHaveAttribute(
      "style",
      /width: 25%/,
    );
  });

  test("@smoke a bound slider drives the fill", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      chartSection,
      "Progress Bound Value",
    );
    const bar = panel.getByRole("progressbar", { name: "Upload progress" });

    await expect(bar).toHaveAttribute("aria-valuenow", "35.00");

    await panel.locator("input[type=range]").fill("80");

    await expect(bar).toHaveAttribute("aria-valuenow", "80.00");
    await expect(bar.locator(":scope > div.h-full")).toHaveAttribute(
      "style",
      /width: 80%/,
    );
  });

  test("the indeterminate bar animates without a percent", async ({ page }) => {
    const panel = await openExamplePanel(
      page,
      chartSection,
      "Progress Bound Value",
    );
    const bar = panel.locator("[role=progressbar]").filter({
      hasText: "Syncing",
    });

    await expect(bar).not.toHaveAttribute("aria-valuenow", /.+/);
    await expect(bar.locator(":scope > div.h-full")).toHaveAttribute(
      "style",
      /animation: ui-progress-bar-indeterminate/,
    );
  });
});

test.describe("Chart gauge", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html${chartSection}`);
  });

  test("@smoke the gauge draws its bands, needle, and center label", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Gauge");
    const gauge = panel.getByRole("meter", { name: "CPU load" });

    await expect(gauge).toHaveAttribute("aria-valuemin", "0");
    await expect(gauge).toHaveAttribute("aria-valuemax", "100");
    await expect(gauge).toHaveAttribute("aria-valuenow", "65.00");
    await expect(gauge.locator("svg")).toHaveAttribute(
      "viewBox",
      "0 0 200 200",
    );

    const bandPaths = gauge.locator("svg path");
    await expect(bandPaths).toHaveCount(3);
    await expect(bandPaths.nth(0)).toHaveClass(/stroke-emerald-500/);
    await expect(bandPaths.nth(1)).toHaveClass(/stroke-amber-500/);
    await expect(bandPaths.nth(2)).toHaveClass(/stroke-red-500/);

    await expect(gauge.locator("svg line")).toHaveAttribute(
      "transform",
      "rotate(234 100 100)",
    );
    await expect(gauge.locator("svg text")).toHaveText(["65", "CPU"]);
  });

  test("@browser the semicircle gauge keeps the needle in the top half", async ({
    page,
  }) => {
    const panel = await openExamplePanel(
      page,
      chartSection,
      "Gauge Semicircle",
    );
    const gauge = panel.getByRole("meter");

    await expect(gauge).toHaveAttribute("aria-valuenow", "72.00");
    await expect(gauge.locator("svg")).toHaveAttribute(
      "viewBox",
      "0 0 200 150",
    );
    const needle = gauge.locator("svg line");
    await expect(needle).toHaveAttribute("transform", "rotate(39.6 100 100)");

    const needlePosition = await needle.evaluate((needleElement) => {
      const screenMatrix = needleElement.getScreenCTM();
      const makePoint = (x, y) => {
        const point = needleElement.ownerSVGElement.createSVGPoint();
        point.x = x;
        point.y = y;
        return point;
      };
      const pivot = makePoint(100, 100).matrixTransform(screenMatrix);
      const tip = makePoint(100, 45).matrixTransform(screenMatrix);
      return { pivotY: pivot.y, tipY: tip.y };
    });
    expect(needlePosition.tipY).toBeLessThan(needlePosition.pivotY);
  });

  test("@browser band color tokens resolve to their stroke colors", async ({
    page,
  }) => {
    const panel = await openExamplePanel(page, chartSection, "Gauge");
    const bandPaths = panel
      .getByRole("meter", { name: "CPU load" })
      .locator("svg path");

    await expect(bandPaths.nth(0)).toHaveCSS("stroke", "rgb(16, 185, 129)");
    await expect(bandPaths.nth(1)).toHaveCSS("stroke", "rgb(245, 158, 11)");
    await expect(bandPaths.nth(2)).toHaveCSS("stroke", "rgb(239, 68, 68)");
  });
});

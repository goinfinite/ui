UiToolset.RegisterAlpineState(() => {
  const chartSeriesPalette = [
    "#3987e5",
    "#d95926",
    "#199e70",
    "#c98500",
    "#d55181",
    "#008300",
    "#9085e9",
    "#e66767",
  ];
  const chartPaletteExceededColor = "#737373";
  const chartColorProbeUnresolvedColor = "rgb(1, 2, 3)";
  const chartColorProbeMaxAttempts = 10;
  const chartColorProbeRetryDelayMs = 16;
  const chartMissingValueText = "–";
  const chartBarMaxThicknessPx = 24;
  const chartStackedBarMaxThicknessPx = 40;
  const chartLineDashByStyle = {
    solid: [1, 0],
    dashed: [6, 4],
    dotted: [2, 3],
  };
  const chartFallbackGapColor = "#171717";
  const chartFallbackTickLabelColor = "#a3a3a3";
  const chartFallbackTitleColor = "#d4d4d4";
  const chartFallbackLineColor = "rgba(250, 250, 250, 0.08)";
  const chartFallbackCrosshairColor = "rgba(250, 250, 250, 0.3)";
  const chartFallbackCenterLabelSliceColor = "#a3a3a3";
  const chartFallbackCenterLabelValueColor = "#fafafa";
  const chartFallbackColorByToken = {
    "neutral-900": chartFallbackGapColor,
    "neutral-400": chartFallbackTickLabelColor,
    "neutral-300": chartFallbackTitleColor,
    "neutral-50/8": chartFallbackLineColor,
    "neutral-50/30": chartFallbackCrosshairColor,
  };

  const chartColorNormalizationCanvas = document.createElement("canvas");
  chartColorNormalizationCanvas.width = 1;
  chartColorNormalizationCanvas.height = 1;
  const chartColorNormalizationContext =
    chartColorNormalizationCanvas.getContext("2d");

  function expressionStringQuoter(rawText) {
    return JSON.stringify(String(rawText));
  }

  function chartColorNormalizer(rawColor) {
    if (chartColorNormalizationContext === null) {
      return rawColor;
    }
    chartColorNormalizationContext.clearRect(0, 0, 1, 1);
    chartColorNormalizationContext.fillStyle = rawColor;
    chartColorNormalizationContext.fillRect(0, 0, 1, 1);
    const [red, green, blue, alpha] =
      chartColorNormalizationContext.getImageData(0, 0, 1, 1).data;
    if (alpha === 255) {
      return `rgb(${red}, ${green}, ${blue})`;
    }
    return `rgba(${red}, ${green}, ${blue}, ${alpha / 255})`;
  }

  function chartRowKeyAccessBuilder(rowKey) {
    return `datum[${expressionStringQuoter(rowKey)}]`;
  }

  function chartXKeyExpressionBuilder(xAxis) {
    const xKeyAccess = chartRowKeyAccessBuilder(xAxis.rowKey);
    if (xAxis.scaleType === "time") {
      return `time(toDate(${xKeyAccess}))`;
    }
    if (xAxis.scaleType === "number") {
      return `toNumber(${xKeyAccess})`;
    }
    return `${xKeyAccess} + ''`;
  }

  function chartRowSortFieldResolver(xAxis) {
    if (xAxis.scaleType === "category") {
      return "chartRowOrder";
    }
    return "chartX";
  }

  function chartStackTransformsBuilder(chartConfig, seriesRowKeys) {
    if (chartConfig.cartesian.stackMode === "none") {
      return [
        { type: "formula", as: "chartYBase", expr: "0" },
        { type: "formula", as: "chartYTop", expr: "datum.chartValue" },
      ];
    }
    const stackOffset =
      chartConfig.cartesian.stackMode === "percent" ? "normalize" : "zero";
    return [
      {
        type: "formula",
        as: "chartSeriesOrder",
        expr: `indexof(${JSON.stringify(seriesRowKeys)}, datum.chartSeriesRowKey)`,
      },
      {
        type: "formula",
        as: "chartStackValue",
        expr: "isValid(datum.chartValue) ? toNumber(datum.chartValue) : 0",
      },
      {
        type: "stack",
        groupby: ["chartX"],
        field: "chartStackValue",
        sort: { field: "chartSeriesOrder" },
        offset: stackOffset,
        as: ["chartYBase", "chartYTop"],
      },
    ];
  }

  function chartCartesianDataBuilder(chartConfig, seriesRowKeys) {
    return [
      {
        name: "rows",
        values: [],
        transform: [
          {
            type: "formula",
            as: "chartX",
            expr: chartXKeyExpressionBuilder(chartConfig.xAxis),
          },
          {
            type: "collect",
            sort: { field: chartRowSortFieldResolver(chartConfig.xAxis) },
          },
        ],
      },
      {
        name: "points",
        source: "rows",
        transform: [
          {
            type: "fold",
            fields: seriesRowKeys,
            as: ["chartSeriesRowKey", "chartValue"],
          },
          ...chartStackTransformsBuilder(chartConfig, seriesRowKeys),
        ],
      },
      {
        name: "validPoints",
        source: "points",
        transform: [{ type: "filter", expr: "isValid(datum.chartValue)" }],
      },
      {
        name: "hoveredRows",
        source: "rows",
        transform: [{ type: "filter", expr: "datum.chartX === hoveredX" }],
      },
      {
        name: "hoveredPoints",
        source: "validPoints",
        transform: [{ type: "filter", expr: "datum.chartX === hoveredX" }],
      },
    ];
  }

  function chartXScaleBuilder(chartConfig) {
    const rowOrderedDomain = {
      data: "rows",
      field: "chartX",
      sort: { field: chartRowSortFieldResolver(chartConfig.xAxis), op: "min" },
    };
    if (chartConfig.chartType === "bar") {
      return {
        name: "x",
        type: "band",
        domain: rowOrderedDomain,
        range: "width",
        paddingInner: 0.2,
        paddingOuter: 0.1,
      };
    }
    const continuousDomain = { data: "rows", field: "chartX" };
    if (chartConfig.xAxis.scaleType === "time") {
      return {
        name: "x",
        type: "time",
        domain: continuousDomain,
        range: "width",
      };
    }
    if (chartConfig.xAxis.scaleType === "number") {
      return {
        name: "x",
        type: "linear",
        domain: continuousDomain,
        range: "width",
        zero: false,
      };
    }
    return {
      name: "x",
      type: "point",
      domain: rowOrderedDomain,
      range: "width",
      padding: 0.5,
    };
  }

  function chartYScaleBuilder(chartConfig) {
    const yScale = {
      name: "y",
      type: "linear",
      domain: { data: "points", fields: ["chartYBase", "chartYTop"] },
      range: "height",
      zero: true,
      nice: true,
    };
    const yAxis = chartConfig.cartesian.yAxis;
    if (yAxis.minValue !== "") {
      yScale.domainMin = Number(yAxis.minValue);
      yScale.nice = false;
    }
    if (yAxis.maxValue !== "") {
      yScale.domainMax = Number(yAxis.maxValue);
      yScale.nice = false;
    }
    return yScale;
  }

  function chartPaletteColorResolver(colorIndex) {
    if (colorIndex < chartSeriesPalette.length) {
      return chartSeriesPalette[colorIndex];
    }
    return chartPaletteExceededColor;
  }

  function chartTokenColorResolver(colorByToken, colorToken, fallbackColor) {
    return colorByToken[colorToken] ?? fallbackColor;
  }

  function chartSeriesColorsResolver(chartSeries, colorByToken) {
    return chartSeries.map((seriesEntry, seriesIndex) =>
      chartTokenColorResolver(
        colorByToken,
        seriesEntry.identityColor,
        chartPaletteColorResolver(seriesIndex),
      ),
    );
  }

  function chartCartesianScalesBuilder(
    chartConfig,
    seriesRowKeys,
    colorByToken,
  ) {
    return [
      chartXScaleBuilder(chartConfig),
      chartYScaleBuilder(chartConfig),
      {
        name: "color",
        type: "ordinal",
        domain: seriesRowKeys,
        range: chartSeriesColorsResolver(chartConfig.series, colorByToken),
      },
      {
        name: "seriesLabel",
        type: "ordinal",
        domain: seriesRowKeys,
        range: chartConfig.series.map((seriesEntry) => seriesEntry.label),
      },
      {
        name: "seriesLineDash",
        type: "ordinal",
        domain: seriesRowKeys,
        range: chartConfig.series.map(
          (seriesEntry) => chartLineDashByStyle[seriesEntry.lineDashStyle],
        ),
      },
    ];
  }

  function chartAxisAppearanceBuilder(axisSettings, colorByToken) {
    const axisAppearance = {
      titleColor: chartTokenColorResolver(
        colorByToken,
        axisSettings.titleColor,
        chartFallbackTitleColor,
      ),
      labelColor: chartTokenColorResolver(
        colorByToken,
        axisSettings.tickLabelColor,
        chartFallbackTickLabelColor,
      ),
      domainColor: chartTokenColorResolver(
        colorByToken,
        axisSettings.baselineColor,
        chartFallbackLineColor,
      ),
      tickColor: chartTokenColorResolver(
        colorByToken,
        axisSettings.baselineColor,
        chartFallbackLineColor,
      ),
      gridColor: chartTokenColorResolver(
        colorByToken,
        axisSettings.gridLineColor,
        chartFallbackLineColor,
      ),
    };
    if (axisSettings.title !== "") {
      axisAppearance.title = axisSettings.title;
    }
    if (axisSettings.isHidden) {
      axisAppearance.labels = false;
      axisAppearance.ticks = false;
      axisAppearance.domain = false;
      delete axisAppearance.title;
    }
    return axisAppearance;
  }

  function chartBandLabelEncodingBuilder(xAxis, bandLabelFormat) {
    const formatLiteral = expressionStringQuoter(bandLabelFormat);
    const formatFunction = xAxis.scaleType === "time" ? "timeFormat" : "format";
    return {
      labels: {
        update: {
          text: { signal: `${formatFunction}(datum.value, ${formatLiteral})` },
        },
      },
    };
  }

  function chartXAxisBuilder(chartConfig, colorByToken) {
    const xAxisSettings = chartConfig.xAxis;
    const xAxis = {
      orient: "bottom",
      scale: "x",
      tickCount: 6,
      grid: xAxisSettings.isGridLineVisible,
      ...chartAxisAppearanceBuilder(xAxisSettings, colorByToken),
    };
    const isContinuousScale = xAxisSettings.scaleType !== "category";
    const tickLabelFormat = xAxisSettings.tickLabelFormat;
    if (isContinuousScale && tickLabelFormat !== "") {
      xAxis.format = tickLabelFormat;
      xAxis.formatType = xAxisSettings.scaleType;
    }
    const bandLabelFormat = tickLabelFormat || chartConfig.tooltip.titleFormat;
    const isFormattedBand =
      chartConfig.chartType === "bar" &&
      isContinuousScale &&
      bandLabelFormat !== "";
    if (isFormattedBand) {
      delete xAxis.format;
      delete xAxis.formatType;
      xAxis.encode = chartBandLabelEncodingBuilder(
        xAxisSettings,
        bandLabelFormat,
      );
    }
    if (xAxisSettings.isTickStartEndOnly) {
      xAxis.values = {
        signal:
          "length(domain('x')) ? [domain('x')[0], peek(domain('x'))] : []",
      };
    }
    return xAxis;
  }

  function chartYAxisBuilder(chartConfig, colorByToken) {
    const yAxisSettings = chartConfig.cartesian.yAxis;
    return {
      orient: yAxisSettings.position,
      scale: "y",
      grid: !yAxisSettings.isGridLineHidden,
      domain: false,
      ticks: false,
      tickCount: yAxisSettings.tickCount,
      format: yAxisSettings.tickLabelFormat,
      ...chartAxisAppearanceBuilder(yAxisSettings, colorByToken),
    };
  }

  function chartAreaFillBuilder(areaFillStyle) {
    if (areaFillStyle === "solid") {
      return {
        fill: { scale: "color", field: "chartSeriesRowKey" },
        fillOpacity: { value: 0.15 },
      };
    }
    const seriesColor = "rgb(scale('color', datum.chartSeriesRowKey))";
    const fadedColorExpression = (opacityRatio) =>
      `'rgba(' + ${seriesColor}.r + ', ' + ${seriesColor}.g + ', ' + ${seriesColor}.b + ', ' + (${seriesColor}.opacity * ${opacityRatio}) + ')'`;
    return {
      fill: {
        signal: `{gradient: 'linear', x1: 0, y1: 0, x2: 0, y2: 1, stops: [{offset: 0, color: ${fadedColorExpression(0.35)}}, {offset: 1, color: ${fadedColorExpression(0)}}]}`,
      },
    };
  }

  function chartClampedBaselineResolver(baselineValueExpression) {
    return {
      signal: `scale('y', clamp(${baselineValueExpression}, domain('y')[0], domain('y')[1]))`,
    };
  }

  function chartAreaBaselineResolver(chartConfig) {
    if (chartConfig.cartesian.stackMode === "none") {
      return chartClampedBaselineResolver("0");
    }
    return chartClampedBaselineResolver("datum.chartYBase");
  }

  function chartSeriesMarksBuilder(chartConfig, pointPosition) {
    const seriesMarks = [];
    const hasAreaFill =
      chartConfig.chartType === "area" &&
      chartConfig.cartesian.areaFillStyle !== "none";
    if (hasAreaFill) {
      seriesMarks.push({
        type: "area",
        from: { data: "seriesPoints" },
        encode: {
          enter: {
            interpolate: { value: "monotone" },
            ...chartAreaFillBuilder(chartConfig.cartesian.areaFillStyle),
          },
          update: {
            ...pointPosition,
            y2: chartAreaBaselineResolver(chartConfig),
          },
        },
      });
    }
    seriesMarks.push({
      type: "line",
      from: { data: "seriesPoints" },
      encode: {
        enter: {
          interpolate: { value: "monotone" },
          stroke: { scale: "color", field: "chartSeriesRowKey" },
          strokeDash: { scale: "seriesLineDash", field: "chartSeriesRowKey" },
          strokeWidth: { value: 2 },
          strokeCap: { value: "round" },
          strokeJoin: { value: "round" },
        },
        update: pointPosition,
      },
    });
    return seriesMarks;
  }

  function chartXTitleExpressionBuilder(chartConfig, xKeyAccess) {
    const titleFormat = chartConfig.tooltip.titleFormat;
    const formatLiteral = expressionStringQuoter(titleFormat);
    if (chartConfig.xAxis.scaleType === "time") {
      return `timeFormat(${xKeyAccess}, ${formatLiteral})`;
    }
    if (chartConfig.xAxis.scaleType === "number" && titleFormat !== "") {
      return `format(${xKeyAccess}, ${formatLiteral})`;
    }
    return `${xKeyAccess} + ''`;
  }

  function chartValueTextExpressionBuilder(chartConfig, valueAccess) {
    const valueFormat = expressionStringQuoter(chartConfig.tooltip.valueFormat);
    const missingText = expressionStringQuoter(chartMissingValueText);
    return `isValid(${valueAccess}) ? format(${valueAccess}, ${valueFormat}) : ${missingText}`;
  }

  function chartRowTooltipExpressionBuilder(
    chartConfig,
    rowAccess,
    seriesColors,
  ) {
    const tooltipRows = chartConfig.series.map((seriesEntry, seriesIndex) => {
      const valueAccess = `${rowAccess}[${expressionStringQuoter(seriesEntry.rowKey)}]`;
      const labelLiteral = expressionStringQuoter(seriesEntry.label);
      const colorLiteral = expressionStringQuoter(seriesColors[seriesIndex]);
      const valueText = chartValueTextExpressionBuilder(
        chartConfig,
        valueAccess,
      );
      return `{label: ${labelLiteral}, color: ${colorLiteral}, value: ${valueText}}`;
    });
    const titleExpression = chartXTitleExpressionBuilder(
      chartConfig,
      `${rowAccess}.chartX`,
    );
    return `{title: ${titleExpression}, row: ${rowAccess}, rows: [${tooltipRows.join(", ")}]}`;
  }

  function chartPointMarkerBuilder(dataName, markerSize, surfaceColor) {
    return {
      type: "symbol",
      from: { data: dataName },
      interactive: false,
      encode: {
        update: {
          x: { scale: "x", field: "chartX" },
          y: { scale: "y", field: "chartYTop" },
          size: { value: markerSize },
          fill: { scale: "color", field: "chartSeriesRowKey" },
          stroke: { value: surfaceColor },
          strokeWidth: { value: 2 },
        },
      },
    };
  }

  function chartHoverLayerBuilder(chartConfig, colorByToken, seriesColors) {
    const surfaceColor = chartTokenColorResolver(
      colorByToken,
      chartConfig.gapColor,
      chartFallbackGapColor,
    );
    const hoverMarks = [];
    if (!chartConfig.tooltip.isCrosshairHidden) {
      hoverMarks.push({
        type: "rule",
        from: { data: "hoveredRows" },
        interactive: false,
        encode: {
          update: {
            x: { scale: "x", field: "chartX" },
            y: { value: 0 },
            y2: { signal: "height" },
            stroke: {
              value: chartTokenColorResolver(
                colorByToken,
                chartConfig.tooltip.crosshairColor,
                chartFallbackCrosshairColor,
              ),
            },
            strokeWidth: { value: 1 },
          },
        },
      });
    }
    if (chartConfig.cartesian.markerDisplayMode === "always") {
      hoverMarks.push(chartPointMarkerBuilder("validPoints", 36, surfaceColor));
    }
    if (chartConfig.cartesian.markerDisplayMode !== "hidden") {
      hoverMarks.push(
        chartPointMarkerBuilder("hoveredPoints", 64, surfaceColor),
      );
    }
    hoverMarks.push(
      {
        name: "hoverAnchors",
        type: "symbol",
        from: { data: "rows" },
        interactive: false,
        encode: {
          update: {
            x: { scale: "x", field: "chartX" },
            y: { signal: "height / 2" },
            size: { value: 1 },
            fill: { value: "transparent" },
          },
        },
      },
      {
        name: "hoverCells",
        type: "path",
        from: { data: "hoverAnchors" },
        encode: {
          update: {
            fill: { value: "transparent" },
            tooltip: {
              signal: chartRowTooltipExpressionBuilder(
                chartConfig,
                "datum.datum",
                seriesColors,
              ),
            },
          },
        },
        transform: [
          {
            type: "voronoi",
            x: "datum.x",
            y: "datum.y",
            size: [{ signal: "width" }, { signal: "height" }],
          },
        ],
      },
    );
    return hoverMarks;
  }

  function chartLineMarksBuilder(chartConfig, colorByToken, seriesColors) {
    const pointPosition = {
      x: { scale: "x", field: "chartX" },
      y: { scale: "y", field: "chartYTop" },
      defined: { signal: "isValid(datum.chartValue)" },
    };
    return [
      {
        type: "group",
        from: {
          facet: {
            name: "seriesPoints",
            data: "points",
            groupby: "chartSeriesRowKey",
          },
        },
        marks: chartSeriesMarksBuilder(chartConfig, pointPosition),
      },
      ...chartHoverLayerBuilder(chartConfig, colorByToken, seriesColors),
    ];
  }

  function chartBarTooltipExpressionBuilder(chartConfig) {
    const titleExpression = chartXTitleExpressionBuilder(
      chartConfig,
      "datum.chartX",
    );
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    return `{title: ${titleExpression}, row: datum, rows: [{label: scale('seriesLabel', datum.chartSeriesRowKey), color: scale('color', datum.chartSeriesRowKey), value: ${valueText}}]}`;
  }

  function chartStackedBarMarksBuilder(chartConfig, colorByToken) {
    return [
      {
        type: "rect",
        from: { data: "validPoints" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartSeriesRowKey" },
            stroke: {
              value: chartTokenColorResolver(
                colorByToken,
                chartConfig.gapColor,
                chartFallbackGapColor,
              ),
            },
            strokeWidth: { value: 2 },
            cornerRadius: { value: 4 },
          },
          update: {
            xc: { scale: "x", field: "chartX", band: 0.5 },
            width: {
              signal: `min(bandwidth('x'), ${chartStackedBarMaxThicknessPx})`,
            },
            y: { scale: "y", field: "chartYTop" },
            y2: chartClampedBaselineResolver("datum.chartYBase"),
            fillOpacity: { value: 1 },
            tooltip: { signal: chartBarTooltipExpressionBuilder(chartConfig) },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
  }

  function chartGroupedBarMarksBuilder(chartConfig, seriesRowKeys) {
    return [
      {
        type: "group",
        from: {
          facet: { name: "barGroup", data: "validPoints", groupby: "chartX" },
        },
        encode: {
          update: {
            x: { scale: "x", field: "chartX" },
            width: { scale: "x", band: 1 },
          },
        },
        scales: [
          {
            name: "seriesOffset",
            type: "band",
            domain: seriesRowKeys,
            range: { signal: "[0, bandwidth('x')]" },
            paddingInner: 0.15,
          },
        ],
        marks: [
          {
            type: "rect",
            from: { data: "barGroup" },
            encode: {
              enter: {
                fill: { scale: "color", field: "chartSeriesRowKey" },
                cornerRadiusTopLeft: { value: 4 },
                cornerRadiusTopRight: { value: 4 },
              },
              update: {
                xc: {
                  scale: "seriesOffset",
                  field: "chartSeriesRowKey",
                  band: 0.5,
                },
                width: {
                  signal: `min(bandwidth('seriesOffset'), ${chartBarMaxThicknessPx})`,
                },
                y: { scale: "y", field: "chartYTop" },
                y2: chartClampedBaselineResolver("0"),
                fillOpacity: { value: 1 },
                tooltip: {
                  signal: chartBarTooltipExpressionBuilder(chartConfig),
                },
              },
              hover: { fillOpacity: { value: 0.8 } },
            },
          },
        ],
      },
    ];
  }

  function chartCartesianMarksBuilder(
    chartConfig,
    seriesRowKeys,
    colorByToken,
  ) {
    const seriesColors = chartSeriesColorsResolver(
      chartConfig.series,
      colorByToken,
    );
    if (chartConfig.chartType !== "bar") {
      return chartLineMarksBuilder(chartConfig, colorByToken, seriesColors);
    }
    if (chartConfig.cartesian.stackMode === "none") {
      return chartGroupedBarMarksBuilder(chartConfig, seriesRowKeys);
    }
    return chartStackedBarMarksBuilder(chartConfig, colorByToken);
  }

  function chartBaseSpecBuilder(fontFamily) {
    return {
      $schema: "https://vega.github.io/schema/vega/v6.json",
      width: 300,
      height: 200,
      padding: 4,
      autosize: { type: "fit", contains: "padding" },
      background: null,
      config: {
        text: { font: fontFamily },
        axis: {
          labelFont: fontFamily,
          labelFontSize: 11,
          titleFont: fontFamily,
          titleFontSize: 11,
          titleFontWeight: 600,
          gridWidth: 1,
          labelOverlap: true,
        },
      },
    };
  }

  function chartCartesianSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const seriesRowKeys = chartConfig.series.map(
      (seriesEntry) => seriesEntry.rowKey,
    );
    const chartSpec = chartBaseSpecBuilder(fontFamily);
    chartSpec.signals = [
      {
        name: "hoveredX",
        value: null,
        on: [
          { events: "@hoverCells:pointerover", update: "datum.datum.chartX" },
          { events: "@hoverCells:pointerout", update: "null" },
        ],
      },
    ];
    chartSpec.data = chartCartesianDataBuilder(chartConfig, seriesRowKeys);
    chartSpec.scales = chartCartesianScalesBuilder(
      chartConfig,
      seriesRowKeys,
      colorByToken,
    );
    chartSpec.axes = [
      chartXAxisBuilder(chartConfig, colorByToken),
      chartYAxisBuilder(chartConfig, colorByToken),
    ];
    chartSpec.marks = chartCartesianMarksBuilder(
      chartConfig,
      seriesRowKeys,
      colorByToken,
    );
    return chartSpec;
  }

  function chartDonutSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const valueSeries = chartConfig.series[0];
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    const donutTooltipExpression = `{title: datum.chartCategory, row: datum, rows: [{label: ${expressionStringQuoter(valueSeries.label)}, color: scale('color', datum.chartCategory), value: ${valueText}}]}`;
    const holeRatio = chartConfig.donut.holePercent / 100;
    const donutConfig = chartConfig.donut;
    const isCenterLabelShown = donutConfig.centerLabelDisplayMode !== "hidden";

    const chartSpec = chartBaseSpecBuilder(fontFamily);
    chartSpec.signals = [
      { name: "chartCategories", value: [] },
      { name: "chartCategoryColors", value: [] },
    ];
    if (isCenterLabelShown) {
      chartSpec.signals.push(
        { name: "centerLabelAlwaysSlice", value: "" },
        {
          name: "centerLabelSlice",
          value: "",
          on: [
            {
              events: "@donutArcs:pointerover",
              update: "datum.chartCategory",
            },
            {
              events: "@donutArcs:pointerout",
              update: "centerLabelAlwaysSlice",
            },
          ],
        },
      );
    }
    chartSpec.data = [
      {
        name: "rows",
        values: [],
        transform: [
          {
            type: "formula",
            as: "chartCategory",
            expr: `${chartRowKeyAccessBuilder(chartConfig.xAxis.rowKey)} + ''`,
          },
          {
            type: "formula",
            as: "chartValue",
            expr: `toNumber(${chartRowKeyAccessBuilder(valueSeries.rowKey)})`,
          },
          { type: "pie", field: "chartValue", sort: false },
        ],
      },
    ];
    chartSpec.scales = [
      {
        name: "color",
        type: "ordinal",
        domain: { signal: "chartCategories" },
        range: { signal: "chartCategoryColors" },
      },
    ];
    chartSpec.marks = [
      {
        name: "donutArcs",
        type: "arc",
        from: { data: "rows" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartCategory" },
            stroke: {
              value: chartTokenColorResolver(
                colorByToken,
                chartConfig.gapColor,
                chartFallbackGapColor,
              ),
            },
            strokeWidth: { value: 2 },
            cornerRadius: { value: 4 },
          },
          update: {
            x: { signal: "width / 2" },
            y: { signal: "height / 2" },
            startAngle: { field: "startAngle" },
            endAngle: { field: "endAngle" },
            innerRadius: { signal: `min(width, height) / 2 * ${holeRatio}` },
            outerRadius: { signal: "min(width, height) / 2" },
            fillOpacity: { value: 1 },
            tooltip: { signal: donutTooltipExpression },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
    if (isCenterLabelShown) {
      const centerLabelScale = `min(width, height) / 2 * ${holeRatio} * ${donutConfig.centerLabelSizeMultiplier}`;
      const valueSuffixLiteral = expressionStringQuoter(
        donutConfig.centerLabelValueSuffix,
      );
      const centerLabelValueText = chartValueTextExpressionBuilder(
        chartConfig,
        "datum.chartValue",
      );
      const centerLabelLimit = {
        signal: `min(width, height) * ${holeRatio} * 0.9`,
      };
      chartSpec.marks.push(
        {
          type: "text",
          from: { data: "rows" },
          interactive: false,
          encode: {
            update: {
              x: { signal: "width / 2" },
              y: { signal: "height / 2 - 3" },
              align: { value: "center" },
              baseline: { value: "bottom" },
              text: { signal: "datum.chartCategory" },
              fill: {
                value: chartTokenColorResolver(
                  colorByToken,
                  donutConfig.centerLabelSliceColor,
                  chartFallbackCenterLabelSliceColor,
                ),
              },
              fontSize: { signal: `${centerLabelScale} * 0.20` },
              limit: centerLabelLimit,
              opacity: {
                signal: "datum.chartCategory === centerLabelSlice ? 1 : 0",
              },
            },
          },
        },
        {
          type: "text",
          from: { data: "rows" },
          interactive: false,
          encode: {
            update: {
              x: { signal: "width / 2" },
              y: { signal: "height / 2 + 3" },
              align: { value: "center" },
              baseline: { value: "top" },
              text: {
                signal: `(${centerLabelValueText}) + ${valueSuffixLiteral}`,
              },
              fill: {
                value: chartTokenColorResolver(
                  colorByToken,
                  donutConfig.centerLabelValueColor,
                  chartFallbackCenterLabelValueColor,
                ),
              },
              fontSize: { signal: `${centerLabelScale} * 0.34` },
              limit: centerLabelLimit,
              opacity: {
                signal: "datum.chartCategory === centerLabelSlice ? 1 : 0",
              },
            },
          },
        },
      );
    }
    return chartSpec;
  }

  function chartDonutCategoriesReader(chartConfig, chartRows) {
    const donutCategories = new Set();
    for (const chartRow of chartRows) {
      donutCategories.add(String(chartRow[chartConfig.xAxis.rowKey]));
    }
    return Array.from(donutCategories);
  }

  function chartDonutAlwaysSliceResolver(chartConfig, donutCategories) {
    const alwaysSlice = chartConfig.donut.centerLabelAlwaysSlice;
    if (alwaysSlice === "") {
      return donutCategories.length > 0 ? donutCategories[0] : "";
    }
    if (donutCategories.includes(alwaysSlice)) {
      return alwaysSlice;
    }
    console.error(`ChartDonutCenterLabelSliceNotFound: ${alwaysSlice}`);
    return donutCategories.length > 0 ? donutCategories[0] : "";
  }

  function chartCategoryColorsResolver(
    chartConfig,
    donutCategories,
    colorByToken,
  ) {
    return donutCategories.map((donutCategory, categoryIndex) =>
      chartTokenColorResolver(
        colorByToken,
        chartConfig.donut.sliceIdentityColors[donutCategory],
        chartPaletteColorResolver(categoryIndex),
      ),
    );
  }

  function chartColorProbesReader(probeContainer) {
    const colorByToken = {};
    const unresolvedTokens = [];
    for (const probeElement of probeContainer.querySelectorAll(
      "[data-chart-color-token]",
    )) {
      const colorToken = probeElement.dataset.chartColorToken;
      const probedColor = getComputedStyle(probeElement).color;
      if (
        probedColor === "" ||
        probedColor === chartColorProbeUnresolvedColor
      ) {
        unresolvedTokens.push(colorToken);
        continue;
      }
      colorByToken[colorToken] = chartColorNormalizer(probedColor);
    }
    return { colorByToken, unresolvedTokens };
  }

  function chartRetryDelayWaiter() {
    return new Promise((resolveDelay) =>
      setTimeout(resolveDelay, chartColorProbeRetryDelayMs),
    );
  }

  async function chartColorsResolver(probeContainer) {
    let probeResult = chartColorProbesReader(probeContainer);
    for (
      let probeAttempt = 1;
      probeAttempt < chartColorProbeMaxAttempts &&
      probeResult.unresolvedTokens.length > 0;
      probeAttempt++
    ) {
      await chartRetryDelayWaiter();
      probeResult = chartColorProbesReader(probeContainer);
    }
    if (probeResult.unresolvedTokens.length > 0) {
      console.error(
        `ChartColorTokenUnresolved: ${probeResult.unresolvedTokens.join(", ")}`,
      );
    }
    return { ...chartFallbackColorByToken, ...probeResult.colorByToken };
  }

  const chartSpecBuilderByType = { donut: chartDonutSpecBuilder };

  Alpine.data("chart", (configScriptId) => {
    let chartView = null;
    let chartResizeObserver = null;
    let chartColorByToken = { ...chartFallbackColorByToken };
    let hasReportedPaletteExceeded = false;
    let hasBeenDestroyed = false;

    return {
      chartConfig: null,
      chartRows: [],
      chartLegendEntries: [],
      chartTooltip: {
        isVisible: false,
        isFlipped: false,
        isContentHtmlShown: false,
        contentHtml: "",
        title: "",
        rows: [],
      },

      async init() {
        try {
          const configElement = document.getElementById(configScriptId);
          this.chartConfig = JSON.parse(configElement?.textContent || "null");
        } catch (parseError) {
          console.error(`ChartInvalidConfigJson: ${parseError.message}`);
          return;
        }
        if (!this.chartConfig) {
          console.error(`ChartConfigNotFound: ${configScriptId}`);
          return;
        }
        let initialRows = this.chartConfig.rows ?? [];
        const rowsStatePath = this.chartConfig.rowsStatePath;
        if (rowsStatePath) {
          initialRows = Alpine.evaluate(this.$el, rowsStatePath);
          this.$watch(rowsStatePath, (watchedRows) =>
            this.replaceChartRows(watchedRows),
          );
        }
        this.chartRows = this.chartRowsTrimmer(initialRows);
        chartColorByToken = await chartColorsResolver(
          this.$refs.chartColorProbes,
        );
        if (hasBeenDestroyed) {
          return;
        }
        this.chartLegendEntries = this.chartLegendEntriesResolver();
        if (!window.vega) {
          console.error("ChartVegaLibraryNotLoaded");
          return;
        }

        const fontFamily = getComputedStyle(this.$el).fontFamily;
        const chartSpecBuilder =
          chartSpecBuilderByType[this.chartConfig.chartType] ??
          chartCartesianSpecBuilder;
        const chartSpec = chartSpecBuilder(
          this.chartConfig,
          fontFamily,
          chartColorByToken,
        );
        chartView = new vega.View(vega.parse(chartSpec), {
          renderer: "svg",
          container: this.$refs.chartCanvas,
          hover: true,
        });
        chartView.tooltip((_tooltipHandler, pointerEvent, _item, value) =>
          this.showChartTooltip(pointerEvent, value),
        );

        chartResizeObserver = new ResizeObserver(() => this.resizeChartView());
        chartResizeObserver.observe(this.$refs.chartCanvas);
        this.renderChartRows();
      },

      destroy() {
        hasBeenDestroyed = true;
        chartResizeObserver?.disconnect();
        chartView?.finalize();
      },

      chartRowsTrimmer(incomingRows) {
        const copiedRows = [];
        for (const chartRow of incomingRows ?? []) {
          copiedRows.push({ ...chartRow });
        }
        const retentionLimit = this.chartConfig?.maxRows ?? 0;
        if (retentionLimit === 0 || copiedRows.length <= retentionLimit) {
          return copiedRows;
        }
        return copiedRows.slice(copiedRows.length - retentionLimit);
      },

      reportPaletteExceeded(legendEntries) {
        const exceededEntries = legendEntries.filter(
          (legendEntry) => legendEntry.color === chartPaletteExceededColor,
        );
        if (exceededEntries.length === 0 || hasReportedPaletteExceeded) {
          return;
        }
        hasReportedPaletteExceeded = true;
        console.error(
          `ChartSeriesPaletteExceeded: ${chartSeriesPalette.length + exceededEntries.length}`,
        );
      },

      chartLegendEntriesResolver() {
        let legendLabels = this.chartConfig.series.map(
          (seriesEntry) => seriesEntry.label,
        );
        let legendColors = chartSeriesColorsResolver(
          this.chartConfig.series,
          chartColorByToken,
        );
        if (this.chartConfig.chartType === "donut") {
          legendLabels = chartDonutCategoriesReader(
            this.chartConfig,
            this.chartRows,
          );
          legendColors = chartCategoryColorsResolver(
            this.chartConfig,
            legendLabels,
            chartColorByToken,
          );
        }
        const legendEntries = legendLabels.map((legendLabel, entryIndex) => ({
          label: legendLabel,
          color: legendColors[entryIndex],
        }));
        this.reportPaletteExceeded(legendEntries);
        return legendEntries;
      },

      renderChartRows() {
        const vegaRows = this.chartRows.map((chartRow, rowIndex) => ({
          ...chartRow,
          chartRowOrder: rowIndex,
        }));
        if (this.chartConfig.chartType === "donut") {
          const donutCategories = chartDonutCategoriesReader(
            this.chartConfig,
            vegaRows,
          );
          chartView.signal("chartCategories", donutCategories);
          chartView.signal(
            "chartCategoryColors",
            chartCategoryColorsResolver(
              this.chartConfig,
              donutCategories,
              chartColorByToken,
            ),
          );
          if (this.chartConfig.donut.centerLabelDisplayMode !== "hidden") {
            let alwaysSlice = "";
            if (this.chartConfig.donut.centerLabelDisplayMode === "always") {
              alwaysSlice = chartDonutAlwaysSliceResolver(
                this.chartConfig,
                donutCategories,
              );
            }
            chartView.signal("centerLabelAlwaysSlice", alwaysSlice);
            chartView.signal("centerLabelSlice", alwaysSlice);
          }
        }
        const replacedRows = vega
          .changeset()
          .remove(() => true)
          .insert(vegaRows);
        chartView.change("rows", replacedRows);
        if (this.chartConfig.chartType !== "donut") {
          chartView.signal("hoveredX", null);
        }
        this.chartTooltip.isVisible = false;
        this.resizeChartView();
      },

      replaceChartRows(incomingRows) {
        this.chartRows = this.chartRowsTrimmer(incomingRows);
        this.chartLegendEntries = this.chartLegendEntriesResolver();
        if (!chartView) {
          return;
        }
        this.renderChartRows();
      },

      resizeChartView() {
        const chartCanvas = this.$refs.chartCanvas;
        chartView
          .width(chartCanvas.clientWidth)
          .height(chartCanvas.clientHeight)
          .runAsync()
          .catch((renderError) =>
            console.error(`ChartRenderFailed: ${renderError.message}`),
          );
      },

      showChartTooltip(pointerEvent, tooltipValue) {
        const isTooltipShown =
          tooltipValue && pointerEvent && !this.chartConfig.tooltip.isHidden;
        if (!isTooltipShown) {
          this.chartTooltip.isVisible = false;
          return;
        }
        const tooltipOffsetPx = 14;
        const flipThresholdPx = 240;
        const contentHtml = this.chartTooltipContentHtmlResolver(tooltipValue);
        this.chartTooltip = {
          isVisible: true,
          isFlipped: pointerEvent.clientX > window.innerWidth - flipThresholdPx,
          left: pointerEvent.clientX + tooltipOffsetPx,
          right: window.innerWidth - pointerEvent.clientX + tooltipOffsetPx,
          top: pointerEvent.clientY + tooltipOffsetPx,
          contentHtml,
          isContentHtmlShown: contentHtml !== "",
          title: String(tooltipValue.title ?? ""),
          rows: tooltipValue.rows ?? [],
        };
      },

      chartTooltipContentHtmlResolver(tooltipValue) {
        const contentFuncName = this.chartConfig.tooltip.contentFunc;
        if (contentFuncName === "") {
          return "";
        }
        const contentFunc = window[contentFuncName];
        if (typeof contentFunc !== "function") {
          console.error(`ChartTooltipContentFuncMissing: ${contentFuncName}`);
          return "";
        }
        return String(contentFunc(tooltipValue.row) ?? "");
      },

      chartTooltipStyleResolver() {
        const verticalPosition = `top: ${this.chartTooltip.top}px;`;
        if (this.chartTooltip.isFlipped) {
          return `${verticalPosition} right: ${this.chartTooltip.right}px;`;
        }
        return `${verticalPosition} left: ${this.chartTooltip.left}px;`;
      },
    };
  });
});

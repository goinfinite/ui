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
  const chartFallbackSliceLabelPillColor = "rgba(10, 10, 10, 0.7)";
  const chartPercentFormat = ".0%";
  const chartFallbackColorByToken = {
    "neutral-900": chartFallbackGapColor,
    "neutral-400": chartFallbackTickLabelColor,
    "neutral-300": chartFallbackTitleColor,
    "neutral-50/8": chartFallbackLineColor,
    "neutral-50/30": chartFallbackCrosshairColor,
    "neutral-950/70": chartFallbackSliceLabelPillColor,
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

  function chartIsBarType(chartConfig) {
    return (
      chartConfig.chartType === "bar" ||
      chartConfig.chartType === "horizontalBar"
    );
  }

  function chartIsCategoryColorType(chartConfig) {
    return (
      chartConfig.chartType === "donut" ||
      chartConfig.chartType === "parliament" ||
      chartConfig.chartType === "pie" ||
      chartConfig.chartType === "radialBar" ||
      chartConfig.chartType === "sunburst" ||
      chartConfig.chartType === "waffle"
    );
  }

  function chartIsDonutType(chartConfig) {
    return chartConfig.chartType === "donut" || chartConfig.chartType === "pie";
  }

  function chartIsHorizontalBar(chartConfig) {
    return chartConfig.chartType === "horizontalBar";
  }

  function chartCategoryScaleNameResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return "y";
    }
    return "x";
  }

  function chartValueScaleNameResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return "x";
    }
    return "y";
  }

  function chartCategoryScaleRangeResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return "height";
    }
    return "width";
  }

  function chartValueScaleRangeResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return "width";
    }
    return "height";
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

  function chartCategoryScaleBuilder(chartConfig) {
    const rowOrderedDomain = {
      data: "rows",
      field: "chartX",
      sort: { field: chartRowSortFieldResolver(chartConfig.xAxis), op: "min" },
    };
    if (chartIsBarType(chartConfig)) {
      return {
        name: chartCategoryScaleNameResolver(chartConfig),
        type: "band",
        domain: rowOrderedDomain,
        range: chartCategoryScaleRangeResolver(chartConfig),
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

  function chartValueScaleBuilder(chartConfig) {
    const valueScale = {
      name: chartValueScaleNameResolver(chartConfig),
      type: "linear",
      domain: { data: "points", fields: ["chartYBase", "chartYTop"] },
      range: chartValueScaleRangeResolver(chartConfig),
      zero: true,
      nice: true,
    };
    const yAxis = chartConfig.cartesian.yAxis;
    if (yAxis.minValue !== "") {
      valueScale.domainMin = Number(yAxis.minValue);
      valueScale.nice = false;
    }
    if (yAxis.maxValue !== "") {
      valueScale.domainMax = Number(yAxis.maxValue);
      valueScale.nice = false;
    }
    return valueScale;
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
      chartCategoryScaleBuilder(chartConfig),
      chartValueScaleBuilder(chartConfig),
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

  function chartCategoryAxisOrientResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return "left";
    }
    return "bottom";
  }

  function chartValueAxisOrientResolver(chartConfig) {
    if (!chartIsHorizontalBar(chartConfig)) {
      return chartConfig.cartesian.yAxis.position;
    }
    if (chartConfig.cartesian.yAxis.position === "right") {
      return "top";
    }
    return "bottom";
  }

  function chartCategoryAxisBuilder(chartConfig, colorByToken) {
    const xAxisSettings = chartConfig.xAxis;
    const categoryScaleName = chartCategoryScaleNameResolver(chartConfig);
    const categoryAxis = {
      orient: chartCategoryAxisOrientResolver(chartConfig),
      scale: categoryScaleName,
      tickCount: 6,
      grid: xAxisSettings.isGridLineVisible,
      ...chartAxisAppearanceBuilder(xAxisSettings, colorByToken),
    };
    const isContinuousScale = xAxisSettings.scaleType !== "category";
    const tickLabelFormat = xAxisSettings.tickLabelFormat;
    if (isContinuousScale && tickLabelFormat !== "") {
      categoryAxis.format = tickLabelFormat;
      categoryAxis.formatType = xAxisSettings.scaleType;
    }
    const bandLabelFormat = tickLabelFormat || chartConfig.tooltip.titleFormat;
    const isFormattedBand =
      chartIsBarType(chartConfig) &&
      isContinuousScale &&
      bandLabelFormat !== "";
    if (isFormattedBand) {
      delete categoryAxis.format;
      delete categoryAxis.formatType;
      categoryAxis.encode = chartBandLabelEncodingBuilder(
        xAxisSettings,
        bandLabelFormat,
      );
    }
    if (xAxisSettings.isTickStartEndOnly) {
      categoryAxis.values = {
        signal: `length(domain('${categoryScaleName}')) ? [domain('${categoryScaleName}')[0], peek(domain('${categoryScaleName}'))] : []`,
      };
    }
    return categoryAxis;
  }

  function chartValueAxisBuilder(chartConfig, colorByToken) {
    const yAxisSettings = chartConfig.cartesian.yAxis;
    return {
      orient: chartValueAxisOrientResolver(chartConfig),
      scale: chartValueScaleNameResolver(chartConfig),
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

  function chartClampedBaselineResolver(chartConfig, baselineValueExpression) {
    const valueScaleName = chartValueScaleNameResolver(chartConfig);
    return {
      signal: `scale('${valueScaleName}', clamp(${baselineValueExpression}, domain('${valueScaleName}')[0], domain('${valueScaleName}')[1]))`,
    };
  }

  function chartAreaBaselineResolver(chartConfig) {
    if (chartConfig.cartesian.stackMode === "none") {
      return chartClampedBaselineResolver(chartConfig, "0");
    }
    return chartClampedBaselineResolver(chartConfig, "datum.chartYBase");
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

  function chartStackedBarUpdateResolver(chartConfig) {
    const valueScaleName = chartValueScaleNameResolver(chartConfig);
    const baselineResolver = chartClampedBaselineResolver(
      chartConfig,
      "datum.chartYBase",
    );
    if (chartIsHorizontalBar(chartConfig)) {
      return {
        yc: { scale: "y", field: "chartX", band: 0.5 },
        height: {
          signal: `min(bandwidth('y'), ${chartStackedBarMaxThicknessPx})`,
        },
        x: { scale: valueScaleName, field: "chartYTop" },
        x2: baselineResolver,
      };
    }
    return {
      xc: { scale: "x", field: "chartX", band: 0.5 },
      width: {
        signal: `min(bandwidth('x'), ${chartStackedBarMaxThicknessPx})`,
      },
      y: { scale: valueScaleName, field: "chartYTop" },
      y2: baselineResolver,
    };
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
            ...chartStackedBarUpdateResolver(chartConfig),
            fillOpacity: { value: 1 },
            tooltip: { signal: chartBarTooltipExpressionBuilder(chartConfig) },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
  }

  function chartGroupedBarGroupEncodingResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return {
        y: { scale: "y", field: "chartX" },
        height: { scale: "y", band: 1 },
      };
    }
    return {
      x: { scale: "x", field: "chartX" },
      width: { scale: "x", band: 1 },
    };
  }

  function chartSeriesOffsetRangeResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return { signal: "[0, bandwidth('y')]" };
    }
    return { signal: "[0, bandwidth('x')]" };
  }

  function chartBarEndCornerRadiiResolver(chartConfig) {
    if (chartIsHorizontalBar(chartConfig)) {
      return {
        cornerRadiusTopRight: { value: 4 },
        cornerRadiusBottomRight: { value: 4 },
      };
    }
    return {
      cornerRadiusTopLeft: { value: 4 },
      cornerRadiusTopRight: { value: 4 },
    };
  }

  function chartGroupedBarUpdateResolver(chartConfig) {
    const valueScaleName = chartValueScaleNameResolver(chartConfig);
    const baselineResolver = chartClampedBaselineResolver(chartConfig, "0");
    if (chartIsHorizontalBar(chartConfig)) {
      return {
        yc: { scale: "seriesOffset", field: "chartSeriesRowKey", band: 0.5 },
        height: {
          signal: `min(bandwidth('seriesOffset'), ${chartBarMaxThicknessPx})`,
        },
        x: { scale: valueScaleName, field: "chartYTop" },
        x2: baselineResolver,
      };
    }
    return {
      xc: { scale: "seriesOffset", field: "chartSeriesRowKey", band: 0.5 },
      width: {
        signal: `min(bandwidth('seriesOffset'), ${chartBarMaxThicknessPx})`,
      },
      y: { scale: valueScaleName, field: "chartYTop" },
      y2: baselineResolver,
    };
  }

  function chartGroupedBarMarksBuilder(chartConfig, seriesRowKeys) {
    return [
      {
        type: "group",
        from: {
          facet: { name: "barGroup", data: "validPoints", groupby: "chartX" },
        },
        encode: {
          update: chartGroupedBarGroupEncodingResolver(chartConfig),
        },
        scales: [
          {
            name: "seriesOffset",
            type: "band",
            domain: seriesRowKeys,
            range: chartSeriesOffsetRangeResolver(chartConfig),
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
                ...chartBarEndCornerRadiiResolver(chartConfig),
              },
              update: {
                ...chartGroupedBarUpdateResolver(chartConfig),
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
    if (!chartIsBarType(chartConfig)) {
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
      chartCategoryAxisBuilder(chartConfig, colorByToken),
      chartValueAxisBuilder(chartConfig, colorByToken),
    ];
    chartSpec.marks = chartCartesianMarksBuilder(
      chartConfig,
      seriesRowKeys,
      colorByToken,
    );
    return chartSpec;
  }

  function chartCategoryTooltipExpressionBuilder(
    chartConfig,
    tooltipFieldExpressions,
  ) {
    const valueSeries = chartConfig.series[0];
    return `{title: ${tooltipFieldExpressions.title}, row: datum, rows: [{label: ${expressionStringQuoter(valueSeries.label)}, color: scale('color', ${tooltipFieldExpressions.color}), value: ${tooltipFieldExpressions.value}}]}`;
  }

  function chartCategoryBaseSpecBuilder(
    chartConfig,
    fontFamily,
    rowTransforms,
  ) {
    const valueSeries = chartConfig.series[0];
    const chartSpec = chartBaseSpecBuilder(fontFamily);
    chartSpec.signals = [
      { name: "chartCategories", value: [] },
      { name: "chartCategoryColors", value: [] },
    ];
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
          ...rowTransforms,
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
    return chartSpec;
  }

  function chartGapStrokeEncodingBuilder(chartConfig, colorByToken) {
    return {
      stroke: {
        value: chartTokenColorResolver(
          colorByToken,
          chartConfig.gapColor,
          chartFallbackGapColor,
        ),
      },
      strokeWidth: { value: 2 },
    };
  }

  function chartLabelTextWidthExpression(textExpression, fontSize) {
    const characterWidthRatio = 0.62;
    const horizontalPaddingPx = 12;
    return `length(toString(${textExpression})) * ${fontSize * characterWidthRatio} + ${horizontalPaddingPx}`;
  }

  function chartInFillLabelMarksBuilder(labelMarks) {
    const {
      dataName,
      xExpression,
      yExpression,
      textExpression,
      fontSize,
      textColor,
      pillColor,
      opacityExpression,
    } = labelMarks;
    const pillHeight = fontSize + 8;
    const widthExpression = chartLabelTextWidthExpression(
      textExpression,
      fontSize,
    );
    const pillEncode = {
      x: { signal: `(${xExpression}) - (${widthExpression}) / 2` },
      y: { signal: `(${yExpression}) - ${pillHeight} / 2` },
      width: { signal: widthExpression },
      height: { value: pillHeight },
      cornerRadius: { value: pillHeight / 2 },
      fill: { value: pillColor },
    };
    const textEncode = {
      x: { signal: xExpression },
      y: { signal: yExpression },
      align: { value: "center" },
      baseline: { value: "middle" },
      text: { signal: textExpression },
      fill: { value: textColor },
      fontSize: { value: fontSize },
      fontWeight: { value: "bold" },
    };
    if (opacityExpression) {
      pillEncode.opacity = { signal: opacityExpression };
      textEncode.opacity = { signal: opacityExpression };
    }
    return [
      {
        type: "rect",
        from: { data: dataName },
        interactive: false,
        encode: { update: pillEncode },
      },
      {
        type: "text",
        from: { data: dataName },
        interactive: false,
        encode: { update: textEncode },
      },
    ];
  }

  function chartSliceSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    const donutTooltipExpression = chartCategoryTooltipExpressionBuilder(
      chartConfig,
      {
        title: "datum.chartCategory",
        color: "datum.chartCategory",
        value: valueText,
      },
    );
    const holeRatio = chartConfig.slices.holePercent / 100;
    const slicesConfig = chartConfig.slices;
    const isCenterLabelShown = slicesConfig.centerLabelDisplayMode !== "hidden";

    const chartSpec = chartCategoryBaseSpecBuilder(chartConfig, fontFamily, [
      { type: "pie", field: "chartValue", sort: false },
    ]);
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
    chartSpec.marks = [
      {
        name: "donutArcs",
        type: "arc",
        from: { data: "rows" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartCategory" },
            ...chartGapStrokeEncodingBuilder(chartConfig, colorByToken),
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
    if (chartConfig.chartType === "pie" && !slicesConfig.isSliceLabelHidden) {
      const pieMidAngleExpression = "(datum.startAngle + datum.endAngle) / 2";
      const pieMidRadiusExpression = "min(width, height) / 2 * 0.62";
      const piePercentExpression = `format((datum.endAngle - datum.startAngle) / (2 * PI), '${chartPercentFormat}')`;
      chartSpec.marks.push(
        ...chartInFillLabelMarksBuilder({
          dataName: "rows",
          xExpression: `width / 2 + ${pieMidRadiusExpression} * sin(${pieMidAngleExpression})`,
          yExpression: `height / 2 - ${pieMidRadiusExpression} * cos(${pieMidAngleExpression})`,
          textExpression: piePercentExpression,
          fontSize: 12,
          textColor: chartTokenColorResolver(
            colorByToken,
            slicesConfig.sliceLabelColor,
            chartFallbackCenterLabelValueColor,
          ),
          pillColor: chartTokenColorResolver(
            colorByToken,
            slicesConfig.sliceLabelPillBackgroundColor,
            chartFallbackSliceLabelPillColor,
          ),
          opacityExpression: "datum.endAngle - datum.startAngle > 0.25 ? 1 : 0",
        }),
      );
    }
    if (isCenterLabelShown) {
      const centerLabelScale = `min(width, height) / 2 * ${holeRatio} * ${slicesConfig.centerLabelSizeMultiplier}`;
      const valueSuffixLiteral = expressionStringQuoter(
        slicesConfig.centerLabelValueSuffix,
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
                  slicesConfig.centerLabelSliceColor,
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
                  slicesConfig.centerLabelValueColor,
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

  function chartRadialBarSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    const radialBarTooltipExpression = chartCategoryTooltipExpressionBuilder(
      chartConfig,
      {
        title: "datum.chartCategory",
        color: "datum.chartCategory",
        value: valueText,
      },
    );
    const holeRatio = chartConfig.slices.holePercent / 100;
    const slicesConfig = chartConfig.slices;

    const chartSpec = chartCategoryBaseSpecBuilder(chartConfig, fontFamily, [
      {
        type: "joinaggregate",
        fields: ["chartValue"],
        ops: ["max"],
        as: ["chartValueMax"],
      },
    ]);
    chartSpec.scales.push(
      {
        name: "categoryTheta",
        type: "band",
        domain: { data: "rows", field: "chartCategory" },
        range: { signal: "[0, 2 * PI]" },
        paddingInner: 0.15,
      },
      {
        name: "valueRadius",
        type: "linear",
        domain: { data: "rows", field: "chartValue" },
        range: {
          signal: `[min(width, height) / 2 * ${holeRatio}, min(width, height) / 2]`,
        },
        zero: true,
      },
    );
    chartSpec.marks = [
      {
        name: "radialBars",
        type: "arc",
        from: { data: "rows" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartCategory" },
            ...chartGapStrokeEncodingBuilder(chartConfig, colorByToken),
            cornerRadius: { value: 4 },
          },
          update: {
            x: { signal: "width / 2" },
            y: { signal: "height / 2" },
            startAngle: { scale: "categoryTheta", field: "chartCategory" },
            endAngle: {
              scale: "categoryTheta",
              field: "chartCategory",
              band: 1,
            },
            innerRadius: { signal: `min(width, height) / 2 * ${holeRatio}` },
            outerRadius: { scale: "valueRadius", field: "chartValue" },
            fillOpacity: { value: 1 },
            tooltip: { signal: radialBarTooltipExpression },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
    if (!slicesConfig.isSliceLabelHidden) {
      const radialMidAngleExpression =
        "scale('categoryTheta', datum.chartCategory) + bandwidth('categoryTheta') / 2";
      const radialMidRadiusExpression = `(min(width, height) / 2 * ${holeRatio} + scale('valueRadius', datum.chartValue)) / 2`;
      const radialPercentExpression =
        "datum.chartValueMax > 0 ? datum.chartValue / datum.chartValueMax : 0";
      chartSpec.marks.push(
        ...chartInFillLabelMarksBuilder({
          dataName: "rows",
          xExpression: `width / 2 + ${radialMidRadiusExpression} * sin(${radialMidAngleExpression})`,
          yExpression: `height / 2 - ${radialMidRadiusExpression} * cos(${radialMidAngleExpression})`,
          textExpression: `format(${radialPercentExpression}, '${chartPercentFormat}')`,
          fontSize: 11,
          textColor: chartTokenColorResolver(
            colorByToken,
            slicesConfig.sliceLabelColor,
            chartFallbackCenterLabelValueColor,
          ),
          pillColor: chartTokenColorResolver(
            colorByToken,
            slicesConfig.sliceLabelPillBackgroundColor,
            chartFallbackSliceLabelPillColor,
          ),
        }),
      );
    }
    return chartSpec;
  }

  function chartWaffleSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    const waffleTooltipExpression = chartCategoryTooltipExpressionBuilder(
      chartConfig,
      {
        title: "datum.chartCategory",
        color: "datum.chartCategory",
        value: valueText,
      },
    );
    const waffleConfig = chartConfig.waffle;
    const waffleRowCount = Math.ceil(
      waffleConfig.cellCount / waffleConfig.columnCount,
    );
    const waffleCellInnerRatio = 1 - waffleConfig.cellGapPercent / 100;

    const chartSpec = chartCategoryBaseSpecBuilder(chartConfig, fontFamily, [
      {
        type: "joinaggregate",
        fields: ["chartValue"],
        ops: ["sum"],
        as: ["chartTotal"],
      },
      {
        type: "window",
        ops: ["sum"],
        fields: ["chartValue"],
        as: ["chartCumulative"],
      },
      {
        type: "formula",
        as: "chartCellEnd",
        expr: `datum.chartTotal > 0 ? round(datum.chartCumulative / datum.chartTotal * ${waffleConfig.cellCount}) : 0`,
      },
      {
        type: "window",
        ops: ["lag"],
        fields: ["chartCellEnd"],
        as: ["chartPreviousCellEnd"],
      },
      {
        type: "formula",
        as: "chartCellStart",
        expr: "isNumber(datum.chartPreviousCellEnd) ? datum.chartPreviousCellEnd : 0",
      },
      {
        type: "formula",
        as: "chartCells",
        expr: "sequence(datum.chartCellStart, datum.chartCellEnd)",
      },
      { type: "flatten", fields: ["chartCells"], as: ["chartCellIndex"] },
      {
        type: "formula",
        as: "chartCellColumn",
        expr: `datum.chartCellIndex % ${waffleConfig.columnCount}`,
      },
      {
        type: "formula",
        as: "chartCellRow",
        expr: `floor(datum.chartCellIndex / ${waffleConfig.columnCount})`,
      },
    ]);
    chartSpec.signals.push(
      {
        name: "waffleCellSize",
        update: `min(width / ${waffleConfig.columnCount}, height / ${waffleRowCount})`,
      },
      {
        name: "waffleCellInner",
        update: `waffleCellSize * ${waffleCellInnerRatio}`,
      },
      {
        name: "waffleCellInset",
        update: "(waffleCellSize - waffleCellInner) / 2",
      },
      {
        name: "waffleGridWidth",
        update: `waffleCellSize * ${waffleConfig.columnCount}`,
      },
      {
        name: "waffleGridHeight",
        update: `waffleCellSize * ${waffleRowCount}`,
      },
      { name: "waffleGridLeft", update: "(width - waffleGridWidth) / 2" },
      { name: "waffleGridTop", update: "(height - waffleGridHeight) / 2" },
    );
    const waffleCellCornerRadius =
      waffleConfig.cellShape === "circle"
        ? { signal: "waffleCellInner / 2" }
        : { value: 0 };
    chartSpec.marks = [
      {
        name: "waffleCells",
        type: "rect",
        from: { data: "rows" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartCategory" },
            ...chartGapStrokeEncodingBuilder(chartConfig, colorByToken),
          },
          update: {
            x: {
              signal:
                "waffleGridLeft + datum.chartCellColumn * waffleCellSize + waffleCellInset",
            },
            y: {
              signal:
                "waffleGridTop + datum.chartCellRow * waffleCellSize + waffleCellInset",
            },
            width: { signal: "waffleCellInner" },
            height: { signal: "waffleCellInner" },
            cornerRadius: waffleCellCornerRadius,
            fillOpacity: { value: 1 },
            tooltip: { signal: waffleTooltipExpression },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
    return chartSpec;
  }

  function chartSunburstSpecBuilder(chartConfig, fontFamily, colorByToken) {
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.depth === 1 ? datum.chartNodeTotal : datum.chartNodeValue",
    );
    const sunburstTooltipExpression = chartCategoryTooltipExpressionBuilder(
      chartConfig,
      {
        title:
          "datum.chartNodeChild ? datum.chartNodeChild : datum.chartNodeKey",
        color: "datum.chartNodeCategory",
        value: valueText,
      },
    );
    const childText = `(${chartRowKeyAccessBuilder(chartConfig.sunburst.childRowKey)} == null ? '' : ${chartRowKeyAccessBuilder(chartConfig.sunburst.childRowKey)} + '')`;
    const rootKeyLiteral = expressionStringQuoter("chart-sunburst-root");

    const chartSpec = chartCategoryBaseSpecBuilder(chartConfig, fontFamily, [
      {
        type: "formula",
        as: "chartNodeList",
        expr:
          `[{chartNodeKey: ${rootKeyLiteral}, chartNodeParent: '', chartNodeChild: '', chartNodeValue: 0, chartNodeIsLeaf: 0}, ` +
          `{chartNodeKey: datum.chartCategory, chartNodeParent: ${rootKeyLiteral}, chartNodeChild: '', chartNodeValue: 0, chartNodeIsLeaf: 0}, ` +
          `${childText} === '' ` +
          `? {chartNodeKey: datum.chartCategory, chartNodeParent: ${rootKeyLiteral}, chartNodeChild: '', chartNodeValue: datum.chartValue, chartNodeIsLeaf: 0} ` +
          `: {chartNodeKey: datum.chartCategory.length + ':' + datum.chartCategory + (${childText}).length + ':' + ${childText}, chartNodeParent: datum.chartCategory, chartNodeChild: ${childText}, chartNodeValue: datum.chartValue, chartNodeIsLeaf: 1}]`,
      },
      { type: "flatten", fields: ["chartNodeList"], as: ["chartNode"] },
      {
        type: "formula",
        as: "chartNodeKey",
        expr: "datum.chartNode.chartNodeKey",
      },
      {
        type: "formula",
        as: "chartNodeParent",
        expr: "datum.chartNode.chartNodeParent",
      },
      {
        type: "formula",
        as: "chartNodeChild",
        expr: "datum.chartNode.chartNodeChild",
      },
      {
        type: "formula",
        as: "chartNodeValue",
        expr: "datum.chartNode.chartNodeValue",
      },
      {
        type: "formula",
        as: "chartNodeIsLeaf",
        expr: "datum.chartNode.chartNodeIsLeaf",
      },
      {
        type: "aggregate",
        groupby: ["chartNodeKey"],
        fields: [
          "chartNodeParent",
          "chartNodeChild",
          "chartNodeValue",
          "chartNodeIsLeaf",
        ],
        ops: ["min", "min", "sum", "max"],
        as: [
          "chartNodeParent",
          "chartNodeChild",
          "chartNodeValue",
          "chartNodeIsLeaf",
        ],
      },
      {
        type: "formula",
        as: "chartNodeCategory",
        expr: "datum.chartNodeIsLeaf > 0 ? datum.chartNodeParent : datum.chartNodeKey",
      },
      {
        type: "window",
        groupby: ["chartNodeCategory"],
        ops: ["sum"],
        fields: ["chartNodeValue"],
        frame: [null, null],
        as: ["chartNodeTotal"],
      },
      { type: "stratify", key: "chartNodeKey", parentKey: "chartNodeParent" },
      {
        type: "partition",
        field: "chartNodeValue",
        size: [{ signal: "2 * PI" }, { signal: "sunburstRadius" }],
      },
    ]);
    chartSpec.signals.push(
      {
        name: "sunburstRadius",
        update: "min(width, height) / 2 - 4",
      },
      { name: "sunburstTotal", value: 0 },
    );
    chartSpec.data.push({
      name: "arcs",
      source: "rows",
      transform: [{ type: "filter", expr: "datum.depth > 0" }],
    });
    chartSpec.marks = [
      {
        name: "sunburstArcs",
        type: "arc",
        from: { data: "arcs" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartNodeCategory" },
            ...chartGapStrokeEncodingBuilder(chartConfig, colorByToken),
            cornerRadius: { value: 4 },
          },
          update: {
            x: { signal: "width / 2" },
            y: { signal: "height / 2" },
            startAngle: { field: "x0" },
            endAngle: { field: "x1" },
            innerRadius: { field: "y0" },
            outerRadius: { field: "y1" },
            fillOpacity: { signal: "datum.depth === 1 ? 1 : 0.6" },
            tooltip: { signal: sunburstTooltipExpression },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
    if (!chartConfig.slices.isSliceLabelHidden) {
      const sunburstPercentExpression =
        "datum.depth === 1 ? (sunburstTotal > 0 ? datum.chartNodeTotal / sunburstTotal : 0) : (datum.chartNodeTotal > 0 ? datum.chartNodeValue / datum.chartNodeTotal : 0)";
      chartSpec.marks.push(
        ...chartInFillLabelMarksBuilder({
          dataName: "arcs",
          xExpression:
            "width / 2 + ((datum.y0 + datum.y1) / 2) * sin((datum.x0 + datum.x1) / 2)",
          yExpression:
            "height / 2 - ((datum.y0 + datum.y1) / 2) * cos((datum.x0 + datum.x1) / 2)",
          textExpression: `format(${sunburstPercentExpression}, '${chartPercentFormat}')`,
          fontSize: 10,
          textColor: chartTokenColorResolver(
            colorByToken,
            chartConfig.slices.sliceLabelColor,
            chartFallbackCenterLabelValueColor,
          ),
          pillColor: chartTokenColorResolver(
            colorByToken,
            chartConfig.slices.sliceLabelPillBackgroundColor,
            chartFallbackSliceLabelPillColor,
          ),
          opacityExpression: "datum.x1 - datum.x0 > 0.25 ? 1 : 0",
        }),
      );
    }
    return chartSpec;
  }

  function chartParliamentSpecBuilder(chartConfig, fontFamily, _colorByToken) {
    const valueText = chartValueTextExpressionBuilder(
      chartConfig,
      "datum.chartValue",
    );
    const parliamentTooltipExpression = chartCategoryTooltipExpressionBuilder(
      chartConfig,
      {
        title: "datum.chartCategory",
        color: "datum.chartCategory",
        value: valueText,
      },
    );
    const rowCount = chartConfig.parliament.rowCount;
    const parliamentDotDiameterRatio =
      (1 - chartConfig.parliament.dotGapPercent / 100) *
      (chartConfig.parliament.dotSizePercent / 100);
    const seatRowBoundaryExpression = (rowIndex) =>
      `round(datum.chartSeatsTotal * ${(rowIndex + 1) * (rowIndex + 1)} / ${rowCount * rowCount})`;
    const seatRowExpression = [
      ...[...Array(rowCount - 1).keys()].map(
        (rowIndex) =>
          `datum.chartSeatOrdinal < ${seatRowBoundaryExpression(rowIndex)} ? ${rowIndex} :`,
      ),
      `${rowCount - 1}`,
    ].join(" ");
    const seatRowStartExpression = [
      ...[...Array(rowCount - 1).keys()].map(
        (rowIndex) =>
          `datum.chartSeatOrdinal < ${seatRowBoundaryExpression(rowIndex)} ? ${rowIndex === 0 ? "0" : seatRowBoundaryExpression(rowIndex - 1)} :`,
      ),
      seatRowBoundaryExpression(rowCount - 2),
    ].join(" ");
    const seatRowEndExpression = [
      ...[...Array(rowCount - 1).keys()].map(
        (rowIndex) =>
          `datum.chartSeatRow === ${rowIndex} ? ${seatRowBoundaryExpression(rowIndex)} :`,
      ),
      "datum.chartSeatsTotal",
    ].join(" ");

    const chartSpec = chartCategoryBaseSpecBuilder(chartConfig, fontFamily, [
      {
        type: "joinaggregate",
        fields: ["chartValue"],
        ops: ["sum"],
        as: ["chartTotal"],
      },
      {
        type: "formula",
        as: "chartSeatsTotal",
        expr: "max(1, round(datum.chartTotal))",
      },
      {
        type: "window",
        ops: ["sum"],
        fields: ["chartValue"],
        as: ["chartCumulative"],
      },
      {
        type: "formula",
        as: "chartSeatEnd",
        expr: "datum.chartTotal > 0 ? round(datum.chartCumulative / datum.chartTotal * datum.chartSeatsTotal) : 0",
      },
      {
        type: "window",
        ops: ["lag"],
        fields: ["chartSeatEnd"],
        as: ["chartPreviousSeatEnd"],
      },
      {
        type: "formula",
        as: "chartSeatStart",
        expr: "isNumber(datum.chartPreviousSeatEnd) ? datum.chartPreviousSeatEnd : 0",
      },
      {
        type: "formula",
        as: "chartSeatOrdinals",
        expr: "sequence(datum.chartSeatStart, datum.chartSeatEnd)",
      },
      {
        type: "flatten",
        fields: ["chartSeatOrdinals"],
        as: ["chartSeatOrdinal"],
      },
      { type: "formula", as: "chartSeatRow", expr: seatRowExpression },
      {
        type: "formula",
        as: "chartSeatRowStart",
        expr: seatRowStartExpression,
      },
      { type: "formula", as: "chartSeatRowEnd", expr: seatRowEndExpression },
      {
        type: "formula",
        as: "chartSeatAngle",
        expr: "PI * (datum.chartSeatOrdinal - datum.chartSeatRowStart + 0.5) / (datum.chartSeatRowEnd - datum.chartSeatRowStart)",
      },
      {
        type: "formula",
        as: "chartSeatRadiusFraction",
        expr: `(2 * datum.chartSeatRow + 1) / ${2 * rowCount}`,
      },
      {
        type: "formula",
        as: "chartSeatRowSeats",
        expr: "max(1, datum.chartSeatRowEnd - datum.chartSeatRowStart)",
      },
      {
        type: "formula",
        as: "chartSeatSpacing",
        expr: `min(parliamentRadius / ${rowCount}, datum.chartSeatRadiusFraction * parliamentRadius * PI / datum.chartSeatRowSeats)`,
      },
      {
        type: "formula",
        as: "chartSeatDotArea",
        expr: `PI / 4 * pow(datum.chartSeatSpacing * ${parliamentDotDiameterRatio}, 2)`,
      },
    ]);
    const parliamentDotShape = chartConfig.parliament.dotShape;
    chartSpec.signals.push(
      { name: "parliamentCenterX", update: "width / 2" },
      { name: "parliamentCenterY", update: "height - 4" },
      { name: "parliamentRadius", update: "min(width / 2, height - 8) - 4" },
    );
    chartSpec.marks = [
      {
        name: "parliamentSeats",
        type: "symbol",
        from: { data: "rows" },
        encode: {
          enter: {
            fill: { scale: "color", field: "chartCategory" },
            shape: { value: parliamentDotShape },
          },
          update: {
            x: {
              signal:
                "parliamentCenterX - parliamentRadius * datum.chartSeatRadiusFraction * cos(datum.chartSeatAngle)",
            },
            y: {
              signal:
                "parliamentCenterY - parliamentRadius * datum.chartSeatRadiusFraction * sin(datum.chartSeatAngle)",
            },
            size: {
              signal: "datum.chartSeatDotArea",
            },
            fillOpacity: { value: 1 },
            tooltip: { signal: parliamentTooltipExpression },
          },
          hover: { fillOpacity: { value: 0.8 } },
        },
      },
    ];
    return chartSpec;
  }

  function chartCategoryNamesReader(chartConfig, chartRows) {
    const categoryNames = new Set();
    for (const chartRow of chartRows) {
      categoryNames.add(String(chartRow[chartConfig.xAxis.rowKey]));
    }
    return Array.from(categoryNames);
  }

  function chartCenterLabelAlwaysSliceResolver(chartConfig, sliceNames) {
    const alwaysSlice = chartConfig.slices.centerLabelAlwaysSlice;
    if (alwaysSlice === "") {
      return sliceNames.length > 0 ? sliceNames[0] : "";
    }
    if (sliceNames.includes(alwaysSlice)) {
      return alwaysSlice;
    }
    console.error(`ChartCenterLabelSliceNotFound: ${alwaysSlice}`);
    return sliceNames.length > 0 ? sliceNames[0] : "";
  }

  function chartCategoryColorsResolver(
    chartConfig,
    categoryNames,
    colorByToken,
  ) {
    return categoryNames.map((categoryName, categoryIndex) =>
      chartTokenColorResolver(
        colorByToken,
        chartConfig.slices.sliceIdentityColors[categoryName],
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

  const chartSpecBuilderByType = {
    donut: chartSliceSpecBuilder,
    horizontalBar: chartCartesianSpecBuilder,
    parliament: chartParliamentSpecBuilder,
    pie: chartSliceSpecBuilder,
    radialBar: chartRadialBarSpecBuilder,
    sunburst: chartSunburstSpecBuilder,
    waffle: chartWaffleSpecBuilder,
  };

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
        if (chartIsCategoryColorType(this.chartConfig)) {
          legendLabels = chartCategoryNamesReader(
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

      chartSunburstTotalResolver() {
        const valueSeriesRowKey = this.chartConfig.series[0].rowKey;
        let totalValue = 0;
        for (const chartRow of this.chartRows) {
          const rowValue = Number(chartRow[valueSeriesRowKey]);
          if (!Number.isNaN(rowValue)) {
            totalValue += rowValue;
          }
        }
        return totalValue;
      },

      renderChartRows() {
        const vegaRows = this.chartRows.map((chartRow, rowIndex) => ({
          ...chartRow,
          chartRowOrder: rowIndex,
        }));
        if (chartIsCategoryColorType(this.chartConfig)) {
          const categoryNames = chartCategoryNamesReader(
            this.chartConfig,
            vegaRows,
          );
          chartView.signal("chartCategories", categoryNames);
          chartView.signal(
            "chartCategoryColors",
            chartCategoryColorsResolver(
              this.chartConfig,
              categoryNames,
              chartColorByToken,
            ),
          );
          if (
            this.chartConfig.slices.centerLabelDisplayMode !== "hidden" &&
            chartIsDonutType(this.chartConfig)
          ) {
            let alwaysSlice = "";
            if (this.chartConfig.slices.centerLabelDisplayMode === "always") {
              alwaysSlice = chartCenterLabelAlwaysSliceResolver(
                this.chartConfig,
                categoryNames,
              );
            }
            chartView.signal("centerLabelAlwaysSlice", alwaysSlice);
            chartView.signal("centerLabelSlice", alwaysSlice);
          }
          if (this.chartConfig.chartType === "sunburst") {
            chartView.signal(
              "sunburstTotal",
              this.chartSunburstTotalResolver(),
            );
          }
        }
        const replacedRows = vega
          .changeset()
          .remove(() => true)
          .insert(vegaRows);
        chartView.change("rows", replacedRows);
        if (!chartIsCategoryColorType(this.chartConfig)) {
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

  Alpine.data(
    "chartProgress",
    (
      initialValue,
      minValue,
      maxValue,
      thresholds,
      defaultColor,
      labelPrefix,
      labelSuffix,
    ) => ({
      progressValue: Number(initialValue) || 0,
      minValue: Number(minValue) || 0,
      maxValue: Number.isFinite(Number(maxValue)) ? Number(maxValue) : 100,
      thresholds: Array.isArray(thresholds) ? thresholds : [],
      defaultColor,
      labelPrefix: typeof labelPrefix === "string" ? labelPrefix : "",
      labelSuffix: typeof labelSuffix === "string" ? labelSuffix : "",

      get percent() {
        const valueSpan = this.maxValue - this.minValue;
        if (!(valueSpan > 0)) {
          return 0;
        }
        const ratio = (this.progressValue - this.minValue) / valueSpan;
        return Math.min(Math.max(ratio, 0), 1) * 100;
      },

      get valueAriaText() {
        const clampedValue = Math.min(
          Math.max(this.progressValue, this.minValue),
          this.maxValue,
        );
        return (Math.round(clampedValue * 100) / 100).toFixed(2);
      },

      get percentWidth() {
        return `${Math.round(this.percent * 10) / 10}%`;
      },

      get percentLabel() {
        return `${this.labelPrefix}${this.percentWidth}${this.labelSuffix}`;
      },

      get fillColorClass() {
        for (const threshold of this.thresholds) {
          if (this.percent >= threshold.minimumPercent) {
            return `bg-${threshold.backgroundColor}`;
          }
        }
        return `bg-${this.defaultColor}`;
      },

      syncValueFromPath(boundValue) {
        const numericValue = Number(boundValue);
        this.progressValue = Number.isNaN(numericValue) ? 0 : numericValue;
      },
    }),
  );
});

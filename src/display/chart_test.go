package uiDisplay

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"
)

func TestChartConfigResolver(t *testing.T) {
	cpuSeries := []ChartSeriesSettings{{RowKey: "cpu"}}
	timeAxis := ChartXAxisSettings{RowKey: "time"}
	testCases := []struct {
		name                     string
		settings                 ChartSettings
		expectedErrPrefix        string
		expectedChartType        ChartType
		expectedSeriesLabel      string
		expectedYTickLabelFormat string
		expectedTooltipFormats   [2]string
	}{
		{
			name:                     "DefaultsFillTypeScaleAndLabel",
			settings:                 ChartSettings{XAxis: timeAxis, Series: cpuSeries},
			expectedChartType:        ChartTypeLine,
			expectedSeriesLabel:      "cpu",
			expectedYTickLabelFormat: chartDefaultValueFormat,
			expectedTooltipFormats:   [2]string{"", chartDefaultValueFormat},
		},
		{
			name: "TimeScaleDefaultsTooltipTitleFormat",
			settings: ChartSettings{
				Type:   ChartTypeArea,
				XAxis:  ChartXAxisSettings{RowKey: "time", ScaleType: ChartScaleTypeTime},
				Series: []ChartSeriesSettings{{RowKey: "cpu", Label: "CPU"}},
				Cartesian: ChartCartesianSettings{
					YAxis: ChartYAxisSettings{TickLabelFormat: ".0%", MinValue: "0", MaxValue: "1"},
				},
			},
			expectedChartType:        ChartTypeArea,
			expectedSeriesLabel:      "CPU",
			expectedYTickLabelFormat: ".0%",
			expectedTooltipFormats:   [2]string{chartDefaultTimeTitleFormat, ".0%"},
		},
		{
			name: "PercentStackKeepsRawTooltipValues",
			settings: ChartSettings{
				Type: ChartTypeBar, XAxis: timeAxis, Series: cpuSeries,
				Cartesian: ChartCartesianSettings{StackMode: ChartStackModePercent},
			},
			expectedChartType:        ChartTypeBar,
			expectedSeriesLabel:      "cpu",
			expectedYTickLabelFormat: chartDefaultPercentFormat,
			expectedTooltipFormats:   [2]string{"", chartDefaultValueFormat},
		},
		{
			name:              "MissingSeriesFails",
			settings:          ChartSettings{XAxis: timeAxis},
			expectedErrPrefix: "ChartSeriesMissing",
		},
		{
			name:              "MissingSeriesRowKeyFails",
			settings:          ChartSettings{XAxis: timeAxis, Series: []ChartSeriesSettings{{Label: "CPU"}}},
			expectedErrPrefix: "ChartSeriesRowKeyMissing",
		},
		{
			name:              "MissingXAxisRowKeyFails",
			settings:          ChartSettings{Series: cpuSeries},
			expectedErrPrefix: "ChartXAxisRowKeyMissing",
		},
		{
			name:              "UnknownTypeFails",
			settings:          ChartSettings{Type: "pie", XAxis: timeAxis, Series: cpuSeries},
			expectedErrPrefix: "ChartInvalidType",
		},
		{
			name: "UnknownStackModeFails",
			settings: ChartSettings{
				XAxis: timeAxis, Series: cpuSeries,
				Cartesian: ChartCartesianSettings{StackMode: "stream"},
			},
			expectedErrPrefix: "ChartInvalidStackMode",
		},
		{
			name: "UnknownLineDashStyleFails",
			settings: ChartSettings{
				XAxis:  timeAxis,
				Series: []ChartSeriesSettings{{RowKey: "cpu", LineDashStyle: "wavy"}},
			},
			expectedErrPrefix: "ChartInvalidSeriesLineDashStyle",
		},
		{
			name: "InvalidYValueBoundFails",
			settings: ChartSettings{
				XAxis: timeAxis, Series: cpuSeries,
				Cartesian: ChartCartesianSettings{
					YAxis: ChartYAxisSettings{MinValue: "zero"},
				},
			},
			expectedErrPrefix: "ChartInvalidYAxisValueBound",
		},
		{
			name: "OversizedDonutHoleFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: timeAxis, Series: cpuSeries,
				Donut: ChartDonutSettings{HolePercent: 95},
			},
			expectedErrPrefix: "ChartInvalidDonutHolePercent",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resolvedConfig, resolveErr := testCase.settings.configResolver()
			if testCase.expectedErrPrefix != "" {
				if resolveErr == nil || !strings.HasPrefix(resolveErr.Error(), testCase.expectedErrPrefix) {
					t.Fatalf("ExpectedError %q, got %v", testCase.expectedErrPrefix, resolveErr)
				}
				return
			}
			if resolveErr != nil {
				t.Fatalf("UnexpectedChartConfigError: %v", resolveErr)
			}
			if resolvedConfig.ChartType != testCase.expectedChartType {
				t.Errorf("ChartTypeMismatch: want %q, got %q", testCase.expectedChartType, resolvedConfig.ChartType)
			}
			if resolvedConfig.Series[0].Label != testCase.expectedSeriesLabel {
				t.Errorf("SeriesLabelMismatch: want %q, got %q", testCase.expectedSeriesLabel, resolvedConfig.Series[0].Label)
			}
			if resolvedConfig.Series[0].LineDashStyle != ChartLineDashStyleSolid {
				t.Errorf("LineDashStyleDefaultMismatch: got %q", resolvedConfig.Series[0].LineDashStyle)
			}
			if resolvedConfig.Cartesian.YAxis.TickLabelFormat != testCase.expectedYTickLabelFormat {
				t.Errorf("YTickLabelFormatMismatch: want %q, got %q", testCase.expectedYTickLabelFormat, resolvedConfig.Cartesian.YAxis.TickLabelFormat)
			}
			resolvedTooltipFormats := [2]string{
				resolvedConfig.Tooltip.TitleFormat, resolvedConfig.Tooltip.ValueFormat,
			}
			if resolvedTooltipFormats != testCase.expectedTooltipFormats {
				t.Errorf("TooltipFormatsMismatch: want %q, got %q", testCase.expectedTooltipFormats, resolvedTooltipFormats)
			}
		})
	}
}

func TestChartColorTokensCollector(t *testing.T) {
	resolvedConfig, resolveErr := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time", TickLabelColor: "sky-300"},
		Series: []ChartSeriesSettings{
			{RowKey: "cpu", IdentityColor: "sky-300"},
			{RowKey: "memory"},
		},
		Donut: ChartDonutSettings{SliceIdentityColors: map[string]string{"nginx": "[#4c8dff]"}},
	}.configResolver()
	if resolveErr != nil {
		t.Fatalf("UnexpectedChartConfigError: %v", resolveErr)
	}
	colorTokens := resolvedConfig.colorTokensCollector()
	for _, expectedToken := range []string{"sky-300", "[#4c8dff]", "neutral-900", "neutral-50/30"} {
		if !slices.Contains(colorTokens, expectedToken) {
			t.Errorf("ColorTokenMissing: %q", expectedToken)
		}
	}
	if slices.Contains(colorTokens, "") {
		t.Error("EmptyColorTokenCollected")
	}
	skyTokenCount := 0
	for _, colorToken := range colorTokens {
		if colorToken == "sky-300" {
			skyTokenCount++
		}
	}
	if skyTokenCount != 1 {
		t.Errorf("DuplicateColorToken: sky-300 appears %d times", skyTokenCount)
	}
}

func TestChartRendersParts(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type:  ChartTypeBar,
		XAxis: ChartXAxisSettings{RowKey: "month", Title: "Month"},
		Series: []ChartSeriesSettings{
			{RowKey: "received", Label: "Received", IdentityColor: "sky-500"},
			{RowKey: "transmitted", Label: "Transmitted"},
		},
		Legend: ChartLegendSettings{
			Position:    ChartLegendPositionBottom,
			Alignment:   ChartLegendAlignmentEnd,
			DisplayMode: ChartLegendDisplayModeAlways,
		},
		Tooltip: ChartTooltipSettings{BackgroundColor: "sky-950"},
		Rows:    []map[string]any{{"month": "Jan", "received": 1, "transmitted": 2}},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`<th scope="col">Month</th>`,
		`<th scope="col">Received</th>`,
		`aria-label="Received, Transmitted"`,
		`x-text="chartRow[&#34;transmitted&#34;]"`,
		`"month":"Jan"`,
		`data-chart-color-token="sky-500"`,
		`order-last`,
		`justify-end`,
		`x-show="chartLegendEntries.length &gt; 0"`,
		`bg-sky-950`,
		`No data`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartHidesLegend(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		XAxis:  ChartXAxisSettings{RowKey: "month"},
		Series: []ChartSeriesSettings{{RowKey: "received"}, {RowKey: "transmitted"}},
		Legend: ChartLegendSettings{DisplayMode: ChartLegendDisplayModeHidden},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	if strings.Contains(buffer.String(), "legendEntry") {
		t.Error("HiddenLegendRendered")
	}
}

func TestChartRenderFailsOnInvalidSettings(t *testing.T) {
	for _, invalidSettings := range []ChartSettings{
		{XAxis: ChartXAxisSettings{RowKey: "month"}},
		{
			XAxis:  ChartXAxisSettings{RowKey: "month"},
			Series: []ChartSeriesSettings{{RowKey: "received"}},
			Legend: ChartLegendSettings{Position: "left"},
		},
	} {
		renderErr := Chart(invalidSettings).Render(context.Background(), &bytes.Buffer{})
		if renderErr == nil {
			t.Errorf("ExpectedRenderError: %+v", invalidSettings)
		}
	}
}

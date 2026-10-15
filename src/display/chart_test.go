package uiDisplay

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
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
				Type: ChartTypeBar, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
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
				XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
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
				XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Cartesian: ChartCartesianSettings{
					YAxis: ChartYAxisSettings{MinValue: "zero"},
				},
			},
			expectedErrPrefix: "ChartInvalidYAxisValueBound",
		},
		{
			name: "OversizedDonutHoleFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Donut: ChartDonutSettings{HolePercent: 95},
			},
			expectedErrPrefix: "ChartInvalidDonutHolePercent",
		},
		{
			name: "UnknownCenterLabelDisplayModeFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Donut: ChartDonutSettings{CenterLabelDisplayMode: "sometimes"},
			},
			expectedErrPrefix: "ChartInvalidDonutCenterLabelDisplayMode",
		},
		{
			name: "UnknownCenterLabelSizeFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Donut: ChartDonutSettings{CenterLabelSize: "huge"},
			},
			expectedErrPrefix: "ChartInvalidDonutCenterLabelSize",
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

func TestChartYTickCountDefaults(t *testing.T) {
	defaultedConfig, defaultedErr := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
	}.configResolver()
	if defaultedErr != nil {
		t.Fatalf("UnexpectedChartConfigError: %v", defaultedErr)
	}
	if defaultedConfig.Cartesian.YAxis.TickCount != chartDefaultYTickCount {
		t.Errorf(
			"DefaultTickCountMismatch: want %d, got %d",
			chartDefaultYTickCount, defaultedConfig.Cartesian.YAxis.TickCount,
		)
	}

	providedConfig, providedErr := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Cartesian: ChartCartesianSettings{YAxis: ChartYAxisSettings{TickCount: 8}},
	}.configResolver()
	if providedErr != nil {
		t.Fatalf("UnexpectedChartConfigError: %v", providedErr)
	}
	if providedConfig.Cartesian.YAxis.TickCount != 8 {
		t.Errorf("ProvidedTickCountMismatch: got %d", providedConfig.Cartesian.YAxis.TickCount)
	}
}

func TestChartDonutHolePercentBoundary(t *testing.T) {
	atMaxSettings := ChartSettings{
		Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Donut: ChartDonutSettings{HolePercent: chartMaxDonutHolePercent},
	}
	if _, resolveErr := atMaxSettings.configResolver(); resolveErr != nil {
		t.Errorf("AtMaxHolePercentShouldPass: %v", resolveErr)
	}

	overMaxSettings := ChartSettings{
		Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Donut: ChartDonutSettings{HolePercent: chartMaxDonutHolePercent + 1},
	}
	_, overMaxErr := overMaxSettings.configResolver()
	if overMaxErr == nil || !strings.HasPrefix(overMaxErr.Error(), "ChartInvalidDonutHolePercent") {
		t.Errorf("OverMaxHolePercentShouldFail: %v", overMaxErr)
	}
}

func TestChartInvalidLegendAlignmentFails(t *testing.T) {
	invalidSettings := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Legend: ChartLegendSettings{Alignment: "middle"},
	}
	renderErr := Chart(invalidSettings).Render(context.Background(), &bytes.Buffer{})
	if renderErr == nil || !strings.HasPrefix(renderErr.Error(), "ChartInvalidLegendAlignment") {
		t.Errorf("InvalidLegendAlignmentShouldFail: %v", renderErr)
	}
}

func TestChartTooltipRingClassDefaults(t *testing.T) {
	defaultedSettings := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
	}
	if actual := defaultedSettings.tooltipClassesResolver(); !strings.Contains(actual, "ring-1 ring-neutral-50/10") {
		t.Errorf("DefaultTooltipRingMissing: %q", actual)
	}

	providedSettings := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Tooltip: ChartTooltipSettings{
			RingColor: "sky-400", RingThickness: uiToolset.RingThicknessMd,
		},
	}
	if actual := providedSettings.tooltipClassesResolver(); strings.Contains(actual, "ring-neutral-50/10") {
		t.Errorf("ProvidedTooltipRingOverridden: %q", actual)
	}
}

func TestChartColorTokensCollector(t *testing.T) {
	resolvedConfig, resolveErr := ChartSettings{
		XAxis: ChartXAxisSettings{RowKey: "time", TickLabelColor: "sky-300"},
		Series: []ChartSeriesSettings{
			{RowKey: "cpu", IdentityColor: "sky-300"},
			{RowKey: "memory"},
		},
		Donut: ChartDonutSettings{
			SliceIdentityColors:    map[string]string{"nginx": "[#4c8dff]"},
			CenterLabelDisplayMode: ChartDonutCenterLabelDisplayModeAlways,
			CenterLabelSliceColor:  "amber-300",
			CenterLabelValueColor:  "amber-100",
		},
	}.configResolver()
	if resolveErr != nil {
		t.Fatalf("UnexpectedChartConfigError: %v", resolveErr)
	}
	colorTokens := resolvedConfig.colorTokensCollector()
	for _, expectedToken := range []string{"sky-300", "[#4c8dff]", "neutral-900", "neutral-50/30", "amber-300", "amber-100"} {
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

func TestChartHiddenCenterLabelSkipsItsColorTokens(t *testing.T) {
	resolvedConfig, resolveErr := ChartSettings{
		Type:   ChartTypeDonut,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Donut: ChartDonutSettings{
			CenterLabelSliceColor: "amber-300",
			CenterLabelValueColor: "amber-100",
		},
	}.configResolver()
	if resolveErr != nil {
		t.Fatalf("UnexpectedChartConfigError: %v", resolveErr)
	}
	for _, unexpectedToken := range []string{"amber-300", "amber-100"} {
		if slices.Contains(resolvedConfig.colorTokensCollector(), unexpectedToken) {
			t.Errorf("HiddenCenterLabelColorTokenCollected: %q", unexpectedToken)
		}
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

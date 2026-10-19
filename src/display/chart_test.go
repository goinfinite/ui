package uiDisplay

import (
	"bytes"
	"context"
	"math"
	"slices"
	"strconv"
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
			settings:          ChartSettings{Type: "bubble", XAxis: timeAxis, Series: cpuSeries},
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
				Slices: ChartSliceSettings{HolePercent: 95},
			},
			expectedErrPrefix: "ChartInvalidSliceHolePercent",
		},
		{
			name: "UnknownCenterLabelDisplayModeFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Slices: ChartSliceSettings{CenterLabelDisplayMode: "sometimes"},
			},
			expectedErrPrefix: "ChartInvalidCenterLabelDisplayMode",
		},
		{
			name: "UnknownCenterLabelSizeFails",
			settings: ChartSettings{
				Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
				Slices: ChartSliceSettings{CenterLabelSize: "huge"},
			},
			expectedErrPrefix: "ChartInvalidCenterLabelSize",
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
		Slices: ChartSliceSettings{HolePercent: chartMaxSliceHolePercent},
	}
	if _, resolveErr := atMaxSettings.configResolver(); resolveErr != nil {
		t.Errorf("AtMaxHolePercentShouldPass: %v", resolveErr)
	}

	overMaxSettings := ChartSettings{
		Type: ChartTypeDonut, XAxis: ChartXAxisSettings{RowKey: "time"}, Series: []ChartSeriesSettings{{RowKey: "cpu"}},
		Slices: ChartSliceSettings{HolePercent: chartMaxSliceHolePercent + 1},
	}
	_, overMaxErr := overMaxSettings.configResolver()
	if overMaxErr == nil || !strings.HasPrefix(overMaxErr.Error(), "ChartInvalidSliceHolePercent") {
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
		Slices: ChartSliceSettings{
			SliceIdentityColors:    map[string]string{"nginx": "[#4c8dff]"},
			CenterLabelDisplayMode: ChartCenterLabelDisplayModeAlways,
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
		Slices: ChartSliceSettings{
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
		{
			Type:   ChartTypeDonut,
			XAxis:  ChartXAxisSettings{RowKey: "browser"},
			Series: []ChartSeriesSettings{{RowKey: "sessions"}},
			Slices: ChartSliceSettings{HolePercent: 95},
		},
	} {
		renderErr := Chart(invalidSettings).Render(context.Background(), &bytes.Buffer{})
		if renderErr == nil {
			t.Errorf("ExpectedRenderError: %+v", invalidSettings)
		}
	}
}

func TestChartHorizontalBarConfigType(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeHorizontalBar,
		XAxis:  ChartXAxisSettings{RowKey: "month"},
		Series: []ChartSeriesSettings{{RowKey: "requests"}},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.ChartType != ChartTypeHorizontalBar {
		t.Errorf("ChartType = %q, expected %q",
			resolvedConfig.ChartType, ChartTypeHorizontalBar)
	}
}

func TestChartRadialBarConfigType(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeRadialBar,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Slices: ChartSliceSettings{HolePercent: 40},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.ChartType != ChartTypeRadialBar {
		t.Errorf("ChartType = %q, expected %q",
			resolvedConfig.ChartType, ChartTypeRadialBar)
	}
	if resolvedConfig.Slices.HolePercent != 40 {
		t.Errorf("HolePercent = %d, expected 40 (radialBar must keep the donut settings)",
			resolvedConfig.Slices.HolePercent)
	}
}

func TestChartWaffleConfigType(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeWaffle,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.ChartType != ChartTypeWaffle {
		t.Errorf("ChartType = %q, expected %q",
			resolvedConfig.ChartType, ChartTypeWaffle)
	}
	if resolvedConfig.Waffle.CellCount != 100 || resolvedConfig.Waffle.ColumnCount != 10 {
		t.Errorf("Waffle = %+v, expected 100 cells in 10 columns",
			resolvedConfig.Waffle)
	}
}

func TestChartWaffleKeepsConfiguredGrid(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeWaffle,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Waffle: ChartWaffleSettings{CellCount: 25, ColumnCount: 5},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Waffle.CellCount != 25 || resolvedConfig.Waffle.ColumnCount != 5 {
		t.Errorf("Waffle = %+v, expected 25 cells in 5 columns",
			resolvedConfig.Waffle)
	}
}

func TestChartSunburstConfigType(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:     ChartTypeSunburst,
		XAxis:    ChartXAxisSettings{RowKey: "service"},
		Series:   []ChartSeriesSettings{{RowKey: "usedGib"}},
		Sunburst: ChartSunburstSettings{ChildRowKey: "dataset"},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.ChartType != ChartTypeSunburst {
		t.Errorf("ChartType = %q, expected %q",
			resolvedConfig.ChartType, ChartTypeSunburst)
	}
	if resolvedConfig.Sunburst.ChildRowKey != "dataset" {
		t.Errorf("ChildRowKey = %q, expected %q",
			resolvedConfig.Sunburst.ChildRowKey, "dataset")
	}
}

func TestChartSunburstRequiresChildRowKey(t *testing.T) {
	_, configErr := ChartSettings{
		Type:   ChartTypeSunburst,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
	}.configResolver()
	if configErr == nil || configErr.Error() != "ChartSunburstChildRowKeyMissing" {
		t.Fatalf("ConfigResolverError = %v, expected ChartSunburstChildRowKeyMissing",
			configErr)
	}
}

func TestChartParliamentConfigType(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:    ChartTypeParliament,
		XAxis:   ChartXAxisSettings{RowKey: "party"},
		Series:  []ChartSeriesSettings{{RowKey: "seats"}},
		Rows:    []map[string]any{{"party": "Center", "seats": 52}},
		Tooltip: ChartTooltipSettings{ValueFormat: ",.0f"},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.ChartType != ChartTypeParliament {
		t.Errorf("ChartType = %q, expected %q",
			resolvedConfig.ChartType, ChartTypeParliament)
	}
	if resolvedConfig.Parliament.RowCount != chartDefaultParliamentRowCount {
		t.Errorf("RowCount = %d, expected %d",
			resolvedConfig.Parliament.RowCount, chartDefaultParliamentRowCount)
	}
}

func TestChartParliamentKeepsExplicitRowCount(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:       ChartTypeParliament,
		XAxis:      ChartXAxisSettings{RowKey: "party"},
		Series:     []ChartSeriesSettings{{RowKey: "seats"}},
		Rows:       []map[string]any{{"party": "Center", "seats": 52}},
		Parliament: ChartParliamentSettings{RowCount: 8},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Parliament.RowCount != 8 {
		t.Errorf("RowCount = %d, expected %d",
			resolvedConfig.Parliament.RowCount, 8)
	}
}

func TestChartPieForcesZeroHoleAndHiddenCenterLabel(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypePie,
		XAxis:  ChartXAxisSettings{RowKey: "browser"},
		Series: []ChartSeriesSettings{{RowKey: "sessions"}},
		Slices: ChartSliceSettings{
			HolePercent:            60,
			CenterLabelDisplayMode: ChartCenterLabelDisplayModeAlways,
		},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Slices.HolePercent != 0 {
		t.Errorf("HolePercent = %d, expected 0", resolvedConfig.Slices.HolePercent)
	}
	if resolvedConfig.Slices.CenterLabelDisplayMode != ChartCenterLabelDisplayModeHidden {
		t.Errorf("CenterLabelDisplayMode = %q, expected %q",
			resolvedConfig.Slices.CenterLabelDisplayMode,
			ChartCenterLabelDisplayModeHidden)
	}
}

func TestChartDonutKeepsHoleAndCenterLabelSettings(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeDonut,
		XAxis:  ChartXAxisSettings{RowKey: "browser"},
		Series: []ChartSeriesSettings{{RowKey: "sessions"}},
		Slices: ChartSliceSettings{CenterLabelDisplayMode: ChartCenterLabelDisplayModeAlways},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Slices.HolePercent != chartDefaultSliceHolePercent {
		t.Errorf("HolePercent = %d, expected %d",
			resolvedConfig.Slices.HolePercent, chartDefaultSliceHolePercent)
	}
	if resolvedConfig.Slices.CenterLabelDisplayMode != ChartCenterLabelDisplayModeAlways {
		t.Errorf("CenterLabelDisplayMode = %q, expected %q",
			resolvedConfig.Slices.CenterLabelDisplayMode,
			ChartCenterLabelDisplayModeAlways)
	}
}

func TestChartProgressBoundsResolver(t *testing.T) {
	tests := []struct {
		minValue    float64
		maxValue    float64
		expectedMin float64
		expectedMax float64
	}{
		{0, 0, 0, 100},
		{0, 10, 0, 10},
		{10, 0, 10, 110},
		{-50, 50, -50, 50},
		{5, 5, 5, 105},
	}
	for _, test := range tests {
		progressBounds := ChartProgressSettings{
			MinValue: test.minValue,
			MaxValue: test.maxValue,
		}.boundsResolver()
		if progressBounds.Min != test.expectedMin || progressBounds.Max != test.expectedMax {
			t.Errorf("boundsResolver(%v, %v) = {%v, %v}, expected {%v, %v}",
				test.minValue, test.maxValue,
				progressBounds.Min, progressBounds.Max,
				test.expectedMin, test.expectedMax)
		}
	}
}

func TestChartProgressPercentResolver(t *testing.T) {
	tests := []struct {
		value    float64
		minValue float64
		maxValue float64
		expected float64
	}{
		{50, 0, 100, 50},
		{150, 0, 100, 100},
		{-10, 0, 100, 0},
		{25, 0, 50, 50},
		{3, 2, 4, 50},
		{0, 0, 0, 0},
	}
	for _, test := range tests {
		result := ChartProgressSettings{
			Value:    test.value,
			MinValue: test.minValue,
			MaxValue: test.maxValue,
		}.percentResolver()
		if result != test.expected {
			t.Errorf("percentResolver(value=%v, min=%v, max=%v) = %v, expected %v",
				test.value, test.minValue, test.maxValue, result, test.expected)
		}
	}
}

func TestChartProgressThresholdsResolver(t *testing.T) {
	customThresholds := []ChartProgressThresholdSettings{
		{MinimumPercent: 80, BackgroundColor: "emerald-500"},
	}
	resolvedCustom := ChartProgressSettings{Thresholds: customThresholds}.thresholdsResolver()
	if len(resolvedCustom) != 1 {
		t.Fatalf("thresholdsResolver custom count = %d, expected 1", len(resolvedCustom))
	}
	if resolvedCustom[0].MinimumPercent != 80 || resolvedCustom[0].BackgroundColor != "emerald-500" {
		t.Errorf("thresholdsResolver custom = %+v, expected the custom threshold", resolvedCustom[0])
	}

	resolvedDefaults := ChartProgressSettings{}.thresholdsResolver()
	expectedDefaults := []chartProgressThreshold{
		{MinimumPercent: 90, BackgroundColor: "red-800"},
		{MinimumPercent: 85, BackgroundColor: "red-600"},
		{MinimumPercent: 75, BackgroundColor: "amber-600"},
		{MinimumPercent: 50, BackgroundColor: "amber-300"},
	}
	if len(resolvedDefaults) != len(expectedDefaults) {
		t.Fatalf("thresholdsResolver default count = %d, expected %d",
			len(resolvedDefaults), len(expectedDefaults))
	}
	for index, expectedThreshold := range expectedDefaults {
		if resolvedDefaults[index] != expectedThreshold {
			t.Errorf("thresholdsResolver default[%d] = %+v, expected %+v",
				index, resolvedDefaults[index], expectedThreshold)
		}
	}
}

func TestChartProgressFillColorResolver(t *testing.T) {
	tests := []struct {
		percent      float64
		defaultColor string
		thresholds   []ChartProgressThresholdSettings
		expected     string
	}{
		{95, "", nil, "bg-red-800"},
		{90, "", nil, "bg-red-800"},
		{89, "", nil, "bg-red-600"},
		{85, "", nil, "bg-red-600"},
		{80, "", nil, "bg-amber-600"},
		{75, "", nil, "bg-amber-600"},
		{60, "", nil, "bg-amber-300"},
		{50, "", nil, "bg-amber-300"},
		{20, "", nil, "bg-cyan-400"},
		{20, "violet-500", nil, "bg-violet-500"},
		{85, "", []ChartProgressThresholdSettings{{MinimumPercent: 80, BackgroundColor: "emerald-500"}}, "bg-emerald-500"},
		{50, "", []ChartProgressThresholdSettings{{MinimumPercent: 80, BackgroundColor: "emerald-500"}}, "bg-cyan-400"},
	}
	for _, test := range tests {
		result := ChartProgressSettings{
			DefaultColor: test.defaultColor,
			Thresholds:   test.thresholds,
		}.fillColorResolver(test.percent)
		if result != test.expected {
			t.Errorf("fillColorResolver(%v) = %q, expected %q",
				test.percent, result, test.expected)
		}
	}
}

func TestChartProgressPercentTextResolver(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{0, "0%"},
		{50, "50%"},
		{100, "100%"},
		{33.333333, "33.3%"},
		{42.55, "42.6%"},
	}
	for _, test := range tests {
		result := chartProgressPercentTextResolver(test.percent)
		if result != test.expected {
			t.Errorf("percentTextResolver(%v) = %q, expected %q",
				test.percent, result, test.expected)
		}
	}
}

func TestChartProgressSizeTokensResolver(t *testing.T) {
	tests := []struct {
		size                string
		expectedTrackHeight string
		expectedLabelSize   string
	}{
		{ChartProgressSizeXs, "h-4", "text-[0.625rem]"},
		{ChartProgressSizeSm, "h-6", "text-xs"},
		{ChartProgressSizeMd, "h-8", "text-sm"},
		{"", "h-6", "text-xs"},
		{"unknown", "h-6", "text-xs"},
	}
	for _, test := range tests {
		result := chartProgressSizeTokensResolver(test.size)
		if result.TrackHeight != test.expectedTrackHeight {
			t.Errorf("sizeTokensResolver(%q) trackHeight = %q, expected %q",
				test.size, result.TrackHeight, test.expectedTrackHeight)
		}
		if result.LabelTextSize != test.expectedLabelSize {
			t.Errorf("sizeTokensResolver(%q) labelTextSize = %q, expected %q",
				test.size, result.LabelTextSize, test.expectedLabelSize)
		}
	}
}

func TestChartProgressLabelTextSizeResolver(t *testing.T) {
	tests := []struct {
		configuredClass string
		size            string
		expected        string
	}{
		{"", "", "text-xs"},
		{"", ChartProgressSizeXs, "text-[0.625rem]"},
		{"", ChartProgressSizeSm, "text-xs"},
		{"", ChartProgressSizeMd, "text-sm"},
		{"text-lg", ChartProgressSizeSm, "text-lg"},
	}
	for _, test := range tests {
		result := chartProgressLabelTextSizeResolver(test.configuredClass, test.size)
		if result != test.expected {
			t.Errorf("labelTextSizeResolver(%q, %q) = %q, expected %q",
				test.configuredClass, test.size, result, test.expected)
		}
	}
}

func TestChartProgressRenderDeterminate(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type:      ChartTypeProgress,
		AriaLabel: "Upload progress",
		Progress: ChartProgressSettings{
			Value:    25,
			MaxValue: 50,
			Label:    "Half",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`role="progressbar"`,
		`aria-valuemin="0"`,
		`aria-valuemax="50"`,
		`aria-valuenow="25.00"`,
		`aria-label="Upload progress"`,
		`style="width: 50%;"`,
		`bg-amber-300`,
		`bg-neutral-950/70`,
		`>Half</span>`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "x-data") {
		t.Error("StaticChartProgressMustNotUseAlpine")
	}
	if strings.Contains(renderedHtml, "chartCanvas") {
		t.Error("ChartProgressMustNotRenderVegaCanvas")
	}
}

func TestChartProgressRenderPercentLabel(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			Value:                  42,
			ShouldShowPercentLabel: true,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	if !strings.Contains(buffer.String(), `>42%</span>`) {
		t.Error("RenderedHtmlMissingPercentLabel")
	}
}

func TestChartProgressRenderBound(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type:      ChartTypeProgress,
		AriaLabel: "Upload progress",
		Progress:  ChartProgressSettings{ValueOneWayStatePath: "uploadPercent"},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-effect="syncValueFromPath(uploadPercent)"`,
		`:aria-valuenow="valueAriaText"`,
		`:class="fillColorClass"`,
		`:style="&#39;width: &#39; + percentWidth"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if !strings.Contains(renderedHtml, `chartProgress(0,0,100`) {
		t.Error("RenderedHtmlMissingChartProgressState")
	}
}

func TestChartProgressRenderIndeterminate(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			IsIndeterminate: true,
			Label:           "Working",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`ui-progress-bar-indeterminate`,
		`style="animation: ui-progress-bar-indeterminate 1.5s ease-in-out infinite;"`,
		`w-2/5`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "aria-valuenow") {
		t.Error("IndeterminateChartProgressMustNotExposeAriaValueNow")
	}
}

func TestChartGaugeAngleResolver(t *testing.T) {
	tests := []struct {
		percent      float64
		isSemicircle bool
		expected     float64
	}{
		{0, false, 0},
		{25, false, 90},
		{50, false, 180},
		{100, false, 360},
		{0, true, -90},
		{25, true, -45},
		{50, true, 0},
		{100, true, 90},
	}
	for _, test := range tests {
		result := chartGaugeAngleResolver(test.percent, test.isSemicircle)
		if result != test.expected {
			t.Errorf("chartGaugeAngleResolver(%v, %v) = %v, expected %v",
				test.percent, test.isSemicircle, result, test.expected)
		}
	}
}

func TestChartGaugeArcPathResolver(t *testing.T) {
	tests := []struct {
		minimumPercent float64
		maximumPercent float64
		isSemicircle   bool
		expected       string
	}{
		{0, 60, false, "M100.00 18.00 A82.00 82.00 0 1 1 51.80 166.34"},
		{60, 80, false, "M51.80 166.34 A82.00 82.00 0 0 1 22.01 74.66"},
		{80, 100, false, "M22.01 74.66 A82.00 82.00 0 0 1 100.00 18.00"},
		{0, 50, true, "M18.00 100.00 A82.00 82.00 0 0 1 100.00 18.00"},
	}
	for _, test := range tests {
		result := chartGaugeArcPathResolver(
			chartGaugeAngleResolver(test.minimumPercent, test.isSemicircle),
			chartGaugeAngleResolver(test.maximumPercent, test.isSemicircle),
			chartGaugeOuterRadius-chartGaugeDefaultArcThickness/2,
		)
		if result != test.expected {
			t.Errorf("arcPath(%v, %v, semicircle=%v) = %q, expected %q",
				test.minimumPercent, test.maximumPercent, test.isSemicircle,
				result, test.expected)
		}
	}
}

func TestChartGaugeBandValidationFails(t *testing.T) {
	invalidBands := []ChartGaugeBandSettings{
		{MinimumPercent: 70, MaximumPercent: 40},
		{MinimumPercent: -10, MaximumPercent: 50},
		{MinimumPercent: 50, MaximumPercent: 101},
		{MinimumPercent: 60, MaximumPercent: 60},
		{MinimumPercent: math.NaN(), MaximumPercent: 40},
	}
	for _, invalidBand := range invalidBands {
		var buffer bytes.Buffer
		renderErr := Chart(ChartSettings{
			Type:  ChartTypeGauge,
			Gauge: ChartGaugeSettings{Bands: []ChartGaugeBandSettings{invalidBand}},
		}).Render(context.Background(), &buffer)
		if renderErr == nil {
			t.Errorf("GaugeBandMustFail: %+v", invalidBand)
			continue
		}
		if !strings.Contains(renderErr.Error(), "ChartInvalidGaugeBandPercent") {
			t.Errorf("GaugeBandErrorTextMismatch: %v", renderErr)
		}
	}
}

func TestChartGaugeClampedValueResolver(t *testing.T) {
	tests := []struct {
		value    float64
		minValue float64
		maxValue float64
		expected float64
	}{
		{65, 0, 100, 65},
		{12, 8, 24, 12},
		{30, 8, 24, 24},
		{-5, 8, 24, 8},
		{150, 0, 100, 100},
		{0, 0, 0, 0},
	}
	for _, test := range tests {
		result := ChartGaugeSettings{
			Value:    test.value,
			MinValue: test.minValue,
			MaxValue: test.maxValue,
		}.clampedValueResolver()
		if result != test.expected {
			t.Errorf("clampedValueResolver(value=%v, min=%v, max=%v) = %v, expected %v",
				test.value, test.minValue, test.maxValue, result, test.expected)
		}
	}
}

func TestChartGaugeTextSettingsResolver(t *testing.T) {
	defaultSettings, defaultErr := ChartGaugeSettings{}.textSettingsResolver()
	if defaultErr != nil {
		t.Fatalf("TextSettingsResolverFailed: %v", defaultErr)
	}
	expectedDefaults := chartGaugeTextSettings{
		LabelSize: chartGaugeDefaultLabelTextSize,
		Offset:    0,
		ValueSize: chartGaugeDefaultValueTextSize,
	}
	if defaultSettings != expectedDefaults {
		t.Errorf("DefaultTextSettings = %+v, expected %+v",
			defaultSettings, expectedDefaults)
	}
	customSettings, customErr := ChartGaugeSettings{
		LabelTextSize: 16,
		TextOffset:    20,
		ValueTextSize: 36,
	}.textSettingsResolver()
	if customErr != nil {
		t.Fatalf("TextSettingsResolverFailed: %v", customErr)
	}
	expectedCustom := chartGaugeTextSettings{LabelSize: 16, Offset: 20, ValueSize: 36}
	if customSettings != expectedCustom {
		t.Errorf("CustomTextSettings = %+v, expected %+v",
			customSettings, expectedCustom)
	}
	invalidSettings := []struct {
		settings       ChartGaugeSettings
		expectedPrefix string
	}{
		{ChartGaugeSettings{ValueTextSize: 4}, "ChartInvalidGaugeValueTextSize"},
		{ChartGaugeSettings{ValueTextSize: 80}, "ChartInvalidGaugeValueTextSize"},
		{ChartGaugeSettings{ValueTextSize: math.NaN()}, "ChartInvalidGaugeValueTextSize"},
		{ChartGaugeSettings{LabelTextSize: 4}, "ChartInvalidGaugeLabelTextSize"},
		{ChartGaugeSettings{LabelTextSize: 80}, "ChartInvalidGaugeLabelTextSize"},
		{ChartGaugeSettings{LabelTextSize: math.NaN()}, "ChartInvalidGaugeLabelTextSize"},
		{ChartGaugeSettings{TextOffset: -1}, "ChartInvalidGaugeTextOffset"},
		{ChartGaugeSettings{TextOffset: 81}, "ChartInvalidGaugeTextOffset"},
		{ChartGaugeSettings{TextOffset: math.NaN()}, "ChartInvalidGaugeTextOffset"},
	}
	for _, invalid := range invalidSettings {
		_, settingsErr := invalid.settings.textSettingsResolver()
		if settingsErr == nil ||
			!strings.HasPrefix(settingsErr.Error(), invalid.expectedPrefix) {
			t.Errorf("TextSettingsMustFail (%+v): %v",
				invalid.settings, settingsErr)
		}
	}
}

func TestChartGaugeTextLayoutResolver(t *testing.T) {
	defaultTextSettings := chartGaugeTextSettings{
		LabelSize: chartGaugeDefaultLabelTextSize,
		ValueSize: chartGaugeDefaultValueTextSize,
	}
	tests := []struct {
		hasLabel     bool
		isSemicircle bool
		textSettings chartGaugeTextSettings
		expected     chartGaugeTextLayout
	}{
		{false, false, defaultTextSettings,
			chartGaugeTextLayout{ValueY: "134", Bottom: 146}},
		{false, true, defaultTextSettings,
			chartGaugeTextLayout{ValueY: "124", Bottom: 136}},
		{true, false, defaultTextSettings,
			chartGaugeTextLayout{ValueY: "134", LabelY: "156", Bottom: 162}},
		{true, true, defaultTextSettings,
			chartGaugeTextLayout{ValueY: "122", LabelY: "142", Bottom: 148}},
		{true, true, chartGaugeTextSettings{LabelSize: 16, Offset: 20, ValueSize: 36},
			chartGaugeTextLayout{ValueY: "142", LabelY: "162", Bottom: 170}},
		{true, true, chartGaugeTextSettings{LabelSize: 6, Offset: 0, ValueSize: 72},
			chartGaugeTextLayout{ValueY: "122", LabelY: "142", Bottom: 158}},
	}
	for _, test := range tests {
		result := chartGaugeTextLayoutResolver(
			test.hasLabel, test.isSemicircle, test.textSettings,
		)
		if result != test.expected {
			t.Errorf("chartGaugeTextLayoutResolver(%v, %v, %+v) = %+v, expected %+v",
				test.hasLabel, test.isSemicircle, test.textSettings,
				result, test.expected)
		}
	}
}

func TestChartGaugeRenderFullCircle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type:      ChartTypeGauge,
		AriaLabel: "CPU load",
		Gauge: ChartGaugeSettings{
			Value: 65,
			Bands: []ChartGaugeBandSettings{
				{MinimumPercent: 0, MaximumPercent: 60, Color: "emerald-500"},
				{MinimumPercent: 60, MaximumPercent: 80, Color: "amber-500"},
				{MinimumPercent: 80, MaximumPercent: 100, Color: "red-500"},
			},
			Label: "CPU",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`role="meter"`,
		`aria-valuemin="0"`,
		`aria-valuemax="100"`,
		`aria-valuenow="65.00"`,
		`aria-label="CPU load"`,
		`viewBox="0 0 200 200"`,
		`class="stroke-neutral-50/15"`,
		`d="M100.00 18.00 A82.00 82.00 0 1 1 51.80 166.34"`,
		`class="stroke-emerald-500"`,
		`d="M51.80 166.34 A82.00 82.00 0 0 1 22.01 74.66"`,
		`class="stroke-amber-500"`,
		`d="M22.01 74.66 A82.00 82.00 0 0 1 100.00 18.00"`,
		`class="stroke-red-500"`,
		`transform="rotate(234 100 100)"`,
		`y2="45"`,
		`y="134"`,
		`y="156"`,
		`>65</text>`,
		`>CPU</text>`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "chartCanvas") || strings.Contains(renderedHtml, "x-data") {
		t.Error("ChartGaugeMustNotRenderVegaOrAlpine")
	}
}

func TestChartGaugeRenderSemicircle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeGauge,
		Gauge: ChartGaugeSettings{
			Value:        72,
			IsSemicircle: true,
			Label:        "Memory",
			NeedleColor:  "sky-400",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`viewBox="0 0 200 150"`,
		`aria-valuenow="72.00"`,
		`class="stroke-sky-400"`,
		`transform="rotate(39.6 100 100)"`,
		`y="122"`,
		`y="142"`,
		`>72</text>`,
		`>Memory</text>`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartGaugeRenderCustomTextSizeAndOffset(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeGauge,
		Gauge: ChartGaugeSettings{
			Value:         72,
			IsSemicircle:  true,
			Label:         "Memory",
			LabelTextSize: 16,
			TextOffset:    20,
			ValueTextSize: 36,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`viewBox="0 0 200 172"`,
		`font-size="36"`,
		`font-size="16"`,
		`y="142"`,
		`y="162"`,
		`>72</text>`,
		`>Memory</text>`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartGaugeRenderValueTextSetsViewBoxBottom(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeGauge,
		Gauge: ChartGaugeSettings{
			Value:         65,
			IsSemicircle:  true,
			Label:         "X",
			LabelTextSize: 6,
			ValueTextSize: 72,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	if !strings.Contains(buffer.String(), `viewBox="0 0 200 160"`) {
		t.Errorf("RenderedHtmlMissingExpandedViewBox: %s", buffer.String())
	}
}

func TestChartGaugeRenderClampsValueToBounds(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeGauge,
		Gauge: ChartGaugeSettings{
			Value:    30,
			MinValue: 8,
			MaxValue: 24,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`aria-valuemin="8"`,
		`aria-valuemax="24"`,
		`aria-valuenow="24.00"`,
		`transform="rotate(360 100 100)"`,
		`>24</text>`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartParliamentDotShapeAndSize(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:    ChartTypeParliament,
		XAxis:   ChartXAxisSettings{RowKey: "party"},
		Series:  []ChartSeriesSettings{{RowKey: "seats"}},
		Rows:    []map[string]any{{"party": "Center", "seats": 52}},
		Tooltip: ChartTooltipSettings{ValueFormat: ",.0f"},
		Parliament: ChartParliamentSettings{
			DotShape:       ChartDotShapeTriangleUp,
			DotSizePercent: 150,
		},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Parliament.DotShape != "triangle-up" {
		t.Errorf("DotShape = %q, expected triangle-up", resolvedConfig.Parliament.DotShape)
	}
	if resolvedConfig.Parliament.DotSizePercent != 150 {
		t.Errorf("DotSizePercent = %v, expected 150", resolvedConfig.Parliament.DotSizePercent)
	}
}

func TestChartParliamentDotDefaults(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:    ChartTypeParliament,
		XAxis:   ChartXAxisSettings{RowKey: "party"},
		Series:  []ChartSeriesSettings{{RowKey: "seats"}},
		Rows:    []map[string]any{{"party": "Center", "seats": 52}},
		Tooltip: ChartTooltipSettings{ValueFormat: ",.0f"},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Parliament.DotShape != ChartDotShapeCircle {
		t.Errorf("DotShape = %q, expected circle", resolvedConfig.Parliament.DotShape)
	}
	if resolvedConfig.Parliament.DotGapPercent != chartDefaultParliamentDotGapPercent {
		t.Errorf("DotGapPercent = %v, expected %v",
			resolvedConfig.Parliament.DotGapPercent,
			chartDefaultParliamentDotGapPercent)
	}
	if resolvedConfig.Parliament.DotSizePercent != chartDefaultParliamentDotSizePercent {
		t.Errorf("DotSizePercent = %v, expected %v",
			resolvedConfig.Parliament.DotSizePercent,
			chartDefaultParliamentDotSizePercent)
	}
}

func TestChartParliamentKeepsExplicitDotGap(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:    ChartTypeParliament,
		XAxis:   ChartXAxisSettings{RowKey: "party"},
		Series:  []ChartSeriesSettings{{RowKey: "seats"}},
		Rows:    []map[string]any{{"party": "Center", "seats": 52}},
		Tooltip: ChartTooltipSettings{ValueFormat: ",.0f"},
		Parliament: ChartParliamentSettings{
			DotGapPercent: 5,
		},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Parliament.DotGapPercent != 5 {
		t.Errorf("DotGapPercent = %v, expected 5",
			resolvedConfig.Parliament.DotGapPercent)
	}
}

func TestChartParliamentInvalidDotSettingsFail(t *testing.T) {
	invalidSettings := []ChartParliamentSettings{
		{DotShape: "hexagon"},
		{DotGapPercent: -5},
		{DotGapPercent: 99},
		{DotGapPercent: math.NaN()},
		{DotSizePercent: 5},
		{DotSizePercent: 250},
		{DotSizePercent: math.NaN()},
	}
	expectedErrorPrefixes := []string{
		"ChartInvalidParliamentDotShape",
		"ChartInvalidParliamentDotGapPercent",
		"ChartInvalidParliamentDotGapPercent",
		"ChartInvalidParliamentDotGapPercent",
		"ChartInvalidParliamentDotSizePercent",
		"ChartInvalidParliamentDotSizePercent",
		"ChartInvalidParliamentDotSizePercent",
	}
	for settingIndex, invalidSetting := range invalidSettings {
		_, configErr := ChartSettings{
			Type:       ChartTypeParliament,
			XAxis:      ChartXAxisSettings{RowKey: "party"},
			Series:     []ChartSeriesSettings{{RowKey: "seats"}},
			Parliament: invalidSetting,
		}.configResolver()
		if configErr == nil ||
			!strings.HasPrefix(configErr.Error(), expectedErrorPrefixes[settingIndex]) {
			t.Errorf("ParliamentSettingsMustFail (%+v): %v", invalidSetting, configErr)
		}
	}
}

func TestChartWaffleCellShapeAndGap(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeWaffle,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Waffle: ChartWaffleSettings{
			CellGapPercent: 25,
			CellShape:      ChartWaffleCellShapeCircle,
		},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Waffle.CellShape != ChartWaffleCellShapeCircle {
		t.Errorf("CellShape = %q, expected circle", resolvedConfig.Waffle.CellShape)
	}
	if resolvedConfig.Waffle.CellGapPercent != 25 {
		t.Errorf("CellGapPercent = %v, expected 25", resolvedConfig.Waffle.CellGapPercent)
	}
}

func TestChartWaffleCellDefaults(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeWaffle,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Waffle.CellShape != ChartWaffleCellShapeSquare {
		t.Errorf("CellShape = %q, expected square", resolvedConfig.Waffle.CellShape)
	}
	if resolvedConfig.Waffle.CellGapPercent != chartDefaultWaffleCellGapPercent {
		t.Errorf("CellGapPercent = %v, expected %v",
			resolvedConfig.Waffle.CellGapPercent, chartDefaultWaffleCellGapPercent)
	}
}

func TestChartWaffleInvalidCellSettingsFail(t *testing.T) {
	invalidSettings := []ChartWaffleSettings{
		{CellShape: "triangle"},
		{CellGapPercent: -5},
		{CellGapPercent: 75},
		{CellGapPercent: math.NaN()},
	}
	expectedErrorPrefixes := []string{
		"ChartInvalidWaffleCellShape",
		"ChartInvalidWaffleCellGapPercent",
		"ChartInvalidWaffleCellGapPercent",
		"ChartInvalidWaffleCellGapPercent",
	}
	for settingIndex, invalidSetting := range invalidSettings {
		_, configErr := ChartSettings{
			Type:   ChartTypeWaffle,
			XAxis:  ChartXAxisSettings{RowKey: "service"},
			Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
			Waffle: invalidSetting,
		}.configResolver()
		if configErr == nil ||
			!strings.HasPrefix(configErr.Error(), expectedErrorPrefixes[settingIndex]) {
			t.Errorf("WaffleSettingsMustFail (%+v): %v", invalidSetting, configErr)
		}
	}
}

func TestChartSliceLabelColorDefaultsAndCollection(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypeRadialBar,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	if resolvedConfig.Slices.SliceLabelColor != "neutral-50" {
		t.Errorf("SliceLabelColor = %q, expected neutral-50",
			resolvedConfig.Slices.SliceLabelColor)
	}
	if resolvedConfig.Slices.SliceLabelPillBackgroundColor != chartDefaultSliceLabelPillColor {
		t.Errorf("SliceLabelPillBackgroundColor = %q, expected %q",
			resolvedConfig.Slices.SliceLabelPillBackgroundColor,
			chartDefaultSliceLabelPillColor)
	}
	tokens := resolvedConfig.colorTokensCollector()
	if !slices.Contains(tokens, "neutral-50") {
		t.Errorf("ColorTokens = %v, expected to contain neutral-50", tokens)
	}
	if !slices.Contains(tokens, chartDefaultSliceLabelPillColor) {
		t.Errorf("ColorTokens = %v, expected to contain %q",
			tokens, chartDefaultSliceLabelPillColor)
	}
	hiddenConfig, hiddenConfigErr := ChartSettings{
		Type:   ChartTypeSunburst,
		XAxis:  ChartXAxisSettings{RowKey: "team"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Slices: ChartSliceSettings{
			IsSliceLabelHidden: true,
			SliceLabelColor:    "emerald-200",
		},
		Sunburst: ChartSunburstSettings{ChildRowKey: "service"},
	}.configResolver()
	if hiddenConfigErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", hiddenConfigErr)
	}
	hiddenTokens := hiddenConfig.colorTokensCollector()
	if slices.Contains(hiddenTokens, "emerald-200") {
		t.Errorf("ColorTokens = %v, hidden slice label color must not be collected",
			hiddenTokens)
	}
	if slices.Contains(hiddenTokens, chartDefaultSliceLabelPillColor) {
		t.Errorf("ColorTokens = %v, hidden slice label pill color must not be collected",
			hiddenTokens)
	}
}

func TestChartPieSliceLabelPillColorIsCollected(t *testing.T) {
	resolvedConfig, configErr := ChartSettings{
		Type:   ChartTypePie,
		XAxis:  ChartXAxisSettings{RowKey: "service"},
		Series: []ChartSeriesSettings{{RowKey: "usedGib"}},
		Slices: ChartSliceSettings{SliceLabelPillBackgroundColor: "sky-950/80"},
	}.configResolver()
	if configErr != nil {
		t.Fatalf("ConfigResolverFailed: %v", configErr)
	}
	tokens := resolvedConfig.colorTokensCollector()
	if !slices.Contains(tokens, "sky-950/80") {
		t.Errorf("ColorTokens = %v, expected to contain sky-950/80", tokens)
	}
}

func TestChartGaugeLayoutResolver(t *testing.T) {
	defaultLayout, defaultErr := ChartGaugeSettings{}.layoutResolver()
	if defaultErr != nil {
		t.Fatalf("LayoutResolverFailed: %v", defaultErr)
	}
	expectedDefaultLayout := chartGaugeLayout{
		ArcRadius:       chartGaugeOuterRadius - chartGaugeDefaultArcThickness/2,
		ArcThickness:    chartGaugeDefaultArcThickness,
		NeedleLength:    chartGaugeDefaultNeedleLength,
		NeedleThickness: chartGaugeDefaultNeedleThickness,
		PivotRadius:     chartGaugeDefaultPivotRadius,
	}
	if defaultLayout != expectedDefaultLayout {
		t.Errorf("DefaultLayout = %+v, expected %+v", defaultLayout, expectedDefaultLayout)
	}
	customLayout, customErr := ChartGaugeSettings{
		ArcThickness:    20,
		NeedleLength:    60,
		NeedleThickness: 5,
		PivotRadius:     8,
	}.layoutResolver()
	if customErr != nil {
		t.Fatalf("LayoutResolverFailed: %v", customErr)
	}
	if customLayout.ArcRadius != chartGaugeOuterRadius-10 {
		t.Errorf("ArcRadius = %v, expected %v",
			customLayout.ArcRadius, chartGaugeOuterRadius-10)
	}
	if customLayout.NeedleThickness != 5 || customLayout.PivotRadius != 8 {
		t.Errorf("CustomLayout = %+v", customLayout)
	}
}

func TestChartGaugeInvalidLayoutSettingsFail(t *testing.T) {
	invalidSettings := []ChartGaugeSettings{
		{ArcThickness: 2},
		{ArcThickness: 40},
		{ArcThickness: math.NaN()},
		{NeedleLength: 5},
		{NeedleLength: 100},
		{NeedleThickness: 0.5},
		{NeedleThickness: 20},
		{PivotRadius: 25},
	}
	expectedErrorPrefixes := []string{
		"ChartInvalidGaugeArcThickness",
		"ChartInvalidGaugeArcThickness",
		"ChartInvalidGaugeArcThickness",
		"ChartInvalidGaugeNeedleLength",
		"ChartInvalidGaugeNeedleLength",
		"ChartInvalidGaugeNeedleThickness",
		"ChartInvalidGaugeNeedleThickness",
		"ChartInvalidGaugePivotRadius",
	}
	for settingIndex, invalidSetting := range invalidSettings {
		_, layoutErr := invalidSetting.layoutResolver()
		if layoutErr == nil ||
			!strings.HasPrefix(layoutErr.Error(), expectedErrorPrefixes[settingIndex]) {
			t.Errorf("GaugeLayoutMustFail (%+v): %v", invalidSetting, layoutErr)
		}
	}
}

func TestChartGaugeTicksResolver(t *testing.T) {
	noTicks := ChartGaugeSettings{Value: 10}.ticksResolver(
		chartGaugeLayout{}, ChartGaugeTickPositionOutside,
	)
	if len(noTicks) != 0 {
		t.Errorf("TicksResolver(0) = %d ticks, expected 0", len(noTicks))
	}
	gaugeSettings := ChartGaugeSettings{
		Value:        65,
		IsSemicircle: true,
		MaxValue:     100,
		TickCount:    4,
	}
	layout, layoutErr := gaugeSettings.layoutResolver()
	if layoutErr != nil {
		t.Fatalf("LayoutResolverFailed: %v", layoutErr)
	}
	ticks := gaugeSettings.ticksResolver(layout, ChartGaugeTickPositionOutside)
	if len(ticks) != 5 {
		t.Fatalf("TicksResolver(4) = %d ticks, expected 5", len(ticks))
	}
	expectedLabels := []string{"0", "25", "50", "75", "100"}
	for tickIndex, tick := range ticks {
		if tick.Label != expectedLabels[tickIndex] {
			t.Errorf("Tick[%d].Label = %q, expected %q",
				tickIndex, tick.Label, expectedLabels[tickIndex])
		}
		if tick.LabelX == "" || tick.LabelY == "" {
			t.Errorf("Tick[%d] missing label position", tickIndex)
		}
	}
	fullCircleSettings := ChartGaugeSettings{
		Value:     65,
		MaxValue:  100,
		TickCount: 4,
	}
	fullCircleLayout, layoutErr := fullCircleSettings.layoutResolver()
	if layoutErr != nil {
		t.Fatalf("LayoutResolverFailed: %v", layoutErr)
	}
	fullCircleTicks := fullCircleSettings.ticksResolver(
		fullCircleLayout, ChartGaugeTickPositionOutside,
	)
	if len(fullCircleTicks) != 4 {
		t.Fatalf("TicksResolver(full circle, 4) = %d ticks, expected 4",
			len(fullCircleTicks))
	}
	expectedFullCircleLabels := []string{"0", "25", "50", "75"}
	for tickIndex, tick := range fullCircleTicks {
		if tick.Label != expectedFullCircleLabels[tickIndex] {
			t.Errorf("FullCircleTick[%d].Label = %q, expected %q",
				tickIndex, tick.Label, expectedFullCircleLabels[tickIndex])
		}
	}
}

func TestChartGaugeTickPositionResolver(t *testing.T) {
	tests := []struct {
		tickPosition ChartGaugeTickPosition
		expected     ChartGaugeTickPosition
		expectErr    bool
	}{
		{"", ChartGaugeTickPositionOutside, false},
		{ChartGaugeTickPositionOutside, ChartGaugeTickPositionOutside, false},
		{ChartGaugeTickPositionInside, ChartGaugeTickPositionInside, false},
		{"middle", "", true},
	}
	for _, test := range tests {
		result, resultErr := ChartGaugeSettings{
			TickPosition: test.tickPosition,
		}.tickPositionResolver()
		if test.expectErr {
			if resultErr == nil {
				t.Errorf("tickPositionResolver(%q) must fail", test.tickPosition)
			}
			continue
		}
		if resultErr != nil {
			t.Errorf("tickPositionResolver(%q) failed: %v", test.tickPosition, resultErr)
			continue
		}
		if result != test.expected {
			t.Errorf("tickPositionResolver(%q) = %q, expected %q",
				test.tickPosition, result, test.expected)
		}
	}
}

func TestChartGaugeTickPositionsDiffer(t *testing.T) {
	gaugeSettings := ChartGaugeSettings{
		Value:        65,
		IsSemicircle: true,
		MaxValue:     100,
		TickCount:    2,
	}
	layout, layoutErr := gaugeSettings.layoutResolver()
	if layoutErr != nil {
		t.Fatalf("LayoutResolverFailed: %v", layoutErr)
	}
	outsideTicks := gaugeSettings.ticksResolver(
		layout, ChartGaugeTickPositionOutside,
	)
	insideTicks := gaugeSettings.ticksResolver(
		layout, ChartGaugeTickPositionInside,
	)
	if len(outsideTicks) != len(insideTicks) {
		t.Fatalf("TickCount outside=%d inside=%d, expected the same",
			len(outsideTicks), len(insideTicks))
	}
	outsideStartX, _ := strconv.ParseFloat(outsideTicks[0].LineX1, 64)
	insideStartX, _ := strconv.ParseFloat(insideTicks[0].LineX1, 64)
	if outsideStartX >= insideStartX {
		t.Errorf("Outside tick must sit farther from the centre: outside=%v inside=%v",
			outsideStartX, insideStartX)
	}
	outsideLabelX, _ := strconv.ParseFloat(outsideTicks[0].LabelX, 64)
	insideLabelX, _ := strconv.ParseFloat(insideTicks[0].LabelX, 64)
	if outsideLabelX >= insideLabelX {
		t.Errorf("Outside tick label must sit farther from the centre: outside=%v inside=%v",
			outsideLabelX, insideLabelX)
	}
}

func TestChartGaugeInsideTickLabelSkipsCenterTextArea(t *testing.T) {
	lowerTick := chartGaugeTick{LabelX: "100", LabelY: "154"}
	upperTick := chartGaugeTick{LabelX: "100", LabelY: "46"}
	if chartGaugeTickLabelShouldRender(
		lowerTick, ChartGaugeTickPositionInside, false,
	) {
		t.Error("FullCircleInsideTickLabelMustSkipTheCenterTextArea")
	}
	if !chartGaugeTickLabelShouldRender(
		upperTick, ChartGaugeTickPositionInside, false,
	) {
		t.Error("FullCircleInsideTickLabelMustRenderAboveThePivot")
	}
	if !chartGaugeTickLabelShouldRender(
		lowerTick, ChartGaugeTickPositionInside, true,
	) {
		t.Error("SemicircleInsideTickLabelMustRenderBelowThePivot")
	}
	if !chartGaugeTickLabelShouldRender(
		lowerTick, ChartGaugeTickPositionOutside, false,
	) {
		t.Error("OutsideTickLabelMustRenderBelowThePivot")
	}
}

func TestChartGaugeViewBoxResolver(t *testing.T) {
	tests := []struct {
		isSemicircle    bool
		hasOutsideTicks bool
		textBottom      float64
		textWidth       float64
		expected        string
	}{
		{false, false, 0, 0, "0 0 200 200"},
		{true, false, 0, 0, "0 0 200 150"},
		{false, true, 0, 0, "-20 -20 240 240"},
		{true, true, 0, 0, "-20 -20 240 190"},
		{true, false, 180, 0, "0 0 200 182"},
		{false, false, 210, 0, "0 0 200 212"},
		{true, true, 180, 0, "-20 -20 240 202"},
		{true, false, 0, 300, "-52 0 304 150"},
		{true, true, 0, 300, "-52 -20 304 190"},
	}
	for _, test := range tests {
		result := chartGaugeViewBoxResolver(
			test.isSemicircle, test.hasOutsideTicks, test.textBottom, test.textWidth,
		)
		if result != test.expected {
			t.Errorf("viewBoxResolver(semicircle=%v, outsideTicks=%v, textBottom=%v, textWidth=%v) = %q, expected %q",
				test.isSemicircle, test.hasOutsideTicks, test.textBottom,
				test.textWidth, result, test.expected)
		}
	}
}

func TestChartGaugeTextWidthResolver(t *testing.T) {
	defaultTextSettings := chartGaugeTextSettings{
		LabelSize: chartGaugeDefaultLabelTextSize,
		ValueSize: chartGaugeDefaultValueTextSize,
	}
	valueOnlyWidth := chartGaugeTextWidthResolver(
		"65", "", false, defaultTextSettings,
	)
	if valueOnlyWidth <= 0 {
		t.Errorf("TextWidthResolver value-only = %v, expected a positive width",
			valueOnlyWidth)
	}
	labelWidth := chartGaugeTextWidthResolver(
		"65", "A long label", true, defaultTextSettings,
	)
	if labelWidth <= valueOnlyWidth {
		t.Errorf("TextWidthResolver = %v, expected the longer label width over %v",
			labelWidth, valueOnlyWidth)
	}
}

func TestChartGaugeRenderTicksAndThickness(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeGauge,
		Gauge: ChartGaugeSettings{
			Value:           65,
			IsSemicircle:    true,
			ArcThickness:    20,
			NeedleThickness: 5,
			NeedleLength:    60,
			TickCount:       2,
			TickColor:       "neutral-50/50",
			TickLabelColor:  "amber-200",
			ValueColor:      "emerald-300",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`stroke-width="20"`,
		`stroke-width="5"`,
		`y2="40"`,
		`class="stroke-neutral-50/50"`,
		`class="fill-amber-200"`,
		`class="fill-emerald-300"`,
		`A80.00 80.00`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartProgressRenderStriped(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			Value:       55,
			IsStriped:   true,
			StripeColor: "neutral-50/40",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`text-neutral-50/40`,
		`repeating-linear-gradient(135deg, currentColor 0 0.375rem, transparent 0.375rem 0.75rem)`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartProgressRenderLabelAffixesAndSizeClasses(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			Value:                  82,
			LabelPrefix:            "Used ",
			LabelSuffix:            " of capacity",
			LabelTextSizeClass:     "text-lg",
			MaxHeightClass:         "max-h-10",
			MinHeightClass:         "min-h-4",
			MinWidthClass:          "min-w-40",
			MaxWidthClass:          "max-w-96",
			ShouldShowPercentLabel: true,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`>Used 82% of capacity</span>`,
		`style="width: 82%;"`,
		`text-lg`,
		`min-w-40`,
		`max-w-96`,
		`min-h-4`,
		`max-h-10`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "width: Used") {
		t.Error("StaticFillWidthMustNotContainLabelAffixes")
	}
}

func TestChartProgressRenderDefaultLabelPill(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			Value:                  42,
			ShouldShowPercentLabel: true,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`bg-neutral-950/70`,
		`rounded-full`,
		`py-0.5`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestChartProgressIndeterminateHidesPercentLabel(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Chart(ChartSettings{
		Type: ChartTypeProgress,
		Progress: ChartProgressSettings{
			Value:                  42,
			IsIndeterminate:        true,
			ShouldShowPercentLabel: true,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ChartRenderFailed: %v", renderErr)
	}
	if strings.Contains(buffer.String(), ">42%</span>") {
		t.Error("IndeterminateChartProgressMustNotShowPercentLabel")
	}
}

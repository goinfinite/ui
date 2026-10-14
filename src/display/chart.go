package uiDisplay

import (
	_ "embed"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/a-h/templ"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

//go:embed chartState.js
var chartAlpineState string

var chartAlpineStateOnce = templ.NewOnceHandle(
	templ.WithComponent(uiToolset.MinifierTemplateJs(&chartAlpineState)),
)

var chartConfigScriptIdCounter atomic.Uint64

func generateChartConfigScriptId() string {
	return fmt.Sprintf("chart-config-%d", chartConfigScriptIdCounter.Add(1))
}

type ChartType string

const (
	ChartTypeArea  ChartType = "area"
	ChartTypeBar   ChartType = "bar"
	ChartTypeDonut ChartType = "donut"
	ChartTypeLine  ChartType = "line"
)

type ChartScaleType string

const (
	ChartScaleTypeCategory ChartScaleType = "category"
	ChartScaleTypeNumber   ChartScaleType = "number"
	ChartScaleTypeTime     ChartScaleType = "time"
)

type ChartStackMode string

const (
	ChartStackModeNone    ChartStackMode = "none"
	ChartStackModePercent ChartStackMode = "percent"
	ChartStackModeStacked ChartStackMode = "stacked"
)

type ChartAreaFillStyle string

const (
	ChartAreaFillStyleGradient ChartAreaFillStyle = "gradient"
	ChartAreaFillStyleNone     ChartAreaFillStyle = "none"
	ChartAreaFillStyleSolid    ChartAreaFillStyle = "solid"
)

type ChartMarkerDisplayMode string

const (
	ChartMarkerDisplayModeAlways ChartMarkerDisplayMode = "always"
	ChartMarkerDisplayModeHidden ChartMarkerDisplayMode = "hidden"
	ChartMarkerDisplayModeHover  ChartMarkerDisplayMode = "hover"
)

type ChartLineDashStyle string

const (
	ChartLineDashStyleDashed ChartLineDashStyle = "dashed"
	ChartLineDashStyleDotted ChartLineDashStyle = "dotted"
	ChartLineDashStyleSolid  ChartLineDashStyle = "solid"
)

type ChartYAxisPosition string

const (
	ChartYAxisPositionLeft  ChartYAxisPosition = "left"
	ChartYAxisPositionRight ChartYAxisPosition = "right"
)

type ChartLegendDisplayMode string

const (
	ChartLegendDisplayModeAlways ChartLegendDisplayMode = "always"
	ChartLegendDisplayModeAuto   ChartLegendDisplayMode = "auto"
	ChartLegendDisplayModeHidden ChartLegendDisplayMode = "hidden"
)

type ChartLegendPosition string

const (
	ChartLegendPositionBottom ChartLegendPosition = "bottom"
	ChartLegendPositionTop    ChartLegendPosition = "top"
)

type ChartLegendAlignment string

const (
	ChartLegendAlignmentCenter ChartLegendAlignment = "center"
	ChartLegendAlignmentEnd    ChartLegendAlignment = "end"
	ChartLegendAlignmentStart  ChartLegendAlignment = "start"
)

type ChartDonutCenterLabelDisplayMode string

const (
	ChartDonutCenterLabelDisplayModeAlways ChartDonutCenterLabelDisplayMode = "always"
	ChartDonutCenterLabelDisplayModeHidden ChartDonutCenterLabelDisplayMode = "hidden"
	ChartDonutCenterLabelDisplayModeHover  ChartDonutCenterLabelDisplayMode = "hover"
)

type ChartDonutCenterLabelSize string

const (
	ChartDonutCenterLabelSizeXs ChartDonutCenterLabelSize = "xs"
	ChartDonutCenterLabelSizeSm ChartDonutCenterLabelSize = "sm"
	ChartDonutCenterLabelSizeMd ChartDonutCenterLabelSize = "md"
	ChartDonutCenterLabelSizeLg ChartDonutCenterLabelSize = "lg"
	ChartDonutCenterLabelSizeXl ChartDonutCenterLabelSize = "xl"
)

const (
	chartDefaultTimeTitleFormat       string = "%b %d, %H:%M:%S"
	chartDefaultValueFormat           string = ",.2~f"
	chartDefaultPercentFormat         string = ".0%"
	chartDefaultDonutHolePercent      uint8  = 60
	chartMaxDonutHolePercent          uint8  = 90
	chartDefaultYTickCount            uint   = 5
	chartDefaultCenterLabelSliceColor string = "neutral-400"
	chartDefaultCenterLabelValueColor string = "neutral-50"
)

func chartCenterLabelSizeMultiplierResolver(
	centerLabelSize ChartDonutCenterLabelSize,
) string {
	switch centerLabelSize {
	case ChartDonutCenterLabelSizeXs:
		return "0.75"
	case ChartDonutCenterLabelSizeSm:
		return "0.875"
	case ChartDonutCenterLabelSizeMd:
		return "1"
	case ChartDonutCenterLabelSizeLg:
		return "1.15"
	case ChartDonutCenterLabelSizeXl:
		return "1.3"
	}
	return "1"
}

type ChartSeriesSettings struct {
	RowKey string
	Label  string

	// OptionalFields
	IdentityColor string
	LineDashStyle ChartLineDashStyle
}

type ChartXAxisSettings struct {
	RowKey string

	// OptionalFields
	BaselineColor      string
	GridLineColor      string
	IsGridLineVisible  bool
	IsHidden           bool
	IsTickStartEndOnly bool
	ScaleType          ChartScaleType
	TickLabelColor     string
	TickLabelFormat    string
	Title              string
	TitleColor         string
}

type ChartYAxisSettings struct {
	// OptionalFields
	BaselineColor    string
	GridLineColor    string
	IsGridLineHidden bool
	IsHidden         bool
	MaxValue         string
	MinValue         string
	Position         ChartYAxisPosition
	TickCount        uint
	TickLabelColor   string
	TickLabelFormat  string
	Title            string
	TitleColor       string
}

type ChartLegendSettings struct {
	// OptionalFields
	Alignment   ChartLegendAlignment
	DisplayMode ChartLegendDisplayMode
	GapSize     string
	Position    ChartLegendPosition
	SwatchShape string
	TextColor   string
}

type ChartTooltipSettings struct {
	// OptionalFields
	BackgroundColor   string
	BorderRadius      string
	ContentFunc       string
	CrosshairColor    string
	IsCrosshairHidden bool
	IsHidden          bool
	MaxWidthClass     string
	MinWidthClass     string
	PaddingSize       string
	RingColor         string
	RingThickness     string
	SeriesLabelColor  string
	ShadowSize        string
	TitleColor        string
	TitleFormat       string
	ValueColor        string
	ValueFormat       string
}

type ChartCartesianSettings struct {
	// OptionalFields
	AreaFillStyle     ChartAreaFillStyle
	MarkerDisplayMode ChartMarkerDisplayMode
	StackMode         ChartStackMode
	YAxis             ChartYAxisSettings
}

type ChartDonutSettings struct {
	// OptionalFields
	CenterLabelAlwaysSlice string
	CenterLabelDisplayMode ChartDonutCenterLabelDisplayMode
	CenterLabelSize        ChartDonutCenterLabelSize
	CenterLabelSliceColor  string
	CenterLabelValueColor  string
	CenterLabelValueSuffix string
	HolePercent            uint8
	SliceIdentityColors    map[string]string
}

type ChartSettings struct {
	Type   ChartType
	XAxis  ChartXAxisSettings
	Series []ChartSeriesSettings

	// OptionalFields
	AriaLabel           string
	Cartesian           ChartCartesianSettings
	Donut               ChartDonutSettings
	EmptyState          templ.Component
	GapColor            string
	HeightClass         string
	Legend              ChartLegendSettings
	MaxRows             uint
	Rows                any
	RowsOneWayStatePath string
	Tooltip             ChartTooltipSettings
}

type chartSeriesConfig struct {
	RowKey        string             `json:"rowKey"`
	Label         string             `json:"label"`
	IdentityColor string             `json:"identityColor"`
	LineDashStyle ChartLineDashStyle `json:"lineDashStyle"`
}

type chartXAxisConfig struct {
	RowKey             string         `json:"rowKey"`
	BaselineColor      string         `json:"baselineColor"`
	GridLineColor      string         `json:"gridLineColor"`
	IsGridLineVisible  bool           `json:"isGridLineVisible"`
	IsHidden           bool           `json:"isHidden"`
	IsTickStartEndOnly bool           `json:"isTickStartEndOnly"`
	ScaleType          ChartScaleType `json:"scaleType"`
	TickLabelColor     string         `json:"tickLabelColor"`
	TickLabelFormat    string         `json:"tickLabelFormat"`
	Title              string         `json:"title"`
	TitleColor         string         `json:"titleColor"`
}

type chartYAxisConfig struct {
	BaselineColor    string             `json:"baselineColor"`
	GridLineColor    string             `json:"gridLineColor"`
	IsGridLineHidden bool               `json:"isGridLineHidden"`
	IsHidden         bool               `json:"isHidden"`
	MaxValue         string             `json:"maxValue"`
	MinValue         string             `json:"minValue"`
	Position         ChartYAxisPosition `json:"position"`
	TickCount        uint               `json:"tickCount"`
	TickLabelColor   string             `json:"tickLabelColor"`
	TickLabelFormat  string             `json:"tickLabelFormat"`
	Title            string             `json:"title"`
	TitleColor       string             `json:"titleColor"`
}

type chartTooltipConfig struct {
	ContentFunc       string `json:"contentFunc"`
	CrosshairColor    string `json:"crosshairColor"`
	IsCrosshairHidden bool   `json:"isCrosshairHidden"`
	IsHidden          bool   `json:"isHidden"`
	TitleFormat       string `json:"titleFormat"`
	ValueFormat       string `json:"valueFormat"`
}

type chartCartesianConfig struct {
	AreaFillStyle     ChartAreaFillStyle     `json:"areaFillStyle"`
	MarkerDisplayMode ChartMarkerDisplayMode `json:"markerDisplayMode"`
	StackMode         ChartStackMode         `json:"stackMode"`
	YAxis             chartYAxisConfig       `json:"yAxis"`
}

type chartDonutConfig struct {
	CenterLabelAlwaysSlice    string                           `json:"centerLabelAlwaysSlice"`
	CenterLabelDisplayMode    ChartDonutCenterLabelDisplayMode `json:"centerLabelDisplayMode"`
	CenterLabelSizeMultiplier string                           `json:"centerLabelSizeMultiplier"`
	CenterLabelSliceColor     string                           `json:"centerLabelSliceColor"`
	CenterLabelValueColor     string                           `json:"centerLabelValueColor"`
	CenterLabelValueSuffix    string                           `json:"centerLabelValueSuffix"`
	HolePercent               uint8                            `json:"holePercent"`
	SliceIdentityColors       map[string]string                `json:"sliceIdentityColors"`
}

type chartConfig struct {
	ChartType     ChartType            `json:"chartType"`
	XAxis         chartXAxisConfig     `json:"xAxis"`
	Series        []chartSeriesConfig  `json:"series"`
	Tooltip       chartTooltipConfig   `json:"tooltip"`
	Cartesian     chartCartesianConfig `json:"cartesian"`
	Donut         chartDonutConfig     `json:"donut"`
	GapColor      string               `json:"gapColor"`
	Rows          any                  `json:"rows"`
	RowsStatePath string               `json:"rowsStatePath"`
	MaxRows       uint                 `json:"maxRows"`
}

func chartOptionResolver[Option ~string](
	optionValue Option, defaultValue Option, validValues []Option, errorKey string,
) (Option, error) {
	if optionValue == "" {
		return defaultValue, nil
	}
	if !slices.Contains(validValues, optionValue) {
		return optionValue, fmt.Errorf("%s: %s", errorKey, optionValue)
	}
	return optionValue, nil
}

func chartTextDefaulter(text, defaultText string) string {
	if text == "" {
		return defaultText
	}
	return text
}

func chartXAxisResolver(xAxis ChartXAxisSettings) (chartXAxisConfig, error) {
	if xAxis.RowKey == "" {
		return chartXAxisConfig{}, errors.New("ChartXAxisRowKeyMissing")
	}
	scaleType, scaleTypeErr := chartOptionResolver(
		xAxis.ScaleType, ChartScaleTypeCategory,
		[]ChartScaleType{ChartScaleTypeCategory, ChartScaleTypeNumber, ChartScaleTypeTime},
		"ChartInvalidXAxisScaleType",
	)
	if scaleTypeErr != nil {
		return chartXAxisConfig{}, scaleTypeErr
	}
	return chartXAxisConfig{
		RowKey:             xAxis.RowKey,
		BaselineColor:      chartTextDefaulter(xAxis.BaselineColor, "neutral-50/8"),
		GridLineColor:      chartTextDefaulter(xAxis.GridLineColor, "neutral-50/8"),
		IsGridLineVisible:  xAxis.IsGridLineVisible,
		IsHidden:           xAxis.IsHidden,
		IsTickStartEndOnly: xAxis.IsTickStartEndOnly,
		ScaleType:          scaleType,
		TickLabelColor:     chartTextDefaulter(xAxis.TickLabelColor, "neutral-400"),
		TickLabelFormat:    xAxis.TickLabelFormat,
		Title:              xAxis.Title,
		TitleColor:         chartTextDefaulter(xAxis.TitleColor, "neutral-300"),
	}, nil
}

func chartYValueBoundValidator(boundText string) error {
	if boundText == "" {
		return nil
	}
	_, parseErr := strconv.ParseFloat(boundText, 64)
	if parseErr != nil {
		return fmt.Errorf("ChartInvalidYAxisValueBound (%s): %w", boundText, parseErr)
	}
	return nil
}

func chartYAxisResolver(
	yAxis ChartYAxisSettings, stackMode ChartStackMode,
) (chartYAxisConfig, error) {
	position, positionErr := chartOptionResolver(
		yAxis.Position, ChartYAxisPositionLeft,
		[]ChartYAxisPosition{ChartYAxisPositionLeft, ChartYAxisPositionRight},
		"ChartInvalidYAxisPosition",
	)
	if positionErr != nil {
		return chartYAxisConfig{}, positionErr
	}
	for _, boundText := range []string{yAxis.MinValue, yAxis.MaxValue} {
		boundErr := chartYValueBoundValidator(boundText)
		if boundErr != nil {
			return chartYAxisConfig{}, boundErr
		}
	}
	defaultTickLabelFormat := chartDefaultValueFormat
	if stackMode == ChartStackModePercent {
		defaultTickLabelFormat = chartDefaultPercentFormat
	}
	tickCount := yAxis.TickCount
	if tickCount == 0 {
		tickCount = chartDefaultYTickCount
	}
	return chartYAxisConfig{
		BaselineColor:    chartTextDefaulter(yAxis.BaselineColor, "neutral-50/8"),
		GridLineColor:    chartTextDefaulter(yAxis.GridLineColor, "neutral-50/8"),
		IsGridLineHidden: yAxis.IsGridLineHidden,
		IsHidden:         yAxis.IsHidden,
		MaxValue:         yAxis.MaxValue,
		MinValue:         yAxis.MinValue,
		Position:         position,
		TickCount:        tickCount,
		TickLabelColor:   chartTextDefaulter(yAxis.TickLabelColor, "neutral-400"),
		TickLabelFormat:  chartTextDefaulter(yAxis.TickLabelFormat, defaultTickLabelFormat),
		Title:            yAxis.Title,
		TitleColor:       chartTextDefaulter(yAxis.TitleColor, "neutral-300"),
	}, nil
}

func chartSeriesResolver(series []ChartSeriesSettings) ([]chartSeriesConfig, error) {
	if len(series) == 0 {
		return nil, errors.New("ChartSeriesMissing")
	}
	resolvedSeries := []chartSeriesConfig{}
	for _, seriesEntry := range series {
		if seriesEntry.RowKey == "" {
			return nil, errors.New("ChartSeriesRowKeyMissing")
		}
		lineDashStyle, lineDashErr := chartOptionResolver(
			seriesEntry.LineDashStyle, ChartLineDashStyleSolid,
			[]ChartLineDashStyle{
				ChartLineDashStyleDashed, ChartLineDashStyleDotted, ChartLineDashStyleSolid,
			},
			"ChartInvalidSeriesLineDashStyle",
		)
		if lineDashErr != nil {
			return nil, lineDashErr
		}
		resolvedSeries = append(resolvedSeries, chartSeriesConfig{
			RowKey:        seriesEntry.RowKey,
			Label:         chartTextDefaulter(seriesEntry.Label, seriesEntry.RowKey),
			IdentityColor: seriesEntry.IdentityColor,
			LineDashStyle: lineDashStyle,
		})
	}
	return resolvedSeries, nil
}

func chartTooltipResolver(
	tooltip ChartTooltipSettings, xAxis chartXAxisConfig, yAxis chartYAxisConfig,
	stackMode ChartStackMode,
) chartTooltipConfig {
	defaultTitleFormat := xAxis.TickLabelFormat
	if xAxis.ScaleType == ChartScaleTypeTime {
		defaultTitleFormat = chartDefaultTimeTitleFormat
	}
	defaultValueFormat := yAxis.TickLabelFormat
	if stackMode == ChartStackModePercent {
		defaultValueFormat = chartDefaultValueFormat
	}
	return chartTooltipConfig{
		ContentFunc:       tooltip.ContentFunc,
		CrosshairColor:    chartTextDefaulter(tooltip.CrosshairColor, "neutral-50/30"),
		IsCrosshairHidden: tooltip.IsCrosshairHidden,
		IsHidden:          tooltip.IsHidden,
		TitleFormat:       chartTextDefaulter(tooltip.TitleFormat, defaultTitleFormat),
		ValueFormat:       chartTextDefaulter(tooltip.ValueFormat, defaultValueFormat),
	}
}

func chartCartesianResolver(
	cartesian ChartCartesianSettings,
) (chartCartesianConfig, error) {
	stackMode, stackModeErr := chartOptionResolver(
		cartesian.StackMode, ChartStackModeNone,
		[]ChartStackMode{ChartStackModeNone, ChartStackModePercent, ChartStackModeStacked},
		"ChartInvalidStackMode",
	)
	if stackModeErr != nil {
		return chartCartesianConfig{}, stackModeErr
	}
	areaFillStyle, areaFillErr := chartOptionResolver(
		cartesian.AreaFillStyle, ChartAreaFillStyleGradient,
		[]ChartAreaFillStyle{
			ChartAreaFillStyleGradient, ChartAreaFillStyleNone, ChartAreaFillStyleSolid,
		},
		"ChartInvalidAreaFillStyle",
	)
	if areaFillErr != nil {
		return chartCartesianConfig{}, areaFillErr
	}
	markerDisplayMode, markerDisplayErr := chartOptionResolver(
		cartesian.MarkerDisplayMode, ChartMarkerDisplayModeHover,
		[]ChartMarkerDisplayMode{
			ChartMarkerDisplayModeAlways, ChartMarkerDisplayModeHidden,
			ChartMarkerDisplayModeHover,
		},
		"ChartInvalidMarkerDisplayMode",
	)
	if markerDisplayErr != nil {
		return chartCartesianConfig{}, markerDisplayErr
	}
	yAxis, yAxisErr := chartYAxisResolver(cartesian.YAxis, stackMode)
	if yAxisErr != nil {
		return chartCartesianConfig{}, yAxisErr
	}
	return chartCartesianConfig{
		AreaFillStyle:     areaFillStyle,
		MarkerDisplayMode: markerDisplayMode,
		StackMode:         stackMode,
		YAxis:             yAxis,
	}, nil
}

func chartDonutResolver(donut ChartDonutSettings) (chartDonutConfig, error) {
	holePercent := donut.HolePercent
	if holePercent == 0 {
		holePercent = chartDefaultDonutHolePercent
	}
	if holePercent > chartMaxDonutHolePercent {
		return chartDonutConfig{}, fmt.Errorf(
			"ChartInvalidDonutHolePercent: %d", holePercent,
		)
	}
	centerLabelDisplayMode, centerLabelDisplayModeErr := chartOptionResolver(
		donut.CenterLabelDisplayMode, ChartDonutCenterLabelDisplayModeHidden,
		[]ChartDonutCenterLabelDisplayMode{
			ChartDonutCenterLabelDisplayModeAlways, ChartDonutCenterLabelDisplayModeHidden,
			ChartDonutCenterLabelDisplayModeHover,
		},
		"ChartInvalidDonutCenterLabelDisplayMode",
	)
	if centerLabelDisplayModeErr != nil {
		return chartDonutConfig{}, centerLabelDisplayModeErr
	}
	centerLabelSize, centerLabelSizeErr := chartOptionResolver(
		donut.CenterLabelSize, ChartDonutCenterLabelSizeMd,
		[]ChartDonutCenterLabelSize{
			ChartDonutCenterLabelSizeXs, ChartDonutCenterLabelSizeSm,
			ChartDonutCenterLabelSizeMd, ChartDonutCenterLabelSizeLg,
			ChartDonutCenterLabelSizeXl,
		},
		"ChartInvalidDonutCenterLabelSize",
	)
	if centerLabelSizeErr != nil {
		return chartDonutConfig{}, centerLabelSizeErr
	}
	sliceIdentityColors := donut.SliceIdentityColors
	if sliceIdentityColors == nil {
		sliceIdentityColors = map[string]string{}
	}
	return chartDonutConfig{
		CenterLabelAlwaysSlice:    donut.CenterLabelAlwaysSlice,
		CenterLabelDisplayMode:    centerLabelDisplayMode,
		CenterLabelSizeMultiplier: chartCenterLabelSizeMultiplierResolver(centerLabelSize),
		CenterLabelSliceColor:     chartTextDefaulter(donut.CenterLabelSliceColor, chartDefaultCenterLabelSliceColor),
		CenterLabelValueColor:     chartTextDefaulter(donut.CenterLabelValueColor, chartDefaultCenterLabelValueColor),
		CenterLabelValueSuffix:    donut.CenterLabelValueSuffix,
		HolePercent:               holePercent,
		SliceIdentityColors:       sliceIdentityColors,
	}, nil
}

func (componentSettings ChartSettings) configResolver() (chartConfig, error) {
	chartType, chartTypeErr := chartOptionResolver(
		componentSettings.Type, ChartTypeLine,
		[]ChartType{ChartTypeArea, ChartTypeBar, ChartTypeDonut, ChartTypeLine},
		"ChartInvalidType",
	)
	if chartTypeErr != nil {
		return chartConfig{}, chartTypeErr
	}
	xAxis, xAxisErr := chartXAxisResolver(componentSettings.XAxis)
	if xAxisErr != nil {
		return chartConfig{}, xAxisErr
	}
	series, seriesErr := chartSeriesResolver(componentSettings.Series)
	if seriesErr != nil {
		return chartConfig{}, seriesErr
	}
	cartesian, cartesianErr := chartCartesianResolver(componentSettings.Cartesian)
	if cartesianErr != nil {
		return chartConfig{}, cartesianErr
	}
	donut, donutErr := chartDonutResolver(componentSettings.Donut)
	if donutErr != nil {
		return chartConfig{}, donutErr
	}
	rows := componentSettings.Rows
	if rows == nil {
		rows = []any{}
	}
	return chartConfig{
		ChartType: chartType,
		XAxis:     xAxis,
		Series:    series,
		Tooltip: chartTooltipResolver(
			componentSettings.Tooltip, xAxis, cartesian.YAxis, cartesian.StackMode,
		),
		Cartesian:     cartesian,
		Donut:         donut,
		GapColor:      chartTextDefaulter(componentSettings.GapColor, "neutral-900"),
		Rows:          rows,
		RowsStatePath: componentSettings.RowsOneWayStatePath,
		MaxRows:       componentSettings.MaxRows,
	}, nil
}

func (componentSettings ChartSettings) legendResolver() (ChartLegendSettings, error) {
	resolvedLegend := componentSettings.Legend
	displayMode, displayModeErr := chartOptionResolver(
		resolvedLegend.DisplayMode, ChartLegendDisplayModeAuto,
		[]ChartLegendDisplayMode{
			ChartLegendDisplayModeAlways, ChartLegendDisplayModeAuto,
			ChartLegendDisplayModeHidden,
		},
		"ChartInvalidLegendDisplayMode",
	)
	if displayModeErr != nil {
		return resolvedLegend, displayModeErr
	}
	resolvedLegend.DisplayMode = displayMode
	position, positionErr := chartOptionResolver(
		resolvedLegend.Position, ChartLegendPositionTop,
		[]ChartLegendPosition{ChartLegendPositionBottom, ChartLegendPositionTop},
		"ChartInvalidLegendPosition",
	)
	if positionErr != nil {
		return resolvedLegend, positionErr
	}
	resolvedLegend.Position = position
	alignment, alignmentErr := chartOptionResolver(
		resolvedLegend.Alignment, ChartLegendAlignmentStart,
		[]ChartLegendAlignment{
			ChartLegendAlignmentCenter, ChartLegendAlignmentEnd, ChartLegendAlignmentStart,
		},
		"ChartInvalidLegendAlignment",
	)
	if alignmentErr != nil {
		return resolvedLegend, alignmentErr
	}
	resolvedLegend.Alignment = alignment
	return resolvedLegend, nil
}

func (resolvedLegend ChartLegendSettings) visibilityExpressionResolver() string {
	if resolvedLegend.DisplayMode == ChartLegendDisplayModeAlways {
		return "chartLegendEntries.length > 0"
	}
	return "chartLegendEntries.length > 1"
}

func (resolvedLegend ChartLegendSettings) listClassesResolver() string {
	justifyClassByAlignment := map[ChartLegendAlignment]string{
		ChartLegendAlignmentCenter: "justify-center",
		ChartLegendAlignmentEnd:    "justify-end",
		ChartLegendAlignmentStart:  "justify-start",
	}
	orderClass := "order-first"
	if resolvedLegend.Position == ChartLegendPositionBottom {
		orderClass = "order-last"
	}
	return uiToolset.NonEmptyClassJoiner([]string{
		"flex flex-row flex-wrap text-xs",
		orderClass,
		justifyClassByAlignment[resolvedLegend.Alignment],
		uiToolset.GapClassResolver(resolvedLegend.GapSize, "gap-x-4 gap-y-1"),
		uiToolset.TextColorClassResolver(resolvedLegend.TextColor, "text-neutral-300"),
	})
}

func (resolvedLegend ChartLegendSettings) swatchClassesResolver() string {
	return "h-2.5 w-2.5 shrink-0 " +
		uiToolset.ShapeClassResolver(resolvedLegend.SwatchShape, "rounded-full")
}

func (componentSettings ChartSettings) tooltipClassesResolver() string {
	tooltip := componentSettings.Tooltip
	ringClasses := uiToolset.CompactRingClassResolver(tooltip.RingColor, tooltip.RingThickness)
	if ringClasses == "" {
		ringClasses = "ring-1 ring-neutral-50/10"
	}
	return uiToolset.NonEmptyClassJoiner([]string{
		"pointer-events-none fixed z-50 flex flex-col gap-1 text-xs",
		uiToolset.BackgroundColorClassResolver(tooltip.BackgroundColor, "bg-neutral-800/95"),
		uiToolset.BorderRadiusClassResolver(tooltip.BorderRadius, "rounded-md"),
		uiToolset.PaddingClassResolver(tooltip.PaddingSize, "p-2"),
		uiToolset.ShadowClassResolver(tooltip.ShadowSize, "shadow-lg"),
		ringClasses,
		chartTextDefaulter(tooltip.MinWidthClass, "min-w-36"),
		tooltip.MaxWidthClass,
	})
}

func (resolvedConfig chartConfig) colorTokensCollector() []string {
	colorTokens := []string{
		resolvedConfig.GapColor,
		resolvedConfig.Tooltip.CrosshairColor,
		resolvedConfig.XAxis.TitleColor,
		resolvedConfig.XAxis.TickLabelColor,
		resolvedConfig.XAxis.BaselineColor,
		resolvedConfig.XAxis.GridLineColor,
		resolvedConfig.Cartesian.YAxis.TitleColor,
		resolvedConfig.Cartesian.YAxis.TickLabelColor,
		resolvedConfig.Cartesian.YAxis.BaselineColor,
		resolvedConfig.Cartesian.YAxis.GridLineColor,
	}
	if resolvedConfig.Donut.CenterLabelDisplayMode != ChartDonutCenterLabelDisplayModeHidden {
		colorTokens = append(
			colorTokens,
			resolvedConfig.Donut.CenterLabelSliceColor,
			resolvedConfig.Donut.CenterLabelValueColor,
		)
	}
	for _, seriesEntry := range resolvedConfig.Series {
		colorTokens = append(colorTokens, seriesEntry.IdentityColor)
	}
	for _, sliceColorToken := range resolvedConfig.Donut.SliceIdentityColors {
		colorTokens = append(colorTokens, sliceColorToken)
	}
	uniqueTokens := []string{}
	for _, colorToken := range colorTokens {
		if colorToken == "" || slices.Contains(uniqueTokens, colorToken) {
			continue
		}
		uniqueTokens = append(uniqueTokens, colorToken)
	}
	slices.Sort(uniqueTokens)
	return uniqueTokens
}

func (resolvedConfig chartConfig) seriesLabelsJoiner() string {
	seriesLabels := []string{}
	for _, seriesEntry := range resolvedConfig.Series {
		seriesLabels = append(seriesLabels, seriesEntry.Label)
	}
	return strings.Join(seriesLabels, ", ")
}

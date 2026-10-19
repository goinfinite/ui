package uiDisplay

import (
	_ "embed"
	"errors"
	"fmt"
	"math"
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

const (
	ChartProgressSizeXs string = "xs"
	ChartProgressSizeSm string = "sm"
	ChartProgressSizeMd string = "md"
)

const chartProgressIndeterminateStyle string = "<style type=\"text/css\">@keyframes ui-progress-bar-indeterminate{0%{left:-40%}100%{left:100%}}</style>"

var chartProgressIndeterminateStyleOnce = templ.NewOnceHandle(
	templ.WithComponent(templ.Raw(chartProgressIndeterminateStyle)),
)

type ChartType string

const (
	ChartTypeArea          ChartType = "area"
	ChartTypeBar           ChartType = "bar"
	ChartTypeDonut         ChartType = "donut"
	ChartTypeGauge         ChartType = "gauge"
	ChartTypeHorizontalBar ChartType = "horizontalBar"
	ChartTypeLine          ChartType = "line"
	ChartTypeParliament    ChartType = "parliament"
	ChartTypePie           ChartType = "pie"
	ChartTypeProgress      ChartType = "progress"
	ChartTypeRadialBar     ChartType = "radialBar"
	ChartTypeSunburst      ChartType = "sunburst"
	ChartTypeWaffle        ChartType = "waffle"
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

type ChartCenterLabelDisplayMode string

const (
	ChartCenterLabelDisplayModeAlways ChartCenterLabelDisplayMode = "always"
	ChartCenterLabelDisplayModeHidden ChartCenterLabelDisplayMode = "hidden"
	ChartCenterLabelDisplayModeHover  ChartCenterLabelDisplayMode = "hover"
)

type ChartCenterLabelSize string

const (
	ChartCenterLabelSizeXs ChartCenterLabelSize = "xs"
	ChartCenterLabelSizeSm ChartCenterLabelSize = "sm"
	ChartCenterLabelSizeMd ChartCenterLabelSize = "md"
	ChartCenterLabelSizeLg ChartCenterLabelSize = "lg"
	ChartCenterLabelSizeXl ChartCenterLabelSize = "xl"
)

const (
	chartDefaultTimeTitleFormat    string = "%b %d, %H:%M:%S"
	chartDefaultValueFormat        string = ",.2~f"
	chartDefaultPercentFormat      string = ".0%"
	chartDefaultSliceHolePercent   uint8  = 60
	chartMaxSliceHolePercent       uint8  = 90
	chartDefaultYTickCount         uint   = 5
	chartDefaultWaffleCellCount    uint   = 100
	chartDefaultWaffleColumnCount  uint   = 10
	chartDefaultParliamentRowCount uint   = 6

	chartDefaultWaffleCellGapPercent     float64 = 4
	chartMaxWaffleCellGapPercent         float64 = 60
	chartDefaultParliamentDotGapPercent  float64 = 25
	chartMaxParliamentDotGapPercent      float64 = 95
	chartDefaultParliamentDotSizePercent float64 = 100
	chartMinParliamentDotSizePercent     float64 = 10
	chartMaxParliamentDotSizePercent     float64 = 200
	chartDefaultCenterLabelSliceColor    string  = "neutral-400"
	chartDefaultCenterLabelValueColor    string  = "neutral-50"
	chartDefaultSliceLabelPillColor      string  = "neutral-950/70"
)

func chartCenterLabelSizeMultiplierResolver(
	centerLabelSize ChartCenterLabelSize,
) string {
	switch centerLabelSize {
	case ChartCenterLabelSizeXs:
		return "0.75"
	case ChartCenterLabelSizeSm:
		return "0.875"
	case ChartCenterLabelSizeMd:
		return "1"
	case ChartCenterLabelSizeLg:
		return "1.15"
	case ChartCenterLabelSizeXl:
		return "1.3"
	}
	return "1"
}

type ChartProgressThresholdSettings struct {
	MinimumPercent  float64
	BackgroundColor string
}

type ChartProgressSettings struct {
	// OptionalFields
	DefaultColor           string
	IsIndeterminate        bool
	IsStriped              bool
	Label                  string
	LabelBackgroundColor   string
	LabelOneWayStatePath   string
	LabelPrefix            string
	LabelSuffix            string
	LabelTextColor         string
	LabelTextSizeClass     string
	MaxHeightClass         string
	MaxValue               float64
	MaxWidthClass          string
	MinHeightClass         string
	MinValue               float64
	MinWidthClass          string
	ShouldShowPercentLabel bool
	Size                   string
	StripeColor            string
	Thresholds             []ChartProgressThresholdSettings
	TrackBackgroundColor   string
	Value                  float64
	ValueOneWayStatePath   string
}

type chartProgressThreshold struct {
	MinimumPercent  float64 `json:"minimumPercent"`
	BackgroundColor string  `json:"backgroundColor"`
}

type chartValueBounds struct {
	Min float64
	Max float64
}

func chartValueBoundsNormalizer(minValue, maxValue float64) chartValueBounds {
	if maxValue > minValue {
		return chartValueBounds{Min: minValue, Max: maxValue}
	}
	return chartValueBounds{Min: minValue, Max: minValue + 100}
}

func chartValuePercentResolver(value, minValue, maxValue float64) float64 {
	valueSpan := maxValue - minValue
	if valueSpan <= 0 {
		return 0
	}
	ratio := (value - minValue) / valueSpan
	clampedRatio := math.Min(math.Max(ratio, 0), 1)
	return clampedRatio * 100
}

func (componentSettings ChartProgressSettings) boundsResolver() chartValueBounds {
	return chartValueBoundsNormalizer(componentSettings.MinValue, componentSettings.MaxValue)
}

func (componentSettings ChartProgressSettings) percentResolver() float64 {
	progressBounds := componentSettings.boundsResolver()
	return chartValuePercentResolver(
		componentSettings.Value, progressBounds.Min, progressBounds.Max,
	)
}

func (componentSettings ChartProgressSettings) clampedValueResolver() float64 {
	progressBounds := componentSettings.boundsResolver()
	return math.Min(math.Max(componentSettings.Value, progressBounds.Min), progressBounds.Max)
}

func (componentSettings ChartProgressSettings) thresholdsResolver() []chartProgressThreshold {
	if len(componentSettings.Thresholds) > 0 {
		resolvedThresholds := make(
			[]chartProgressThreshold, 0, len(componentSettings.Thresholds),
		)
		for _, threshold := range componentSettings.Thresholds {
			resolvedThresholds = append(
				resolvedThresholds, chartProgressThreshold(threshold),
			)
		}
		return resolvedThresholds
	}
	return []chartProgressThreshold{
		{MinimumPercent: 90, BackgroundColor: "red-800"},
		{MinimumPercent: 85, BackgroundColor: "red-600"},
		{MinimumPercent: 75, BackgroundColor: "amber-600"},
		{MinimumPercent: 50, BackgroundColor: "amber-300"},
	}
}

func (componentSettings ChartProgressSettings) defaultColorResolver() string {
	if componentSettings.DefaultColor != "" {
		return componentSettings.DefaultColor
	}
	return "cyan-400"
}

func (componentSettings ChartProgressSettings) fillColorResolver(percent float64) string {
	for _, threshold := range componentSettings.thresholdsResolver() {
		if percent >= threshold.MinimumPercent {
			return "bg-" + threshold.BackgroundColor
		}
	}
	return "bg-" + componentSettings.defaultColorResolver()
}

type chartProgressSizeTokens struct {
	TrackHeight   string
	LabelTextSize string
}

func chartProgressSizeTokensResolver(size string) chartProgressSizeTokens {
	switch size {
	case ChartProgressSizeXs:
		return chartProgressSizeTokens{TrackHeight: "h-4", LabelTextSize: "text-[0.625rem]"}
	case ChartProgressSizeMd:
		return chartProgressSizeTokens{TrackHeight: "h-8", LabelTextSize: "text-sm"}
	}
	return chartProgressSizeTokens{TrackHeight: "h-6", LabelTextSize: "text-xs"}
}

func chartProgressLabelTextSizeResolver(configuredClass, size string) string {
	if configuredClass != "" {
		return configuredClass
	}
	return chartProgressSizeTokensResolver(size).LabelTextSize
}

func chartProgressPercentTextResolver(percent float64) string {
	return chartRoundedValueTextResolver(percent) + "%"
}

func chartRoundedValueTextResolver(value float64) string {
	roundedValue := math.Round(value*10) / 10
	return strconv.FormatFloat(roundedValue, 'f', -1, 64)
}

const (
	chartGaugeOuterRadius            float64 = 90
	chartGaugeDefaultArcThickness    float64 = 16
	chartGaugeMinArcThickness        float64 = 4
	chartGaugeMaxArcThickness        float64 = 36
	chartGaugeDefaultNeedleLength    float64 = 55
	chartGaugeMinNeedleLength        float64 = 10
	chartGaugeMaxNeedleLength        float64 = 90
	chartGaugeDefaultNeedleThickness float64 = 3
	chartGaugeMinNeedleThickness     float64 = 1
	chartGaugeMaxNeedleThickness     float64 = 12
	chartGaugeDefaultPivotRadius     float64 = 5
	chartGaugeMaxPivotRadius         float64 = 20
	chartGaugeDefaultValueTextSize   float64 = 24
	chartGaugeDefaultLabelTextSize   float64 = 12
	chartGaugeMinTextSize            float64 = 6
	chartGaugeMaxTextSize            float64 = 72
	chartGaugeMaxTextOffset          float64 = 80
	chartGaugeTickLength             float64 = 8
	chartGaugeTickLabelGap           float64 = 10
	chartGaugeTickGap                float64 = 2
)

type ChartGaugeBandSettings struct {
	MinimumPercent float64
	MaximumPercent float64

	// OptionalFields
	Color string
}

type ChartGaugeTickPosition string

const (
	ChartGaugeTickPositionInside  ChartGaugeTickPosition = "inside"
	ChartGaugeTickPositionOutside ChartGaugeTickPosition = "outside"
)

type ChartGaugeSettings struct {
	// OptionalFields
	ArcThickness    float64
	Bands           []ChartGaugeBandSettings
	IsSemicircle    bool
	Label           string
	LabelColor      string
	LabelTextSize   float64
	MaxValue        float64
	MinValue        float64
	NeedleColor     string
	NeedleLength    float64
	NeedleThickness float64
	PivotRadius     float64
	TextOffset      float64
	TickColor       string
	TickCount       uint
	TickLabelColor  string
	TickPosition    ChartGaugeTickPosition
	TrackColor      string
	Value           float64
	ValueColor      string
	ValueTextSize   float64
}

type ChartWaffleCellShape string

const (
	ChartWaffleCellShapeCircle ChartWaffleCellShape = "circle"
	ChartWaffleCellShapeSquare ChartWaffleCellShape = "square"
)

type ChartWaffleSettings struct {
	// OptionalFields
	CellCount      uint
	CellGapPercent float64
	CellShape      ChartWaffleCellShape
	ColumnCount    uint
}

type ChartSunburstSettings struct {
	ChildRowKey string
}

type ChartDotShape string

const (
	ChartDotShapeCircle       ChartDotShape = "circle"
	ChartDotShapeCross        ChartDotShape = "cross"
	ChartDotShapeDiamond      ChartDotShape = "diamond"
	ChartDotShapeSquare       ChartDotShape = "square"
	ChartDotShapeTriangleDown ChartDotShape = "triangleDown"
	ChartDotShapeTriangleUp   ChartDotShape = "triangleUp"
)

type ChartParliamentSettings struct {
	// OptionalFields
	DotGapPercent  float64
	DotShape       ChartDotShape
	DotSizePercent float64
	RowCount       uint
}

type chartGaugeBand struct {
	ArcPath      string
	IsFullCircle bool
	Color        string
}

type chartGaugeTextSettings struct {
	LabelSize float64
	Offset    float64
	ValueSize float64
}

type chartGaugeTextLayout struct {
	ValueY string
	LabelY string
	Bottom float64
}

type chartGaugeLayout struct {
	ArcRadius       float64
	ArcThickness    float64
	NeedleLength    float64
	NeedleThickness float64
	PivotRadius     float64
}

type chartGaugeTick struct {
	LineX1 string
	LineY1 string
	LineX2 string
	LineY2 string
	LabelX string
	LabelY string
	Label  string
}

func chartGaugeValueDefaulter(value, defaultValue float64) float64 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func chartGaugeRangeValidator(
	value, minValue, maxValue float64, errorKey string,
) error {
	if math.IsNaN(value) || value < minValue || value > maxValue {
		return fmt.Errorf("%s: %v", errorKey, value)
	}
	return nil
}

func (componentSettings ChartGaugeSettings) layoutResolver() (chartGaugeLayout, error) {
	arcThickness := chartGaugeValueDefaulter(
		componentSettings.ArcThickness, chartGaugeDefaultArcThickness,
	)
	arcThicknessErr := chartGaugeRangeValidator(
		arcThickness, chartGaugeMinArcThickness, chartGaugeMaxArcThickness,
		"ChartInvalidGaugeArcThickness",
	)
	if arcThicknessErr != nil {
		return chartGaugeLayout{}, arcThicknessErr
	}
	needleLength := chartGaugeValueDefaulter(
		componentSettings.NeedleLength, chartGaugeDefaultNeedleLength,
	)
	needleLengthErr := chartGaugeRangeValidator(
		needleLength, chartGaugeMinNeedleLength, chartGaugeMaxNeedleLength,
		"ChartInvalidGaugeNeedleLength",
	)
	if needleLengthErr != nil {
		return chartGaugeLayout{}, needleLengthErr
	}
	needleThickness := chartGaugeValueDefaulter(
		componentSettings.NeedleThickness, chartGaugeDefaultNeedleThickness,
	)
	needleThicknessErr := chartGaugeRangeValidator(
		needleThickness, chartGaugeMinNeedleThickness, chartGaugeMaxNeedleThickness,
		"ChartInvalidGaugeNeedleThickness",
	)
	if needleThicknessErr != nil {
		return chartGaugeLayout{}, needleThicknessErr
	}
	pivotRadius := chartGaugeValueDefaulter(
		componentSettings.PivotRadius, chartGaugeDefaultPivotRadius,
	)
	pivotRadiusErr := chartGaugeRangeValidator(
		pivotRadius, 0, chartGaugeMaxPivotRadius, "ChartInvalidGaugePivotRadius",
	)
	if pivotRadiusErr != nil {
		return chartGaugeLayout{}, pivotRadiusErr
	}
	return chartGaugeLayout{
		ArcRadius:       chartGaugeOuterRadius - arcThickness/2,
		ArcThickness:    arcThickness,
		NeedleLength:    needleLength,
		NeedleThickness: needleThickness,
		PivotRadius:     pivotRadius,
	}, nil
}

func chartGaugeTextSizeResolver(
	size, defaultSize float64, errorKey string,
) (float64, error) {
	if size == 0 {
		return defaultSize, nil
	}
	if math.IsNaN(size) || size < chartGaugeMinTextSize || size > chartGaugeMaxTextSize {
		return 0, fmt.Errorf("%s: %v", errorKey, size)
	}
	return size, nil
}

func (componentSettings ChartGaugeSettings) textSettingsResolver() (
	chartGaugeTextSettings, error,
) {
	valueSize, valueSizeErr := chartGaugeTextSizeResolver(
		componentSettings.ValueTextSize, chartGaugeDefaultValueTextSize,
		"ChartInvalidGaugeValueTextSize",
	)
	if valueSizeErr != nil {
		return chartGaugeTextSettings{}, valueSizeErr
	}
	labelSize, labelSizeErr := chartGaugeTextSizeResolver(
		componentSettings.LabelTextSize, chartGaugeDefaultLabelTextSize,
		"ChartInvalidGaugeLabelTextSize",
	)
	if labelSizeErr != nil {
		return chartGaugeTextSettings{}, labelSizeErr
	}
	textOffset := componentSettings.TextOffset
	if math.IsNaN(textOffset) ||
		textOffset < 0 || textOffset > chartGaugeMaxTextOffset {
		return chartGaugeTextSettings{}, fmt.Errorf(
			"ChartInvalidGaugeTextOffset: %v", textOffset,
		)
	}
	return chartGaugeTextSettings{
		LabelSize: labelSize,
		Offset:    textOffset,
		ValueSize: valueSize,
	}, nil
}

func (componentSettings ChartGaugeSettings) tickPositionResolver() (
	ChartGaugeTickPosition, error,
) {
	return chartOptionResolver(
		componentSettings.TickPosition, ChartGaugeTickPositionOutside,
		[]ChartGaugeTickPosition{
			ChartGaugeTickPositionInside, ChartGaugeTickPositionOutside,
		},
		"ChartInvalidGaugeTickPosition",
	)
}

func (componentSettings ChartGaugeSettings) boundsResolver() chartValueBounds {
	return chartValueBoundsNormalizer(componentSettings.MinValue, componentSettings.MaxValue)
}

func (componentSettings ChartGaugeSettings) percentResolver() float64 {
	gaugeBounds := componentSettings.boundsResolver()
	return chartValuePercentResolver(
		componentSettings.Value, gaugeBounds.Min, gaugeBounds.Max,
	)
}

func (componentSettings ChartGaugeSettings) clampedValueResolver() float64 {
	gaugeBounds := componentSettings.boundsResolver()
	return math.Min(math.Max(componentSettings.Value, gaugeBounds.Min), gaugeBounds.Max)
}

func chartGaugeAngleResolver(percent float64, isSemicircle bool) float64 {
	if isSemicircle {
		return -90 + percent*1.8
	}
	return percent * 3.6
}

func chartGaugePointResolver(angleDeg float64, radius float64) (float64, float64) {
	angleRad := angleDeg * math.Pi / 180
	return 100 + radius*math.Sin(angleRad), 100 - radius*math.Cos(angleRad)
}

func chartGaugeArcPathResolver(startAngleDeg, endAngleDeg, radius float64) string {
	startX, startY := chartGaugePointResolver(startAngleDeg, radius)
	endX, endY := chartGaugePointResolver(endAngleDeg, radius)
	largeArcFlag := 0
	if endAngleDeg-startAngleDeg > 180 {
		largeArcFlag = 1
	}
	return fmt.Sprintf(
		"M%.2f %.2f A%.2f %.2f 0 %d 1 %.2f %.2f",
		startX, startY, radius, radius, largeArcFlag, endX, endY,
	)
}

func chartGaugeBandShapeResolver(
	minimumPercent, maximumPercent float64, isSemicircle bool, arcRadius float64,
) chartGaugeBand {
	isFullCircle := !isSemicircle && minimumPercent == 0 && maximumPercent == 100
	gaugeBand := chartGaugeBand{IsFullCircle: isFullCircle}
	if !isFullCircle {
		gaugeBand.ArcPath = chartGaugeArcPathResolver(
			chartGaugeAngleResolver(minimumPercent, isSemicircle),
			chartGaugeAngleResolver(maximumPercent, isSemicircle),
			arcRadius,
		)
	}
	return gaugeBand
}

func (componentSettings ChartGaugeSettings) trackResolver(
	layout chartGaugeLayout,
) chartGaugeBand {
	return chartGaugeBandShapeResolver(
		0, 100, componentSettings.IsSemicircle, layout.ArcRadius,
	)
}

func (componentSettings ChartGaugeSettings) bandsResolver(
	layout chartGaugeLayout,
) ([]chartGaugeBand, error) {
	resolvedBands := []chartGaugeBand{}
	for _, band := range componentSettings.Bands {
		isValidBand := band.MinimumPercent >= 0 && band.MaximumPercent <= 100 &&
			band.MinimumPercent < band.MaximumPercent
		if !isValidBand {
			return nil, fmt.Errorf(
				"ChartInvalidGaugeBandPercent: %v-%v", band.MinimumPercent, band.MaximumPercent,
			)
		}
		resolvedBand := chartGaugeBandShapeResolver(
			band.MinimumPercent, band.MaximumPercent, componentSettings.IsSemicircle,
			layout.ArcRadius,
		)
		resolvedBand.Color = chartTextDefaulter(band.Color, "neutral-700")
		resolvedBands = append(resolvedBands, resolvedBand)
	}
	return resolvedBands, nil
}

func chartGaugeNeedleTransformResolver(percent float64, isSemicircle bool) string {
	needleAngle := chartGaugeAngleResolver(percent, isSemicircle)
	roundedAngle := math.Round(needleAngle*100) / 100
	return fmt.Sprintf(
		"rotate(%s 100 100)", strconv.FormatFloat(roundedAngle, 'f', -1, 64),
	)
}

func chartGaugeViewBoxResolver(
	isSemicircle bool, hasOutsideTicks bool, textBottom, textWidth float64,
) string {
	minX, minY := 0.0, 0.0
	width, height := 200.0, 150.0
	if !isSemicircle {
		height = 200
	}
	if hasOutsideTicks {
		minX, minY = -20, -20
		width = 240
		if isSemicircle {
			height = 190
		} else {
			height = 240
		}
	}
	textWidthWithMargin := textWidth + 4
	if textWidthWithMargin > width {
		widthExpansion := textWidthWithMargin - width
		minX -= widthExpansion / 2
		width += widthExpansion
	}
	viewBoxBottom := minY + height
	textBottomWithMargin := textBottom + 2
	if textBottomWithMargin > viewBoxBottom {
		height += textBottomWithMargin - viewBoxBottom
	}
	return chartGaugeCoordinateTextResolver(minX) + " " +
		chartGaugeCoordinateTextResolver(minY) + " " +
		chartGaugeCoordinateTextResolver(width) + " " +
		chartGaugeCoordinateTextResolver(height)
}

func chartGaugeTextLayoutResolver(
	hasLabel, isSemicircle bool, textSettings chartGaugeTextSettings,
) chartGaugeTextLayout {
	valueY := 134.0
	labelY := 156.0
	if isSemicircle {
		valueY = 122
		labelY = 142
	}
	if !hasLabel {
		valueY = 134
		if isSemicircle {
			valueY = 124
		}
	}
	valueY += textSettings.Offset
	labelY += textSettings.Offset
	valueBottom := valueY + textSettings.ValueSize/2
	textLayout := chartGaugeTextLayout{
		ValueY: chartGaugeCoordinateTextResolver(valueY),
		Bottom: valueBottom,
	}
	if !hasLabel {
		return textLayout
	}
	labelBottom := labelY + textSettings.LabelSize/2
	textLayout.LabelY = chartGaugeCoordinateTextResolver(labelY)
	textLayout.Bottom = math.Max(valueBottom, labelBottom)
	return textLayout
}

func chartGaugeTextWidthEstimate(text string, fontSize float64) float64 {
	return float64(len([]rune(text))) * fontSize * 0.62
}

func chartGaugeTextWidthResolver(
	valueText, label string, hasLabel bool, textSettings chartGaugeTextSettings,
) float64 {
	textWidth := chartGaugeTextWidthEstimate(valueText, textSettings.ValueSize)
	if !hasLabel {
		return textWidth
	}
	return math.Max(
		textWidth, chartGaugeTextWidthEstimate(label, textSettings.LabelSize),
	)
}

func chartGaugeTickLabelShouldRender(
	tick chartGaugeTick, tickPosition ChartGaugeTickPosition, isSemicircle bool,
) bool {
	if isSemicircle || tickPosition != ChartGaugeTickPositionInside {
		return true
	}
	labelY, parseErr := strconv.ParseFloat(tick.LabelY, 64)
	if parseErr != nil {
		return true
	}
	return labelY <= 100
}

func chartGaugeCoordinateTextResolver(coordinate float64) string {
	roundedCoordinate := math.Round(coordinate*100) / 100
	return strconv.FormatFloat(roundedCoordinate, 'f', -1, 64)
}

func (componentSettings ChartGaugeSettings) ticksResolver(
	layout chartGaugeLayout, tickPosition ChartGaugeTickPosition,
) []chartGaugeTick {
	tickCount := componentSettings.TickCount
	if tickCount == 0 {
		return []chartGaugeTick{}
	}
	gaugeBounds := componentSettings.boundsResolver()
	arcOuterRadius := layout.ArcRadius + layout.ArcThickness/2
	arcInnerRadius := layout.ArcRadius - layout.ArcThickness/2
	tickStartRadius := arcInnerRadius - chartGaugeTickGap
	tickEndRadius := tickStartRadius - chartGaugeTickLength
	tickLabelRadius := tickEndRadius - chartGaugeTickLabelGap
	if tickPosition == ChartGaugeTickPositionOutside {
		tickStartRadius = arcOuterRadius + chartGaugeTickGap
		tickEndRadius = tickStartRadius + chartGaugeTickLength
		tickLabelRadius = tickEndRadius + chartGaugeTickLabelGap
	}
	lastTickIndex := tickCount
	if !componentSettings.IsSemicircle {
		lastTickIndex = tickCount - 1
	}
	resolvedTicks := make([]chartGaugeTick, 0, lastTickIndex+1)
	for tickIndex := uint(0); tickIndex <= lastTickIndex; tickIndex++ {
		tickRatio := float64(tickIndex) / float64(tickCount)
		tickAngle := chartGaugeAngleResolver(
			tickRatio*100, componentSettings.IsSemicircle,
		)
		lineX1, lineY1 := chartGaugePointResolver(tickAngle, tickStartRadius)
		lineX2, lineY2 := chartGaugePointResolver(tickAngle, tickEndRadius)
		labelX, labelY := chartGaugePointResolver(tickAngle, tickLabelRadius)
		tickValue := gaugeBounds.Min + (gaugeBounds.Max-gaugeBounds.Min)*tickRatio
		resolvedTicks = append(resolvedTicks, chartGaugeTick{
			LineX1: chartGaugeCoordinateTextResolver(lineX1),
			LineY1: chartGaugeCoordinateTextResolver(lineY1),
			LineX2: chartGaugeCoordinateTextResolver(lineX2),
			LineY2: chartGaugeCoordinateTextResolver(lineY2),
			LabelX: chartGaugeCoordinateTextResolver(labelX),
			LabelY: chartGaugeCoordinateTextResolver(labelY),
			Label:  chartRoundedValueTextResolver(tickValue),
		})
	}
	return resolvedTicks
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

type ChartSliceSettings struct {
	// OptionalFields
	CenterLabelAlwaysSlice        string
	CenterLabelDisplayMode        ChartCenterLabelDisplayMode
	CenterLabelSize               ChartCenterLabelSize
	CenterLabelSliceColor         string
	CenterLabelValueColor         string
	CenterLabelValueSuffix        string
	HolePercent                   uint8
	IsSliceLabelHidden            bool
	SliceIdentityColors           map[string]string
	SliceLabelColor               string
	SliceLabelPillBackgroundColor string
}

type ChartSettings struct {
	Type   ChartType
	XAxis  ChartXAxisSettings
	Series []ChartSeriesSettings

	// OptionalFields
	AriaLabel           string
	Cartesian           ChartCartesianSettings
	EmptyState          templ.Component
	GapColor            string
	Gauge               ChartGaugeSettings
	HeightClass         string
	Legend              ChartLegendSettings
	MaxRows             uint
	Parliament          ChartParliamentSettings
	Progress            ChartProgressSettings
	Rows                any
	RowsOneWayStatePath string
	Slices              ChartSliceSettings
	Sunburst            ChartSunburstSettings
	Tooltip             ChartTooltipSettings
	Waffle              ChartWaffleSettings
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

type chartSliceConfig struct {
	CenterLabelAlwaysSlice        string                      `json:"centerLabelAlwaysSlice"`
	CenterLabelDisplayMode        ChartCenterLabelDisplayMode `json:"centerLabelDisplayMode"`
	CenterLabelSizeMultiplier     string                      `json:"centerLabelSizeMultiplier"`
	CenterLabelSliceColor         string                      `json:"centerLabelSliceColor"`
	CenterLabelValueColor         string                      `json:"centerLabelValueColor"`
	CenterLabelValueSuffix        string                      `json:"centerLabelValueSuffix"`
	HolePercent                   uint8                       `json:"holePercent"`
	IsSliceLabelHidden            bool                        `json:"isSliceLabelHidden"`
	SliceIdentityColors           map[string]string           `json:"sliceIdentityColors"`
	SliceLabelColor               string                      `json:"sliceLabelColor"`
	SliceLabelPillBackgroundColor string                      `json:"sliceLabelPillBackgroundColor"`
}

type chartWaffleConfig struct {
	CellCount      uint                 `json:"cellCount"`
	CellGapPercent float64              `json:"cellGapPercent"`
	CellShape      ChartWaffleCellShape `json:"cellShape"`
	ColumnCount    uint                 `json:"columnCount"`
}

type chartSunburstConfig struct {
	ChildRowKey string `json:"childRowKey"`
}

type chartParliamentConfig struct {
	DotGapPercent  float64       `json:"dotGapPercent"`
	DotShape       ChartDotShape `json:"dotShape"`
	DotSizePercent float64       `json:"dotSizePercent"`
	RowCount       uint          `json:"rowCount"`
}

type chartConfig struct {
	ChartType     ChartType             `json:"chartType"`
	XAxis         chartXAxisConfig      `json:"xAxis"`
	Series        []chartSeriesConfig   `json:"series"`
	Tooltip       chartTooltipConfig    `json:"tooltip"`
	Cartesian     chartCartesianConfig  `json:"cartesian"`
	Parliament    chartParliamentConfig `json:"parliament"`
	Slices        chartSliceConfig      `json:"slices"`
	Sunburst      chartSunburstConfig   `json:"sunburst"`
	Waffle        chartWaffleConfig     `json:"waffle"`
	GapColor      string                `json:"gapColor"`
	Rows          any                   `json:"rows"`
	RowsStatePath string                `json:"rowsStatePath"`
	MaxRows       uint                  `json:"maxRows"`
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

func chartSliceHolePercentResolver(
	configuredHolePercent uint8, isPie bool,
) (uint8, error) {
	if isPie {
		return 0, nil
	}
	if configuredHolePercent == 0 {
		return chartDefaultSliceHolePercent, nil
	}
	if configuredHolePercent > chartMaxSliceHolePercent {
		return 0, fmt.Errorf(
			"ChartInvalidSliceHolePercent: %d", configuredHolePercent,
		)
	}
	return configuredHolePercent, nil
}

func chartSliceResolver(
	slices ChartSliceSettings, chartType ChartType,
) (chartSliceConfig, error) {
	isPie := chartType == ChartTypePie
	holePercent, holePercentErr := chartSliceHolePercentResolver(
		slices.HolePercent, isPie,
	)
	if holePercentErr != nil {
		return chartSliceConfig{}, holePercentErr
	}
	centerLabelDisplayMode, centerLabelDisplayModeErr := chartOptionResolver(
		slices.CenterLabelDisplayMode, ChartCenterLabelDisplayModeHidden,
		[]ChartCenterLabelDisplayMode{
			ChartCenterLabelDisplayModeAlways, ChartCenterLabelDisplayModeHidden,
			ChartCenterLabelDisplayModeHover,
		},
		"ChartInvalidCenterLabelDisplayMode",
	)
	if centerLabelDisplayModeErr != nil {
		return chartSliceConfig{}, centerLabelDisplayModeErr
	}
	if isPie {
		centerLabelDisplayMode = ChartCenterLabelDisplayModeHidden
	}
	centerLabelSize, centerLabelSizeErr := chartOptionResolver(
		slices.CenterLabelSize, ChartCenterLabelSizeMd,
		[]ChartCenterLabelSize{
			ChartCenterLabelSizeXs, ChartCenterLabelSizeSm,
			ChartCenterLabelSizeMd, ChartCenterLabelSizeLg,
			ChartCenterLabelSizeXl,
		},
		"ChartInvalidCenterLabelSize",
	)
	if centerLabelSizeErr != nil {
		return chartSliceConfig{}, centerLabelSizeErr
	}
	sliceIdentityColors := slices.SliceIdentityColors
	if sliceIdentityColors == nil {
		sliceIdentityColors = map[string]string{}
	}
	return chartSliceConfig{
		CenterLabelAlwaysSlice:        slices.CenterLabelAlwaysSlice,
		CenterLabelDisplayMode:        centerLabelDisplayMode,
		CenterLabelSizeMultiplier:     chartCenterLabelSizeMultiplierResolver(centerLabelSize),
		CenterLabelSliceColor:         chartTextDefaulter(slices.CenterLabelSliceColor, chartDefaultCenterLabelSliceColor),
		CenterLabelValueColor:         chartTextDefaulter(slices.CenterLabelValueColor, chartDefaultCenterLabelValueColor),
		CenterLabelValueSuffix:        slices.CenterLabelValueSuffix,
		HolePercent:                   holePercent,
		IsSliceLabelHidden:            slices.IsSliceLabelHidden,
		SliceIdentityColors:           sliceIdentityColors,
		SliceLabelColor:               chartTextDefaulter(slices.SliceLabelColor, "neutral-50"),
		SliceLabelPillBackgroundColor: chartTextDefaulter(slices.SliceLabelPillBackgroundColor, chartDefaultSliceLabelPillColor),
	}, nil
}

func (componentSettings ChartSettings) chartTypeResolver() (ChartType, error) {
	return chartOptionResolver(
		componentSettings.Type, ChartTypeLine,
		[]ChartType{
			ChartTypeArea, ChartTypeBar, ChartTypeDonut, ChartTypeGauge,
			ChartTypeHorizontalBar, ChartTypeLine, ChartTypeParliament,
			ChartTypePie, ChartTypeProgress, ChartTypeRadialBar,
			ChartTypeSunburst, ChartTypeWaffle,
		},
		"ChartInvalidType",
	)
}

func chartWaffleGridResolver(waffle ChartWaffleSettings) (chartWaffleConfig, error) {
	cellShape, cellShapeErr := chartOptionResolver(
		waffle.CellShape, ChartWaffleCellShapeSquare,
		[]ChartWaffleCellShape{ChartWaffleCellShapeCircle, ChartWaffleCellShapeSquare},
		"ChartInvalidWaffleCellShape",
	)
	if cellShapeErr != nil {
		return chartWaffleConfig{}, cellShapeErr
	}
	cellGapPercent := waffle.CellGapPercent
	if cellGapPercent == 0 {
		cellGapPercent = chartDefaultWaffleCellGapPercent
	}
	if math.IsNaN(cellGapPercent) ||
		cellGapPercent < 0 || cellGapPercent > chartMaxWaffleCellGapPercent {
		return chartWaffleConfig{}, fmt.Errorf(
			"ChartInvalidWaffleCellGapPercent: %v", cellGapPercent,
		)
	}
	cellCount := waffle.CellCount
	if cellCount == 0 {
		cellCount = chartDefaultWaffleCellCount
	}
	columnCount := waffle.ColumnCount
	if columnCount == 0 {
		columnCount = chartDefaultWaffleColumnCount
	}
	return chartWaffleConfig{
		CellCount:      cellCount,
		CellGapPercent: cellGapPercent,
		CellShape:      cellShape,
		ColumnCount:    columnCount,
	}, nil
}

func chartSunburstResolver(
	sunburst ChartSunburstSettings, chartType ChartType,
) (chartSunburstConfig, error) {
	if chartType == ChartTypeSunburst && sunburst.ChildRowKey == "" {
		return chartSunburstConfig{}, errors.New("ChartSunburstChildRowKeyMissing")
	}
	return chartSunburstConfig(sunburst), nil
}

func chartVegaSymbolShapeResolver(dotShape ChartDotShape) ChartDotShape {
	switch dotShape {
	case ChartDotShapeTriangleUp:
		return "triangle-up"
	case ChartDotShapeTriangleDown:
		return "triangle-down"
	}
	return dotShape
}

func chartParliamentLayoutResolver(
	parliament ChartParliamentSettings,
) (chartParliamentConfig, error) {
	dotShape, dotShapeErr := chartOptionResolver(
		parliament.DotShape, ChartDotShapeCircle,
		[]ChartDotShape{
			ChartDotShapeCircle, ChartDotShapeCross, ChartDotShapeDiamond,
			ChartDotShapeSquare, ChartDotShapeTriangleDown, ChartDotShapeTriangleUp,
		},
		"ChartInvalidParliamentDotShape",
	)
	if dotShapeErr != nil {
		return chartParliamentConfig{}, dotShapeErr
	}
	dotGapPercent := parliament.DotGapPercent
	if dotGapPercent == 0 {
		dotGapPercent = chartDefaultParliamentDotGapPercent
	}
	if math.IsNaN(dotGapPercent) ||
		dotGapPercent < 0 || dotGapPercent > chartMaxParliamentDotGapPercent {
		return chartParliamentConfig{}, fmt.Errorf(
			"ChartInvalidParliamentDotGapPercent: %v", dotGapPercent,
		)
	}
	dotSizePercent := parliament.DotSizePercent
	if dotSizePercent == 0 {
		dotSizePercent = chartDefaultParliamentDotSizePercent
	}
	if math.IsNaN(dotSizePercent) ||
		dotSizePercent < chartMinParliamentDotSizePercent ||
		dotSizePercent > chartMaxParliamentDotSizePercent {
		return chartParliamentConfig{}, fmt.Errorf(
			"ChartInvalidParliamentDotSizePercent: %v", dotSizePercent,
		)
	}
	rowCount := parliament.RowCount
	if rowCount == 0 {
		rowCount = chartDefaultParliamentRowCount
	}
	return chartParliamentConfig{
		DotGapPercent:  dotGapPercent,
		DotShape:       chartVegaSymbolShapeResolver(dotShape),
		DotSizePercent: dotSizePercent,
		RowCount:       rowCount,
	}, nil
}

func (componentSettings ChartSettings) configResolver() (chartConfig, error) {
	chartType, chartTypeErr := componentSettings.chartTypeResolver()
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
	slices, slicesErr := chartSliceResolver(componentSettings.Slices, chartType)
	if slicesErr != nil {
		return chartConfig{}, slicesErr
	}
	sunburst, sunburstErr := chartSunburstResolver(
		componentSettings.Sunburst, chartType,
	)
	if sunburstErr != nil {
		return chartConfig{}, sunburstErr
	}
	parliament, parliamentErr := chartParliamentLayoutResolver(
		componentSettings.Parliament,
	)
	if parliamentErr != nil {
		return chartConfig{}, parliamentErr
	}
	waffle, waffleErr := chartWaffleGridResolver(componentSettings.Waffle)
	if waffleErr != nil {
		return chartConfig{}, waffleErr
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
		Parliament:    parliament,
		Slices:        slices,
		Sunburst:      sunburst,
		Waffle:        waffle,
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
	if resolvedConfig.Slices.CenterLabelDisplayMode != ChartCenterLabelDisplayModeHidden {
		colorTokens = append(
			colorTokens,
			resolvedConfig.Slices.CenterLabelSliceColor,
			resolvedConfig.Slices.CenterLabelValueColor,
		)
	}
	isSliceLabelShown := !resolvedConfig.Slices.IsSliceLabelHidden &&
		(resolvedConfig.ChartType == ChartTypePie ||
			resolvedConfig.ChartType == ChartTypeRadialBar ||
			resolvedConfig.ChartType == ChartTypeSunburst)
	if isSliceLabelShown {
		colorTokens = append(
			colorTokens,
			resolvedConfig.Slices.SliceLabelColor,
			resolvedConfig.Slices.SliceLabelPillBackgroundColor,
		)
	}
	for _, seriesEntry := range resolvedConfig.Series {
		colorTokens = append(colorTokens, seriesEntry.IdentityColor)
	}
	for _, sliceColorToken := range resolvedConfig.Slices.SliceIdentityColors {
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

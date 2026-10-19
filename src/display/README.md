# Display

The display package shows content, state, and page chrome. Part of [Infinite UI](../../README.md).

## Accordion

Collapsible sections built on `<details>` and `<summary>`.

```go
@uiDisplay.Accordion(uiDisplay.AccordionSettings{
    Items: []uiDisplay.AccordionItemSettings{
        {Title: "Section 1", Content: AccordionDemoSection1(), Icon: "ph-info"},
        {Title: "Section 2", Content: AccordionDemoSection2()},
    },
})
```

- Each item takes a `Title`, a `Content` component, and an optional `Icon`.
- `BorderRadius` rounds only the outer corners: the top of the first item and the bottom of the last item.
- `ContentBackgroundColor` sets the color behind open content (default `neutral-900`). Use the color of the surface that holds the accordion. The content then shows as a cutout of that surface.
- `IsSingleOpen` keeps one item open. The items share one name, so an open item closes the others.
- `IsFirstItemOpen` opens the first item at render.
- `TextCase` takes a `uiToolset.TextCase*` value for the titles.

## Alert

A status message with a title, a description, and a variation.

```go
@uiDisplay.Alert(uiDisplay.AlertSettings{
    Title:       "Alert Title",
    Description: "This is an alert message.",
    Variation:   uiDisplay.AlertVariationInfo,
    Size:        uiDisplay.AlertSizeMd,
    IsCloseable: true,
})
```

- `Variation` takes `AlertVariationSuccess`, `AlertVariationWarning`, `AlertVariationError`, or `AlertVariationInfo`.
- `ContentHtml` replaces the title and the description with custom markup.
- `Title` and `Description` render as raw HTML. Give them only trusted content, never request data.
- `IsCloseable` shows a close button. `DisplayTimeoutSeconds` and `OnTimeoutFunc` close the alert after a delay.
- The `*OneWayStatePath` fields bind the title and the description to Alpine state.

## Chart

Line, area, bar, horizontal bar, pie, donut, radial bar, sunburst, parliament, and waffle charts drawn with [Vega](https://vega.github.io/vega/), plus a progress bar and a gauge drawn without it. The chart reads the global `vega` and does not load it. Load Vega once per page, usually with `@uiImport.HeadTagsVega()` in the layout head. Without Vega, the chart shows the legend and the table and logs `ChartVegaLibraryNotLoaded`. A page that uses only the progress or gauge type needs no Vega at all.

```go
@uiDisplay.Chart(uiDisplay.ChartSettings{
    Type: uiDisplay.ChartTypeArea,
    XAxis: uiDisplay.ChartXAxisSettings{
        RowKey:    "time",
        ScaleType: uiDisplay.ChartScaleTypeTime,
    },
    Series: []uiDisplay.ChartSeriesSettings{
        {RowKey: "cpu", Label: "CPU"},
        {RowKey: "memory", Label: "Memory", IdentityColor: "sky-400"},
    },

    // OptionalFields
    Cartesian: uiDisplay.ChartCartesianSettings{
        YAxis: uiDisplay.ChartYAxisSettings{
            TickLabelFormat: ".0%",
            MinValue:        "0",
            MaxValue:        "1",
        },
    },
    RowsOneWayStatePath: "resourceUsageRows",
    MaxRows:             20,
})
```

An unknown option value, a missing `RowKey`, or a bound that is not a number stops the render with an error.

### Rows

- A row is one object. `XAxis.RowKey` names the key with the X value. Each series `RowKey` names the key with that series value. Example: `{time: 1760360000000, cpu: 0.42, memory: 0.61}`.
- `Rows` sets the initial rows from Go. Any value that encodes to a JSON array of objects works.
- `RowsOneWayStatePath` binds the rows to Alpine state. The chart draws again when the array or a row changes.
- Give `RowsOneWayStatePath` from code, never from request data. The component puts it into a client-side expression.
- `MaxRows` keeps only the newest rows. `0` keeps every row.

### Type and marks

- `Type` takes `ChartTypeLine` (the default), `ChartTypeArea`, `ChartTypeBar`, `ChartTypeHorizontalBar`, `ChartTypePie`, `ChartTypeDonut`, `ChartTypeRadialBar`, `ChartTypeSunburst`, `ChartTypeParliament`, `ChartTypeWaffle`, `ChartTypeProgress`, or `ChartTypeGauge`.
- A line, area, or bar chart reads `Cartesian`. A progress chart reads `Progress`, and a gauge chart reads `Gauge`; both skip `XAxis` and `Series`.
- The donut-like variations read `Slices`. The donut, the pie, and the radial bar have no group of their own: `Type` selects the variation and `Slices` holds every setting. The sunburst, the parliament, and the waffle read `Slices` and their own group for the extra fields.
- `Cartesian.StackMode` takes `ChartStackModeNone` (the default), `ChartStackModeStacked`, or `ChartStackModePercent`. A stacked chart piles the series. A percent chart shows each series as a share of the total, and the Y axis format changes to `".0%"`. Without a stack, a bar chart puts the series side by side.
- `Cartesian.AreaFillStyle` takes `ChartAreaFillStyleGradient` (the default), `ChartAreaFillStyleSolid`, or `ChartAreaFillStyleNone`.
- `Cartesian.MarkerDisplayMode` takes `ChartMarkerDisplayModeHover` (the default), `ChartMarkerDisplayModeAlways`, or `ChartMarkerDisplayModeHidden`.
- A donut reads the first series. `XAxis.RowKey` names each slice. `Slices.HolePercent` sets the hole size (default 60, maximum 90).
- A pie is a donut with no hole. It reads the same rows as the donut, ignores `Slices.HolePercent`, and hides the center label. Each slice shows its share of the total on a rounded pill at the slice midpoint.
- A radial bar is a circular bar chart. It reads the same rows as the donut. Each category gets one bar that grows from `Slices.HolePercent` to the outer edge, and the legend lists the categories in row order. Each bar shows its share of the largest value at its midpoint on a rounded pill.
- A sunburst is a donut with a second ring. `XAxis.RowKey` names the category and `Sunburst.ChildRowKey` the child. The inner ring holds one arc per category and the outer ring one arc per child. A category arc spans the sum of its children. Each arc shows its percentage at its midpoint on a rounded pill: a category arc its share of the total and a child arc its share of its category.
- The pie, the radial bar, and the sunburst labels take `Slices.IsSliceLabelHidden` to remove them, `Slices.SliceLabelColor` to paint the text, and `Slices.SliceLabelPillBackgroundColor` to paint the pill behind the text (default `neutral-950/70`).
- A parliament draws one dot per seat over concentric semicircle rows. `XAxis.RowKey` names the party and the first series holds the seat count. Parties fill the seats as contiguous blocks, ordered by row and then left to right. `Parliament.RowCount` sets the number of rows (default 6). `Parliament.DotShape` takes `ChartDotShapeCircle` (the default), `ChartDotShapeSquare`, `ChartDotShapeDiamond`, `ChartDotShapeCross`, `ChartDotShapeTriangleUp`, or `ChartDotShapeTriangleDown`. `Parliament.DotGapPercent` sets the empty share of the seat-to-seat space up to 95 (default 25; zero selects the default), and `Parliament.DotSizePercent` scales each dot from 10 to 200 on top of that space (default 100).
- A waffle shows the parts of a whole as a grid of equal cells. Each category gets as many cells as its share of the total. `Waffle.CellCount` sets the cell total (default 100) and `Waffle.ColumnCount` the grid width (default 10). A zero count selects the default. `Waffle.CellShape` takes `ChartWaffleCellShapeSquare` (the default) or `ChartWaffleCellShapeCircle`, and `Waffle.CellGapPercent` sets the space between cells from 0 to 60 (default 4). The cells tile the grid exactly, so no cell is lost to rounding.
- `Slices.CenterLabelDisplayMode` shows the slice name and its value in the hole of a donut or a pie. It takes `ChartCenterLabelDisplayModeHidden` (the default), `ChartCenterLabelDisplayModeHover`, or `ChartCenterLabelDisplayModeAlways`. Always keeps the at-rest slice. A hover replaces it. Hover mode shows the label only while a slice is hovered. A radial bar ignores the center label.
- `Slices.CenterLabelAlwaysSlice` names the at-rest slice. An empty value shows the first slice. An unknown name logs `ChartCenterLabelSliceNotFound` and shows the first slice. `Slices.CenterLabelSliceColor` and `Slices.CenterLabelValueColor` set the two label colors. The value takes `Tooltip.ValueFormat`.
- `Slices.CenterLabelSize` takes `ChartCenterLabelSizeXs` through `ChartCenterLabelSizeXl` (default `Md`). It scales both label lines together. `Slices.CenterLabelValueSuffix` follows the center label value only, for example `" GiB"`. The tooltip keeps its own value.
- `GapColor` sets the color of the gaps between bars, slices, and markers (default `neutral-900`). Match it to the background behind the chart.

### Colors

- Every color field takes a Tailwind color token, for example `"sky-400"`, `"red-500/20"`, or `"[#4c8dff]"`. The chart reads the real color from a hidden element with that class.
- `IdentityColor` is the one color of a series. The line, the area, the bar, the marker, the legend swatch, and the tooltip swatch all use it, so the legend always matches the plot.
- Without `IdentityColor`, each series gets the next color of a fixed palette of eight. The palette passes color-vision checks on the dark surface. A ninth series shows gray and logs `ChartSeriesPaletteExceeded`, so group extra series.
- `Slices.SliceIdentityColors` maps a slice name to a token. Other slices use the palette.
- A token that the page cannot resolve logs `ChartColorTokenUnresolved`.

### Series

- `RowKey` is required. An empty `Label` shows the `RowKey`.
- `LineDashStyle` takes `ChartLineDashStyleSolid` (the default), `ChartLineDashStyleDashed`, or `ChartLineDashStyleDotted`.

### Axes

- `XAxis.ScaleType` takes `ChartScaleTypeCategory` (the default), `ChartScaleTypeTime`, or `ChartScaleTypeNumber`. A category axis keeps the row order. A time axis takes Unix milliseconds or ISO 8601 strings and sorts by time. A number axis sorts by value.
- `TickLabelFormat` takes a [d3 time format](https://d3js.org/d3-time-format) on a time axis and a [d3 number format](https://d3js.org/d3-format) on a number axis. A category axis ignores it and shows the row values as they are. The Y default is `",.2~f"`.
- `Title` shows a title on the axis. `XAxis.Title` also names the first column of the screen-reader table.
- `TitleColor`, `TickLabelColor`, `BaselineColor`, and `GridLineColor` set the axis colors.
- `IsHidden` hides the labels, the ticks, and the baseline. The grid stays.
- `XAxis.IsTickStartEndOnly` shows only the first and the last X label.
- `XAxis.IsGridLineVisible` shows the X grid, which is off by default. `YAxis.IsGridLineHidden` hides the Y grid.
- `YAxis.MinValue` and `YAxis.MaxValue` fix the value range. Empty values fit the data.
- `YAxis.TickCount` sets the number of Y labels (default 5). `YAxis.Position` takes `ChartYAxisPositionLeft` (the default) or `ChartYAxisPositionRight`.
- A horizontal bar draws the category axis down the left edge and the value axis along the bottom. `XAxis` still names the category and `Cartesian.YAxis` still names the value axis. `YAxis.Position` right moves the value labels to the top.

### Legend

- `DisplayMode` takes `ChartLegendDisplayModeAuto` (the default, two or more entries), `ChartLegendDisplayModeAlways`, or `ChartLegendDisplayModeHidden`.
- `Position` takes `ChartLegendPositionTop` (the default) or `ChartLegendPositionBottom`.
- `Alignment` takes `ChartLegendAlignmentStart` (the default), `ChartLegendAlignmentCenter`, or `ChartLegendAlignmentEnd`.
- `SwatchShape` takes a `uiToolset.Shape*` constant (default circular). `TextColor` and `GapSize` style the entries.

### Tooltip

- On a line or area chart, hover shows a crosshair and each series value at that point. A bar or a slice shows its own value.
- `TitleFormat` formats the X value. A time axis defaults to `"%b %d, %H:%M:%S"`. A category axis ignores it and shows the row value as it is. `ValueFormat` formats the series values. It defaults to the Y tick label format, and to `",.2~f"` on a percent stack.
- `ContentFunc` names a global JavaScript function that returns the tooltip HTML for the hovered row. It receives the row object and replaces the default title and rows. The returned HTML is trusted and inserted as is: never return request data without sanitizing it first. A name that is not a function logs `ChartTooltipContentFuncMissing` and falls back to the default.
- `BackgroundColor`, `TitleColor`, `SeriesLabelColor`, and `ValueColor` set the colors. `RingColor`, `RingThickness`, `ShadowSize`, `PaddingSize`, `BorderRadius`, `MinWidthClass`, and `MaxWidthClass` style the surface.
- `CrosshairColor` sets the crosshair color. `IsCrosshairHidden` hides it. `IsHidden` hides the tooltip.

### Layout and accessibility

- `HeightClass` sets the plot height (default `h-64`).
- `EmptyState` takes a component that shows over the plot when there are no rows. The default is the text `No data`.
- `AriaLabel` names the plot and the screen-reader table of rows. Without it, the series labels name the plot. On a progress chart it names the bar, and on a gauge chart it names the meter.

### Progress

`ChartTypeProgress` draws an HTML progress bar and no Vega canvas. It reads `Progress` and skips `XAxis`, `Series`, and the Vega sections. A page that uses only the progress type needs no Vega head tag.

- `Progress.Value` is the number the bar shows. `Progress.MinValue` and `Progress.MaxValue` set the value range; a value outside the range clamps to the nearest end. The default range is 0 to 100.
- `Progress.Thresholds` are checked in list order, and the first entry whose `MinimumPercent` has been reached paints the fill. List them from highest to lowest percent. Without thresholds, the fill takes red-800 at 90%, red-600 at 85%, amber-600 at 75%, amber-300 at 50%, and cyan-400 below.
- `Progress.DefaultColor` paints every percent below the lowest threshold.
- `Progress.IsIndeterminate` hides the percent and animates a sliding segment, for work of unknown length.
- `Progress.IsStriped` draws diagonal stripes over the fill in `Progress.StripeColor` (default `neutral-50/20`).
- `Progress.ShouldShowPercentLabel` shows a percent label. `Progress.Label` replaces the text with fixed content, and `Progress.LabelOneWayStatePath` reads the text from a state path. `Progress.LabelPrefix` and `Progress.LabelSuffix` wrap the percent text.
- `Progress.ValueOneWayStatePath` makes the bar follow a state path. The bar reads the path and never writes back.
- `Progress.Size` takes `ChartProgressSizeXs`, `ChartProgressSizeSm` (default), or `ChartProgressSizeMd` and sets the track height and the label text size.
- `Progress.MinWidthClass`, `Progress.MaxWidthClass`, `Progress.MinHeightClass`, and `Progress.MaxHeightClass` constrain the bar; it fills its container by default with a `min-w-20` floor.
- `Progress.TrackBackgroundColor`, `Progress.DefaultColor`, `Progress.LabelBackgroundColor`, and `Progress.LabelTextColor` take color tokens. The label sits on a rounded pill; the pill takes `LabelBackgroundColor`, which defaults to `neutral-950/70`.
- `Progress.LabelTextSizeClass` takes a Tailwind text size class for the label alone, for example `text-lg`. Without it, the label follows `Progress.Size`.
- Give the state paths from code, never from request data. The component puts them into client-side expressions.

### Gauge

`ChartTypeGauge` draws a needle over a track, with optional colored bands and a center label. It renders as plain SVG with the `meter` role and no Vega canvas. It reads `Gauge` and skips `XAxis`, `Series`, and the Vega sections. A page that uses only the gauge type needs no Vega head tag.

- `Gauge.Value` is the number the needle shows. `Gauge.MinValue` and `Gauge.MaxValue` set the value range; a value outside the range clamps to the nearest end. The default range is 0 to 100. The center text and the accessible value show the clamped value; the needle follows the same clamp.
- `Gauge.Bands` draws colored arcs between percent marks, in list order. Each band spans `MinimumPercent` to `MaximumPercent`, where 0 is the start of the arc and 100 its end. A band outside that range, or with the minimum at or above the maximum, fails the render with `ChartInvalidGaugeBandPercent`. A band color defaults to `neutral-700`.
- `Gauge.IsSemicircle` draws the top half only.
- `Gauge.Label` shows fixed text under the value. `Gauge.ValueColor` and `Gauge.LabelColor` paint the two center texts (defaults `neutral-50` and `neutral-400`). `Gauge.ValueTextSize` and `Gauge.LabelTextSize` set the two font sizes (defaults 24 and 12, allowed 6 to 72). `Gauge.TextOffset` moves both texts down, away from the pivot (default 0, allowed 0 to 80); the viewBox grows with the offset and the larger of the two text sizes, and widens for a long value or label, so the text never clips.
- `Gauge.ArcThickness` sets the track and band width in the 200-unit viewBox (default 16, allowed 4 to 36).
- `Gauge.NeedleLength`, `Gauge.NeedleThickness`, and `Gauge.PivotRadius` shape the needle: the length accepts 10 to 90 (default 55), the thickness 1 to 12 (default 3), and the pivot radius up to 20 (default 5). A zero value selects its default, and a value outside its range fails the render.
- `Gauge.TickCount` draws that many evenly spaced tick intervals across `MinValue` to `MaxValue`, each with a label; zero draws none. `Gauge.TickPosition` places them `ChartGaugeTickPositionOutside` (the default) or `ChartGaugeTickPositionInside` the track. Outside ticks widen the viewBox so the labels stay in view. On a full-circle gauge, inside ticks label the top half only, so the center value and label keep their space. `Gauge.TickColor` and `Gauge.TickLabelColor` paint them.
- `Gauge.TrackColor` and `Gauge.NeedleColor` take color tokens. The defaults are `neutral-50/15` and `neutral-50`. `NeedleColor` paints the needle and its pivot.
- The center value and label sit below the pivot with clearance, so the needle does not touch them.
- Give the color tokens and the label from code, never from request data. The component puts them into class names and text.

### Waffle

`ChartTypeWaffle` draws the parts of a whole as a grid of equal cells on a Vega canvas. It reads the first series and `XAxis.RowKey`, and it colors the cells like the donut, through `Slices.SliceIdentityColors` and the palette. The legend lists the categories in row order. It ignores `Slices.HolePercent` and the center label.

- `Waffle.CellCount` sets the number of cells (default 100). `Waffle.ColumnCount` sets the number of cells per row (default 10). A zero count selects the default.
- `Waffle.CellShape` takes `ChartWaffleCellShapeSquare` (the default) or `ChartWaffleCellShapeCircle`. `Waffle.CellGapPercent` sets the space between cells up to 60 (default 4; zero selects the default). A value outside the range fails the render.
- Each category ends on a rounded cumulative boundary, so the categories tile every cell exactly.
- The grid centers itself in the plot, whatever the height and the width of the canvas.

### Sunburst

`ChartTypeSunburst` draws a two-level hierarchy as two concentric rings on a Vega canvas. `XAxis.RowKey` names the category, `Sunburst.ChildRowKey` names the child of each row, and the first series holds the value. The inner ring holds one arc per category, and the outer ring one arc per child. A category arc spans the sum of its children, so the rings stay aligned. The legend lists the categories in row order.

- `Sunburst.ChildRowKey` is required. A missing name fails the render with `ChartSunburstChildRowKeyMissing`.
- A row with an empty child draws its category arc only.
- A child draws in its category color at 60 percent opacity. The colors come from `Slices.SliceIdentityColors` and the palette, like the donut.
- The category tooltip shows the category total; a child tooltip shows the child value.

### Parliament

`ChartTypeParliament` draws one dot per seat over concentric semicircle rows on a Vega canvas. It reads the first series and `XAxis.RowKey`, and it colors the dots like the donut, through `Slices.SliceIdentityColors` and the palette. The legend lists the parties in row order. It ignores `Slices.HolePercent` and the center label.

- `Parliament.RowCount` sets the number of rows (default 6). A zero count selects the default.
- `Parliament.DotShape` takes `ChartDotShapeCircle` (the default), `ChartDotShapeSquare`, `ChartDotShapeDiamond`, `ChartDotShapeCross`, `ChartDotShapeTriangleUp`, or `ChartDotShapeTriangleDown`. `Parliament.DotSizePercent` scales each dot, accepting 10 to 200 (default 100; zero selects the default). A bad shape or size fails the render.
- The seats fill row by row, from the innermost row out and from left to right. Each party occupies a contiguous block of seats.
- A seat count with a fraction rounds to whole seats at each party boundary, so the parties tile every seat exactly.

## CloakLoading

A full-screen loader that covers the page until CSS and JavaScript load. Put it first in `<body>`.

```go
@uiDisplay.CloakLoading(uiDisplay.CloakLoadingSettings{
    HideDelaySeconds: "1",
    TextMessage:      "Loading...",
    Icon:             "ph-compass-rose",
})
```

- It uses inline styles, so it renders before the stylesheet loads.
- It hides `HideDelaySeconds` after Alpine starts.
- The icon size takes a `CloakLoadingIconSize*` constant. The animation takes a `CloakLoadingAnimationName*` constant.

## CollapsibleText

A block of text that truncates to a line limit with an inline Show more affordance. The affordance renders only when the text exceeds the limit, and it sits directly after the ellipsis on the truncated line.

```go
@uiDisplay.CollapsibleText(uiDisplay.CollapsibleTextSettings{
    Text: "A paragraph long enough to exceed three lines on a narrow card.",

    // OptionalFields
    MaxVisibleLines: 3,
    ExpandLabel:     "show details",
    CollapseLabel:   "show summary",
    ExpandIcon:      "ph-caret-down",
    CollapseIcon:    "ph-caret-up",
    FlankIconLeft:   "ph-file-text",
    FlankIconRight:  "ph-caret-up-down",
    IconColor:       "neutral-400",
    Size:            uiDisplay.CollapsibleTextSizeSm,
    TextColor:       "neutral-300",
    ToggleColor:     "secondary-400",
})
```

- `MaxVisibleLines` caps the visible lines; values below 1 default to 3. The component measures the rendered text, truncates it at the last word that fits, and renders the affordance inline after the ellipsis. The affordance renders only when the text overflows.
- `ExpandLabel` and `CollapseLabel` replace the affordance text. The defaults are lowercase: "show more" and "show less".
- `ExpandIcon` and `CollapseIcon` set the affordance icon for each state. The defaults are `ph-caret-down` and `ph-caret-up`.
- `Size` takes `CollapsibleTextSizeXs`, `CollapsibleTextSizeSm` (default), `CollapsibleTextSizeMd`, or `CollapsibleTextSizeLg`. It scales the text, the affordance, and the flanking icons together.
- `FlankIconLeft` and `FlankIconRight` place a Phosphor icon at each side of the text, rendered in the bold weight. `IconColor` takes a color token with a neutral-400 default.
- `TextColor` and `ToggleColor` take color tokens. The defaults are neutral-300 and secondary-400.

## ConfirmationDialog and Presets

Confirmation dialogs built on Modal. Four presets set the icon, tone, and text: `ConfirmDialog`, `WarningDialog`, `CriticalDialog`, and `DeleteDialog`. `CriticalDialog` and `DeleteDialog` ask the user to type the target name and have no fullscreen mode.

```go
@uiDisplay.DeleteDialog(uiDisplay.ConfirmationDialogSettings{
    IsVisibleTwoWayStatePath: "isDeleteDialogVisible",
    OnConfirmFunc:            "deleteRecord()",

    // OptionalFields
    TargetNameStatePath: "record.name",
    TargetIdStatePath:   "record.id",
})
```

- `Size` takes `ConfirmationDialogSizeXs` (40%) through `ConfirmationDialogSizeXxl` (90%), as viewport percentages on both axes. The default is `ConfirmationDialogSizeSm` (50%). The dialog has no fullscreen size and no resize control.
- `WidthPercent` and `HeightPercent` replace the box, as in Modal. Without `HeightPercent`, the height fits the content up to 85%.
- The header takes `HeaderTitle`, `HeaderSubHeading`, their `*OneWayStatePath` and `*Color` fields, `HeaderSize`, `TextCase`, and the `HeaderIcon*` fields. Each preset sets the title and icon defaults.
- `HeaderIconPosition` defaults to `HeaderIconPositionTop`, with the icon above the title. `HeaderIconPositionLeft` puts it beside the title.
- `ActionsContent` adds buttons next to the close button. `HeaderContent` replaces the header. `TextColor` sets the text color.
- `TargetNameStatePath` and `TargetIdStatePath` put live values in the message. The name shows as a dashed chip and the id as `#<id>`.
- `IsTypeToConfirmEnabled` turns on the typed check. The confirm button enables only when the typed text matches the target name, or the id when no name path is set. An empty target keeps the button disabled.
- `TypeToConfirmValueTwoWayStatePath` keeps the typed text in your Alpine state. By default, the dialog keeps it.
- Each typed-check dialog keeps its own text, but all use the input name `typedConfirmationValue`. Do not put two of them in one submitted form.
- Give the state paths from code, never from request data. The component puts them into client-side expressions.
- `Message` replaces the note under the question. `MessageContent` replaces the full message block.
- `CancelFunc` closes the dialog by default. `OnConfirmFunc` runs as written, so add the close when the action succeeds.
- The cancel button has an outline. The confirm button has the preset tone fill.
- `HeaderTitle`, `ConfirmButtonLabel`, `ConfirmButtonIcon`, and `CancelButtonLabel` replace the preset text.

## HeaderBlock

The header row of PageHeading, Card, and ConfirmationDialog: icon, title, sub-heading, and an actions slot on the right.

```go
@uiDisplay.HeaderBlock(uiDisplay.HeaderBlockSettings{
    HeaderTitle: "Network",

    // OptionalFields
    HeaderSubHeading:   "Interfaces and routes",
    HeaderIcon:         "ph-globe",
    HeaderIconPosition: uiDisplay.HeaderIconPositionLeft,
    HeadingLevel:       uiDisplay.HeaderBlockHeadingLevelTwo,
})
```

- `HeadingLevel` takes `HeaderBlockHeadingLevelOne` or `HeaderBlockHeadingLevelTwo` (the default).
- `HeaderSize` takes `HeaderSizeXs` through `HeaderSizeXl`. It scales the title, the icon, and the sub-heading.
- The `HeaderIcon*` fields match the `HeaderIcon` settings. `HeaderIconPosition` takes `HeaderIconPositionLeft` or `HeaderIconPositionTop`.
- `HeaderTitleColor`, `HeaderSubHeadingColor`, and `TextColor` take color tokens.
- `TextCase` takes a `uiToolset.TextCase*` value for the title and the sub-heading. `Card`, `PageHeading`, and `ConfirmationDialog` pass it on.

## HeaderIcon

An icon with an optional background chip.

```go
@uiDisplay.HeaderIcon(uiDisplay.HeaderIconSettings{
    Name: "ph-cube",

    // OptionalFields
    Color:           "red-950",
    BackgroundColor: "neutral-300",
    BorderRadius:    uiToolset.BorderRadiusXl,
    PaddingSize:     uiToolset.PaddingSizeLg,
})
```

- `Color` sets the icon color. `BackgroundColor` sets the chip color.
- `HeaderBlock` places the icon through `HeaderIconPosition`.
- `BorderRadius` takes `uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull`. `PaddingSize` takes `uiToolset.PaddingSizeNone` through `uiToolset.PaddingSizeXl`.

## LoadingOverlay

An overlay for requests in progress. It binds to a state path or follows the HTMX `htmx-request` class.

```go
@uiDisplay.LoadingOverlay(uiDisplay.LoadingOverlaySettings{
    IsLoadingOneWayStatePath: "isLoading",
    Icon:                     "ph-compass-rose",
    IconSize:                 uiDisplay.LoadingOverlayIconSizeMd,
    AnimationName:            uiDisplay.LoadingOverlayAnimationNameSpin,
})
```

Set `hx-indicator="#loading-overlay"` on HTMX elements, or call `UiToolset.ToggleLoadingOverlay()`. `Id` gives a second overlay its own id.

## Modal

A window dialog with header, middle, and footer slots.

```go
@uiDisplay.Modal(uiDisplay.ModalSettings{
    MiddleContent:            ModalDemoContent(),
    Title:                    "Modal Title",
    InitialSize:              uiDisplay.ModalSizeMd,
    IsVisibleTwoWayStatePath: "isModalVisible",
})
```

- The standard header shows when `Title` is set and `HeaderContent` is nil. It has the enlarge, reduce, and close buttons. A custom `HeaderContent` keeps the same buttons beside it.
- `InitialSize` sets the panel as a viewport percentage on both axes: `ModalSizeXs` (40%), `ModalSizeSm` (50%), `ModalSizeMd` (60%), `ModalSizeLg` (70%), `ModalSizeXl` (80%), `ModalSizeXxl` (90%), and `ModalSizeFull` (100%). The panel follows the screen shape. It also scales the header title and padding.
- `WidthPercent` and `HeightPercent` replace one axis with a whole percentage from 1 to 100. With both set, the resize buttons do not show.
- `MinWidthPercent`, `MaxWidthPercent`, `MinHeightPercent`, and `MaxHeightPercent` set limits on each axis. The panel keeps its size inside these limits.
- `PossibleSizes` lists the sizes in order (default md, lg, xl, xxl). The list starts at `InitialSize`. An `InitialSize` outside the list falls to the first entry.
- `CurrentSizeTwoWayStatePath` binds the live size. When the field is empty and the panel can resize, the modal makes its own path.
- Enlarge moves one size up, and reduce moves one size down. Reduce hides at the initial size, and enlarge hides at the largest size.
- `ShouldHeightMatchContent` sets the height to the content, up to 85%.
- `IsMiddleContentScrollDisabled` stops the middle scroll, so the content controls its overflow. Use it for a full-bleed panel such as a terminal.
- `IsUnresizable` and `IsUncloseable` remove those buttons. The resize buttons also hide when fewer than two sizes are available.
- The close button and a backdrop click set the visibility path to false and run `OnCloseFunc`. `IsUncloseable` stops both, so the app must close the modal.
- `OnResizeFunc` runs on each resize click.

## Tag

A label with an outer ring and an inner background.

```go
@uiDisplay.Tag(uiDisplay.TagSettings{
    OuterLeftIcon:  "ph-info",
    OuterLeftLabel: "Info",
    InnerIcon:      "ph-warning",
    InnerLabel:     "Warning",
    Size:           uiDisplay.TagSizeXs,
})
```

- The `OuterLeft*`, `OuterRight*`, and `Inner*` fields take an icon, a static label, or a `*OneWayStatePath` for live text.
- `InnerValueOneWayStatePath` binds a live value without the lowercase change. Use it with a static `OuterLeftLabel` for a label-value tag.
- `Size` takes `TagSizeTiny` through `TagSizeXl`. `TagSizeTiny` fits table cells.
- `TextCase` takes a `uiToolset.TextCase*` value for the labels.
- `OnRemoveFunc` adds a remove button. `RemoveButtonLabel` names it.
- `OuterRadius` and `InnerRadius` take `uiToolset.BorderRadius*` values (default `uiToolset.BorderRadiusMd`).
- `OuterBackgroundColor`, `OuterRingColor`, and `OuterTextColor` take color tokens, for example `"neutral-50/10"`.

## Toast

A notification popup that an Alpine store controls.

```go
@uiDisplay.Toast(uiDisplay.ToastSettings{
    TextColor:          "neutral-50",
    Size:               uiDisplay.ToastSizeMd,
    AutoDismissSeconds: 10,
})
```

- Render one Toast per page.
- Call `$store.toast.displayToast(message, type)` or `$store.toast.displayToastWithApiResponse(apiResponse, httpStatusCode)`. Call `$store.toast.clearToast()` to hide it.
- The background defaults to `neutral-800/90`. The ring defaults to `ring-2` in the tone color: green for `success`, yellow for `partialSuccess`, and red for `danger`. `BackgroundColor`, `RingColor`, and `RingThickness` change them.
- With HTMX on the page, the Toast listens to `htmx:afterRequest`. A JSON response with a `readableMessage`, a `humanReadableMessage` outcome, or a string `body` shows a toast. The HTTP status sets the tone: 2xx success, 207 partial success, 4xx and higher danger.
- The toast closes after 10 seconds by default.

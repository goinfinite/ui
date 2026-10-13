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

Line, area, bar, and donut charts drawn with [Vega](https://vega.github.io/vega/). The chart reads the global `vega` and does not load it. Load Vega once per page, usually with `@uiImport.HeadTagsVega()` in the layout head. Without Vega, the chart shows the legend and the table and logs `ChartVegaLibraryNotLoaded`.

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
- `MaxRows` keeps only the newest rows. `0` keeps every row.

### Type and marks

- `Type` takes `ChartTypeLine` (the default), `ChartTypeArea`, `ChartTypeBar`, or `ChartTypeDonut`.
- A line, area, or bar chart reads `Cartesian`. A donut reads `Donut`.
- `Cartesian.StackMode` takes `ChartStackModeNone` (the default), `ChartStackModeStacked`, or `ChartStackModePercent`. A stacked chart piles the series. A percent chart shows each series as a share of the total, and the Y axis format changes to `".0%"`. Without a stack, a bar chart puts the series side by side.
- `Cartesian.AreaFillStyle` takes `ChartAreaFillStyleGradient` (the default), `ChartAreaFillStyleSolid`, or `ChartAreaFillStyleNone`.
- `Cartesian.MarkerDisplayMode` takes `ChartMarkerDisplayModeHover` (the default), `ChartMarkerDisplayModeAlways`, or `ChartMarkerDisplayModeHidden`.
- A donut reads the first series. `XAxis.RowKey` names each slice. `Donut.HolePercent` sets the hole size (default 60, maximum 90).
- `GapColor` sets the color of the gaps between bars, slices, and markers (default `neutral-900`). Match it to the background behind the chart.

### Colors

- Every color field takes a Tailwind color token, for example `"sky-400"`, `"red-500/20"`, or `"[#4c8dff]"`. The chart reads the real color from a hidden element with that class.
- `IdentityColor` is the one color of a series. The line, the area, the bar, the marker, the legend swatch, and the tooltip swatch all use it, so the legend always matches the plot.
- Without `IdentityColor`, each series gets the next color of a fixed palette of eight. The palette passes color-vision checks on the dark surface. A ninth series shows gray and logs `ChartSeriesPaletteExceeded`, so group extra series.
- `Donut.SliceIdentityColors` maps a slice name to a token. Other slices use the palette.
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

### Legend

- `DisplayMode` takes `ChartLegendDisplayModeAuto` (the default, two or more entries), `ChartLegendDisplayModeAlways`, or `ChartLegendDisplayModeHidden`.
- `Position` takes `ChartLegendPositionTop` (the default) or `ChartLegendPositionBottom`.
- `Alignment` takes `ChartLegendAlignmentStart` (the default), `ChartLegendAlignmentCenter`, or `ChartLegendAlignmentEnd`.
- `SwatchShape` takes a `uiToolset.Shape*` constant (default circular). `TextColor` and `GapSize` style the entries.

### Tooltip

- On a line or area chart, hover shows a crosshair and each series value at that point. A bar or a slice shows its own value.
- `TitleFormat` formats the X value. A time axis defaults to `"%b %d, %H:%M:%S"`. A category axis ignores it and shows the row value as it is. `ValueFormat` formats the series values. It defaults to the Y tick label format, and to `",.2~f"` on a percent stack.
- `ContentFunc` names a global JavaScript function that returns the tooltip HTML for the hovered row. It receives the row object and replaces the default title and rows. A name that is not a function logs `ChartTooltipContentFuncMissing` and falls back to the default.
- `BackgroundColor`, `TitleColor`, `SeriesLabelColor`, and `ValueColor` set the colors. `RingColor`, `RingThickness`, `ShadowSize`, `PaddingSize`, `BorderRadius`, `MinWidthClass`, and `MaxWidthClass` style the surface.
- `CrosshairColor` sets the crosshair color. `IsCrosshairHidden` hides it. `IsHidden` hides the tooltip.

### Layout and accessibility

- `HeightClass` sets the plot height (default `h-64`).
- `EmptyState` takes a component that shows over the plot when there are no rows. The default is the text `No data`.
- `AriaLabel` names the plot and the screen-reader table of rows. Without it, the series labels name the plot.

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

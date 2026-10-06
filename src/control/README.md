# Control

Control layer of Infinite UI. It triggers actions and picks values. Components here do not collect form data. Part of [Infinite UI](../../README.md).

## Button

Button with optional label, icons, tooltip, ring, and shape variants.

```go
@uiControl.Button(uiControl.ButtonSettings{
    Label:       "Click me",
    IconLeft:    "ph-info",
    OnClickFunc: "alert('Button clicked!')",
})
```

- `Size` accepts the `uiControl.ButtonSize*` constants.
- `Shape` accepts `ButtonShapeCircular`, `ButtonShapeRounded`, or `ButtonShapeSquare`.
- Icons use the `ph-bold` weight. Pass the icon name alone, for example `"ph-floppy-disk"`.
- `IsDisabled` sets a static disabled state. `IsDisabledOneWayStatePath` disables the button while the path is truthy and dims it.
- `IsVisibleOneWayStatePath` hides the button while the path is falsy.
- `IsSubmit` renders `type="submit"` so the button submits its form. Every button renders a `type`; the default is `button`, so a button inside a form does not submit it.
- `TooltipContent` renders a text tooltip. `TooltipContentHtml` accepts a component, and `TooltipContentOneWayStatePath` binds live text. `TooltipPosition` accepts the `uiToolset.TooltipPosition*` constants. The tooltip renders in a fixed layer on the document body. A scroll container cannot clip it or gain a scrollbar from it.
- `BackgroundColor`, `TextColor`, and `RingColor` accept Tailwind color tokens, for example `"red-500/20"`. The `*Hover` fields set the hover state.
- `TextCase` accepts a `uiToolset.TextCase*` value and transforms the label. The default, `TextCaseNone`, leaves it as typed.

## RangeSlider

Slider with one or two thumbs, value bubbles, tick marks, and gradient tracks.

```go
@uiControl.RangeSlider(uiControl.RangeSliderSettings{
    ThumbValueTwoWayStatePath: "sliderValue",
    TrackStartValue:           "0",
    TrackEndValue:             "100",
    TrackSteps:                "1",
})
```

- The track is configured with `TrackStartValue`, `TrackEndValue`, `TrackSteps`, and the optional `TrackStartValueLabel*` and `TrackEndValueLabel*` fields.
- `TrackTooltipEnabled` previews the value under the pointer while it moves along the track. The preview snaps to `TrackSteps` and clamps to the track bounds. `TrackTooltipPosition` accepts `RangeSliderTrackTooltipPositionTop` or `RangeSliderTrackTooltipPositionBottom`; the `TrackTooltipBackgroundColor`, `TrackTooltipTextColor`, `TrackTooltipRingColor`, and `TrackTooltipRingThickness` fields style it.
- Single-thumb mode binds `ThumbValueTwoWayStatePath`. `ThumbValueBubbleEnabled` shows the current value above the thumb. `ThumbValueBubbleDisplayMode` accepts `RangeSliderThumbValueBubbleDisplayModeAlways` (the default) or `RangeSliderThumbValueBubbleDisplayModeHover`, which shows the bubble only while the slider is hovered or dragged.
- `ThumbDualValueModeEnabled` adds a second thumb. Bind it with `ThumbUpperValueTwoWayStatePath`. The upper thumb mirrors the lower thumb settings through the `ThumbUpper*` fields.
- `ThumbAriaLabel` and `ThumbUpperAriaLabel` name the thumbs for screen readers. A thumb falls back to its visible label, then to a default name.
- Thumb corrections run at init and on change: an out-of-bounds value is clamped, and in dual mode the thumbs keep one step of separation.

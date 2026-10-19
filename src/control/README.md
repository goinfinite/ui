# Control

The control package triggers actions and picks values. Its components do not collect form data. Part of [Infinite UI](../../README.md).

## Button

A button with an optional label, icons, tooltip, ring, and shape.

```go
@uiControl.Button(uiControl.ButtonSettings{
    Label:       "Click me",
    IconLeft:    "ph-info",
    OnClickFunc: "alert('Button clicked!')",
})
```

- `Size` takes a `uiControl.ButtonSize*` constant.
- `PaddingSize` takes a `uiToolset.PaddingSize*` constant and replaces the padding from `Size`. For a fixed square button, use `uiToolset.PaddingSizeNone` with `MinWidthClass` and `MinHeightClass`.
- `MinWidthClass`, `MaxWidthClass`, `MinHeightClass`, and `MaxHeightClass` take Tailwind classes that limit the button box.
- `Shape` takes `uiToolset.ShapeCircular`, `uiToolset.ShapeRounded`, or `uiToolset.ShapeSquare`. The default is `rounded-md`.
- The icons use the `ph-bold` weight. Give the icon name only, for example `"ph-floppy-disk"`.
- `IconSizeClass` and `TextSizeClass` take a Tailwind text size class and size the icons and the label on their own, replacing the `Size` tokens.
- `AriaLabel` sets a static accessible name. `AriaLabelOneWayStatePath` binds the name to Alpine state. `AriaExpandedOneWayStatePath` binds `aria-expanded`, for a button that opens a controlled region.
- `IsDisabled` disables the button. `IsDisabledOneWayStatePath` disables and dims it while the path is truthy.
- `IsVisibleOneWayStatePath` hides the button while the path is falsy.
- Every button has a `type`. The default is `button`, so a button in a form does not submit it. `IsSubmit` sets `type="submit"`.
- `TooltipContent` shows a text tooltip. `TooltipContentHtml` takes a component, and `TooltipContentOneWayStatePath` binds live text. `TooltipPosition` takes a `uiToolset.TooltipPosition*` constant. The `Tooltip*WidthClass` and `Tooltip*HeightClass` fields limit the size. The tooltip renders on the document body, so a scroll container cannot clip it.
- `BackgroundColor`, `TextColor`, and `RingColor` take Tailwind color tokens, for example `"red-500/20"`. The `*Hover` fields set the hover colors. `BackgroundColorActive` sets the pressed color.
- `TextCase` takes a `uiToolset.TextCase*` value for the label.

### Action states

`ShouldTrackAction` makes the button follow the action it starts. While the action runs, the button shows a spinner and `PendingLabel`, disables itself, and sets `aria-busy`. When the action ends, the button shows `SucceededLabel` or `FailedLabel` for two seconds. Each label defaults to `Label`.

`PendingIcon`, `SucceededIcon`, and `FailedIcon` replace the action icons. They default to `ph-spinner`, `ph-check-circle`, and `ph-warning-circle`. `PendingIconAnimation` takes a `uiToolset.Animation*` constant and defaults to `AnimationSpin`. `AnimationNone` drops the animation. Any other value passes through as-is, so a custom animation must include its full class name, for example `"animate-wiggle"`.

Without `IsSubmit`, the action is the promise that `OnClickFunc` returns. No form is necessary.

- A resolved promise shows success. A rejected promise shows failure and logs `ButtonActionFailed`.
- `UiToolset.JsonAjax(...)` returns a promise of this kind.
- `OnClickFunc` must be one expression, because the button reads its result. Use `"saveProfile()"`, not `"isOpen = false; saveProfile()"`.
- A result that is not a promise does not change the state.

With `IsSubmit`, the action is the form that holds the button.

- The button refuses a second submit while pending. This includes the Enter key.
- An HTMX request shows the result. A native POST shows no result, because the page navigates away.
- A form with a `target` other than `_self` does not show the native pending state.
- A native response that does not navigate, such as a download or a 204, keeps the button pending. Use HTMX or `IsPendingOneWayStatePath` for these forms.
- A cancelled HTMX request does not start the pending state.
- A page restored from the back-forward cache returns to idle.
- HTMX requests from other elements in the form do not change the button.

`IsPendingOneWayStatePath` shows the pending state while the path is truthy, on any button. Use it when one request blocks several buttons.

## RangeSlider

A slider with one or two thumbs, value bubbles, tick marks, and gradient tracks.

```go
@uiControl.RangeSlider(uiControl.RangeSliderSettings{
    ThumbValueTwoWayStatePath: "sliderValue",
    TrackStartValue:           "0",
    TrackEndValue:             "100",
    TrackSteps:                "1",

    // OptionalFields
    Thumb: uiControl.RangeSliderThumbSettings{
        ValueBubbleEnabled: true,
    },
    Track: uiControl.RangeSliderTrackSettings{
        TicksEnabled: true,
    },
})
```

- `ThumbValueTwoWayStatePath` binds the thumb. `ThumbValue` places it without a binding.
- `TrackStartValue`, `TrackEndValue`, and `TrackSteps` set the track. `TrackStartValueOneWayStatePath` and `TrackEndValueOneWayStatePath` read the bounds from Alpine state instead.
- `Thumb` holds every thumb appearance field: `Shape`, `BackgroundColor`, `TextColor`, `Label`, `LabelIcon`, `LabelOneWayStatePath`, `AriaLabel`, and the `ValueBubble*` fields.
- `Thumb.ValueBubbleEnabled` shows the value above the thumb. `ValueBubbleDisplayMode` takes `RangeSliderThumbValueBubbleDisplayModeAlways` (the default) or `RangeSliderThumbValueBubbleDisplayModeHover`. The other `ValueBubble*` fields style it.
- `Track` holds every track appearance field: the `BackgroundColor*` and `FillColor*` colors, the `StartValueLabel*` and `EndValueLabel*` bounds labels, the `TicksEnabled` and `TickStep` marks, and the `Tooltip*` preview.
- `Track.TooltipEnabled` shows the value under the pointer. The value snaps to `TrackSteps` and stays inside the track. `TooltipPosition` takes `RangeSliderTrackTooltipPositionTop` or `RangeSliderTrackTooltipPositionBottom`.
- `ThumbDualValueModeEnabled` adds a second thumb, bound with `ThumbUpperValueTwoWayStatePath`. `ThumbUpper` takes the same fields as `Thumb`.
- A thumb without an `AriaLabel` uses its visible label, then a default name.
- At init and on each change, the slider moves an out-of-bounds value back inside the track. In dual mode, the thumbs stay one step apart.

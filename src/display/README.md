# Display

Display layer of Infinite UI. It shows content, state, and page chrome. Part of [Infinite UI](../../README.md).

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

Each item takes a `Title`, a `Content` component, and an optional `Icon`.

## Alert

Status message with title, description, and variation presets.

```go
@uiDisplay.Alert(uiDisplay.AlertSettings{
    Title:       "Alert Title",
    Description: "This is an alert message.",
    Variation:   uiDisplay.AlertVariationInfo,
    Size:        uiDisplay.AlertSizeMd,
    IsCloseable: true,
})
```

- `Variation` accepts `AlertVariationSuccess`, `AlertVariationWarning`, `AlertVariationError`, or `AlertVariationInfo`.
- `ContentHtml` replaces the title and description with custom markup.
- `IsCloseable` shows a close button. `DisplayTimeoutSeconds` and `OnTimeoutFunc` auto-dismiss the alert.
- Title and description support live state through the `*OneWayStatePath` fields.

## CloakLoading

Full-screen loader that covers the page before CSS and JavaScript load. Place it at the end of `<body>`.

```go
@uiDisplay.CloakLoading(uiDisplay.CloakLoadingSettings{
    HideDelaySeconds: "1",
    TextMessage:      "Loading...",
    Icon:             "ph-compass-rose",
})
```

It uses inline styles on purpose, so it renders before the stylesheet arrives. `HideDelaySeconds` hides it after the given delay. Icon sizes accept the `CloakLoadingIconSize*` constants, and animations accept the `CloakLoadingAnimationName*` constants.

## LoadingOverlay

Overlay for in-flight requests. It binds to a state path or reacts to the HTMX `htmx-request` class.

```go
@uiDisplay.LoadingOverlay(uiDisplay.LoadingOverlaySettings{
    IsLoadingOneWayStatePath: "isLoading",
    Icon:                     "ph-compass-rose",
    IconSize:                 uiDisplay.LoadingOverlayIconSizeMd,
    AnimationName:            uiDisplay.LoadingOverlayAnimationNameSpin,
})
```

Set `hx-indicator="#loading-overlay"` on HTMX elements, or toggle it with `UiToolset.ToggleLoadingOverlay()`.

## Modal

Dialog with header, middle, and footer slots.

```go
@uiDisplay.Modal(uiDisplay.ModalSettings{
    MiddleContent:            ModalDemoContent(),
    Title:                    "Modal Title",
    Size:                     uiDisplay.ModalSizeMd,
    IsVisibleTwoWayStatePath: "isModalVisible",
})
```

- The standard header appears when `Title` is set and `HeaderContent` is nil. It carries resize and close buttons.
- `IsUnresizable` and `IsUncloseable` remove those buttons.
- The close button and a backdrop click set the visibility path to false and run `OnCloseFunc`. `IsUncloseable` stops both; the app then closes the modal by setting the visibility path to false.
- `OnResizeFunc` runs when the resize button is clicked.

## Tag

Composite label with an outer ring and an inner background.

```go
@uiDisplay.Tag(uiDisplay.TagSettings{
    OuterLeftIcon:  "ph-info",
    OuterLeftLabel: "Info",
    InnerIcon:      "ph-warning",
    InnerLabel:     "Warning",
    Size:           uiDisplay.TagSizeXs,
})
```

- `OuterLeft*`, `OuterRight*`, and `Inner*` fields take an icon, a static label, or a `*OneWayStatePath` for live text.
- `InnerValueOneWayStatePath` binds a live value in the inner segment without the lowercase transform, for client-side value tags: set a static `OuterLeftLabel` as the label and the path as the value.
- `Size` accepts `TagSizeTiny` through `TagSizeXl`; `TagSizeTiny` renders a dense layout for table cells.
- Set `OnRemoveFunc` to render a removable chip. `RemoveButtonLabel` names the remove button.
- `OuterBackgroundColor`, `OuterRingColor`, and `OuterTextColor` accept color tokens, for example `"neutral-50/10"`.

## Toast

Notification popup driven by an Alpine global store.

```go
@uiDisplay.Toast(uiDisplay.ToastSettings{
    BackgroundColor:    "neutral-800",
    TextColor:          "neutral-50",
    Size:               uiDisplay.ToastSizeMd,
    AutoDismissSeconds: 10,
})
```

Render one Toast per page. Then call `$store.toast.displayToast(message, type)` or `$store.toast.displayToastWithApiResponse(apiResponse, httpStatusCode)`. Call `$store.toast.clearToast()` to hide it early.

When HTMX is present, the Toast listens to `htmx:afterRequest`. A JSON response with a `readableMessage`, a `humanReadableMessage` outcome, or a string `body` displays a toast automatically. The styling follows the HTTP status: 2xx success, 207 partial success, 4xx and above danger. Auto-dismiss defaults to 10 seconds.

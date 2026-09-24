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

`BorderRadius` rounds only the outer edges: the first item's top corners and the last item's bottom corners. Middle items stay square.

`ContentBackgroundColor` sets the color behind an open item's content. The default is `neutral-900`. Match it to the surface that holds the accordion. The content then reads as a cutout of that surface, with rounded inner edges and padding.

Set `IsSingleOpen` to keep one item open at a time. The component groups the items under a shared name, so opening one closes the rest.

`TextCase` accepts a `uiToolset.TextCase*` value and transforms the item titles. The default, `TextCaseNone`, leaves them as typed.

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
- `Title` and `Description` render as raw HTML. Pass only trusted content. Never pass request data into these fields.
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

## ConfirmationDialog and Presets

Confirmation dialogs over Modal with four presets: `ConfirmDialog`, `WarningDialog`, `CriticalDialog`, and `DeleteDialog`. Each preset supplies its icon, tone, and confirmation copy. `CriticalDialog` and `DeleteDialog` enable the type-to-confirm gate and have no fullscreen mode.

```go
@uiDisplay.DeleteDialog(uiDisplay.ConfirmationDialogSettings{
    IsVisibleTwoWayStatePath: "isDeleteDialogVisible",
    OnConfirmFunc:            "deleteRecord()",

    // OptionalFields
    TargetNameStatePath: "record.name",
    TargetIdStatePath:   "record.id",
})
```

- `Size` accepts `ConfirmationDialogSizeXs` (40%) through `ConfirmationDialogSizeXxl` (90%) as viewport percentages on both axes, like Modal. There is no fullscreen size and no resize control. `WidthPercent` and `HeightPercent` override the panel box like Modal. The default height hugs the content with an `85%` ceiling; set `HeightPercent` to pin it.
- The header block carries `HeaderTitle`, `HeaderSubHeading`, `HeaderTitleOneWayStatePath`, `HeaderSubHeadingOneWayStatePath`, `HeaderTitleColor`, `HeaderSubHeadingColor`, `HeaderSize`, `TextCase`, and the `HeaderIcon*` icon controls. `ActionsContent` adds buttons next to the close button. `HeaderContent` replaces the whole header block. `TextColor` sets the dialog text color. Each preset fills the title and icon defaults.
- `HeaderIconPosition` defaults to `HeaderIconPositionTop`, which centers the icon above the title. Set `HeaderIconPositionLeft` to place it beside the title.
- `TargetNameStatePath` and `TargetIdStatePath` bind live values into the message. The name renders as a dashed chip and the id as `#<id>`.
- `IsTypeToConfirmEnabled` turns on the type-to-confirm gate. `CriticalDialog` and `DeleteDialog` set it. The confirm button enables only when the typed value matches the target name (or the target id when no name path is set). An empty or missing target keeps the confirm button disabled. `TypeToConfirmValueTwoWayStatePath` moves the typed value into your own Alpine state; the default keeps it inside the dialog.
- Two type-to-confirm dialogs on one page each hold their own typed value, but both fields use the input name `typedConfirmationValue`. Keep them out of one shared submitted form.
- Pass state paths from code, never from request data. The component embeds them into client-side expressions.
- `Message` replaces the note under the question. `MessageContent` replaces the whole message block.
- `CancelFunc` defaults to closing the dialog. `OnConfirmFunc` runs as written, so append the close when the action succeeds.
- The cancel button renders outlined with a light ring. The confirm button carries the preset tone fill, so the eye lands on it first.
- `HeaderTitle`, `ConfirmButtonLabel`, `ConfirmButtonIcon`, and `CancelButtonLabel` override the preset copy.

## HeaderBlock

Shared header row behind PageHeading, Card, and ConfirmationDialog: icon, heading, sub-heading, and a right-aligned actions slot.

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

- `HeadingLevel` accepts `HeaderBlockHeadingLevelOne` or `HeaderBlockHeadingLevelTwo`; the default is two.
- `HeaderSize` accepts `HeaderSizeXs` through `HeaderSizeXl` and scales the title, icon, and sub-heading together.
- The `HeaderIcon*` fields match `HeaderIcon`. `HeaderIconPosition` accepts `HeaderIconPositionLeft` or `HeaderIconPositionTop`.
- `HeaderTitleColor`, `HeaderSubHeadingColor`, and `TextColor` take color tokens.
- `TextCase` accepts a `uiToolset.TextCase*` value and transforms the title and sub-heading. The default, `TextCaseNone`, leaves them as typed. `Card`, `PageHeading`, and `ConfirmationDialog` forward it.

## HeaderIcon

Icon with an optional background chip.

```go
@uiDisplay.HeaderIcon(uiDisplay.HeaderIconSettings{
    Name: "ph-cube",

    // OptionalFields
    Color:           "red-950",
    BackgroundColor: "neutral-300",
    BorderRadius:    uiDisplay.HeaderIconBorderRadiusXl,
    PaddingSize:     uiDisplay.HeaderIconPaddingSizeLg,
})
```

- `Color` sets the icon text color. `BackgroundColor` sets the chip behind it. Icon placement (left of the title or above it) belongs to `HeaderBlock`'s `HeaderIconPosition`.
- `BorderRadius` accepts `HeaderIconBorderRadiusNone` through `HeaderIconBorderRadiusXl`. `PaddingSize` accepts `HeaderIconPaddingSizeNone` through `HeaderIconPaddingSizeXl`.

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

Window-style dialog with header, middle, and footer slots.

```go
@uiDisplay.Modal(uiDisplay.ModalSettings{
    MiddleContent:            ModalDemoContent(),
    Title:                    "Modal Title",
    InitialSize:              uiDisplay.ModalSizeMd,
    IsVisibleTwoWayStatePath: "isModalVisible",
})
```

- The standard header appears when `Title` is set and `HeaderContent` is nil. It carries the enlarge, reduce, and close buttons. A custom `HeaderContent` still gets the same action cluster beside it.
- `InitialSize` sets the panel box as a viewport percentage on both axes: `ModalSizeXs` (40%), `ModalSizeSm` (50%), `ModalSizeMd` (60%), `ModalSizeLg` (70%), `ModalSizeXl` (80%), `ModalSizeXxl` (90%), and `ModalSizeFull` (100%). Width and height grow together, so the box follows the screen's aspect ratio and a widescreen gets more width than height. `ModalSizeXxl` is the near-full size for wide content such as a web terminal. InitialSize also scales the header title and padding.
- `WidthPercent` and `HeightPercent` replace the box on that axis with a viewport percentage. Set one to change a single dimension, or both to pin the panel. Values are whole percentages, 1 through 100; 100 renders a full axis. When both are set, the resize controls retire.
- `PossibleSizes` sets the ordered sizes the panel can take; the default is md, lg, xl, xxl. `InitialSize` is the first reachable entry: any entry before it is dropped, and an `InitialSize` outside the list falls to the first entry.
- `CurrentSizeTwoWayStatePath` binds the live size so the enlarge and reduce controls can drive it. The modal creates an internal path when the field is empty and the panel is resizable.
- Enlarge steps up one entry, reduce steps back one. Reduce hides at the initial size and enlarge hides at the largest reachable size.
- `IsHeightContentSized` renders the panel at its content height with an `85%` ceiling, ignoring the size height. Use it for dialogs that must stay short.
- `IsMiddleContentScrollDisabled` stops the middle region from scrolling and lets the content own the overflow, for a full-bleed panel such as a terminal.
- `IsUnresizable` and `IsUncloseable` remove those buttons. The resize controls also retire when the panel holds fewer than two reachable sizes.
- The close button and a backdrop click set the visibility path to false and run `OnCloseFunc`. `IsUncloseable` stops both; the app then closes the modal by setting the visibility path to false.
- `OnResizeFunc` runs when either resize control is clicked.

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
- `TextCase` accepts a `uiToolset.TextCase*` value and transforms the label segments. The default, `TextCaseNone`, leaves them as typed.
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

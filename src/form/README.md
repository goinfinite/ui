# Form

Form layer of Infinite UI. It collects user input: text, choices, and switches. Every component binds to Alpine.js state; field components render their label as a floating legend. Part of [Infinite UI](../../README.md).

Every labelled form component accepts `TextCase` with a `uiToolset.TextCase*` value: `TextCaseNone` (the default, as typed), `TextCaseLower`, `TextCaseUpper`, or `TextCaseCapitalize`. It transforms the label, and the placeholder where the component renders one.

## CheckboxInput

Checkbox with a label, bound to a boolean or an array state path.

```go
@uiForm.CheckboxInput(uiForm.CheckboxInputSettings{
    Label:           "Accept the terms",
    TwoWayStatePath: "hasAcceptedTerms",

    // OptionalFields
    InputName: "hasAcceptedTerms",
    Shape:     uiToolset.ShapeRounded,
})
```

`Shape` accepts `uiToolset.ShapeSquare`, `uiToolset.ShapeRounded` (the default), and `uiToolset.ShapeCircular`. `Size` accepts the `CheckboxInputSize*` constants. `CheckedColor`, `UncheckedColor`, and `FocusRingColor` set the box colors; the box border previews `CheckedColor` on hover. `LabelPosition` accepts `CheckboxInputLabelPositionLeft` and `CheckboxInputLabelPositionRight`. `IsChecked`, `IsDisabled`, and `IsRequired` render static states. `IsCheckedOneWayStatePath`, `IsDisabledOneWayStatePath`, and `IndeterminateOneWayStatePath` bind them to Alpine state. `AriaLabel` names the checkbox when no visible label is present.

`IsInvalid` and `IsInvalidOneWayStatePath` flag the error state. The box border turns to `ErrorColor` (default `"red-500"`), and the input carries `aria-invalid`. `ErrorMessage` and `ErrorMessageOneWayStatePath` render the message below the control. With `InputId` set, the input links to it through `aria-describedby`. When an invalid path is set, the message shows only while the path is true.

## InputField

Single-line input with a floating label, optional affixes, and an optional hint.

```go
@uiForm.InputField(uiForm.InputFieldSettings{
    InputType:       uiForm.InputTypeText,
    InputName:       "name",
    Label:           "Name",
    TwoWayStatePath: "name",
    IsRequired:      true,
})
```

`InputType` accepts the `uiForm.InputType*` constants. `IconLeft` takes a Phosphor icon class and renders it inside the field before the input, for example `"ph-magnifying-glass"`. An affix renders a static value (`AffixLeftValue`, `AffixRightValue`) or binds a state path (`AffixLeftStatePath`, `AffixRightStatePath`). Set `AffixLeftWidthPercent` or `AffixRightWidthPercent` to fix an affix to a percentage of the field width, for example `25` or `50`. `AffixRightComponent` renders any templ component at the right edge, after the right affix. `PasswordInput` uses it for its action buttons. `InputTypeExpression` supplies the type from an Alpine expression and overrides `InputType` at runtime, for example `"isPasswordVisible ? 'text' : 'password'"`. Number inputs accept `InputNumberMin`, `InputNumberMax`, and `InputNumberStep`. `Size` accepts the `InputFieldSize*` constants and defaults to `md`. It scales the input, the affixes, and the floating label together. `TextCase` accepts a `uiToolset.TextCase*` value. It transforms the floating label and the placeholder. The default, `TextCaseNone`, leaves the text as typed.

## PasswordInput

Password field with a reveal toggle, an optional random password generator, and an optional strength meter. It composes `InputField`, so every `InputField` setting travels through.

```go
@uiForm.PasswordInput(uiForm.PasswordInputSettings{
    InputName:       "password",
    Label:           "Password",
    TwoWayStatePath: "password",

    // OptionalFields
    ShouldShowGenerateButton: true,
    ShouldShowStrengthMeter:  true,
})
```

`Rules` configures the generator and the meter: `MinLength` (default 6), `MaxLength` (default 64), `GenerationLength` (default 16), and one `ShouldInclude*Chars` flag per character class. When no class is selected, every class is required. The generator runs on `UiToolset.CreateRandomPassword(options)` and guarantees at least one character from each selected class. The meter shows a percentage bar and a criteria checklist that updates as the user types; the checklist shows only the enabled criteria. The generate button copies the new password to the clipboard and shows a success toast. The toast text follows the browser language and falls back to English; the page must mount `Toast`. `ActionButtonStyle` accepts `PasswordInputActionButtonStyleBoxed` (the default) or `PasswordInputActionButtonStylePlain` for bare icons. `StrengthMeterColor`, `StrengthCriteriaTextColor`, `StrengthCriteriaFulfilledIconColor`, and `StrengthCriteriaUnfulfilledIconColor` take a color token and repaint the meter bar, the checklist text, and the fulfilled and unfulfilled icons.

## InlineCheckboxGroup

A row or column of checkboxes under one shared label. The label notches into the top border of the field, matching `InlineRadioGroup` and the other notched fields.

```go
@uiForm.InlineCheckboxGroup(uiForm.InlineCheckboxGroupSettings{
    Label: "Fruits",
    InputSettings: []uiForm.CheckboxInputSettings{
        {Label: "Apple", Value: "apple", TwoWayStatePath: "selectedFruits", InputName: "selectedFruits"},
        {Label: "Banana", Value: "banana", TwoWayStatePath: "selectedFruits", InputName: "selectedFruits"},
    },
})
```

Each `CheckboxInputSettings` carries its own size, shape, colors, and errors. When the shared state path holds an array, Alpine adds and removes the checked values. `Orientation` accepts `InlineCheckboxGroupOrientationHorizontal` (the default) and `InlineCheckboxGroupOrientationVertical`. `MaxVisibleOptions` caps the visible option rows and scrolls the rest. `TextCase` transforms only the shared label.

## TextArea

Multiline input with expand, copy, and clear actions.

```go
@uiForm.TextArea(uiForm.TextAreaSettings{
    Label:           "Description",
    TwoWayStatePath: "description",
})
```

Set `IsCode` for monospace text. `IsReadOnly` sets the native `readonly` attribute. `IsRequired` adds the required marker to the label.

## SelectInput

Dropdown select. It is not a native `<select>`: screen-reader-only radio inputs and an Alpine-driven list carry the value. Three option modes exist:

1. `FlatOptions []string` for plain values.
2. `LabelValueOptions []SelectLabelValueOption` for label and value pairs.
3. `LabelValueOptions` with `LabelHtml` for rich option content.

```go
@uiForm.SelectInput(uiForm.SelectInputSettings{
    InputName:       "status",
    Label:           "Status",
    FlatOptions:     []string{"running", "stopped"},
    TwoWayStatePath: "status",
})
```

`IsMultiSelect` renders checkbox rows, keeps the dropdown open after each toggle, and binds an array on `TwoWayStatePath`; the trigger joins the selected labels. `InputId` prefixes the option checkbox ids in multi mode. `ShouldIncludeBlankOption` adds a blank entry in single mode; multi mode shows the clear control whenever the array holds a value. `IsDisabledOneWayStatePath` disables the trigger while the path is truthy. `OnChangeFunc` runs after a change, a toggle, or a clear. `Size` accepts the `SelectInputSize*` constants and defaults to `md`. `DropdownMinHeightClass` and `DropdownMaxHeightClass` cap the list height (default `max-h-60`). `MaxVisibleOptions` caps it to a row count and wins over both.

## SearchableSelectInput

Dropdown with a filter box. Typing narrows the local options in place. The list caps its height and scrolls. A clear button empties the selection. Use it past roughly ten options, where a native select becomes hard to scan.

```go
@uiForm.SearchableSelectInput(uiForm.SearchableSelectInputSettings{
    InputName:       "country",
    Label:           "Country",
    FlatOptions:     []string{"Argentina", "Brazil", "Chile"},
    TwoWayStatePath: "country",
})
```

It shares the `FlatOptions` and `LabelValueOptions` modes with `SelectInput`. The filter matches the label, case-insensitively. `IsMultiSelect` binds an array and keeps the dropdown open after each toggle. `SelectionDisplay` accepts `SearchableSelectInputSelectionDisplayText` (the default) or `SearchableSelectInputSelectionDisplayTags` for removable tags inside the field. The text display joins the selected labels. With tags, `ShouldAllowCustomValues` turns a typed value into a tag on Enter, and Backspace on an empty input removes the last tag. Each tag shows the option label while the bound array holds the stored value. Tags apply to multi mode only. `TagOuterBackgroundColor` (default `secondary-500`) paints the tag and its border, `TagOuterTextColor` (default `neutral-50`) paints its text, and `TagInnerBackgroundColor` (default `neutral-50/10`) paints the label bubble. `DropdownMinHeightClass` and `DropdownMaxHeightClass` cap the list (default `max-h-60`). `MaxVisibleOptions` caps it to a row count and wins over both. `OnChangeFunc` runs on selection, on a multi-select toggle, on a tag change, and on clear. Single mode writes one hidden input. Multi mode writes one hidden input per selected value, so a form submission carries every value under `InputName`.

## RemoteSearchableSelectInput

Searchable select whose options come from a URL as the user types. It debounces the requests and waits for the minimum query length. It discards stale responses and shows loading, empty, and error states.

```go
@uiForm.RemoteSearchableSelectInput(uiForm.RemoteSearchableSelectInputSettings{
    InputName:       "country",
    Label:           "Country",
    OptionsUrl:      "/api/countries/search",
    TwoWayStatePath: "country",
})
```

`OptionsQueryParam` names the query parameter (default `q`). `MinQueryLength` gates the first request (default 3). `DebounceMs` waits before each request (default 300). Below the minimum query length, including an empty query, the dropdown shows the minimum-length prompt. `DropdownMinHeightClass`, `DropdownMaxHeightClass`, and `MaxVisibleOptions` size the list like the local variant. The endpoint receives `GET OptionsUrl?<query param>=<typed text>`. It returns `{ "body": [{ "label": "...", "value": "..." }] }`. A bare JSON array of strings or `{label, value}` objects also works. `IsMultiSelect` binds an array and caches the selected labels while the query changes. `OnChangeFunc` runs on selection, on a multi-select toggle, and on clear.

## RadioInput

Single radio option. Group several inputs by sharing one `TwoWayStatePath`.

```go
@uiForm.RadioInput(uiForm.RadioInputSettings{
    Label:           "Option 1",
    StateValue:      "option1",
    TwoWayStatePath: "selectedOption",
})
```

## InlineRadioGroup

A row of radios with one shared label. The label notches into the top border of the field, matching `InputField`, `SelectInput`, and the other notched fields.

```go
@uiForm.InlineRadioGroup(uiForm.InlineRadioGroupSettings{
    Label: "Select an option",
    InputSettings: []uiForm.RadioInputSettings{
        {Label: "Option 1", StateValue: "option1", TwoWayStatePath: "groupSelection"},
        {Label: "Option 2", StateValue: "option2", TwoWayStatePath: "groupSelection"},
    },
})
```

`Orientation` accepts `InlineRadioGroupOrientationHorizontal` (the default) and `InlineRadioGroupOrientationVertical`. `MaxVisibleOptions` caps the visible option rows and scrolls the rest. `InlineRadioGroup`'s `TextCase` transforms only the shared label. Each option label takes its own `RadioInputSettings.TextCase`.

## ToggleSwitch

Boolean switch. Bind one state path, or set `CustomValue` to collect one value into an array state path.

```go
@uiForm.ToggleSwitch(uiForm.ToggleSwitchSettings{
    Label:           "Enable notifications",
    TwoWayStatePath: "isNotificationsEnabled",
})
```

`LabelPosition` accepts `ToggleSwitchLabelPositionLeft` and `ToggleSwitchLabelPositionRight`. The `*Color` settings override the track and thumb colors.

## InputHint

Shared hint helper for the fields above. It renders no field of its own.

- `HintDisplay` accepts `InputHintDisplayTooltip` for an icon with a tooltip, or `InputHintDisplayDescription` for a line below the field.
- `HintValue` carries static text; `HintStatePath` binds live text.
- `HintDisplayTooltipMinWidthClass`, `HintDisplayTooltipMaxWidthClass`, `HintDisplayTooltipMinHeightClass`, and `HintDisplayTooltipMaxHeightClass` cap the hint tooltip on every field that offers hints. `HintDisplayTooltipBackgroundColor` paints it.
- The tooltip display renders through the shared `uiToolset` tooltip, so it teleports to the document body and a scroll container or a modal cannot clip it.
- `InputHintTooltip` and `InputHintDescription` can also render directly, when a hint sits next to custom markup.

## Non-obvious behaviors

- `SelectInput` renders a `templ.JSONScript` block for label-value options. The block feeds the selected label lookup and the multi-select summary.
- `SearchableSelectInput` and `RemoteSearchableSelectInput` render one `templ.JSONScript` config block and share one Alpine state, `searchableSelectInput`. The config carries the normalized options, the multi-select and custom-value flags, and the remote settings.
- The floating legend collapses while the field is empty, so the empty field shows the label as a placeholder.
- `InputField` hides its empty legend with `display: none`. Chrome reserves scroll space for a zero-sized legend, so the opacity-based collapse the other field components use can phantom-scroll an overflow container.
- A dropdown opens upward when the space below the trigger is too small for it. The space stops at the bottom of the viewport and at the bottom of every clipping ancestor, so the list never opens downward into an `overflow-hidden` panel such as the carousel surface. `UiToolset.SelectDropdown.openUpwardResolver` makes the decision for `SelectInput` and the searchable fields.
- `UiToolset.CreateRandomPassword` accepts an options object: `length`, `minLength`, `maxLength`, and one `include*` flag per character class. With no arguments it keeps the original 16-character behavior.
- `InputName` sets the key in an HTMX form submission. `InputId` is optional.

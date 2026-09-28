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
    Shape:     uiForm.CheckboxInputShapeRounded,
})
```

`Shape` accepts `CheckboxInputShapeSquare`, `CheckboxInputShapeRounded` (the default), and `CheckboxInputShapeCircular`. `Size` accepts the `CheckboxInputSize*` constants. `CheckedColor`, `UncheckedColor`, and `FocusRingColor` set the box colors. `LabelPosition` accepts `CheckboxInputLabelPositionLeft` and `CheckboxInputLabelPositionRight`. `IsChecked`, `IsDisabled`, and `IsRequired` render static states; `IsCheckedOneWayStatePath`, `IsDisabledOneWayStatePath`, and `IndeterminateOneWayStatePath` bind them to Alpine state. `AriaLabel` names the checkbox when no visible label is present.

`IsInvalid` and `IsInvalidOneWayStatePath` flag the error state: the box border turns to `ErrorColor` (default `"red-500"`), and the input carries `aria-invalid`. `ErrorMessage` and `ErrorMessageOneWayStatePath` render the message below the control; with `InputId` set, the input links to it through `aria-describedby`. When an invalid path is set, the message shows only while the path is true.

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

`InputType` accepts the `uiForm.InputType*` constants. An affix renders a static value (`AffixLeftValue`, `AffixRightValue`) or binds a state path (`AffixLeftStatePath`, `AffixRightStatePath`). Set `AffixLeftWidthPercent` or `AffixRightWidthPercent` to fix an affix to a percentage of the field width, for example `25` or `50`. Number inputs accept `InputNumberMin`, `InputNumberMax`, and `InputNumberStep`. `Size` accepts the `InputFieldSize*` constants and defaults to `md`; it scales the input, the affixes, and the floating label together. `TextCase` accepts a `uiToolset.TextCase*` value; it transforms the floating label and the placeholder. The default, `TextCaseNone`, leaves the text as typed.

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

`ShouldIncludeBlankOption` adds a clear entry. `IsDisabledOneWayStatePath` disables the trigger while the path is truthy. `OnChangeFunc` runs after a change. `Size` accepts the `SelectInputSize*` constants and defaults to `md`.

## MultiSelectInput

Dropdown that holds a list of selected values on one state path.

```go
@uiForm.MultiSelectInput(uiForm.MultiSelectInputSettings{
    InputName:       "countries",
    Label:           "Countries",
    FlatOptions:     []string{"Argentina", "Brazil", "Chile"},
    TwoWayStatePath: "countries",
})
```

It shares the `FlatOptions` and `LabelValueOptions` modes with `SelectInput`. `OnChangeFunc` runs after a checkbox toggle or a clear.

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

A row of radios with one shared label.

```go
@uiForm.InlineRadioGroup(uiForm.InlineRadioGroupSettings{
    Label: "Select an option",
    InputSettings: []uiForm.RadioInputSettings{
        {Label: "Option 1", StateValue: "option1", TwoWayStatePath: "groupSelection"},
        {Label: "Option 2", StateValue: "option2", TwoWayStatePath: "groupSelection"},
    },
})
```

`InlineRadioGroup`'s `TextCase` transforms only the shared label. Each option label takes its own `RadioInputSettings.TextCase`.

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
- `InputHintTooltip` and `InputHintDescription` can also render directly, when a hint sits next to custom markup.

## Non-obvious behaviors

- `SelectInput` and `MultiSelectInput` render a `templ.JSONScript` block for label-value options. The block feeds the selected label lookup.
- The floating legend collapses while the field is empty, so the empty field shows the label as a placeholder.
- `InputField` hides its empty legend with `display: none`. Chrome reserves scroll space for a zero-sized legend, so the opacity-based collapse the other field components use can phantom-scroll an overflow container.
- A `SelectInput` dropdown opens upward when it would overflow the bottom of the viewport.
- `InputName` sets the key in an HTMX form submission. `InputId` is optional.

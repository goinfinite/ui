# Form

The form package collects user input: text, choices, numbers, files, and switches. Each component binds to Alpine.js state. A field shows its label as a floating legend. Part of [Infinite UI](../../README.md).

Each labelled component takes `TextCase` with a `uiToolset.TextCase*` value: `TextCaseNone` (the default), `TextCaseLower`, `TextCaseUpper`, or `TextCaseCapitalize`. It changes the label and, where present, the placeholder.

## CheckboxInput

A checkbox with a label. It binds to a boolean or to an array.

```go
@uiForm.CheckboxInput(uiForm.CheckboxInputSettings{
    Label:           "Accept the terms",
    TwoWayStatePath: "hasAcceptedTerms",

    // OptionalFields
    InputName: "hasAcceptedTerms",
    Shape:     uiToolset.ShapeRounded,
})
```

- `Shape` takes `uiToolset.ShapeSquare`, `uiToolset.ShapeRounded` (the default), or `uiToolset.ShapeCircular`. `Size` takes a `CheckboxInputSize*` constant.
- `CheckedColor`, `UncheckedColor`, and `FocusRingColor` set the box colors. On hover, the border shows `CheckedColor`.
- `LabelPosition` takes `CheckboxInputLabelPositionLeft` or `CheckboxInputLabelPositionRight`.
- `IsChecked`, `IsDisabled`, and `IsRequired` set static states. `IsCheckedOneWayStatePath`, `IsDisabledOneWayStatePath`, and `IndeterminateOneWayStatePath` bind them to Alpine state.
- `AriaLabel` names a checkbox that has no visible label. `AriaLabelOneWayStatePath` binds that name to an Alpine path.
- `LabelOneWayStatePath` renders the visible label from an Alpine path. `Label` remains the static fallback.
- `IsInvalid` and `IsInvalidOneWayStatePath` set the error state. The border changes to `ErrorColor` (default `"red-500"`), and the input gets `aria-invalid`.
- `ErrorMessage` and `ErrorMessageOneWayStatePath` show a message below the box. With an invalid path, the message shows only while the path is true. With `InputId`, the input points to the message through `aria-describedby`.

## CronSchedule

A cron expression control. A mode switch picks between a predefined-schedule dropdown and a five-field custom builder: minute, hour, day, month, weekday. Both modes write one cron string to `TwoWayStatePath`.

```go
@uiForm.CronSchedule(uiForm.CronScheduleSettings{
    Label:           "Cron Schedule",
    TwoWayStatePath: "schedule",
    InputName:       "schedule",
})
```

- The predefined dropdown carries the macro entries (hourly, weekly, monthly, annually, on reboot) and common fixed expressions.
- Each custom field takes a preset from its own select or a typed value. The select shows `Custom` while the typed value matches no preset.
- Each field validates its cron syntax (minute 0-59, hour 0-23, day 1-31, month 1-12, weekday 0-7, with `*`, ranges, lists, and steps). An invalid field turns red and the bound path clears until the schedule is valid again.
- A preview below the fields reads the resulting schedule in plain words, for example `At 09:30, on weekdays`.
- The bound path seeds the control: a preset expression opens in predefined mode, and a custom expression opens in custom mode with the fields filled.
- `Size` takes an `InputFieldSize*` constant. `TextCase` changes the label, the field labels, and the preview.
- `Field` takes a `CronScheduleFieldStyle` with the outer `BackgroundColor`, `BorderColor`, and `BorderRadius`. `Dropdown` takes a `CronScheduleDropdownStyle` with `BackgroundColor`, `BorderColor`, and `BorderRadius` for the predefined list and the five field lists. `Preview` takes a `CronSchedulePreviewStyle` with `TextColor` for the plain-words line and its clock icon, and `CodeBackgroundColor` and `CodeTextColor` for the expression chip. `RadioStyle` takes a `RadioInputStyle` for the mode radios. The predefined trigger and the five custom field surfaces keep their neutral style.

## DateTimePicker

A text input that opens a calendar popover. Modes cover date, time, and datetime. The picker writes the selected value to a state path as an ISO 8601 string.

```go
@uiForm.DateTimePicker(uiForm.DateTimePickerSettings{
    Label:           "Start Date",
    TwoWayStatePath: "startDate",

    // OptionalFields
    Mode:           uiForm.DateTimePickerModeDate,
    MinValue:       "2026-01-01",
    MaxValue:       "2026-12-31",
    Locale:         "en-US",
    FirstDayOfWeek: "sunday",
})
```

- `Mode` takes `DateTimePickerModeDate` (default), `DateTimePickerModeDateTime`, or `DateTimePickerModeTime`.
- The picker stores ISO 8601 values: `YYYY-MM-DD` in date mode, `YYYY-MM-DDTHH:MM` in datetime mode, and `HH:MM` in time mode. The input shows the value in the locale format.
- Typing an ISO 8601 value commits it. Any other text reverts to the stored value on blur or Enter. A typed offset is dropped and the configured timezone supplies the offset again.
- `MinValue` and `MaxValue` disable the days outside the bounds. Use the same ISO format as the mode.
- `Locale` takes a BCP 47 tag for the month names, the weekday names, and the input format (default `en-US`). `FirstDayOfWeek` shifts the week start by weekday name, for example `"monday"`.
- `IsRange` shows two inputs split by a divider, one for the begin value and one for the end value. Each side carries a caption from `RangeStartLabel` and `RangeEndLabel` (default `Start` and `End`). Each input edits its own side and commits on its own; the active input stays bright, and the calendar edits the active side. The start writes to `TwoWayStatePath`, and the end writes to `RangeEndTwoWayStatePath`. The calendar paints the span as one continuous band and rounds only the two ends. Clicking a day before the start swaps the ends, and a typed value that would cross the other side swaps the ends too.
- `Value` and `RangeEndValue` seed the picker when no state path is set.
- The picker seeds its drafts from the bound paths at mount and writes each change back. A path written from outside the picker stays until the next user change.
- Give the state paths from code, never from request data. The component puts them into client-side expressions.
- `InputName` adds a hidden input that submits the ISO value with the form. In range mode the end input takes `RangeEndInputName` or `InputName + "-end"`.
- The time selector sits centered in the popover and renders through `SelectInput`. `TimeSelectorSize` takes a `SelectInputSize*` constant. `TimeSelector` takes a `DateTimePickerTimeSelectorStyle` group with `BackgroundColor`, `BorderColor`, `BorderRadius`, `TextColor`, `DropdownBackgroundColor`, `DropdownBorderColor`, and `DropdownBorderRadius`.
- `Timezone` takes an IANA name, for example `"America/Sao_Paulo"`. In datetime mode, setting it shows a timezone select beside the time when the field is wide enough and wraps it below when it is not, and writes the value with its UTC offset, for example `2026-10-16T09:30-03:00`. Changing the timezone recomposes the current drafts. `TimezoneOptions` replaces the built-in zone list, and a configured zone outside the list is added at the top. `TimezoneTwoWayStatePath` seeds the selected zone from the app state and writes each change back. `Timezone` supplies the zone when the path is empty. `TimezoneInputName` adds a hidden input for form submission.
- Arrow keys move day focus by one day (left/right) or one week (up/down). Home and End jump to the week edges. Escape closes the popover from anywhere on the page.
- `Size` takes an `InputFieldSize*` constant.
- `LabelOneWayStatePath` renders the visible label and the accessible names from an Alpine path. In range mode it also drives the Start and End names. `Label` remains the static fallback.
- `TextCase` takes a `uiToolset.TextCase*` value. It changes the label, the range captions, the inputs, the time selector labels, and the Done button.
- `CalendarButtonStyle` takes `InputActionButtonStyleBoxed` (the default) or `InputActionButtonStylePlain` for the calendar toggle.
- The Done action renders with the `ph-check` icon and lowercase text by default; a set `TextCase` overrides the lowercase default. `DoneButtonLabel`, `DoneButtonIconLeft`, `DoneButtonIconRight`, `DoneButtonSize`, and `DoneButtonStyle` customize it. `DoneButtonStyle` takes `BackgroundColor`, `BackgroundColorHover`, `IconSizeClass`, `RingColor`, `RingColorHover`, `RingThickness`, `Shape`, `TextColor`, `TextColorHover`, and `TextSizeClass`. `IconSizeClass` and `TextSizeClass` take a Tailwind text size class and size the icon and the label on their own; the default icon size follows the button size.
- `HintValue` and `HintStatePath` carry the hint text, and `HintDisplay` picks the display mode: `InputHintDisplayTooltip` (default when set) or `InputHintDisplayDescription`. The tooltip settings mirror the InputField, including `HintIconStyle`.
- `Style` takes a `DateTimePickerStyle` group. The field fields are `FieldBackgroundColor`, `FieldBorderColor`, `FieldBorderRadius`, `TextColor`, `PlaceholderColor`, `CalendarButtonIconColor`, `RangeDividerColor`, and `RangeLabelTextColor`; the calendar surface fields are `CalendarBackgroundColor`, `CalendarBorderColor`, and `CalendarBorderRadius`; the calendar day fields are `DayTextColor`, `DayMutedTextColor`, `DayDisabledTextColor`, `DayHoverBackgroundColor`, `DaySelectedBackgroundColor`, `DaySelectedTextColor`, `DayRangeBackgroundColor`, `DayRangeTextColor`, and `DayTodayRingColor`. Every field takes a Tailwind token, and the defaults match the neutral and secondary tones.
- The calendar navigation, the toggle, and the Done action render through `uiControl.Button`.

## InputField

A single-line input with a floating label, optional affixes, and an optional hint.

```go
@uiForm.InputField(uiForm.InputFieldSettings{
    InputType:       uiForm.InputTypeText,
    InputName:       "name",
    Label:           "Name",
    TwoWayStatePath: "name",
    IsRequired:      true,
})
```

- `InputType` takes a `uiForm.InputType*` constant. `InputTypeExpression` sets the type from an Alpine expression and replaces `InputType`, for example `"isPasswordVisible ? 'text' : 'password'"`.
- `IconLeft` shows a Phosphor icon before the input, for example `"ph-magnifying-glass"`.
- `LabelOneWayStatePath` renders the legend, the placeholder, and the accessible name from an Alpine path. `Label` remains the static fallback.
- An affix shows a static value (`AffixLeftValue`, `AffixRightValue`) or a state path (`AffixLeftStatePath`, `AffixRightStatePath`). `AffixLeftWidthPercent` and `AffixRightWidthPercent` fix the affix width, for example `25`.
- `AffixRightComponent` renders a templ component at the right edge. `PasswordInput` uses it for its buttons.
- A number input takes `InputNumberMin`, `InputNumberMax`, and `InputNumberStep`.
- `IsInvalidOneWayStatePath` paints the field border red and sets `aria-invalid` while the path is truthy.
- `Size` takes an `InputFieldSize*` constant (default `md`). It scales the input, the affixes, and the label.

## FileUploadInput

A drop zone that also opens a file picker. It lists each file with a type icon, its size, and a remove button. The files submit with the form under `InputName`.

```go
<form hx-post="/api/v1/files/" hx-encoding="multipart/form-data" hx-swap="none">
    @uiForm.FileUploadInput(uiForm.FileUploadInputSettings{
        InputName: "archiveFiles",
        Label:     "Drop archives here or click to pick",

        // OptionalFields
        AcceptedTypes:    []string{".tar", ".zip", ".xz", ".br"},
        MaxFileCount:     3,
        MaxFileSizeBytes: 5242880,
    })
</form>
```

- Set `hx-encoding="multipart/form-data"` on an HTMX form, or `enctype="multipart/form-data"` on a native form.
- A visually hidden file input holds the listed files. The server reads them as one multipart field. The input stays focusable: Tab reaches the zone, and Enter or Space opens the picker.
- `AcceptedTypes` takes a slice in the `accept` syntax: extensions (`.zip`), MIME families (`image/*`), and MIME types. The component joins the slice for the `accept` attribute. A dropped file skips `accept`, so the component checks each file against the same list.
- `MaxFileSizeBytes` and `MaxFileCount` reject files over the limit. Zero means no limit. With `MaxFileCount: 1`, the picker selects one file, and a new file replaces the pending one.
- The component also rejects duplicates. Each rejected file shows its name and reason in the browser language (English, Portuguese, Spanish, or Chinese).
- The limits run only in the browser. The server must check them again.
- During an HTMX request, each file shows its share of the upload progress. A native POST shows no progress, because the browser reports none.
- After a successful request, the files show as uploaded and leave the input, so the next submit sends only new files. After a failed request, the files stay, so the next submit sends them again.
- The form `reset` event clears the list. `Description` shows a hint below the `Label`. `IsRequired` marks the label and sets `required`.
- `LabelOneWayStatePath` and `DescriptionOneWayStatePath` bind the visible label and hint to Alpine paths. `Label` and `Description` remain the static fallbacks.
- `Size` takes `FileUploadInputSizeXs` through `FileUploadInputSizeXl` and scales the zone, the icon, the rows, and the text together.
- `Zone`, `Row`, `Progress`, and `Rejected` take the colors and the radius of each part. Every color field takes a Tailwind token.
- `Zone.Icon` takes a Phosphor icon name for the drop zone (default `ph-tray-arrow-up`). The row icon follows the file type and is not configurable.

## MultiColumnRepeatableFieldset

A repeatable list of rows built from column definitions. Each cell renders through the matching input component: text and number cells through `InputField`, select cells through `SelectInput`, and checkbox cells through `CheckboxInput`. The user adds, duplicates, moves, and removes rows, or switches to a bulk textarea. The list writes itself to a state path as JSON.

```go
@uiForm.MultiColumnRepeatableFieldset(uiForm.MultiColumnRepeatableFieldsetSettings{
    TwoWayStatePath: "teamMembers",
    Columns: []uiForm.MultiColumnRepeatableFieldsetColumn{
        {Id: "name", Label: "Name", InputType: uiForm.MultiColumnRepeatableFieldsetColumnTypeText, IsRequired: true},
        {Id: "role", Label: "Role", InputType: uiForm.MultiColumnRepeatableFieldsetColumnTypeSelect,
            Options: []uiForm.SelectLabelValueOption{{Label: "Owner", Value: "owner"}}},
    },

    // OptionalFields
    AddRowLabel:     "Add member",
    Dropdown:        uiForm.MultiColumnRepeatableFieldsetDropdownStyle{BackgroundColor: "neutral-950"},
    EmptyStateLabel: "No members yet",
    InputName:       "teamMembers",
    Size:            uiForm.InputFieldSizeSm,
    TextCase:        "capitalize",
})
```

- `InputType` takes `MultiColumnRepeatableFieldsetColumnTypeText`, `MultiColumnRepeatableFieldsetColumnTypeNumber`, `MultiColumnRepeatableFieldsetColumnTypeSelect`, or `MultiColumnRepeatableFieldsetColumnTypeCheckbox`. An unknown type renders a text cell.
- Number columns take `Min`, `Max`, and `Step`. Select columns take `Options`; the picker is clearable and shows the label as its placeholder.
- A column's `LabelOneWayStatePath` binds the header, the text or select cell label, and the checkbox accessible name to Alpine state. `Label` remains the static fallback.
- `IsRequired` marks the cell required: the field carries the required and aria-invalid wiring and paints a red border while empty. It does not block the row actions.
- `Size` takes an `InputFieldSize*` constant and applies to every cell. `TextCase` changes the cell label case.
- Each cell type carries a style group. `Dropdown` takes `BackgroundColor`, `BorderColor`, and `BorderRadius` for the select option list. `Select` takes `SelectedIconColor` and a `Trigger` surface of `BackgroundColor`, `BorderColor`, `BorderRadius`, and `TextColor`. `Checkbox` takes `CheckedColor`, `UncheckedColor`, `FocusRingColor`, and `Shape`.
- `NewRowFunc` names a JavaScript expression that returns a new row object. The add row action merges the returned values into the column defaults, so each added row can carry distinct values. Invalid or unknown entries fall back to the column defaults and log `MultiColumnRepeatableFieldsetNewRowFuncFailed`. Give the expression from code, never from request data; the component evaluates it in the browser.
- `BulkEdit.IsEnabled` adds a button that switches the table to one textarea line per row and back. `Delimiter` defaults to `=` for two columns and `|` for three or more columns, and every textarea edit writes the bound array, so a pasted block submits with the form. Each line splits on the first `len(Columns)-1` unescaped delimiters, so the last column keeps any remaining delimiter text, and each part is trimmed. A backslash escapes the delimiter, a newline, or a literal backslash: type `\=` for a literal `=` when the delimiter is `=`, type `\|` for a literal `|` when the delimiter is `|`, type `\n` for a newline, and type `\\` for a literal backslash. A delimiter must be one character and must not be `n`, `\`, a carriage return, or a line feed; otherwise the default applies. Bulk edits update existing row objects: rows whose serialized text is unchanged keep their identity, and other rows keep properties outside the configured columns when the row count is unchanged. Blank lines are ignored, and a text-only row with every cell empty is skipped when the textarea opens. Checkbox cells read `true`, `1`, `yes`, or `on` as checked; select cells take the option value. `ButtonLabel`, `TableButtonLabel`, `ButtonIcon`, `TableButtonIcon`, `TextAreaLabel`, and `HintValue` replace the default texts and icons. The default button labels are lowercase.
- `AddRowButton` and `BulkEditButton` style the action buttons. Each takes `BackgroundColor`, `BackgroundColorHover`, `Shape`, `TextColor`, and `TextColorHover`; `AddRowButton` also takes `IconLeft` and `IconRight`. `AddRowLabel` replaces the add button text.
- The move-up and move-down buttons disable on the edge rows. `EmptyStateLabel` replaces the empty-state text.
- `InputName` adds a hidden input that submits the rows as JSON with the form.
- The rows read and write through `TwoWayStatePath`. Each edit writes the whole array back. An empty array shows the empty state.
- Give the state paths from code, never from request data. The component puts them into client-side expressions.

## NumberRangeInputField

A number box and a slider on one state path. The user types the value or drags it.

```go
@uiForm.NumberRangeInputField(uiForm.NumberRangeInputFieldSettings{
    InputName:       "cpuQuota",
    Label:           "CPU Quota",
    TwoWayStatePath: "cpuQuota",

    // OptionalFields
    Min:        "10",
    Max:        "400",
    Step:       "10",
    UnitSuffix: "%",
})
```

- `Min`, `Max`, and `Step` default to `"0"`, `"100"`, and `"1"`. They set the box limits and the slider track.
- `LabelOneWayStatePath` forwards to the number box label and the slider accessible name. `Label` remains the static fallback.
- The box keeps a draft while the user types, so a partial value does not reach the path.
- On `change` (blur or Enter), the field keeps the draft between `Min` and `Max`, moves it to the nearest step, and writes a number to the path. An empty or invalid draft returns to the path value.
- The slider writes the path directly, and the box follows.
- `UnitSuffix` shows as a right affix. `Size` takes an `InputFieldSize*` constant and scales the box and the slider.
- `SliderDisplay` takes `NumberRangeInputFieldSliderDisplayAlways` (the default), `NumberRangeInputFieldSliderDisplayResponsive`, or `NumberRangeInputFieldSliderDisplayHidden`. Responsive hides the slider while the field container is narrower than 24rem.
- `Thumb` and `Track` take the `uiControl.RangeSliderThumbSettings` and `uiControl.RangeSliderTrackSettings` values, so the slider takes every RangeSlider appearance field. The field sets the value, the bounds, the size, and the aria label itself.
- `IsValueLabelsHidden` hides the track end value labels, which show by default.
- The box submits under `InputName`.

## PasswordInput

A password field with a reveal button, an optional password generator, and an optional strength meter. It uses `InputField`, so every `InputField` setting applies.

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

- `Rules` configures the generator and the meter: `MinLength` (default 6), `MaxLength` (default 64), `GenerationLength` (default 16), and one `ShouldInclude*Chars` flag for each character class. With no class selected, all classes apply.
- The generator uses `UiToolset.CreateRandomPassword(options)` and includes at least one character of each selected class.
- While the field has focus and Caps Lock is on, a `Caps Lock is on` line shows below the field. The line clears on focusout. `CapsLockWarningText` replaces the text, `CapsLockWarningTextOneWayStatePath` binds it to Alpine state, `CapsLockWarningTextColor` takes a color token, and `CapsLockWarningTextCase` takes a `uiToolset.TextCase*` value. `ShouldHideCapsLockWarning` removes the warning and its key listeners.
- The generate button copies the password and shows a toast in the browser language, with English as the fallback. The page must mount `Toast`.
- The meter shows a bar and a checklist of the enabled criteria. Both update as the user types.
- The reveal and generate buttons are `uiControl.Button` squares sized to the field, each with a tooltip and an accessible name. `ActionButtonStyle` takes `PasswordInputActionButtonStyleBoxed` (the default) or `PasswordInputActionButtonStylePlain`.
- `LabelOneWayStatePath` forwards to the `InputField` label, placeholder, and accessible-name path.
- `StrengthMeterColor`, `StrengthCriteriaTextColor`, `StrengthCriteriaFulfilledIconColor`, and `StrengthCriteriaUnfulfilledIconColor` take color tokens.

## InlineCheckboxGroup

A row or column of checkboxes under one label. The label sits in the top border, like the other notched fields.

```go
@uiForm.InlineCheckboxGroup(uiForm.InlineCheckboxGroupSettings{
    Label: "Fruits",
    InputSettings: []uiForm.CheckboxInputSettings{
        {Label: "Apple", Value: "apple", TwoWayStatePath: "selectedFruits", InputName: "selectedFruits"},
        {Label: "Banana", Value: "banana", TwoWayStatePath: "selectedFruits", InputName: "selectedFruits"},
    },
})
```

- Each `CheckboxInputSettings` keeps its own size, shape, colors, and errors.
- All inputs must use the same array `TwoWayStatePath`. Alpine adds and removes the checked values.
- `Orientation` takes `InlineCheckboxGroupOrientationHorizontal` (the default) or `InlineCheckboxGroupOrientationVertical`.
- `MaxVisibleOptions` limits the visible rows and scrolls the rest.
- `TextCase` changes only the group label. `LabelOneWayStatePath` binds the group label to Alpine state. `Label` remains the static fallback.
- On hover or focus, the group shows select-all and deselect-all `FloatingButton` actions. Select-all skips disabled options. `FloatingButtonStyle` restyles them. `ShouldHideFloatingButtons` removes them. Vertical rows reserve the resolved action width, so a large action never covers the labels.

## FloatingButton

A small icon button over a field. It shows on hover or focus. `InlineCheckboxGroup` and `TextArea` use it.

```go
@uiForm.FloatingButton(uiForm.FloatingButtonSettings{
    AriaLabel:   "Clear text area content",
    Icon:        "ph-x-circle",
    OnClickFunc: "description = ''",
})
```

- `Icon` sets a static icon. `IconOneWayStatePath` binds the icon to Alpine state.
- `Style` takes a `FloatingButtonStyle` with background, ring, and text colors, plus `Size` and `Shape`. `Size` takes a `uiControl.ButtonSize*` constant (default `xs`). An empty field keeps the neutral default.
- The hosting field reserves the resolved action width: `InlineCheckboxGroup` and `TextArea` read the same button-size width contract, so a `ButtonSizeLg` or `ButtonSizeXl` action never covers the field content.

## TextArea

A multiline input with expand, copy, and clear actions.

```go
@uiForm.TextArea(uiForm.TextAreaSettings{
    Label:           "Description",
    TwoWayStatePath: "description",
})
```

- `IsCode` uses a monospace font. `IsReadOnly` sets `readonly` and drops the clear action, so a read-only field never mutates its bound state. `IsRequired` marks the label and sets `required`.
- `LabelOneWayStatePath` renders the label, the placeholder, and the accessible name from an Alpine path. `Label` remains the static fallback.
- On hover or focus, the field shows the expand, copy, and clear `FloatingButton` actions. A read-only field shows expand and copy only. `FloatingButtonStyle` restyles them. `ShouldHideFloatingButtons` removes them, and the text uses the full width. The text reserves the resolved action width, so a large action never covers the content.

## SelectInput

A dropdown select. It is not a native `<select>`: hidden radio inputs and an Alpine list hold the value. It has four option modes:

1. `FlatOptions []string` for plain values.
2. `LabelValueOptions []SelectLabelValueOption` for label and value pairs.
3. `LabelValueOptions` with `LabelHtml` for rich option content.
4. `OptionsOneWayStatePath` for an Alpine array built at render time. Each entry is a plain value or an object with `value` and `label`.

```go
@uiForm.SelectInput(uiForm.SelectInputSettings{
    InputName:       "status",
    Label:           "Status",
    FlatOptions:     []string{"running", "stopped"},
    TwoWayStatePath: "status",
})
```

- `LabelOneWayStatePath` renders the legend and the empty trigger text from an Alpine path. `Label` remains the static fallback.
- `IsMultiSelect` shows checkbox rows, keeps the list open after each toggle, and binds an array. The trigger joins the selected labels. `InputId` prefixes the checkbox ids.
- `OptionsOneWayStatePath` replaces the static option lists. The trigger resolves the selected label from the same array, in single and multi mode. A zero-valued option renders as `0`; in multi mode, a selected value with no matching option shows its string form.
- `ShouldIncludeBlankOption` adds a blank entry in single mode. In multi mode, a clear control shows while the array holds a value.
- `IsDisabledOneWayStatePath` disables the trigger while the path is truthy. `IsInvalidOneWayStatePath` paints the field border red and sets `aria-invalid` on the trigger while the path is truthy.
- `OnChangeFunc` runs after a change, a toggle, or a clear.
- `Size` takes a `SelectInputSize*` constant (default `md`).
- The trigger opens on click, Enter, Space, or ArrowDown; the checked row takes focus, Enter commits, and Escape closes with focus return. A pick in single mode closes the dropdown.
- `DropdownMinHeightClass` and `DropdownMaxHeightClass` limit the list height (default `max-h-60`). `MaxVisibleOptions` sets a row count and replaces both.
- `Trigger` takes a `SelectInputTriggerStyle` with `BackgroundColor`, `BorderColor`, `BorderRadius`, and `TextColor` for the field surface. `SelectedIconColor` adds a checked-row icon and its focus ring. `DropdownBackgroundColor`, `DropdownBorderColor`, and `DropdownBorderRadius` shape the list surface.

## DualSelectInput

One control that combines a select for the unit or operator with a text or number input for the value. Both halves share one label and one border, for entries like "port 80" or "equals 5".

```go
@uiForm.DualSelectInput(uiForm.DualSelectInputSettings{
    Label:   "Port",
    Options: portOptions,

    // OptionalFields
    SelectTwoWayStatePath: "portUnit",
    ValueTwoWayStatePath:  "portValue",
    InputName:             "portValue",
    SelectName:            "portUnit",
})
```

- With `SelectTwoWayStatePath` empty, both halves bind `ValueTwoWayStatePath`: a preset pick writes the slot, and a typed value that matches no option shows `Custom` on the select.
- `InputType` takes a `uiForm.InputType*` constant (default `text`). `SelectName` and `InputName` submit the two halves as separate form fields.
- `LabelOneWayStatePath` renders the legend, the empty trigger text, the placeholder, and the accessible name from an Alpine path. `Label` remains the static fallback.
- The select half is not a native `<select>`: it opens the shared custom dropdown. The trigger opens on click, Enter, Space, or ArrowDown; the checked row takes focus, Enter commits, and Escape closes with focus return.
- The value half carries a darker background than the select half, so the two halves read as separate controls.
- `IsInvalidOneWayStatePath` marks both halves and the border while the path is truthy.
- `Size` takes an `InputFieldSize*` constant. `TextCase` changes the label.
- `Field` takes a `DualSelectInputFieldStyle` with `BackgroundColor`, `BorderColor`, `BorderRadius`, `DividerColor` (the line between the halves), and `TextColor`. An empty field keeps the neutral default.
- `Dropdown` takes a `DualSelectInputDropdownStyle` with `BackgroundColor`, `BorderColor`, and `BorderRadius` for the option list surface.

## SearchableSelectInput

A dropdown with a filter box. Use it for more than about ten options. The list scrolls, and a clear button empties the selection.

```go
@uiForm.SearchableSelectInput(uiForm.SearchableSelectInputSettings{
    InputName:       "country",
    Label:           "Country",
    FlatOptions:     []string{"Argentina", "Brazil", "Chile"},
    TwoWayStatePath: "country",
})
```

- It takes the `FlatOptions` and `LabelValueOptions` modes of `SelectInput`. The filter matches the label and ignores case.
- `LabelOneWayStatePath` renders the legend, the empty placeholder, and the accessible name from an Alpine path. `Label` remains the static fallback.
- `IsMultiSelect` binds an array and keeps the list open after each toggle.
- `SelectionDisplay` takes `SearchableSelectInputSelectionDisplayText` (the default, which joins the labels) or `SearchableSelectInputSelectionDisplayTags`.
- Tags apply only in multi mode. Each tag shows the option label, and the array holds the value. Backspace in an empty input removes the last tag. `ShouldAllowCustomValues` turns typed text into a tag on Enter.
- `TagOuterBackgroundColor` (default `secondary-500`), `TagOuterTextColor` (default `neutral-50`), and `TagInnerBackgroundColor` (default `neutral-50/10`) color the tags.
- The list height settings work as in `SelectInput`.
- `OnChangeFunc` runs on a selection, a toggle, a tag change, and a clear.
- Single mode writes one hidden input. Multi mode writes one hidden input for each value, so the form submits every value under `InputName`.

## RemoteSearchableSelectInput

A searchable select that loads its options from a URL while the user types. It shows loading, empty, and error states.

```go
@uiForm.RemoteSearchableSelectInput(uiForm.RemoteSearchableSelectInputSettings{
    InputName:       "country",
    Label:           "Country",
    OptionsUrl:      "/api/countries/search",
    TwoWayStatePath: "country",
})
```

- `LabelOneWayStatePath` renders the legend, the empty placeholder, and the accessible name from an Alpine path. `Label` remains the static fallback.
- The endpoint receives `GET OptionsUrl?<query param>=<typed text>`. It returns `{ "body": [{ "label": "...", "value": "..." }] }`, or a bare array of strings or `{label, value}` objects.
- `OptionsQueryParam` names the parameter (default `q`). `DebounceMs` sets the wait before each request (default 300). The component ignores a response that a newer request replaced.
- `MinQueryLength` sets the shortest query that starts a request (default 3). A shorter query shows a prompt.
- The list height settings work as in `SelectInput`.
- `IsMultiSelect` binds an array and keeps the selected labels while the query changes.
- `InitialOptions` gives labels for values that the state already holds. The list shows them until the user types. A fetched option replaces an initial option with the same value.
- `OnChangeFunc` runs on a selection, a toggle, and a clear.

## RadioInput

One radio option. Inputs on the same `TwoWayStatePath` form a group.

```go
@uiForm.RadioInput(uiForm.RadioInputSettings{
    Label:           "Option 1",
    StateValue:      "option1",
    TwoWayStatePath: "selectedOption",
})
```

- `Size` takes a `RadioInputSize*` constant. `TextCase` changes the label. `LabelOneWayStatePath` binds the visible label to Alpine state. `Label` remains the static fallback.
- `Style` takes a `RadioInputStyle` with `TextColor` for the idle label (default `neutral-100`), `ActiveColor` for the checked and hover color (default `secondary-500`), and `FocusRingColor` for the keyboard focus ring (default: `ActiveColor`).

## InlineRadioGroup

A row of radios under one label. The label sits in the top border, like the other notched fields.

```go
@uiForm.InlineRadioGroup(uiForm.InlineRadioGroupSettings{
    Label: "Select an option",
    InputSettings: []uiForm.RadioInputSettings{
        {Label: "Option 1", StateValue: "option1", TwoWayStatePath: "groupSelection"},
        {Label: "Option 2", StateValue: "option2", TwoWayStatePath: "groupSelection"},
    },
})
```

- `Orientation` takes `InlineRadioGroupOrientationHorizontal` (the default) or `InlineRadioGroupOrientationVertical`.
- `MaxVisibleOptions` limits the visible rows and scrolls the rest.
- The group `TextCase` changes only the group label. `LabelOneWayStatePath` binds the group label to Alpine state. Each option takes its own `RadioInputSettings.TextCase`.

## ToggleSwitch

An on-off switch. It binds to a boolean. With `CustomValue`, it adds one value to an array instead.

```go
@uiForm.ToggleSwitch(uiForm.ToggleSwitchSettings{
    Label:           "Enable notifications",
    TwoWayStatePath: "isNotificationsEnabled",
})
```

`LabelPosition` takes `ToggleSwitchLabelPositionLeft` or `ToggleSwitchLabelPositionRight`. `LabelOneWayStatePath` binds the label to Alpine state. `Label` remains the static fallback. The `*Color` settings change the track and thumb colors.

## InputHint

The hint helper that the fields above share. It has no field of its own.

- `HintDisplay` takes `InputHintDisplayTooltip` (an icon with a tooltip) or `InputHintDisplayDescription` (a line below the field).
- `HintValue` sets static text. `HintStatePath` binds live text.
- The `HintDisplayTooltip*WidthClass` and `HintDisplayTooltip*HeightClass` fields limit the tooltip size. `HintDisplayTooltipBackgroundColor` sets its color.
- The tooltip renders on the document body, so a scroll container or a modal cannot clip it.
- `InputHintTooltip` and `InputHintDescription` can render on their own, next to custom markup.

## Non-obvious behaviors

- `SelectInput` renders a `templ.JSONScript` block for label-value options. The block feeds the label lookup and the multi-select summary.
- `SearchableSelectInput` and `RemoteSearchableSelectInput` share the Alpine state `searchableSelectInput`. Each renders one `templ.JSONScript` config with the options, the flags, and the remote settings.
- The floating legend collapses while the field is empty, so the label shows as a placeholder.
- `InputField` hides its empty legend with `display: none`. Chrome keeps scroll space for a zero-size legend, so the opacity collapse of the other fields can cause a false scroll.
- A dropdown opens upward when the space below cannot hold it. The space ends at the viewport and at each clipping ancestor. When neither side fits, the dropdown uses the larger side. `UiToolset.SelectDropdown.openUpwardResolver` makes this decision for `SelectInput` and the searchable fields.
- `UiToolset.CreateRandomPassword` takes an options object: `length`, `minLength`, `maxLength`, and one `include*` flag for each character class. Without arguments, it returns 16 characters.
- `InputName` sets the key in a form submission. `InputId` is optional.

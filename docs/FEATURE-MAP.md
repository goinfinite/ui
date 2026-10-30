# Feature Map

> Auto-maintained index of every user-facing feature and the code path that implements it. Updated alongside the code — not after the fact.

## Text Input Field

Single-line text input with configurable type (text, email, number, date, password, etc.), with support for labels, a runtime label state path, label text case, hints, required indicators, a size scale, optional prefix/suffix affixes, an arbitrary right-edge component slot, and a runtime type expression.

**Flow:**

1. `src/form/inputField.templ` — Component definition with InputFieldSettings struct exposing InputType, Label, LabelOneWayStatePath, TextCase, TwoWayStatePath, Value, Size, optional affixes, affix width percentages, AffixRightComponent for an arbitrary right-edge component, and InputTypeExpression for a runtime type; LabelOneWayStatePath renders the legend, the placeholder, and the accessible name with Alpine bindings, while Label remains the static fallback
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/inputField_templ.go` — Compiled templ output rendering HTML input with Alpine.js binding and Tailwind styling

---

## Password Input

Password field with a reveal toggle, an optional random password generator, an optional strength meter, and a Caps Lock hint. The meter shows a percentage bar plus a criteria checklist that updates as the user types. The hint shows below the field while it has focus and Caps Lock is on.

**Flow:**

1. `src/form/passwordInput.templ` — PasswordInputSettings, the rules resolver (length range, generation length, character classes), the criteria-item resolver, the reveal and generate action buttons rendered as uiControl.Button with padding and min-width/min-height constraints and a boxed or plain style (ActionButtonStyle), the Caps Lock hint with its key and focus handlers, custom text, one-way text path, text color, text case, and hide setting, the meter with its color settings, the LabelOneWayStatePath passthrough, and the InputField composition
2. `src/form/passwordInputState.js` — Alpine data component: reveal flag, Caps Lock flag and its key-event and focusout updates, strength percentage, per-criterion flags, the generator call, and the clipboard copy with a language-matched toast that ends with an exclamation mark
3. `src/import/toolset/randomPassword.js` — `createRandomPassword(options)` with configurable length and character classes; guarantees one character from each selected class
4. `src/form/passwordInput_test.go` — Table-driven tests for the rules and criteria resolvers plus the rendered actions and meter
5. `src/form/passwordInput_templ.go` — Compiled output

---

## Multi-line Text Area

Multi-line text input with five height steps (h-12/24/36/48/60), expand-to-3x toggle, floating action icons (expand, copy, clear) anchored to the text line, and optional hints rendered as a focusable tooltip or a description line. `FloatingButtonStyle` repaints the action icons and changes their size and shape; `ShouldHideFloatingButtons` removes them and drops the reserved right padding. The text reserves the resolved action width, so a large action never covers the content. A read-only field shows expand and copy only, so it never mutates its bound state, and its negative top margin resets when it is the first child, so a scroll container never clips the legend. LabelOneWayStatePath binds the legend, the placeholder, and the accessible name to Alpine state.

**Flow:**

1. `src/form/textArea.templ` — Component definition with TextAreaSettings struct and the LabelOneWayStatePath legend, placeholder, and accessible-name binding; heights and icon positions supplied mutually exclusively via Alpine `:class`
2. `src/form/floatingButton.templ` — Shared floating icon button that resolves the style defaults and renders a Button
3. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
4. `src/form/textArea_templ.go` — Compiled output rendering textarea element with styling and Alpine.js integration

---

## Select Dropdown

Dropdown select component with support for flat string options, label-value pairs, or options from an Alpine array, a blank option, and a multi-select mode. Single mode renders native radio rows; multi mode renders the shared styled checkbox rows, keeps the dropdown open after each toggle, and joins the selected labels in the trigger. The trigger opens on click, Enter, Space, or ArrowDown; the checked row takes focus, Enter commits, and Escape or a pick closes with focus return.

**Flow:**

1. `src/form/selectInput.templ` — Component definition with SelectInputSettings struct (including IsMultiSelect, InputId, LabelOneWayStatePath, OptionsOneWayStatePath, the Trigger style group, SelectedIconColor, and the dropdown border and radius), SelectLabelValueOption data structure, the option checkbox and option icon renderers, the dynamic x-for option rows, and the dropdown height and MaxVisibleOptions settings; LabelOneWayStatePath binds the legend and empty trigger text, OptionsOneWayStatePath replaces the static option lists and reads plain values or value/label objects from Alpine state, the selected-display presence and empty-state checks keep a zero value visible, dynamic single and multi select keep the option's original value type, the loop aliases derive from the bound state path and OnChangeFunc, and each single-select radio handles its own change assignment so the event bubbles
2. `src/form/selectInputState.js` — Alpine.js data component for the select dropdown: the shared dropdown controller (open toggle, checked-row focus, close with focus return, and commit), the multi-select summary formatter, the typed multi-select toggle that stores original option values, the strict-identity option matcher, and the dynamic option label resolvers for single and multi mode, including the string fallback for unmatched multi-select values
3. `src/import/toolset/selectDropdown.js` — Shared open-upward decision that measures the space above and below the trigger, bounded by the viewport and every clipping ancestor, so the list never opens downward into an overflow-hidden container; also the shared dropdown controller that drives the open, the focus, the close with focus return, and the commit for every custom dropdown
4. `src/form/selectDropdown.go` — Shared dropdown surface classes with the optional border color and radius, the dropdown height classes resolver, the MaxVisibleOptions height style with the row height of each dropdown family, and the label-value option values resolver
5. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
6. `src/form/selectInput_templ.go` — Compiled output rendering select with native radio or checkbox options, an embedded JSON script for label lookup and the multi-select summary, and Alpine.js state management

Supports optional hint text rendered either as a focusable info-icon tooltip inside the dropdown row or as a description line below the fieldset.

---

## Dual Select Input

One control that combines a select for the unit or operator with a text or number input for the value, sharing one label and one border. The select half opens the shared custom dropdown, and the value half carries a darker background so the halves read as separate controls. A preset pick writes the shared path, and a typed value that matches no option shows a Custom row in the list and Custom on the trigger.

**Flow:**

1. `src/form/dualSelectInput.templ` — Component definition with DualSelectInputSettings struct, the LabelOneWayStatePath legend, empty trigger, placeholder, and accessible-name binding, the DualSelectInputFieldStyle group (background, border, radius, divider, and text colors), and the DualSelectInputDropdownStyle group (background, border, and radius); the custom dropdown trigger, the option rows, the Custom row for unmatched values, and the value input with the darker background and the field right radius; the value input reaches the field's top border and bottom edge through the shared legend-band size token, and the fieldset stretches to fill a taller grid row; the select and the input bind the same ValueTwoWayStatePath unless SelectTwoWayStatePath splits them; the invalid path marks both halves and the border
2. `src/form/dualSelectInputState.js` — Alpine data component: the options JSON parse, the unmatched-value check, the selected-label lookup with the Custom fallback, and the shared dropdown controller
3. `src/form/inputField.templ` — Shared size tokens for the label, the padding, and the text size
4. `src/form/dualSelectInput_test.go` — Tests for the rendered field style, the rendered dropdown style, the rendered value input legend band, the label state path, and the rendered shared and split bindings
5. `src/form/dualSelectInput_templ.go` — Compiled output

---

## Searchable Select

Dropdown with a filter box that narrows local options in place, a clear button, a No matches row, and a max-height scroll list. Single mode binds one value; multi mode binds an array and shows the joined labels or removable tags. With tags, a typed value becomes a tag when custom values are enabled; Enter adds, Backspace removes the last tag, and each tag carries a remove button.

**Flow:**

1. `src/form/searchableSelectInput.templ` — SearchableSelectInputSettings with the LabelOneWayStatePath legend, empty placeholder, and accessible-name binding, and the shared shell, trigger, filter input, action icons, selected labels, selected tags with tag color settings, local option rows, and local dropdown
2. `src/form/searchableSelectInput.go` — The client config types and resolvers (normalized items, multi-select and custom-value flags, remote settings, the selection presentation), plus the selection, toggle, clear, custom-value, tag-remove, tag-backspace, and selected-state expression builders
3. `src/form/searchableSelectInputState.js` — Shared Alpine data component: filter, visibility, No matches, selected-label lookup, and the remote fetch lifecycle
4. `src/import/toolset/selectDropdown.js` — Shared open-upward decision that measures the space above and below the trigger, bounded by the viewport and every clipping ancestor
5. `src/form/selectDropdown.go` — Shared dropdown surface classes with the optional border color and radius, and the open-list height settings
6. `src/form/searchableSelectInput_test.go` — Tests for the item and presentation resolvers, the expression builders, and the rendered dropdown
7. `src/form/searchableSelectInput_templ.go` — Compiled output

Single mode writes one hidden input; multi mode writes one hidden input per selected value under `InputName`.

---

## Remote Searchable Select

Searchable select whose options come from a URL as the user types, with debounce, minimum query length, stale-response discard, a loading indicator, and error and empty states.

**Flow:**

1. `src/form/remoteSearchableSelectInput.templ` — RemoteSearchableSelectInputSettings with the LabelOneWayStatePath forwarding, the remote config resolver, the dropdown height and MaxVisibleOptions settings, the InitialOptions that label bound values until the user types, and the remote option rows and dropdown with the minimum-length prompt, sharing the searchable shell, trigger, and filter input
2. `src/form/searchableSelectInputState.js` — Shared Alpine data component owning the remote fetch lifecycle: request id for stale discard, loading and error flags, and the debounced request
3. `src/form/remoteSearchableSelectInput_test.go` — Tests for the remote config resolver
4. `src/form/remoteSearchableSelectInput_templ.go` — Compiled output

The endpoint receives `GET OptionsUrl?<query param>=<typed text>` and returns `{ "body": [{ "label": "...", "value": "..." }] }`; a bare array of strings or label-value objects also works.

---

## Checkbox Input

Checkbox that selects one value, bound to a boolean or an array Alpine.js state, with configurable shapes, sizes, colors, label position, and disabled state.

**Flow:**

1. `src/form/checkboxInput.templ` — Component definition with CheckboxInputSettings for binding, shape, size, colors, indeterminate, disabled, label position, the LabelOneWayStatePath label binding, and the AriaLabelOneWayStatePath accessible-name binding
2. `src/form/checkboxInput_templ.go` — Compiled output rendering the native input, the styled box, and the check or dash icon

The check renders through color inheritance: the box carries the checked color and the icon inherits it. The indeterminate state is set on the native input through `x-effect` and renders as a dash.

---

## Cron Schedule

Combined control for a cron expression. A mode switch picks between a predefined-schedule dropdown and a five-field custom builder (minute, hour, day, month, weekday). Both write one cron string to the bound path. Each field validates its cron syntax with the field marked red while invalid, an invalid schedule clears the bound path, and a preview reads the resulting schedule in plain words.

**Flow:**

1. `src/form/cronSchedule.templ` — Component definition with CronScheduleSettings struct and the CronScheduleFieldStyle (outer background, border, and radius), CronScheduleDropdownStyle (both option list surfaces), CronSchedulePreviewStyle (plain-words text color, icon color, and expression chip colors), and RadioInputStyle (mode radio colors) groups; the mode radios, the predefined SelectInput, the five DualSelectInput fields, the preview row that follows Size and TextCase, and the hidden form input
2. `src/form/cronSchedule.go` — The predefined option list with the macro entries and the corrected Every 5 days value, and the field descriptors with the per-field invalid expressions; the option values reach the Alpine state through the shared selectDropdown.go resolver
3. `src/form/cronScheduleState.js` — Alpine data component: the mode and field state, the cron field validator for stars, steps, ranges, and lists with per-field bounds, the custom expression composer, the human-readable preview builder, and the init that seeds the control from the bound path
4. `src/form/cronSchedule_test.go` — Tests for the option list including the Every 5 days value, the field descriptors, the rendered bound state, the style groups, and the size and text case on the preview
5. `src/form/cronSchedule_templ.go` — Compiled output

---

## Toggle Switch

Switch component that binds either a boolean Alpine.js state or a custom value in an Alpine.js array, with configurable sizes, colors, required state, and disabled state.

**Flow:**

1. `src/form/toggleSwitch.templ` — Component definition with ToggleSwitchSettings for binding, sizing, colors, disabled state, label position, and the LabelOneWayStatePath label binding
2. `src/form/toggleSwitch_templ.go` — Compiled output rendering the hidden form field, checkbox input, and switch track

With `CustomValue`, Alpine.js adds or removes the value from an array. Without it, Alpine.js treats the checkbox as a boolean state.

---

## Radio Button Input

Single radio button component for choice selection within a group, with label and state binding support. The RadioInputStyle group sets the idle text color, the checked and hover color, and the keyboard focus ring.

**Flow:**

1. `src/form/radioInput.templ` — Component definition with RadioInputSettings struct, the RadioInputStyle group, and the LabelOneWayStatePath label binding
2. `src/form/radioInput_test.go` — Tests for the rendered default style, the rendered custom style, and the label state path
3. `src/form/radioInput_templ.go` — Compiled output rendering input[type="radio"] with Tailwind styling

---

## Inline Radio Group

Radio button group for presenting multiple mutually exclusive options in a row or a column under one notched label.

**Flow:**

1. `src/form/inlineRadioGroup.templ` — Component definition with InlineRadioGroupSettings struct, the orientation resolver, and the LabelOneWayStatePath group-label binding
2. `src/form/inlineRadioGroup_templ.go` — Compiled output rendering multiple radio inputs horizontally or vertically

---

## Inline Checkbox Group

Checkbox group for presenting several choices in a row or a column under one notched label, bound to one array state path. The group reveals select-all and deselect-all floating actions on hover or focus; `FloatingButtonStyle` repaints them and changes their size and shape, and `ShouldHideFloatingButtons` removes them. Vertical rows reserve the resolved action width, so a large action never covers the labels.

**Flow:**

1. `src/form/inlineCheckboxGroup.templ` — Component definition with InlineCheckboxGroupSettings struct, the orientation resolver, and the LabelOneWayStatePath group-label binding, composing CheckboxInput per setting and FloatingButton for the select-all actions
2. `src/form/floatingButton.templ` — Shared floating icon button that resolves the style defaults and renders a Button
3. `src/form/inlineCheckboxGroup_templ.go` — Compiled output rendering the checkboxes horizontally or vertically

---

## Floating Button

Small icon button that floats over a field and reveals on hover or focus. `InlineCheckboxGroup` and `TextArea` render their actions with it. `FloatingButtonStyle` overrides the background, ring, and text colors, the size, and the shape; every field falls back to the neutral default.

**Flow:**

1. `src/form/floatingButton.templ` — Component definition with FloatingButtonSettings and FloatingButtonStyle structs; resolves each style default and renders a Button
2. `src/form/floatingButton.go` — The size resolver that defaults to the extra-small Button, and the action width resolver that asks the Button contract for the icon-only width the fields reserve
3. `src/form/floatingButton_templ.go` — Compiled output rendering the Button with the resolved style

---

## Button

Interactive button component with customizable label, icons (left/right/top using Phosphor), size/shape variants, colors, disabled state, click handlers, an optional hover tooltip, action states, and an optional countdown. Every button renders a `type` attribute: `button` by default, `submit` when `IsSubmit` is set. `FormId` adds the HTML form attribute, so a submit button outside a form can submit it. With `ShouldTrackAction`, the button follows the action it triggers: the promise its OnClickFunc returns, or, with `IsSubmit`, the native POST or HTMX request of its form. While the action runs, it shows an icon and the pending label and disables itself; when the action ends, it shows a succeeded or failed result. `PendingIcon`, `SucceededIcon`, and `FailedIcon` replace the default action icons, and `PendingIconAnimation` takes a `uiToolset.Animation*` constant for the pending icon animation. `CountdownSeconds` keeps the button disabled and shows the remaining seconds in the label; `CountdownTriggerOneWayStatePath` restarts the countdown when an Alpine value changes.

**Flow:**

1. `src/control/button.templ` — Component definition with ButtonSettings struct supporting OnClickFunc handlers, icon binding (including IconTop and its state path), FormId, a padding override with min/max width and height constraints, a static or Alpine-bound label and aria label, the type attribute, the optional tooltip, the countdown settings, the action labels with their one-way state paths, and per-state icons; `Button` wraps `buttonFrame` in one button state scope when an action display or countdown is present, and a tracked non-submit click passes its result to `trackActionPromise`
2. `src/control/buttonSize.go` — The size token resolver (padding, text size, gap, icon size, and the padding and icon widths in rem) and the exported icon-only width contract that the floating action fields reserve
3. `src/control/buttonSize_test.go` — Tests for the width contract and a rendered-class check that the contract matches the Button output
4. `src/control/buttonState.js` — Alpine data component that tracks the click promise or the owning form submit and HTMX request events, resolves a form named by the button FormId, including a form mounted after init, guards double submission, chooses the idle, pending, succeeded, or failed label expression, resolves the icon for the same state from the map the component passes in, holds the countdown remaining seconds, starts and stops the interval, exposes isCounting, and restarts when the trigger value changes
5. `src/toolset/tooltipState.js` — Shared Alpine data component that teleports a tooltip to the document body, positions it with fixed coordinates on hover or focus, and flips it to stay inside the viewport
6. `src/toolset/tooltip.go` — The once-handle that embeds the shared tooltip state, plus TooltipSurfaceClassesResolver for the shared fixed-layer tooltip classes and the shared TooltipPosition* placement values
7. `src/control/button_templ.go` — Compiled output rendering the button with Alpine.js event binding and the teleported tooltip

---

## Number Range Input Field

Number box and slider bound to one state path. The box edits a local draft and commits on change: it clamps the draft to min/max, snaps it to the step, and writes a number to the path. The slider writes the path directly. The slider renders always, hides below a 24rem container, or never renders. The slider takes its appearance from the uiControl RangeSlider thumb and track settings.

**Flow:**

1. `src/form/numberRangeInputField.templ` — Settings, bounds defaults, the box and slider sub-templates, the slider display layout, the LabelOneWayStatePath forwarding to the number box and slider accessible name, and the thumb and track appearance passthrough
2. `src/form/numberRangeInputFieldState.js` — Alpine draft sync and the clamp-and-snap commit
3. `src/form/inputField.templ` — Renders the number box, the unit affix, and the floating label
4. `src/control/rangeSlider.templ` — Renders the slider and clamps the path on writes

---

## File Upload Input

Drop zone with click-to-pick that lists each file with a type icon, size, status, and remove button, rejects files that break the size, count, type, or duplicate rules, and submits the files with the form under InputName. During an HTMX request from the owning form, each file shows its share of the upload progress, then an uploaded or failed state.

**Flow:**

1. `src/form/fileUploadInput.templ` — Settings, the upload limits config, the size scale, the drop zone over a visually hidden, focusable file input, the LabelOneWayStatePath and DescriptionOneWayStatePath bindings, the file list with a remove Button, the rejection list, and the Zone, Row, Progress, and Rejected style groups
2. `src/form/fileUploadInputState.js` — Limit checks, the DataTransfer sync of the non-uploaded files into the input, the per-file progress split, the result marking, the reset clear, and the language-matched messages
3. `demo/fileUploadDemoRouting.js` — Demo-only XHR patch that replays upload progress without network traffic

---

## Date Time Picker

Text input that opens a calendar popover for date, time, or datetime values. LabelOneWayStatePath binds the visible label, the single-input accessible name, and the range Start and End accessible names to Alpine state. The value moves through the bound state path as an ISO 8601 string. An optional range mode splits the field into captioned Start and End sides with a divider, each editable and committed on its own, and paints the calendar span as a continuous band. Min and max bound the selectable days, locale and first-day settings drive the month and weekday names, and a typed ISO value commits from the keyboard. The centered time row renders through SelectInput, with optional timezone selection that sits beside the time when the field is wide enough and writes the UTC offset into the value. Optional hints render as a focusable tooltip or a description line. DateTimePickerStyle paints the field, the calendar surface, the range captions and divider, and every calendar day state; TimeSelector and DoneButtonStyle paint the time controls and the Done action; TextCase reaches the labels, inputs, time selects, and Done action; and the calendar navigation, the toggle, and the Done action render through the library Button.

**Flow:**

1. `src/form/dateTimePicker.templ` — Component definition with DateTimePickerSettings struct, the LabelOneWayStatePath binding, the DateTimePickerStyle, DateTimePickerTimeSelectorStyle, and DateTimePickerDoneButtonStyle groups, the calendar and SelectInput time popover markup, the timezone select, the captioned range sides, the hint rendering through InputHintTooltip and InputHintDescription, and the data and writeback expression builders that pass the bound values into the Alpine factory
2. `src/form/dateTimePickerState.js` — Alpine data component: the view navigation, the draft and per-side typed-value sync, the side activation, the range order normalization, the bounds checks, the locale formatting, the typed commit and revert, the timezone offset and wall-clock resolution, the banded day classes, and the keyboard day focus
3. `src/form/dateTimePickerOptions.go` — The hour, minute, and timezone option lists
4. `src/form/dateTimePicker_test.go` — Renders, path wiring, mode and locale settings, bounds, range captions and divider, time selector rendering, timezone options, style config, and Done button customization tests
5. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
6. `src/form/inputActionButton.go` — Shared boxed and plain action-button recipe behind the calendar toggle and the PasswordInput actions
7. `src/control/button.templ` — The calendar navigation, toggle, and Done buttons
8. `src/form/dateTimePicker_templ.go` — Compiled output

---

## Multi Column Repeatable Fieldset

Repeatable table of rows built from column definitions. Each cell renders through the matching input component: text and number cells through InputField, select cells through SelectInput, and checkbox cells through CheckboxInput, with a shared Size and TextCase. A column's LabelOneWayStatePath binds the header, the text or select cell label, and the checkbox accessible name to Alpine state. Each cell type carries its own style group: Checkbox (checked and unchecked colors, focus ring, shape), Select (trigger surface and checked-option icon color), and Dropdown (the option list surface). AddRowButton and BulkEditButton style the action buttons with color, shape, and icon, and the default button labels are lowercase. `NewRowFunc` names the expression that builds a new row, so the add row action can seed distinct values. `BulkEdit` adds a toggle that swaps the table for one textarea line per row, with a delimiter that defaults to `=` for two columns and `|` for three or more, backslash escaping for delimiters, newlines, and literal backslashes, single-character delimiter validation, text-matched row updates that keep properties outside the columns, leave unchanged cell values alone, and keep row order when an edit duplicates a line, text-only empty rows skipped, and each edit parsed back into the bound array. The user adds, duplicates, moves, and removes rows. The whole array reads and writes through the bound state path and submits as JSON through a hidden input.

**Flow:**

1. `src/form/multiColumnRepeatableFieldset.templ` — Component definition with the settings and column structs, the LabelOneWayStatePath header and cell binding, the cell and button style groups, the table markup, the cell sub-template mapped to InputField, SelectInput, or CheckboxInput, the row action buttons with the edge-row disabling, the empty-state row, the optional BulkEdit button and TextArea with lowercase default labels, and the hidden JSON input
2. `src/form/multiColumnRepeatableFieldsetState.js` — Alpine data component: the default row builder with the optional NewRowFunc evaluation, the append, duplicate, move, and remove operations, the bulk-edit toggle with the escape-aware delimiter serialize and parse, the same-index-first text-matched row reconciliation that preserves extra properties, skips unchanged cells, and keeps row order when an edit duplicates a line, the per-state icon overrides, and the required-cell validation expressions
3. `src/form/multiColumnRepeatableFieldset_test.go` — Renders, wiring, column input types, empty state, hidden input, NewRowFunc, style group, button style, and bulk-edit tests
4. `src/form/multiColumnRepeatableFieldset_templ.go` — Compiled output

---

## Range Slider

Input range slider control with min/max constraints, step values, bidirectional Alpine.js state binding for numeric input, optional tick marks, value bubbles with an always or hover display mode, an optional track hover tooltip that previews the value under the pointer, and static or one-way-state accessible names for single and dual thumbs.

**Flow:**

1. `src/control/rangeSlider.templ` — Component definition with RangeSliderSettings struct supporting single and dual thumb state paths and the AriaLabelOneWayStatePath thumb binding
2. `src/control/rangeSlider_templ.go` — Compiled output rendering input[type="range"] with custom styling and state binding

---

## Alert/Notification Box

Dismissible alert component with title, description, icons (left/right), and variations for success, warning, error, and info states.

**Flow:**

1. `src/display/alert.templ` — Component definition with AlertSettings struct exposing variants, sizes, and icon options
2. `src/display/alert_templ.go` — Compiled output rendering alert box with Tailwind styling and optional close button

---

## Chart

Line, area, bar, horizontal bar, pie, donut, radial bar, sunburst, parliament, and waffle charts drawn with plain Vega, plus a progress bar and a gauge drawn without Vega. The settings split by chart kind: the shared settings hold the X axis, the series, the legend, and the tooltip, while the type-specific options live in a Cartesian, a Slices, a Waffle, a Sunburst, a Parliament, a Progress, or a Gauge group. Slices is the group every donut-like variation reads: the donut, the pie, the radial bar, the sunburst, the parliament, and the waffle, so the donut stays one type constant among them. A pie is a donut with a forced zero hole and a hidden center label, and each slice labels its share of the total at its midpoint. A radial bar is a circular bar chart that grows each category from the Slices hole to the outer edge; each bar labels its share of the largest value at its midpoint. A sunburst draws a category ring and a child ring over one synthetic hierarchy root, where each category arc spans the sum of its children; each arc labels its share of the total, and a child arc its share of its category. Every in-fill slice label sits on a rounded pill whose background is a customizable token. A parliament dots one seat per value over concentric semicircle rows in contiguous party blocks, with a configurable dot shape, dot gap, and dot size. A waffle tiles a grid of equal cells by each category's share of the total, with a configurable cell shape and gap. A horizontal bar swaps the category and value axes while `XAxis` keeps naming the category. Colors are Tailwind tokens. The Vega types support stacking, area fills, markers, line dashes, an empty state, a screen-reader table, and live rows from an Alpine path with a row cap; the progress type supports thresholds, bounds, min and max size classes, a label text size class, an indeterminate animation, diagonal stripes, label affixes, a percent or fixed label, a one-way bound value path, and a label pill with a customizable background; the gauge type supports colored bands, a semicircle mode, a center label with configurable value and label text sizes and a text offset away from the pivot, tick marks with labels placed outside or inside the track, configurable arc and needle geometry, and custom colors including the needle.

**Flow:**

1. `src/display/chart.go` — The shared settings and the per-kind Cartesian, Slices (the slice label color and pill background), Waffle (cell count, column count, cell shape, cell gap), Sunburst (the child row key), Parliament (row count, dot shape, dot gap, dot size), Progress, and Gauge (including the tick position) groups, the private wire structs with their JSON tags, the type resolver, the config resolver with defaults and option checks, the legend and tooltip classes, the progress and gauge resolvers, and the color token collector
2. `src/display/chart.templ` — The progress and gauge markup branches, then the JSON config script, legend, canvas, empty state, tooltip, color probes, and table markup
3. `src/display/chartState.js` — Color probe reader, Vega spec builders per chart type with orientation-aware scales, axes, and bar channels, the pie, radial bar, and sunburst percentage labels over their pill backgrounds, the parliament dot and waffle cell shape channels, the progress Alpine data, the validated palette, the Vega view lifecycle outside Alpine state, the rows watch, and the tooltip handler
4. `src/import/import.templ` — HeadTagsVega loads Vega once for the page, usually from the layout head
5. `demo/chartDemo.templ` and `demo/chartDemoState.js` — Live streaming area chart and the line, bar, horizontal bar, stacked, donut (with the center label), pie, radial bar, sunburst, parliament, parliament dot shape, waffle, waffle shape and gap, progress, progress thresholds, progress bound value, gauge, gauge semicircle, gauge ticks, appearance, and empty examples

---

## Modal/Dialog

Overlay modal dialog with header (title), body content, footer, viewport-percentage sizing, an animated size transition, backdrop, and close/resize handlers using Alpine.js visibility binding.

**Flow:**

1. `src/display/modal.templ` — Component definition with ModalSettings struct supporting IsVisibleTwoWayStatePath, the percentage InitialSize scale (xs 40% through xxl 90% and full), the WidthPercent and HeightPercent overrides, and the PossibleSizes enlarge and reduce range
2. `src/display/modal.go` — Size resolvers, the percentage class builder, the resizable resolver, and the enlarge and reduce expression builders
3. `src/display/modal_templ.go` — Compiled output rendering backdrop and modal box with Alpine.js visibility and event management

---

## Confirmation Dialog Presets

Confirm, warning, critical, and delete dialog presets over Modal. Each preset supplies icon, tone, and confirmation copy. Delete and critical gate the confirm action behind typing the target name or id.

**Flow:**

1. `src/display/confirmationDialog.templ` — ConfirmationDialogSettings struct, the confirmationDialog engine, and the ConfirmDialog, WarningDialog, CriticalDialog, and DeleteDialog presets
2. `src/display/confirmationDialog.go` — Tone resolver, target match path resolver, the type-to-confirm field label resolver, and the type-to-confirm disabled expression builder
3. `src/display/confirmationDialog_test.go` — Table-driven tests for the tone, match path, and disabled expression helpers
4. `src/display/headerBlock.templ` — Header block with title, sub-heading, icon, and actions slot used by the dialog header
5. `src/display/headerIcon.templ` — Icon chip with color, background, padding, and radius settings
6. `src/form/inputField.templ` — Text input used for the type-to-confirm field
7. `src/control/button.templ` — Cancel and confirm buttons; the confirm button carries the disabled expression

---

## Header Block

Shared header row (title, sub-heading, icon, actions) behind page headings, cards, and confirmation dialogs, with size and color control and left or top icon placement.

**Flow:**

1. `src/display/headerBlock.templ` — HeaderBlockSettings struct and the HeaderBlock component with the HeaderSize scale and HeadingLevel setting
2. `src/display/headerIcon.templ` — HeaderIconSettings struct and the HeaderIcon chip component
3. `src/display/headerBlock_templ.go` — Compiled output rendering the header row
4. `src/display/headerIcon_templ.go` — Compiled output rendering the icon chip

---

## Sidebar Navigation

Collapsible vertical sidebar panel for navigation with sections and state management for collapse/expand behavior.

**Flow:**

1. `src/structural/sidebar.templ` — Component definition with SidebarSettings struct
2. `src/structural/sidebar_templ.go` — Compiled output rendering sidebar with Alpine.js state and collapsible sections

---

## Tag/Badge

Small label/badge component for categorization and tagging with customizable size, color, optional icons, and an optional remove button for filter chips.

**Flow:**

1. `src/display/tag.templ` — Component definition with TagSettings struct; `OnRemoveFunc` renders a named remove button; `OuterRadius` and `InnerRadius` default to `rounded-md`
2. `src/display/tag_templ.go` — Compiled output rendering small badge element with Tailwind styling

---

## Page Heading

One heading component over HeaderBlock. A Level setting picks the variant: page renders the h1 page header with page spacing, section renders the h2 section heading with icon-chip defaults. Both carry an optional description line and a right-aligned action slot.

**Flow:**

1. `src/structural/pageHeading.templ` — PageHeadingSettings struct and the PageHeading component; wraps HeaderBlock at the h1 or h2 level by Level. `Description` fills the sub-heading line; `HeaderSubHeading` is the legacy name.
2. `src/structural/pageHeading_templ.go` — Compiled output rendering the heading block

---

## Card

Surface container with an optional header block and middle and footer content slots, with customizable border radius (including square edges), padding, gap, shadow, ring, and colors.

**Flow:**

1. `src/structural/card.templ` — CardSettings struct and the Card component; MiddleContent and FooterContent slots and an optional header block
2. `src/structural/card_templ.go` — Compiled output rendering the card surface

---

## Pagination

Page controls with a live readout, a page-number strip, first/previous/next/last buttons, and an items-per-page selector. Binds the page number and page size to Alpine.js state paths and calls `OnChangeFunc` after every change. The page number is zero-based; `ShouldUseOneBasedPageDisplay` labels the strip from 1 without changing the state. The page count derives from the item total and the bound page size, so the strip and the controls react when the page size changes. `IsHiddenWhenSinglePage` hides only the page-number controls while the records fit on one page; the readout and the items-per-page selector stay visible.

**Flow:**

1. `src/structural/pagination.templ` — Component definition with PaginationSettings; renders the readout, page-number strip, buttons, and items-per-page selector
2. `src/structural/pagination.go` — Builds the readout and page-count expressions from the state paths and item total, and holds the shared items-per-page choices and items-per-page resolvers used by DataTable and Carousel
3. `src/structural/paginationState.js` — Builds the page-number strip from the current page and page total, marking the current page with `aria-current`
4. `src/structural/pagination_test.go` — Table-driven tests for the readout, page-count, last-page-number, and items-per-page expressions
5. `src/structural/pagination_templ.go` — Compiled output rendering the navigation landmark

---

## Tabs

Horizontal and vertical tab headers bound to one selected value, with per-tab label, value, optional icon, and optional badge count. Content panels switch with Alpine `x-show`. Optional URL-hash syncing, tab side placement, icon position, and surface customization.

**Flow:**

1. `src/structural/tabs.templ` — Component definition with TabSettings and TabItemSettings; renders the tab list and the panels, with tab list and content surface settings, and caps the vertical tab list to MaxVisibleTabs when set
2. `src/structural/tabs.go` — The tab id prefix generator, the selected, click, keyboard, and URL-hash sync expression builders, and the orientation, side, alignment, border radius, aria-label, percent class, constraint class, and badge class resolvers
3. `src/structural/tabsState.js` — Alpine data component that measures the tallest tab and caps the vertical tab list to MaxVisibleTabs, re-measuring on list resize
4. `src/toolset/styleTokens.go` — Shared token-to-class resolvers that Tabs and the other components call
5. `src/structural/tabs_test.go` — Table-driven tests for the expression builders and resolvers, and the MaxVisibleTabs render behavior
6. `src/structural/tabs_templ.go` — Compiled output rendering the ARIA tablist with roving tabindex, arrow and Home/End keyboard handlers, and the x-show panels

---

## Filter Bar

Standalone filter bar that renders one editor per declared filter (text contains, enum select, multi-enum select, number range, date range), shows active filters as removable chips while at least one filter is active, and resets everything with clear-all. A multi-enum filter holds an array and the refresh URL repeats its parameter once per selected value.

**Flow:**

1. `src/structural/filterBar.templ` — Component definition with FilterBarSettings and FilterSettings; renders editors and chips bound to a values object; the FilterChipStyle struct carries each chip color token, the chip radius (default `md`), and the label text case (default lowercase)
2. `src/structural/filterBarState.js` — Alpine component with chip visibility, chip label, single-filter reset, clear-all, and any-active helpers
3. `src/form/selectInput.templ` — Multi-select dropdown editor (IsMultiSelect) for the multi-enum kind
4. `src/display/tag.templ` — Removable Tag variant used for the chips
5. `src/structural/filterBar_templ.go` — Compiled output

---

## Data Table

Generic server-driven table taking column definitions and rows. Adds sortable headers that pin to the scroll viewport by default, row selection with bulk actions, the filter bar, a default search box, header action slots, a pagination footer, loading and error states, and refresh from a query URL template. Uses `htmx.ajax` when HTMX is present and a `fetch` fallback otherwise.

**Flow:**

1. `src/structural/dataTable.templ` — Generic component (DataTable[T]) with DataTableSettings and DataTableColumnSettings; renders the filter bar, toolbar, table region, error row, and pagination; FilterChipStyle forwards to the shared filter bar; the settings script renders inside the component root; the scroll viewport carries `data-ui-data-table-scroll`
2. `src/structural/dataTable.go` — DataTableSettings and DataTableColumnSettings types with their methods (client settings, id, page size, pagination label, density, text case, checkbox shape, checkbox size, row stripe, root class, sticky header class, scroll container class), the named density and sort-direction types, and the URL placeholder constants. The shared TextAlignment type serves the DataTable column alignment; the shared HorizontalAlignment type serves the search box. PageNumber is zero-based and ShouldUseOneBasedPageDisplay changes only the strip labels. The `Initial*` fields seed client state; `HeaderClass`, `CellClass`, `RowClassResolver`, and `IsStriped` carry styling; `StickyHeaderBackgroundColor` paints the pinned header and `ItemsPerPageDropdownBackgroundColor` paints the items-per-page menu; the root never clips a child dropdown, and the scroll container rounds its top corners when the table body touches the root's top edge; `MinWidthClass` and `MaxWidthClass` size the surface; `MinHeightClass`, `MaxHeightClass`, and `MaxVisibleRows` size the scroll viewport; `IsHeaderStatic` unpins the header; `CheckboxCheckedColor` and `CheckboxUncheckedColor` carry the selection color
3. `src/structural/dataTableState.js` — Alpine data component: debounces refreshes, swaps the region carrying `data-ui-data-table`, owns selection and sort helpers, and measures the scroll viewport to show `MaxVisibleRows` rows; builds the refresh URL through `UiToolset.ServerFragmentRefreshComponent.refreshUrlBuilder`
4. `src/structural/serverFragmentRefreshComponentState.js` — Shared refresh helpers, all under `UiToolset.ServerFragmentRefreshComponent`: the refresh URL builder, the fragment fetcher, and the lifecycle (settings resolver, event watcher, error handler attacher, handler detacher, request debouncer, failure handler, loading clearer)
5. `src/structural/pagination.templ` — Table footer pagination with the items-per-page selector and its configurable menu background
6. `src/structural/filterBar.templ` — Table filter bar
7. `src/structural/dataTable_templ.go` — Compiled output

---

## Carousel

Multi-item slider over a server page. The server returns a chunk of items; the client slides a window of `ItemsPerView` inside that chunk. The visible count follows Tailwind's named breakpoints measured against the browser window. Adds prev/next arrows (outside or inside), dot indicators (top or bottom), a search box (top or bottom), swipe, optional autoplay with pause-on-hover, the filter bar below the track, a pagination footer, and loading and error states. Items can wrap in `CarouselItemTooltip`, which shows a hover tooltip that teleports to the document body so a scroll container or a modal cannot clip it. Refresh reuses the DataTable query URL template and the shared refresh URL builder, and swaps the track and pagination regions from one response through `htmx.ajax` when HTMX is present and `fetch` otherwise.

**Flow:**

1. `src/structural/carousel.templ` — Generic component (Carousel[T]) with CarouselSettings; renders the search box, arrows, track, dots, error state, filter bar, and pagination; FilterChipStyle forwards to the shared filter bar; the track sits in `data-ui-carousel` and the footer in `data-ui-carousel-pagination`; the settings script renders inside the component root
2. `src/structural/carousel.go` — CarouselSettings and the CarouselItemsPerViewSettings type, the surface/arrows/dots/item token constants and class resolvers, plus the id, items-per-view, items-per-page, search-box position, autoplay interval, aria-label, and color resolvers and the client-settings builder
3. `src/structural/carouselItemTooltip.templ` — CarouselItemTooltip: wraps one item in a hover and focus trigger scope and renders a tooltip that teleports to the document body, with content, position, color, and compact-ring settings
4. `src/structural/carouselItemTooltip_test.go` — Render tests for the teleport markup, the aria-describedby wiring, the plain wrapper without content, and the customization tokens
5. `src/structural/alignment.go` — The HorizontalAlignment and TextAlignment types and their class methods
6. `src/toolset/styleTokens.go` — Shared token-to-class resolvers that Carousel, Tabs, Card, Modal, Alert, HeaderIcon, Tag, and Toast call
7. `src/toolset/tooltip.go` and `src/toolset/tooltipState.js` — The shared tooltip once-handle, surface classes, position values, and Alpine data component that CarouselItemTooltip, the Button, and the InputHint render
8. `src/structural/carouselState.js` — Alpine data component: clamps the window, derives the responsive items-per-view from the named breakpoints against the window width, applies the track transform, handles prev/next/dot/swipe input, runs the autoplay timer with pause-on-hover, sets off-window items inert, and refreshes through `htmx.ajax` when HTMX is present and `fetch` otherwise, swapping both regions
9. `src/structural/serverFragmentRefreshComponentState.js` — Shared refresh URL builder for the query template, the placeholders, and the filter values, plus the refresh lifecycle helpers under `UiToolset.ServerFragmentRefreshComponent`
10. `src/structural/pagination.go` — The shared ItemsPerPage type and items-per-page resolvers
11. `src/structural/pagination.templ` — Chunk footer pagination with the items-per-page selector and its configurable menu background
12. `src/structural/filterBar.templ` — Chunk filter bar
13. `src/structural/carousel_templ.go` — Compiled output

---

## Toast Notification

Dismissible notification toast component with title, description, and Alpine.js state management for visibility and auto-dismiss.

**Flow:**

1. `src/display/toast.templ` — Component definition with ToastSettings struct exposing optional AutoDismissSeconds (defaults to 10s), a `neutral-800/90` background default, and a tone-colored ring default when RingColor is unset
2. `src/display/toastState.js` — Alpine toast state, the type-to-ring-color resolver, and HTMX response handling
3. `src/display/toast_templ.go` — Compiled output rendering toast element with Alpine.js binding and timer logic
4. `src/import/toolset/apiResponse.js` — API response message and outcome resolution
5. `src/import/toolset/jsonAjax.js` — JsonAjax response handling delegated to the toast store

---

## Loading Overlay

Full-screen overlay with loading spinner indicator and a fade transition, used to block interaction during asynchronous operations.

**Flow:**

1. `src/display/loadingOverlay.templ` — Component definition with LoadingOverlaySettings struct
2. `src/display/loadingOverlay_templ.go` — Compiled output rendering fixed overlay with spinner animation

---

## Cloak Loading

Overlay to hide/obscure content during loading. It renders a fixed full-viewport layer and hides itself after the configured delay.

**Flow:**

1. `src/display/cloakLoading.templ` — Component definition with CloakLoadingSettings struct
2. `src/display/cloakLoading_templ.go` — Compiled output rendering semi-transparent overlay over target content

---

## Accordion

Collapsible section component for grouping content into expandable panels. The configured radius rounds only the first item's top corners and the last item's bottom corners; middle items stay square. `IsSingleOpen` groups the items so one stays open at a time. `IsFirstItemOpen` renders the first item expanded and leaves the rest collapsed.

**Flow:**

1. `src/display/accordion.templ` — Component definition with AccordionSettings struct
2. `src/display/accordion_templ.go` — Compiled output rendering accordion structure with Alpine.js collapse behavior

---

## Collapsible Text

Text block truncated to a line limit with an inline lowercase show-more affordance after the ellipsis, rendered only when the text overflows. Optional Phosphor icons flank the text, the toggle icons take one name per state, and the size scale covers the text, the affordance, and the icons together.

**Flow:**

1. `src/display/collapsibleText.templ` — Component definition with CollapsibleTextSettings struct, the truncation target, the inline affordance with its per-state toggle icons, and the flanking icons
2. `src/display/collapsibleTextState.js` — Alpine data component: the display text and toggle label getters, the measured truncation against the line budget, the word-boundary cut, and the expand and collapse state
3. `src/display/collapsibleText_test.go` — Renders, size and label resolvers, icons, and color tests
4. `src/display/collapsibleText_templ.go` — Compiled output

---

## CDN Dependencies Setup

Consolidated export of all third-party library dependencies (Alpine.js, UnoCSS runtime with Tailwind-compatible presets, Phosphor Icons, Google Fonts, HTMX, Vega) as composable templ components for easy inclusion in application head tags.

**Flow:**

1. `src/import/import.templ` — Templ components exporting CDN links with SRI hashes (HeadTagsMinimal, HeadTagsFull, and specialized imports)
2. `demo/demoIndex.templ` — Usage example showing @uiImport.HeadTagsFull() in HTML head section

---

## JavaScript Toolset

Bundled utility functions for client-side operations: random password generation, Unix timestamp formatting, the select-dropdown open-upward decision, loading overlay toggle, JSON AJAX requests, API response message resolution, and Alpine.js lifecycle hooks.

**Flow:**

1. `src/import/toolset/index.js` — UiToolset assembly; sibling files own one concern each: Alpine state registration, loading overlay, API response resolution, JsonAjax, password generation, Unix timestamp formatting, and the select-dropdown open-upward resolver
2. `src/import/toolset/unixTimestamp.js` — UnixTimestamp.ToIsoDate and UnixTimestamp.ToRfc1123; both take Unix seconds, format in UTC, and reject a non-number or a year outside 0000-9999 with InvalidUnixTimestamp
3. `src/import/import.templ` — HeadTagsToolset() component concatenating the toolset files and embedding the minified result in a script tag via MinifierTemplateJs()
4. `src/toolset/minifier.go` — esbuild-based minifier called by MinifierTemplateJs() to minify JS before rendering

---

## CSS/JavaScript Minification

Utility for minifying JavaScript and CSS at compile time or runtime using esbuild, with configurable options and error fallback.

**Flow:**

1. `src/toolset/minifier.go` — Minifier() function wrapping esbuild transform API with support for JavaScript and CSS content types
2. `src/import/import.templ` — Uses MinifierTemplateJs() and MinifierTemplateCss() wrapper functions to minify the embedded toolset files and inline styles
3. `src/toolset/minifier_test.go` — Tests validating minification behavior and error handling

---

## Style Tokens

Shared token-to-class resolvers for components that expose a styling token. Each resolver maps a token to the matching Tailwind class and takes the fallback class the caller wants when the token is empty or unknown.

**Flow:**

1. `src/toolset/styleTokens.go` — the shared `BorderRadius*`, `Shape*`, `ShadowSize*`, `RingThickness*`, `PaddingSize*`, `GapSize*`, and `Animation*` token constants, plus `BackgroundColorClassResolver()`, `TextColorClassResolver()`, `BorderColorClassResolver()`, `BorderRadiusTokenResolver()`, `BorderRadiusClassResolver()`, `ShapeClassResolver()`, `ShadowClassResolver()`, `RingThicknessClassResolver()`, `RingClassResolver()`, `CompactRingClassResolver()`, `PaddingClassResolver()`, `CompactPaddingClassResolver()`, `GapClassResolver()`, and `AnimationClassResolver()`
2. `src/structural/carousel.go`, `carousel.templ`, `tabs.go`, `tabs.templ`, `card.templ` — structural callers
3. `src/display/modal.templ`, `alert.templ`, `headerIcon.templ`, `toast.templ`, `accordion.templ`, `tag.templ`, `confirmationDialog.templ` — display callers
4. `src/control/button.templ`, `rangeSlider.templ`, `src/form/checkboxInput.templ` — control and form callers

---

## Component Id

Shared helpers that build the DOM id a component exposes, so one page holds several copies without colliding on element ids, radio group names, or input names. `ComponentIdPrefixGenerator` mints a unique per-render id prefix from a component name and a counter, for components that render once and do not refresh: the tabs and the accordion. `HashComponentIdParts` hashes the inputs that identify a refreshing component so its id stays the same across the page render and every refresh render: the carousel and the data table.

**Flow:**

1. `src/toolset/componentIdPrefix.go` — `NewComponentIdPrefixGenerator()` and `GenerateNext()`
2. `src/structural/tabs.go` and `src/display/accordion.templ` — one generator per component
3. `src/toolset/componentIdHash.go` — `HashComponentIdParts()`
4. `src/structural/carousel.go` and `dataTable.go` — id derived from the component id hash

---

## Text Case

Shared text-case values and a class resolver for components that expose a casing setting. Each component applies the matching CSS transform to its label or title; `TextCaseNone` leaves the text as typed.

**Flow:**

1. `src/toolset/textCase.go` — `TextCaseNone`, `TextCaseLower`, `TextCaseUpper`, `TextCaseCapitalize`, and `TextCaseClassResolver()`
2. Form labels and placeholders — `InputField`, `PasswordInput`, `CheckboxInput`, `RadioInput`, `InlineCheckboxGroup`, `InlineRadioGroup`, `SelectInput`, `DualSelectInput`, `SearchableSelectInput`, `RemoteSearchableSelectInput`, `TextArea`, `ToggleSwitch`, `DateTimePicker`, `NumberRangeInputField`, `MultiColumnRepeatableFieldset`, and `FileUploadInput`
3. `src/control/button.templ` — button label
4. `src/display/tag.templ` — label segments; `src/display/accordion.templ` — item titles
5. `src/display/headerBlock.templ` — title and sub-heading, forwarded by `Card`, `PageHeading`, and `ConfirmationDialog`
6. `src/structural/dataTable.go` — header labels

---

## Static HTML Demo Generation

Build-time HTML generation showcasing all UI components with usage examples and navigation structure.

**Flow:**

1. `demo/demo.go` — `demoGenerator` renders DemoIndex() to docs/index.html, the advanced example panels to one JSON asset per component section in docs/assets/advancedExamples/, and the data table and carousel refresh fragments to one JSON asset each in docs/assets/, writing each output only when its content changed
2. `demo/demoIndex.templ` plus one `demo/*Demo.templ` file per component, framed by `demo/demoExample.templ` — Page structure, sidebar navigation, and every usage example
3. `demo/advancedExamples.go` — Render-time registry that swaps each advanced example panel for a placeholder and renders the panel alone into its section asset
4. `demo/advancedExampleLoader.js` — Browser-side loader that inserts a panel from its section asset on the first open
5. `demo/data.go` — Demo record type, 25 sample rows, filter declarations, table column definitions, and the settings builder that slices one page
6. `demo/dataTableDemoRouting.js` — Browser-side router that serves the requested page's pre-rendered fragment from the JSON asset as the refresh response
7. `src/import/import.templ` — DemoIndex imports HeadTagsFull() for CDN resources

---

## Test Suite

Single entry point for all verification: Go units, Playwright behavioral specs against the demo, axe ratchet, and golden performance budgets.

**Flow:**

1. `tests/tests.sh` — Contract entry point: registry selection, demo server and Lightpanda CDP server lifecycle, cumulative levels, exit 0/1/2
2. `tests/registry.yaml` — Explicit feature:scope/level registry; trusted input
3. `tests/lib/registry.mjs` — Registry reader for list and select modes
4. `tests/lib/serveDemo.mjs` — Localhost server for docs/ that proxies external script, stylesheet, and image URLs through a fetch-once cache, keeping CDN latency out of specs
5. `tests/lib/checkPerformance.mjs` — Compares measured latencies against `tests/golden.yaml` tiers
6. `tests/ui/run.sh` — Playwright mode runner: the fast and standard modes run on Lightpanda; the exhaustive modes (a11y, performance, and the cross-browser modes) run on Chromium, then Firefox, with 1 worker
7. `tests/ui/fixture.js` — The `test` every spec imports; connects the worker browser to Lightpanda over CDP when the endpoint is set
8. `tests/ui/specs/` — Behavioral specs by feature: form, control, display, structural, toolset, a11y, performance; tests that need a full browser engine carry the `@browser` tag
9. `tests/golden.yaml` — Golden budgets for the demo navigation, text area expand, select dropdown open, searchable select filter, number range slider drag, and Button action state

---

# Feature Map

> Auto-maintained index of every user-facing feature and the code path that implements it. Updated alongside the code — not after the fact.

## Text Input Field

Single-line text input with configurable type (text, email, number, date, password, etc.), with support for labels, hints, required indicators, and optional prefix/suffix affixes.

**Flow:**

1. `src/form/inputField.templ` — Component definition with InputFieldSettings struct exposing InputType, Label, TwoWayStatePath, Value, and optional affixes
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/inputField_templ.go` — Compiled templ output rendering HTML input with Alpine.js binding and Tailwind styling

---

## Multi-line Text Area

Multi-line text input with five height steps (h-12/24/36/48/60), expand-to-3x toggle, floating action icons (expand, copy, clear) anchored to the text line, and optional hints rendered as a focusable tooltip or a description line.

**Flow:**

1. `src/form/textArea.templ` — Component definition with TextAreaSettings struct; heights and icon positions supplied mutually exclusively via Alpine `:class`
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/textArea_templ.go` — Compiled output rendering textarea element with styling and Alpine.js integration

---

## Select Dropdown

Dropdown select component with support for flat string options or label-value pairs, with optional grouping and blank option.

**Flow:**

1. `src/form/selectInput.templ` — Component definition with SelectInputSettings struct and SelectLabelValueOption data structure
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/selectInput_templ.go` — Compiled output rendering select with native radio options, an embedded JSON script for label lookup, and Alpine.js state management

Supports optional hint text rendered either as a focusable info-icon tooltip inside the dropdown row or as a description line below the fieldset.

---

## Multi-Select Dropdown

Dropdown component that lets the user select multiple options from a flat list or label-value pairs, binding an array via Alpine.js two-way state path.

**Flow:**

1. `src/form/multiSelectInput.templ` — Component definition with MultiSelectInputSettings struct, reusing SelectLabelValueOption for label-value options
2. `src/form/multiSelectInputState.js` — Alpine.js data component providing the dropdown toggle state
3. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
4. `src/form/multiSelectInput_templ.go` — Compiled output rendering checkbox-based dropdown with embedded JSON script for label-value options and Alpine.js state management

Form submission uses multiple checkboxes sharing the same `name` so the browser sends an array of values. Supports optional hint text rendered either as a focusable info-icon tooltip inside the dropdown row or as a description line below the fieldset.

---

## Checkbox Input

Checkbox that selects one value, bound to a boolean or an array Alpine.js state, with configurable shapes, sizes, colors, label position, and disabled state.

**Flow:**

1. `src/form/checkboxInput.templ` — Component definition with CheckboxInputSettings for binding, shape, size, colors, indeterminate, disabled, and label position
2. `src/form/checkboxInput_templ.go` — Compiled output rendering the native input, the styled box, and the check or dash icon

The check renders through color inheritance: the box carries the checked color and the icon inherits it. The indeterminate state is set on the native input through `x-effect` and renders as a dash.

---

## Toggle Switch

Switch component that binds either a boolean Alpine.js state or a custom value in an Alpine.js array, with configurable sizes, colors, required state, and disabled state.

**Flow:**

1. `src/form/toggleSwitch.templ` — Component definition with ToggleSwitchSettings for binding, sizing, colors, disabled state, and label position
2. `src/form/toggleSwitch_templ.go` — Compiled output rendering the hidden form field, checkbox input, and switch track

With `CustomValue`, Alpine.js adds or removes the value from an array. Without it, Alpine.js treats the checkbox as a boolean state.

---

## Radio Button Input

Single radio button component for choice selection within a group, with label and state binding support.

**Flow:**

1. `src/form/radioInput.templ` — Component definition with RadioInputSettings struct
2. `src/form/radioInput_templ.go` — Compiled output rendering input[type="radio"] with Tailwind styling

---

## Inline Radio Group

Horizontal radio button group for presenting multiple mutually exclusive options in a single row.

**Flow:**

1. `src/form/inlineRadioGroup.templ` — Component definition with InlineRadioGroupSettings struct
2. `src/form/inlineRadioGroup_templ.go` — Compiled output rendering multiple radio inputs horizontally

---

## Button

Interactive button component with customizable label, icons (left/right using Phosphor), size/shape variants, colors, disabled state, and click handlers.

**Flow:**

1. `src/control/button.templ` — Component definition with ButtonSettings struct supporting OnClickFunc handlers and icon binding
2. `src/control/button_templ.go` — Compiled output rendering button element with Alpine.js event binding and optional tooltip support

---

## Range Slider

Input range slider control with min/max constraints, step values, bidirectional Alpine.js state binding for numeric input, optional tick marks, and accessible names for single and dual thumbs.

**Flow:**

1. `src/control/rangeSlider.templ` — Component definition with RangeSliderSettings struct supporting single and dual thumb state paths
2. `src/control/rangeSlider_templ.go` — Compiled output rendering input[type="range"] with custom styling and state binding

---

## Alert/Notification Box

Dismissible alert component with title, description, icons (left/right), and variations for success, warning, error, and info states.

**Flow:**

1. `src/display/alert.templ` — Component definition with AlertSettings struct exposing variants, sizes, and icon options
2. `src/display/alert_templ.go` — Compiled output rendering alert box with Tailwind styling and optional close button

---

## Modal/Dialog

Overlay modal dialog with header (title), body content, footer, customizable size, backdrop, and close/resize handlers using Alpine.js visibility binding.

**Flow:**

1. `src/display/modal.templ` — Component definition with ModalSettings struct supporting IsVisibleTwoWayStatePath
2. `src/display/modal_templ.go` — Compiled output rendering backdrop and modal box with Alpine.js visibility and event management

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

1. `src/display/tag.templ` — Component definition with TagSettings struct; `OnRemoveFunc` renders a named remove button
2. `src/display/tag_templ.go` — Compiled output rendering small badge element with Tailwind styling

---

## Pagination

Page controls with a live readout, a page-number strip, first/previous/next/last buttons, and an items-per-page selector. Binds the page number and page size to Alpine.js state paths and calls `OnChangeFunc` after every change.

**Flow:**

1. `src/structural/pagination.templ` — Component definition with PaginationSettings; renders the readout, page-number strip, buttons, and items-per-page selector
2. `src/structural/pagination.go` — Builds the readout expression from the state paths and item total
3. `src/structural/paginationState.js` — Builds the page-number strip from the current page and page total, marking the current page with `aria-current`
4. `src/structural/pagination_test.go` — Table-driven tests for the readout expression
5. `src/structural/pagination_templ.go` — Compiled output rendering the navigation landmark

---

## Filter Bar

Standalone filter bar that renders one editor per declared filter (text contains, enum select, number range, date range), shows active filters as removable chips, and resets everything with clear-all.

**Flow:**

1. `src/structural/filterBar.templ` — Component definition with FilterBarSettings and FilterSettings; renders editors and chips bound to a values object
2. `src/structural/filterBar.go` — Alpine expression builders for chip values, chip visibility, single-filter reset, and clear-all
3. `src/display/tag.templ` — Removable Tag variant used for the chips
4. `src/structural/filterBar_templ.go` — Compiled output

---

## Data Table

Generic server-driven table taking column definitions and rows. Adds sortable headers, row selection with bulk actions, the filter bar, a search slot, header action slots, a pagination footer, loading and error states, and refresh from a URL template. Uses `htmx.ajax` when HTMX is present and a `fetch` fallback otherwise.

**Flow:**

1. `src/structural/dataTable.templ` — Generic component (DataTable[T]) with DataTableSettings and DataTableColumnSettings; renders the filter bar, toolbar, table region, error row, and pagination
2. `src/structural/dataTable.go` — DataTableSettings and DataTableColumnSettings types with their methods (client settings, id, page size, pagination label, density, header case, checkbox shape, checkbox size, row stripe), the named density/alignment/sort-direction/header-text-case/page-size types, and the URL placeholder constants. The alignment type carries the column alignment and justify class methods. The `Initial*` fields seed client state; `HeaderClass`, `CellClass`, `RowClassResolver`, and `IsStriped` carry styling; `CheckboxCheckedColor` and `CheckboxUncheckedColor` carry the selection color
3. `src/structural/dataTableState.js` — Alpine data component: builds the refresh URL from the template, debounces refreshes, swaps the region carrying `data-ui-data-table`, and owns selection and sort helpers
4. `src/structural/pagination.templ` — Table footer pagination
5. `src/structural/filterBar.templ` — Table filter bar
6. `src/structural/dataTable_templ.go` — Compiled output

---

## Toast Notification

Dismissible notification toast component with title, description, and Alpine.js state management for visibility and auto-dismiss.

**Flow:**

1. `src/display/toast.templ` — Component definition with ToastSettings struct exposing optional AutoDismissSeconds (defaults to 10s)
2. `src/display/toastState.js` — Alpine toast state and HTMX response handling
3. `src/display/toast_templ.go` — Compiled output rendering toast element with Alpine.js binding and timer logic
4. `src/import/toolset/apiResponse.js` — API response message and outcome resolution
5. `src/import/toolset/jsonAjax.js` — JsonAjax response handling delegated to the toast store

---

## Loading Overlay

Full-screen overlay with loading spinner indicator, used to block interaction during asynchronous operations.

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

Collapsible section component for grouping content into expandable panels.

**Flow:**

1. `src/display/accordion.templ` — Component definition with AccordionSettings struct
2. `src/display/accordion_templ.go` — Compiled output rendering accordion structure with Alpine.js collapse behavior

---

## CDN Dependencies Setup

Consolidated export of all third-party library dependencies (Alpine.js, UnoCSS runtime with Tailwind-compatible presets, Phosphor Icons, Google Fonts, HTMX) as composable templ components for easy inclusion in application head tags.

**Flow:**

1. `src/import/import.templ` — Templ components exporting CDN links with SRI hashes (HeadTagsMinimal, HeadTagsFull, and specialized imports)
2. `demo/demo.templ` — Usage example showing @uiImport.HeadTagsFull() in HTML head section

---

## JavaScript Toolset

Bundled utility functions for client-side operations: random password generation, loading overlay toggle, JSON AJAX requests, API response message resolution, and Alpine.js lifecycle hooks.

**Flow:**

1. `src/import/toolset/index.js` — UiToolset assembly; sibling files own one concern each: Alpine state registration, loading overlay, API response resolution, JsonAjax, and password generation
2. `src/import/import.templ` — HeadTagsToolset() component concatenating the toolset files and embedding the minified result in a script tag via MinifierTemplateJs()
3. `src/toolset/minifier.go` — esbuild-based minifier called by MinifierTemplateJs() to minify JS before rendering

---

## CSS/JavaScript Minification

Utility for minifying JavaScript and CSS at compile time or runtime using esbuild, with configurable options and error fallback.

**Flow:**

1. `src/toolset/minifier.go` — Minifier() function wrapping esbuild transform API with support for JavaScript and CSS content types
2. `src/import/import.templ` — Uses MinifierTemplateJs() and MinifierTemplateCss() wrapper functions to minify the embedded toolset files and inline styles
3. `src/toolset/minifier_test.go` — Tests validating minification behavior and error handling

---

## Static HTML Demo Generation

Build-time HTML generation showcasing all UI components with usage examples and navigation structure.

**Flow:**

1. `demo/demo.go` — Entrypoint that renders DemoIndex() to docs/index.html and one data table refresh fragment per page plus an all-records fragment to docs/assets/
2. `demo/demo.templ` — Full demo page structure with component usage examples, sidebar navigation, and styling
3. `demo/data.go` — Demo record type, 25 sample rows, filter declarations, table column definitions, and the settings builder that slices one page
4. `demo/dataTableDemoRouting.js` — Browser-side router that rewrites each refresh request to the fragment for the requested page
5. `src/import/import.templ` — DemoIndex imports HeadTagsFull() for CDN resources

---

## Test Suite

Single entry point for all verification: Go units, Playwright behavioral specs against the demo, axe ratchet, and golden performance budgets.

**Flow:**

1. `tests/tests.sh` — Contract entry point: registry selection, demo server lifecycle, cumulative levels, exit 0/1/2
2. `tests/registry.yaml` — Explicit feature:scope/level registry; trusted input
3. `tests/lib/registry.mjs` — Registry reader for list and select modes
4. `tests/lib/serveDemo.mjs` — Localhost server for docs/ that proxies external script, stylesheet, and image URLs through a fetch-once cache, keeping CDN latency out of specs
5. `tests/lib/checkPerformance.mjs` — Compares measured latencies against `tests/golden.yaml` tiers
6. `tests/ui/run.sh` — Playwright mode runner (smoke, standard, a11y, performance, toolset, control, structural-smoke, structural, structural-a11y, cross-browser, toolset-cross-browser, control-cross-browser, structural-cross-browser)
7. `tests/ui/specs/` — Behavioral specs by feature: form, control, structural, toolset, a11y, performance

---

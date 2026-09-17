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

1. `src/control/rangeSlider.templ` — Component definition with RangeSliderSettings struct supporting TwoWayStatePath
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

1. `src/display/sidebar.templ` — Component definition with SidebarSettings struct
2. `src/display/sidebar_templ.go` — Compiled output rendering sidebar with Alpine.js state and collapsible sections

---

## Tag/Badge

Small label/badge component for categorization and tagging with customizable size, color, and optional icons.

**Flow:**

1. `src/display/tag.templ` — Component definition with TagSettings struct
2. `src/display/tag_templ.go` — Compiled output rendering small badge element with Tailwind styling

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

Overlay to hide/obscure content during loading without blocking the entire viewport.

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

1. `demo/demo.go` — Entrypoint that renders DemoIndex() templ component and writes to index.html
2. `demo/demo.templ` — Full demo page structure with component usage examples, sidebar navigation, and styling
3. `src/import/import.templ` — DemoIndex imports HeadTagsFull() for CDN resources

---

## Test Suite

Single entry point for all verification: Go units, Playwright behavioral specs against the demo, axe ratchet, and golden performance budgets.

**Flow:**

1. `tests/tests.sh` — Contract entry point: registry selection, demo server lifecycle, cumulative levels, exit 0/1/2
2. `tests/registry.yaml` — Explicit feature:scope/level registry; trusted input
3. `tests/lib/registry.mjs` — Registry reader for list and select modes
4. `tests/lib/serveDemo.mjs` — Localhost server for docs/ that proxies external script, stylesheet, and image URLs through a fetch-once cache, keeping CDN latency out of specs
5. `tests/lib/checkPerformance.mjs` — Compares measured latencies against `tests/golden.yaml` tiers
6. `tests/ui/run.sh` — Playwright mode runner (smoke, standard, a11y, performance, toolset, control, cross-browser, toolset-cross-browser, control-cross-browser)
7. `tests/ui/specs/` — Behavioral specs by feature: form, control, toolset, a11y, performance

---

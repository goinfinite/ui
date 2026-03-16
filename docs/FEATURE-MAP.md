# Feature Map

> Auto-maintained index of every user-facing feature and the code path that implements it. Updated alongside the code — not after the fact.

## Text Input Field

Single-line text input with configurable type (text, email, number, date, password, etc.), with support for labels, hints, required indicators, and optional prefix/suffix affixes.

**Flow:**

1. `src/form/inputField.templ` — Component definition with InputFieldSettings struct exposing InputType, Label, TwoWayStatePath, Value, and optional affixes
2. `src/form/inputField_templ.go` — Compiled templ output rendering HTML input with Alpine.js binding and Tailwind styling

---

## Multi-line Text Area

Multi-line text input component with configurable rows and full support for form submission and state binding.

**Flow:**

1. `src/form/textArea.templ` — Component definition with TextAreaSettings struct
2. `src/form/textArea_templ.go` — Compiled output rendering textarea element with styling and Alpine.js integration

---

## Select Dropdown

Dropdown select component with support for flat string options or label-value pairs, with optional grouping and blank option.

**Flow:**

1. `src/form/selectInput.templ` — Component definition with SelectInputSettings struct and SelectLabelValueOption data structure
2. `src/form/selectInput_templ.go` — Compiled output rendering select with embedded JSON script for options and Alpine.js state management

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

Input range slider control with min/max constraints, step values, and bidirectional Alpine.js state binding for numeric input.

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

1. `src/display/toast.templ` — Component definition with ToastSettings struct
2. `src/display/toastState.js` — Helper JavaScript for managing toast queue and lifecycle
3. `src/display/toast_templ.go` — Compiled output rendering toast element with Alpine.js binding and timer logic

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

Bundled utility functions for client-side operations: random password generation, loading overlay toggle, JSON AJAX requests, and Alpine.js lifecycle hooks.

**Flow:**

1. `src/import/toolset.js` — JavaScript utility library with UiToolset functions embedded as string
2. `src/import/import.templ` — HeadTagsToolset() component embedding minified toolset.js in a script tag via MinifierTemplateJs()
3. `src/toolset/minifier.go` — esbuild-based minifier called by MinifierTemplateJs() to minify JS before rendering

---

## CSS/JavaScript Minification

Utility for minifying JavaScript and CSS at compile time or runtime using esbuild, with configurable options and error fallback.

**Flow:**

1. `src/toolset/minifier.go` — Minifier() function wrapping esbuild transform API with support for JavaScript and CSS content types
2. `src/import/import.templ` — Uses MinifierTemplateJs() and MinifierTemplateCss() wrapper functions to minify embedded toolset.js and inline styles
3. `src/toolset/minifier_test.go` — Tests validating minification behavior and error handling

---

## Static HTML Demo Generation

Build-time HTML generation showcasing all UI components with usage examples and navigation structure.

**Flow:**

1. `demo/demo.go` — Entrypoint that renders DemoIndex() templ component and writes to index.html
2. `demo/demo.templ` — Full demo page structure with component usage examples, sidebar navigation, and styling
3. `src/import/import.templ` — DemoIndex imports HeadTagsFull() for CDN resources

---

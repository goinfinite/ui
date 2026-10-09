# Feature Map

> Auto-maintained index of every user-facing feature and the code path that implements it. Updated alongside the code — not after the fact.

## Text Input Field

Single-line text input with configurable type (text, email, number, date, password, etc.), with support for labels, label text case, hints, required indicators, a size scale, optional prefix/suffix affixes, an arbitrary right-edge component slot, and a runtime type expression.

**Flow:**

1. `src/form/inputField.templ` — Component definition with InputFieldSettings struct exposing InputType, Label, TextCase, TwoWayStatePath, Value, Size, optional affixes, affix width percentages, AffixRightComponent for an arbitrary right-edge component, and InputTypeExpression for a runtime type
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/inputField_templ.go` — Compiled templ output rendering HTML input with Alpine.js binding and Tailwind styling

---

## Password Input

Password field with a reveal toggle, an optional random password generator, and an optional strength meter. The meter shows a percentage bar plus a criteria checklist that updates as the user types.

**Flow:**

1. `src/form/passwordInput.templ` — PasswordInputSettings, the rules resolver (length range, generation length, character classes), the criteria-item resolver, the reveal and generate action buttons rendered as uiControl.Button with padding and min-width/min-height constraints and a boxed or plain style (ActionButtonStyle), the meter with its color settings, and the InputField composition
2. `src/form/passwordInputState.js` — Alpine data component: reveal flag, strength percentage, per-criterion flags, the generator call, and the clipboard copy with a language-matched toast that ends with an exclamation mark
3. `src/import/toolset/randomPassword.js` — `createRandomPassword(options)` with configurable length and character classes; guarantees one character from each selected class
4. `src/form/passwordInput_test.go` — Table-driven tests for the rules and criteria resolvers plus the rendered actions and meter
5. `src/form/passwordInput_templ.go` — Compiled output

---

## Multi-line Text Area

Multi-line text input with five height steps (h-12/24/36/48/60), expand-to-3x toggle, floating action icons (expand, copy, clear) anchored to the text line, and optional hints rendered as a focusable tooltip or a description line.

**Flow:**

1. `src/form/textArea.templ` — Component definition with TextAreaSettings struct; heights and icon positions supplied mutually exclusively via Alpine `:class`
2. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
3. `src/form/textArea_templ.go` — Compiled output rendering textarea element with styling and Alpine.js integration

---

## Select Dropdown

Dropdown select component with support for flat string options or label-value pairs, a blank option, and a multi-select mode. Single mode renders native radio rows; multi mode renders the shared styled checkbox rows, keeps the dropdown open after each toggle, and joins the selected labels in the trigger.

**Flow:**

1. `src/form/selectInput.templ` — Component definition with SelectInputSettings struct (including IsMultiSelect and InputId), SelectLabelValueOption data structure, the option checkbox renderer, and the dropdown height and MaxVisibleOptions settings
2. `src/form/selectInputState.js` — Alpine.js data component for the select dropdown: open toggle, close, the multi-select summary formatter, and the open-direction decision
3. `src/import/toolset/selectDropdown.js` — Shared open-upward decision that measures the space above and below the trigger, bounded by the viewport and every clipping ancestor, so the list never opens downward into an overflow-hidden container
4. `src/form/selectDropdown.go` — Shared dropdown height classes resolver and the MaxVisibleOptions height style, with the row height of each dropdown family
5. `src/form/inputHint.templ` — Shared hint renderer for the tooltip and description display modes
6. `src/form/selectInput_templ.go` — Compiled output rendering select with native radio or checkbox options, an embedded JSON script for label lookup and the multi-select summary, and Alpine.js state management

Supports optional hint text rendered either as a focusable info-icon tooltip inside the dropdown row or as a description line below the fieldset.

---

## Searchable Select

Dropdown with a filter box that narrows local options in place, a clear button, a No matches row, and a max-height scroll list. Single mode binds one value; multi mode binds an array and shows the joined labels or removable tags. With tags, a typed value becomes a tag when custom values are enabled; Enter adds, Backspace removes the last tag, and each tag carries a remove button.

**Flow:**

1. `src/form/searchableSelectInput.templ` — SearchableSelectInputSettings and the shared shell, trigger, filter input, action icons, selected labels, selected tags with tag color settings, local option rows, and local dropdown
2. `src/form/searchableSelectInput.go` — The client config types and resolvers (normalized items, multi-select and custom-value flags, remote settings, the selection presentation), the dropdown class resolver, plus the selection, toggle, clear, custom-value, tag-remove, tag-backspace, and selected-state expression builders
3. `src/form/searchableSelectInputState.js` — Shared Alpine data component: filter, visibility, No matches, selected-label lookup, and the remote fetch lifecycle
4. `src/import/toolset/selectDropdown.js` — Shared open-upward decision that measures the space above and below the trigger, bounded by the viewport and every clipping ancestor
5. `src/form/searchableSelectInput_test.go` — Tests for the item and presentation resolvers, the expression builders, and the rendered dropdown
6. `src/form/searchableSelectInput_templ.go` — Compiled output

Single mode writes one hidden input; multi mode writes one hidden input per selected value under `InputName`.

---

## Remote Searchable Select

Searchable select whose options come from a URL as the user types, with debounce, minimum query length, stale-response discard, a loading indicator, and error and empty states.

**Flow:**

1. `src/form/remoteSearchableSelectInput.templ` — RemoteSearchableSelectInputSettings, the remote config resolver, the dropdown height and MaxVisibleOptions settings, the InitialOptions that label bound values until the user types, and the remote option rows and dropdown with the minimum-length prompt, sharing the searchable shell, trigger, and filter input
2. `src/form/searchableSelectInputState.js` — Shared Alpine data component owning the remote fetch lifecycle: request id for stale discard, loading and error flags, and the debounced request
3. `src/form/remoteSearchableSelectInput_test.go` — Tests for the remote config resolver
4. `src/form/remoteSearchableSelectInput_templ.go` — Compiled output

The endpoint receives `GET OptionsUrl?<query param>=<typed text>` and returns `{ "body": [{ "label": "...", "value": "..." }] }`; a bare array of strings or label-value objects also works.

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

Radio button group for presenting multiple mutually exclusive options in a row or a column under one notched label.

**Flow:**

1. `src/form/inlineRadioGroup.templ` — Component definition with InlineRadioGroupSettings struct and the orientation resolver
2. `src/form/inlineRadioGroup_templ.go` — Compiled output rendering multiple radio inputs horizontally or vertically

---

## Inline Checkbox Group

Checkbox group for presenting several choices in a row or a column under one notched label, bound to one array state path.

**Flow:**

1. `src/form/inlineCheckboxGroup.templ` — Component definition with InlineCheckboxGroupSettings struct and the orientation resolver, composing CheckboxInput per setting
2. `src/form/inlineCheckboxGroup_templ.go` — Compiled output rendering the checkboxes horizontally or vertically

---

## Button

Interactive button component with customizable label, icons (left/right using Phosphor), size/shape variants, colors, disabled state, click handlers, and an optional hover tooltip. Every button renders a `type` attribute: `button` by default, `submit` when `IsSubmit` is set.

**Flow:**

1. `src/control/button.templ` — Component definition with ButtonSettings struct supporting OnClickFunc handlers, icon binding, a padding override with min/max width and height constraints, a static or Alpine-bound aria label, the type attribute, and the optional tooltip
2. `src/toolset/tooltipState.js` — Shared Alpine data component that teleports a tooltip to the document body, positions it with fixed coordinates on hover or focus, and flips it to stay inside the viewport
3. `src/toolset/tooltip.go` — The once-handle that embeds the shared tooltip state, plus TooltipSurfaceClassesResolver for the shared fixed-layer tooltip classes and the shared TooltipPosition* placement values
4. `src/control/button_templ.go` — Compiled output rendering the button with Alpine.js event binding and the teleported tooltip

---

## Range Slider

Input range slider control with min/max constraints, step values, bidirectional Alpine.js state binding for numeric input, optional tick marks, value bubbles with an always or hover display mode, an optional track hover tooltip that previews the value under the pointer, and accessible names for single and dual thumbs.

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

1. `src/structural/filterBar.templ` — Component definition with FilterBarSettings and FilterSettings; renders editors and chips bound to a values object; the FilterChipStyle struct carries each chip color token and the label text case (default lowercase)
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

1. `src/display/toast.templ` — Component definition with ToastSettings struct exposing optional AutoDismissSeconds (defaults to 10s), a `neutral-800` background default, and a tone-colored ring default when RingColor is unset
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

## CDN Dependencies Setup

Consolidated export of all third-party library dependencies (Alpine.js, UnoCSS runtime with Tailwind-compatible presets, Phosphor Icons, Google Fonts, HTMX) as composable templ components for easy inclusion in application head tags.

**Flow:**

1. `src/import/import.templ` — Templ components exporting CDN links with SRI hashes (HeadTagsMinimal, HeadTagsFull, and specialized imports)
2. `demo/demoIndex.templ` — Usage example showing @uiImport.HeadTagsFull() in HTML head section

---

## JavaScript Toolset

Bundled utility functions for client-side operations: random password generation, the select-dropdown open-upward decision, loading overlay toggle, JSON AJAX requests, API response message resolution, and Alpine.js lifecycle hooks.

**Flow:**

1. `src/import/toolset/index.js` — UiToolset assembly; sibling files own one concern each: Alpine state registration, loading overlay, API response resolution, JsonAjax, password generation, and the select-dropdown open-upward resolver
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

## Style Tokens

Shared token-to-class resolvers for components that expose a styling token. Each resolver maps a token to the matching Tailwind class and takes the fallback class the caller wants when the token is empty or unknown.

**Flow:**

1. `src/toolset/styleTokens.go` — the shared `BorderRadius*`, `Shape*`, `ShadowSize*`, `RingThickness*`, `PaddingSize*`, and `GapSize*` token constants, plus `BackgroundColorClassResolver()`, `TextColorClassResolver()`, `BorderColorClassResolver()`, `BorderRadiusTokenResolver()`, `BorderRadiusClassResolver()`, `ShapeClassResolver()`, `ShadowClassResolver()`, `RingThicknessClassResolver()`, `RingClassResolver()`, `CompactRingClassResolver()`, `PaddingClassResolver()`, `CompactPaddingClassResolver()`, and `GapClassResolver()`
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
2. Form labels and placeholders — `InputField`, `CheckboxInput`, `RadioInput`, `InlineRadioGroup`, `SelectInput`, `TextArea`, `ToggleSwitch`
3. `src/control/button.templ` — button label
4. `src/display/tag.templ` — label segments; `src/display/accordion.templ` — item titles
5. `src/display/headerBlock.templ` — title and sub-heading, forwarded by `Card`, `PageHeading`, and `ConfirmationDialog`
6. `src/structural/dataTable.go` — header labels

---

## Static HTML Demo Generation

Build-time HTML generation showcasing all UI components with usage examples and navigation structure.

**Flow:**

1. `demo/demo.go` — `demoGenerator` renders DemoIndex() to docs/index.html and the data table and carousel refresh fragments to one JSON asset each in docs/assets/, writing each output only when its content changed
2. `demo/demoIndex.templ` plus one `demo/*Demo.templ` file per component, framed by `demo/demoExample.templ` — Page structure, sidebar navigation, and every usage example
3. `demo/data.go` — Demo record type, 25 sample rows, filter declarations, table column definitions, and the settings builder that slices one page
4. `demo/dataTableDemoRouting.js` — Browser-side router that serves the requested page's pre-rendered fragment from the JSON asset as the refresh response
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
7. `tests/ui/specs/` — Behavioral specs by feature: form, control, display, structural, toolset, a11y, performance

---

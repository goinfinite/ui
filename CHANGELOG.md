# Changelog

```log
0.2.5 - 2026/10/07
BREAKING: move the refresh URL builder, the fragment fetcher, and the refresh lifecycle helpers under UiToolset.ServerFragmentRefreshComponent
refactor: share the refresh settings, event wiring, error handling, and cleanup between the carousel and the data table
refactor: rename the shared refresh files and helpers so each name states its subject
test: share the filter settings assertion between the carousel and the data table tests
test: gather the ring and the padding resolver cases in one table each
fix: use the safer bash conditional syntax in the test runners
fix: mark the demo file serve promise as ignored
fix: render the derived carousel and data table id on the component root
fix: resolve the value bubble ring classes from the resolved ring color so an upper-only color renders
fix: fail the carousel htmx refresh when the response does not replace the pagination region
fix: separate the component id hash parts so two different part boundaries cannot collide
fix: reposition the tooltip on a viewport change without revealing it
BREAKING: replace the IsHeaderSticky setting with IsHeaderStatic and pin the data table header by default
BREAKING: replace the IsSearchBoxFullWidth setting with IsSearchBoxCompact and render the search box full width by default
feat: add min and max width and height classes and a max visible rows setting to the data table viewport
feat: add a max visible tabs setting to the tabs component
fix: bind the alert title icon state path classes
fix: keep a date filter input inside its range editor column
feat: forward the tooltip width and height caps through every hint-bearing field
feat: add an icon-left example to the input field demo
feat: add min and max width and height caps to the shared tooltip through TooltipSurfaceSettings
feat: add a magnifier icon to the default carousel and data table search boxes
feat: add an optional leading icon to InputField
fix: clip the carousel item wrapper so a narrow item cannot paint into the track gap
fix: keep the arrow hover feedback when a custom arrows background color is set
fix: size the carousel arrows as fixed squares so they render circular instead of tall
BREAKING: rename ButtonTooltipPosition* and CarouselItemTooltipPosition* to the shared uiToolset.TooltipPosition* constants
feat: render the InputHint tooltip through the shared toolset tooltip so a scroll container or a modal cannot clip it
feat: add the CarouselItemTooltip component with text or component content, placement, and compact-ring settings
feat: add the carousel component with responsive items per view, prev/next arrows, dot indicators, swipe, and optional autoplay with pause-on-hover
feat: add the server pagination and the filter bar to the carousel
feat: add the surface, arrow, dot, and item styling to the carousel
feat: add a top or bottom position to the carousel search box
feat: share the refresh request builder and the filter-value, page-size, query-param, and debounce resolvers across the structural components
refactor: rename the shared refreshQuery files and helper to refreshRequest so the name states the request they build
refactor: rename the DataTableAlignment type to TextAlignment and add HorizontalAlignment so the search box no longer borrows a text type
refactor: replace the CarouselPageSize and DataTablePageSize types with one shared ItemsPerPage type
refactor: name the shared pagination resolvers itemsPerPage and call them directly so the settings structs drop their identical delegating methods
refactor: move the border radius, shadow, ring, padding, gap, and background, text, and border color token resolvers into the toolset and call them from the structural and display components
refactor: name the carousel prev/next buttons arrows and group the items-per-view breakpoints in one settings struct
refactor: derive the carousel and data table id from a stable hash of their identity inputs so the id holds across refreshes, and keep the render counter for the tabs and the accordion
refactor: name the shared component id helpers by the id they build so their files state their purpose
refactor: keep the alignment types in one structural file and share the refresh fragment fetch between the carousel and the data table
refactor: gather the demo output generation as methods on one type and name the conditional write by the action it performs
refactor: add the shape, ring thickness, and compact ring token resolvers and align the border radius token with the class it names
refactor: move the shared hover tooltip into the toolset so the Button, the CarouselItemTooltip, and the InputHint render one mechanism
refactor: turn the select dropdown into an Alpine data component so the open-direction decision stops being a global UiToolset helper
refactor: attach the tooltip scroll and resize listeners only while the tooltip is visible
refactor: delay the tooltip reveal so a passing pointer does not flash it
refactor: render the tooltip through TooltipSurfaceSettings so the size caps travel with the surface
feat: refresh the carousel through htmx when it is present, with the pagination swapped out of band
chore: write the demo page and the refresh assets only when their content changed
fix: read the carousel swipe from the pointer events alone so touch input works
fix: pause the carousel autoplay on hover only when the setting is set
fix: swap the carousel track and pagination together so the totals stay current
fix: discard a stale carousel refresh and keep the filter bar outside the swapped region
fix: apply the lg border-radius token on the carousel surface and item wrapper
fix: keep the default outside carousel arrows on the light background
fix: position the inside carousel arrows with left and right classes instead of an inline style
fix: default the carousel border radius to the class its md token produces
fix: fail the carousel refresh when the response carries no pagination region
fix: read the settings script from the component root so two components with the same query template keep their own settings
fix: swap the carousel pagination only into the component that refreshed
fix: seed no filter in the demo carousel and describe the server example as paging only
test: cover the carousel resolvers, the shared style token and id prefix helpers, the rendered track, the touch swipe, the settings script placement, the pagination swap target, and the accessibility
docs: document the carousel in the structural readme and the feature map
chore(docs): regenerate the demo page

0.2.4 - 2026/10/01
feat: add the tabs component with horizontal and vertical tab headers
feat: add the badge count and the URL-hash syncing to the tabs
feat: add the tab side, the icon position, the surface and content customization, the badge styling, the size constraints, and the alignment to the tabs
test: cover the tabs expression builders and the rendered tab behavior
docs: document the tabs component in the structural readme and the feature map
chore(docs): regenerate the demo page

0.2.3 - 2026/09/30
feat: add the track tooltip to the range slider
feat: add the value bubble display mode to the range slider
feat: label the page strip from one without changing the bound state
feat: make the page number zero-based across pagination and the data table
fix: keep the button tooltip state from aborting on a second embed
fix: hold the cloak loading screen until Alpine finishes initializing
fix: keep the track tooltip preview on the step grid
feat: animate the modal size transition
feat: fade the loading overlay in and out
test: cover the track tooltip, the value bubble, the one-based labels, the modal resize, and the loading overlay markup
test: share the track tooltip hover helper across the slider specs
docs: document the zero-based page number in the component readmes
docs: add the favicon to the demo page
docs: cover the demo page from the first paint with the cloak loading screen
chore(docs): regenerate the demo page

0.2.2 - 2026/09/22
feat: add the button tooltip in the fixed body layer
fix: always render the button type attribute
refactor: render each embedded state script once per page
docs: bind the input field hint demos to Alpine state
docs: bind the affix width demo to Alpine state
docs: bind the sizes and label case demos to Alpine state
refactor: split the demo page into one templ file per component
fix: show the loading overlay while a request is in flight
test: assert the runner CLI contract and cap the test processes
test: assert rendered behavior across the component specs
docs: document the new components in the package readmes
feat: add the text case values and class resolver to the toolset
feat: add the text case setting to the components
feat: scale the input field with its affixes and floating label
feat: add the change callback to the multi select input
fix: center the hint icon on the field box
feat: add multi value support to the filter bar
feat: repeat multi value filters in the refresh query
feat: hide the pagination on a single page
feat: pass the filter and pagination options through the data table
feat: default the modal resize range to medium through near-full
feat: size the modal width and height from viewport percentages
refactor: split the modal resize into separate enlarge and reduce controls
refactor: rename the modal Size to InitialSize and AvailableSizes to PossibleSizes
refactor: move the demo section layout into one shared template
test: cover the sidebar dynamic classes builder and the accordion cutout
refactor: extract the sidebar dynamic classes builder into Go
fix: stop the sidebar float classes from reading an empty attachment path
fix: bind the sidebar demo collapse state to the shared sidebar content
feat: show the accordion item content as a cutout of the parent surface
docs: state the alert title and description raw-html contract
docs: separate the demo live-state readouts from their components
docs: stack the demo meta blocks and put the live example beside them
test: cover the card and page headings with structural specs
test: cover the confirmation dialog presets and type-to-confirm gate
test: cover the modal viewport percentage sizes
docs: show the page headings, card, confirmation dialog presets, and near-full modal in the demo
feat: add the delete and critical type-to-confirm gate to the confirmation dialogs
feat: add the confirm, warning, critical, and delete confirmation dialog presets
feat: make the modal size scale a coupled viewport percentage and add the xxl size
feat: add the card with header block, middle and footer slots, and surface settings
feat: add the page heading with page and section levels
docs: record the v0.2.2 feature entries in the changelog
test: raise the color contrast baseline to 84
test: register the form and display unit test nodes
test: rename the assertion messages to PascalCase
docs: show the checkbox error state and tag variants in the demo
feat: add the tag tiny size, bound inner value, and border-rendered rings
feat: add the checkbox input error state
docs: document the demo server start and stop for agent-browser
test: raise the select dropdown open golden baseline to 50ms
test: assert the stale refresh outcome without request counts
docs: record the pending entries in the changelog
refactor: namespace the button tooltip state methods
feat: add the modal viewport size constraints and rename IsHeightContentSized
feat: refine the confirmation dialog sizing and alignment
refactor: serve the data table refresh from one json asset
chore(docs): regenerate the demo page
fix: center the demo readout text in its box
fix: serve the demo files with a no-store cache header
test: hold the stale refresh in flight before superseding it
test: drive the data table fragment spec from component state
docs: record the pending entries in the changelog
chore: exclude the generated files from sonar
fix: derive the select options script id per instance
feat: add the loading overlay id setting
fix: render the radio input for its label
test: deduplicate the modal specs and drop the terminal content id
test: use double brackets in the shell assertion helpers

0.2.1 - 2026/09/18
fix: serve demo files only from canonical paths under the docs root
fix: label inputs and textareas even when an input id is set
test: wait for the data table refresh before asserting the error is hidden
fix: run the browser specs without npx
fix: reject demo requests that escape the docs root
fix: use a strong hash for the demo asset cache
fix: abort superseded data table refreshes
fix: close the modal on a backdrop press, not a panel drag
fix: lowercase the state-bound tag labels
fix: label form inputs and textareas for screen readers
fix: label the demo radios and drop the duplicate modal id
test: exclude structural specs from the form and cross-browser suites
docs: show the RegisterAlpineState callback in the toolset README
feat: render a default search box in the data table
refactor: rename the data table query url template setting
refactor: rename the data table column render callback to CellRenderer
refactor: rename the data table items per page size choices setting
refactor: derive the data table id from its configuration only
refactor: rename the filter settings type field to kind
refactor: move the filter bar expressions to an alpine state
feat: add an aria label setting to the button
chore: adopt biome for JavaScript, CSS, and SVG lint and formatting
docs: box the live example state readouts in the demo
feat: add data table row and column styling, search box alignment, and checkbox shape, size, and color settings
fix: keep the checkbox and toggle switch label at its content width
refactor: suffix the data table row callback fields with resolver
refactor: prefix the data table initial state settings with initial
docs: capitalize the data table advanced example titles
feat: add checkbox input with shapes, sizes, colors, and indeterminate state
fix: use the checkbox input for the data table row selection
fix: open the select input dropdown upward when there is no room below
fix: compact the pagination buttons and keep extra small icon buttons at 24px
feat: add a data table header text case setting and lowercase the sortable headers
fix: apply the select input size setting and compact the items per page control
docs: order the demo components alphabetically and label the sidebar sections
refactor: return errors from the demo generator and name its artifacts
fix: run the modal close callback on a backdrop click
fix: stop the backdrop click from closing an uncloseable modal
docs: split the README into one README per package
docs: state the data table server requirement in the demo and README
fix: serve the requested page in the data table demo
fix: show page numbers with the current page marked in pagination
fix: scale the tag remove icon and show one sort icon per table header
fix: lowercase buttons, labels, and table headers and shorten form controls
test: exercise table filters, search, and page numbers end to end
docs: describe the structural components in the context files
feat: add data table with server-driven refresh
feat: add filter bar with removable chips
feat: add pagination
feat: add removable tag variant
refactor: move sidebar to the structural package

0.2.0 - 2026/09/16
fix: restrict the demo asset proxy to approved hosts
refactor: remove the unused selectInput InputId setting
feat: normalize range slider initial two-way state
test: bound demo proxy retries and survive asset fetch failures
docs: describe the demo asset proxy in the context files
test: serve demo assets from a local cache to stop CDN stalls
fix: sync range slider values and make select and text area controls accessible
test: add range slider drag specs
fix: make range slider thumbs draggable
test: add hint accessibility specs
fix: make form hint tooltips keyboard and touch accessible
feat: add hints to TextArea
docs: add hint examples to the demo
fix: escape selectInput option values and run OnChangeFunc on change
test: add range slider control specs and register the control suite
feat: add range slider tick marks and accessible thumb names

0.1.9 - 2026/09/15
fix: set the demo document language to English
fix: reject invalid inputs and count concurrent requests in jsonAjax
test: add jsonAjax specs and start the demo server only when a suite needs it
test: fail the performance check when a golden budget has no result
feat: add label position setting to ToggleSwitch
fix: drop top margin that misaligned InlineRadioGroup with sibling form controls
fix: rescale TextArea heights to steps of 12 and soften the large size font
fix: align TextArea action icons with the text inset and widen the gutter so lines clear them
docs: replace goreportcard badge with sonar quality gate in README and demo
docs: bump demo sidebar version label to 0.1.9
chore: update go and deps
chore: pin tool versions in mise
feat: display liaison response messages in toast
docs: add consumer skill file
test: add tests.sh suite with Playwright behavioral specs for form components and axe ratchet
test: add golden performance budgets with click-to-final-state latency specs

0.1.8 - 2026/08/04
feat: add ToggleSwitch with configurable styles and Alpine bindings
docs: expand ToggleSwitch and component usage examples
fix: scope hint tooltips and accordion groups
fix: preserve separate FormData values for array-bound ToggleSwitch
fix: serialize MultiSelectInput demo state as valid JavaScript

0.1.7 - 2026/07/27
chore: bump go and templ deps
chore: update templ generated files to v0.3.1020
feat: selectInput hint display modes (tooltip and description)
feat: multiSelectInput
feat: toast AutoDismissSeconds setting (defaults to 10s, was 4s)
docs: add development build instructions

0.1.6 - 2026/03/16
fix: ensure createRandomPassword meets password VO requirements
chore: update go and deps
chore: update templ to v0.3.1001

0.1.5 - 2025/06/07
feat: max width setting for alert
fix: allow html in alert title and description
chore: update templ to v0.3.898

0.1.4 - 2025/06/07
feat: alert
feat: button IsVisibleOneWayStatePath setting

0.1.3 - 2025/06/05
feat: modal resize and close functions

0.1.2 - 2025/06/05
feat: javascript toolset
feat: esbuild-based minifier

0.1.1 - 2025/06/05
fix: increase bg opacity on button hover
fix: invert closeable and resizable logic on modal

0.1.0 - 2025/06/04
feat: add isResizable to modal
fix: modal padding
chore: update templ to v0.3.894

0.0.9 - 2025/06/03
feat: modal
feat: alpine state for button icons
chore: update templ to v0.3.887

0.0.8 - 2025/05/27
feat: add cloak loading
fix: add suffix seconds to loading overlay duration

0.0.7 - 2025/05/27
feat: add vega to import
feat: add alpine intersect to import
fix: add missing z-index to toast

0.0.6 - 2025/05/27
feat: loading overlay

0.0.5 - 2025/05/26
feat: toast
fix: show sidebar collapse button only on hover

0.0.4 - 2025/05/23
feat: add wrapper to sidebar
fix: sidebar scroll bar width

0.0.3 - 2025/05/23
feat: sidebar
fix: textarea and input field text color replacement
fix: hide focus outline on input and textarea on Firefox

0.0.2 - 2025/05/17
feat: add mozilla fonts
chore: update deps

0.0.1 - 2025/05/09
feat: initial release
```

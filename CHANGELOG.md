# Changelog

```log
0.2.0 - 2026/09/16
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

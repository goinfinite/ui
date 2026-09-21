---
name: ui-usage
description: Use when working in a Go project that imports github.com/goinfinite/ui — presents the component packages, routes you to the right source file, and states the usage conventions.
version: 1.1.0
lastUpdated: 2026-09-18
---

## Purpose

Infinite UI _(UI)_ is the shared component library for Go projects in the Infinite ecosystem. It provides templ components, head-tag imports, and JS/Go toolsets. The components build on Alpine.js, Tailwind-compatible CSS, and Phosphor Icons. The UnoCSS runtime serves the CSS. The components are tested and shared across projects. Reusing them beats hand-rolling markup. Before you write an input, a modal, or a toast by hand, check whether UI already provides it. This skill maps the packages and points you to the source file that documents each component.

## Procedure

### 1. Locate the installed module

Run this command from the project root:

```sh
go list -m -f '{{.Dir}}' github.com/goinfinite/ui
```

The output is the module directory inside the Go module cache. Read the source there. Do not copy files into the project. Do not edit the cache.

### 2. Pick the package and open its components

Each package holds one kind of component. Pick the package that matches the task:

- Collecting user input → `src/form/`
- Triggering an action or picking a range → `src/control/`
- Showing content, state, or chrome → `src/display/`
- Displaying and paging structured data → `src/structural/`
- Bootstrapping the page (CDN head tags and the JS toolset) → `src/import/`
- Minifying JS/CSS at build or render time → `src/toolset/`

The package README owns the component list and the usage snippets. Read it first, then only the `.templ` files that apply. Each component file defines the component and its settings struct together.

### 3. Read deeper only when needed

These references describe UI's internals. Consult them only when the task reaches the matching subsystem:

- `docs/FEATURE-MAP.md` — end-to-end flows for a UI feature, including which JS state file backs each component.
- `demo/demo.templ` — every component rendered with real settings. It is the fastest way to see a usage example.
- `https://ui.demo.goinfinite.net/` — the same demo running live. Open it with your browser tool when you have one. It shows the rendered behavior that source files hide: animations, dropdown state, and styling.
- `.context.md` in the package you import — constraints on that package's files.
- `README.md` in the module root — the library overview, installation, and conventions.

The module directory matches the version in `go.mod`, so its documentation matches the API you compile against.

### 4. Bootstrap the page

Every component assumes the page loads Alpine.js, the CSS engine, and Phosphor Icons. Include the head tags in your layout's `<head>`:

- `@uiImport.HeadTagsMinimal()` — Alpine, UnoCSS runtime, Phosphor Icons.
- `@uiImport.HeadTagsFull()` — adds Google Fonts, HTMX, and the JS toolset.

Granular variants exist (`HeadTagsAlpineJs()`, `HeadTagsHtmx()`, `HeadTagsToolset()`, ...) for selective inclusion. The UnoCSS runtime serves Tailwind. Put custom `primary` and `secondary` colors in the `window.__unocss` theme config on the page. Set that config in a plain `<script>` before the head tags, so the runtime reads it at startup. A `tailwind.config.js` file has no effect on the runtime.

### 5. Follow the usage conventions

- Pass one settings struct per component. Required fields sit at the top, optional fields below. Example:

  ```go
  @uiForm.InputField(uiForm.InputFieldSettings{
      InputType:  uiForm.InputTypeText,
      InputName:  "name",
      Label:      "Name",
      IsRequired: true,
  })
  ```

- Bind state with `TwoWayStatePath` (`x-model`) or `OneWayStatePath` (`x-bind`). The path points into the nearest `x-data` object: `"user.name"` for nested objects, `"users[0]"` for arrays.
- Call functions with a `*Func` field. The value is a full call with parentheses, for example `OnClickFunc: "myFunction('hi')"`. The function can live on the parent `x-data` object, `window`, or `document`.
- Set `InputName` on every form input. HTMX submissions build `FormData` from the `name` attribute. `InputId` only locates the element in the DOM.
- Import the packages plainly. The package names are already prefixed: `uiForm`, `uiControl`, `uiDisplay`, `uiImport`, `uiStructural`, `uiToolset`. No import alias is needed.
- Ship Alpine state as a sibling `.js` file. Embed it with `go:embed`, render it with `@uiToolset.MinifierTemplateJs(&stateVar)`, and register it inside `UiToolset.RegisterAlpineState(() => { Alpine.data(...) })` in the JS file.

## Guardrails

- Never hand-roll a component UI already provides. Extend or wrap the existing component instead.
- Never copy UI code into the project. Import it.
- Never edit files inside the Go module cache. The cache is read-only input.
- Do not assume a component is missing because it is not listed here. Read the package directory in the module cache first.
- Do not add a JavaScript dependency for behavior the JS toolset already offers (`ToggleLoadingOverlay`, `JsonAjax`, `ResolveApiResponseDisplay`, `CreateRandomPassword`).
- Do not call `Alpine.data()` or `Alpine.store()` at file scope. Register them inside `UiToolset.RegisterAlpineState()` to avoid duplicate registration when navigating between pages.
- Do not place `CloakLoading` anywhere but the end of `<body>`. It uses inline styles, so it renders before CSS loads.
- Do not fight Toast's automatic behavior. With a Toast on the page, HTMX JSON responses carrying `readableMessage` trigger toasts on their own.

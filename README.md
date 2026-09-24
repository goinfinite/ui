# [Infinite UI](https://github.com/goinfinite/ui) &middot; [![Demo](https://img.shields.io/badge/demo-233876)](https://ui.demo.goinfinite.net/) [![/r/goinfinite](https://img.shields.io/badge/%2Fr%2Fgoinfinite-FF4500?logo=reddit&logoColor=ffffff)](https://www.reddit.com/r/goinfinite/) [![Discussions](https://img.shields.io/badge/discussions-751A3D?logo=github)](https://github.com/orgs/goinfinite/discussions) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goinfinite_ui&metric=alert_status)](https://sonarcloud.io/project/overview?id=goinfinite_ui) [![License](https://img.shields.io/badge/license-MIT-teal.svg)](https://github.com/goinfinite/ui/blob/main/LICENSE.md)

Infinite UI is a collection of reusable components for building elegant user interfaces in Go with [a-h/templ](https://github.com/a-h/templ), [Alpine.js](https://github.com/alpinejs/alpine), [Tailwind CSS](https://tailwindcss.com) via [UnoCSS](https://unocss.dev/integrations/runtime), and [Phosphor Icons](https://phosphoricons.com/). It standardizes the common cases and leaves the custom cases open.

> [!TIP]
> **Working with an AI agent?** Point it to [`SKILL.md`](SKILL.md) before it writes code that imports Infinite UI. The skill routes the agent to the right package and its README.

> [!IMPORTANT]
> **Human Reviewed**: AI models assist development, but senior developers review every line for coherence, readability, and maintainability.

## Installation

Infinite UI requires Go 1.27.1 or later. Install it with:

```bash
go get github.com/goinfinite/ui@latest
```

Your `<head>` must load Alpine.js, a Tailwind-compatible CSS engine, and Phosphor Icons. Use `@uiImport.HeadTagsMinimal()` or `@uiImport.HeadTagsFull()`; the full variant adds Google Fonts, HTMX, and the JavaScript toolset. Granular variants exist for each subset.

Put the custom `primary` and `secondary` colors in the `window.__unocss` theme config on the page, in a plain `<script>` before the head tags. A `tailwind.config.js` file has no effect on the UnoCSS runtime.

If you have not installed the template engine yet:

```bash
go get github.com/a-h/templ
```

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Usage

Every component accepts one settings struct. Required fields come first, then a comment line, then the optional fields. This avoids pointers, which templ does not support well.

```go
@uiForm.InputField(uiForm.InputFieldSettings{
    InputType:  uiForm.InputTypeText,
    InputName:  "name",
    Label:      "Name",
    IsRequired: true,
})
```

Components connect to Alpine.js state through their `*StatePath` fields:

- A `TwoWayStatePath` field binds with `x-model`.
- A `OneWayStatePath` field binds with `x-bind`.

The path names a property in the nearest `x-data` object. Use dot notation for nested objects and bracket notation for arrays.

Components run functions on events through their `On*Func` fields. The value is a full call, for example `OnClickFunc: "myFunction('Hello World!')"` or `OnClickFunc: "myFunction()"`. The function can live in the parent `x-data` object, on `window`, or on `document`.

Components use a neutral color palette with reduced opacity values, for example `bg-neutral-50/10`. This keeps them compatible with any design system that provides `primary` and `secondary` colors.

HTMX is optional. Components work without it; the data table falls back to `fetch`.

## Packages

Each package has its own README with the component list and usage snippets:

- **[Form](src/form/README.md)** — text, choice, and switch inputs.
- **[Control](src/control/README.md)** — buttons and sliders.
- **[Display](src/display/README.md)** — alerts, confirmation dialogs, modals, header blocks, tags, toasts, and page chrome.
- **[Structural](src/structural/README.md)** — the card, page headings, server-driven data table, filter bar, pagination, and sidebar.
- **[Toolset](src/toolset/README.md)** — the JavaScript `UiToolset` and the Go minifier.

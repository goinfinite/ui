# [Infinite UI](https://github.com/goinfinite/ui) &middot; [![Demo](https://img.shields.io/badge/demo-233876)](https://ui.demo.goinfinite.net/) [![/r/goinfinite](https://img.shields.io/badge/%2Fr%2Fgoinfinite-FF4500?logo=reddit&logoColor=ffffff)](https://www.reddit.com/r/goinfinite/) [![Discussions](https://img.shields.io/badge/discussions-751A3D?logo=github)](https://github.com/orgs/goinfinite/discussions) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=goinfinite_ui&metric=alert_status)](https://sonarcloud.io/project/overview?id=goinfinite_ui) [![License](https://img.shields.io/badge/license-MIT-teal.svg)](https://github.com/goinfinite/ui/blob/main/LICENSE.md)

Infinite UI is a set of reusable Go components. It uses [a-h/templ](https://github.com/a-h/templ), [Alpine.js](https://github.com/alpinejs/alpine), [Tailwind CSS](https://tailwindcss.com) through [UnoCSS](https://unocss.dev/integrations/runtime), and [Phosphor Icons](https://phosphoricons.com/). The components cover the common cases and leave the custom cases open.

> [!TIP]
> Working with an AI agent? Point it to [`SKILL.md`](SKILL.md) before it writes code that imports Infinite UI. The skill sends the agent to the correct package README.

> [!IMPORTANT]
> **Human Reviewed**: AI models assist development, but senior developers review every line for coherence, readability, and maintainability.

## Installation

Infinite UI requires Go 1.27.1 or later:

```bash
go get github.com/goinfinite/ui@latest
go get github.com/a-h/templ
```

The page `<head>` must load Alpine.js, a Tailwind-compatible CSS engine, and Phosphor Icons. Use `@uiImport.HeadTagsMinimal()` or `@uiImport.HeadTagsFull()`. The full variant adds Google Fonts, HTMX, and the JavaScript toolset. A granular variant exists for each library. `Chart` needs Vega once per page. Load it with `@uiImport.HeadTagsVega()`, usually in the layout head.

Set the `primary` and `secondary` colors in `window.__unocss`, in a plain `<script>` before the head tags. The UnoCSS runtime ignores `tailwind.config.js`.

[CHANGELOG.md](CHANGELOG.md) holds the release history.

## Usage

Each component takes one settings struct. The required fields come first. The optional fields follow the `// OptionalFields` line. The struct avoids pointers, because templ handles them badly.

```go
@uiForm.InputField(uiForm.InputFieldSettings{
    InputType:  uiForm.InputTypeText,
    InputName:  "name",
    Label:      "Name",
    IsRequired: true,
})
```

The `*StatePath` fields connect a component to Alpine.js state. A `TwoWayStatePath` binds with `x-model`. A `OneWayStatePath` binds with `x-bind`. The path names a property of the nearest `x-data` object, with dots for nested objects and brackets for arrays.

The `On*Func` fields run a function on an event. The value is a full call, for example `OnClickFunc: "save('draft')"`. The function can live on the parent `x-data` object, on `window`, or on `document`.

The components use neutral colors with opacity, for example `bg-neutral-50/10`. They fit any design system that sets `primary` and `secondary`.

HTMX is optional. Without it, the data table uses `fetch`.

## Packages

Each package README lists its components with usage snippets:

- [Form](src/form/README.md): text, choice, number range, file upload, and switch inputs.
- [Control](src/control/README.md): buttons with action states, and sliders.
- [Display](src/display/README.md): alerts, charts, dialogs, modals, header blocks, tags, toasts, and loading screens.
- [Structural](src/structural/README.md): card, page heading, data table, carousel, filter bar, pagination, sidebar, and tabs.
- [Toolset](src/toolset/README.md): the JavaScript `UiToolset` and the Go helpers.

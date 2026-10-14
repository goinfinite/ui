# Toolset

The toolset package holds the JavaScript `UiToolset` object, a Go minifier, and the shared Go text and style helpers. Part of [Infinite UI](../../README.md).

## Style tokens

Resolvers that turn a token into a Tailwind class. Each resolver takes a fallback class for an empty or unknown token. The fixed-value constants (`BorderRadius*`, `Shape*`, `ShadowSize*`, `RingThickness*`, `PaddingSize*`, `GapSize*`, and `Animation*`) live here. Color tokens are free-form strings.

- `BackgroundColorClassResolver(backgroundColor, fallbackClass)` returns `bg-` plus the token.
- `TextColorClassResolver(textColor, fallbackClass)` returns `text-` plus the token.
- `BorderColorClassResolver(borderColor, fallbackClass)` returns `border-` plus the token.
- `BorderRadiusTokenResolver(borderRadius, fallbackToken)` returns the token when it is a `BorderRadius*` value.
- `BorderRadiusClassResolver(borderRadius, fallbackClass)` returns the class with the token name, for example `rounded-md` for `md`.
- `ShapeClassResolver(shape, fallbackClass)` returns `rounded-full`, `rounded`, or `rounded-none` for `ShapeCircular`, `ShapeRounded`, or `ShapeSquare`.
- `ShadowClassResolver(shadowSize, fallbackClass)` returns the `shadow-*` class for `none` through `xl`.
- `RingThicknessClassResolver(ringThickness, fallbackClass)` returns the `ring-*` class for `xs` through `xl`.
- `RingClassResolver(ringColor, ringThickness)` returns the thickness and color classes. It returns an empty string when either value is missing.
- `CompactRingClassResolver(ringColor, ringThickness)` uses the thinner tooltip scale, `ring-0.5` through `ring-2.5`. An empty thickness gives `ring-1`. An empty color gives an empty string.
- `PaddingClassResolver(paddingSize, fallbackClass)` returns the `p-*` class for `none` through `xl`.
- `CompactPaddingClassResolver(paddingSize, fallbackClass)` uses a smaller `p-*` scale for chips and list items.
- `GapClassResolver(gapSize, fallbackClass)` returns the `gap-*` class for `none` through `xl`.
- `AnimationClassResolver(animation, fallbackClass)` returns the `animate-*` class for `AnimationSpin`, `AnimationPulse`, `AnimationBounce`, or `AnimationPing`. `AnimationNone` returns an empty string, an empty token returns the fallback, and any other value passes through as a custom class.
- `NonEmptyClassJoiner(classList)` joins the non-empty entries with one space.

## Component id

Helpers that build a component DOM id. With a unique id, one page can hold several copies of a component without duplicate element ids, radio group names, or input names.

- `NewComponentIdPrefixGenerator(componentName)` creates a generator. `GenerateNext()` returns `<componentName>-<n>`, with `n` one higher on each call. Use it for components that render once, such as tabs and the accordion.
- `HashComponentIdParts(idParts...)` returns a stable hash of the parts. A component that refreshes uses it, so its id stays the same on the page render and on each refresh.

## Text case

Text case values and their class resolver.

- `TextCaseNone` keeps the text as typed (the default).
- `TextCaseLower`, `TextCaseUpper`, and `TextCaseCapitalize` change the case.
- `TextCaseClassResolver(textCase)` returns the Tailwind class. It returns an empty string for `TextCaseNone` and unknown values.

## Tooltip

Shared parts for hover tooltips. The component sets `x-data` to the `tooltip` Alpine component, puts `x-ref="trigger"` on the trigger, and renders the tooltip in `<template x-teleport="body">`. The tooltip renders on the document body, so a scroll container, a moving track, or a modal cannot clip it.

- `TooltipAlpineStateOnce` embeds `tooltipState.js` once. Render `@uiToolset.TooltipAlpineStateOnce.Once()` above each component that opens a tooltip.
- `TooltipSurfaceClassesResolver(TooltipSurfaceSettings{...})` returns the tooltip classes. Empty colors give `bg-neutral-800/95` and `text-neutral-50`. The ring uses the compact scale. `MinWidthClass`, `MaxWidthClass`, `MinHeightClass`, and `MaxHeightClass` limit the size. `MaxWidthClass` defaults to `max-w-96`, so long text wraps.
- `TooltipPositionTop`, `TooltipPositionBottom`, `TooltipPositionLeft`, and `TooltipPositionRight` set the placement.
- The Alpine component takes a position (empty means top). The trigger calls `showTooltip()` and `hideTooltip()` on its mouse and focus events.
- `showTooltip()` waits 200ms, so a pointer that crosses the trigger does not flash the tooltip. `hideTooltip()` cancels a waiting reveal.
- Keep each declaration of `tooltipState.js` inside the `RegisterAlpineState` callback. A second embed must not declare a top-level `const` again, because that stops the script with a SyntaxError.

## JavaScript toolset

The JavaScript files live in `src/import/toolset/`. Include them with `@uiImport.HeadTagsFullJs()` or `@uiImport.HeadTagsToolset()`. Call the utilities through the `UiToolset` object.

- `UiToolset.CreateRandomPassword(options)` returns a random password. `options` takes `length` (default 16), `minLength` (default 6), `maxLength` (default 64), and one `include*` flag for each character class.
- `UiToolset.ResolveApiResponseDisplay(apiResponse, httpStatusCode)` returns the message and the outcome (`success`, `partialSuccess`, or `error`) of an Infinite API response.
- `UiToolset.ToggleLoadingOverlay()` shows or hides the element with the id `loading-overlay`.
- `UiToolset.JsonAjax(method, url, payload, toast)` sends a JSON request. It shows the loading overlay and can show a toast from the response. It returns a promise that rejects on failure.
- `UiToolset.RegisterAlpineState(stateFunction)` runs `stateFunction` on `alpine:init`, or at once when Alpine already started. Call `Alpine.data()` inside it, so a page change does not register the same state twice.

```js
UiToolset.RegisterAlpineState(() => {
  Alpine.data("dataTable", (settingsScriptId) => ({ ... }));
});
```

With the Toast component on the page, an API response shows its message. Use this envelope:

```json
{
    "status": 200,
    "readableMessage": "Operation completed successfully.",
    "body": {}
}
```

- `readableMessage` is a string. `body` can hold any JSON value. The HTTP status sets the toast tone.
- A `humanReadableMessage` object with `error`, `partialSuccess`, and `success` fields also works.
- The toast uses a non-empty `readableMessage` first, then the matching `humanReadableMessage` field, then a string `body`. With no message, the toast stays hidden.

## Go minifier

The minifier uses [esbuild](https://github.com/evanw/esbuild) to shrink JavaScript and CSS before the page gets them. `MinifierSettings` changes the defaults when they break your code.

- `Minifier(MinifierSettings)` minifies JavaScript or CSS and returns a string pointer.
- `MinifierJs(unminifiedContent)` minifies JavaScript. `MinifierCss(unminifiedContent)` minifies CSS.
- `MinifierTemplate(MinifierSettings)` returns a `templ.Component`.
- `MinifierTemplateJs(unminifiedContent)` returns a `<script>` tag. `MinifierTemplateCss(unminifiedContent)` returns a `<style>` tag.

Embed the source file and render it minified:

```go
//go:embed accountsState.js
var accountsAlpineState string
```

```templ
@uiToolset.MinifierTemplateJs(&accountsAlpineState)
@uiToolset.MinifierTemplateCss(&accountsCustomCss)
```

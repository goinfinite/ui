# Toolset

Toolset of Infinite UI. It provides a JavaScript utility object and a Go minifier. Part of [Infinite UI](../../README.md).

## JavaScript toolset

The JavaScript files live in `src/import/toolset/`. Include them with `@uiImport.HeadTagsFullJs()` or `@uiImport.HeadTagsToolset()`.

Call the utilities through the `UiToolset` object.

- `UiToolset.CreateRandomPassword()`: Creates a random password of length 16 characters.
- `UiToolset.ResolveApiResponseDisplay(apiResponse, httpStatusCode)`: Resolves the message and outcome (`success`, `partialSuccess`, or `error`) from an Infinite API response envelope.
- `UiToolset.ToggleLoadingOverlay()`: Toggles the loading overlay element with the id `loading-overlay`.
- `UiToolset.JsonAjax(method, url, payload, toast)`: Makes a JSON AJAX request. It shows the loading overlay and can display a toast from the response envelope.
- `UiToolset.RegisterAlpineState(stateFunction)`: Registers `stateFunction` to run on `alpine:init`, or immediately when Alpine is already initialized. Call `Alpine.data()` inside the callback. Use it instead of raw `Alpine.data()` calls so navigation does not register the same listener twice.

```js
UiToolset.RegisterAlpineState(() => {
  Alpine.data("dataTable", (settingsScriptId) => ({ ... }));
});
```

When the Toast component is present, API responses can display messages with the recommended envelope:

```json
{
    "status": 200,
    "readableMessage": "Operation completed successfully.",
    "body": {}
}
```

`readableMessage` is a string. `body` can hold any JSON value. Toast styling follows the HTTP status. A `humanReadableMessage` object with `error`, `partialSuccess`, and `success` fields is also supported. `readableMessage` takes precedence when it is not empty. Otherwise the matching `humanReadableMessage` field is used. A string `body` is the last fallback. The toast stays hidden when the response carries no message.

## Go minifier

Based on [esbuild](https://github.com/evanw/esbuild), the minifier shrinks JavaScript and CSS before the page receives them. The `MinifierSettings` struct overrides the defaults when they break your code.

- `Minifier(MinifierSettings)`: Minifies JavaScript or CSS and returns a string pointer.
- `MinifierJs(unminifiedContent)`: Minifies JavaScript.
- `MinifierCss(unminifiedContent)`: Minifies CSS.
- `MinifierTemplate(MinifierSettings)`: Returns a `templ.Component` instead of a string.
- `MinifierTemplateJs(unminifiedContent)`: Wraps the result in a `<script>` tag.
- `MinifierTemplateCss(unminifiedContent)`: Wraps the result in a `<style>` tag.

Embed the source file and inject it minified:

```go
//go:embed accountsState.js
var accountsAlpineState string
```

```templ
@uiToolset.MinifierTemplateJs(&accountsAlpineState)
@uiToolset.MinifierTemplateCss(&accountsCustomCss)
```

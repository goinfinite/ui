# [Infinite UI](https://github.com/goinfinite/ui) &middot; [![Demo](https://img.shields.io/badge/demo-233876)](https://ui.demo.goinfinite.net/) [![/r/goinfinite](https://img.shields.io/badge/%2Fr%2Fgoinfinite-FF4500?logo=reddit&logoColor=ffffff)](https://www.reddit.com/r/goinfinite/) [![Discussions](https://img.shields.io/badge/discussions-751A3D?logo=github)](https://github.com/orgs/goinfinite/discussions) [![Report Card](https://img.shields.io/badge/report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/goinfinite/ui) [![License](https://img.shields.io/badge/license-MIT-teal.svg)](https://github.com/goinfinite/ui/blob/main/LICENSE.md)

Infinite UI is a collection of reusable components for building elegant user interfaces in Go with [a-h/templ](https://github.com/a-h/templ), [Alpine.js](https://github.com/alpinejs/alpine), [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss), [Phosphor Icons](https://phosphoricons.com/) and the occasional additional JavaScript libraries when necessary.

The library is engineered for developer efficiency and ease of implementation, offering a balance between standardization and customization to accelerate development workflows while maintaining high-quality, responsive interfaces.

## Features

- **Reusable Components**: A set of pre-built components that can be easily integrated into your projects.
- **Responsive Design**: Components are designed to be responsive and work well on various screen sizes.
- **Lightweight**: Built with performance in mind, ensuring fast load times and smooth interactions.
- **HTMX Ready**: Although not required, Infinite UI is designed to work well with HTMX.

## Installation

To use Infinite UI in your project, you can install it using Go modules. Run the following command in your terminal:

```bash
go get github.com/goinfinite/ui@latest
```

After installing Infinite UI, make sure your `<head>` component is including [Alpine.js](https://alpinejs.dev/essentials/installation), [Tailwind CSS](https://unocss.dev/integrations/runtime) and [Phosphor Icons](https://github.com/phosphor-icons/homepage?tab=readme-ov-file#vanilla-web).

One way to install the necessary scripts and styles is to use the `@uiImport.HeadTagsMinimal()` or `@uiImport.HeadTagsFull()` components. The only difference between the two is that `HeadTagsFull()` also includes [Google Fonts](https://fonts.google.com/) and [HTMX](https://htmx.org/).

If you haven't yet installed the template engine `a-h/templ`, you can do so by running the following command:

```bash
go get github.com/a-h/templ
```

## Usage

Infinite UI components are designed with a consistent interface pattern. Each component accepts a settings struct that contains all configuration options.

The settings structs are organized with required fields at the top, followed by optional fields (separated by a comment line). This approach eliminates the need for pointers, which aren't well-supported in templ.

Component settings may accept constant values from predefined types or simple string values depending on the component's requirements.

### Usage Example

```go
@uiForm.InputField(uiForm.InputFieldSettings{
    InputType:  uiForm.InputTypeText,
    InputName:  "name",
    Label:      "Name",
    IsRequired: true,
})
```

### Styling

Infinite UI components are designed with a neutral color palette using reduced opacity values (e.g., `bg-neutral-50/10`). This approach ensures compatibility across various design systems.

For hover and other states, please configure your Tailwind CSS with custom `primary` and `secondary` color schemes:

- Follow the [Tailwind CSS color customization guide](https://tailwindcss.com/docs/colors#customizing-your-colors);
- If you're using alternative atomic CSS engines like UnoCSS, please refer to their specific configuration documentation.

### Alpine.js States

Components leverage Alpine.js for state management by offering:

- `TwoWayStatePath`: Path for bidirectional data binding using `x-model` directive;
- `OneWayStatePath`: Path for read-only data binding using `x-bind` directive;

These common settings struct fields allow you to connect component states to your application's data model seamlessly. The path is a string representing the x-data object property to bind to.

For example, `x-data="{ name: 'John' }"` would use `OneWayStatePath: "name"` or `TwoWayStatePath: "name"`.

- For nested objects, use dot notation: `x-data="{ user: { name: 'John' } }"` would use `OneWayStatePath: "user.name"` or `TwoWayStatePath: "user.name"`;
- For arrays, use bracket notation: `x-data="{ users: ['John', 'Jane'] }"` would use `OneWayStatePath: "users[0]"` or `TwoWayStatePath: "users[0]"`;

### Function Calling

Components can execute functions when specific events occur. These functions can be defined within the parent component's x-data object or broader scopes like `window` or `document`.

Two patterns are available for function invocation, depending on component requirements:

1. **`Func` suffix fields** (e.g., `OnClickFunc`, `OnChangeFunc`):

- Functions must include parameters or empty parentheses;
- Example: `OnClickFunc: "myFunction('Hello World!')"` or `OnClickFunc: "myFunction()"`;
- Suitable when you need to specify exact parameters during event triggers.

2. **`FuncName` suffix fields** (e.g., `OnClickFuncName`, `OnChangeFuncName`):

- Specify only the function name without parameters;
- Example: `OnClickFuncName: "myFunction"`;
- The component will automatically pass appropriate parameters;
- Refer to component documentation for parameter requirements.

The `On` prefix (e.g., `OnClick`, `OnChange`) indicates event-driven execution, though function calling is not limited to events in all components.

### HTMX Ready

Infinite UI does not require HTMX to function, but it is designed to work well with it. HTMX is a JavaScript library that allows you to make AJAX requests and update parts of your page without reloading the entire page.

Form submissions done with HTMX usually utilize the `FormData` object to send form data to the server. Infinite UI allows you to set the `InputName` field to specify the keys of the `FormData` object. There is also `InputId` but `id` attributes are not required for form submissions, it's for locating the element in the DOM if necessary.

### JavaScript Toolset

The `toolset.js` file is a collection of JavaScript utility functions that can be used in your project. It is included in the `@uiImport.HeadTagsFullJs()` component or you can include it manually using `@uiImport.ToolsetJs()`.

To use the toolset, you can access it through the `UiToolset` object or `window.UiToolset` object.

- `UiToolset.CreateRandomPassword()`: Creates a random password of length 16 characters.
- `UiToolset.ToggleLoadingOverlay()`: Toggles the loading overlay element with the id `loading-overlay`.
- `UiToolset.JsonAjax()`: Makes a JSON AJAX request to the specified URL.¹
- `UiToolset.RegisterAlpineState()`: Registers a function to be called when Alpine.js is initialized or when it is already initialized. Useful for avoiding registering repeated addEventListeners for the same method, like when transiting through pages.

_¹Available when Alpine.js was already initialized._

### Go(lang) Toolset

Besides the JavaScript toolset, Infinite UI also provides a Go(lang) toolset to address common UI related tasks.

#### Minifier

Based on [esbuild](https://github.com/evanw/esbuild), the minifier is a tool to minify JavaScript and CSS content. Use the MinifierSettings struct if the defaults are breaking your code.

- `Minifier(MinifierSettings)`: Minifies JavaScript and CSS content.
- `MinifierJs(unminifiedContent)`: Minifier but takes a string as input and returns a string.
- `MinifierCss(unminifiedContent)`: MinifierJs but for CSS content.
- `MinifierTemplate(MinifierSettings)`: Minifier but returns a templ.Component instead of a string.
- `MinifierJsTemplate(unminifiedContent)`: MinifierJs but returns a templ.Component.
- `MinifierCssTemplate(unminifiedContent)`: MinifierCss but returns a templ.Component.

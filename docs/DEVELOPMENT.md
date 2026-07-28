# Development

## Building the Demo

The demo is a static HTML file generated from the `.templ` sources. To regenerate it, run the following command from the project root:

```bash
templ generate && cd docs; go run ../demo/*.go; cd ..
```

This will:

1. Run `templ generate` to compile every `.templ` file in the repository into its `_templ.go` counterpart.
2. Switch to the `docs/` directory.
3. Run the `demo` package, which renders the demo page (`DemoIndex()`) and writes it to `docs/index.html`.
4. Switch back to the project root.

The generated `docs/index.html` is the file served at [ui.demo.goinfinite.net](https://ui.demo.goinfinite.net/).

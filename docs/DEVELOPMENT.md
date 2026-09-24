# Development

## Building the Demo

The demo is a static HTML file generated from the `.templ` sources. To regenerate it, run the following command from the project root:

```bash
templ generate && (cd docs && go run ../demo/*.go)
```

This will:

1. Run `templ generate` to compile every `.templ` file in the repository into its `_templ.go` counterpart.
2. Switch to the `docs/` directory.
3. Run the `demo` package, which renders the demo page (`DemoIndex()`) and writes it to `docs/index.html`.
4. Switch back to the project root.

The generated `docs/index.html` is the file served at [ui.demo.goinfinite.net](https://ui.demo.goinfinite.net/).

## Publishing the Demo

GitHub Pages serves the demo. The Pages source is the `docs/` folder of the main branch — a GitHub Pages convention that keeps the published site out of the repository root. `docs/CNAME` binds the custom domain `ui.demo.goinfinite.net`; Pages reads that file automatically.

Publishing is plain git: commit the regenerated `docs/index.html` and push to main. Pages rebuilds and serves the new file within about a minute.

## Running the Tests

The entry point is `tests/tests.sh`. It serves `docs/` on a local port, runs the selected suites, and tears the server down.

```bash
./tests/tests.sh                       # standard level, all features
./tests/tests.sh --level=fast          # quick feedback while editing
./tests/tests.sh --level=exhaustive    # cross-browser, before PR or release
./tests/tests.sh list                  # show the registered test tree
```

Exit codes: 0 pass, 1 test failure, 2 runner error. Levels are cumulative.

UI tests are Playwright specs under `tests/ui/specs/`, organized by feature (`form`, `display`, `control`, `structural`, `toolset`, `a11y`, `performance`). They assert rendered behavior — geometry, visibility, computed styles, and rendered text — against the demo page. Install dependencies with `npm install` and browsers with `npx playwright install` inside `tests/`. WebKit is registered in `playwright.config.js` but disabled in `run.sh` because this host lacks its system libraries; enable it in CI containers that provide them.

The accessibility scope runs axe-core against the demo and fails on regressions past the counts recorded in `tests/ui/a11y-baseline.json`. Lower the baseline as components get fixed. Raise it only deliberately: `UPDATE_A11Y_BASELINE=1 ./tests/tests.sh --scope=accessibility`.

The performance scope measures click-to-final-state latency and navigation LCP in `tests/ui/specs/performance/`, writes `tests/performance-results.json`, and compares it against `tests/golden.yaml`. Tiers: `info` and `warn` print only; `fatal_ms` fails the run. New budgets ship without `fatal_ms` — promoting one is a deliberate human edit. The runner never rewrites `golden.yaml`.

New suites are registered in `tests/registry.yaml` with a feature, scope, and level. `tests.sh` executes registered `run` commands through `bash -c`, so treat `registry.yaml` as trusted input: only repo writers may change it, exactly like the scripts themselves.

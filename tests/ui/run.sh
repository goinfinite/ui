#!/usr/bin/env bash
# @description  Runs Playwright specs for one registered mode against the served demo.
# @usage        tests/ui/run.sh <smoke|standard|a11y|performance|toolset|control|cross-browser|toolset-cross-browser|control-cross-browser>
# @output       Playwright test report, exit code forwarded from the run.
# @requires     bash v4+, node v24+, DEMO_URL pointing at the demo page
# @version      0.3.0
# @updated      2026-09-15
set -euo pipefail

#
## Invocation
#
validateInvocation() {
	case "$1" in
		smoke | standard | a11y | performance | toolset | control | cross-browser | toolset-cross-browser | control-cross-browser) return 0 ;;
	esac
	echo "usage: tests/ui/run.sh <smoke|standard|a11y|performance|toolset|control|cross-browser|toolset-cross-browser|control-cross-browser>" >&2
	exit 2
}

requireDemoUrl() {
	if [ -z "${DEMO_URL:-}" ]; then
		echo "run.sh: DEMO_URL is unset; enter the suite through tests/tests.sh" >&2
		exit 2
	fi
}

#
## SpecSelection
#
runSelectedSpecs() {
	case "$1" in
		smoke) npx playwright test --project=chromium --grep @smoke ;;
		standard) npx playwright test --project=chromium --grep-invert "@smoke|@a11y|@perf|@toolset|@control" ;;
		a11y) npx playwright test --project=chromium --grep @a11y ;;
		performance) npx playwright test --project=chromium --grep @perf ;;
		toolset) npx playwright test --project=chromium --grep @toolset ;;
		control) npx playwright test --project=chromium --grep @control ;;
		cross-browser)
			# WebKit is defined in playwright.config.js but needs system libraries
			# this host lacks; enable it with --project=webkit where they exist.
			npx playwright test --project=firefox --grep-invert "@a11y|@perf|@toolset|@control"
			;;
		toolset-cross-browser) npx playwright test --project=firefox --grep @toolset ;;
		control-cross-browser) npx playwright test --project=firefox --grep @control ;;
	esac
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then return 0; fi

validateInvocation "${1:-}"
requireDemoUrl
cd "$(dirname "$0")"
runSelectedSpecs "$1"

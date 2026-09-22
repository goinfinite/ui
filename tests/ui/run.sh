#!/usr/bin/env bash
# @description  Runs Playwright specs for one registered mode against the served demo.
# @usage        tests/ui/run.sh <smoke|standard|a11y|performance|toolset|control|structural-smoke|structural|structural-a11y|cross-browser|toolset-cross-browser|control-cross-browser|structural-cross-browser>
# @output       Playwright test report, exit code forwarded from the run.
# @requires     bash v4+, node v24+, DEMO_URL pointing at the demo page
# @version      0.4.0
# @updated      2026-09-18
set -euo pipefail

#
## Invocation
#
validateInvocation() {
	case "$1" in
		smoke | standard | a11y | performance | toolset | control | structural-smoke | structural | structural-a11y | cross-browser | toolset-cross-browser | control-cross-browser | structural-cross-browser) return 0 ;;
	esac
	echo "usage: tests/ui/run.sh <smoke|standard|a11y|performance|toolset|control|structural-smoke|structural|structural-a11y|cross-browser|toolset-cross-browser|control-cross-browser|structural-cross-browser>" >&2
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
		smoke) ../node_modules/.bin/playwright test --project=chromium --grep @smoke ;;
		standard) ../node_modules/.bin/playwright test --project=chromium --grep-invert "@smoke|@a11y|@perf|@toolset|@control|@structural" ;;
		a11y) ../node_modules/.bin/playwright test --project=chromium --grep @a11y ;;
		performance) ../node_modules/.bin/playwright test --project=chromium --grep @perf ;;
		toolset) ../node_modules/.bin/playwright test --project=chromium --grep @toolset ;;
		control) ../node_modules/.bin/playwright test --project=chromium --grep @control ;;
		structural-smoke) ../node_modules/.bin/playwright test --project=chromium --grep "@structural.*@smoke" ;;
		structural) ../node_modules/.bin/playwright test --project=chromium --grep "@structural" --grep-invert "@smoke|@a11y" ;;
		structural-a11y) ../node_modules/.bin/playwright test --project=chromium --grep "@structural.*@a11y" ;;
		cross-browser)
			# WebKit is defined in playwright.config.js but needs system libraries
			# this host lacks; enable it with --project=webkit where they exist.
			../node_modules/.bin/playwright test --project=firefox --grep-invert "@a11y|@perf|@toolset|@control|@structural"
			;;
		toolset-cross-browser) ../node_modules/.bin/playwright test --project=firefox --grep @toolset ;;
		control-cross-browser) ../node_modules/.bin/playwright test --project=firefox --grep @control ;;
		structural-cross-browser) ../node_modules/.bin/playwright test --project=firefox --grep "@structural" --grep-invert "@a11y" ;;
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

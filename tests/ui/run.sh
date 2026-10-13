#!/usr/bin/env bash
# @description  Runs Playwright specs for one registered mode against the served demo.
# @usage        tests/ui/run.sh <smoke|standard|toolset|control|structural-smoke|structural|structural-a11y|a11y|performance|cross-browser|toolset-cross-browser|control-cross-browser|structural-cross-browser>
# @output       Playwright test report, exit code forwarded from the run.
# @requires     bash v4+, node v24+, DEMO_URL pointing at the demo page, UI_TEST_LIGHTPANDA_CDP_URL for the Lightpanda modes
# @version      0.5.1
# @updated      2026-10-12
set -euo pipefail

lightpandaModes=(smoke standard toolset control structural-smoke structural structural-a11y)
realEngineModes=(a11y performance cross-browser toolset-cross-browser control-cross-browser structural-cross-browser)

#
## Invocation
#
isModeListed() {
	local mode="$1"
	shift
	local listedMode
	for listedMode in "$@"; do
		if [[ "$listedMode" == "$mode" ]]; then
			return 0
		fi
	done
	return 1
}

modeRequiresLightpanda() {
	isModeListed "$1" "${lightpandaModes[@]}"
}

validateInvocation() {
	if isModeListed "$1" "${lightpandaModes[@]}" "${realEngineModes[@]}"; then
		return 0
	fi
	local modeChoices
	modeChoices="$(IFS='|' && echo "${lightpandaModes[*]}|${realEngineModes[*]}")"
	echo "usage: tests/ui/run.sh <$modeChoices>" >&2
	exit 2
}

requireDemoUrl() {
	if [[ -z "${DEMO_URL:-}" ]]; then
		echo "run.sh: DEMO_URL is unset; enter the suite through tests/tests.sh" >&2
		exit 2
	fi
}

requireLightpandaEndpoint() {
	if [[ -z "${UI_TEST_LIGHTPANDA_CDP_URL:-}" ]]; then
		echo "run.sh: UI_TEST_LIGHTPANDA_CDP_URL is unset; enter the suite through tests/tests.sh" >&2
		exit 2
	fi
}

#
## SpecSelection
#
runLightpandaSpecs() {
	local mode="$1"
	requireLightpandaEndpoint
	local playwrightCommand=(../node_modules/.bin/playwright test --project=lightpanda --output="../test-results/$mode")
	case "$mode" in
		smoke) "${playwrightCommand[@]}" --grep @smoke ;;
		standard) "${playwrightCommand[@]}" --grep-invert "@smoke|@a11y|@perf|@toolset|@control|@structural" ;;
		toolset) "${playwrightCommand[@]}" --grep @toolset ;;
		control) "${playwrightCommand[@]}" --grep @control ;;
		structural-smoke) "${playwrightCommand[@]}" --grep "@structural.*@smoke" ;;
		structural) "${playwrightCommand[@]}" --grep "@structural" --grep-invert "@smoke|@a11y" ;;
		structural-a11y) "${playwrightCommand[@]}" --grep "@structural.*@a11y" ;;
	esac
}

runRealEngineSpecs() {
	local mode="$1"
	local playwrightCommand=(env -u UI_TEST_LIGHTPANDA_CDP_URL ../node_modules/.bin/playwright test --workers=1 --output="../test-results/$mode")
	local realEngines=(--project=chromium --project=firefox)
	case "$mode" in
		a11y) "${playwrightCommand[@]}" --project=chromium --grep @a11y ;;
		performance) "${playwrightCommand[@]}" --project=chromium --grep @perf ;;
		cross-browser) "${playwrightCommand[@]}" "${realEngines[@]}" --grep-invert "@a11y|@perf|@toolset|@control|@structural" ;;
		toolset-cross-browser) "${playwrightCommand[@]}" "${realEngines[@]}" --grep @toolset ;;
		control-cross-browser) "${playwrightCommand[@]}" "${realEngines[@]}" --grep @control ;;
		structural-cross-browser) "${playwrightCommand[@]}" "${realEngines[@]}" --grep "@structural" --grep-invert "@a11y" ;;
	esac
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then return 0; fi

validateInvocation "${1:-}"
requireDemoUrl
cd "$(dirname "$0")"
specRunner=runRealEngineSpecs
if modeRequiresLightpanda "$1"; then
	specRunner=runLightpandaSpecs
fi
"$specRunner" "$1"

#!/usr/bin/env bash
# @description  Test entry point: selects registered suites, serves the demo and the Lightpanda browser, runs them, reports.
# @usage        tests/tests.sh [list] [--level=<fast|standard|exhaustive>] [--feature=<name>[,<name>...]] [--scope=<name>[,<name>...]]
# @output       Per-node progress on stdout, failures on stderr, final RESULT line, exit 0 pass / 1 failure / 2 runner error.
# @requires     bash v4+, node v24+, curl, timeout, go, lightpanda
# @version      0.5.1
# @updated      2026-10-12
# @see          .agents/skills/design-testing-strategy/SKILL.md
set -euo pipefail

#
## Invocation
#
parseInvocation() {
	local level="" feature="" scope="" listMode="false"
	for argument in "$@"; do
		case "$argument" in
			list) listMode="true" ;;
			--level=*) level="${argument#*=}" ;;
			--feature=*) feature="${argument#*=}" ;;
			--scope=*) scope="${argument#*=}" ;;
			*)
				echo "tests.sh: unknown argument: $argument" >&2
				exit 2
				;;
		esac
	done
	echo "$level|$feature|$scope|$listMode"
}

buildRegistryFilters() {
	local level="$1" feature="$2" scope="$3"
	if [[ -n "$level" ]]; then
		echo "--level=$level"
	fi
	if [[ -n "$feature" ]]; then
		echo "--feature=$feature"
	fi
	if [[ -n "$scope" ]]; then
		echo "--scope=$scope"
	fi
}

#
## BackgroundServer
#
waitForBackgroundServer() {
	local probeUrl="$1" serverPid="$2" failureMessage="$3"
	local readyDeadline=$((SECONDS + 10))
	while [[ "$SECONDS" -lt "$readyDeadline" ]]; do
		if kill -0 "$serverPid" 2>/dev/null &&
			curl -fs --max-time 2 -o /dev/null "$probeUrl"; then
			return 0
		fi
		sleep 0.2
	done
	echo "tests.sh: $failureMessage" >&2
	exit 2
}

#
## DemoServer
#
startDemoServer() {
	local demoPort="$1"
	node tests/lib/serveDemo.mjs "$demoPort" >/dev/null 2>&1 &
	echo "$!"
}

#
## LightpandaServer
#
lightpandaCdpMessageCeilingBytes=$((16 * 1024 * 1024))

requireLightpandaBinary() {
	if ! command -v lightpanda >/dev/null 2>&1; then
		echo "MissingRequiredTool: lightpanda (run 'mise install')" >&2
		exit 2
	fi
}

uiRunnerModeFinder() {
	local nodeCommand="$1"
	local commandAfterRunner="${nodeCommand#*tests/ui/run.sh }"
	if [[ "$commandAfterRunner" == "$nodeCommand" ]]; then
		return 1
	fi
	echo "${commandAfterRunner%% *}"
}

selectionRequiresLightpanda() {
	local selection="$1"
	while IFS=$'\t' read -r _ nodeCommand _; do
		local mode
		mode="$(uiRunnerModeFinder "$nodeCommand")" || continue
		if modeRequiresLightpanda "$mode"; then
			return 0
		fi
	done <<<"$selection"
	return 1
}

refuseOccupiedLightpandaPort() {
	local cdpUrl="$1"
	if curl -fs --max-time 2 -o /dev/null "$cdpUrl/json/version"; then
		echo "tests.sh: port ${cdpUrl##*:} already answers; stop the browser that holds it" >&2
		exit 2
	fi
}

startLightpandaServer() {
	local cdpPort="$1"
	lightpanda serve --host 127.0.0.1 --port "$cdpPort" \
		--cdp-max-message-size "$lightpandaCdpMessageCeilingBytes" \
		--disable-metrics \
		>"${TMPDIR:-/tmp}/lightpanda-${cdpPort}.log" 2>&1 &
	echo "$!"
}

#
## SuiteExecution
#
selectRegisteredNodes() {
	local selection
	selection="$(node tests/lib/registry.mjs select "$@")"
	if [[ -z "$selection" ]]; then
		echo "tests.sh: no tests selected" >&2
		exit 2
	fi
	echo "$selection"
}

selectionRequiresDemo() {
	local selection="$1"
	while IFS=$'\t' read -r _ _ nodeRequiresDemo; do
		if [[ "$nodeRequiresDemo" = "true" ]]; then
			return 0
		fi
	done <<<"$selection"
	return 1
}

runNode() {
	local nodeLabel="$1" nodeCommand="$2"
	echo ""
	echo "== $nodeLabel: $nodeCommand"
	timeout 1800 bash -c "$nodeCommand"
}

reportResult() {
	local passedCount="$1" failedCount="$2" timeoutCount="$3" elapsedSeconds="$4"
	echo ""
	if [[ "$timeoutCount" -gt 0 ]]; then
		echo "RESULT: TIMEOUT ($passedCount passed, $failedCount failed, ${elapsedSeconds}s)"
		exit 2
	fi
	if [[ "$failedCount" -gt 0 ]]; then
		echo "RESULT: FAIL ($passedCount passed, $failedCount failed, ${elapsedSeconds}s)"
		exit 1
	fi
	echo "RESULT: PASS ($passedCount passed, ${elapsedSeconds}s)"
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then return 0; fi

cd "$(dirname "$0")/.."
source tests/ui/run.sh

invocation="$(parseInvocation "$@")"
IFS='|' read -r level feature scope listMode <<<"$invocation"
mapfile -t registryFilters < <(buildRegistryFilters "$level" "$feature" "$scope")

if [[ "$listMode" = "true" ]]; then
	node tests/lib/registry.mjs list "${registryFilters[@]}"
	exit 0
fi

selectedNodes="$(selectRegisteredNodes "${registryFilters[@]}")"

if selectionRequiresDemo "$selectedNodes"; then
	demoPort="${UI_TEST_PORT:-8377}"
	demoServerPid=$(startDemoServer "$demoPort")
	trap 'kill "$demoServerPid" 2>/dev/null || true' EXIT
	export DEMO_URL="http://localhost:$demoPort"
	waitForBackgroundServer "$DEMO_URL/index.html" "$demoServerPid" \
		"could not serve docs/ on port $demoPort"
fi

if selectionRequiresLightpanda "$selectedNodes"; then
	requireLightpandaBinary
	lightpandaPort="${UI_TEST_LIGHTPANDA_PORT:-9377}"
	export UI_TEST_LIGHTPANDA_CDP_URL="http://127.0.0.1:$lightpandaPort"
	refuseOccupiedLightpandaPort "$UI_TEST_LIGHTPANDA_CDP_URL"
	lightpandaServerPid=$(startLightpandaServer "$lightpandaPort")
	trap 'kill "${demoServerPid:-}" "$lightpandaServerPid" 2>/dev/null || true' EXIT
	waitForBackgroundServer "$UI_TEST_LIGHTPANDA_CDP_URL/json/version" "$lightpandaServerPid" \
		"could not start lightpanda on port $lightpandaPort"
fi

startedAt=$SECONDS
passedCount=0
timeoutCount=0
totalCount=0
while IFS=$'\t' read -r nodeLabel nodeCommand _; do
	totalCount=$((totalCount + 1))
	exitCode=0
	runNode "$nodeLabel" "$nodeCommand" || exitCode=$?
	case "$exitCode" in
		0) passedCount=$((passedCount + 1)) ;;
		124)
			timeoutCount=$((timeoutCount + 1))
			echo "TIMEOUT: $nodeLabel exceeded the 30-minute ceiling" >&2
			;;
		*) echo "FAIL: $nodeLabel" >&2 ;;
	esac
done <<<"$selectedNodes"
reportResult "$passedCount" "$((totalCount - passedCount))" "$timeoutCount" "$((SECONDS - startedAt))"

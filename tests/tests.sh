#!/usr/bin/env bash
# @description  Test entry point: selects registered suites, serves the demo, runs them, reports.
# @usage        tests/tests.sh [list] [--level=<fast|standard|exhaustive>] [--feature=<name>[,<name>...]] [--scope=<name>[,<name>...]]
# @output       Per-node progress on stdout, failures on stderr, final RESULT line, exit 0 pass / 1 failure / 2 runner error.
# @requires     bash v4+, node v24+, curl, timeout, go
# @version      0.2.0
# @updated      2026-09-14
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
	if [ -n "$level" ]; then
		echo "--level=$level"
	fi
	if [ -n "$feature" ]; then
		echo "--feature=$feature"
	fi
	if [ -n "$scope" ]; then
		echo "--scope=$scope"
	fi
}

#
## DemoServer
#
startDemoServer() {
	local demoPort="$1"
	node tests/lib/serveDemo.mjs "$demoPort" >/dev/null 2>&1 &
	echo "$!"
}

waitForDemoServer() {
	local demoUrl="$1"
	local readyDeadline=$((SECONDS + 10))
	while [ "$SECONDS" -lt "$readyDeadline" ]; do
		if curl -s -o /dev/null "$demoUrl/index.html"; then
			return 0
		fi
		sleep 0.2
	done
	echo "tests.sh: could not serve docs/ on port ${demoUrl##*:}" >&2
	exit 2
}

#
## SuiteExecution
#
selectRegisteredNodes() {
	local selection
	selection="$(node tests/lib/registry.mjs select "$@")"
	if [ -z "$selection" ]; then
		echo "tests.sh: no tests selected" >&2
		exit 2
	fi
	echo "$selection"
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
	if [ "$timeoutCount" -gt 0 ]; then
		echo "RESULT: TIMEOUT ($passedCount passed, $failedCount failed, ${elapsedSeconds}s)"
		exit 2
	fi
	if [ "$failedCount" -gt 0 ]; then
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

IFS='|' read -r level feature scope listMode <<<"$(parseInvocation "$@")"
mapfile -t registryFilters < <(buildRegistryFilters "$level" "$feature" "$scope")

if [ "$listMode" = "true" ]; then
	node tests/lib/registry.mjs list "${registryFilters[@]}"
	exit 0
fi

selectedNodes="$(selectRegisteredNodes "${registryFilters[@]}")"
demoPort="${UI_TEST_PORT:-8377}"
demoServerPid=$(startDemoServer "$demoPort")
export DEMO_URL="http://localhost:$demoPort"
trap 'kill "$demoServerPid" 2>/dev/null || true' EXIT
waitForDemoServer "$DEMO_URL"

startedAt=$SECONDS
passedCount=0
timeoutCount=0
totalCount=0
while IFS=$'\t' read -r nodeLabel nodeCommand; do
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

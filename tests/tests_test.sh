#!/usr/bin/env bash
# @description  Companion tests for tests.sh: verifies invocation parsing and filter building.
# @usage        bash tests/tests_test.sh
# @output       PASS lines per assertion; exits 1 on the first failure.
# @requires     bash v4+
# @version      0.1.0
# @updated      2026-09-14
set -euo pipefail

source "$(dirname "$0")/tests.sh"

#
## Assertions
#
assertEquals() {
	local label="$1" expected="$2" actual="$3"
	if [ "$expected" != "$actual" ]; then
		echo "AssertionFailed $label: expected '$expected', got '$actual'" >&2
		exit 1
	fi
	echo "PASS $label"
}

#
## Runtime
#
IFS='|' read -r level feature scope listMode <<<"$(parseInvocation --level=fast --feature=form)"
assertEquals "parses level and feature" "fast|form||false" "$level|$feature|$scope|$listMode"

IFS='|' read -r level feature scope listMode <<<"$(parseInvocation list --scope=ui)"
assertEquals "parses list subcommand with scope" "||ui|true" "$level|$feature|$scope|$listMode"

IFS='|' read -r level feature scope listMode <<<"$(parseInvocation)"
assertEquals "defaults to no filters" "|||false" "$level|$feature|$scope|$listMode"

exitStatus=0
(parseInvocation --bogus >/dev/null 2>&1) || exitStatus=$?
assertEquals "rejects unknown argument with exit 2" "2" "$exitStatus"

filters="$(buildRegistryFilters fast "form,toolset" "" | tr '\n' ' ')"
assertEquals "builds registry filters" "--level=fast --feature=form,toolset " "$filters"

filters="$(buildRegistryFilters "" "" "")"
assertEquals "empty filters produce no arguments" "" "$filters"

#!/usr/bin/env bash
# @description  Assertion helpers shared by the tests/ companion scripts.
# @usage        source tests/lib/assertions.sh
# @output       PASS lines per assertion; failures on stderr.
# @requires     bash v4+
# @version      0.1.0
# @updated      2026-09-24
# @see          .agents/rules/code/shell.md

assertionFailureCount=0

assertEquals() {
	local label="$1" expected="$2" actual="$3"
	if [ "$expected" != "$actual" ]; then
		echo "AssertionFailed $label: expected '$expected', got '$actual'" >&2
		assertionFailureCount=$((assertionFailureCount + 1))
		return 0
	fi
	echo "PASS $label"
}

assertExitCode() {
	local label="$1" expected="$2"
	shift 2
	local actual=0
	"$@" >/dev/null 2>&1 || actual=$?
	assertEquals "$label" "$expected" "$actual"
}

assertContains() {
	local label="$1" needle="$2" haystack="$3"
	case "$haystack" in
		*"$needle"*)
			echo "PASS $label"
			;;
		*)
			echo "AssertionFailed $label: '$needle' missing from output" >&2
			assertionFailureCount=$((assertionFailureCount + 1))
			;;
	esac
}

assertEveryLineMatches() {
	local label="$1" pattern="$2" output="$3"
	local mismatchedCount
	mismatchedCount="$(printf '%s\n' "$output" | grep -cv "$pattern" || true)"
	assertEquals "$label" "0" "$mismatchedCount"
}

exitWithAssertionStatus() {
	if [ "$assertionFailureCount" -gt 0 ]; then
		exit 1
	fi
	exit 0
}

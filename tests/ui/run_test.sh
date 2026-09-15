#!/usr/bin/env bash
# @description  Companion tests for run.sh: verifies mode validation and DEMO_URL requirement.
# @usage        bash tests/ui/run_test.sh
# @output       PASS lines per assertion; exits 1 on the first failure.
# @requires     bash v4+
# @version      0.2.0
# @updated      2026-09-15
set -euo pipefail

source "$(dirname "$0")/run.sh"

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
exitStatus=0
(validateInvocation smoke) || exitStatus=$?
assertEquals "accepts registered mode" "0" "$exitStatus"

exitStatus=0
(validateInvocation toolset) || exitStatus=$?
assertEquals "accepts toolset mode" "0" "$exitStatus"

exitStatus=0
(validateInvocation bogus-mode >/dev/null 2>&1) || exitStatus=$?
assertEquals "rejects unknown mode with exit 2" "2" "$exitStatus"

exitStatus=0
(validateInvocation "" >/dev/null 2>&1) || exitStatus=$?
assertEquals "rejects missing mode with exit 2" "2" "$exitStatus"

exitStatus=0
(DEMO_URL="http://localhost:8377" requireDemoUrl) || exitStatus=$?
assertEquals "accepts set DEMO_URL" "0" "$exitStatus"

exitStatus=0
(unset DEMO_URL; requireDemoUrl >/dev/null 2>&1) || exitStatus=$?
assertEquals "refuses unset DEMO_URL with exit 2" "2" "$exitStatus"

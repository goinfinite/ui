#!/usr/bin/env bash
# @description  CLI contract for run.sh: rejected invocations and the DEMO_URL requirement.
# @usage        bash tests/ui/run_test.sh
# @output       PASS lines per assertion; exits 1 if any assertion fails.
# @requires     bash v4+, node v24+
# @version      0.4.0
# @updated      2026-09-24
set -uo pipefail

cd "$(dirname "$0")/../.." || exit 2

source tests/lib/assertions.sh

#
## Runtime
#
# A set DEMO_URL keeps the DEMO_URL guard out of the way, so these two
# invocations reach the mode validation they claim to test.
assertExitCode "rejects an unknown mode with exit 2" "2" \
	env DEMO_URL=http://localhost:8377 bash tests/ui/run.sh bogus-mode
assertExitCode "rejects a missing mode with exit 2" "2" \
	env DEMO_URL=http://localhost:8377 bash tests/ui/run.sh ""
assertExitCode "refuses an unset DEMO_URL with exit 2" "2" \
	env -u DEMO_URL bash tests/ui/run.sh smoke

exitWithAssertionStatus

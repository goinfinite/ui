#!/usr/bin/env bash
# @description  CLI contract for mutate.sh: rejected invocations and the gremlins binary guard.
# @usage        bash tests/mutate_test.sh
# @output       PASS lines per assertion; exits 1 if any assertion fails.
# @requires     bash v4+
# @version      0.1.0
# @updated      2026-10-15
set -uo pipefail

cd "$(dirname "$0")/.." || exit 2

source tests/lib/assertions.sh

#
## Runtime
#
assertExitCode "rejects a missing feature with exit 2" "2" \
	bash tests/mutate.sh ""
assertExitCode "rejects an unknown feature with exit 2" "2" \
	bash tests/mutate.sh bogus-feature
assertContains "names the accepted features in the usage hint" \
	"<control|display|form|structural|toolset>" \
	"$(bash tests/mutate.sh 2>&1)"

exitWithAssertionStatus
#!/usr/bin/env bash
# @description  CLI contract for tests.sh: registry filters, rejected invocations, and demo-server decisions.
# @usage        bash tests/tests_test.sh
# @output       PASS lines per assertion; exits 1 if any assertion fails.
# @requires     bash v4+, node v24+, go
# @version      0.3.0
# @updated      2026-09-24
set -uo pipefail

cd "$(dirname "$0")/.." || exit 2

source tests/lib/assertions.sh

#
## RegistryFilters
#
unfilteredList="$(bash tests/tests.sh list)"
assertContains "prints the unit tree" "runner:unit/fast" "$unfilteredList"
assertContains "prints the ui tree" ":ui/" "$unfilteredList"

formFastNodes="$(bash tests/tests.sh list --level=fast --feature=form)"
assertEveryLineMatches "keeps only the form feature" "^form:" "$formFastNodes"
assertEveryLineMatches "keeps only the fast level" "/fast[[:space:]]" "$formFastNodes"

uiNodes="$(bash tests/tests.sh list --scope=ui)"
assertEveryLineMatches "keeps only the ui scope" ":ui/" "$uiNodes"

#
## RejectedInvocations
#
assertExitCode "rejects an unknown argument with exit 2" "2" bash tests/tests.sh --bogus
assertExitCode "rejects an unknown feature with exit 2" "2" bash tests/tests.sh list --feature=nonexistent
assertExitCode "rejects an empty selection with exit 2" "2" bash tests/tests.sh --feature=runner --scope=ui
assertEquals "prints no nodes when a filter matches nothing" "" "$(bash tests/tests.sh list --feature=runner --scope=ui)"

#
## DemoServerDecision
#
# An unusable port number makes the serve step fail before any node runs.
unusablePortNumber=99999
assertExitCode "runs a demo-free node without serving the demo" "0" \
	env "UI_TEST_PORT=$unusablePortNumber" bash tests/tests.sh --feature=display --scope=unit --level=fast
assertExitCode "serves the demo before a demo-requiring node" "2" \
	env "UI_TEST_PORT=$unusablePortNumber" bash tests/tests.sh --feature=form --scope=ui --level=fast

exitWithAssertionStatus

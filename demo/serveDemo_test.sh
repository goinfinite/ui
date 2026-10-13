#!/usr/bin/env bash
# @description  Behavior of serveDemo.sh: rejected ports and actions, and the start, status, stop lifecycle.
# @usage        bash demo/serveDemo_test.sh
# @output       PASS lines per assertion; exits 1 if any assertion fails.
# @requires     bash v4+, node v24+, curl
# @version      0.1.1
# @updated      2026-10-13
set -uo pipefail

cd "$(dirname "$0")/.." || exit 2

source tests/lib/assertions.sh
source demo/serveDemo.sh
# The sourced script turns on errexit; an assertion run keeps going after a failure.
set +e

lifecyclePort=8461
trap 'stopDemoServer "$lifecyclePort"' EXIT

#
## RejectedInvocations
#
assertExitCode "rejects an unknown action with exit 2" "2" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerAction serve'
assertExitCode "rejects a non-numeric port with exit 2" "2" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerPort abc'
assertExitCode "rejects a port above 65535 with exit 2" "2" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerPort 70000'
assertExitCode "rejects the test suite port with exit 2" "2" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerPort 8377'
assertExitCode "rejects a zero-padded test suite port with exit 2" "2" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerPort 08377'
assertExitCode "accepts the default port" "0" \
	bash -c 'source demo/serveDemo.sh && validateDemoServerPort 8399'

#
## Lifecycle
#
stopDemoServer "$lifecyclePort"
assertEquals "reports a stopped server" "stopped" \
	"$(reportDemoServerStatus "$lifecyclePort")"

startedUrl="$(startDemoServer "$lifecyclePort")"
assertEquals "start prints the demo url" "http://127.0.0.1:${lifecyclePort}/" "$startedUrl"
assertEquals "reports a running server" "running http://127.0.0.1:${lifecyclePort}/" \
	"$(reportDemoServerStatus "$lifecyclePort")"
assertEquals "a second start reuses the running server" "$startedUrl" \
	"$(startDemoServer "$lifecyclePort")"

stopDemoServer "$lifecyclePort"
assertEquals "stop ends the server before it returns" "stopped" \
	"$(reportDemoServerStatus "$lifecyclePort")"

exitWithAssertionStatus

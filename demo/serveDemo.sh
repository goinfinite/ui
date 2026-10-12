#!/usr/bin/env bash
# @description  Starts, stops, or reports the local demo server that serves docs/ for manual inspection.
# @usage        demo/serveDemo.sh <start|stop|status> [port]
# @output       The demo URL on start, the server state on status, nothing on stop.
# @requires     bash v4+, node v24+, curl
# @version      0.1.0
# @updated      2026-10-12
# @see          tests/lib/serveDemo.mjs
set -euo pipefail

defaultDemoPort=8399
testSuiteDemoPort=8377

#
## Invocation
#
validateDemoServerAction() {
	case "$1" in
		start | stop | status) return 0 ;;
	esac
	echo "usage: demo/serveDemo.sh <start|stop|status> [port]" >&2
	exit 2
}

validateDemoServerPort() {
	local port="$1"
	if [[ ! "$port" =~ ^[1-9][0-9]{0,4}$ ]] || ((port > 65535)); then
		echo "InvalidDemoServerPort: $port" >&2
		exit 2
	fi
	if ((port == testSuiteDemoPort)); then
		echo "DemoServerPortReserved: tests/tests.sh owns port $testSuiteDemoPort" >&2
		exit 2
	fi
}

#
## ServerProcess
#
demoServerProcessPatternBuilder() {
	local port="$1"
	echo "serveDemo\.mjs ${port}\$"
}

demoServerPidFinder() {
	local port="$1"
	pgrep -f "$(demoServerProcessPatternBuilder "$port")" || true
}

demoServerUrlBuilder() {
	local port="$1"
	echo "http://127.0.0.1:${port}/"
}

isDemoServerAnswering() {
	local port="$1"
	local statusCode
	statusCode="$(curl -s -o /dev/null -w '%{http_code}' "$(demoServerUrlBuilder "$port")" || true)"
	[[ "$statusCode" == "200" ]]
}

waitForDemoServer() {
	local port="$1"
	local attemptsLeft=20
	while ((attemptsLeft > 0)); do
		if isDemoServerAnswering "$port"; then
			return 0
		fi
		attemptsLeft=$((attemptsLeft - 1))
		sleep 0.5
	done
	echo "DemoServerNotReady: port $port did not answer 200 within 10 seconds" >&2
	exit 1
}

#
## Actions
#
startDemoServer() {
	local port="$1"
	local projectRoot logFile
	if [[ -n "$(demoServerPidFinder "$port")" ]]; then
		demoServerUrlBuilder "$port"
		return 0
	fi
	projectRoot="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
	logFile="${TMPDIR:-/tmp}/serveDemo-${port}.log"
	(node "$projectRoot/tests/lib/serveDemo.mjs" "$port" >"$logFile" 2>&1 &)
	waitForDemoServer "$port"
	demoServerUrlBuilder "$port"
}

stopDemoServer() {
	local port="$1"
	pkill -f "$(demoServerProcessPatternBuilder "$port")" || true
}

reportDemoServerStatus() {
	local port="$1"
	if [[ -z "$(demoServerPidFinder "$port")" ]]; then
		echo "stopped"
		return 0
	fi
	echo "running $(demoServerUrlBuilder "$port")"
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then return 0; fi

demoServerAction="${1:-}"
demoServerPort="${2:-$defaultDemoPort}"
validateDemoServerAction "$demoServerAction"
validateDemoServerPort "$demoServerPort"

case "$demoServerAction" in
	start) startDemoServer "$demoServerPort" ;;
	stop) stopDemoServer "$demoServerPort" ;;
	status) reportDemoServerStatus "$demoServerPort" ;;
esac

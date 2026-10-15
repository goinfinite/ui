#!/usr/bin/env bash
# @description  Runs Gremlins mutation testing over one feature's Go packages and fails below the efficacy floor.
# @usage        tests/mutate.sh <control|display|form|structural|toolset>
# @output       Gremlins per-mutant progress, a summary line per package, exit 0 pass / 1 below floor / 2 bad invocation.
# @requires     bash v4+, go, gremlins
# @version      0.1.0
# @updated      2026-10-15
# @see          .agents/skills/design-testing-strategy/SKILL.md
set -uo pipefail

readonly mutationEfficacyFloorPercent=80
readonly mutationCoverageFloorPercent=80
# Gremlins derives its per-mutant timeout from the baseline test time. A fast
# baseline yields a timeout too short for compilation, so mutants time out and
# score as if they survived. The coefficient widens the window to cover build time.
readonly mutationTimeoutCoefficient=20
readonly mutationWorkerCount=3
readonly generatedFilePattern='_templ\.go'

#
## Validation
#
requireKnownFeature() {
	local feature="$1"
	case "$feature" in
		control | display | form | structural | toolset) return 0 ;;
		*)
			echo "mutate.sh: unknown feature: $feature" >&2
			exit 2
			;;
	esac
}

featurePackagePathResolver() {
	local feature="$1"
	echo "./src/$feature/"
}

requireGremlinsBinary() {
	if ! command -v gremlins >/dev/null 2>&1; then
		echo "MissingRequiredTool: gremlins (run 'mise install')" >&2
		exit 2
	fi
}

#
## Mutation
#
runMutationTesting() {
	local packagePath="$1"
	gremlins unleash "$packagePath" \
		--workers "$mutationWorkerCount" \
		--integration \
		--exclude-files "$generatedFilePattern" \
		--timeout-coefficient "$mutationTimeoutCoefficient" \
		--threshold-efficacy "$mutationEfficacyFloorPercent" \
		--threshold-mcover "$mutationCoverageFloorPercent"
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then return 0; fi

feature="${1:-}"
if [[ -z "$feature" ]]; then
	echo "usage: tests/mutate.sh <control|display|form|structural|toolset>" >&2
	exit 2
fi

requireGremlinsBinary
requireKnownFeature "$feature"
packagePath="$(featurePackagePathResolver "$feature")"
runMutationTesting "$packagePath"

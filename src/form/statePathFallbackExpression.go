package uiForm

import "strconv"

func statePathFallbackExpressionBuilder(statePath, fallback string) string {
	if statePath == "" {
		return ""
	}
	return "(" + statePath + ") || " + strconv.Quote(fallback)
}

package uiStructural

import (
	_ "embed"
	"maps"

	"github.com/a-h/templ"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

//go:embed serverFragmentRefreshComponentState.js
var serverFragmentRefreshComponentAlpineState string

var serverFragmentRefreshComponentAlpineStateOnce = templ.NewOnceHandle(
	templ.WithComponent(uiToolset.MinifierTemplateJs(&serverFragmentRefreshComponentAlpineState)),
)

const defaultRefreshDebounceMs uint = 300

func filterQueryParamNamesResolver(
	filters []FilterSettings,
) map[string]string {
	queryParamNames := map[string]string{}
	for _, filter := range filters {
		queryParamName := filter.QueryParamName
		if queryParamName == "" {
			queryParamName = filter.Key
		}
		queryParamNames[filter.Key] = queryParamName
	}
	return queryParamNames
}

func refreshDebounceResolver(refreshDebounceMs uint) uint {
	if refreshDebounceMs == 0 {
		return defaultRefreshDebounceMs
	}
	return refreshDebounceMs
}

func initialFilterValuesResolver(
	filters []FilterSettings, providedValues map[string]any,
) map[string]any {
	initialValues := map[string]any{}
	for _, filter := range filters {
		switch filter.Kind {
		case FilterKindNumberRange, FilterKindDateRange:
			initialValues[filter.Key] = map[string]string{"min": "", "max": ""}
		default:
			initialValues[filter.Key] = ""
		}
	}
	maps.Copy(initialValues, providedValues)
	return initialValues
}

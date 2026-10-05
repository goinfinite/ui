package uiStructural

import (
	"reflect"
	"testing"
)

func TestInitialFilterValuesResolver(t *testing.T) {
	testCases := []struct {
		name           string
		filters        []FilterSettings
		providedValues map[string]any
		expectedValues map[string]any
	}{
		{
			name: "defaults for scalar and range filters",
			filters: []FilterSettings{
				{Key: "name", Kind: FilterKindTextContains},
				{Key: "cpu", Kind: FilterKindNumberRange},
			},
			expectedValues: map[string]any{
				"name": "",
				"cpu":  map[string]string{"min": "", "max": ""},
			},
		},
		{
			name: "provided values override defaults",
			filters: []FilterSettings{
				{Key: "status", Kind: FilterKindEnumSelect},
			},
			providedValues: map[string]any{"status": "running"},
			expectedValues: map[string]any{"status": "running"},
		},
		{
			name:           "no filters and no provided values",
			filters:        []FilterSettings{},
			expectedValues: map[string]any{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualValues := initialFilterValuesResolver(
				testCase.filters, testCase.providedValues,
			)
			if !reflect.DeepEqual(actualValues, testCase.expectedValues) {
				t.Errorf(
					"InitialFilterValuesMismatch: got %v, want %v",
					actualValues, testCase.expectedValues,
				)
			}
		})
	}
}

func TestFilterQueryParamNamesResolver(t *testing.T) {
	testCases := []struct {
		name               string
		filters            []FilterSettings
		expectedParamNames map[string]string
	}{
		{
			name: "key is the default query parameter name",
			filters: []FilterSettings{
				{Key: "status", Kind: FilterKindEnumSelect},
			},
			expectedParamNames: map[string]string{"status": "status"},
		},
		{
			name: "provided query parameter name wins",
			filters: []FilterSettings{
				{Key: "cpu", Kind: FilterKindNumberRange, QueryParamName: "cores"},
			},
			expectedParamNames: map[string]string{"cpu": "cores"},
		},
		{
			name:               "no filters",
			filters:            []FilterSettings{},
			expectedParamNames: map[string]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualParamNames := filterQueryParamNamesResolver(testCase.filters)
			if !reflect.DeepEqual(actualParamNames, testCase.expectedParamNames) {
				t.Errorf(
					"QueryParamNamesMismatch: got %v, want %v",
					actualParamNames, testCase.expectedParamNames,
				)
			}
		})
	}
}

func TestRefreshDebounceResolver(t *testing.T) {
	if actual := refreshDebounceResolver(0); actual != defaultRefreshDebounceMs {
		t.Errorf(
			"RefreshDebounceMismatch: got %d, want %d",
			actual, defaultRefreshDebounceMs,
		)
	}
	if actual := refreshDebounceResolver(750); actual != 750 {
		t.Errorf("RefreshDebounceMismatch: got %d, want 750", actual)
	}
}

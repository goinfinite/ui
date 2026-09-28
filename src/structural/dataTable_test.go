package uiStructural

import (
	"reflect"
	"testing"

	uiForm "github.com/goinfinite/ui/src/form"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

type dataTableTestRecord struct {
	Id   string
	Name string
}

func TestDataTableInitialFilterValuesResolver(t *testing.T) {
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
			settings := DataTableSettings[dataTableTestRecord]{
				Filters:             testCase.filters,
				InitialFilterValues: testCase.providedValues,
			}
			actualValues := settings.initialFilterValuesResolver()
			if !reflect.DeepEqual(actualValues, testCase.expectedValues) {
				t.Errorf(
					"InitialFilterValuesMismatch: got %v, want %v",
					actualValues, testCase.expectedValues,
				)
			}
		})
	}
}

func TestDataTableClientSettingsResolver(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		Filters: []FilterSettings{
			{Key: "status", Kind: FilterKindEnumSelect},
			{Key: "cpu", Kind: FilterKindNumberRange, QueryParamName: "cores"},
		},
		InitialFilterValues:  map[string]any{"status": "running"},
		RefreshOnEvents:      []string{"refresh:records"},
		InitialSearchQuery:   "alpha",
		InitialSortDirection: DataTableSortDirectionAsc,
		InitialSortKey:       "name",
		QueryUrlTemplate:     "/records?page={pageNumber}",
	}
	clientSettings := settings.clientSettingsResolver(5, 2)

	expectedQueryParamNames := map[string]string{"status": "status", "cpu": "cores"}
	if !reflect.DeepEqual(clientSettings.FilterQueryParamNames, expectedQueryParamNames) {
		t.Errorf(
			"FilterQueryParamNamesMismatch: got %v, want %v",
			clientSettings.FilterQueryParamNames, expectedQueryParamNames,
		)
	}
	expectedFilterValues := map[string]any{
		"status": "running",
		"cpu":    map[string]string{"min": "", "max": ""},
	}
	if !reflect.DeepEqual(clientSettings.InitialState.FilterValues, expectedFilterValues) {
		t.Errorf(
			"FilterValuesMismatch: got %v, want %v",
			clientSettings.InitialState.FilterValues, expectedFilterValues,
		)
	}
	if clientSettings.InitialState.ItemsPerPage != 5 {
		t.Errorf("ItemsPerPageMismatch: got %d, want 5", clientSettings.InitialState.ItemsPerPage)
	}
	if clientSettings.InitialState.PageNumber != 2 {
		t.Errorf("PageNumberMismatch: got %d, want 2", clientSettings.InitialState.PageNumber)
	}
	if clientSettings.InitialState.SearchQuery != "alpha" {
		t.Errorf("SearchQueryMismatch: got %q, want %q", clientSettings.InitialState.SearchQuery, "alpha")
	}
	if clientSettings.InitialState.SortDirection != "asc" {
		t.Errorf("SortDirectionMismatch: got %q, want %q", clientSettings.InitialState.SortDirection, "asc")
	}
	if clientSettings.InitialState.SortKey != "name" {
		t.Errorf("SortKeyMismatch: got %q, want %q", clientSettings.InitialState.SortKey, "name")
	}
	if clientSettings.RefreshDebounceMs != dataTableDefaultRefreshDebounceMs {
		t.Errorf(
			"RefreshDebounceMsMismatch: got %d, want %d",
			clientSettings.RefreshDebounceMs, dataTableDefaultRefreshDebounceMs,
		)
	}
	if clientSettings.QueryUrlTemplate != "/records?page={pageNumber}" {
		t.Errorf("QueryUrlTemplateMismatch: got %q, want %q", clientSettings.QueryUrlTemplate, "/records?page={pageNumber}")
	}

	overrideSettings := settings
	overrideSettings.RefreshDebounceMs = 750
	overrideClientSettings := overrideSettings.clientSettingsResolver(5, 1)
	if overrideClientSettings.RefreshDebounceMs != 750 {
		t.Errorf("RefreshDebounceMsMismatch: got %d, want 750", overrideClientSettings.RefreshDebounceMs)
	}
}

func TestDataTableAlignmentClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		alignment         DataTableAlignment
		expectedClassName string
	}{
		{name: "left", alignment: DataTableAlignmentLeft, expectedClassName: "text-left"},
		{name: "center", alignment: DataTableAlignmentCenter, expectedClassName: "text-center"},
		{name: "right", alignment: DataTableAlignmentRight, expectedClassName: "text-right"},
		{name: "default", alignment: "", expectedClassName: "text-left"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := testCase.alignment.alignmentClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"AlignmentClassMismatch(%q): got %q, want %q",
					testCase.alignment, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestDataTableJustifyClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		alignment         DataTableAlignment
		expectedClassName string
	}{
		{name: "left", alignment: DataTableAlignmentLeft, expectedClassName: "justify-start"},
		{name: "center", alignment: DataTableAlignmentCenter, expectedClassName: "justify-center"},
		{name: "right", alignment: DataTableAlignmentRight, expectedClassName: "justify-end"},
		{name: "default", alignment: "", expectedClassName: "justify-start"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := testCase.alignment.justifyClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"JustifyClassMismatch(%q): got %q, want %q",
					testCase.alignment, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestDataTableTextCaseClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		textCase          string
		expectedClassName string
	}{
		{name: "lower", textCase: uiToolset.TextCaseLower, expectedClassName: "lowercase"},
		{name: "upper", textCase: uiToolset.TextCaseUpper, expectedClassName: "uppercase"},
		{name: "capitalize", textCase: uiToolset.TextCaseCapitalize, expectedClassName: "capitalize"},
		{name: "default", textCase: "", expectedClassName: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			settings := DataTableSettings[dataTableTestRecord]{
				TextCase: testCase.textCase,
			}
			actualClassName := settings.textCaseClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"TextCaseClassMismatch(%q): got %q, want %q",
					testCase.textCase, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestDataTableCheckboxShapeResolver(t *testing.T) {
	testCases := []struct {
		name          string
		providedShape string
		expectedShape string
	}{
		{name: "default", providedShape: "", expectedShape: uiForm.CheckboxInputShapeSquare},
		{name: "provided shape wins", providedShape: uiForm.CheckboxInputShapeCircular, expectedShape: uiForm.CheckboxInputShapeCircular},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			settings := DataTableSettings[dataTableTestRecord]{
				CheckboxShape: testCase.providedShape,
			}
			actualShape := settings.checkboxShapeResolver()
			if actualShape != testCase.expectedShape {
				t.Errorf(
					"CheckboxShapeMismatch(%q): got %q, want %q",
					testCase.providedShape, actualShape, testCase.expectedShape,
				)
			}
		})
	}
}

func TestDataTableCheckboxSizeResolver(t *testing.T) {
	testCases := []struct {
		name         string
		providedSize string
		expectedSize string
	}{
		{name: "default", providedSize: "", expectedSize: uiForm.CheckboxInputSizeMd},
		{name: "provided size wins", providedSize: uiForm.CheckboxInputSizeSm, expectedSize: uiForm.CheckboxInputSizeSm},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			settings := DataTableSettings[dataTableTestRecord]{
				CheckboxSize: testCase.providedSize,
			}
			actualSize := settings.checkboxSizeResolver()
			if actualSize != testCase.expectedSize {
				t.Errorf(
					"CheckboxSizeMismatch(%q): got %q, want %q",
					testCase.providedSize, actualSize, testCase.expectedSize,
				)
			}
		})
	}
}

func TestDataTableRowStripeClassesResolver(t *testing.T) {
	stripedSettings := DataTableSettings[dataTableTestRecord]{IsStriped: true}
	expectedClasses := "odd:bg-neutral-50/5 odd:hover:bg-neutral-50/10"
	actualStripedClasses := stripedSettings.rowStripeClassesResolver()
	if actualStripedClasses != expectedClasses {
		t.Errorf(
			"RowStripeClassesMismatch: got %q, want %q",
			actualStripedClasses, expectedClasses,
		)
	}

	plainSettings := DataTableSettings[dataTableTestRecord]{}
	actualPlainClasses := plainSettings.rowStripeClassesResolver()
	if actualPlainClasses != "" {
		t.Errorf("RowStripeClassesNotEmpty: %q", actualPlainClasses)
	}
}

func TestDataTablePaginationAriaLabelResolver(t *testing.T) {
	testCases := []struct {
		name          string
		providedLabel string
		expectedLabel string
	}{
		{
			name:          "default label",
			providedLabel: "",
			expectedLabel: "Table pagination",
		},
		{
			name:          "provided label wins",
			providedLabel: "Servers table pagination",
			expectedLabel: "Servers table pagination",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			settings := DataTableSettings[dataTableTestRecord]{
				PaginationAriaLabel: testCase.providedLabel,
			}
			actualLabel := settings.paginationAriaLabelResolver()
			if actualLabel != testCase.expectedLabel {
				t.Errorf(
					"PaginationAriaLabelMismatch(%q): got %q, want %q",
					testCase.providedLabel, actualLabel, testCase.expectedLabel,
				)
			}
		})
	}
}

func TestDataTableItemsPerPageSizeChoicesResolver(t *testing.T) {
	providedSettings := DataTableSettings[dataTableTestRecord]{
		ItemsPerPageSizeChoices: []DataTablePageSize{10, 20},
	}
	actualProvidedChoices := providedSettings.itemsPerPageSizeChoicesResolver()
	expectedProvidedChoices := []uint{10, 20}
	if !reflect.DeepEqual(actualProvidedChoices, expectedProvidedChoices) {
		t.Errorf(
			"ItemsPerPageSizeChoicesMismatch: got %v, want %v",
			actualProvidedChoices, expectedProvidedChoices,
		)
	}

	emptySettings := DataTableSettings[dataTableTestRecord]{}
	actualDefaultChoices := emptySettings.itemsPerPageSizeChoicesResolver()
	expectedDefaultChoices := paginationDefaultItemsPerPageSizeChoices
	if !reflect.DeepEqual(actualDefaultChoices, expectedDefaultChoices) {
		t.Errorf(
			"ItemsPerPageSizeChoicesMismatch: got %v, want %v",
			actualDefaultChoices, expectedDefaultChoices,
		)
	}
}

func TestDataTableItemsPerPageResolver(t *testing.T) {
	testCases := []struct {
		name                    string
		itemsPerPageSizeChoices []uint
		itemsPerPage            DataTablePageSize
		expectedItemsPerPage    uint
	}{
		{
			name:                    "explicit items per page wins",
			itemsPerPageSizeChoices: []uint{5, 10, 30, 50},
			itemsPerPage:            30,
			expectedItemsPerPage:    30,
		},
		{
			name:                    "zero falls back to the first choice",
			itemsPerPageSizeChoices: []uint{5, 10, 30, 50},
			itemsPerPage:            0,
			expectedItemsPerPage:    5,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			settings := DataTableSettings[dataTableTestRecord]{
				ItemsPerPage: testCase.itemsPerPage,
			}
			actualItemsPerPage := settings.itemsPerPageResolver(
				testCase.itemsPerPageSizeChoices,
			)
			if actualItemsPerPage != testCase.expectedItemsPerPage {
				t.Errorf(
					"ItemsPerPageMismatch: got %d, want %d",
					actualItemsPerPage, testCase.expectedItemsPerPage,
				)
			}
		})
	}
}

func TestDataTablePageNumberResolver(t *testing.T) {
	defaultSettings := DataTableSettings[dataTableTestRecord]{}
	actualDefaultPageNumber := defaultSettings.pageNumberResolver()
	if actualDefaultPageNumber != 1 {
		t.Errorf("PageNumberMismatch: got %d, want 1", actualDefaultPageNumber)
	}

	providedSettings := DataTableSettings[dataTableTestRecord]{PageNumber: 4}
	actualProvidedPageNumber := providedSettings.pageNumberResolver()
	if actualProvidedPageNumber != 4 {
		t.Errorf("PageNumberMismatch: got %d, want 4", actualProvidedPageNumber)
	}
}

func TestDataTableIdResolver(t *testing.T) {
	explicitIdSettings := DataTableSettings[dataTableTestRecord]{
		Id:      "records-table",
		Columns: []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
	}
	actualExplicitId := explicitIdSettings.idResolver()
	if actualExplicitId != "records-table" {
		t.Errorf("IdMismatch: got %q, want %q", actualExplicitId, "records-table")
	}

	settings := DataTableSettings[dataTableTestRecord]{
		QueryUrlTemplate: "/records?page={pageNumber}",
		Columns:          []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
	}
	derivedId := settings.idResolver()
	repeatedId := settings.idResolver()
	if repeatedId != derivedId {
		t.Errorf("IdResolverNotStable: %q != %q", repeatedId, derivedId)
	}

	rowVariantSettings := settings
	rowVariantSettings.Rows = []dataTableTestRecord{{Id: "one", Name: "alpha"}}
	rowVariantId := rowVariantSettings.idResolver()
	if rowVariantId != derivedId {
		t.Errorf(
			"IdChangedWithDifferentRows: %q != %q",
			rowVariantId, derivedId,
		)
	}

	filteredSettings := settings
	filteredSettings.Filters = []FilterSettings{
		{Key: "status", Kind: FilterKindEnumSelect},
	}
	filteredId := filteredSettings.idResolver()
	if filteredId == derivedId {
		t.Errorf("IdCollidedForDifferentFilters: %q", filteredId)
	}

	sortedSettings := settings
	sortedSettings.Columns = []DataTableColumnSettings[dataTableTestRecord]{
		{Label: "Name", SortKey: "name"},
	}
	sortedId := sortedSettings.idResolver()
	if sortedId == derivedId {
		t.Errorf("IdCollidedForDifferentSortKeys: %q", sortedId)
	}
}

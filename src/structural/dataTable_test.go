package uiStructural

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	uiForm "github.com/goinfinite/ui/src/form"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

type dataTableTestRecord struct {
	Id   string
	Name string
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
		MaxVisibleRows:       3,
		QueryUrlTemplate:     "/records?page={pageNumber}",
	}
	clientSettings := settings.clientSettingsResolver(5, 2)

	assertFilterSettings(
		t,
		clientSettings.FilterQueryParamNames,
		clientSettings.InitialState.FilterValues,
	)
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
	if clientSettings.RefreshDebounceMs != defaultRefreshDebounceMs {
		t.Errorf(
			"RefreshDebounceMsMismatch: got %d, want %d",
			clientSettings.RefreshDebounceMs, defaultRefreshDebounceMs,
		)
	}
	if clientSettings.QueryUrlTemplate != "/records?page={pageNumber}" {
		t.Errorf("QueryUrlTemplateMismatch: got %q, want %q", clientSettings.QueryUrlTemplate, "/records?page={pageNumber}")
	}
	if clientSettings.MaxVisibleRows != 3 {
		t.Errorf("MaxVisibleRowsMismatch: got %d, want 3", clientSettings.MaxVisibleRows)
	}

	overrideSettings := settings
	overrideSettings.RefreshDebounceMs = 750
	overrideClientSettings := overrideSettings.clientSettingsResolver(5, 1)
	if overrideClientSettings.RefreshDebounceMs != 750 {
		t.Errorf("RefreshDebounceMsMismatch: got %d, want 750", overrideClientSettings.RefreshDebounceMs)
	}
}

func TestDataTableDensityPaddingResolvers(t *testing.T) {
	denseSettings := DataTableSettings[dataTableTestRecord]{
		Density: DataTableDensityDense,
	}
	if actual := denseSettings.cellPaddingClassesResolver(); actual != "px-2 py-1.5" {
		t.Errorf("DenseCellPaddingMismatch: got %q", actual)
	}
	if actual := denseSettings.headerPaddingClassesResolver(); actual != "px-2 py-1.5" {
		t.Errorf("DenseHeaderPaddingMismatch: got %q", actual)
	}

	comfortableSettings := DataTableSettings[dataTableTestRecord]{
		Density: DataTableDensityComfortable,
	}
	if actual := comfortableSettings.cellPaddingClassesResolver(); actual != "px-3 py-3" {
		t.Errorf("ComfortableCellPaddingMismatch: got %q", actual)
	}
	if actual := comfortableSettings.headerPaddingClassesResolver(); actual != "px-3 py-2" {
		t.Errorf("ComfortableHeaderPaddingMismatch: got %q", actual)
	}
}

func dataTableRenderedComponent(t *testing.T, component templ.Component) string {
	t.Helper()
	if component == nil {
		return ""
	}
	var buffer bytes.Buffer
	if renderErr := component.Render(context.Background(), &buffer); renderErr != nil {
		t.Fatalf("ComponentRenderFailed: %v", renderErr)
	}
	return buffer.String()
}

func TestDataTableSearchBoxResolver(t *testing.T) {
	providedSearchBox := templ.Raw("<input name=\"provided-search\"/>")
	providedSettings := DataTableSettings[dataTableTestRecord]{
		SearchBox:        providedSearchBox,
		QueryUrlTemplate: "/records?search={search}",
	}
	providedRendered := dataTableRenderedComponent(t, providedSettings.searchBoxResolver())
	if providedRendered != "<input name=\"provided-search\"/>" {
		t.Errorf("ProvidedSearchBoxNotReturned: %q", providedRendered)
	}

	templatedSettings := DataTableSettings[dataTableTestRecord]{
		QueryUrlTemplate: "/records?search={search}",
	}
	if actual := templatedSettings.searchBoxResolver(); actual == nil {
		t.Errorf("TemplatedQueryUrlShouldYieldDefaultSearchBox")
	}

	plainSettings := DataTableSettings[dataTableTestRecord]{
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	if actual := plainSettings.searchBoxResolver(); actual != nil {
		t.Errorf("PlainQueryUrlShouldYieldNoSearchBox")
	}
}

func TestDataTableCheckboxShapeResolver(t *testing.T) {
	testCases := []struct {
		name          string
		providedShape string
		expectedShape string
	}{
		{name: "default", providedShape: "", expectedShape: uiToolset.ShapeSquare},
		{name: "provided shape wins", providedShape: uiToolset.ShapeCircular, expectedShape: uiToolset.ShapeCircular},
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

func TestDataTableRootClassesResolver(t *testing.T) {
	defaultSettings := DataTableSettings[dataTableTestRecord]{}
	defaultClasses := defaultSettings.rootClassesResolver()
	if !strings.Contains(defaultClasses, "rounded-md") {
		t.Errorf("RootClassesMissingBase: %q", defaultClasses)
	}
	if strings.Contains(defaultClasses, "overflow-hidden") {
		t.Errorf("RootClassesClipChildDropdowns: %q", defaultClasses)
	}

	sizedSettings := DataTableSettings[dataTableTestRecord]{
		MinWidthClass: "min-w-96",
		MaxWidthClass: "max-w-7xl",
	}
	sizedClasses := sizedSettings.rootClassesResolver()
	for _, expectedClass := range []string{"min-w-96", "max-w-7xl"} {
		if !strings.Contains(sizedClasses, expectedClass) {
			t.Errorf("RootClassesMissing(%q): %q", expectedClass, sizedClasses)
		}
	}
}

func TestDataTableStickyHeaderClassesResolver(t *testing.T) {
	defaultSettings := DataTableSettings[dataTableTestRecord]{}
	defaultClasses := defaultSettings.stickyHeaderClassesResolver()
	for _, expectedClass := range []string{"sticky", "top-0", "z-10", "backdrop-blur-md", "bg-neutral-950/20"} {
		if !strings.Contains(defaultClasses, expectedClass) {
			t.Errorf("StickyHeaderClassesMissing(%q): %q", expectedClass, defaultClasses)
		}
	}

	coloredSettings := DataTableSettings[dataTableTestRecord]{
		StickyHeaderBackgroundColor: "neutral-50/5",
	}
	coloredClasses := coloredSettings.stickyHeaderClassesResolver()
	if !strings.Contains(coloredClasses, "bg-neutral-50/5") {
		t.Errorf("StickyHeaderClassesMissingCustomBackground: %q", coloredClasses)
	}
	if strings.Contains(coloredClasses, "bg-neutral-950/20") {
		t.Errorf("StickyHeaderClassesKeepDefaultBackground: %q", coloredClasses)
	}

	staticSettings := DataTableSettings[dataTableTestRecord]{IsHeaderStatic: true}
	if actualClasses := staticSettings.stickyHeaderClassesResolver(); actualClasses != "" {
		t.Errorf("StaticHeaderClassesNotEmpty: %q", actualClasses)
	}
}

func TestDataTableRoundsScrollContainerTopWhenBodyIsFirstChild(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		Id:      "records-table",
		Columns: []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
	}
	scrollContainerClassPattern := regexp.MustCompile(`data-ui-data-table-scroll class="([^"]*)"`)
	classFromRender := func(settings DataTableSettings[dataTableTestRecord]) string {
		var buffer bytes.Buffer
		renderErr := DataTable(settings).Render(context.Background(), &buffer)
		if renderErr != nil {
			t.Fatalf("DataTableRenderFailed: %v", renderErr)
		}
		matches := scrollContainerClassPattern.FindStringSubmatch(buffer.String())
		if len(matches) < 2 {
			t.Fatalf("RenderedHtmlMissingScrollContainer")
		}
		return matches[1]
	}

	if !strings.Contains(classFromRender(settings), "rounded-t-md") {
		t.Errorf("ScrollContainerMissingTopRadiusOnFirstChildBody")
	}

	settings.Filters = []FilterSettings{{Key: "status", Label: "Status", Kind: FilterKindEnumSelect}}
	if strings.Contains(classFromRender(settings), "rounded-t-md") {
		t.Errorf("ScrollContainerKeepsTopRadiusWithFilterBar")
	}
}

func TestDataTableRendersItemsPerPageDropdownBackground(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		Id:         "records-table",
		Columns:    []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
		ItemsTotal: 240,
	}
	var buffer bytes.Buffer
	renderErr := DataTable(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DataTableRenderFailed: %v", renderErr)
	}
	if !strings.Contains(buffer.String(), "bg-neutral-800/95") {
		t.Errorf("RenderedHtmlMissingDefaultItemsPerPageDropdownBackground")
	}

	settings.ItemsPerPageDropdownBackgroundColor = "emerald-900"
	coloredBuffer := bytes.Buffer{}
	renderErr = DataTable(settings).Render(context.Background(), &coloredBuffer)
	if renderErr != nil {
		t.Fatalf("DataTableRenderFailed: %v", renderErr)
	}
	coloredHtml := coloredBuffer.String()
	if !strings.Contains(coloredHtml, "bg-emerald-900") {
		t.Errorf("RenderedHtmlMissingCustomItemsPerPageDropdownBackground")
	}
	if !strings.Contains(coloredHtml, "bg-neutral-950/20") {
		t.Errorf("RenderedHtmlMissingStickyHeaderBackground")
	}
}

func TestDataTableScrollContainerClassesResolver(t *testing.T) {
	defaultSettings := DataTableSettings[dataTableTestRecord]{}
	defaultClasses := defaultSettings.scrollContainerClassesResolver()
	for _, expectedClass := range []string{"overflow-auto", "max-h-128"} {
		if !strings.Contains(defaultClasses, expectedClass) {
			t.Errorf("ScrollContainerClassesMissing(%q): %q", expectedClass, defaultClasses)
		}
	}

	staticSettings := DataTableSettings[dataTableTestRecord]{IsHeaderStatic: true}
	staticClasses := staticSettings.scrollContainerClassesResolver()
	if !strings.Contains(staticClasses, "overflow-x-auto") {
		t.Errorf("StaticScrollContainerClassesMissingOverflowX: %q", staticClasses)
	}
	if strings.Contains(staticClasses, "max-h-128") {
		t.Errorf("StaticScrollContainerClassesCarryDefaultMaxHeight: %q", staticClasses)
	}

	explicitHeightSettings := DataTableSettings[dataTableTestRecord]{
		MaxHeightClass: "max-h-64",
		MinHeightClass: "min-h-32",
	}
	explicitHeightClasses := explicitHeightSettings.scrollContainerClassesResolver()
	for _, expectedClass := range []string{"max-h-64", "min-h-32"} {
		if !strings.Contains(explicitHeightClasses, expectedClass) {
			t.Errorf("ScrollContainerClassesMissing(%q): %q", expectedClass, explicitHeightClasses)
		}
	}
	if strings.Contains(explicitHeightClasses, "max-h-128") {
		t.Errorf("ScrollContainerClassesCarryDefaultMaxHeight: %q", explicitHeightClasses)
	}

	measuredRowsSettings := DataTableSettings[dataTableTestRecord]{MaxVisibleRows: 6}
	measuredRowsClasses := measuredRowsSettings.scrollContainerClassesResolver()
	if !strings.Contains(measuredRowsClasses, "max-h-128") {
		t.Errorf("MeasuredRowsFallbackMissing: %q", measuredRowsClasses)
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
			"IdChangedWithDifferentRows: %q != %q", rowVariantId, derivedId,
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

	if !strings.HasPrefix(derivedId, "dataTable-") {
		t.Errorf("IdMissingComponentPrefix: got %q, want %q", derivedId, "dataTable-")
	}
}

func TestDataTableWiresOneBasedPageDisplayIntoPageStripBuilder(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		ItemsPerPage:                 10,
		ItemsTotal:                   240,
		ShouldUseOneBasedPageDisplay: true,
	}
	var buffer bytes.Buffer
	renderErr := DataTable(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DataTableRenderFailed: %v", renderErr)
	}
	oneBasedStripCallPattern := regexp.MustCompile(
		`PaginationPageStripBuilder\(pageNumber, .+, 1\)`,
	)
	if !oneBasedStripCallPattern.MatchString(buffer.String()) {
		t.Errorf("RenderedHtmlMissingOneBasedPageStripCall")
	}
}

func TestDataTableRendersItsSettingsScriptInsideItsRoot(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		Id:               "records-table",
		Columns:          []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := DataTable(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DataTableRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	rootIndex := strings.Index(renderedHtml, `id="records-table"`)
	settingsIndex := strings.Index(renderedHtml, `id="records-table-settings"`)
	tableIndex := strings.Index(renderedHtml, `<div data-ui-data-table>`)
	if rootIndex < 0 || settingsIndex < 0 || tableIndex < 0 {
		t.Fatalf(
			"RenderedHtmlMissingMarker: root=%d settings=%d table=%d",
			rootIndex, settingsIndex, tableIndex,
		)
	}
	if settingsIndex < rootIndex || settingsIndex > tableIndex {
		t.Errorf(
			"SettingsScriptRenderedOutsideTheRoot: root=%d settings=%d table=%d",
			rootIndex, settingsIndex, tableIndex,
		)
	}
}

func TestDataTableRefreshRegionRendersOnlyTheSwappedRegion(t *testing.T) {
	settings := DataTableSettings[dataTableTestRecord]{
		Id:      "records-table",
		Columns: []DataTableColumnSettings[dataTableTestRecord]{{Label: "Name"}},
		Filters: []FilterSettings{
			{Key: "status", Label: "Status", Kind: FilterKindEnumSelect},
		},
		QueryUrlTemplate: "/records?page={pageNumber}&search={searchQuery}",
	}
	var buffer bytes.Buffer
	renderErr := DataTableRefreshRegion(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DataTableRefreshRegionRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()

	if !strings.HasPrefix(renderedHtml, `<div data-ui-data-table>`) {
		t.Errorf("RefreshRegionDoesNotStartWithTheSwappedRegion: %.60q", renderedHtml)
	}
	if !strings.Contains(renderedHtml, `name="records-table-itemsPerPage"`) {
		t.Errorf("RefreshRegionMissingPagination")
	}
	unexpectedMarkers := []string{
		`id="records-table"`,
		`id="records-table-settings"`,
		`Alpine.data("dataTable"`,
		`name="search"`,
		`name="records-table-status"`,
	}
	for _, unexpectedMarker := range unexpectedMarkers {
		if strings.Contains(renderedHtml, unexpectedMarker) {
			t.Errorf("RefreshRegionRendersOutsideMarkup: %s", unexpectedMarker)
		}
	}
}

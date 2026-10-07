package uiStructural

import (
	"bytes"
	"context"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

type carouselTestItem struct {
	Id   string
	Name string
}

func carouselTestRenderer(item carouselTestItem) templ.Component {
	return templ.Raw("<span>" + item.Name + "</span>")
}

func TestCarouselClientSettingsResolver(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		ItemsPerView: CarouselItemsPerViewSettings{Base: 1, Sm: 2, Lg: 3},
		Filters: []FilterSettings{
			{Key: "status", Kind: FilterKindEnumSelect},
			{Key: "cpu", Kind: FilterKindNumberRange, QueryParamName: "cores"},
		},
		InitialFilterValues:     map[string]any{"status": "running"},
		InitialSearchQuery:      "alpha",
		IsAutoplay:              true,
		IsAutoplayPausedOnHover: true,
		IsSwipeEnabled:          true,
		QueryUrlTemplate:        "/records?page={pageNumber}",
		RefreshOnEvents:         []string{"refresh:records"},
	}
	clientSettings := settings.clientSettingsResolver(6, 2)

	expectedItemsPerView := CarouselItemsPerViewSettings{Base: 1, Sm: 2, Lg: 3}
	if !reflect.DeepEqual(clientSettings.ItemsPerView, expectedItemsPerView) {
		t.Errorf(
			"ItemsPerViewMismatch: got %+v, want %+v",
			clientSettings.ItemsPerView, expectedItemsPerView,
		)
	}
	if clientSettings.AutoplayIntervalMs != carouselDefaultAutoplayIntervalMs {
		t.Errorf(
			"AutoplayIntervalMsMismatch: got %d, want %d",
			clientSettings.AutoplayIntervalMs, carouselDefaultAutoplayIntervalMs,
		)
	}
	if !clientSettings.IsAutoplay || !clientSettings.IsAutoplayPausedOnHover ||
		!clientSettings.IsSwipeEnabled {
		t.Errorf("AutoplayOrSwipeFlagsMismatch: %+v", clientSettings)
	}
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
	if clientSettings.InitialState.ItemsPerPage != 6 {
		t.Errorf(
			"ItemsPerPageMismatch: got %d, want 6",
			clientSettings.InitialState.ItemsPerPage,
		)
	}
	if clientSettings.InitialState.PageNumber != 2 {
		t.Errorf(
			"PageNumberMismatch: got %d, want 2",
			clientSettings.InitialState.PageNumber,
		)
	}
	if clientSettings.InitialState.SearchQuery != "alpha" {
		t.Errorf(
			"SearchQueryMismatch: got %q, want %q",
			clientSettings.InitialState.SearchQuery, "alpha",
		)
	}
	if clientSettings.RefreshDebounceMs != defaultRefreshDebounceMs {
		t.Errorf(
			"RefreshDebounceMsMismatch: got %d, want %d",
			clientSettings.RefreshDebounceMs, defaultRefreshDebounceMs,
		)
	}

	overrideSettings := settings
	overrideSettings.RefreshDebounceMs = 750
	overrideSettings.AutoplayIntervalMs = 2500
	overrideSettings.ItemsPerView = CarouselItemsPerViewSettings{Base: 4}
	overrideClientSettings := overrideSettings.clientSettingsResolver(6, 1)
	if overrideClientSettings.RefreshDebounceMs != 750 {
		t.Errorf(
			"RefreshDebounceMsMismatch: got %d, want 750",
			overrideClientSettings.RefreshDebounceMs,
		)
	}
	if overrideClientSettings.AutoplayIntervalMs != 2500 {
		t.Errorf(
			"AutoplayIntervalMsMismatch: got %d, want 2500",
			overrideClientSettings.AutoplayIntervalMs,
		)
	}
	if overrideClientSettings.ItemsPerView.Base != 4 {
		t.Errorf(
			"ItemsPerViewBaseMismatch: got %d, want 4",
			overrideClientSettings.ItemsPerView.Base,
		)
	}
}

func TestCarouselIdResolver(t *testing.T) {
	explicitSettings := CarouselSettings[carouselTestItem]{Id: "records-carousel"}
	if actualId := explicitSettings.idResolver(); actualId != "records-carousel" {
		t.Errorf("IdMismatch: got %q, want %q", actualId, "records-carousel")
	}

	settings := CarouselSettings[carouselTestItem]{
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	derivedId := settings.idResolver()
	repeatedId := settings.idResolver()
	if repeatedId != derivedId {
		t.Errorf("IdResolverNotStable: %q != %q", repeatedId, derivedId)
	}
	if !strings.HasPrefix(derivedId, "carousel-") {
		t.Errorf("IdMissingComponentPrefix: got %q, want %q", derivedId, "carousel-")
	}

	filteredSettings := settings
	filteredSettings.Filters = []FilterSettings{
		{Key: "status", Kind: FilterKindEnumSelect},
	}
	filteredId := filteredSettings.idResolver()
	if filteredId == derivedId {
		t.Errorf("IdCollidedForDifferentFilters: %q", filteredId)
	}
}

func TestCarouselItemsPerViewResolver(t *testing.T) {
	providedSettings := CarouselSettings[carouselTestItem]{
		ItemsPerView: CarouselItemsPerViewSettings{Base: 4, Sm: 6},
	}
	actual := providedSettings.itemsPerViewResolver()
	expected := CarouselItemsPerViewSettings{Base: 4, Sm: 6}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ItemsPerViewMismatch: got %+v, want %+v", actual, expected)
	}

	emptySettings := CarouselSettings[carouselTestItem]{}
	if base := emptySettings.itemsPerViewResolver().Base; base != carouselDefaultItemsPerView {
		t.Errorf(
			"ItemsPerViewBaseMismatch: got %d, want %d",
			base, carouselDefaultItemsPerView,
		)
	}
}

func TestCarouselItemsPerPageSizeChoicesResolver(t *testing.T) {
	providedSettings := CarouselSettings[carouselTestItem]{
		ItemsPerPageSizeChoices: []ItemsPerPage{6, 12},
	}
	actualProvidedChoices := itemsPerPageSizeChoicesResolver(
		providedSettings.ItemsPerPageSizeChoices,
	)
	expectedProvidedChoices := []uint{6, 12}
	if !reflect.DeepEqual(actualProvidedChoices, expectedProvidedChoices) {
		t.Errorf(
			"ItemsPerPageSizeChoicesMismatch: got %v, want %v",
			actualProvidedChoices, expectedProvidedChoices,
		)
	}

	emptySettings := CarouselSettings[carouselTestItem]{}
	actualDefaultChoices := itemsPerPageSizeChoicesResolver(
		emptySettings.ItemsPerPageSizeChoices,
	)
	if !reflect.DeepEqual(actualDefaultChoices, paginationDefaultItemsPerPageSizeChoices) {
		t.Errorf(
			"ItemsPerPageSizeChoicesMismatch: got %v, want %v",
			actualDefaultChoices, paginationDefaultItemsPerPageSizeChoices,
		)
	}
}

func TestCarouselItemsPerPageResolver(t *testing.T) {
	testCases := []struct {
		name                    string
		itemsPerPageSizeChoices []uint
		itemsPerPage            ItemsPerPage
		expectedItemsPerPage    uint
	}{
		{
			name:                    "explicit items per page wins",
			itemsPerPageSizeChoices: []uint{6, 12},
			itemsPerPage:            12,
			expectedItemsPerPage:    12,
		},
		{
			name:                    "zero falls back to the first choice",
			itemsPerPageSizeChoices: []uint{6, 12},
			itemsPerPage:            0,
			expectedItemsPerPage:    6,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualItemsPerPage := itemsPerPageResolver(
				testCase.itemsPerPage, testCase.itemsPerPageSizeChoices,
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

func TestCarouselDotsColorResolvers(t *testing.T) {
	defaultSettings := CarouselSettings[carouselTestItem]{}
	if actual := defaultSettings.dotsActiveColorResolver(); actual != carouselDefaultDotsActiveColor {
		t.Errorf("DotsActiveColorMismatch: got %q", actual)
	}
	if actual := defaultSettings.dotsInactiveColorResolver(); actual != carouselDefaultDotsInactiveColor {
		t.Errorf("DotsInactiveColorMismatch: got %q", actual)
	}

	providedSettings := CarouselSettings[carouselTestItem]{
		DotsActiveColor:   "emerald-500",
		DotsInactiveColor: "neutral-50/30",
	}
	if actual := providedSettings.dotsActiveColorResolver(); actual != "emerald-500" {
		t.Errorf("DotsActiveColorMismatch: got %q", actual)
	}
	if actual := providedSettings.dotsInactiveColorResolver(); actual != "neutral-50/30" {
		t.Errorf("DotsInactiveColorMismatch: got %q", actual)
	}
}

func TestCarouselSearchBoxAlignmentResolver(t *testing.T) {
	defaultSettings := CarouselSettings[carouselTestItem]{}
	if actual := defaultSettings.searchBoxAlignmentResolver(); actual != HorizontalAlignmentCenter {
		t.Errorf("SearchBoxAlignmentMismatch: got %q, want center", actual)
	}

	providedSettings := CarouselSettings[carouselTestItem]{
		SearchBoxAlignment: HorizontalAlignmentRight,
	}
	if actual := providedSettings.searchBoxAlignmentResolver(); actual != HorizontalAlignmentRight {
		t.Errorf("SearchBoxAlignmentMismatch: got %q, want right", actual)
	}
}

func TestCarouselSearchBoxPositionResolver(t *testing.T) {
	defaultSettings := CarouselSettings[carouselTestItem]{}
	if actual := defaultSettings.searchBoxPositionResolver(); actual != CarouselSearchBoxPositionTop {
		t.Errorf("SearchBoxPositionMismatch: got %q, want top", actual)
	}

	providedSettings := CarouselSettings[carouselTestItem]{
		SearchBoxPosition: CarouselSearchBoxPositionBottom,
	}
	if actual := providedSettings.searchBoxPositionResolver(); actual != CarouselSearchBoxPositionBottom {
		t.Errorf("SearchBoxPositionMismatch: got %q, want bottom", actual)
	}
}

func TestCarouselSurfaceClassesResolver(t *testing.T) {
	defaultSettings := CarouselSettings[carouselTestItem]{}
	defaultClassFields := strings.Fields(defaultSettings.surfaceClassesResolver())
	for _, expectedClass := range []string{"bg-neutral-50/2.5", "rounded"} {
		if !slices.Contains(defaultClassFields, expectedClass) {
			t.Errorf(
				"SurfaceClassesMissingDefault: %q in %q",
				expectedClass, defaultClassFields,
			)
		}
	}

	lgSettings := CarouselSettings[carouselTestItem]{
		BorderRadius: CarouselBorderRadiusLg,
	}
	if lgClasses := lgSettings.surfaceClassesResolver(); !strings.Contains(lgClasses, "rounded-lg") {
		t.Errorf("SurfaceClassesMissingLg: %q", lgClasses)
	}

	providedSettings := CarouselSettings[carouselTestItem]{
		BackgroundColor: "neutral-800/50",
		BorderRadius:    CarouselBorderRadiusXl,
		ShadowSize:      CarouselShadowSizeLg,
		RingColor:       "secondary-500/30",
		RingThickness:   CarouselRingThicknessMd,
	}
	providedClasses := providedSettings.surfaceClassesResolver()
	for _, expectedClass := range []string{
		"bg-neutral-800/50", "rounded-xl", "shadow-lg", "ring-2", "ring-secondary-500/30",
	} {
		if !strings.Contains(providedClasses, expectedClass) {
			t.Errorf("SurfaceClassesMissing: %q in %q", expectedClass, providedClasses)
		}
	}
}

func TestCarouselArrowsClassesResolver(t *testing.T) {
	outsideClasses := carouselArrowsClassesResolver(
		CarouselArrowsPositionOutside, CarouselArrowsShapeRounded,
		CarouselArrowsSizeSm, "", "",
	)
	for _, expectedClass := range []string{"rounded", "h-7", "w-7", "bg-neutral-50/7.5"} {
		if !strings.Contains(outsideClasses, expectedClass) {
			t.Errorf("OutsideArrowsClassesMissing: %q in %q", expectedClass, outsideClasses)
		}
	}

	insideClasses := carouselArrowsClassesResolver(
		CarouselArrowsPositionInside, CarouselArrowsShapeSquare,
		CarouselArrowsSizeLg, "secondary-500/20", "secondary-100",
	)
	for _, expectedClass := range []string{
		"rounded-none", "h-9", "w-9", "bg-secondary-500/20",
		"hover:brightness-125", "text-secondary-100",
	} {
		if !strings.Contains(insideClasses, expectedClass) {
			t.Errorf("InsideArrowsClassesMissing: %q in %q", expectedClass, insideClasses)
		}
	}
}

func TestCarouselDotsSizeClassesResolver(t *testing.T) {
	testCases := []struct {
		name             string
		dotsSize         string
		expectedHeight   string
		expectedInactive string
		expectedActive   string
	}{
		{name: "default", dotsSize: "", expectedHeight: "h-2", expectedInactive: "w-2", expectedActive: "w-5"},
		{name: "sm", dotsSize: CarouselDotsSizeSm, expectedHeight: "h-1.5", expectedInactive: "w-1.5", expectedActive: "w-4"},
		{name: "lg", dotsSize: CarouselDotsSizeLg, expectedHeight: "h-2.5", expectedInactive: "w-2.5", expectedActive: "w-6"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			height, inactive, active := carouselDotsSizeClassesResolver(testCase.dotsSize)
			if height != testCase.expectedHeight || inactive != testCase.expectedInactive ||
				active != testCase.expectedActive {
				t.Errorf(
					"DotsSizeClassesMismatch: got (%s,%s,%s), want (%s,%s,%s)",
					height, inactive, active,
					testCase.expectedHeight, testCase.expectedInactive, testCase.expectedActive,
				)
			}
		})
	}
}

func TestCarouselItemClassesResolver(t *testing.T) {
	emptyClasses := carouselItemClassesResolver("", "", "", "", "", "")
	if emptyClasses != "" {
		t.Errorf("ItemClassesNotEmpty: %q", emptyClasses)
	}

	providedClasses := carouselItemClassesResolver(
		"neutral-50/5", "lg", CarouselPaddingSizeMd,
		"neutral-50/10", CarouselRingThicknessXs, CarouselShadowSizeSm,
	)
	for _, expectedClass := range []string{"bg-neutral-50/5", "rounded-lg", "p-5", "ring-1", "ring-neutral-50/10", "shadow-sm"} {
		if !strings.Contains(providedClasses, expectedClass) {
			t.Errorf("ItemClassesMissing: %q in %q", expectedClass, providedClasses)
		}
	}
}

func TestCarouselPaginationAriaLabelResolver(t *testing.T) {
	defaultSettings := CarouselSettings[carouselTestItem]{}
	if actual := defaultSettings.paginationAriaLabelResolver(); actual != "Carousel pagination" {
		t.Errorf("PaginationAriaLabelMismatch: got %q", actual)
	}

	providedSettings := CarouselSettings[carouselTestItem]{
		PaginationAriaLabel: "Servers carousel pagination",
	}
	if actual := providedSettings.paginationAriaLabelResolver(); actual != "Servers carousel pagination" {
		t.Errorf("PaginationAriaLabelMismatch: got %q", actual)
	}
}

func TestCarouselRendersOneTrackItemPerItem(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		ItemRenderer: carouselTestRenderer,
		Items: []carouselTestItem{
			{Id: "one", Name: "alpha"},
			{Id: "two", Name: "bravo"},
			{Id: "three", Name: "charlie"},
		},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := Carousel(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CarouselRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedName := range []string{"alpha", "bravo", "charlie"} {
		if !regexp.MustCompile(`<span>` + expectedName + `</span>`).MatchString(renderedHtml) {
			t.Errorf("RenderedHtmlMissingItem: %s", expectedName)
		}
	}
	if !regexp.MustCompile(`data-ui-carousel-track`).MatchString(renderedHtml) {
		t.Errorf("RenderedHtmlMissingTrack")
	}
	if !strings.Contains(renderedHtml, `data-ui-carousel-pagination`) {
		t.Errorf("RenderedHtmlMissingPaginationRegion")
	}
	trackIndex := strings.Index(renderedHtml, `<div data-ui-carousel>`)
	paginationIndex := strings.Index(renderedHtml, `data-ui-carousel-pagination class=`)
	if paginationIndex < trackIndex {
		t.Errorf("PaginationRegionRenderedBeforeTrackRegion")
	}
}

func TestCarouselRendersHtmxRefreshMarkup(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		Id:               "records-carousel",
		ItemRenderer:     carouselTestRenderer,
		Items:            []carouselTestItem{{Id: "one", Name: "alpha"}},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := Carousel(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CarouselRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if !strings.Contains(renderedHtml, `hx-sync="this:replace"`) {
		t.Errorf("RenderedHtmlMissingRootSyncStrategy")
	}
	if !strings.Contains(
		renderedHtml, `hx-swap-oob="outerHTML:#records-carousel-pagination"`,
	) {
		t.Errorf("RenderedHtmlMissingPaginationOutOfBandSwap")
	}
}

func TestCarouselPaginationOutOfBandSwapTargetsItsOwnComponent(t *testing.T) {
	var buffer bytes.Buffer
	for _, carouselId := range []string{"first-carousel", "second-carousel"} {
		settings := CarouselSettings[carouselTestItem]{
			Id:               carouselId,
			ItemRenderer:     carouselTestRenderer,
			Items:            []carouselTestItem{{Id: "one", Name: "alpha"}},
			QueryUrlTemplate: "/records?page={pageNumber}",
		}
		renderErr := Carousel(settings).Render(context.Background(), &buffer)
		if renderErr != nil {
			t.Fatalf("CarouselRenderFailed: %v", renderErr)
		}
	}
	renderedHtml := buffer.String()
	if strings.Contains(renderedHtml, `outerHTML:[data-ui-carousel-pagination]`) {
		t.Errorf("PaginationOutOfBandSwapIsShared")
	}
	for _, carouselId := range []string{"first-carousel", "second-carousel"} {
		paginationId := carouselId + "-pagination"
		if !strings.Contains(renderedHtml, `id="`+paginationId+`"`) {
			t.Errorf("PaginationIdMissing: %s", paginationId)
		}
		if !strings.Contains(
			renderedHtml, `hx-swap-oob="outerHTML:#`+paginationId+`"`,
		) {
			t.Errorf("PaginationOutOfBandSwapDoesNotTargetOwnComponent: %s",
				paginationId)
		}
	}
}

func TestCarouselRendersEmptyStateWithoutATrack(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		ItemRenderer:     carouselTestRenderer,
		Items:            []carouselTestItem{},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := Carousel(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CarouselRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if !regexp.MustCompile(`No items found\.`).MatchString(renderedHtml) {
		t.Errorf("RenderedHtmlMissingDefaultEmptyState")
	}
	if regexp.MustCompile(`data-ui-carousel-track class=`).MatchString(renderedHtml) {
		t.Errorf("RenderedHtmlUnexpectedTrackForEmptyItems")
	}
}

func TestCarouselRendersFilterBarBetweenTrackAndPagination(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		ItemRenderer: carouselTestRenderer,
		Items:        []carouselTestItem{{Id: "one", Name: "alpha"}},
		Filters: []FilterSettings{
			{Key: "status", Label: "Status", Kind: FilterKindEnumSelect},
		},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := Carousel(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CarouselRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	trackIndex := strings.Index(renderedHtml, `<div data-ui-carousel>`)
	filterIndex := strings.Index(renderedHtml, "Clear filters")
	paginationIndex := strings.Index(renderedHtml, `data-ui-carousel-pagination class=`)
	if trackIndex >= filterIndex || filterIndex >= paginationIndex {
		t.Errorf(
			"FilterBarNotBetweenTrackAndPagination: track=%d filter=%d pagination=%d",
			trackIndex, filterIndex, paginationIndex,
		)
	}
}

func TestCarouselRendersItsSettingsScriptInsideItsRoot(t *testing.T) {
	settings := CarouselSettings[carouselTestItem]{
		Id:               "records-carousel",
		ItemRenderer:     carouselTestRenderer,
		Items:            []carouselTestItem{{Id: "one", Name: "alpha"}},
		QueryUrlTemplate: "/records?page={pageNumber}",
	}
	var buffer bytes.Buffer
	renderErr := Carousel(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CarouselRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	rootIndex := strings.Index(renderedHtml, `id="records-carousel"`)
	settingsIndex := strings.Index(renderedHtml, `id="records-carousel-settings"`)
	trackIndex := strings.Index(renderedHtml, `<div data-ui-carousel>`)
	if rootIndex < 0 || settingsIndex < 0 || trackIndex < 0 {
		t.Fatalf(
			"RenderedHtmlMissingMarker: root=%d settings=%d track=%d",
			rootIndex, settingsIndex, trackIndex,
		)
	}
	if settingsIndex < rootIndex || settingsIndex > trackIndex {
		t.Errorf(
			"SettingsScriptRenderedOutsideTheRoot: root=%d settings=%d track=%d",
			rootIndex, settingsIndex, trackIndex,
		)
	}
}

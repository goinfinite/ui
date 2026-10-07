package uiStructural

import (
	"fmt"

	"github.com/a-h/templ"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

const (
	CarouselUrlPlaceholderPageNumber   string = DataTableUrlPlaceholderPageNumber
	CarouselUrlPlaceholderItemsPerPage string = DataTableUrlPlaceholderItemsPerPage
	CarouselUrlPlaceholderSearch       string = DataTableUrlPlaceholderSearch

	CarouselArrowsPositionInside  string = "inside"
	CarouselArrowsPositionOutside string = "outside"

	CarouselDotsPositionTop    string = "top"
	CarouselDotsPositionBottom string = "bottom"

	CarouselSearchBoxPositionTop    string = "top"
	CarouselSearchBoxPositionBottom string = "bottom"

	CarouselBorderRadiusNone string = "none"
	CarouselBorderRadiusXs   string = "xs"
	CarouselBorderRadiusSm   string = "sm"
	CarouselBorderRadiusMd   string = "md"
	CarouselBorderRadiusLg   string = "lg"
	CarouselBorderRadiusXl   string = "xl"

	CarouselShadowSizeNone string = "none"
	CarouselShadowSizeXs   string = "xs"
	CarouselShadowSizeSm   string = "sm"
	CarouselShadowSizeMd   string = "md"
	CarouselShadowSizeLg   string = "lg"
	CarouselShadowSizeXl   string = "xl"

	CarouselRingThicknessXs string = "xs"
	CarouselRingThicknessSm string = "sm"
	CarouselRingThicknessMd string = "md"
	CarouselRingThicknessLg string = "lg"
	CarouselRingThicknessXl string = "xl"

	CarouselPaddingSizeNone string = "none"
	CarouselPaddingSizeXs   string = "xs"
	CarouselPaddingSizeSm   string = "sm"
	CarouselPaddingSizeMd   string = "md"
	CarouselPaddingSizeLg   string = "lg"
	CarouselPaddingSizeXl   string = "xl"

	CarouselGapSizeNone string = "none"
	CarouselGapSizeXs   string = "xs"
	CarouselGapSizeSm   string = "sm"
	CarouselGapSizeMd   string = "md"
	CarouselGapSizeLg   string = "lg"
	CarouselGapSizeXl   string = "xl"

	CarouselArrowsSizeSm string = "sm"
	CarouselArrowsSizeMd string = "md"
	CarouselArrowsSizeLg string = "lg"

	CarouselArrowsShapeCircular string = "circular"
	CarouselArrowsShapeRounded  string = "rounded"
	CarouselArrowsShapeSquare   string = "square"

	CarouselDotsSizeSm string = "sm"
	CarouselDotsSizeMd string = "md"
	CarouselDotsSizeLg string = "lg"

	carouselDefaultPaginationAriaLabel string = "Carousel pagination"
	carouselDefaultAutoplayIntervalMs  uint   = 4000
	carouselDefaultItemsPerView        uint   = 1
	carouselDefaultDotsActiveColor     string = "secondary-500"
	carouselDefaultDotsInactiveColor   string = "neutral-50/20"
)

type CarouselItemsPerViewSettings struct {
	Base  uint `json:"base"`
	Sm    uint `json:"sm"`
	Md    uint `json:"md"`
	Lg    uint `json:"lg"`
	Xl    uint `json:"xl"`
	TwoXl uint `json:"twoXl"`
}

type CarouselSettings[Item any] struct {
	ItemRenderer func(item Item) templ.Component
	Items        []Item

	// OptionalFields
	ArrowsBackgroundColor            string
	ArrowsIconColor                  string
	ArrowsPosition                   string
	ArrowsShape                      string
	ArrowsSize                       string
	AutoplayIntervalMs               uint
	BackgroundColor                  string
	BorderRadius                     string
	DotsActiveColor                  string
	DotsInactiveColor                string
	DotsPosition                     string
	DotsSize                         string
	EmptyState                       templ.Component
	FilterDropdownBackgroundColor    string
	Filters                          []FilterSettings
	GapSize                          string
	Id                               string
	InitialFilterValues              map[string]any
	InitialSearchQuery               string
	IsAutoplay                       bool
	IsAutoplayPausedOnHover          bool
	IsPaginationHiddenWhenSinglePage bool
	IsSearchBoxCompact               bool
	IsSwipeEnabled                   bool
	ItemBackgroundColor              string
	ItemBorderRadius                 string
	ItemPaddingSize                  string
	ItemRingColor                    string
	ItemRingThickness                string
	ItemShadowSize                   string
	ItemsPerPage                     ItemsPerPage
	ItemsPerPageSizeChoices          []ItemsPerPage
	ItemsPerView                     CarouselItemsPerViewSettings
	ItemsTotal                       uint
	PageNumber                       uint
	PaddingSize                      string
	PaginationAriaLabel              string
	PagesTotal                       uint
	QueryUrlTemplate                 string
	RefreshDebounceMs                uint
	RefreshOnEvents                  []string
	RingColor                        string
	RingThickness                    string
	SearchBox                        templ.Component
	SearchBoxAlignment               HorizontalAlignment
	SearchBoxPosition                string
	ShadowSize                       string
	ShouldUseOneBasedPageDisplay     bool
	TextColor                        string
}

type carouselInitialState struct {
	FilterValues map[string]any `json:"filterValues"`
	ItemsPerPage uint           `json:"itemsPerPage"`
	PageNumber   uint           `json:"pageNumber"`
	SearchQuery  string         `json:"searchQuery"`
}

type carouselClientSettings struct {
	AutoplayIntervalMs      uint                         `json:"autoplayIntervalMs"`
	FilterQueryParamNames   map[string]string            `json:"filterQueryParamNames"`
	InitialState            carouselInitialState         `json:"initialState"`
	IsAutoplay              bool                         `json:"isAutoplay"`
	IsAutoplayPausedOnHover bool                         `json:"isAutoplayPausedOnHover"`
	IsSwipeEnabled          bool                         `json:"isSwipeEnabled"`
	ItemsPerView            CarouselItemsPerViewSettings `json:"itemsPerView"`
	QueryUrlTemplate        string                       `json:"queryUrlTemplate"`
	RefreshDebounceMs       uint                         `json:"refreshDebounceMs"`
	RefreshOnEvents         []string                     `json:"refreshOnEvents"`
}

func (settings CarouselSettings[Item]) itemsPerViewResolver() CarouselItemsPerViewSettings {
	resolved := settings.ItemsPerView
	if resolved.Base == 0 {
		resolved.Base = carouselDefaultItemsPerView
	}
	return resolved
}

func (settings CarouselSettings[Item]) autoplayIntervalMsResolver() uint {
	if settings.AutoplayIntervalMs > 0 {
		return settings.AutoplayIntervalMs
	}
	return carouselDefaultAutoplayIntervalMs
}

func (settings CarouselSettings[Item]) clientSettingsResolver(
	itemsPerPage, pageNumber uint,
) carouselClientSettings {
	return carouselClientSettings{
		AutoplayIntervalMs:    settings.autoplayIntervalMsResolver(),
		FilterQueryParamNames: filterQueryParamNamesResolver(settings.Filters),
		InitialState: carouselInitialState{
			FilterValues: initialFilterValuesResolver(
				settings.Filters, settings.InitialFilterValues,
			),
			ItemsPerPage: itemsPerPage,
			PageNumber:   pageNumber,
			SearchQuery:  settings.InitialSearchQuery,
		},
		IsAutoplay:              settings.IsAutoplay,
		IsAutoplayPausedOnHover: settings.IsAutoplayPausedOnHover,
		IsSwipeEnabled:          settings.IsSwipeEnabled,
		ItemsPerView:            settings.itemsPerViewResolver(),
		QueryUrlTemplate:        settings.QueryUrlTemplate,
		RefreshDebounceMs:       refreshDebounceResolver(settings.RefreshDebounceMs),
		RefreshOnEvents:         settings.RefreshOnEvents,
	}
}

func (settings CarouselSettings[Item]) carouselIdHashResolver() uint64 {
	idParts := []string{settings.QueryUrlTemplate}
	for _, filter := range settings.Filters {
		idParts = append(idParts, filter.Key, filter.QueryParamName)
	}
	return uiToolset.HashComponentIdParts(idParts...)
}

func (settings CarouselSettings[Item]) idResolver() string {
	if settings.Id != "" {
		return settings.Id
	}
	return fmt.Sprintf("carousel-%x", settings.carouselIdHashResolver())
}

func (settings CarouselSettings[Item]) paginationAriaLabelResolver() string {
	if settings.PaginationAriaLabel != "" {
		return settings.PaginationAriaLabel
	}
	return carouselDefaultPaginationAriaLabel
}

func (settings CarouselSettings[Item]) searchBoxAlignmentResolver() HorizontalAlignment {
	if settings.SearchBoxAlignment != "" {
		return settings.SearchBoxAlignment
	}
	return HorizontalAlignmentCenter
}

func (settings CarouselSettings[Item]) searchBoxPositionResolver() string {
	if settings.SearchBoxPosition != "" {
		return settings.SearchBoxPosition
	}
	return CarouselSearchBoxPositionTop
}

func (settings CarouselSettings[Item]) dotsActiveColorResolver() string {
	if settings.DotsActiveColor != "" {
		return settings.DotsActiveColor
	}
	return carouselDefaultDotsActiveColor
}

func (settings CarouselSettings[Item]) dotsInactiveColorResolver() string {
	if settings.DotsInactiveColor != "" {
		return settings.DotsInactiveColor
	}
	return carouselDefaultDotsInactiveColor
}

func (settings CarouselSettings[Item]) surfaceClassesResolver() string {
	surfaceClasses := "flex w-full min-w-0 flex-col overflow-hidden"
	surfaceClasses += " " + uiToolset.BackgroundColorClassResolver(
		settings.BackgroundColor, "bg-neutral-50/2.5",
	)
	surfaceClasses += " " + uiToolset.TextColorClassResolver(settings.TextColor, "")
	surfaceClasses += " " + uiToolset.BorderRadiusClassResolver(
		settings.BorderRadius, "rounded",
	)
	surfaceClasses += " " + uiToolset.ShadowClassResolver(settings.ShadowSize, "")
	surfaceClasses += " " + uiToolset.RingClassResolver(
		settings.RingColor, settings.RingThickness,
	)
	return surfaceClasses
}

func carouselItemClassesResolver(
	backgroundColor, borderRadius, paddingSize,
	ringColor, ringThickness, shadowSize string,
) string {
	itemClasses := ""
	if backgroundColor != "" {
		itemClasses += uiToolset.BackgroundColorClassResolver(backgroundColor, "") + " "
	}
	if borderRadius != "" {
		itemClasses += uiToolset.BorderRadiusClassResolver(borderRadius, "rounded") + " "
	}
	if paddingSize != "" {
		itemClasses += uiToolset.PaddingClassResolver(paddingSize, "p-5") + " "
	}
	if ringColor != "" {
		itemClasses += uiToolset.RingClassResolver(ringColor, ringThickness) + " "
	}
	if shadowSize != "" {
		itemClasses += uiToolset.ShadowClassResolver(shadowSize, "") + " "
	}
	return itemClasses
}

func carouselArrowPositionClassesResolver(
	isInside bool, insideOffsetClass string,
) string {
	if isInside {
		return "absolute top-1/2 z-10 -translate-y-1/2 " + insideOffsetClass
	}
	return "shrink-0 self-center"
}

func carouselArrowsClassesResolver(
	position, shape, size, backgroundColor, iconColor string,
) string {
	shapeClass := uiToolset.ShapeClassResolver(shape, "rounded-full")
	sizeClass := "h-8 w-8 text-base"
	switch size {
	case CarouselArrowsSizeSm:
		sizeClass = "h-7 w-7 text-sm"
	case CarouselArrowsSizeLg:
		sizeClass = "h-9 w-9 text-lg"
	}
	backgroundClass := "bg-neutral-50/7.5 hover:bg-neutral-50/12.5"
	if position == CarouselArrowsPositionInside {
		backgroundClass = "bg-neutral-900/70 hover:bg-neutral-900/90"
	}
	if backgroundColor != "" {
		backgroundClass = uiToolset.BackgroundColorClassResolver(backgroundColor, backgroundClass) + " hover:brightness-125"
	}
	iconColorClass := uiToolset.TextColorClassResolver(iconColor, "text-neutral-50")
	return "flex shrink-0 items-center justify-center " +
		shapeClass + " " + sizeClass + " " + backgroundClass + " " + iconColorClass
}

func carouselDotsSizeClassesResolver(
	dotsSize string,
) (heightClass, inactiveClass, activeClass string) {
	switch dotsSize {
	case CarouselDotsSizeSm:
		return "h-1.5", "w-1.5", "w-4"
	case CarouselDotsSizeLg:
		return "h-2.5", "w-2.5", "w-6"
	}
	return "h-2", "w-2", "w-5"
}

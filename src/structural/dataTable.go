package uiStructural

import (
	"fmt"

	"github.com/a-h/templ"
	uiForm "github.com/goinfinite/ui/src/form"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

type DataTableDensity string

const (
	DataTableDensityComfortable DataTableDensity = "comfortable"
	DataTableDensityDense       DataTableDensity = "dense"
)

type DataTableSortDirection string

const (
	DataTableSortDirectionAsc  DataTableSortDirection = "asc"
	DataTableSortDirectionDesc DataTableSortDirection = "desc"
)

const (
	DataTableUrlPlaceholderPageNumber    string = "{pageNumber}"
	DataTableUrlPlaceholderItemsPerPage  string = "{itemsPerPage}"
	DataTableUrlPlaceholderSortKey       string = "{sortKey}"
	DataTableUrlPlaceholderSortDirection string = "{sortDirection}"
	DataTableUrlPlaceholderSearch        string = "{search}"

	dataTableDefaultPaginationAriaLabel string = "Table pagination"
)

type DataTableColumnSettings[Row any] struct {
	Label        string
	CellRenderer func(row Row) templ.Component

	// OptionalFields
	Alignment     TextAlignment
	CellClass     string
	MaxWidthClass string
	MinWidthClass string
	SortKey       string
	WidthPercent  uint8
}

type DataTableSettings[Row any] struct {
	Columns []DataTableColumnSettings[Row]
	Rows    []Row

	// OptionalFields
	BulkActions                      templ.Component
	CheckboxCheckedColor             string
	CheckboxShape                    string
	CheckboxSize                     string
	CheckboxUncheckedColor           string
	Density                          DataTableDensity
	EmptyState                       templ.Component
	FilterDropdownBackgroundColor    string
	Filters                          []FilterSettings
	HeaderActions                    templ.Component
	HeaderClass                      string
	Id                               string
	InitialFilterValues              map[string]any
	InitialSearchQuery               string
	InitialSortDirection             DataTableSortDirection
	InitialSortKey                   string
	IsHeaderSticky                   bool
	IsPaginationHiddenWhenSinglePage bool
	IsStriped                        bool
	ItemsPerPage                     ItemsPerPage
	ItemsPerPageSizeChoices          []ItemsPerPage
	ItemsTotal                       uint
	PageNumber                       uint
	PaginationAriaLabel              string
	PagesTotal                       uint
	QueryUrlTemplate                 string
	RefreshDebounceMs                uint
	RefreshOnEvents                  []string
	RowClassResolver                 func(row Row) string
	RowIdResolver                    func(row Row) string
	RowLabelResolver                 func(row Row) string
	SearchBox                        templ.Component
	SearchBoxAlignment               HorizontalAlignment
	ShouldUseOneBasedPageDisplay     bool
	TextCase                         string
}

type dataTableInitialState struct {
	FilterValues  map[string]any `json:"filterValues"`
	ItemsPerPage  uint           `json:"itemsPerPage"`
	PageNumber    uint           `json:"pageNumber"`
	SearchQuery   string         `json:"searchQuery"`
	SortDirection string         `json:"sortDirection"`
	SortKey       string         `json:"sortKey"`
}

type dataTableClientSettings struct {
	FilterQueryParamNames map[string]string     `json:"filterQueryParamNames"`
	InitialState          dataTableInitialState `json:"initialState"`
	RefreshDebounceMs     uint                  `json:"refreshDebounceMs"`
	RefreshOnEvents       []string              `json:"refreshOnEvents"`
	QueryUrlTemplate      string                `json:"queryUrlTemplate"`
}

func (settings DataTableSettings[Row]) clientSettingsResolver(
	itemsPerPage, pageNumber uint,
) dataTableClientSettings {
	return dataTableClientSettings{
		FilterQueryParamNames: filterQueryParamNamesResolver(settings.Filters),
		InitialState: dataTableInitialState{
			FilterValues: initialFilterValuesResolver(
				settings.Filters, settings.InitialFilterValues,
			),
			ItemsPerPage:  itemsPerPage,
			PageNumber:    pageNumber,
			SearchQuery:   settings.InitialSearchQuery,
			SortDirection: string(settings.InitialSortDirection),
			SortKey:       settings.InitialSortKey,
		},
		RefreshDebounceMs: refreshDebounceResolver(settings.RefreshDebounceMs),
		RefreshOnEvents:   settings.RefreshOnEvents,
		QueryUrlTemplate:  settings.QueryUrlTemplate,
	}
}

func (settings DataTableSettings[Row]) tableIdHashResolver() uint64 {
	idParts := []string{settings.QueryUrlTemplate}
	for _, column := range settings.Columns {
		idParts = append(idParts, column.Label, column.SortKey)
	}
	for _, filter := range settings.Filters {
		idParts = append(idParts, filter.Key, filter.QueryParamName)
	}
	return uiToolset.HashComponentIdParts(idParts...)
}

func (settings DataTableSettings[Row]) idResolver() string {
	if settings.Id != "" {
		return settings.Id
	}
	return fmt.Sprintf("dataTable-%x", settings.tableIdHashResolver())
}

func (settings DataTableSettings[Row]) paginationAriaLabelResolver() string {
	if settings.PaginationAriaLabel != "" {
		return settings.PaginationAriaLabel
	}
	return dataTableDefaultPaginationAriaLabel
}

func (settings DataTableSettings[Row]) textCaseClassResolver() string {
	return uiToolset.TextCaseClassResolver(settings.TextCase)
}

func (settings DataTableSettings[Row]) cellPaddingClassesResolver() string {
	if settings.Density == DataTableDensityDense {
		return "px-2 py-1.5"
	}
	return "px-3 py-3"
}

func (settings DataTableSettings[Row]) headerPaddingClassesResolver() string {
	if settings.Density == DataTableDensityDense {
		return "px-2 py-1.5"
	}
	return "px-3 py-2"
}

func (settings DataTableSettings[Row]) checkboxShapeResolver() string {
	if settings.CheckboxShape != "" {
		return settings.CheckboxShape
	}
	return uiForm.CheckboxInputShapeSquare
}

func (settings DataTableSettings[Row]) checkboxSizeResolver() string {
	if settings.CheckboxSize != "" {
		return settings.CheckboxSize
	}
	return uiForm.CheckboxInputSizeMd
}

func (settings DataTableSettings[Row]) rowStripeClassesResolver() string {
	if settings.IsStriped {
		return "odd:bg-neutral-50/5 odd:hover:bg-neutral-50/10"
	}
	return ""
}

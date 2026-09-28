package uiStructural

import (
	"fmt"
	"hash/fnv"
	"maps"

	"github.com/a-h/templ"
	uiForm "github.com/goinfinite/ui/src/form"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

type DataTableDensity string

const (
	DataTableDensityComfortable DataTableDensity = "comfortable"
	DataTableDensityDense       DataTableDensity = "dense"
)

type DataTableAlignment string

const (
	DataTableAlignmentLeft   DataTableAlignment = "left"
	DataTableAlignmentCenter DataTableAlignment = "center"
	DataTableAlignmentRight  DataTableAlignment = "right"
)

func (alignment DataTableAlignment) alignmentClassResolver() string {
	switch alignment {
	case DataTableAlignmentCenter:
		return "text-center"
	case DataTableAlignmentRight:
		return "text-right"
	}
	return "text-left"
}

func (alignment DataTableAlignment) justifyClassResolver() string {
	switch alignment {
	case DataTableAlignmentCenter:
		return "justify-center"
	case DataTableAlignmentRight:
		return "justify-end"
	}
	return "justify-start"
}

type DataTableSortDirection string

const (
	DataTableSortDirectionAsc  DataTableSortDirection = "asc"
	DataTableSortDirectionDesc DataTableSortDirection = "desc"
)

type DataTablePageSize uint

const (
	DataTableUrlPlaceholderPageNumber    string = "{pageNumber}"
	DataTableUrlPlaceholderItemsPerPage  string = "{itemsPerPage}"
	DataTableUrlPlaceholderSortKey       string = "{sortKey}"
	DataTableUrlPlaceholderSortDirection string = "{sortDirection}"
	DataTableUrlPlaceholderSearch        string = "{search}"

	dataTableDefaultPaginationAriaLabel string = "Table pagination"
	dataTableDefaultRefreshDebounceMs   uint   = 300
)

type DataTableColumnSettings[Row any] struct {
	Label        string
	CellRenderer func(row Row) templ.Component

	// OptionalFields
	Alignment     DataTableAlignment
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
	ItemsPerPage                     DataTablePageSize
	ItemsPerPageSizeChoices          []DataTablePageSize
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
	SearchBoxAlignment               DataTableAlignment
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

func (settings DataTableSettings[Row]) initialFilterValuesResolver() map[string]any {
	initialValues := map[string]any{}
	for _, filter := range settings.Filters {
		switch filter.Kind {
		case FilterKindNumberRange, FilterKindDateRange:
			initialValues[filter.Key] = map[string]string{"min": "", "max": ""}
		default:
			initialValues[filter.Key] = ""
		}
	}
	maps.Copy(initialValues, settings.InitialFilterValues)
	return initialValues
}

func (settings DataTableSettings[Row]) clientSettingsResolver(
	itemsPerPage, pageNumber uint,
) dataTableClientSettings {
	filterQueryParamNames := map[string]string{}
	for _, filter := range settings.Filters {
		queryParamName := filter.QueryParamName
		if queryParamName == "" {
			queryParamName = filter.Key
		}
		filterQueryParamNames[filter.Key] = queryParamName
	}
	refreshDebounceMs := settings.RefreshDebounceMs
	if refreshDebounceMs == 0 {
		refreshDebounceMs = dataTableDefaultRefreshDebounceMs
	}
	return dataTableClientSettings{
		FilterQueryParamNames: filterQueryParamNames,
		InitialState: dataTableInitialState{
			FilterValues:  settings.initialFilterValuesResolver(),
			ItemsPerPage:  itemsPerPage,
			PageNumber:    pageNumber,
			SearchQuery:   settings.InitialSearchQuery,
			SortDirection: string(settings.InitialSortDirection),
			SortKey:       settings.InitialSortKey,
		},
		RefreshDebounceMs: refreshDebounceMs,
		RefreshOnEvents:   settings.RefreshOnEvents,
		QueryUrlTemplate:  settings.QueryUrlTemplate,
	}
}

func (settings DataTableSettings[Row]) tableIdentityHashResolver() uint64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(settings.QueryUrlTemplate))
	for _, column := range settings.Columns {
		hasher.Write([]byte(column.Label))
		hasher.Write([]byte(column.SortKey))
	}
	for _, filter := range settings.Filters {
		hasher.Write([]byte(filter.Key))
		hasher.Write([]byte(filter.QueryParamName))
	}
	return hasher.Sum64()
}

func (settings DataTableSettings[Row]) idResolver() string {
	if settings.Id != "" {
		return settings.Id
	}
	return fmt.Sprintf("dataTable-%x", settings.tableIdentityHashResolver())
}

func (settings DataTableSettings[Row]) itemsPerPageSizeChoicesResolver() []uint {
	if len(settings.ItemsPerPageSizeChoices) == 0 {
		return paginationDefaultItemsPerPageSizeChoices
	}
	itemsPerPageSizeChoices := make([]uint, len(settings.ItemsPerPageSizeChoices))
	for index, pageSize := range settings.ItemsPerPageSizeChoices {
		itemsPerPageSizeChoices[index] = uint(pageSize)
	}
	return itemsPerPageSizeChoices
}

func (settings DataTableSettings[Row]) itemsPerPageResolver(
	itemsPerPageSizeChoices []uint,
) uint {
	if settings.ItemsPerPage > 0 {
		return uint(settings.ItemsPerPage)
	}
	return itemsPerPageSizeChoices[0]
}

func (settings DataTableSettings[Row]) pageNumberResolver() uint {
	if settings.PageNumber > 0 {
		return settings.PageNumber
	}
	return 1
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

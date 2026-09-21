package main

import (
	"fmt"

	"github.com/a-h/templ"
	uiForm "github.com/goinfinite/ui/src/form"
	uiStructural "github.com/goinfinite/ui/src/structural"
)

const dataTableDemoItemsPerPage uiStructural.DataTablePageSize = 5

var dataTableDemoStaticPageSizeOptions = []uiStructural.DataTablePageSize{5, 10, 25, 50}

var dataTableDemoRecords = buildDemoServerRecords()

func demoTablePageCountResolver(
	itemsTotal uint, itemsPerPage uiStructural.DataTablePageSize,
) uint {
	itemsPerPageCount := uint(itemsPerPage)
	return (itemsTotal + itemsPerPageCount - 1) / itemsPerPageCount
}

var dataTableDemoPagesTotal = demoTablePageCountResolver(
	uint(len(dataTableDemoRecords)), dataTableDemoItemsPerPage,
)

type DemoServerRecord struct {
	Id          string
	Name        string
	Status      string
	CpuCores    string
	CreatedAt   string
	Description string
}

func buildDemoServerRecords() []DemoServerRecord {
	names := []string{
		"alpha", "bravo", "charlie", "delta", "echo",
		"foxtrot", "golf", "hotel", "india", "juliett",
		"kilo", "lima", "mike", "november", "oscar",
		"papa", "quebec", "romeo", "sierra", "tango",
		"uniform", "victor", "whiskey", "xray", "yankee",
	}
	descriptions := []string{
		"Handles the public API traffic for the primary region and drains connections during rolling deploys.",
		"Runs the nightly batch jobs, the reporting pipeline, and the weekly archive export.",
		"Hosts the internal dashboard, the metrics collector, and the alerting rules engine.",
		"Serves the static assets, the image resizing service, and the signed download links.",
		"Runs the message queue workers, the retry scheduler, and the dead letter processor.",
	}
	statuses := []string{"running", "running", "stopped"}
	cpuCores := []string{"2", "4", "1", "8", "2", "16"}
	records := make([]DemoServerRecord, 0, len(names))
	for index, name := range names {
		records = append(records, DemoServerRecord{
			Id:          fmt.Sprintf("srv-%02d", index+1),
			Name:        name,
			Description: descriptions[index%len(descriptions)],
			Status:      statuses[index%len(statuses)],
			CpuCores:    cpuCores[index%len(cpuCores)],
			CreatedAt:   fmt.Sprintf("2026-09-%02d", index+1),
		})
	}
	return records
}

func buildDemoRecordFilters() []uiStructural.FilterSettings {
	return []uiStructural.FilterSettings{
		{Key: "name", Label: "Name", Type: uiStructural.FilterTypeTextContains},
		{
			Key:   "status",
			Label: "Status",
			Type:  uiStructural.FilterTypeEnumSelect,
			Options: []uiForm.SelectLabelValueOption{
				{Label: "running", Value: "running"},
				{Label: "stopped", Value: "stopped"},
			},
		},
		{Key: "cpu", Label: "CPU", Type: uiStructural.FilterTypeNumberRange},
	}
}

func buildDemoServerTableColumns() []uiStructural.DataTableColumnSettings[DemoServerRecord] {
	return []uiStructural.DataTableColumnSettings[DemoServerRecord]{
		{
			Label:        "Name",
			SortKey:      "name",
			WidthPercent: 30,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.Name)
			},
		},
		{
			Label:        "Status",
			WidthPercent: 20,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerStatusCell(record.Status)
			},
		},
		{
			Label:        "CPU cores",
			SortKey:      "cpuCores",
			Alignment:    uiStructural.DataTableAlignmentRight,
			WidthPercent: 15,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.CpuCores)
			},
		},
		{
			Label:        "Created at",
			SortKey:      "createdAt",
			WidthPercent: 20,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.CreatedAt)
			},
		},
	}
}

func buildDemoPlainTableColumns() []uiStructural.DataTableColumnSettings[DemoServerRecord] {
	return []uiStructural.DataTableColumnSettings[DemoServerRecord]{
		{
			Label: "Name",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.Name)
			},
		},
		{
			Label: "Status",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerStatusCell(record.Status)
			},
		},
		{
			Label:     "CPU cores",
			Alignment: uiStructural.DataTableAlignmentRight,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.CpuCores)
			},
		},
	}
}

func buildDemoColumnShowcaseColumns() []uiStructural.DataTableColumnSettings[DemoServerRecord] {
	return []uiStructural.DataTableColumnSettings[DemoServerRecord]{
		{
			Label:        "Name",
			SortKey:      "name",
			WidthPercent: 20,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.Name)
			},
		},
		{
			Label:         "Description",
			MaxWidthClass: "max-w-72 truncate",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.Description)
			},
		},
		{
			Label:        "Status",
			Alignment:    uiStructural.DataTableAlignmentCenter,
			WidthPercent: 15,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerStatusCell(record.Status)
			},
		},
		{
			Label:        "CPU cores",
			SortKey:      "cpuCores",
			Alignment:    uiStructural.DataTableAlignmentRight,
			WidthPercent: 12,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.CpuCores)
			},
		},
		{
			Label:         "Created at",
			MinWidthClass: "min-w-28",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerTextCell(record.CreatedAt)
			},
		},
	}
}

func buildDemoStylingColumns() []uiStructural.DataTableColumnSettings[DemoServerRecord] {
	return []uiStructural.DataTableColumnSettings[DemoServerRecord]{
		{
			Label:     "Name",
			CellClass: "font-bold",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerPlainTextCell(record.Name)
			},
		},
		{
			Label: "Status",
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerStatusCell(record.Status)
			},
		},
		{
			Label:     "CPU cores",
			Alignment: uiStructural.DataTableAlignmentRight,
			Render: func(record DemoServerRecord) templ.Component {
				return DemoServerPlainTextCell(record.CpuCores)
			},
		},
	}
}

func buildDemoDataTableSettings(
	records []DemoServerRecord, pageNumber uint,
	itemsPerPage uiStructural.DataTablePageSize,
) uiStructural.DataTableSettings[DemoServerRecord] {
	itemsPerPageCount := uint(itemsPerPage)
	firstIndex := (pageNumber - 1) * itemsPerPageCount
	lastIndex := min(firstIndex+itemsPerPageCount, uint(len(records)))
	pagesTotal := demoTablePageCountResolver(uint(len(records)), itemsPerPage)
	return uiStructural.DataTableSettings[DemoServerRecord]{
		Columns: buildDemoServerTableColumns(),
		Rows:    records[firstIndex:lastIndex],
		UrlTemplate: "assets/dataTableDemoRefresh.html" +
			"?page=" + uiStructural.DataTableUrlPlaceholderPageNumber +
			"&itemsPerPage=" + uiStructural.DataTableUrlPlaceholderItemsPerPage +
			"&sort=" + uiStructural.DataTableUrlPlaceholderSortKey +
			"&direction=" + uiStructural.DataTableUrlPlaceholderSortDirection +
			"&search=" + uiStructural.DataTableUrlPlaceholderSearch,
		InitialFilterValues:  map[string]any{"status": "running"},
		Filters:              buildDemoRecordFilters(),
		Id:                   "data-table-demo-table",
		ItemsPerPage:         itemsPerPage,
		ItemsPerPageOptions:  []uiStructural.DataTablePageSize{dataTableDemoItemsPerPage, uiStructural.DataTablePageSize(len(records))},
		ItemsTotal:           uint(len(records)),
		PageNumber:           pageNumber,
		PagesTotal:           pagesTotal,
		RowIdResolver:        func(record DemoServerRecord) string { return record.Id },
		RowLabelResolver:     func(record DemoServerRecord) string { return record.Name },
		InitialSortKey:       "name",
		InitialSortDirection: uiStructural.DataTableSortDirectionAsc,
	}
}

func dataTableDemoStaticPageSizeResolver(itemsTotal uint) uiStructural.DataTablePageSize {
	for _, pageSizeOption := range dataTableDemoStaticPageSizeOptions {
		if uint(pageSizeOption) >= itemsTotal {
			return pageSizeOption
		}
	}
	return uiStructural.DataTablePageSize(itemsTotal)
}

func buildDemoStaticTableSettings(
	id string, paginationAriaLabel string,
	columns []uiStructural.DataTableColumnSettings[DemoServerRecord],
	rows []DemoServerRecord,
) uiStructural.DataTableSettings[DemoServerRecord] {
	itemsTotal := uint(len(rows))
	itemsPerPage := dataTableDemoStaticPageSizeResolver(itemsTotal)
	pagesTotal := demoTablePageCountResolver(itemsTotal, itemsPerPage)
	return uiStructural.DataTableSettings[DemoServerRecord]{
		Columns:             columns,
		Rows:                rows,
		Id:                  id,
		ItemsPerPage:        itemsPerPage,
		ItemsPerPageOptions: []uiStructural.DataTablePageSize{itemsPerPage},
		ItemsTotal:          itemsTotal,
		PageNumber:          1,
		PaginationAriaLabel: paginationAriaLabel,
		PagesTotal:          pagesTotal,
	}
}

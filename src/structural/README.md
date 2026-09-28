# Structural

Structural layer of Infinite UI. It composes form and display components into page-level structures. Part of [Infinite UI](../../README.md).

## Card

Surface container with an optional header block and content slots.

```go
@uiStructural.Card(uiStructural.CardSettings{
    MiddleContent: CardBody(),

    // OptionalFields
    HeaderTitle:      "Card Title",
    HeaderSubHeading: "Card sub-heading",
    HeaderIcon:       "ph-cube",
    FooterContent:    CardFooter(),
})
```

- `MiddleContent` is the card body. `FooterContent` sits below it. `HeaderContent` replaces the whole header block.
- The optional header block carries `HeaderTitle`, `HeaderSubHeading`, `HeaderTitleOneWayStatePath`, `HeaderSubHeadingOneWayStatePath`, `HeaderTitleColor`, `HeaderSubHeadingColor`, `HeaderSize`, and the `HeaderIcon*` icon controls. `ActionsContent` puts buttons in the header row, right-aligned.
- `BorderRadius` accepts `CardBorderRadiusNone` through `CardBorderRadiusXl`; the default is `CardBorderRadiusLg`. Use `CardBorderRadiusNone` for square edges.
- `PaddingSize` accepts `CardPaddingSizeNone` through `CardPaddingSizeXl`; the default is `CardPaddingSizeMd`.
- `GapSize` accepts `CardGapSizeNone` through `CardGapSizeXl`; the default is `CardGapSizeMd`, which renders `gap-3` between the header, body, and footer.
- `ShadowSize`, `RingColor`, and `RingThickness` follow the same token scales as Modal and Alert. `BackgroundColor` and `TextColor` take color tokens. `TextCase` takes a `uiToolset.TextCase*` value and transforms the header title and sub-heading.

## DataTable

`@uiStructural.DataTable` renders rows from your data and refreshes them from your server. Every sort, page, filter, or search change requests the URL template you provide. The table uses `htmx.ajax` when HTMX is present and falls back to `fetch` otherwise.

The component requires a server that answers each request. The static demo serves fixed pages, so sorting, filtering, and search update the request only. Serve it over HTTP; browsers block refresh requests from `file://` pages.

```go
@uiStructural.DataTable(uiStructural.DataTableSettings[Record]{
    Columns: columns,
    Rows:    records,
    QueryUrlTemplate: "/records?page=" + uiStructural.DataTableUrlPlaceholderPageNumber +
        "&sort=" + uiStructural.DataTableUrlPlaceholderSortKey +
        "&direction=" + uiStructural.DataTableUrlPlaceholderSortDirection +
        "&search=" + uiStructural.DataTableUrlPlaceholderSearch,

    // OptionalFields
    Filters:              filters,
    RowIdResolver:        func(record Record) string { return record.Id },
    RefreshOnEvents:      []string{"refresh:records-table"},
    InitialSortKey:       "name",
    InitialSortDirection: uiStructural.DataTableSortDirectionAsc,
})
```

Each column takes a `Label`, a `CellRenderer` function, and optional `SortKey`, `Alignment`, `WidthPercent`, width classes, and `CellClass`. `Alignment` takes a `DataTableAlignment` value: `DataTableAlignmentLeft`, `DataTableAlignmentCenter`, or `DataTableAlignmentRight`. `TextCase` takes a `uiToolset.TextCase*` value and transforms the header labels. The default, `TextCaseNone`, leaves them as typed. `Density` takes a `DataTableDensity` value: `DataTableDensityComfortable` (the default) or `DataTableDensityDense`. `InitialSortDirection` takes a `DataTableSortDirection` value: `DataTableSortDirectionAsc` or `DataTableSortDirectionDesc`. `ItemsPerPage` and each entry in `ItemsPerPageSizeChoices` are `DataTablePageSize` values. Set `PaginationAriaLabel` when a page holds more than one table, so each pagination landmark keeps a unique name.

The `Initial*` fields seed the client state at render time: `InitialFilterValues`, `InitialSearchQuery`, `InitialSortKey`, and `InitialSortDirection`. The server renders the matching rows. `PageNumber` and `ItemsPerPage` also seed the client, but the component reads them to render the pagination readout.

`HeaderClass` adds classes to the header row, `CellClass` adds classes to one column's cells, `RowClassResolver` returns classes for each row from its data, and `IsStriped` adds a zebra stripe. These classes append to elements that already carry base utilities, so when two utilities set the same property the generated stylesheet order decides the winner, not the field order. A cell component that sets its own color wins over the row color, so use `RowClassResolver` for cells that leave the color to the row. When `IsHeaderSticky` is set, the sticky header paints its own background, so a `HeaderClass` background does not show. The table renders a default search box when the query URL template carries the search placeholder; pass `SearchBox` to replace it. `SearchBoxAlignment` takes a `DataTableAlignment` value and places the search box left (the default), center, or right within the toolbar. `CheckboxShape` accepts `uiForm.CheckboxInputShapeSquare` (the default), `uiForm.CheckboxInputShapeRounded`, or `uiForm.CheckboxInputShapeCircular`; `CheckboxSize` accepts the `uiForm.CheckboxInputSize*` values and defaults to the medium size; `CheckboxCheckedColor` and `CheckboxUncheckedColor` take a color token and default to `secondary-500` and `neutral-50/20`.

The query URL template uses fixed placeholders. Build it from the `DataTableUrlPlaceholder*` constants and name the query keys:

```go
QueryUrlTemplate: "/records?page=" + uiStructural.DataTableUrlPlaceholderPageNumber +
    "&size=" + uiStructural.DataTableUrlPlaceholderItemsPerPage +
    "&sort=" + uiStructural.DataTableUrlPlaceholderSortKey +
    "&direction=" + uiStructural.DataTableUrlPlaceholderSortDirection +
    "&q=" + uiStructural.DataTableUrlPlaceholderSearch
```

Filter values append to the URL as `key=value` pairs. Number and date ranges append as `keyMin` and `keyMax`. Empty values are omitted. Set `QueryParamName` on a filter to send a different query key.

The server response must contain one element with the `data-ui-data-table` attribute. The component swaps only that element, so the filter bar, search box, and selection stay in place.

Client state lives in the component root: `pageNumber`, `itemsPerPage`, `sortKey`, `sortDirection`, `searchQuery`, `filterValues`, and `selectedRowIds`. The search box and the bulk action slot bind to those paths.

`RefreshOnEvents` lists window event names. Dispatching one of them refreshes the table. This matches the form-to-display refresh pattern.

A failed refresh shows an inline error row with a retry button. A refresh in flight dims the table and disables the controls.

Pass filter keys, state paths, and the query URL template from code, never from request data. The component embeds them into client-side expressions. `HeaderClass`, `CellClass`, the `RowClassResolver` result, `CheckboxCheckedColor`, and `CheckboxUncheckedColor` become HTML class attributes, so keep untrusted data out of them too.

## FilterBar

Standalone filter bar. It renders one editor per declared filter and shows active filters as removable chips.

```go
@uiStructural.FilterBar(uiStructural.FilterBarSettings{
    Filters: []uiStructural.FilterSettings{
        {Key: "name", Label: "Name", Kind: uiStructural.FilterKindTextContains},
        {Key: "status", Label: "Status", Kind: uiStructural.FilterKindEnumSelect, Options: statusOptions},
        {Key: "cpu", Label: "CPU", Kind: uiStructural.FilterKindNumberRange},
    },
    ValuesTwoWayStatePath: "filterValues",

    // OptionalFields
    OnChangeFunc: "refreshTable()",
})
```

- `Kind` accepts `FilterKindTextContains`, `FilterKindEnumSelect`, `FilterKindMultiEnumSelect`, `FilterKindNumberRange`, or `FilterKindDateRange`.
- `ValuesTwoWayStatePath` names a top-level property on the surrounding Alpine scope, for example `filterValues`.
- Enum filters read `Options`. Multi-enum filters read `Options` and write an array of selected values under the filter key. Range filters write `{min, max}` objects under the filter key.
- A chip appears when its filter holds a value. The chip remove button clears that filter.
- The clear-all button appears when any filter is active.
- Set `EnumSelectRadioGroupNamePrefix` when a page holds more than one filter bar with the same enum keys. The prefix keeps each enum dropdown's radio group name unique. The DataTable prefixes it with the table id.

## PageHeading

Page and section headings over `display.HeaderBlock`. One `Level` setting picks the variant.

```go
@uiStructural.PageHeading(uiStructural.PageHeadingSettings{
    HeaderTitle: "Records",
    Level:       uiStructural.PageHeadingLevelPage,

    // OptionalFields
    Description:    "Manage the server records",
    HeaderIcon:     "ph-table",
    ActionsContent: PageHeadingActions(),
})
```

- `Level` accepts `PageHeadingLevelPage` or `PageHeadingLevelSection`; the default is section.
- `PageHeadingLevelPage` renders an `h1` with page spacing. `HeaderSize` defaults to `HeaderSizeXl`. The icon defaults to a bare `secondary-500` glyph.
- `PageHeadingLevelSection` renders an `h2`. `HeaderSize` defaults to `HeaderSizeLg`. With an icon set and no overrides, the icon renders in a padded rounded chip.
- `Description` fills the line under the title. `HeaderSubHeading` is the legacy name for the same line.
- `ActionsContent` puts buttons in the heading row, right-aligned.
- `HeaderTitle` and the description support live text through `HeaderTitleOneWayStatePath` and `HeaderSubHeadingOneWayStatePath`.
- `HeaderTitleColor`, `HeaderSubHeadingColor`, and `TextColor` take color tokens.
- `TextCase` takes a `uiToolset.TextCase*` value and transforms the title and description.
- The `HeaderIcon*` fields control the icon: `HeaderIconPosition` (`HeaderIconPositionLeft` or `HeaderIconPositionTop`), `HeaderIconColor`, `HeaderIconBackgroundColor`, `HeaderIconBorderRadius`, and `HeaderIconPaddingSize`.

## Pagination

Page controls with a readout, a page-number strip, and an items-per-page selector.

```go
@uiStructural.Pagination(uiStructural.PaginationSettings{
    PageNumberTwoWayStatePath:   "pageNumber",
    ItemsPerPageTwoWayStatePath: "itemsPerPage",
    ItemsTotal:                  240,

    // OptionalFields
    OnChangeFunc: "requestRefresh()",
})
```

- The readout follows the state paths. It shows the current range and the total.
- The component derives the page count from `ItemsTotal` and the bound `itemsPerPage`, so the strip and the controls react when the page size changes. `PagesTotal` is an optional fallback used only when `ItemsTotal` is zero.
- The strip shows the first page, the last page, the pages around the current one, and ellipses for gaps. The current page carries `aria-current="page"`.
- `ItemsPerPageSizeChoices` overrides the default page sizes.
- Set `ItemsPerPageInputName` when a page holds more than one pagination bound to the same state path, so the two items-per-page radio groups stay independent.
- `IsDisabledOneWayStatePath` disables every control while the path is truthy.
- `IsHiddenWhenSinglePage` hides the page-number controls while every record fits on one page. The readout and the items-per-page select stay visible. DataTable forwards it as `IsPaginationHiddenWhenSinglePage`.
- `AriaLabel` names the navigation landmark. Set a distinct label when a page holds more than one pagination.

## Sidebar

Vertical navigation panel. It renders inline or fixed, and it can collapse or slide off canvas.

```go
@uiStructural.Sidebar(uiStructural.SidebarSettings{
    MiddleContent: SidebarNavigation(),

    // OptionalFields
    HeaderContent:                 SidebarHeader(),
    FooterContent:                 SidebarFooter(),
    IsCollapsedTwoWayStatePath:    "isSidebarCollapsed",
    AttachmentMode:                uiStructural.SidebarAttachmentModeInline,
    Side:                          uiStructural.SidebarSideLeft,
})
```

- `BackgroundColor` takes a full Tailwind class, for example `"bg-neutral-800/50"`, unlike the color tokens other components take.
- `AttachmentMode` accepts `SidebarAttachmentModeInline` or `SidebarAttachmentModeFixed`.
- `IsOffCanvas` and `IsOffCanvasTwoWayStatePath` slide the panel over the content.
- `Width` sets the expanded width.
- Pass state paths from code, never from request data. The component embeds them into client-side expressions.

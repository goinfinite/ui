# Structural

The structural package combines form and display components into page structures. Part of [Infinite UI](../../README.md).

## Card

A surface with an optional header and content slots.

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

- `MiddleContent` is the body. `FooterContent` goes below it. `HeaderContent` replaces the header.
- The header takes `HeaderTitle`, `HeaderSubHeading`, their `*OneWayStatePath` and `*Color` fields, `HeaderSize`, and the `HeaderIcon*` fields. `ActionsContent` puts buttons on the right of the header.
- `BorderRadius` takes `uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull` (default `uiToolset.BorderRadiusLg`).
- `PaddingSize` takes `uiToolset.PaddingSizeNone` through `uiToolset.PaddingSizeXl` (default `uiToolset.PaddingSizeMd`).
- `GapSize` takes `uiToolset.GapSizeNone` through `uiToolset.GapSizeXl`. The default, `uiToolset.GapSizeMd`, puts `gap-3` between the header, body, and footer.
- `ShadowSize`, `RingColor`, and `RingThickness` use the same scales as Modal and Alert. `BackgroundColor` and `TextColor` take color tokens.
- `TextCase` takes a `uiToolset.TextCase*` value for the header title and sub-heading.

## Carousel

A carousel that moves a window of items inside one server page. The server returns a chunk of items, and the window moves inside that chunk. The footer pages the chunk with the `Pagination` of DataTable. The filter bar and the search box load a new chunk from your server.

```go
@uiStructural.Carousel(uiStructural.CarouselSettings[Record]{
    ItemRenderer: func(record Record) templ.Component { return CarouselItem(record) },
    Items:        records,
    QueryUrlTemplate: "/records?page=" + uiStructural.CarouselUrlPlaceholderPageNumber +
        "&itemsPerPage=" + uiStructural.CarouselUrlPlaceholderItemsPerPage,

    // OptionalFields
    ItemsPerView:   uiStructural.CarouselItemsPerViewSettings{Base: 1, Sm: 2, Lg: 3},
    IsSwipeEnabled: true,
})
```

### Items and paging

- `ItemRenderer` renders one item. `Items` holds the chunk from the server.
- The carousel sets the item width. Each visible item gets an equal share. Do not set a width on the item component.
- `ItemsPerView` sets the visible count for each Tailwind breakpoint, mobile first. `Base` applies to all widths (one when unset). `Sm` (≥640px), `Md` (≥768px), `Lg` (≥1024px), `Xl` (≥1280px), and `TwoXl` (≥1536px) apply from their width up. An unset breakpoint uses the one below it. The breakpoints use the browser window width.
- `ItemsPerPage` and `ItemsPerPageSizeChoices` page the server chunk. `ItemsTotal` and `PagesTotal` feed the footer, as in DataTable.
- `IsPaginationHiddenWhenSinglePage` hides the page numbers while all items fit on one page. The readout and the items-per-page select stay.
- `EmptyState` shows when the chunk has no items.

### Query, filters, and search

- `QueryUrlTemplate` takes the placeholders `CarouselUrlPlaceholderPageNumber`, `CarouselUrlPlaceholderItemsPerPage`, and `CarouselUrlPlaceholderSearch`.
- `Filters` shows the shared `FilterBar` below the track. Filter values go into the URL as `key=value`. Ranges go in as `keyMin` and `keyMax`. `FilterChipStyle` restyles the active-filter chip (see FilterBar).
- The filter bar stays outside the refreshed region, so a refresh does not take focus from a field.
- A search box shows above the track when the template has the search placeholder. `SearchBox` replaces it. `IsSearchBoxCompact` makes it narrow, and `SearchBoxAlignment` places it left, center (the default), or right. `SearchBoxPosition` takes `CarouselSearchBoxPositionTop` (the default) or `CarouselSearchBoxPositionBottom`.

### Motion

- `IsAutoplay` moves the window on a timer. `AutoplayIntervalMs` sets the interval (default 4000). `IsAutoplayPausedOnHover` stops the timer under the pointer.
- `IsSwipeEnabled` moves the window on a horizontal drag of 40 pixels or more.

### Styling

- Surface: `BackgroundColor`, `TextColor`, `BorderRadius`, `PaddingSize`, `GapSize`, `ShadowSize`, `RingColor`, and `RingThickness`. They take the `uiToolset` token scales.
- Arrows: `ArrowsPosition` takes `CarouselArrowsPositionOutside` (the default) or `CarouselArrowsPositionInside`. `ArrowsShape` takes `uiToolset.ShapeCircular` (the default), `uiToolset.ShapeRounded`, or `uiToolset.ShapeSquare`. `ArrowsSize` takes `CarouselArrowsSizeSm`, `Md`, or `Lg` (a 32-pixel square by default). `ArrowsBackgroundColor` and `ArrowsIconColor` set the colors. A custom background still gets lighter on hover.
- Dots: `DotsPosition` takes `CarouselDotsPositionBottom` (the default) or `CarouselDotsPositionTop`. `DotsSize` takes `CarouselDotsSizeSm`, `Md`, or `Lg`. `DotsActiveColor` and `DotsInactiveColor` set the colors.
- Items: `ItemBackgroundColor`, `ItemBorderRadius`, `ItemPaddingSize`, `ItemRingColor`, `ItemRingThickness`, and `ItemShadowSize` style the item wrapper, so a plain renderer still looks like a card.

### Refresh

- `RefreshOnEvents` lists window event names. Each event refreshes the carousel.
- The refresh uses `htmx.ajax` when HTMX is on the page, and `fetch` when it is not.
- One response updates two regions: the body (`data-ui-carousel`, with arrows, track, and dots) and the pagination (`data-ui-carousel-pagination`). The response must contain both.
- The pagination id is the component id plus `-pagination`. `hx-swap-oob` targets that id, so HTMX swaps it out of band. The root has `hx-sync`, so a new refresh stops the one in progress.
- Answer the refresh with `@uiStructural.CarouselRefreshRegions(settings)`. It renders only the two regions, so the response stays small. A full `Carousel` response also works.
- The root holds the client state: `pageNumber`, `itemsPerPage`, `searchQuery`, `filterValues`, `windowStart`, `itemsPerView`, and `itemsCount`.
- A failed refresh shows an error with a retry button.

### Accessibility and trust

- The arrows have accessible names. The dots have `aria-label` and `aria-current`. Items outside the window are `inert`, so they leave the tab order.
- Give state paths, breakpoint values, filter keys, the URL template, and the color inputs from code, never from request data. The component puts them into client-side expressions and class attributes.

## CarouselItemTooltip

Wraps one carousel item and shows a tooltip on hover or on focus.

```go
templ CarouselItem(record Record) {
    @uiStructural.CarouselItemTooltip(uiStructural.CarouselItemTooltipSettings{
        Content: record.Description,

        // OptionalFields
        Position:      uiToolset.TooltipPositionBottom,
        RingColor:     "secondary-500/40",
        RingThickness: uiToolset.RingThicknessXs,
    }) {
        <div class="flex h-full flex-col gap-2">
            <span class="font-bold">{ record.Name }</span>
            <p class="text-neutral-400">{ record.Description }</p>
        </div>
    }
}
```

- `Content` is the tooltip text. `ContentHtml` takes a component. Without either, the item has no tooltip.
- `Position` takes a `uiToolset.TooltipPosition*` value (default top).
- `BackgroundColor` and `TextColor` take color tokens (default `bg-neutral-800/95` with `text-neutral-50`).
- `RingColor` takes a color token. `RingThickness` takes a `uiToolset.RingThickness*` value on the compact scale, from `ring-0.5` (xs) to `ring-2.5` (xl).
- `MinWidthClass`, `MaxWidthClass`, `MinHeightClass`, and `MaxHeightClass` limit the size. Long text wraps at `max-w-96`.
- The tooltip renders on the document body, so the track, a scroll container, or a modal cannot clip it. Alpine removes it with the item, so a refresh does not leave tooltips behind.

## DataTable

A table of your rows that refreshes from your server. Each sort, page, filter, or search change requests your URL template. It uses `htmx.ajax` when HTMX is on the page, and `fetch` when it is not.

The table needs a server that answers each request. The static demo serves fixed pages, so sort, filter, and search change only the request. Serve the page over HTTP, because browsers block refresh requests from `file://`.

```go
@uiStructural.DataTable(uiStructural.DataTableSettings[Record]{
    Columns: columns,
    Rows:    records,
    QueryUrlTemplate: "/records?page=" + uiStructural.DataTableUrlPlaceholderPageNumber +
        "&size=" + uiStructural.DataTableUrlPlaceholderItemsPerPage +
        "&sort=" + uiStructural.DataTableUrlPlaceholderSortKey +
        "&direction=" + uiStructural.DataTableUrlPlaceholderSortDirection +
        "&q=" + uiStructural.DataTableUrlPlaceholderSearch,

    // OptionalFields
    Filters:              filters,
    RowIdResolver:        func(record Record) string { return record.Id },
    RefreshOnEvents:      []string{"refresh:records-table"},
    InitialSortKey:       "name",
    InitialSortDirection: uiStructural.DataTableSortDirectionAsc,
})
```

### Columns and state

- A column takes a `Label`, a `CellRenderer`, and optional `SortKey`, `Alignment`, `WidthPercent`, width classes, and `CellClass`. `Alignment` takes `TextAlignmentLeft`, `TextAlignmentCenter`, or `TextAlignmentRight`.
- `TextCase` takes a `uiToolset.TextCase*` value for the header labels.
- `Density` takes `DataTableDensityComfortable` (the default) or `DataTableDensityDense`.
- `InitialSortDirection` takes `DataTableSortDirectionAsc` or `DataTableSortDirectionDesc`.
- The `Initial*` fields set the client state at render: `InitialFilterValues`, `InitialSearchQuery`, `InitialSortKey`, and `InitialSortDirection`. The server renders the matching rows.
- `PageNumber` and `ItemsPerPage` also set the client state and the pagination readout. `ItemsPerPage` and the `ItemsPerPageSizeChoices` entries are `uiStructural.ItemsPerPage` values.
- `PageNumber` starts at 0 for the first page. `ShouldUseOneBasedPageDisplay` changes only the labels.
- `PaginationAriaLabel` names the pagination. Use a different name for each table on a page.

### Row and header styling

- `HeaderClass` adds classes to the header row. `CellClass` adds classes to the cells of one column. `RowClassResolver` returns classes for each row. `IsStriped` adds stripes.
- These classes add to the base classes. When two classes set the same property, the stylesheet order decides, not the field order.
- A cell component that sets its own color wins over the row color. For row colors, leave the cell color empty.
- The header stays at the top of the scroll area by default, with a blurred `bg-neutral-950/20` surface. `StickyHeaderBackgroundColor` replaces the surface. `IsHeaderStatic` lets the header scroll.
- `ItemsPerPageDropdownBackgroundColor` sets the items-per-page menu color (default `bg-neutral-800/95`). A dropdown in the table shows above the header and outside the table edge.
- `MinWidthClass` and `MaxWidthClass` limit the table width. `MinHeightClass` and `MaxHeightClass` limit the scroll height (default `max-h-128`).
- `MaxVisibleRows` measures the header and the tallest row and shows that many rows. It replaces the height classes.
- A search box shows when the template has the search placeholder. `SearchBox` replaces it. `IsSearchBoxCompact` makes it narrow, and `SearchBoxAlignment` places it left (the default), center, or right.
- `CheckboxShape` takes `uiToolset.ShapeSquare` (the default), `uiToolset.ShapeRounded`, or `uiToolset.ShapeCircular`. `CheckboxSize` takes a `uiForm.CheckboxInputSize*` value (default medium). `CheckboxCheckedColor` (default `secondary-500`) and `CheckboxUncheckedColor` (default `neutral-50/20`) take color tokens.

### Query and filters

- Build `QueryUrlTemplate` from the `DataTableUrlPlaceholder*` constants and your own query keys.
- `DataTableUrlPlaceholderPageNumber` sends the page index that starts at 0.
- Filter values go into the URL as `key=value`. Ranges go in as `keyMin` and `keyMax`. Empty values do not go in. `QueryParamName` on a filter sets a different key.
- `FilterChipStyle` restyles the active-filter chip (see FilterBar).

### Refresh

- The response must contain one element with `data-ui-data-table`. The table replaces only that element, so the filter bar, the search box, and the selection stay.
- Answer the refresh with `@uiStructural.DataTableRefreshRegion(settings)`. It renders only that element, so the response stays small. A full `DataTable` response also works.
- The root holds the client state: `pageNumber`, `itemsPerPage`, `sortKey`, `sortDirection`, `searchQuery`, `filterValues`, and `selectedRowIds`. The search box and the bulk action slot bind to these paths.
- `RefreshOnEvents` lists window event names. Each event refreshes the table.
- During a refresh, the table dims and disables its controls. A failed refresh shows an error row with a retry button.

### Trust

- Give filter keys, state paths, the URL template, and `Id` from code, never from request data. The component puts them into client-side expressions and the root id.
- The class fields and the `RowClassResolver` result go into HTML class attributes. Keep untrusted data out of them.

## FilterBar

A filter bar. It shows one editor for each filter and shows active filters as removable chips.

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

- `Kind` takes `FilterKindTextContains`, `FilterKindEnumSelect`, `FilterKindMultiEnumSelect`, `FilterKindNumberRange`, or `FilterKindDateRange`.
- `ValuesTwoWayStatePath` names a top-level property of the Alpine scope, for example `filterValues`.
- Enum filters read `Options`. A multi-enum filter writes an array of values. A range filter writes a `{min, max}` object.
- A chip shows for each filter with a value. Its remove button clears that filter. The chip row shows only while a filter is active.
- The clear-all button shows while any filter is active.
- `ChipStyle` restyles the chip. `OuterBackgroundColor`, `OuterRingColor`, `OuterTextColor`, `InnerBackgroundColor`, and `InnerTextColor` take color tokens. `Radius` takes a `uiToolset.BorderRadius*` value (default `uiToolset.BorderRadiusMd`). `TextCase` takes a `uiToolset.TextCase*` value (default `uiToolset.TextCaseLower`).
- Set `EnumSelectRadioGroupNamePrefix` when a page has two filter bars with the same enum keys. The DataTable sets it from the table id.

## PageHeading

Page and section headings built on `display.HeaderBlock`. `Level` selects the variant.

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

- `Level` takes `PageHeadingLevelPage` or `PageHeadingLevelSection` (the default).
- A page heading is an `h1` with page spacing. `HeaderSize` defaults to `HeaderSizeXl`. The icon defaults to a `secondary-500` glyph.
- A section heading is an `h2`. `HeaderSize` defaults to `HeaderSizeLg`. The icon shows in a rounded chip.
- `Description` sets the line below the title. `HeaderSubHeading` is the old name for the same line.
- `ActionsContent` puts buttons on the right of the heading.
- `HeaderTitleOneWayStatePath` and `HeaderSubHeadingOneWayStatePath` bind live text.
- `HeaderTitleColor`, `HeaderSubHeadingColor`, and `TextColor` take color tokens.
- `TextCase` takes a `uiToolset.TextCase*` value for the title and the description.
- The `HeaderIcon*` fields set the icon: `HeaderIconPosition`, `HeaderIconColor`, `HeaderIconBackgroundColor`, `HeaderIconBorderRadius`, and `HeaderIconPaddingSize`.

## Pagination

Page controls with a readout, page numbers, and an items-per-page select.

```go
@uiStructural.Pagination(uiStructural.PaginationSettings{
    PageNumberTwoWayStatePath:   "pageNumber",
    ItemsPerPageTwoWayStatePath: "itemsPerPage",
    ItemsTotal:                  240,

    // OptionalFields
    OnChangeFunc: "requestRefresh()",
})
```

- The bound page number starts at 0. The readout shows the item range and the total, starting at 1.
- `ShouldUseOneBasedPageDisplay` labels the pages from 1. The state still starts at 0.
- The page count comes from `ItemsTotal` and the bound `itemsPerPage`. `PagesTotal` applies only when `ItemsTotal` is zero.
- The page numbers show the first page, the last page, the pages near the current one, and ellipses. The current page has `aria-current="page"`.
- `ItemsPerPageSizeChoices` replaces the default choices.
- `DropdownBackgroundColor` sets the menu color (default `bg-neutral-800/95`). DataTable passes it as `ItemsPerPageDropdownBackgroundColor`.
- `ItemsPerPageInputName` keeps two paginations on one state path independent.
- `IsDisabledOneWayStatePath` disables all controls while the path is truthy.
- `IsHiddenWhenSinglePage` hides the page numbers while all records fit on one page. DataTable passes it as `IsPaginationHiddenWhenSinglePage`.
- `AriaLabel` names the navigation. Use a different name for each pagination on a page.

## Sidebar

A vertical navigation panel. It shows inline or fixed, and it can collapse or slide off the screen.

```go
@uiStructural.Sidebar(uiStructural.SidebarSettings{
    MiddleContent: SidebarNavigation(),

    // OptionalFields
    HeaderContent:              SidebarHeader(),
    FooterContent:              SidebarFooter(),
    IsCollapsedTwoWayStatePath: "isSidebarCollapsed",
    AttachmentMode:             uiStructural.SidebarAttachmentModeInline,
    Side:                       uiStructural.SidebarSideLeft,
})
```

- `BackgroundColor` takes a full Tailwind class, for example `"bg-neutral-800/50"`. Other components take color tokens.
- `AttachmentMode` takes `SidebarAttachmentModeInline` or `SidebarAttachmentModeFixed`.
- `IsOffCanvas` and `IsOffCanvasTwoWayStatePath` slide the panel over the content.
- `Width` sets the open width.
- Give state paths from code, never from request data.

## Tabs

Horizontal or vertical tabs bound to one selected value. Each tab has a label, a value, and an optional icon and count badge. Alpine `x-show` switches the panels.

```go
@uiStructural.Tabs(uiStructural.TabSettings{
    Items: []uiStructural.TabItemSettings{
        {Label: "General", Value: "general", Icon: "ph-gear", Content: GeneralPanel()},
        {Label: "Security", Value: "security", Icon: "ph-shield-check", BadgeCount: 3, Content: SecurityPanel()},
    },
    SelectedValueTwoWayStatePath: "selectedTab",

    // OptionalFields
    Orientation: uiStructural.TabOrientationVertical,
})
```

### Selection

- `SelectedValueTwoWayStatePath` holds the active value. A click writes the tab `Value` to it.
- Start the path with one of the `Value`s. With no match, no tab is active, but the first tab stays reachable by keyboard.
- Each `Content` renders in a panel that stays in the DOM, so inputs keep their state.
- `IsUrlHashSynced` reads the URL hash on load and on each change, and writes it on click. A link to `#security` opens that tab.

### Layout

- `Orientation` takes `TabOrientationHorizontal` (the default) or `TabOrientationVertical`.
- For vertical tabs, `Side` takes `TabSideLeft` (the default) or `TabSideRight`, and `Alignment` takes `TabAlignmentTop` (the default), `TabAlignmentCenter`, or `TabAlignmentBottom`.
- `MaxWidthPercent`, `MinWidthPercent`, `MaxHeightPercent`, and `MinHeightPercent` limit the tab list as whole percentages from 1 to 100. The `*WidthClass` and `*HeightClass` fields take Tailwind classes for the same limits.
- `MaxVisibleTabs` measures the tallest tab and shows that many vertical tabs. It replaces the height limits.
- A list longer than its limits scrolls.

### Icons and badges

- `Icon` shows a Phosphor icon. `IconPosition` takes `TabIconPositionLeft` (the default) or `TabIconPositionTop`.
- `BadgeCount` shows a static count. `BadgeCountOneWayStatePath` shows a live count and hides at zero. The `Badge*` fields style the badge.

### Styling

- `BackgroundColor` and `TextColor` take color tokens for the list and the unselected labels. `BorderRadius` defaults to `uiToolset.BorderRadiusMd`.
- `PaddingSize` insets the list (default none). `GapSize` sets the space between tabs (default `uiToolset.GapSizeXs`).
- `RingColor` and `RingThickness` frame the list when both are set. `ShadowSize` takes a `uiToolset.ShadowSize*` value.
- The `Content*` fields style the panel area with the same tokens. `ContentPaddingSize` replaces the default space between the panel and the list.
- `TextCase` takes a `uiToolset.TextCase*` value for the labels.

### Accessibility and trust

- The tabs follow the ARIA tabs pattern: `role="tablist"`, `role="tab"`, `role="tabpanel"`, `aria-selected`, `aria-controls`, and `aria-labelledby`.
- Arrow keys move between tabs: Left and Right for horizontal tabs, Up and Down for vertical tabs. Home and End go to the first and last tab.
- `AriaLabel` names the tab list. Use a different name for each tab group on a page.
- Give state paths, tab values, `Icon`, and the color inputs from code, never from request data.

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
- `BorderRadius` accepts `uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull`; the default is `uiToolset.BorderRadiusLg`. Use `uiToolset.BorderRadiusNone` for square edges.
- `PaddingSize` accepts `uiToolset.PaddingSizeNone` through `uiToolset.PaddingSizeXl`; the default is `uiToolset.PaddingSizeMd`.
- `GapSize` accepts `uiToolset.GapSizeNone` through `uiToolset.GapSizeXl`; the default is `uiToolset.GapSizeMd`, which renders `gap-3` between the header, body, and footer.
- `ShadowSize`, `RingColor`, and `RingThickness` follow the same token scales as Modal and Alert. `BackgroundColor` and `TextColor` take color tokens. `TextCase` takes a `uiToolset.TextCase*` value and transforms the header title and sub-heading.

## Carousel

`@uiStructural.Carousel` slides a window of items inside a server page. The server returns a chunk of items; the window moves inside that chunk. The footer pages the chunk with the same `Pagination` component the DataTable uses, and the filter bar and search box refresh the chunk from your server.

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

- `ItemRenderer` renders one item. `Items` holds the chunk the server returned.
- The carousel owns the item width. Each item fills an equal share of the visible area, so the visible count fits at every window size. Do not set a width on the item component; the carousel overrides it.
- The visible count follows Tailwind's named breakpoints, mobile-first. `ItemsPerView` is a `CarouselItemsPerViewSettings` value: `Base` is the base (one when unset), and `Sm` (≥640px), `Md` (≥768px), `Lg` (≥1024px), `Xl` (≥1280px), and `TwoXl` (≥1536px) each raise the count from their width up. A breakpoint you leave unset inherits the lower one, exactly like a missing `md:` class. The count is measured against the browser window, the same basis the rest of the library uses.
- `ItemsPerPage` and `ItemsPerPageSizeChoices` page the server chunk. `ItemsTotal` and `PagesTotal` feed the footer readout and page count, exactly like DataTable.
- `QueryUrlTemplate` uses the same placeholders as DataTable: `CarouselUrlPlaceholderPageNumber`, `CarouselUrlPlaceholderItemsPerPage`, and `CarouselUrlPlaceholderSearch`.
- `Filters` renders the shared `FilterBar` below the track, before the pagination. Filter values append to the URL as `key=value` pairs; number and date ranges append as `keyMin` and `keyMax`. The filter bar stays outside the swapped region, so a refresh cannot steal focus from a field you are typing in. `FilterChipStyle` forwards to the shared FilterBar and repaints the active-filter chip; see FilterBar for the fields.
- A search box renders centered above the track when the query URL template carries the search placeholder. It carries a magnifier icon and stretches to the available space. Pass `SearchBox` to replace it. `IsSearchBoxCompact` fixes the box to a narrow width so `SearchBoxAlignment` can place it left, center (the default), or right. `SearchBoxPosition` takes `CarouselSearchBoxPositionTop` (the default) or `CarouselSearchBoxPositionBottom` and places the bar above or below the track.
- `IsAutoplay` advances the window on a timer. `AutoplayIntervalMs` sets the interval, 4000 when unset. `IsAutoplayPausedOnHover` stops the timer while the pointer is over the carousel.
- `IsSwipeEnabled` moves the window on a horizontal drag or swipe. The threshold is 40 pixels.
- `IsPaginationHiddenWhenSinglePage` hides the page-number controls while every item fits on one page. The readout and the items-per-page selector stay.
- Surface styling: `BackgroundColor`, `TextColor`, `BorderRadius` (`uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull`), `PaddingSize` and `ItemPaddingSize` (`uiToolset.PaddingSizeNone` through `uiToolset.PaddingSizeXl`), `GapSize` (`uiToolset.GapSizeNone` through `uiToolset.GapSizeXl`), `ShadowSize` (`uiToolset.ShadowSizeNone` through `uiToolset.ShadowSizeXl`), and `RingColor` with `RingThickness` (`uiToolset.RingThicknessXs` through `uiToolset.RingThicknessXl`).
- Arrow styling: `ArrowsPosition` (`CarouselArrowsPositionOutside`, the default, or `CarouselArrowsPositionInside` to overlay the track), `ArrowsShape` (`uiToolset.ShapeCircular`, the default, `uiToolset.ShapeRounded`, or `uiToolset.ShapeSquare`), `ArrowsSize` (`CarouselArrowsSizeSm/Md/Lg`; the arrow is a fixed-size square, 32 pixels on the default), `ArrowsBackgroundColor`, and `ArrowsIconColor`. A custom background still brightens on hover.
- Dot styling: `DotsPosition` (`CarouselDotsPositionBottom`, the default, or `CarouselDotsPositionTop`), `DotsSize` (`CarouselDotsSizeSm/Md/Lg`), `DotsActiveColor`, and `DotsInactiveColor`.
- Item styling: `ItemBackgroundColor`, `ItemBorderRadius` (`uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull`), `ItemPaddingSize`, `ItemRingColor`, `ItemRingThickness`, and `ItemShadowSize` paint the item wrapper, so a plain renderer still gets a card.
- `EmptyState` renders when the chunk holds no items. `RefreshOnEvents` lists window event names; dispatching one refreshes the carousel.
- The refresh uses `htmx.ajax` when HTMX is present and falls back to `fetch` otherwise, and swaps two regions from one response: the carousel body (`data-ui-carousel`, holding the arrows, track, and dots) and the pagination (`data-ui-carousel-pagination`). The response must contain both. The pagination carries an `id` of the component id plus `-pagination`, and `hx-swap-oob` targets that id, so HTMX swaps it out of band into the component that refreshed while the body takes the selected swap, and the root carries `hx-sync` so a new refresh aborts the one in flight.
- Client state lives in the component root: `pageNumber`, `itemsPerPage`, `searchQuery`, `filterValues`, `windowStart`, `itemsPerView`, and `itemsCount`. A failed refresh shows an inline error with a retry button.
- The arrows carry accessible names, the dots carry `aria-label` and `aria-current`, off-window items are set `inert` so they leave the tab order, and the track is keyboard-operable through the prev and next buttons.
- Pass state paths, the items-per-view values, filter keys, the query URL template, and the class-attribute inputs (`BackgroundColor`, `TextColor`, `RingColor`, `ArrowsBackgroundColor`, `ArrowsIconColor`, `DotsActiveColor`, `DotsInactiveColor`, `ItemBackgroundColor`, `ItemRingColor`, and `FilterDropdownBackgroundColor`) from code, never from request data. The component embeds them into client-side expressions and class attributes.

## CarouselItemTooltip

Wraps one carousel item and shows a tooltip on hover or on focus inside the item.

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

- `Content` is the tooltip text. `ContentHtml` accepts a component instead. Without both, the wrapper renders the item with no tooltip.
- `Position` accepts `uiToolset.TooltipPositionTop` (the default), `uiToolset.TooltipPositionBottom`, `uiToolset.TooltipPositionLeft`, or `uiToolset.TooltipPositionRight`.
- `BackgroundColor` and `TextColor` take color tokens. The defaults render `bg-neutral-800/95` on `text-neutral-50`.
- `RingColor` takes a color token and `RingThickness` takes the `uiToolset.RingThickness*` values. The tooltip ring uses the compact ring scale: `xs` renders `ring-0.5` through `xl` renders `ring-2.5`.
- `MinWidthClass`, `MaxWidthClass`, `MinHeightClass`, and `MaxHeightClass` cap the tooltip size. A long content wraps at the default `max-w-96`.
- The tooltip teleports to the document body with fixed coordinates, so the track viewport, a scroll container, or a modal cannot clip it. Alpine removes the teleported node when the item leaves the DOM, so a refresh does not leak tooltips.

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

Each column takes a `Label`, a `CellRenderer` function, and optional `SortKey`, `Alignment`, `WidthPercent`, width classes, and `CellClass`. `Alignment` takes a `TextAlignment` value: `TextAlignmentLeft`, `TextAlignmentCenter`, or `TextAlignmentRight`. `TextCase` takes a `uiToolset.TextCase*` value and transforms the header labels. The default, `TextCaseNone`, leaves them as typed. `Density` takes a `DataTableDensity` value: `DataTableDensityComfortable` (the default) or `DataTableDensityDense`. `InitialSortDirection` takes a `DataTableSortDirection` value: `DataTableSortDirectionAsc` or `DataTableSortDirectionDesc`. `ItemsPerPage` and each entry in `ItemsPerPageSizeChoices` are `uiStructural.ItemsPerPage` values. Set `PaginationAriaLabel` when a page holds more than one table, so each pagination landmark keeps a unique name.

The `Initial*` fields seed the client state at render time: `InitialFilterValues`, `InitialSearchQuery`, `InitialSortKey`, and `InitialSortDirection`. The server renders the matching rows. `PageNumber` and `ItemsPerPage` also seed the client, but the component reads them to render the pagination readout. `PageNumber` is zero-based: the first page is 0, and the zero value is the first page. `ShouldUseOneBasedPageDisplay` changes only the labels, not the state.

`HeaderClass` adds classes to the header row, `CellClass` adds classes to one column's cells, `RowClassResolver` returns classes for each row from its data, and `IsStriped` adds a zebra stripe. These classes append to elements that already carry base utilities, so when two utilities set the same property the generated stylesheet order decides the winner, not the field order. A cell component that sets its own color wins over the row color, so use `RowClassResolver` for cells that leave the color to the row. The header pins to the top of the scroll viewport by default and paints a blurred, translucent `bg-neutral-950/20` surface on every header cell, so the rows dim as they pass under it and a `HeaderClass` background stays behind it. Set `StickyHeaderBackgroundColor` to replace that surface and `IsHeaderStatic` to let the header scroll with the body. Set `ItemsPerPageDropdownBackgroundColor` to change the items-per-page menu background; the default is `bg-neutral-800/95`. A dropdown inside the table paints above the pinned header and is not clipped by the table surface. `MinWidthClass` and `MaxWidthClass` cap the table surface width. `MinHeightClass` and `MaxHeightClass` cap the scroll viewport height; without them the default cap is `max-h-128`. `MaxVisibleRows` measures the header and the tallest data row at runtime and caps the viewport to that many rows, so the body scrolls under the header; it wins over both height classes. The table renders a default search box when the query URL template carries the search placeholder; it carries a magnifier icon and stretches to the available space. Pass `SearchBox` to replace it. `IsSearchBoxCompact` fixes the box to a narrow width so `SearchBoxAlignment` can place it left (the default), center, or right within the toolbar. `CheckboxShape` accepts `uiToolset.ShapeSquare` (the default), `uiToolset.ShapeRounded`, or `uiToolset.ShapeCircular`; `CheckboxSize` accepts the `uiForm.CheckboxInputSize*` values and defaults to the medium size; `CheckboxCheckedColor` and `CheckboxUncheckedColor` take a color token and default to `secondary-500` and `neutral-50/20`.

The query URL template uses fixed placeholders. Build it from the `DataTableUrlPlaceholder*` constants and name the query keys:

```go
QueryUrlTemplate: "/records?page=" + uiStructural.DataTableUrlPlaceholderPageNumber +
    "&size=" + uiStructural.DataTableUrlPlaceholderItemsPerPage +
    "&sort=" + uiStructural.DataTableUrlPlaceholderSortKey +
    "&direction=" + uiStructural.DataTableUrlPlaceholderSortDirection +
    "&q=" + uiStructural.DataTableUrlPlaceholderSearch
```

`DataTableUrlPlaceholderPageNumber` emits the zero-based page index, where the first page is 0. APIs that treat the first page as 0 receive the value unchanged.

Filter values append to the URL as `key=value` pairs. Number and date ranges append as `keyMin` and `keyMax`. Empty values are omitted. Set `QueryParamName` on a filter to send a different query key.

`FilterChipStyle` forwards to the shared FilterBar and repaints the active-filter chip; see FilterBar for the fields.

The server response must contain one element with the `data-ui-data-table` attribute. The component swaps only that element, so the filter bar, search box, and selection stay in place.

Client state lives in the component root: `pageNumber`, `itemsPerPage`, `sortKey`, `sortDirection`, `searchQuery`, `filterValues`, and `selectedRowIds`. The search box and the bulk action slot bind to those paths.

`RefreshOnEvents` lists window event names. Dispatching one of them refreshes the table. This matches the form-to-display refresh pattern.

A failed refresh shows an inline error row with a retry button. A refresh in flight dims the table and disables the controls.

Pass filter keys, state paths, the query URL template, and `Id` from code, never from request data. The component embeds them into client-side expressions and the root id. `HeaderClass`, `CellClass`, `MinWidthClass`, `MaxWidthClass`, the `RowClassResolver` result, `FilterDropdownBackgroundColor`, `ItemsPerPageDropdownBackgroundColor`, `StickyHeaderBackgroundColor`, `CheckboxCheckedColor`, and `CheckboxUncheckedColor` become HTML class attributes, so keep untrusted data out of them too.

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
- A chip appears when its filter holds a value. The chip remove button clears that filter. The chip row renders only while at least one filter is active, so an idle filter bar keeps its padding symmetric.
- The clear-all button appears when any filter is active.
- `ChipStyle` repaints the active-filter chip. `OuterBackgroundColor`, `OuterRingColor`, `OuterTextColor`, `InnerBackgroundColor`, and `InnerTextColor` take color tokens. `Radius` takes a `uiToolset.BorderRadius*` value and rounds both chip layers; the default is `uiToolset.BorderRadiusMd`. `TextCase` takes a `uiToolset.TextCase*` value and transforms both chip labels; the default is `uiToolset.TextCaseLower`. Every field falls back to the neutral default.
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

- The bound page number is zero-based: the first page is 0. The readout shows the current item range and the total, so it starts at 1.
- `ShouldUseOneBasedPageDisplay` labels the strip from 1 while the bound state stays zero-based. The default labels the first page 0.
- The component derives the page count from `ItemsTotal` and the bound `itemsPerPage`, so the strip and the controls react when the items-per-page value changes. `PagesTotal` is an optional fallback used only when `ItemsTotal` is zero.
- The strip shows the first page, the last page, the pages around the current one, and ellipses for gaps. The current page carries `aria-current="page"`.
- `ItemsPerPageSizeChoices` overrides the default items-per-page choices.
- `DropdownBackgroundColor` paints the items-per-page menu; the default is `bg-neutral-800/95`. DataTable forwards it as `ItemsPerPageDropdownBackgroundColor`.
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

## Tabs

Horizontal and vertical tab headers bound to one selected value. Each tab carries a label, a value, and an optional icon and badge count. The panels switch with Alpine `x-show`.

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

- `SelectedValueTwoWayStatePath` names the Alpine property that holds the active tab value. Clicking a tab writes its `Value` to that path.
- Seed the bound property with one of the item `Value`s. When no value matches, no tab is active; the first tab stays reachable by keyboard so the group never traps focus.
- Each item's `Content` renders in a panel with `x-show` bound to the selected value. Panels stay in the DOM, so their inputs keep state across switches.
- `Orientation` accepts `TabOrientationHorizontal` (the default) or `TabOrientationVertical`. The vertical variant lays the tab list as a column beside the panels, for sidebar-style layouts.
- `Side` accepts `TabSideLeft` (the default) or `TabSideRight`, and only applies with `TabOrientationVertical`. It places the tab list on the right side of the panels.
- `Alignment` accepts `TabAlignmentTop` (the default), `TabAlignmentCenter`, or `TabAlignmentBottom`, and only applies with `TabOrientationVertical`. It aligns the tab list with the tab content.
- `MaxWidthPercent`, `MinWidthPercent`, `MaxHeightPercent`, and `MinHeightPercent` bound the tab list on each axis as whole percentages of the tab area, 1 through 100; 100 renders the full class. `MaxWidthClass`, `MinWidthClass`, `MaxHeightClass`, and `MinHeightClass` take arbitrary Tailwind classes for the same bounds. `MaxVisibleTabs` measures the tallest tab at runtime and caps the vertical list to that many tabs; it wins over the height bounds and applies only to the vertical orientation. The horizontal list scrolls on the x axis and the vertical list scrolls on the y axis when the tabs exceed the bounds.
- Each item's `Icon` renders a Phosphor icon. `IconPosition` accepts `TabIconPositionLeft` (the default) or `TabIconPositionTop`; the top position stacks the icon above the label.
- `BadgeCount` renders a static count on the tab. `BadgeCountOneWayStatePath` renders a count from Alpine state and hides itself at zero. `BadgeBackgroundColor`, `BadgeTextColor`, `BadgeRingColor`, `BadgeRingThickness`, and `BadgeBorderRadius` style the badge; the defaults render a neutral badge.
- `IsUrlHashSynced` reads the URL hash on load and on every hash change, and writes the hash on click, so a link to `#security` opens the Security tab.
- The component follows the ARIA tabs pattern: `role="tablist"`, `role="tab"`, `role="tabpanel"`, `aria-selected`, `aria-controls`, and `aria-labelledby`. Arrow keys move between tabs, Home and End jump to the ends. Horizontal tabs use Left and Right; vertical tabs use Up and Down.
- `AriaLabel` names the tab list. Set a distinct label when a page holds more than one tab group.
- `BackgroundColor` takes a color token, for example `"neutral-50/5"`, and paints the tab list. `BorderRadius` accepts `uiToolset.BorderRadiusNone` through `uiToolset.BorderRadiusFull` (the default is `uiToolset.BorderRadiusMd`) and rounds the tab edges. `TextColor` takes a color token and tints the unselected tab labels.
- `PaddingSize` accepts `uiToolset.PaddingSizeNone` through `uiToolset.PaddingSizeXl` and insets the tab list; the default leaves the list flush. `GapSize` accepts `uiToolset.GapSizeNone` through `uiToolset.GapSizeXl` and sets the gap between the tabs; the default is `uiToolset.GapSizeXs`.
- `RingColor` and `RingThickness` (`uiToolset.RingThicknessXs` through `uiToolset.RingThicknessXl`) frame the tab list when both are set. `ShadowSize` accepts `uiToolset.ShadowSizeNone` through `uiToolset.ShadowSizeXl`.
- The `Content*` fields style the panel area: `ContentBackgroundColor`, `ContentPaddingSize`, `ContentBorderRadius`, `ContentRingColor`, `ContentRingThickness`, `ContentShadowSize`, and `ContentTextColor`. They take the same tokens as the tab list settings. `ContentPaddingSize` replaces the default inset that separates the content from the tab list; the default leaves the content flush.
- `TextCase` takes a `uiToolset.TextCase*` value and transforms the tab labels.
- Pass state paths, tab values, and the class-attribute inputs (`Icon`, `BackgroundColor`, `TextColor`, and the `Badge*` color and ring inputs) from code, never from request data. The component embeds them into client-side expressions and class attributes.

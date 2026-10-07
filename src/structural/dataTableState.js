UiToolset.RegisterAlpineState(() => {
  const unsortedSortIconClass =
    "ph-caret-up-down opacity-0 group-hover:opacity-40 group-focus-visible:opacity-40";

  Alpine.data("dataTable", (settingsScriptId) => ({
    queryUrlTemplate: "",
    refreshDebounceMs: 300,
    refreshOnEvents: [],
    filterQueryParamNames: {},
    pageNumber: 0,
    itemsPerPage: 5,
    sortKey: "",
    sortDirection: "",
    searchQuery: "",
    filterValues: {},
    maxVisibleRows: 0,
    scrollViewportMaxHeightPx: 0,
    scrollViewportResizeObserver: null,
    selectedRowIds: [],
    isLoading: false,
    hasRefreshError: false,
    liveMessage: "",
    refreshTimeoutId: null,
    refreshEventHandlers: [],
    rootElement: null,
    refreshAbortController: null,

    tableRegion() {
      return this.rootElement.querySelector("[data-ui-data-table]");
    },

    tableBody() {
      return this.rootElement.querySelector("tbody");
    },

    scrollViewport() {
      return this.rootElement.querySelector("[data-ui-data-table-scroll]");
    },

    scrollViewportTable() {
      return this.scrollViewport()?.querySelector("table") ?? null;
    },

    applyScrollViewportMaxHeight() {
      if (this.maxVisibleRows === 0) {
        this.scrollViewportMaxHeightPx = 0;
        return;
      }
      const viewport = this.scrollViewport();
      const tableHeader = viewport?.querySelector("thead");
      const dataRows = viewport
        ? Array.from(viewport.querySelectorAll("tbody tr")).filter(
            (row) => !row.querySelector("td[colspan]"),
          )
        : [];
      if (!tableHeader || dataRows.length === 0) {
        this.scrollViewportMaxHeightPx = 0;
        return;
      }
      const rowHeight = Math.max(
        ...dataRows.map((row) => row.getBoundingClientRect().height),
      );
      if (rowHeight === 0) {
        return;
      }
      const visibleRows = Math.min(this.maxVisibleRows, dataRows.length);
      this.scrollViewportMaxHeightPx = Math.ceil(
        tableHeader.getBoundingClientRect().height + rowHeight * visibleRows,
      );
    },

    observeScrollViewportTable() {
      this.scrollViewportResizeObserver?.disconnect();
      const table = this.scrollViewportTable();
      if (!table || typeof ResizeObserver === "undefined") {
        return;
      }
      this.scrollViewportResizeObserver = new ResizeObserver(() =>
        this.applyScrollViewportMaxHeight(),
      );
      this.scrollViewportResizeObserver.observe(table);
    },

    refreshScrollViewportLayout() {
      this.applyScrollViewportMaxHeight();
      this.observeScrollViewportTable();
    },

    currentPageRowIds() {
      const rowCheckboxes =
        this.tableBody()?.querySelectorAll("input[type=checkbox]") ?? [];
      return Array.from(rowCheckboxes).map((checkbox) => checkbox.value);
    },

    isCurrentPageFullySelected() {
      const pageRowIds = this.currentPageRowIds();
      return (
        pageRowIds.length > 0 &&
        pageRowIds.every((rowId) => this.selectedRowIds.includes(rowId))
      );
    },

    isCurrentPagePartiallySelected() {
      const pageRowIds = this.currentPageRowIds();
      const selectedCount = pageRowIds.filter((rowId) =>
        this.selectedRowIds.includes(rowId),
      ).length;
      return selectedCount > 0 && selectedCount < pageRowIds.length;
    },

    toggleCurrentPageSelection() {
      const pageRowIds = this.currentPageRowIds();
      if (this.isCurrentPageFullySelected()) {
        this.selectedRowIds = this.selectedRowIds.filter(
          (rowId) => !pageRowIds.includes(rowId),
        );
        return;
      }
      this.selectedRowIds = [
        ...new Set([...this.selectedRowIds, ...pageRowIds]),
      ];
    },

    ariaSortFor(columnSortKey) {
      if (this.sortKey !== columnSortKey) {
        return "none";
      }
      if (this.sortDirection === "asc") {
        return "ascending";
      }
      if (this.sortDirection === "desc") {
        return "descending";
      }
      return "none";
    },

    sortIconClassFor(columnSortKey) {
      if (this.sortKey !== columnSortKey) {
        return unsortedSortIconClass;
      }
      if (this.sortDirection === "asc") {
        return "ph-caret-up";
      }
      if (this.sortDirection === "desc") {
        return "ph-caret-down";
      }
      return unsortedSortIconClass;
    },

    async fetchTableRegion(refreshUrl, abortSignal) {
      const responseDocument =
        await UiToolset.ServerFragmentRefreshComponent.fragmentFetcher(
          refreshUrl,
          abortSignal,
        );
      const freshTableRegion = responseDocument.querySelector(
        "[data-ui-data-table]",
      );
      if (!freshTableRegion) {
        throw new Error("Response is missing the data-ui-data-table fragment");
      }
      return freshTableRegion;
    },

    async refresh() {
      if (!this.queryUrlTemplate) {
        return;
      }
      clearTimeout(this.refreshTimeoutId);
      this.refreshAbortController?.abort();
      const abortController = new AbortController();
      this.refreshAbortController = abortController;
      this.isLoading = true;
      this.hasRefreshError = false;
      const refreshUrl =
        UiToolset.ServerFragmentRefreshComponent.refreshUrlBuilder(this);

      try {
        if (window.htmx?.ajax) {
          await window.htmx.ajax("GET", refreshUrl, {
            source: this.tableRegion(),
            target: this.tableRegion(),
            select: "[data-ui-data-table]",
            swap: "outerHTML",
          });
          if (abortController.signal.aborted) {
            return;
          }
          if (!this.hasRefreshError) {
            this.liveMessage = "Table refreshed";
          }
          this.refreshScrollViewportLayout();
          return;
        }
        const freshTableRegion = await this.fetchTableRegion(
          refreshUrl,
          abortController.signal,
        );
        if (abortController.signal.aborted) {
          return;
        }
        this.tableRegion().replaceWith(freshTableRegion);
        this.liveMessage = "Table refreshed";
        this.refreshScrollViewportLayout();
      } catch (refreshError) {
        UiToolset.ServerFragmentRefreshComponent.failureHandler(
          this,
          abortController,
          refreshError,
          "DataTable",
        );
      } finally {
        UiToolset.ServerFragmentRefreshComponent.loadingClearer(
          this,
          abortController,
        );
      }
    },

    requestRefresh() {
      UiToolset.ServerFragmentRefreshComponent.requestDebouncer(this);
    },

    refreshFromFirstPage() {
      this.pageNumber = 0;
      UiToolset.ServerFragmentRefreshComponent.requestDebouncer(this);
    },

    toggleSort(columnSortKey) {
      this.pageNumber = 0;
      if (this.sortKey !== columnSortKey) {
        this.sortKey = columnSortKey;
        this.sortDirection = "asc";
        this.requestRefresh();
        return;
      }
      if (this.sortDirection === "asc") {
        this.sortDirection = "desc";
        this.requestRefresh();
        return;
      }
      this.sortKey = "";
      this.sortDirection = "";
      this.requestRefresh();
    },

    destroy() {
      UiToolset.ServerFragmentRefreshComponent.handlersDetacher(this);
      this.scrollViewportResizeObserver?.disconnect();
    },

    init() {
      this.rootElement = this.$el;
      const clientSettings =
        UiToolset.ServerFragmentRefreshComponent.settingsResolver(
          this,
          settingsScriptId,
          "DataTable",
        );
      if (!clientSettings) {
        return;
      }

      const initialState = clientSettings.initialState || {};
      this.pageNumber = initialState.pageNumber ?? 0;
      this.itemsPerPage = initialState.itemsPerPage ?? 5;
      this.sortKey = initialState.sortKey ?? "";
      this.sortDirection = initialState.sortDirection ?? "";
      this.searchQuery = initialState.searchQuery ?? "";
      this.filterValues = initialState.filterValues ?? {};
      this.maxVisibleRows = clientSettings.maxVisibleRows ?? 0;

      UiToolset.ServerFragmentRefreshComponent.eventsWatcher(this);

      this.$watch("selectedRowIds", () => {
        this.liveMessage =
          this.selectedRowIds.length === 1
            ? "1 row selected"
            : `${this.selectedRowIds.length} rows selected`;
      });

      UiToolset.ServerFragmentRefreshComponent.errorHandlersAttacher(this);

      this.refreshScrollViewportLayout();
    },
  }));
});

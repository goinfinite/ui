UiToolset.RegisterAlpineState(() => {
  Alpine.data("tabs", (maxVisibleTabs) => ({
    maxVisibleTabs: 0,
    tabsViewportMaxHeightPx: 0,
    viewportResizeObserver: null,

    init() {
      this.maxVisibleTabs = Number(maxVisibleTabs) || 0;
      if (this.maxVisibleTabs === 0) {
        return;
      }
      this.applyViewportMaxHeight();
      this.observeTabList();
    },

    destroy() {
      this.viewportResizeObserver?.disconnect();
    },

    tabList() {
      return this.$el;
    },

    tabButtons() {
      return Array.from(this.tabList().querySelectorAll("[role=tab]"));
    },

    observeTabList() {
      this.viewportResizeObserver?.disconnect();
      if (typeof ResizeObserver === "undefined") {
        return;
      }
      this.viewportResizeObserver = new ResizeObserver(() =>
        this.applyViewportMaxHeight(),
      );
      this.viewportResizeObserver.observe(this.tabList());
    },

    applyViewportMaxHeight() {
      const tabList = this.tabList();
      const tabButtons = this.tabButtons();
      if (tabButtons.length === 0) {
        this.tabsViewportMaxHeightPx = 0;
        return;
      }
      const tabHeight = Math.max(
        ...tabButtons.map(
          (tabButton) => tabButton.getBoundingClientRect().height,
        ),
      );
      if (tabHeight === 0) {
        return;
      }
      const tabListStyles = getComputedStyle(tabList);
      const rowGap = Number.parseFloat(tabListStyles.rowGap) || 0;
      const verticalPadding =
        (Number.parseFloat(tabListStyles.paddingTop) || 0) +
        (Number.parseFloat(tabListStyles.paddingBottom) || 0);
      const visibleTabCount = Math.min(this.maxVisibleTabs, tabButtons.length);
      this.tabsViewportMaxHeightPx = Math.ceil(
        tabHeight * visibleTabCount +
          rowGap * (visibleTabCount - 1) +
          verticalPadding,
      );
    },
  }));
});

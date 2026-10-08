UiToolset.RegisterAlpineState(() => {
  Alpine.data("carousel", (settingsScriptId) => ({
    queryUrlTemplate: "",
    refreshDebounceMs: 300,
    refreshOnEvents: [],
    filterQueryParamNames: {},
    pageNumber: 0,
    itemsPerPage: 5,
    searchQuery: "",
    filterValues: {},
    isLoading: false,
    hasRefreshError: false,
    liveMessage: "",
    refreshTimeoutId: null,
    refreshEventHandlers: [],
    rootElement: null,
    refreshAbortController: null,

    itemsPerViewSettings: { base: 1 },
    itemsPerView: 1,
    windowStart: 0,
    itemsCount: 0,
    isAutoplay: false,
    isAutoplayPausedOnHover: false,
    isSwipeEnabled: false,
    autoplayIntervalMs: 4000,
    autoplayIntervalId: null,
    isPointerOver: false,
    swipeStartX: null,
    viewportWidth: 0,
    resizeObserver: null,
    resizeHandler: null,

    get dotIndexes() {
      const windowCount = this.itemsCount - this.itemsPerView + 1;
      if (windowCount < 1) {
        return [];
      }
      return Array.from({ length: windowCount }, (_, index) => index);
    },

    get canSlidePrevious() {
      return this.windowStart > 0;
    },

    get canSlideNext() {
      return this.windowStart < this.maxWindowStart();
    },

    carouselRegion() {
      return this.rootElement.querySelector("[data-ui-carousel]");
    },

    paginationRegion() {
      return this.rootElement.querySelector("[data-ui-carousel-pagination]");
    },

    trackElement() {
      return this.rootElement.querySelector("[data-ui-carousel-track]");
    },

    viewportElement() {
      return this.rootElement.querySelector("[x-ref=viewport]");
    },

    readItemsCount() {
      const trackElement = this.trackElement();
      if (!trackElement) {
        return 0;
      }
      return trackElement.children.length;
    },

    maxWindowStart() {
      return Math.max(0, this.itemsCount - this.itemsPerView);
    },

    hasMultipleWindows() {
      return this.itemsCount > this.itemsPerView;
    },

    windowStartResolver(candidateStart) {
      return Math.max(0, Math.min(candidateStart, this.maxWindowStart()));
    },

    itemsPerViewResolver() {
      const breakpoints = [
        { minWidthPx: 640, itemsPerView: this.itemsPerViewSettings.sm },
        { minWidthPx: 768, itemsPerView: this.itemsPerViewSettings.md },
        { minWidthPx: 1024, itemsPerView: this.itemsPerViewSettings.lg },
        { minWidthPx: 1280, itemsPerView: this.itemsPerViewSettings.xl },
        { minWidthPx: 1536, itemsPerView: this.itemsPerViewSettings.twoXl },
      ];
      let resolvedItemsPerView = this.itemsPerViewSettings.base || 1;
      for (const breakpoint of breakpoints) {
        if (
          breakpoint.itemsPerView > 0 &&
          window.innerWidth >= breakpoint.minWidthPx
        ) {
          resolvedItemsPerView = breakpoint.itemsPerView;
        }
      }
      return resolvedItemsPerView;
    },

    updateItemsPerView() {
      const viewportElement = this.viewportElement();
      if (!viewportElement) {
        return;
      }
      this.viewportWidth = viewportElement.clientWidth;
      this.itemsPerView = this.itemsPerViewResolver();
      this.windowStart = this.windowStartResolver(this.windowStart);
    },

    updateTrackTransform() {
      const trackElement = this.trackElement();
      if (!trackElement || trackElement.children.length === 0) {
        return;
      }
      const trackStyle = window.getComputedStyle(trackElement);
      const columnGap = Number.parseFloat(trackStyle.columnGap) || 0;
      const visibleGapTotal = columnGap * (this.itemsPerView - 1);
      const itemWidth =
        (this.viewportWidth - visibleGapTotal) / this.itemsPerView;
      for (
        let itemIndex = 0;
        itemIndex < trackElement.children.length;
        itemIndex++
      ) {
        const itemElement = trackElement.children[itemIndex];
        itemElement.style.width = `${itemWidth}px`;
        const isInWindow =
          itemIndex >= this.windowStart &&
          itemIndex < this.windowStart + this.itemsPerView;
        itemElement.inert = !isInWindow;
      }
      const singleItemStep = itemWidth + columnGap;
      const translateX = this.windowStart * singleItemStep;
      trackElement.style.transform = `translateX(-${translateX}px)`;
    },

    updateCarouselLayout() {
      this.updateItemsPerView();
      this.updateTrackTransform();
    },

    goToWindow(candidateStart) {
      this.windowStart = this.windowStartResolver(candidateStart);
      this.updateTrackTransform();
    },

    slideNext() {
      this.goToWindow(this.windowStart + 1);
    },

    slidePrevious() {
      this.goToWindow(this.windowStart - 1);
    },

    startSwipe(event) {
      if (!this.isSwipeEnabled) {
        return;
      }
      this.swipeStartX = event.clientX;
    },

    finishSwipe(event) {
      if (!this.isSwipeEnabled || this.swipeStartX === null) {
        return;
      }
      const swipeDistance = event.clientX - this.swipeStartX;
      this.swipeStartX = null;
      const swipeThresholdPx = 40;
      if (Math.abs(swipeDistance) < swipeThresholdPx) {
        return;
      }
      if (swipeDistance < 0) {
        this.slideNext();
        return;
      }
      this.slidePrevious();
    },

    cancelSwipe() {
      this.swipeStartX = null;
    },

    enterPointerRegion() {
      this.isPointerOver = true;
      if (this.isAutoplayPausedOnHover) {
        this.stopAutoplay();
      }
    },

    leavePointerRegion() {
      this.isPointerOver = false;
      if (this.isAutoplayPausedOnHover) {
        this.startAutoplay();
      }
    },

    startAutoplay() {
      if (!this.isAutoplay || !this.hasMultipleWindows()) {
        return;
      }
      if (this.isAutoplayPausedOnHover && this.isPointerOver) {
        return;
      }
      this.stopAutoplay();
      this.autoplayIntervalId = setInterval(() => {
        const nextStart =
          this.windowStart >= this.maxWindowStart() ? 0 : this.windowStart + 1;
        this.goToWindow(nextStart);
      }, this.autoplayIntervalMs);
    },

    stopAutoplay() {
      if (this.autoplayIntervalId === null) {
        return;
      }
      clearInterval(this.autoplayIntervalId);
      this.autoplayIntervalId = null;
    },

    applyRefreshedDocument(responseDocument) {
      const freshCarouselRegion =
        responseDocument.querySelector("[data-ui-carousel]");
      if (!freshCarouselRegion) {
        throw new Error("Response is missing the data-ui-carousel fragment");
      }
      const freshPaginationRegion = responseDocument.querySelector(
        "[data-ui-carousel-pagination]",
      );
      if (!freshPaginationRegion) {
        throw new Error(
          "Response is missing the data-ui-carousel-pagination fragment",
        );
      }
      const currentPaginationRegion = this.paginationRegion();
      this.carouselRegion().replaceWith(freshCarouselRegion);
      if (currentPaginationRegion) {
        currentPaginationRegion.replaceWith(freshPaginationRegion);
      }
    },

    afterRegionReplaced() {
      this.observeViewport();
      this.itemsCount = this.readItemsCount();
      this.windowStart = 0;
      this.updateCarouselLayout();
      this.startAutoplay();
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
          const previousPaginationRegion = this.paginationRegion();
          await window.htmx.ajax("GET", refreshUrl, {
            source: this.carouselRegion(),
            target: this.carouselRegion(),
            select: "[data-ui-carousel]",
            swap: "outerHTML",
          });
          if (abortController.signal.aborted) {
            return;
          }
          if (this.hasRefreshError) {
            return;
          }
          const refreshedPaginationRegion = this.paginationRegion();
          if (
            !refreshedPaginationRegion ||
            refreshedPaginationRegion === previousPaginationRegion
          ) {
            throw new Error("PaginationRegionNotReplaced");
          }
          this.liveMessage = "Carousel refreshed";
          this.afterRegionReplaced();
          return;
        }
        const responseDocument =
          await UiToolset.ServerFragmentRefreshComponent.fragmentFetcher(
            refreshUrl,
            abortController.signal,
          );
        if (abortController.signal.aborted) {
          return;
        }
        this.applyRefreshedDocument(responseDocument);
        this.liveMessage = "Carousel refreshed";
        this.afterRegionReplaced();
      } catch (refreshError) {
        UiToolset.ServerFragmentRefreshComponent.failureHandler(
          this,
          abortController,
          refreshError,
          "Carousel",
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

    destroy() {
      UiToolset.ServerFragmentRefreshComponent.handlersDetacher(this);
      if (this.resizeHandler) {
        window.removeEventListener("resize", this.resizeHandler);
      }
      this.resizeObserver?.disconnect();
      this.stopAutoplay();
    },

    init() {
      this.rootElement = this.$el;
      const clientSettings =
        UiToolset.ServerFragmentRefreshComponent.settingsResolver(
          this,
          settingsScriptId,
          "Carousel",
        );
      if (!clientSettings) {
        return;
      }
      this.itemsPerViewSettings = clientSettings.itemsPerView || { base: 1 };
      this.isAutoplay = clientSettings.isAutoplay || false;
      this.isAutoplayPausedOnHover =
        clientSettings.isAutoplayPausedOnHover || false;
      this.isSwipeEnabled = clientSettings.isSwipeEnabled || false;
      this.autoplayIntervalMs = clientSettings.autoplayIntervalMs || 4000;

      const initialState = clientSettings.initialState || {};
      this.pageNumber = initialState.pageNumber ?? 0;
      this.itemsPerPage = initialState.itemsPerPage ?? 5;
      this.searchQuery = initialState.searchQuery ?? "";
      this.filterValues = initialState.filterValues ?? {};

      UiToolset.ServerFragmentRefreshComponent.eventsWatcher(this);

      this.resizeObserver = new ResizeObserver(() =>
        this.updateCarouselLayout(),
      );
      this.observeViewport();
      this.resizeHandler = () => this.updateCarouselLayout();
      window.addEventListener("resize", this.resizeHandler);

      this.itemsCount = this.readItemsCount();
      this.updateCarouselLayout();
      this.startAutoplay();

      UiToolset.ServerFragmentRefreshComponent.errorHandlersAttacher(this);
    },

    observeViewport() {
      this.resizeObserver.disconnect();
      const viewportElement = this.viewportElement();
      if (viewportElement) {
        this.resizeObserver.observe(viewportElement);
      }
    },
  }));
});

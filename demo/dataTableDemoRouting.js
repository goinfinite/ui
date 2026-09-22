UiToolset.RegisterAlpineState(() => {
  const refreshFile = "dataTableDemoRefresh.html";
  const defaultItemsPerPage = 5;

  function resolveFragmentPath(requestUrl) {
    const request = new URL(requestUrl, window.location.href);
    const itemsPerPage = Number(
      request.searchParams.get("itemsPerPage") ?? defaultItemsPerPage,
    );
    if (itemsPerPage > defaultItemsPerPage) {
      return `assets/dataTableDemoRefreshAll.html${request.search}`;
    }
    const pageNumber = Number(request.searchParams.get("page") ?? "1");
    return `assets/dataTableDemoRefreshPage${pageNumber}.html${request.search}`;
  }

  document.addEventListener("htmx:configRequest", (event) => {
    if (!event.detail.path.includes(refreshFile)) {
      return;
    }
    event.detail.path = resolveFragmentPath(event.detail.path);
  });

  const originalFetch = window.fetch.bind(window);
  window.fetch = (input, init) => {
    if (typeof input !== "string" || !input.includes(refreshFile)) {
      return originalFetch(input, init);
    }
    return originalFetch(resolveFragmentPath(input), init);
  };
});

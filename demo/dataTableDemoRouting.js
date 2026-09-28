UiToolset.RegisterAlpineState(() => {
  const refreshAssetFile = "dataTableDemoRefreshFragments.json";
  const defaultItemsPerPage = 5;
  const defaultPageNumber = 1;

  function isRefreshAssetRequest(requestUrl) {
    return (
      typeof requestUrl === "string" && requestUrl.includes(refreshAssetFile)
    );
  }

  function resolveFragmentKey(requestUrl) {
    const request = new URL(requestUrl, window.location.href);
    const itemsPerPage = Number(
      request.searchParams.get("itemsPerPage") ?? defaultItemsPerPage,
    );
    const pageNumber = Number(
      request.searchParams.get("page") ?? defaultPageNumber,
    );
    return `${pageNumber}-${itemsPerPage}`;
  }

  function resolveFragmentHtml(assetBody, requestUrl) {
    const fragmentKey = resolveFragmentKey(requestUrl);
    const fragmentHtml = JSON.parse(assetBody)[fragmentKey];
    if (fragmentHtml === undefined) {
      throw new Error(`DataTableDemoFragmentMissing: ${fragmentKey}`);
    }
    return fragmentHtml;
  }

  const originalFetch = window.fetch.bind(window);
  window.fetch = async (input, init) => {
    const response = await originalFetch(input, init);
    if (!isRefreshAssetRequest(input) || !response.ok) {
      return response;
    }
    return new Response(resolveFragmentHtml(await response.text(), input), {
      status: 200,
      headers: { "content-type": "text/html" },
    });
  };

  document.addEventListener("htmx:beforeSwap", (event) => {
    const requestUrl = event.detail.xhr?.responseURL;
    if (!isRefreshAssetRequest(requestUrl) || event.detail.xhr.status !== 200) {
      return;
    }
    try {
      event.detail.serverResponse = resolveFragmentHtml(
        event.detail.serverResponse,
        requestUrl,
      );
    } catch (resolveError) {
      console.error(resolveError.message);
      event.detail.shouldSwap = false;
      event.detail.isError = true;
    }
  });
});

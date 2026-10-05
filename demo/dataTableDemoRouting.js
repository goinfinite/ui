UiToolset.RegisterAlpineState(() => {
  const refreshAssets = [
    { assetFile: "dataTableDemoRefreshFragments.json", defaultItemsPerPage: 5 },
    { assetFile: "carouselDemoRefreshFragments.json", defaultItemsPerPage: 6 },
  ];
  const defaultPageNumber = 0;

  function refreshAssetForRequest(requestUrl) {
    if (typeof requestUrl !== "string") {
      return undefined;
    }
    return refreshAssets.find((refreshAsset) =>
      requestUrl.includes(refreshAsset.assetFile),
    );
  }

  function resolveFragmentKey(refreshAsset, request) {
    const itemsPerPage = Number(
      request.searchParams.get("itemsPerPage") ??
        refreshAsset.defaultItemsPerPage,
    );
    const pageNumber = Number(
      request.searchParams.get("page") ?? defaultPageNumber,
    );
    return `${pageNumber}-${itemsPerPage}`;
  }

  function resolveFragmentHtml(refreshAsset, assetBody, requestUrl) {
    const request = new URL(requestUrl, window.location.href);
    const fragmentKey = resolveFragmentKey(refreshAsset, request);
    const fragmentHtml = JSON.parse(assetBody)[fragmentKey];
    if (fragmentHtml === undefined) {
      throw new Error(
        `DemoRefreshFragmentMissing (${refreshAsset.assetFile}): ${fragmentKey}`,
      );
    }
    const componentId = request.searchParams.get("componentId");
    if (componentId === null) {
      return fragmentHtml;
    }
    return fragmentHtml.replaceAll(
      "carousel-demo-fragment-carousel",
      componentId,
    );
  }

  const originalFetch = window.fetch.bind(window);
  window.fetch = async (input, init) => {
    const response = await originalFetch(input, init);
    const refreshAsset = refreshAssetForRequest(input);
    if (!refreshAsset || !response.ok) {
      return response;
    }
    return new Response(
      resolveFragmentHtml(refreshAsset, await response.text(), input),
      {
        status: 200,
        headers: { "content-type": "text/html" },
      },
    );
  };

  document.addEventListener("htmx:beforeSwap", (event) => {
    const requestUrl = event.detail.xhr?.responseURL;
    const refreshAsset = refreshAssetForRequest(requestUrl);
    if (!refreshAsset || event.detail.xhr.status !== 200) {
      return;
    }
    try {
      event.detail.serverResponse = resolveFragmentHtml(
        refreshAsset,
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

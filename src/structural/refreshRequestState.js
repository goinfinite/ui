UiToolset.RegisterAlpineState(() => {
  window.UiToolset.BuildRefreshUrl = buildRefreshUrl;
  window.UiToolset.FetchRefreshFragment = fetchRefreshFragment;
});

function resolveRefreshPlaceholderValue(refreshState, placeholderName) {
  switch (placeholderName) {
    case "pageNumber":
      return String(refreshState.pageNumber);
    case "itemsPerPage":
      return String(refreshState.itemsPerPage);
    case "sortKey":
      return refreshState.sortKey ?? "";
    case "sortDirection":
      return refreshState.sortDirection ?? "";
    case "search":
      return refreshState.searchQuery ?? "";
    default:
      return "";
  }
}

function resolveRefreshTemplatePair(refreshState, queryPair) {
  const separatorIndex = queryPair.indexOf("=");
  if (separatorIndex === -1) {
    return queryPair;
  }
  const pairValue = queryPair.slice(separatorIndex + 1);
  const placeholderPattern =
    /^\{(pageNumber|itemsPerPage|sortKey|sortDirection|search)\}$/;
  const placeholderMatch = pairValue.match(placeholderPattern);
  if (!placeholderMatch) {
    return queryPair;
  }
  const placeholderValue = resolveRefreshPlaceholderValue(
    refreshState,
    placeholderMatch[1],
  );
  if (placeholderValue === "") {
    return null;
  }
  const pairKey = queryPair.slice(0, separatorIndex + 1);
  return pairKey + encodeURIComponent(placeholderValue);
}

function isBlankFilterValue(filterValue) {
  return (
    filterValue === null || filterValue === undefined || filterValue === ""
  );
}

function buildRefreshFilterPairs(refreshState) {
  const filterPairs = [];
  const filterValues = refreshState.filterValues ?? {};
  const filterQueryParamNames = refreshState.filterQueryParamNames ?? {};
  for (const [filterKey, filterValue] of Object.entries(filterValues)) {
    const queryParamName = filterQueryParamNames[filterKey] ?? filterKey;
    const encodedQueryParamName = encodeURIComponent(queryParamName);
    if (Array.isArray(filterValue)) {
      for (const singleValue of filterValue) {
        if (isBlankFilterValue(singleValue)) {
          continue;
        }
        filterPairs.push(
          `${encodedQueryParamName}=${encodeURIComponent(singleValue)}`,
        );
      }
      continue;
    }
    if (filterValue && typeof filterValue === "object") {
      if (!isBlankFilterValue(filterValue.min)) {
        const encodedMin = encodeURIComponent(filterValue.min);
        filterPairs.push(`${encodedQueryParamName}Min=${encodedMin}`);
      }
      if (!isBlankFilterValue(filterValue.max)) {
        const encodedMax = encodeURIComponent(filterValue.max);
        filterPairs.push(`${encodedQueryParamName}Max=${encodedMax}`);
      }
      continue;
    }
    if (isBlankFilterValue(filterValue)) {
      continue;
    }
    const encodedValue = encodeURIComponent(filterValue);
    filterPairs.push(`${encodedQueryParamName}=${encodedValue}`);
  }
  return filterPairs;
}

async function fetchRefreshFragment(refreshUrl, abortSignal) {
  const response = await fetch(refreshUrl, {
    headers: { "HX-Request": "true" },
    signal: abortSignal,
  });
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  const responseText = await response.text();
  return new DOMParser().parseFromString(responseText, "text/html");
}

function buildRefreshUrl(refreshState) {
  const [baseUrl, queryString = ""] = refreshState.queryUrlTemplate.split("?");
  const queryPairs = queryString ? queryString.split("&") : [];
  const resolvedPairs = queryPairs
    .map((queryPair) => resolveRefreshTemplatePair(refreshState, queryPair))
    .filter((queryPair) => queryPair !== null);
  const allPairs = [...resolvedPairs, ...buildRefreshFilterPairs(refreshState)];
  if (allPairs.length === 0) {
    return baseUrl;
  }
  return `${baseUrl}?${allPairs.join("&")}`;
}

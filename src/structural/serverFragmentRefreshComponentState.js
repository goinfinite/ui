UiToolset.RegisterAlpineState(() => {
  window.UiToolset.ServerFragmentRefreshComponent = {
    settingsResolver,
    eventsWatcher,
    errorHandlersAttacher,
    handlersDetacher,
    requestDebouncer,
    failureHandler,
    loadingClearer,
    refreshUrlBuilder,
    fragmentFetcher,
  };
});

function settingsResolver(alpineComponent, settingsScriptId, componentLabel) {
  const settingsElement = [...alpineComponent.rootElement.children].find(
    (child) => child.id === settingsScriptId,
  );
  if (!settingsElement) {
    console.error(
      `${componentLabel}SettingsScriptMissing: ${settingsScriptId}`,
    );
    return null;
  }
  let clientSettings;
  try {
    clientSettings = JSON.parse(settingsElement.textContent);
  } catch (parseError) {
    console.error(
      `${componentLabel}SettingsParseFailed: ${parseError.message}`,
    );
    return null;
  }
  alpineComponent.queryUrlTemplate = clientSettings.queryUrlTemplate || "";
  alpineComponent.refreshDebounceMs = clientSettings.refreshDebounceMs || 300;
  alpineComponent.refreshOnEvents = clientSettings.refreshOnEvents || [];
  alpineComponent.filterQueryParamNames =
    clientSettings.filterQueryParamNames || {};
  return clientSettings;
}

function eventsWatcher(alpineComponent) {
  for (const refreshEventName of alpineComponent.refreshOnEvents) {
    const refreshEventHandler = () => alpineComponent.requestRefresh();
    alpineComponent.refreshEventHandlers.push({
      refreshEventName,
      refreshEventHandler,
    });
    window.addEventListener(refreshEventName, refreshEventHandler);
  }
}

function errorHandlersAttacher(alpineComponent) {
  alpineComponent.htmxErrorHandler = () => {
    alpineComponent.hasRefreshError = true;
    alpineComponent.isLoading = false;
  };
  alpineComponent.rootElement.addEventListener(
    "htmx:responseError",
    alpineComponent.htmxErrorHandler,
  );
  alpineComponent.rootElement.addEventListener(
    "htmx:sendError",
    alpineComponent.htmxErrorHandler,
  );
}

function handlersDetacher(alpineComponent) {
  for (const {
    refreshEventName,
    refreshEventHandler,
  } of alpineComponent.refreshEventHandlers) {
    window.removeEventListener(refreshEventName, refreshEventHandler);
  }
  alpineComponent.rootElement.removeEventListener(
    "htmx:responseError",
    alpineComponent.htmxErrorHandler,
  );
  alpineComponent.rootElement.removeEventListener(
    "htmx:sendError",
    alpineComponent.htmxErrorHandler,
  );
  alpineComponent.refreshAbortController?.abort();
  clearTimeout(alpineComponent.refreshTimeoutId);
}

function requestDebouncer(alpineComponent) {
  clearTimeout(alpineComponent.refreshTimeoutId);
  alpineComponent.refreshTimeoutId = setTimeout(
    () => alpineComponent.refresh(),
    alpineComponent.refreshDebounceMs,
  );
}

function failureHandler(
  alpineComponent,
  abortController,
  refreshError,
  componentLabel,
) {
  if (abortController.signal.aborted) {
    return;
  }
  alpineComponent.hasRefreshError = true;
  console.error(
    `${componentLabel}RefreshFailed: ${refreshError?.message ?? refreshError}`,
  );
}

function loadingClearer(alpineComponent, abortController) {
  if (!abortController.signal.aborted) {
    alpineComponent.isLoading = false;
  }
}

function resolvePlaceholderValue(alpineComponent, placeholderName) {
  switch (placeholderName) {
    case "pageNumber":
      return String(alpineComponent.pageNumber);
    case "itemsPerPage":
      return String(alpineComponent.itemsPerPage);
    case "sortKey":
      return alpineComponent.sortKey ?? "";
    case "sortDirection":
      return alpineComponent.sortDirection ?? "";
    case "search":
      return alpineComponent.searchQuery ?? "";
    default:
      return "";
  }
}

function resolveTemplatePair(alpineComponent, queryPair) {
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
  const placeholderValue = resolvePlaceholderValue(
    alpineComponent,
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

function buildFilterPairs(alpineComponent) {
  const filterPairs = [];
  const filterValues = alpineComponent.filterValues ?? {};
  const filterQueryParamNames = alpineComponent.filterQueryParamNames ?? {};
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

function refreshUrlBuilder(alpineComponent) {
  const [baseUrl, queryString = ""] =
    alpineComponent.queryUrlTemplate.split("?");
  const queryPairs = queryString ? queryString.split("&") : [];
  const resolvedPairs = queryPairs
    .map((queryPair) => resolveTemplatePair(alpineComponent, queryPair))
    .filter((queryPair) => queryPair !== null);
  const allPairs = [...resolvedPairs, ...buildFilterPairs(alpineComponent)];
  if (allPairs.length === 0) {
    return baseUrl;
  }
  return `${baseUrl}?${allPairs.join("&")}`;
}

async function fragmentFetcher(refreshUrl, abortSignal) {
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

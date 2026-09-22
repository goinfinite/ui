UiToolset.RegisterAlpineState(() => {
  window.UiToolset.PaginationPageStripBuilder = paginationPageStripBuilder;
});

function paginationShownPageNumbersResolver(currentPageNumber, pagesTotal) {
  const candidatePageNumbers = [
    1,
    currentPageNumber - 1,
    currentPageNumber,
    currentPageNumber + 1,
    pagesTotal,
  ];
  const inRangePageNumbers = candidatePageNumbers.filter(
    (pageNumber) => pageNumber >= 1 && pageNumber <= pagesTotal,
  );
  return [...new Set(inRangePageNumbers)].sort(
    (leftPageNumber, rightPageNumber) => leftPageNumber - rightPageNumber,
  );
}

function paginationPageStripItemBuilder(pageNumber) {
  return {
    key: `page${pageNumber}`,
    label: String(pageNumber),
    pageNumber: pageNumber,
  };
}

function paginationEllipsisStripItemBuilder(previousPageNumber) {
  return {
    key: `ellipsis${previousPageNumber}`,
    label: "…",
    pageNumber: null,
  };
}

function paginationPageStripBuilder(currentPageNumber, pagesTotal) {
  if (pagesTotal < 1) {
    return [];
  }
  const shownPageNumbers = paginationShownPageNumbersResolver(
    currentPageNumber,
    pagesTotal,
  );
  const stripItems = [];
  for (let pageIndex = 0; pageIndex < shownPageNumbers.length; pageIndex++) {
    const pageNumber = shownPageNumbers[pageIndex];
    if (pageIndex > 0) {
      const previousPageNumber = shownPageNumbers[pageIndex - 1];
      const skippedPageCount = pageNumber - previousPageNumber - 1;
      if (skippedPageCount === 1) {
        stripItems.push(paginationPageStripItemBuilder(previousPageNumber + 1));
      }
      if (skippedPageCount > 1) {
        stripItems.push(paginationEllipsisStripItemBuilder(previousPageNumber));
      }
    }
    stripItems.push(paginationPageStripItemBuilder(pageNumber));
  }
  return stripItems;
}

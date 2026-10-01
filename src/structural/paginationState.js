UiToolset.RegisterAlpineState(() => {
  window.UiToolset.PaginationPageStripBuilder = paginationPageStripBuilder;
});

function paginationShownPageNumbersResolver(currentPageNumber, pagesTotal) {
  const candidatePageNumbers = [
    0,
    currentPageNumber - 1,
    currentPageNumber,
    currentPageNumber + 1,
    pagesTotal - 1,
  ];
  const inRangePageNumbers = candidatePageNumbers.filter(
    (pageNumber) => pageNumber >= 0 && pageNumber <= pagesTotal - 1,
  );
  return [...new Set(inRangePageNumbers)].sort(
    (leftPageNumber, rightPageNumber) => leftPageNumber - rightPageNumber,
  );
}

function paginationPageStripItemBuilder(pageNumber, displayBase) {
  return {
    key: `page${pageNumber}`,
    label: String(pageNumber + displayBase),
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

function paginationPageStripBuilder(
  currentPageNumber,
  pagesTotal,
  displayBase,
) {
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
        stripItems.push(
          paginationPageStripItemBuilder(previousPageNumber + 1, displayBase),
        );
      }
      if (skippedPageCount > 1) {
        stripItems.push(paginationEllipsisStripItemBuilder(previousPageNumber));
      }
    }
    stripItems.push(paginationPageStripItemBuilder(pageNumber, displayBase));
  }
  return stripItems;
}

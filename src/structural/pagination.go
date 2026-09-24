package uiStructural

import "strconv"

var paginationDefaultItemsPerPageSizeChoices = []uint{5, 10, 30, 50}

func paginationPagesTotalExpressionBuilder(
	itemsPerPagePath string, itemsTotal, pagesTotal uint,
) string {
	if itemsTotal == 0 {
		return strconv.FormatUint(uint64(pagesTotal), 10)
	}
	itemsTotalText := strconv.FormatUint(uint64(itemsTotal), 10)
	return "Math.max(1, Math.ceil(" + itemsTotalText + " / (" +
		itemsPerPagePath + " || 1)))"
}

func paginationReadoutExpressionBuilder(
	pageNumberPath, itemsPerPagePath string, itemsTotal uint,
) string {
	itemsTotalText := strconv.FormatUint(uint64(itemsTotal), 10)
	if itemsTotal == 0 {
		return `"0 of 0"`
	}
	guardedItemsPerPage := "(" + itemsPerPagePath + " || 1)"
	firstItemExpression := "(" + pageNumberPath + " - 1) * " +
		guardedItemsPerPage + " + 1"
	lastItemExpression := "Math.min(" + pageNumberPath + " * " +
		guardedItemsPerPage + ", " + itemsTotalText + ")"
	return firstItemExpression + ` + "–" + ` + lastItemExpression +
		` + " of ` + itemsTotalText + `"`
}

package uiStructural

import "strconv"

type ItemsPerPage uint

var paginationDefaultItemsPerPageSizeChoices = []uint{5, 10, 30, 50}

func itemsPerPageSizeChoicesResolver(choices []ItemsPerPage) []uint {
	if len(choices) == 0 {
		return paginationDefaultItemsPerPageSizeChoices
	}
	itemsPerPageSizeChoices := make([]uint, len(choices))
	for index, itemsPerPage := range choices {
		itemsPerPageSizeChoices[index] = uint(itemsPerPage)
	}
	return itemsPerPageSizeChoices
}

func itemsPerPageResolver(itemsPerPage ItemsPerPage, choices []uint) uint {
	if itemsPerPage > 0 {
		return uint(itemsPerPage)
	}
	return choices[0]
}

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

func paginationLastPageNumberExpressionBuilder(
	pagesTotalExpression string,
) string {
	return "Math.max(0, " + pagesTotalExpression + " - 1)"
}

func paginationReadoutExpressionBuilder(
	pageNumberPath, itemsPerPagePath string, itemsTotal uint,
) string {
	itemsTotalText := strconv.FormatUint(uint64(itemsTotal), 10)
	if itemsTotal == 0 {
		return `"0 of 0"`
	}
	guardedItemsPerPage := "(" + itemsPerPagePath + " || 1)"
	firstItemExpression := pageNumberPath + " * " + guardedItemsPerPage + " + 1"
	lastItemExpression := "Math.min((" + pageNumberPath + " + 1) * " +
		guardedItemsPerPage + ", " + itemsTotalText + ")"
	return firstItemExpression + ` + "–" + ` + lastItemExpression +
		` + " of ` + itemsTotalText + `"`
}

package uiStructural

import "strconv"

var paginationDefaultItemsPerPageSizeChoices = []uint{5, 10, 30, 50}

func paginationReadoutExpressionBuilder(
	pageNumberPath, itemsPerPagePath string, itemsTotal uint,
) string {
	itemsTotalText := strconv.FormatUint(uint64(itemsTotal), 10)
	if itemsTotal == 0 {
		return `"0 of 0"`
	}
	firstItemExpression := "(" + pageNumberPath + " - 1) * " +
		itemsPerPagePath + " + 1"
	lastItemExpression := "Math.min(" + pageNumberPath + " * " +
		itemsPerPagePath + ", " + itemsTotalText + ")"
	return firstItemExpression + ` + "–" + ` + lastItemExpression +
		` + " of ` + itemsTotalText + `"`
}

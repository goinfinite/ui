package uiStructural

import "testing"

func TestPaginationReadoutExpression(t *testing.T) {
	testCases := []struct {
		name               string
		pageNumberPath     string
		itemsPerPagePath   string
		itemsTotal         uint
		expectedExpression string
	}{
		{
			name:             "standard page",
			pageNumberPath:   "pageNumber",
			itemsPerPagePath: "itemsPerPage",
			itemsTotal:       240,
			expectedExpression: `(pageNumber - 1) * itemsPerPage + 1 + "–" + ` +
				`Math.min(pageNumber * itemsPerPage, 240) + " of 240"`,
		},
		{
			name:               "empty result set",
			pageNumberPath:     "pageNumber",
			itemsPerPagePath:   "itemsPerPage",
			itemsTotal:         0,
			expectedExpression: `"0 of 0"`,
		},
		{
			name:             "dotted state paths",
			pageNumberPath:   "tableState.pageNumber",
			itemsPerPagePath: "tableState.itemsPerPage",
			itemsTotal:       30,
			expectedExpression: `(tableState.pageNumber - 1) * tableState.itemsPerPage + 1 + "–" + ` +
				`Math.min(tableState.pageNumber * tableState.itemsPerPage, 30) + " of 30"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualExpression := paginationReadoutExpressionBuilder(
				testCase.pageNumberPath, testCase.itemsPerPagePath, testCase.itemsTotal,
			)
			if actualExpression != testCase.expectedExpression {
				t.Errorf(
					"paginationReadoutExpressionBuilder(%q, %q, %d) = %q; want %q",
					testCase.pageNumberPath, testCase.itemsPerPagePath,
					testCase.itemsTotal, actualExpression, testCase.expectedExpression,
				)
			}
		})
	}
}

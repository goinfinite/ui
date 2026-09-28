package uiStructural

import "testing"

func TestPaginationPagesTotalExpressionBuilder(t *testing.T) {
	testCases := []struct {
		name             string
		itemsPerPagePath string
		itemsTotal       uint
		pagesTotal       uint
		expected         string
	}{
		{
			name:             "derives the page count from the item total",
			itemsPerPagePath: "itemsPerPage",
			itemsTotal:       240,
			pagesTotal:       24,
			expected:         "Math.max(1, Math.ceil(240 / (itemsPerPage || 1)))",
		},
		{
			name:             "uses the fallback when the item total is unknown",
			itemsPerPagePath: "itemsPerPage",
			itemsTotal:       0,
			pagesTotal:       7,
			expected:         "7",
		},
		{
			name:             "dotted state paths",
			itemsPerPagePath: "tableState.itemsPerPage",
			itemsTotal:       30,
			pagesTotal:       3,
			expected:         "Math.max(1, Math.ceil(30 / (tableState.itemsPerPage || 1)))",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := paginationPagesTotalExpressionBuilder(
				testCase.itemsPerPagePath, testCase.itemsTotal, testCase.pagesTotal,
			)
			if actual != testCase.expected {
				t.Errorf(
					"pagesTotalExpressionBuilder(%q, %d, %d) = %q; want %q",
					testCase.itemsPerPagePath, testCase.itemsTotal,
					testCase.pagesTotal, actual, testCase.expected,
				)
			}
		})
	}
}

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
			expectedExpression: `(pageNumber - 1) * (itemsPerPage || 1) + 1 + "–" + ` +
				`Math.min(pageNumber * (itemsPerPage || 1), 240) + " of 240"`,
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
			expectedExpression: `(tableState.pageNumber - 1) * (tableState.itemsPerPage || 1) + 1 + "–" + ` +
				`Math.min(tableState.pageNumber * (tableState.itemsPerPage || 1), 30) + " of 30"`,
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

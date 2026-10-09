package uiForm

import "strconv"

const (
	selectInputOptionRowHeightRem           = 2.5
	searchableSelectInputOptionRowHeightRem = 2.25
	selectDropdownBorderHeightRem           = 0.125
)

func selectDropdownHeightClassesResolver(
	minHeightClass, maxHeightClass string,
) string {
	maxHeight := "max-h-60"
	if maxHeightClass != "" {
		maxHeight = maxHeightClass
	}
	if minHeightClass == "" {
		return maxHeight
	}
	return minHeightClass + " " + maxHeight
}

func selectDropdownMaxVisibleOptionsStyle(
	maxVisibleOptions uint, optionRowHeightRem float64,
) string {
	if maxVisibleOptions == 0 {
		return ""
	}
	maxHeightRem := float64(maxVisibleOptions)*optionRowHeightRem +
		selectDropdownBorderHeightRem
	return "max-height: " + strconv.FormatFloat(maxHeightRem, 'f', -1, 64) + "rem"
}

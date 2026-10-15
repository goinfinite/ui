package uiForm

import (
	"strconv"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

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

func selectDropdownSurfaceClassesResolver(
	dropdownBackgroundColor, dropdownBorderColor, dropdownBorderRadius,
	dropdownMinHeightClass, dropdownMaxHeightClass string,
) string {
	dropdownBackgroundClass := uiToolset.BackgroundColorClassResolver(
		dropdownBackgroundColor, "bg-neutral-800/95",
	)
	dropdownBorderClass := uiToolset.BorderColorClassResolver(
		dropdownBorderColor, "border-neutral-50/5",
	)
	dropdownBorderRadiusClass := uiToolset.BorderRadiusClassResolver(
		dropdownBorderRadius, "rounded-md",
	)
	dropdownHeightClasses := selectDropdownHeightClassesResolver(
		dropdownMinHeightClass, dropdownMaxHeightClass,
	)
	return "absolute left-0 z-20 w-full overflow-auto " +
		dropdownBorderRadiusClass + " border " + dropdownBorderClass +
		" shadow-lg hover:border-neutral-50/30 " +
		dropdownHeightClasses + " " + dropdownBackgroundClass
}

func selectLabelValueOptionValuesResolver(
	options []SelectLabelValueOption,
) []string {
	values := make([]string, 0, len(options))
	for _, option := range options {
		values = append(values, option.Value)
	}
	return values
}

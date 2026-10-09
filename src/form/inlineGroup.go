package uiForm

import "strconv"

const (
	inlineGroupOrientationHorizontal = "horizontal"
	inlineGroupOrientationVertical   = "vertical"

	inlineGroupOptionRowHeightRem = 1.75
)

func inlineGroupContentClassesResolver(
	orientation string, maxVisibleOptions uint,
) string {
	contentClasses := "flex min-h-8 flex-row items-center gap-2 p-1.5"
	if orientation == inlineGroupOrientationVertical {
		contentClasses = "flex min-h-8 flex-col items-start gap-2 p-1.5"
	}
	if maxVisibleOptions > 0 {
		contentClasses += " overflow-y-auto"
	}
	return contentClasses
}

func inlineGroupMaxVisibleOptionsStyle(maxVisibleOptions uint) string {
	if maxVisibleOptions == 0 {
		return ""
	}
	maxHeightRem := float64(maxVisibleOptions) * inlineGroupOptionRowHeightRem
	return "max-height: " + strconv.FormatFloat(maxHeightRem, 'f', -1, 64) + "rem"
}

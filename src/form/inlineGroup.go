package uiForm

import (
	"encoding/json"
	"log/slog"
	"strconv"
)

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

func inlineCheckboxGroupStatePathResolver(inputSettings []CheckboxInputSettings) string {
	if len(inputSettings) == 0 {
		return ""
	}
	return inputSettings[0].TwoWayStatePath
}

func inlineCheckboxGroupValuesArrayLiteralBuilder(inputSettings []CheckboxInputSettings) string {
	values := make([]string, 0, len(inputSettings))
	seenValues := make(map[string]bool, len(inputSettings))
	for _, settings := range inputSettings {
		if settings.IsDisabled {
			continue
		}
		value := settings.Value
		if value == "" {
			value = checkboxInputDefaultValue
		}
		if seenValues[value] {
			continue
		}
		seenValues[value] = true
		values = append(values, value)
	}
	encodedValues, marshalErr := json.Marshal(values)
	if marshalErr != nil {
		slog.Debug(
			"InlineCheckboxGroupValuesMarshalFailure",
			slog.String("err", marshalErr.Error()),
		)
		return "[]"
	}
	return string(encodedValues)
}

func inlineCheckboxGroupOptionRowClassesResolver(
	orientation string, hasFloatingButtons bool,
) string {
	rowWidthClass := "w-fit"
	if orientation == inlineGroupOrientationVertical {
		rowWidthClass = "w-[calc(100%+0.75rem)]"
		if hasFloatingButtons {
			rowWidthClass = "w-[calc(100%-1.5rem)]"
		}
	}
	return rowWidthClass + " -m-1.5 rounded p-1.5 transition-colors hover:bg-neutral-50/10"
}

package uiDisplay

import "strconv"

const (
	confirmationPresetConfirm  string = "confirm"
	confirmationPresetWarning  string = "warning"
	confirmationPresetCritical string = "critical"
	confirmationPresetDelete   string = "delete"
)

type confirmationTone struct {
	ConfirmButtonIcon                 string
	ConfirmButtonLabel                string
	ConfirmButtonBackgroundColor      string
	ConfirmButtonBackgroundColorHover string
	HeaderIcon                        string
	HeaderIconBackgroundColor         string
	HeaderIconColor                   string
	HeaderTitle                       string
	Note                              string
	PanelBackgroundColor              string
	QuestionAlone                     string
	QuestionWithTarget                string
	TargetChipBorderColor             string
}

func confirmationToneResolver(preset string) confirmationTone {
	tones := map[string]confirmationTone{
		confirmationPresetConfirm: {
			ConfirmButtonIcon:                 "ph-check",
			ConfirmButtonLabel:                "Yes, Proceed!",
			ConfirmButtonBackgroundColor:      "primary-500",
			ConfirmButtonBackgroundColorHover: "primary-400",
			HeaderIcon:                        "ph-question",
			HeaderIconBackgroundColor:         "neutral-50/20",
			HeaderIconColor:                   "neutral-50",
			HeaderTitle:                       "Confirm Action",
			Note:                              "This action needs your confirmation.",
			PanelBackgroundColor:              "primary-900",
			QuestionAlone:                     "Are you sure you want to proceed?",
			QuestionWithTarget:                "Are you sure you want to proceed with",
			TargetChipBorderColor:             "neutral-50/40",
		},
		confirmationPresetWarning: {
			ConfirmButtonIcon:                 "ph-check",
			ConfirmButtonLabel:                "Yes, Proceed!",
			ConfirmButtonBackgroundColor:      "amber-600",
			ConfirmButtonBackgroundColorHover: "amber-500",
			HeaderIcon:                        "ph-warning-circle",
			HeaderIconBackgroundColor:         "amber-600",
			HeaderIconColor:                   "neutral-50",
			HeaderTitle:                       "Warning",
			Note:                              "Review the effects before you continue.",
			PanelBackgroundColor:              "yellow-900",
			QuestionAlone:                     "Are you sure you want to proceed?",
			QuestionWithTarget:                "Are you sure you want to proceed with",
			TargetChipBorderColor:             "amber-500",
		},
		confirmationPresetCritical: {
			ConfirmButtonIcon:                 "ph-warning",
			ConfirmButtonLabel:                "Yes, Proceed!",
			ConfirmButtonBackgroundColor:      "red-800",
			ConfirmButtonBackgroundColorHover: "red-900",
			HeaderIcon:                        "ph-warning",
			HeaderIconBackgroundColor:         "neutral-300",
			HeaderIconColor:                   "red-950",
			HeaderTitle:                       "Critical Action",
			Note:                              "This action cannot be undone.",
			PanelBackgroundColor:              "red-950",
			QuestionAlone:                     "Are you sure you want to proceed?",
			QuestionWithTarget:                "Are you sure you want to proceed with",
			TargetChipBorderColor:             "red-900",
		},
		confirmationPresetDelete: {
			ConfirmButtonIcon:                 "ph-trash",
			ConfirmButtonLabel:                "Yes, Delete!",
			ConfirmButtonBackgroundColor:      "red-800",
			ConfirmButtonBackgroundColorHover: "red-900",
			HeaderIcon:                        "ph-trash",
			HeaderIconBackgroundColor:         "neutral-300",
			HeaderIconColor:                   "red-950",
			HeaderTitle:                       "Delete",
			Note:                              "This process cannot be undone. This will permanently delete the element(s).",
			PanelBackgroundColor:              "red-950",
			QuestionAlone:                     "Are you sure you want to delete?",
			QuestionWithTarget:                "Are you sure you want to delete",
			TargetChipBorderColor:             "red-900",
		},
	}
	tone, isKnownPreset := tones[preset]
	if !isKnownPreset {
		return tones[confirmationPresetConfirm]
	}
	return tone
}

func confirmationTypeToConfirmMatchStatePathResolver(
	targetNameStatePath, targetIdStatePath string,
) string {
	if targetNameStatePath != "" {
		return targetNameStatePath
	}
	return targetIdStatePath
}

func confirmationTypeToConfirmFieldLabelResolver(
	targetNameStatePath, targetIdStatePath, expectedValue string,
) string {
	if targetNameStatePath != "" {
		return "Type the name to confirm"
	}
	if targetIdStatePath != "" {
		return "Type the id to confirm"
	}
	if expectedValue != "" {
		return "Type " + strconv.Quote(expectedValue) + " to confirm"
	}
	return "Type to confirm"
}

func confirmationTypeToConfirmDisabledExpressionBuilder(
	typedConfirmationStatePath, matchStatePath, expectedValue string,
) string {
	if matchStatePath == "" && expectedValue == "" {
		return "true"
	}
	typedExpression := "String(" + typedConfirmationStatePath + " ?? '').trim()"
	if matchStatePath == "" {
		return typedExpression + " !== " + strconv.Quote(expectedValue)
	}
	matchExpression := "String(" + matchStatePath + " ?? '').trim()"
	emptyMatchExpression := matchExpression + " === ''"
	mismatchExpression := typedExpression + " !== " + matchExpression
	return emptyMatchExpression + " || " + mismatchExpression
}

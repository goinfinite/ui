package uiDisplay

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestConfirmationDialogRetainsSmallInitialSize(t *testing.T) {
	settings := ConfirmationDialogSettings{
		Size:                     ConfirmationDialogSizeXs,
		IsVisibleTwoWayStatePath: "isDialogVisible",
	}

	var rendered bytes.Buffer
	err := confirmationDialog(
		settings, confirmationToneResolver(confirmationPresetConfirm),
	).Render(context.Background(), &rendered)
	if err != nil {
		t.Fatalf("DialogRenderFailed: %v", err)
	}

	if !strings.Contains(rendered.String(), "w-[40%]") {
		t.Errorf("RenderedWidthClassMissing: %q", "w-[40%]")
	}
}

func TestConfirmationDialogDefaultsToSmallSize(t *testing.T) {
	settings := ConfirmationDialogSettings{
		IsVisibleTwoWayStatePath: "isDialogVisible",
	}

	var rendered bytes.Buffer
	err := confirmationDialog(
		settings, confirmationToneResolver(confirmationPresetConfirm),
	).Render(context.Background(), &rendered)
	if err != nil {
		t.Fatalf("DialogRenderFailed: %v", err)
	}

	if !strings.Contains(rendered.String(), "w-[50%]") {
		t.Errorf("RenderedWidthClassMissing: %q", "w-[50%]")
	}
}

func TestConfirmationToneResolver(t *testing.T) {
	tests := []struct {
		name                 string
		preset               string
		expectedTitle        string
		expectedIcon         string
		expectedConfirmLabel string
	}{
		{
			name:                 "ConfirmPresetReturnsConfirmTone",
			preset:               confirmationPresetConfirm,
			expectedTitle:        "Confirm Action",
			expectedIcon:         "ph-question",
			expectedConfirmLabel: "Yes, Proceed!",
		},
		{
			name:                 "WarningPresetReturnsWarningTone",
			preset:               confirmationPresetWarning,
			expectedTitle:        "Warning",
			expectedIcon:         "ph-warning-circle",
			expectedConfirmLabel: "Yes, Proceed!",
		},
		{
			name:                 "CriticalPresetReturnsCriticalTone",
			preset:               confirmationPresetCritical,
			expectedTitle:        "Critical Action",
			expectedIcon:         "ph-warning",
			expectedConfirmLabel: "Yes, Proceed!",
		},
		{
			name:                 "DeletePresetReturnsDeleteTone",
			preset:               confirmationPresetDelete,
			expectedTitle:        "Delete",
			expectedIcon:         "ph-trash",
			expectedConfirmLabel: "Yes, Delete!",
		},
		{
			name:                 "UnknownPresetFallsBackToConfirmTone",
			preset:               "unknown",
			expectedTitle:        "Confirm Action",
			expectedIcon:         "ph-question",
			expectedConfirmLabel: "Yes, Proceed!",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tone := confirmationToneResolver(testCase.preset)

			if tone.HeaderTitle != testCase.expectedTitle {
				t.Errorf("HeaderTitleMismatch: got %q, want %q", tone.HeaderTitle, testCase.expectedTitle)
			}
			if tone.HeaderIcon != testCase.expectedIcon {
				t.Errorf("HeaderIconMismatch: got %q, want %q", tone.HeaderIcon, testCase.expectedIcon)
			}
			if tone.ConfirmButtonLabel != testCase.expectedConfirmLabel {
				t.Errorf(
					"ConfirmButtonLabelMismatch: got %q, want %q",
					tone.ConfirmButtonLabel, testCase.expectedConfirmLabel,
				)
			}
		})
	}
}

func TestConfirmationTypeToConfirmMatchStatePathResolver(t *testing.T) {
	tests := []struct {
		name                string
		targetNameStatePath string
		targetIdStatePath   string
		expected            string
	}{
		{
			name:                "NamePathWinsOverIdPath",
			targetNameStatePath: "record.name",
			targetIdStatePath:   "record.id",
			expected:            "record.name",
		},
		{
			name:                "IdPathUsedWhenNamePathIsMissing",
			targetNameStatePath: "",
			targetIdStatePath:   "record.id",
			expected:            "record.id",
		},
		{
			name:                "EmptyWhenNoTargetPaths",
			targetNameStatePath: "",
			targetIdStatePath:   "",
			expected:            "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := confirmationTypeToConfirmMatchStatePathResolver(
				testCase.targetNameStatePath, testCase.targetIdStatePath,
			)

			if resolved != testCase.expected {
				t.Errorf("MatchStatePathMismatch: got %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestConfirmationTypeToConfirmFieldLabelResolver(t *testing.T) {
	tests := []struct {
		name                string
		targetNameStatePath string
		targetIdStatePath   string
		expectedValue       string
		expected            string
	}{
		{
			name:                "NameLabelWins",
			targetNameStatePath: "record.name",
			targetIdStatePath:   "record.id",
			expected:            "Type the name to confirm",
		},
		{
			name:              "IdLabelUsedWhenNamePathIsMissing",
			targetIdStatePath: "record.id",
			expected:          "Type the id to confirm",
		},
		{
			name:          "ExpectedValueLabelWhenNoTargetPaths",
			expectedValue: "proceed",
			expected:      `Type "proceed" to confirm`,
		},
		{
			name:     "GenericLabelWhenNothingIsSet",
			expected: "Type to confirm",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := confirmationTypeToConfirmFieldLabelResolver(
				testCase.targetNameStatePath, testCase.targetIdStatePath,
				testCase.expectedValue,
			)

			if resolved != testCase.expected {
				t.Errorf("FieldLabelMismatch: got %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestConfirmationTypeToConfirmDisabledExpressionBuilder(t *testing.T) {
	tests := []struct {
		name           string
		typedStatePath string
		matchStatePath string
		expectedValue  string
		expected       string
	}{
		{
			name:           "MatchStatePathComparison",
			typedStatePath: "typedConfirmationValue",
			matchStatePath: "record.name",
			expected: "String(record.name ?? '').trim() === '' || " +
				"String(typedConfirmationValue ?? '').trim() !== String(record.name ?? '').trim()",
		},
		{
			name:           "ExpectedValueComparisonWhenNoMatchStatePath",
			typedStatePath: "typedConfirmationValue",
			expectedValue:  "proceed",
			expected:       `String(typedConfirmationValue ?? '').trim() !== "proceed"`,
		},
		{
			name:           "PermanentlyDisabledWhenNoMatchSource",
			typedStatePath: "typedConfirmationValue",
			expected:       "true",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			disabledExpression := confirmationTypeToConfirmDisabledExpressionBuilder(
				testCase.typedStatePath, testCase.matchStatePath, testCase.expectedValue,
			)

			if disabledExpression != testCase.expected {
				t.Errorf(
					"DisabledExpressionMismatch: got %q, want %q",
					disabledExpression, testCase.expected,
				)
			}
		})
	}
}

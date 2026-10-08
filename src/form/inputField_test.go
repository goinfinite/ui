package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestInputFieldRendersHintTooltipSizeClasses(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := InputField(InputFieldSettings{
		InputType:                        InputTypeText,
		InputName:                        "name",
		Label:                            "Name",
		HintValue:                        "This is a helpful hint.",
		HintDisplay:                      InputHintDisplayTooltip,
		HintDisplayTooltipMinWidthClass:  "min-w-40",
		HintDisplayTooltipMaxWidthClass:  "max-w-80",
		HintDisplayTooltipMinHeightClass: "min-h-12",
		HintDisplayTooltipMaxHeightClass: "max-h-64",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("InputFieldRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedClass := range []string{"min-w-40", "max-w-80", "min-h-12", "max-h-64"} {
		if !strings.Contains(renderedHtml, expectedClass) {
			t.Errorf("RenderedHtmlMissingTooltipClass(%q)", expectedClass)
		}
	}
}

func TestInputFieldSizeTokensResolver(t *testing.T) {
	testCases := []struct {
		name           string
		size           string
		expectedTokens inputFieldSizeTokens
	}{
		{
			name: "xs",
			size: InputFieldSizeXs,
			expectedTokens: inputFieldSizeTokens{
				AffixPadding:   "p-1",
				InputPadding:   "p-1 pt-0",
				LeftPadding:    "pl-1",
				LabelOffset:    "0.5",
				LegendTextSize: "text-[0.625rem]",
				TextSize:       "text-xs",
			},
		},
		{
			name: "sm",
			size: InputFieldSizeSm,
			expectedTokens: inputFieldSizeTokens{
				AffixPadding:   "p-1.5",
				InputPadding:   "p-1.5 pt-0",
				LeftPadding:    "pl-1.5",
				LabelOffset:    "1",
				LegendTextSize: "text-[0.625rem]",
				TextSize:       "text-xs",
			},
		},
		{
			name: "lg",
			size: InputFieldSizeLg,
			expectedTokens: inputFieldSizeTokens{
				AffixPadding:   "p-2",
				InputPadding:   "p-2 pt-0",
				LeftPadding:    "pl-2",
				LabelOffset:    "2",
				LegendTextSize: "text-sm",
				TextSize:       "text-base",
			},
		},
		{
			name: "xl",
			size: InputFieldSizeXl,
			expectedTokens: inputFieldSizeTokens{
				AffixPadding:   "p-2.5",
				InputPadding:   "p-2.5 pt-0",
				LeftPadding:    "pl-2.5",
				LabelOffset:    "2.5",
				LegendTextSize: "text-base",
				TextSize:       "text-lg",
			},
		},
		{
			name: "default",
			size: "",
			expectedTokens: inputFieldSizeTokens{
				AffixPadding:   "p-1.5",
				InputPadding:   "p-1.5 pt-0",
				LeftPadding:    "pl-1.5",
				LabelOffset:    "1.5",
				LegendTextSize: "text-xs",
				TextSize:       "text-sm",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualTokens := inputFieldSizeTokensResolver(testCase.size)
			if actualTokens != testCase.expectedTokens {
				t.Errorf(
					"SizeTokensMismatch(%q): got %+v, want %+v",
					testCase.size, actualTokens, testCase.expectedTokens,
				)
			}
		})
	}
}

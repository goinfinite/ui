package uiForm

import "testing"

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

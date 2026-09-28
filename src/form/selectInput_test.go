package uiForm

import (
	"strings"
	"testing"
)

func TestSelectInputSizeClassesResolver(t *testing.T) {
	testCases := []struct {
		name                      string
		size                      string
		expectedLegendClass       string
		expectedTriggerClass      string
		expectedEmptyPaddingClass string
	}{
		{
			name:                      "xs",
			size:                      SelectInputSizeXs,
			expectedLegendClass:       "text-[0.625rem]",
			expectedTriggerClass:      "px-1 pb-0.5 text-xs",
			expectedEmptyPaddingClass: "pt-0.5",
		},
		{
			name:                      "sm",
			size:                      SelectInputSizeSm,
			expectedLegendClass:       "text-[0.625rem]",
			expectedTriggerClass:      "px-1.5 pb-1 text-xs",
			expectedEmptyPaddingClass: "pt-1",
		},
		{
			name:                      "md",
			size:                      SelectInputSizeMd,
			expectedLegendClass:       "text-xs",
			expectedTriggerClass:      "px-1.5 pb-1.5 text-sm",
			expectedEmptyPaddingClass: "pt-1.5",
		},
		{
			name:                      "lg",
			size:                      SelectInputSizeLg,
			expectedLegendClass:       "text-sm",
			expectedTriggerClass:      "px-2 pb-2 text-base",
			expectedEmptyPaddingClass: "pt-2",
		},
		{
			name:                      "xl",
			size:                      SelectInputSizeXl,
			expectedLegendClass:       "text-base",
			expectedTriggerClass:      "px-2.5 pb-2.5 text-lg",
			expectedEmptyPaddingClass: "pt-2.5",
		},
		{
			name:                      "default",
			size:                      "",
			expectedLegendClass:       "text-xs",
			expectedTriggerClass:      "px-1.5 pb-1.5 text-sm",
			expectedEmptyPaddingClass: "pt-1.5",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			legendClass, triggerClass, emptyPaddingClass := selectInputSizeClassesResolver(testCase.size)
			actualClasses := []string{legendClass, triggerClass, emptyPaddingClass}
			expectedClasses := []string{
				testCase.expectedLegendClass, testCase.expectedTriggerClass, testCase.expectedEmptyPaddingClass,
			}
			if strings.Join(actualClasses, "|") != strings.Join(expectedClasses, "|") {
				t.Errorf(
					"SizeClassesMismatch(%q): got %q, want %q",
					testCase.size, actualClasses, expectedClasses,
				)
			}
		})
	}
}

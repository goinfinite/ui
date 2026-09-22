package uiForm

import (
	"strings"
	"testing"
)

func TestCheckboxInputSizeClassesResolver(t *testing.T) {
	testCases := []struct {
		name                string
		size                string
		expectedBoxClasses  string
		expectedIconClasses string
	}{
		{name: "xs", size: CheckboxInputSizeXs, expectedBoxClasses: "h-3.5 w-3.5", expectedIconClasses: "text-[0.5rem]"},
		{name: "sm", size: CheckboxInputSizeSm, expectedBoxClasses: "h-4 w-4", expectedIconClasses: "text-[0.625rem]"},
		{name: "md", size: CheckboxInputSizeMd, expectedBoxClasses: "h-5 w-5", expectedIconClasses: "text-xs"},
		{name: "lg", size: CheckboxInputSizeLg, expectedBoxClasses: "h-6 w-6", expectedIconClasses: "text-sm"},
		{name: "xl", size: CheckboxInputSizeXl, expectedBoxClasses: "h-7 w-7", expectedIconClasses: "text-base"},
		{name: "default", size: "", expectedBoxClasses: "h-5 w-5", expectedIconClasses: "text-xs"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			boxClasses, iconClasses, _ := checkboxInputSizeClassesResolver(testCase.size)
			actualClasses := []string{boxClasses, iconClasses}
			expectedClasses := []string{testCase.expectedBoxClasses, testCase.expectedIconClasses}
			if strings.Join(actualClasses, "|") != strings.Join(expectedClasses, "|") {
				t.Errorf(
					"checkboxInputSizeClassesResolver(%q) = %q; want %q",
					testCase.size, actualClasses, expectedClasses,
				)
			}
		})
	}
}

func TestCheckboxInputShapeClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		shape             string
		expectedClassName string
	}{
		{name: "square", shape: CheckboxInputShapeSquare, expectedClassName: "rounded-none"},
		{name: "rounded", shape: CheckboxInputShapeRounded, expectedClassName: "rounded"},
		{name: "circular", shape: CheckboxInputShapeCircular, expectedClassName: "rounded-full"},
		{name: "default", shape: "", expectedClassName: "rounded"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := checkboxInputShapeClassResolver(testCase.shape)
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"checkboxInputShapeClassResolver(%q) = %q; want %q",
					testCase.shape, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

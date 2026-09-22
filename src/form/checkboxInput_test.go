package uiForm

import (
	"bytes"
	"context"
	"regexp"
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
					"SizeClassesMismatch(%q): got %q, want %q",
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
					"ShapeClassMismatch(%q): got %q, want %q",
					testCase.shape, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestCheckboxInputBoxBorderClassesResolver(t *testing.T) {
	testCases := []struct {
		name            string
		uncheckedColor  string
		checkedColor    string
		errorColor      string
		isInvalid       bool
		expectedClasses string
	}{
		{
			name:           "valid uses unchecked and checked colors",
			uncheckedColor: "neutral-50/20", checkedColor: "secondary-500",
			errorColor:      "red-500",
			isInvalid:       false,
			expectedClasses: "border-neutral-50/20 peer-checked:border-secondary-500",
		},
		{
			name:           "invalid uses the error color on both borders",
			uncheckedColor: "neutral-50/20", checkedColor: "secondary-500",
			errorColor:      "red-500",
			isInvalid:       true,
			expectedClasses: "border-red-500 peer-checked:border-red-500",
		},
		{
			name:           "invalid honors a custom error color",
			uncheckedColor: "neutral-50/20", checkedColor: "secondary-500",
			errorColor:      "amber-500",
			isInvalid:       true,
			expectedClasses: "border-amber-500 peer-checked:border-amber-500",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClasses := checkboxInputBoxBorderClassesResolver(
				testCase.uncheckedColor, testCase.checkedColor,
				testCase.errorColor, testCase.isInvalid,
			)
			if actualClasses != testCase.expectedClasses {
				t.Errorf(
					"BoxBorderClassesMismatch: got %q, want %q",
					actualClasses, testCase.expectedClasses,
				)
			}
		})
	}
}

func TestCheckboxInputRendersStaticErrorState(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CheckboxInput(CheckboxInputSettings{
		InputId:      "terms",
		Label:        "Accept the terms",
		IsInvalid:    true,
		ErrorMessage: "You must accept the terms.",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CheckboxRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"aria-invalid",
		"border-red-500 peer-checked:border-red-500",
		`aria-describedby="terms-error"`,
		`id="terms-error"`,
		"You must accept the terms.",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCheckboxInputRendersBoundErrorState(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CheckboxInput(CheckboxInputSettings{
		Label:                       "Accept the terms",
		IsInvalidOneWayStatePath:    "termsError",
		ErrorMessageOneWayStatePath: "termsErrorMessage",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CheckboxRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`:aria-invalid="termsError"`,
		`!border-red-500 !peer-checked:border-red-500`,
		`x-show="termsError"`,
		`x-text="termsErrorMessage"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCheckboxInputCombinesIndeterminateAndInvalidClassExpressions(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CheckboxInput(CheckboxInputSettings{
		Label:                        "Select all",
		IndeterminateOneWayStatePath: "someSelected",
		IsInvalidOneWayStatePath:     "termsError",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CheckboxRenderFailed: %v", renderErr)
	}
	classExprMatch := regexp.MustCompile(`:class="([^"]*)"`).FindStringSubmatch(buffer.String())
	if classExprMatch == nil {
		t.Fatal("RenderedHtmlMissingDynamicClassExpression")
	}
	classExpr := strings.ReplaceAll(classExprMatch[1], "&#39;", "'")
	expectedFragments := []string{
		"!border-secondary-500 !text-secondary-500",
		"!border-red-500 !peer-checked:border-red-500",
		"someSelected",
		"termsError",
		"+ ' ' +",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(classExpr, expectedFragment) {
			t.Errorf("DynamicClassExpressionMissingFragment: expression %q lacks %q", classExpr, expectedFragment)
		}
	}
}

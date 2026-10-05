package uiStructural

import "testing"

func TestHorizontalAlignmentJustifyClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		alignment         HorizontalAlignment
		expectedClassName string
	}{
		{name: "left", alignment: HorizontalAlignmentLeft, expectedClassName: "justify-start"},
		{name: "center", alignment: HorizontalAlignmentCenter, expectedClassName: "justify-center"},
		{name: "right", alignment: HorizontalAlignmentRight, expectedClassName: "justify-end"},
		{name: "default", alignment: "", expectedClassName: "justify-start"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := testCase.alignment.justifyClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"JustifyClassMismatch(%q): got %q, want %q",
					testCase.alignment, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestTextAlignmentClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		alignment         TextAlignment
		expectedClassName string
	}{
		{name: "left", alignment: TextAlignmentLeft, expectedClassName: "text-left"},
		{name: "center", alignment: TextAlignmentCenter, expectedClassName: "text-center"},
		{name: "right", alignment: TextAlignmentRight, expectedClassName: "text-right"},
		{name: "default", alignment: "", expectedClassName: "text-left"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := testCase.alignment.alignmentClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"AlignmentClassMismatch(%q): got %q, want %q",
					testCase.alignment, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

func TestTextAlignmentJustifyClassResolver(t *testing.T) {
	testCases := []struct {
		name              string
		alignment         TextAlignment
		expectedClassName string
	}{
		{name: "left", alignment: TextAlignmentLeft, expectedClassName: "justify-start"},
		{name: "center", alignment: TextAlignmentCenter, expectedClassName: "justify-center"},
		{name: "right", alignment: TextAlignmentRight, expectedClassName: "justify-end"},
		{name: "default", alignment: "", expectedClassName: "justify-start"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClassName := testCase.alignment.justifyClassResolver()
			if actualClassName != testCase.expectedClassName {
				t.Errorf(
					"JustifyClassMismatch(%q): got %q, want %q",
					testCase.alignment, actualClassName, testCase.expectedClassName,
				)
			}
		})
	}
}

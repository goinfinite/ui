package uiToolset

import "testing"

func TestTextCaseClassResolver(t *testing.T) {
	testCases := []struct {
		name     string
		textCase string
		expected string
	}{
		{name: "none", textCase: TextCaseNone, expected: ""},
		{name: "lower", textCase: TextCaseLower, expected: "lowercase"},
		{name: "upper", textCase: TextCaseUpper, expected: "uppercase"},
		{name: "capitalize", textCase: TextCaseCapitalize, expected: "capitalize"},
		{name: "default", textCase: "", expected: ""},
		{name: "unknown", textCase: "titlecase", expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := TextCaseClassResolver(testCase.textCase)
			if actual != testCase.expected {
				t.Errorf(
					"TextCaseClassMismatch(%q): got %q, want %q",
					testCase.textCase, actual, testCase.expected,
				)
			}
		})
	}
}

package uiToolset

import (
	"strings"
	"testing"
)

func TestMinifierJavaScript(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "simple variable declaration",
			input: "const hello = 'world';\n\nfunction test() {\n  console.log(hello);\n}\n",
		},
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "comments and whitespace",
			input: "// This is a comment\nconst x = 1;  // Another comment\n\n/* Multi-line\n   comment */\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Minifier(MinifierSettings{ContentType: MinifierContentTypeJavaScript, UnminifiedContent: &testCase.input})

			if len(*result) > len(testCase.input) && testCase.input != "" {
				t.Errorf("MinifiedResultShouldBeSmallerThanInput, Got Input: %d Chars, Result: %d Chars", len(testCase.input), len(*result))
			}

			if testCase.input != "" && len(strings.TrimSpace(*result)) == 0 {
				t.Errorf("EmptyResultForNonEmptyInput")
			}

			if testCase.input == "" && len(*result) > 1 {
				t.Errorf("InvalidResultForEmptyInput: %q", *result)
			}
		})
	}
}

func TestMinifierCSS(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "simple css rule",
			input: "body {\n  margin: 0;\n  padding: 10px;\n}\n\n.class {\n  color: red;\n}\n",
		},
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "comments and whitespace",
			input: "/* Header styles */\n.header {\n  display: flex;  /* flex display */\n  gap: 10px;\n}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Minifier(MinifierSettings{ContentType: MinifierContentTypeCss, UnminifiedContent: &testCase.input})

			if len(*result) > len(testCase.input) && testCase.input != "" {
				t.Errorf("MinifiedResultShouldBeSmallerThanInput, Got Input: %d Chars, Result: %d Chars", len(testCase.input), len(*result))
			}

			if testCase.input != "" && len(strings.TrimSpace(*result)) == 0 {
				t.Errorf("EmptyResultForNonEmptyInput")
			}

			if testCase.input == "" && len(*result) > 1 {
				t.Errorf("InvalidResultForEmptyInput: %q", *result)
			}

			if strings.Contains(testCase.input, "/*") && strings.Contains(*result, "/*") {
				t.Errorf("CommentsNotRemovedInOutput")
			}
		})
	}
}

func TestMinifierUnsupportedContentType(t *testing.T) {
	testCases := []struct {
		name        string
		contentType string
		input       string
	}{
		{
			name:        "html content",
			contentType: "text/html",
			input:       "<div>Hello World</div>",
		},
		{
			name:        "plain text",
			contentType: "text/plain",
			input:       "Hello World",
		},
		{
			name:        "empty content type",
			contentType: "",
			input:       "Some content",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Minifier(MinifierSettings{ContentType: testCase.contentType, UnminifiedContent: &testCase.input})
			if *result != testCase.input {
				t.Errorf("OriginalContentNotPreserved, Expected: %q, Got: %q", testCase.input, *result)
			}
		})
	}
}

func TestMinifierWithInvalidInput(t *testing.T) {
	testCases := []struct {
		name        string
		contentType string
		input       string
	}{
		{
			name:        "invalid javascript",
			contentType: MinifierContentTypeJavaScript,
			input:       "const x = ;",
		},
		{
			name:        "invalid css",
			contentType: MinifierContentTypeCss,
			input:       "body { color: }",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			originalInput := testCase.input
			result := Minifier(MinifierSettings{ContentType: testCase.contentType, UnminifiedContent: &testCase.input})

			if *result == "" {
				t.Errorf("EmptyResultForInvalidInput: %q", originalInput)
			}

			if testCase.contentType == MinifierContentTypeJavaScript {
				if !strings.Contains(*result, "const x") {
					t.Errorf("OriginalDeclarationNotPreserved, Got: %q", *result)
				}
			} else if testCase.contentType == MinifierContentTypeCss {
				if !strings.Contains(*result, "body") || !strings.Contains(*result, "color") {
					t.Errorf("OriginalCssElementsNotPreserved, Got: %q", *result)
				}
			}
		})
	}
}

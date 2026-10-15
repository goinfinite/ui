package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestTextAreaRendersPlaceholderStyling(t *testing.T) {
	testCases := []struct {
		name                   string
		settings               TextAreaSettings
		expectedPlaceholderTag string
	}{
		{
			name: "UnboundEmptyHidesPlaceholder",
			settings: TextAreaSettings{
				InputName: "description",
				Label:     "Description",
			},
			expectedPlaceholderTag: "placeholder:opacity-0",
		},
		{
			name: "BoundShowsPlaceholder",
			settings: TextAreaSettings{
				InputName:       "description",
				Label:           "Description",
				TwoWayStatePath: "description",
			},
			expectedPlaceholderTag: "placeholder:opacity-80",
		},
		{
			name: "UnboundValueShowsPlaceholder",
			settings: TextAreaSettings{
				InputName: "description",
				Label:     "Description",
				Value:     "Lorem ipsum.",
			},
			expectedPlaceholderTag: "placeholder:opacity-80",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := TextArea(testCase.settings).Render(context.Background(), &buffer)
			if renderErr != nil {
				t.Fatalf("TextAreaRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			for _, expectedClass := range []string{
				"placeholder:text-neutral-50", testCase.expectedPlaceholderTag,
			} {
				if !strings.Contains(renderedHtml, expectedClass) {
					t.Errorf("RenderedHtmlMissingPlaceholderClass(%q)", expectedClass)
				}
			}
		})
	}
}

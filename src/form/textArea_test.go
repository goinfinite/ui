package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestTextAreaRendersClearButtonByReadOnlyState(t *testing.T) {
	testCases := []struct {
		name            string
		settings        TextAreaSettings
		expectClearIcon bool
	}{
		{
			name: "ReadOnlyBoundFieldHidesClear",
			settings: TextAreaSettings{
				InputName:       "output",
				Label:           "Output",
				TwoWayStatePath: "output",
				IsReadOnly:      true,
			},
			expectClearIcon: false,
		},
		{
			name: "EditableBoundFieldShowsClear",
			settings: TextAreaSettings{
				InputName:       "description",
				Label:           "Description",
				TwoWayStatePath: "description",
			},
			expectClearIcon: true,
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
			hasClearIcon := strings.Contains(renderedHtml, "ph-x-circle")
			if hasClearIcon != testCase.expectClearIcon {
				t.Errorf(
					"ClearButtonMismatch: got %v, want %v",
					hasClearIcon, testCase.expectClearIcon,
				)
			}
		})
	}
}

func TestTextAreaRendersFirstChildMarginReset(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := TextArea(TextAreaSettings{
		InputName: "description",
		Label:     "Description",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("TextAreaRenderFailed: %v", renderErr)
	}
	if !strings.Contains(buffer.String(), "first:mt-0") {
		t.Error("RenderedHtmlMissingFirstChildMarginReset")
	}
}

func TestTextAreaRendersLabelStatePath(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := TextArea(TextAreaSettings{
		InputName:            "description",
		Label:                "Description",
		LabelOneWayStatePath: "fieldLabel",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("TextAreaRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="fieldLabel"`,
		`:aria-label="(fieldLabel) || &#34;Description&#34;"`,
		`:placeholder="(fieldLabel) || &#34;Description&#34;"`,
		"Description",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

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

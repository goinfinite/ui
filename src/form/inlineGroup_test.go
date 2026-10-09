package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestInlineGroupContentClassesResolver(t *testing.T) {
	testCases := []struct {
		name              string
		orientation       string
		maxVisibleOptions uint
		expected          string
	}{
		{
			name:        "default is horizontal",
			orientation: "",
			expected:    "flex min-h-8 flex-row items-center gap-2 p-1.5",
		},
		{
			name:        "vertical",
			orientation: InlineRadioGroupOrientationVertical,
			expected:    "flex min-h-8 flex-col items-start gap-2 p-1.5",
		},
		{
			name:              "max visible options caps the height",
			orientation:       InlineCheckboxGroupOrientationVertical,
			maxVisibleOptions: 4,
			expected:          "flex min-h-8 flex-col items-start gap-2 p-1.5 overflow-y-auto",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualClasses := inlineGroupContentClassesResolver(
				testCase.orientation, testCase.maxVisibleOptions,
			)
			if actualClasses != testCase.expected {
				t.Errorf(
					"ContentClassesMismatch(%q): got %q, want %q",
					testCase.orientation, actualClasses, testCase.expected,
				)
			}
		})
	}
}

func TestInlineGroupMaxVisibleOptionsStyle(t *testing.T) {
	if actualStyle := inlineGroupMaxVisibleOptionsStyle(0); actualStyle != "" {
		t.Errorf("MaxVisibleOptionsStyleNotEmpty: %q", actualStyle)
	}
	if actualStyle := inlineGroupMaxVisibleOptionsStyle(4); actualStyle != "max-height: 7rem" {
		t.Errorf("MaxVisibleOptionsStyleMismatch: %q", actualStyle)
	}
}

func TestInlineRadioGroupRendersOrientationAndMaxVisibleOptions(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := InlineRadioGroup(InlineRadioGroupSettings{
		Label:             "Select an option",
		Orientation:       InlineRadioGroupOrientationVertical,
		MaxVisibleOptions: 3,
		InputSettings: []RadioInputSettings{
			{
				Label:           "Option 1",
				StateValue:      "option1",
				TwoWayStatePath: "groupSelection",
			},
		},
		TwoWayStatePath: "groupSelection",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("InlineRadioGroupRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"flex-col", "overflow-y-auto", "max-height: 5.25rem"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

func TestInlineCheckboxGroupRendersCheckboxesAndOrientation(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := InlineCheckboxGroup(InlineCheckboxGroupSettings{
		Label:             "Fruits",
		Orientation:       InlineCheckboxGroupOrientationVertical,
		MaxVisibleOptions: 2,
		InputSettings: []CheckboxInputSettings{
			{
				Label:           "Apple",
				Value:           "apple",
				TwoWayStatePath: "selectedFruits",
			},
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("InlineCheckboxGroupRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"flex-col", "type=\"checkbox\"", "Apple", "max-height: 3.5rem"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

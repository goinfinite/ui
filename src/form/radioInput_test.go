package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRadioInputRendersDefaultStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := RadioInput(RadioInputSettings{
		Label:           "Option 1",
		StateValue:      "option1",
		TwoWayStatePath: "selectedOption",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("RadioInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"text-neutral-100",
		"hover:text-secondary-500",
		"text-secondary-500 font-bold",
		"peer-focus-visible:ring-secondary-500",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestRadioInputRendersCustomStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := RadioInput(RadioInputSettings{
		Label:           "Option 1",
		StateValue:      "option1",
		TwoWayStatePath: "selectedOption",
		Style: RadioInputStyle{
			ActiveColor:    "violet-400",
			FocusRingColor: "violet-300",
			TextColor:      "violet-100",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("RadioInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"text-violet-100",
		"hover:text-violet-400",
		"text-violet-400 font-bold",
		"peer-focus-visible:ring-violet-300",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "secondary-500") {
		t.Error("CustomStyleRendersDefaultColor")
	}
}

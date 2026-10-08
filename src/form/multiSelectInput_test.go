package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestMultiSelectInputRendersOpaqueDropdownAboveStickyHeaders(t *testing.T) {
	settings := MultiSelectInputSettings{
		InputName:   "status",
		Label:       "Status",
		FlatOptions: []string{"running"},
	}
	var buffer bytes.Buffer
	renderErr := MultiSelectInput(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("MultiSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"z-20", "bg-neutral-800/95"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}

	settings.DropdownBackgroundColor = "emerald-900"
	coloredBuffer := bytes.Buffer{}
	renderErr = MultiSelectInput(settings).Render(context.Background(), &coloredBuffer)
	if renderErr != nil {
		t.Fatalf("MultiSelectInputRenderFailed: %v", renderErr)
	}
	coloredHtml := coloredBuffer.String()
	if !strings.Contains(coloredHtml, "bg-emerald-900") {
		t.Errorf("RenderedHtmlMissingCustomDropdownBackground")
	}
	if strings.Contains(coloredHtml, "bg-neutral-800/95") {
		t.Errorf("RenderedHtmlKeepsDefaultDropdownBackground")
	}
}

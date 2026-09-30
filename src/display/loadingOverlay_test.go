package uiDisplay

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLoadingOverlayHtmxModeRendersHiddenWithHtmxOverride(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := LoadingOverlay(LoadingOverlaySettings{
		Id: "loading-overlay-test",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("LoadingOverlayRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`style="opacity: 0; visibility: hidden`,
		".htmx-request #loading-overlay-test,.htmx-request#loading-overlay-test{opacity:1 !important;visibility:visible !important}",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, ":style=") {
		t.Error("HtmxModeMustNotBindAlpineStyle")
	}
}

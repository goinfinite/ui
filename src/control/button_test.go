package uiControl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestButtonRendersPaddingAndSizeConstraints(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Button(ButtonSettings{
		PaddingSize:              uiToolset.PaddingSizeNone,
		MinWidthClass:            "min-w-7",
		MinHeightClass:           "min-h-7",
		Size:                     ButtonSizeMd,
		IconLeft:                 "ph-binary",
		AriaLabelOneWayStatePath: "isVisible ? 'Hide' : 'Show'",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ButtonRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		"p-0", "min-w-7", "min-h-7",
		`:aria-label="isVisible ? &#39;Hide&#39; : &#39;Show&#39;"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

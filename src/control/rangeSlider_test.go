package uiControl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestThumbRendersUpperOnlyBubbleRingClasses(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Thumb(
		RangeSliderSettings{},
		RangeSliderThumbSettings{
			ValueBubbleEnabled:       true,
			ValueBubbleRingColor:     "red-700",
			ValueBubbleRingThickness: uiToolset.RingThicknessLg,
		},
		true,
	).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ThumbRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if !strings.Contains(renderedHtml, "ring-2 ring-red-700") {
		t.Errorf("RenderedUpperBubbleMissingRingClasses: %s", renderedHtml)
	}
}

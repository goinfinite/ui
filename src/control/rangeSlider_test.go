package uiControl

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestThumbRendersUpperOnlyBubbleRingClasses(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Thumb(RangeSliderSettings{
		ThumbUpperValueBubbleEnabled:       true,
		ThumbUpperValueBubbleRingColor:     "red-700",
		ThumbUpperValueBubbleRingThickness: RangeSliderThumbValueBubbleRingThicknessLg,
	}, true).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ThumbRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if !strings.Contains(renderedHtml, "ring-2 ring-red-700") {
		t.Errorf("RenderedUpperBubbleMissingRingClasses: %s", renderedHtml)
	}
}

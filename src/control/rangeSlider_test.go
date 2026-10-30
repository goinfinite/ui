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

func TestRangeSliderThumbAriaLabelExpression(testingHandle *testing.T) {
	testCases := []struct {
		name          string
		thumbSettings RangeSliderThumbSettings
		isUpper       bool
		expectedValue string
	}{
		{
			name:          "empty path yields no expression",
			thumbSettings: RangeSliderThumbSettings{AriaLabel: "Quantity"},
			expectedValue: "",
		},
		{
			name: "path falls back to the static label",
			thumbSettings: RangeSliderThumbSettings{
				AriaLabel:                "Quantity",
				AriaLabelOneWayStatePath: "quantityLabel",
			},
			expectedValue: `(quantityLabel) || "Quantity"`,
		},
		{
			name: "empty static label falls back to the default name",
			thumbSettings: RangeSliderThumbSettings{
				AriaLabelOneWayStatePath: "quantityLabel",
			},
			expectedValue: `(quantityLabel) || "Value"`,
		},
		{
			name: "upper thumb prefixes the fallback name",
			thumbSettings: RangeSliderThumbSettings{
				Label:                    "Quantity",
				AriaLabelOneWayStatePath: "quantityLabel",
			},
			isUpper:       true,
			expectedValue: `(quantityLabel) || "Upper Quantity"`,
		},
	}
	for _, testCase := range testCases {
		testingHandle.Run(testCase.name, func(subtestHandle *testing.T) {
			actualValue := (RangeSliderSettings{}).resolveThumbAriaLabelExpression(
				testCase.thumbSettings, testCase.isUpper,
			)
			if actualValue != testCase.expectedValue {
				subtestHandle.Errorf(
					"AriaLabelExpressionMismatch: got %q, want %q",
					actualValue, testCase.expectedValue,
				)
			}
		})
	}
}

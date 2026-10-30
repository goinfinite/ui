package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestNumberRangeInputFieldRendersLabelStatePath(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := NumberRangeInputField(NumberRangeInputFieldSettings{
		InputName:            "quantity",
		Label:                "Quantity",
		LabelOneWayStatePath: "quantityLabel",
		TwoWayStatePath:      "quantityValue",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("NumberRangeInputFieldRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="quantityLabel"`,
		`:aria-label="(quantityLabel) || &#34;Quantity&#34;"`,
		`:placeholder="(quantityLabel) || &#34;Quantity&#34;"`,
		"Quantity",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	ariaLabelExpression := `:aria-label="(quantityLabel) || &#34;Quantity&#34;"`
	actualAriaLabelCount := strings.Count(renderedHtml, ariaLabelExpression)
	if actualAriaLabelCount < 2 {
		testingHandle.Errorf(
			"RenderedHtmlAriaLabelCount = %d, want at least 2",
			actualAriaLabelCount,
		)
	}
}

package uiDisplay

import (
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestAccordionEdgeRadiusClassesResolver(t *testing.T) {
	tests := []struct {
		borderRadius string
		expected     string
	}{
		{uiToolset.BorderRadiusNone, ""},
		{uiToolset.BorderRadiusXs, "first:rounded-t-xs last:rounded-b-xs"},
		{uiToolset.BorderRadiusSm, "first:rounded-t-sm last:rounded-b-sm"},
		{uiToolset.BorderRadiusMd, "first:rounded-t-md last:rounded-b-md"},
		{uiToolset.BorderRadiusLg, "first:rounded-t-lg last:rounded-b-lg"},
		{uiToolset.BorderRadiusXl, "first:rounded-t-xl last:rounded-b-xl"},
		{uiToolset.BorderRadius2xl, "first:rounded-t-2xl last:rounded-b-2xl"},
		{uiToolset.BorderRadius3xl, "first:rounded-t-3xl last:rounded-b-3xl"},
		{"", "first:rounded-t-md last:rounded-b-md"},
		{"unknown", "first:rounded-t-md last:rounded-b-md"},
	}
	for _, test := range tests {
		result := accordionEdgeRadiusClassesResolver(test.borderRadius)
		if result != test.expected {
			t.Errorf("edgeRadiusClassesResolver(%q) = %q, expected %q",
				test.borderRadius, result, test.expected)
		}
	}
}

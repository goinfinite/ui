package uiDisplay

import "testing"

func TestAccordionEdgeRadiusClassesResolver(t *testing.T) {
	tests := []struct {
		borderRadius string
		expected     string
	}{
		{AccordionBorderRadiusNone, ""},
		{AccordionBorderRadiusXs, "first:rounded-t-xs last:rounded-b-xs"},
		{AccordionBorderRadiusSm, "first:rounded-t-sm last:rounded-b-sm"},
		{AccordionBorderRadiusMd, "first:rounded-t last:rounded-b"},
		{AccordionBorderRadiusLg, "first:rounded-t-lg last:rounded-b-lg"},
		{AccordionBorderRadiusXl, "first:rounded-t-xl last:rounded-b-xl"},
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

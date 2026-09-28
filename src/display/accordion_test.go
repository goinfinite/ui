package uiDisplay

import "testing"

func TestAccordionPaddingClassResolver(t *testing.T) {
	tests := []struct {
		paddingSize  string
		defaultClass string
		expected     string
	}{
		{AccordionPaddingSizeNone, "p-3", "p-0"},
		{AccordionPaddingSizeXs, "p-3", "p-1"},
		{AccordionPaddingSizeSm, "p-3", "p-1.5"},
		{AccordionPaddingSizeMd, "p-3", "p-2"},
		{AccordionPaddingSizeLg, "p-3", "p-2.5"},
		{AccordionPaddingSizeXl, "p-3", "p-3"},
		{"", "p-3", "p-3"},
		{"unknown", "p-2", "p-2"},
	}
	for _, test := range tests {
		result := accordionPaddingClassResolver(test.paddingSize, test.defaultClass)
		if result != test.expected {
			t.Errorf("paddingClassResolver(%q, %q) = %q, expected %q",
				test.paddingSize, test.defaultClass, result, test.expected)
		}
	}
}

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

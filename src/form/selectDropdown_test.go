package uiForm

import (
	"strings"
	"testing"
)

func TestSelectDropdownSurfaceClassesResolver(t *testing.T) {
	defaultClasses := selectDropdownSurfaceClassesResolver("", "", "", "", "")
	if !strings.Contains(defaultClasses, "bg-neutral-800/95") {
		t.Errorf("DefaultDropdownBackgroundMissing: %q", defaultClasses)
	}
	if !strings.Contains(defaultClasses, "border-neutral-50/5") {
		t.Errorf("DefaultDropdownBorderMissing: %q", defaultClasses)
	}
	if !strings.Contains(defaultClasses, "rounded-md") {
		t.Errorf("DefaultDropdownBorderRadiusMissing: %q", defaultClasses)
	}
	if !strings.Contains(defaultClasses, "max-h-60") {
		t.Errorf("DefaultDropdownMaxHeightMissing: %q", defaultClasses)
	}
	customClasses := selectDropdownSurfaceClassesResolver(
		"emerald-900", "emerald-500/40", "lg", "min-h-24", "max-h-40",
	)
	expectedFragments := []string{
		"bg-emerald-900", "border-emerald-500/40", "rounded-lg",
		"min-h-24", "max-h-40", "absolute", "overflow-auto",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(customClasses, expectedFragment) {
			t.Errorf(
				"DropdownClassesMissingFragment(%q): %q",
				expectedFragment, customClasses,
			)
		}
	}
}

func TestSelectDropdownMaxVisibleOptionsStyle(t *testing.T) {
	if actual := selectDropdownMaxVisibleOptionsStyle(0, 2); actual != "" {
		t.Errorf("ZeroOptionsShouldYieldNoStyle: %q", actual)
	}

	actualStyle := selectDropdownMaxVisibleOptionsStyle(3, 2)
	if actualStyle != "max-height: 6.125rem" {
		t.Errorf("MaxVisibleOptionsStyleMismatch: %q", actualStyle)
	}
}

func TestSelectLabelValueOptionValuesResolver(t *testing.T) {
	values := selectLabelValueOptionValuesResolver([]SelectLabelValueOption{
		{Label: "seconds", Value: "seconds"},
		{Label: "minutes", Value: "minutes"},
	})
	expectedValues := []string{"seconds", "minutes"}
	if len(values) != len(expectedValues) {
		t.Fatalf(
			"OptionValuesLengthMismatch: got %d, want %d",
			len(values), len(expectedValues),
		)
	}
	for valueIndex, expectedValue := range expectedValues {
		if values[valueIndex] != expectedValue {
			t.Errorf(
				"OptionValueMismatch(%d): got %q, want %q",
				valueIndex, values[valueIndex], expectedValue,
			)
		}
	}
}

package uiToolset

import (
	"strings"
	"testing"
)

func TestTooltipSurfaceClassesResolver(t *testing.T) {
	defaultClasses := TooltipSurfaceClassesResolver(TooltipSurfaceSettings{})
	for _, expectedClass := range []string{
		"fixed", "z-100", "pointer-events-none", "max-w-96",
		"bg-neutral-800/95", "text-neutral-50",
	} {
		if !strings.Contains(defaultClasses, expectedClass) {
			t.Errorf(
				"TooltipSurfaceClassesResolverMissingDefault: %q in %q",
				expectedClass, defaultClasses,
			)
		}
	}
	if strings.Contains(defaultClasses, "ring-") {
		t.Errorf("TooltipSurfaceClassesResolverHasRingWithoutColor: %q", defaultClasses)
	}

	customClasses := TooltipSurfaceClassesResolver(TooltipSurfaceSettings{
		BackgroundColor: "primary-900/95",
		TextColor:       "amber-50",
		RingColor:       "primary-500/40",
		RingThickness:   "xs",
		MinWidthClass:   "min-w-32",
		MaxWidthClass:   "max-w-128",
		MinHeightClass:  "min-h-8",
		MaxHeightClass:  "max-h-64",
	})
	for _, expectedClass := range []string{
		"bg-primary-900/95", "text-amber-50", "ring-0.5 ring-primary-500/40",
		"min-w-32", "max-w-128", "min-h-8", "max-h-64",
	} {
		if !strings.Contains(customClasses, expectedClass) {
			t.Errorf(
				"TooltipSurfaceClassesResolverMissingCustom: %q in %q",
				expectedClass, customClasses,
			)
		}
	}
	if strings.Contains(customClasses, "max-w-96") {
		t.Errorf("TooltipSurfaceClassesResolverKeptDefaultMaxWidth: %q", customClasses)
	}
}

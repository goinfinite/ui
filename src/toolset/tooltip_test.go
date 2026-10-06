package uiToolset

import (
	"strings"
	"testing"
)

func TestTooltipSurfaceClassesResolver(t *testing.T) {
	defaultClasses := TooltipSurfaceClassesResolver("", "", "", "")
	for _, expectedClass := range []string{
		"fixed", "z-100", "pointer-events-none",
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

	customClasses := TooltipSurfaceClassesResolver(
		"primary-900/95", "amber-50", "primary-500/40", "xs",
	)
	for _, expectedClass := range []string{
		"bg-primary-900/95", "text-amber-50", "ring-0.5 ring-primary-500/40",
	} {
		if !strings.Contains(customClasses, expectedClass) {
			t.Errorf(
				"TooltipSurfaceClassesResolverMissingCustom: %q in %q",
				expectedClass, customClasses,
			)
		}
	}
}

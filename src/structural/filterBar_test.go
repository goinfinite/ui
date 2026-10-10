package uiStructural

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestFilterBarChipStyleResolverDefaults(t *testing.T) {
	style := filterChipStyleResolver(FilterChipStyle{})
	if style.TextCase != uiToolset.TextCaseLower {
		t.Errorf("ChipStyleTextCaseDefault = %q, want %q", style.TextCase, uiToolset.TextCaseLower)
	}
	if style.Radius != uiToolset.BorderRadiusMd {
		t.Errorf("ChipStyleRadiusDefault = %q, want %q", style.Radius, uiToolset.BorderRadiusMd)
	}
	if style.OuterBackgroundColor != "neutral-50/10" {
		t.Errorf("OuterBackgroundColorDefault = %q", style.OuterBackgroundColor)
	}
	if style.OuterRingColor != "neutral-50/20" {
		t.Errorf("OuterRingColorDefault = %q", style.OuterRingColor)
	}
	if style.OuterTextColor != "neutral-50" {
		t.Errorf("OuterTextColorDefault = %q", style.OuterTextColor)
	}
	if style.InnerBackgroundColor != "neutral-50/20" {
		t.Errorf("InnerBackgroundColorDefault = %q", style.InnerBackgroundColor)
	}
	if style.InnerTextColor != "neutral-50" {
		t.Errorf("InnerTextColorDefault = %q", style.InnerTextColor)
	}
}

func TestFilterBarChipStyleResolverOverrides(t *testing.T) {
	style := filterChipStyleResolver(FilterChipStyle{
		TextCase:             uiToolset.TextCaseUpper,
		Radius:               uiToolset.BorderRadiusFull,
		OuterBackgroundColor: "secondary-500/20",
		OuterRingColor:       "secondary-500/40",
		OuterTextColor:       "neutral-50",
		InnerBackgroundColor: "secondary-500/40",
		InnerTextColor:       "neutral-50",
	})
	if style.TextCase != uiToolset.TextCaseUpper {
		t.Errorf("ChipStyleTextCase = %q, want %q", style.TextCase, uiToolset.TextCaseUpper)
	}
	if style.Radius != uiToolset.BorderRadiusFull {
		t.Errorf("ChipStyleRadius = %q, want %q", style.Radius, uiToolset.BorderRadiusFull)
	}
	if style.OuterBackgroundColor != "secondary-500/20" {
		t.Errorf("OuterBackgroundColor = %q", style.OuterBackgroundColor)
	}
	if style.OuterRingColor != "secondary-500/40" {
		t.Errorf("OuterRingColor = %q", style.OuterRingColor)
	}
	if style.InnerBackgroundColor != "secondary-500/40" {
		t.Errorf("InnerBackgroundColor = %q", style.InnerBackgroundColor)
	}
}

func TestFilterBarChipAppliesLowercaseDefault(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := FilterBar(FilterBarSettings{
		Filters: []FilterSettings{
			{Key: "status", Label: "Status", Kind: FilterKindTextContains},
		},
		ValuesTwoWayStatePath: "filterValues",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("FilterBarRenderFailed: %v", renderErr)
	}
	if !strings.Contains(buffer.String(), "lowercase") {
		t.Error("RenderedChipMissingLowercaseClass")
	}
}
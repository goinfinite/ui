package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestDualSelectInputRendersFieldStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := DualSelectInput(DualSelectInputSettings{
		Label: "Port",
		Options: []SelectLabelValueOption{
			{Label: "port", Value: "port"},
		},
		Field: DualSelectInputFieldStyle{
			BackgroundColor: "violet-950",
			BorderColor:     "violet-500/40",
			BorderRadius:    uiToolset.BorderRadiusLg,
			DividerColor:    "violet-500/30",
			TextColor:       "violet-100",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DualSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"bg-violet-950", "border-violet-500/40", "rounded-lg",
		"border-violet-500/30", "text-violet-100",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestDualSelectInputRendersDropdownStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := DualSelectInput(DualSelectInputSettings{
		Label: "Port",
		Options: []SelectLabelValueOption{
			{Label: "port", Value: "port"},
		},
		Dropdown: DualSelectInputDropdownStyle{
			BackgroundColor: "violet-900",
			BorderColor:     "violet-500/40",
			BorderRadius:    uiToolset.BorderRadiusLg,
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DualSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"bg-violet-900", "border-violet-500/40", "rounded-lg",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestDualSelectInputRendersValueInputLegendBand(t *testing.T) {
	testCases := []struct {
		name            string
		valueStatePath  string
		expectDynamicCd bool
	}{
		{name: "BoundValueUsesDynamicBand", valueStatePath: "portValue", expectDynamicCd: true},
		{name: "UnboundValueUsesStaticBand", valueStatePath: "", expectDynamicCd: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := DualSelectInput(DualSelectInputSettings{
				Label: "Port",
				Options: []SelectLabelValueOption{
					{Label: "port", Value: "port"},
				},
				ValueTwoWayStatePath: testCase.valueStatePath,
			}).Render(context.Background(), &buffer)
			if renderErr != nil {
				t.Fatalf("DualSelectInputRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			if !strings.Contains(renderedHtml, "-mt-1.4 !pt-1.4") {
				t.Error("RenderedHtmlMissingValueInputLegendBandClasses")
			}
			hasDynamicClass := strings.Contains(renderedHtml, `:class="{`)
			if hasDynamicClass != testCase.expectDynamicCd {
				t.Errorf(
					"DynamicClassMismatch: got %v, want %v",
					hasDynamicClass, testCase.expectDynamicCd,
				)
			}
		})
	}
}

func TestDualSelectInputRendersSharedSlotBinding(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := DualSelectInput(DualSelectInputSettings{
		Label: "Minute",
		Options: []SelectLabelValueOption{
			{Label: "Every minute", Value: "*"},
			{Label: "Every 5 minutes", Value: "*/5"},
		},
		ValueTwoWayStatePath: "minute",
		InputName:            "minute",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DualSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-data="dualSelectInput(`,
		`x-ref="selectTrigger"`,
		`x-ref="selectDropdown"`,
		`x-model="minute"`,
		`x-text="selectedLabelResolver(minute)"`,
		`x-model:value="minute"`,
		`hasUnmatchedSelectValue(minute)`,
		`bg-neutral-950/50`,
		`bg-neutral-800/95`,
		`placeholder:text-neutral-50`,
		`placeholder:opacity-80`,
		`rounded-r`,
		`name="minute"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "<select") {
		t.Errorf("RenderedHtmlKeepsNativeSelect")
	}
}

func TestDualSelectInputRendersSeparateHalves(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := DualSelectInput(DualSelectInputSettings{
		Label: "Timeout",
		Options: []SelectLabelValueOption{
			{Label: "seconds", Value: "seconds"},
		},
		SelectTwoWayStatePath: "timeoutUnit",
		ValueTwoWayStatePath:  "timeoutValue",
		InputType:             InputTypeNumber,
		SelectName:            "timeoutUnit",
		InputName:             "timeoutValue",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("DualSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-model="timeoutUnit"`,
		`x-model:value="timeoutValue"`,
		`type="number"`,
		`name="timeoutUnit"`,
		`name="timeoutValue"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestSelectInputSizeClassesResolver(t *testing.T) {
	testCases := []struct {
		name                      string
		size                      string
		expectedLegendClass       string
		expectedTriggerClass      string
		expectedEmptyPaddingClass string
	}{
		{
			name:                      "xs",
			size:                      SelectInputSizeXs,
			expectedLegendClass:       "text-[0.625rem]",
			expectedTriggerClass:      "px-1 text-xs",
			expectedEmptyPaddingClass: "pt-0.5",
		},
		{
			name:                      "sm",
			size:                      SelectInputSizeSm,
			expectedLegendClass:       "text-[0.625rem]",
			expectedTriggerClass:      "px-1.5 text-xs",
			expectedEmptyPaddingClass: "pt-1",
		},
		{
			name:                      "md",
			size:                      SelectInputSizeMd,
			expectedLegendClass:       "text-xs",
			expectedTriggerClass:      "px-1.5 text-sm",
			expectedEmptyPaddingClass: "pt-1.5",
		},
		{
			name:                      "lg",
			size:                      SelectInputSizeLg,
			expectedLegendClass:       "text-sm",
			expectedTriggerClass:      "px-2 text-base",
			expectedEmptyPaddingClass: "pt-2",
		},
		{
			name:                      "xl",
			size:                      SelectInputSizeXl,
			expectedLegendClass:       "text-base",
			expectedTriggerClass:      "px-2.5 text-lg",
			expectedEmptyPaddingClass: "pt-2.5",
		},
		{
			name:                      "default",
			size:                      "",
			expectedLegendClass:       "text-xs",
			expectedTriggerClass:      "px-1.5 text-sm",
			expectedEmptyPaddingClass: "pt-1.5",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			legendClass, triggerClass, emptyPaddingClass := selectInputSizeClassesResolver(testCase.size)
			actualClasses := []string{legendClass, triggerClass, emptyPaddingClass}
			expectedClasses := []string{
				testCase.expectedLegendClass, testCase.expectedTriggerClass, testCase.expectedEmptyPaddingClass,
			}
			if strings.Join(actualClasses, "|") != strings.Join(expectedClasses, "|") {
				t.Errorf(
					"SizeClassesMismatch(%q): got %q, want %q",
					testCase.size, actualClasses, expectedClasses,
				)
			}
		})
	}
}

func TestSelectInputRendersOpaqueDropdownAboveStickyHeaders(t *testing.T) {
	settings := SelectInputSettings{
		InputName:   "status",
		Label:       "Status",
		FlatOptions: []string{"running"},
	}
	var buffer bytes.Buffer
	renderErr := SelectInput(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"z-20", "bg-neutral-800/95"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}

	settings.DropdownBackgroundColor = "emerald-900"
	coloredBuffer := bytes.Buffer{}
	renderErr = SelectInput(settings).Render(context.Background(), &coloredBuffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	coloredHtml := coloredBuffer.String()
	if !strings.Contains(coloredHtml, "bg-emerald-900") {
		t.Errorf("RenderedHtmlMissingCustomDropdownBackground")
	}
	if strings.Contains(coloredHtml, "bg-neutral-800/95") {
		t.Errorf("RenderedHtmlKeepsDefaultDropdownBackground")
	}
}

func TestSelectInputRendersTriggerDropdownAndIconStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		InputName:            "country",
		Label:                "Country",
		FlatOptions:          []string{"Argentina", "Brazil"},
		TwoWayStatePath:      "country",
		SelectedIconColor:    "violet-500",
		DropdownBorderColor:  "violet-500/20",
		DropdownBorderRadius: uiToolset.BorderRadiusLg,
		Trigger: SelectInputTriggerStyle{
			BackgroundColor: "violet-950",
			BorderColor:     "violet-500/40",
			BorderRadius:    uiToolset.BorderRadiusLg,
			TextColor:       "violet-100",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"bg-violet-950", "border-violet-500/40", "rounded-lg", "text-violet-100",
		"border-violet-500/20", "peer-checked:text-violet-500",
		"peer-focus-visible:ring-violet-500", "ph-circle",
		`@keydown.arrow-down.prevent="openWithOptionFocus()"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestSelectInputRendersInvalidStateBinding(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		Label:                    "Role",
		FlatOptions:              []string{"Owner", "Member"},
		TwoWayStatePath:          "roleValue",
		IsInvalidOneWayStatePath: "roleIsInvalid",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`:class="{&#39;mt-0&#39;: !(roleValue), &#39;!border-red-500&#39;: roleIsInvalid}"`,
		`:aria-invalid="roleIsInvalid"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestSelectInputRendersNeutralSurfaceByDefault(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		InputName:   "status",
		Label:       "Status",
		FlatOptions: []string{"running"},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		"bg-neutral-50/7.5", "border-neutral-50/5", "rounded", "text-neutral-100",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "ph-circle") {
		t.Errorf("RenderedHtmlKeepsSelectedIconWithoutColorSetting")
	}
}

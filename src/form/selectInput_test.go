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

func TestSelectInputEmptyStateCheckExpressionBuilder(t *testing.T) {
	testCases := []struct {
		name          string
		statePath     string
		isMultiSelect bool
		expected      string
	}{
		{
			name:      "single select checks explicit empty values",
			statePath: "selectedRole",
			expected:  "(selectedRole === null || selectedRole === undefined || selectedRole === '')",
		},
		{
			name:          "multi select checks its length",
			statePath:     "selectedRoles",
			isMultiSelect: true,
			expected:      "!(selectedRoles.length)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualExpression := selectInputEmptyStateCheckExpressionBuilder(
				testCase.statePath, testCase.isMultiSelect,
			)
			if actualExpression != testCase.expected {
				t.Errorf(
					"EmptyStateCheckExpressionMismatch(%q): got %q, want %q",
					testCase.statePath, actualExpression, testCase.expected,
				)
			}
		})
	}
}

func TestSelectInputOptionAlias(t *testing.T) {
	testCases := []struct {
		name         string
		statePath    string
		onChangeFunc string
		alias        string
		expected     string
	}{
		{
			name:      "matching alias stays unchanged",
			statePath: "selectedRole",
			alias:     "dynamicOption",
			expected:  "dynamicOption",
		},
		{
			name:      "option alias gains a suffix",
			statePath: "dynamicOption",
			alias:     "dynamicOption",
			expected:  "dynamicOption_",
		},
		{
			name:      "option alias prefix gains a suffix",
			statePath: "dynamicOption.value",
			alias:     "dynamicOption",
			expected:  "dynamicOption_",
		},
		{
			name:      "index alias gains a suffix",
			statePath: "optionIndex",
			alias:     "optionIndex",
			expected:  "optionIndex_",
		},
		{
			name:      "index alias prefix gains a suffix",
			statePath: "optionIndex.value",
			alias:     "optionIndex",
			expected:  "optionIndex_",
		},
		{
			name:         "index alias yields to the callback",
			statePath:    "selected",
			onChangeFunc: "lastIndex = optionIndex",
			alias:        "optionIndex",
			expected:     "optionIndex_",
		},
		{
			name:         "option alias yields to the callback",
			statePath:    "selected",
			onChangeFunc: "lastOption = dynamicOption",
			alias:        "dynamicOption",
			expected:     "dynamicOption_",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualAlias := selectInputOptionAlias(
				testCase.statePath, testCase.onChangeFunc, testCase.alias,
			)
			if actualAlias != testCase.expected {
				t.Errorf(
					"OptionAliasMismatch(%q, %q, %q): got %q, want %q",
					testCase.statePath, testCase.onChangeFunc, testCase.alias,
					actualAlias, testCase.expected,
				)
			}
		})
	}
}

func TestSelectInputOptionIconExpressionBuilder(t *testing.T) {
	actualExpression := selectInputOptionIconExpressionBuilder(
		"selectedValue", `"admin"`,
	)
	expectedExpression := `String(selectedValue) === String("admin") ? "ph-check-circle" : "ph-circle"`
	if actualExpression != expectedExpression {
		t.Errorf(
			"IconExpressionMismatch: got %q, want %q",
			actualExpression, expectedExpression,
		)
	}
}

func TestSelectInputDynamicOptionIconExpressionBuilder(t *testing.T) {
	actualExpression := selectInputDynamicOptionIconExpressionBuilder(
		"selectedValue", "dynamicOption?.value ?? dynamicOption",
	)
	expectedExpression := `selectedValue === (dynamicOption?.value ?? dynamicOption) ? "ph-check-circle" : "ph-circle"`
	if actualExpression != expectedExpression {
		t.Errorf(
			"DynamicIconExpressionMismatch: got %q, want %q",
			actualExpression, expectedExpression,
		)
	}
}

func TestSelectInputRendersDynamicOptionsStatePath(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		InputName:              "status",
		Label:                  "Status",
		OptionsOneWayStatePath: "dynamicOptions",
		TwoWayStatePath:        "selectedStatus",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-for="(dynamicOption, optionIndex) in dynamicOptions"`,
		`dynamicOptionsDisplayResolver(dynamicOptions, selectedStatus)`,
		`:value="dynamicOption?.value ?? dynamicOption"`,
		`x-text="dynamicOption?.label ?? dynamicOption?.value ?? dynamicOption"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestSelectInputDynamicOptionDoesNotShadowStatePath(t *testing.T) {
	testCases := []struct {
		name              string
		statePath         string
		expectedFragments []string
	}{
		{
			name:      "option",
			statePath: "option",
			expectedFragments: []string{
				`x-for="(dynamicOption, optionIndex) in dynamicOptions"`,
				`@change="option = dynamicOption?.value ?? dynamicOption"`,
				`option === (dynamicOption?.value ?? dynamicOption)`,
			},
		},
		{
			name:      "dynamicOption",
			statePath: "dynamicOption",
			expectedFragments: []string{
				`x-for="(dynamicOption_, optionIndex) in dynamicOptions"`,
				`@change="dynamicOption = dynamicOption_?.value ?? dynamicOption_"`,
				`dynamicOption === (dynamicOption_?.value ?? dynamicOption_)`,
			},
		},
		{
			name:      "optionIndex",
			statePath: "optionIndex",
			expectedFragments: []string{
				`x-for="(dynamicOption, optionIndex_) in dynamicOptions"`,
				`@change="optionIndex = dynamicOption?.value ?? dynamicOption"`,
				`optionIndex === (dynamicOption?.value ?? dynamicOption)`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := SelectInput(SelectInputSettings{
				InputName:              "status",
				Label:                  "Status",
				OptionsOneWayStatePath: "dynamicOptions",
				TwoWayStatePath:        testCase.statePath,
			}).Render(context.Background(), &buffer)
			if renderErr != nil {
				t.Fatalf("SelectInputRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			for _, expectedFragment := range testCase.expectedFragments {
				if !strings.Contains(renderedHtml, expectedFragment) {
					t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
				}
			}
		})
	}
}

func TestSelectInputMultiSelectOptionIndexAlias(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		InputId:                "statusOptions",
		InputName:              "status",
		Label:                  "Status",
		OptionsOneWayStatePath: "dynamicOptions",
		TwoWayStatePath:        "optionIndex",
		IsMultiSelect:          true,
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-for="(dynamicOption, optionIndex_) in dynamicOptions"`,
		`&#34;statusOptions-&#34; + optionIndex_`,
		`@change="optionIndex = toggleSelectedOptionValue(optionIndex, dynamicOption?.value ?? dynamicOption, $event.target.checked)"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestSelectInputOnChangeFuncAliasDoesNotShadowStatePath(t *testing.T) {
	testCases := []struct {
		name              string
		onChangeFunc      string
		expectedFragments []string
	}{
		{
			name:         "option index callback",
			onChangeFunc: "lastIndex = optionIndex",
			expectedFragments: []string{
				`x-for="(dynamicOption, optionIndex_) in dynamicOptions"`,
				`@change="selected = dynamicOption?.value ?? dynamicOption; lastIndex = optionIndex"`,
			},
		},
		{
			name:         "dynamic option callback",
			onChangeFunc: "lastOption = dynamicOption",
			expectedFragments: []string{
				`x-for="(dynamicOption_, optionIndex) in dynamicOptions"`,
				`@change="selected = dynamicOption_?.value ?? dynamicOption_; lastOption = dynamicOption"`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := SelectInput(SelectInputSettings{
				InputName:              "status",
				Label:                  "Status",
				OptionsOneWayStatePath: "dynamicOptions",
				TwoWayStatePath:        "selected",
				OnChangeFunc:           testCase.onChangeFunc,
			}).Render(context.Background(), &buffer)
			if renderErr != nil {
				t.Fatalf("SelectInputRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			for _, expectedFragment := range testCase.expectedFragments {
				if !strings.Contains(renderedHtml, expectedFragment) {
					t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
				}
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
		`:class="{&#39;mt-0&#39;: (roleValue === null || roleValue === undefined || roleValue === &#39;&#39;), &#39;!border-red-500&#39;: roleIsInvalid}"`,
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

func TestSelectInputRendersLabelStatePath(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := SelectInput(SelectInputSettings{
		InputName:            "role",
		Label:                "Role",
		LabelOneWayStatePath: "roleLabel",
		FlatOptions:          []string{"Owner", "Member"},
		TwoWayStatePath:      "roleValue",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("SelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="roleLabel"`,
		"Role",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

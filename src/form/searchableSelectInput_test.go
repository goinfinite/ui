package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestSearchableSelectInputItemsResolver(t *testing.T) {
	items := searchableSelectInputItemsResolver(
		[]string{"Brazil"},
		[]SelectLabelValueOption{{Label: "Argentina", Value: "AR"}},
	)
	if len(items) != 2 {
		t.Fatalf("ItemsLengthMismatch: got %d, want 2", len(items))
	}
	if items[0].Label != "Brazil" || items[0].Value != "Brazil" ||
		items[0].SearchableText != "Brazil" {
		t.Errorf("FlatOptionMismatch: got %+v", items[0])
	}
	if items[1].Label != "Argentina" || items[1].Value != "AR" ||
		items[1].SearchableText != "Argentina" {
		t.Errorf("LabelValueOptionMismatch: got %+v", items[1])
	}
}

func TestSearchableSelectInputSelectionExpressionBuilders(t *testing.T) {
	testCases := []struct {
		name     string
		actual   string
		expected string
	}{
		{
			name: "single select",
			actual: searchableSelectInputSingleSelectExpressionBuilder(
				"country", "'AR'", "'Argentina'", "; onCountryChange()",
			),
			expected: "selectedLabelCache.set('AR', 'Argentina'); country = 'AR';" +
				" userInput = 'Argentina'; closeDropdown(); onCountryChange()",
		},
		{
			name: "multi toggle",
			actual: searchableSelectInputMultiToggleExpressionBuilder(
				"countries", "'AR'", "'Argentina'", "",
			),
			expected: "selectedLabelCache.set('AR', 'Argentina'); countries =" +
				" countries.includes('AR') ? countries.filter((item) =>" +
				" item !== 'AR') : [...countries, 'AR']; userInput = ''",
		},
		{
			name: "clear multi",
			actual: searchableSelectInputClearExpressionBuilder(
				"countries", true, "",
			),
			expected: "countries = []; userInput = ''",
		},
		{
			name: "clear single",
			actual: searchableSelectInputClearExpressionBuilder(
				"country", false, "",
			),
			expected: "country = ''; userInput = ''",
		},
		{
			name: "custom value add",
			actual: searchableSelectInputCustomValueAddExpressionBuilder(
				"countries", "",
			),
			expected: "if (userInput.trim() !== '') { countries = countries" +
				".includes(userInput.trim()) ? countries : [...countries," +
				" userInput.trim()] } userInput = ''",
		},
		{
			name: "tag remove",
			actual: searchableSelectInputTagRemoveExpressionBuilder(
				"countries", "selectedValue", "; onRemove()",
			),
			expected: "countries = countries.filter((item) => item !==" +
				" selectedValue); onRemove()",
		},
		{
			name: "tag backspace",
			actual: searchableSelectInputTagBackspaceExpressionBuilder(
				"countries", "",
			),
			expected: "if (userInput === '' && countries.length > 0) {" +
				" countries = countries.slice(0, -1) }",
		},
		{
			name: "selected single",
			actual: searchableSelectInputSelectedExpressionBuilder(
				"country", "'AR'", false,
			),
			expected: "String(country) === 'AR'",
		},
		{
			name: "selected multi",
			actual: searchableSelectInputSelectedExpressionBuilder(
				"countries", "'AR'", true,
			),
			expected: "countries.includes('AR')",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.actual != testCase.expected {
				t.Errorf(
					"ExpressionMismatch: got %q, want %q",
					testCase.actual, testCase.expected,
				)
			}
		})
	}
}

func TestSearchableSelectInputSelectionPresentationResolver(t *testing.T) {
	testCases := []struct {
		name                    string
		isMultiSelect           bool
		selectionDisplay        string
		shouldAllowCustomValues bool
		expected                searchableSelectInputSelectionPresentation
	}{
		{
			name:                    "single select always shows text",
			isMultiSelect:           false,
			selectionDisplay:        SearchableSelectInputSelectionDisplayTags,
			shouldAllowCustomValues: true,
			expected: searchableSelectInputSelectionPresentation{
				Display: SearchableSelectInputSelectionDisplayText,
			},
		},
		{
			name:          "multi select defaults to text",
			isMultiSelect: true,
			expected: searchableSelectInputSelectionPresentation{
				Display: SearchableSelectInputSelectionDisplayText,
			},
		},
		{
			name:                    "multi select tags keep custom values",
			isMultiSelect:           true,
			selectionDisplay:        SearchableSelectInputSelectionDisplayTags,
			shouldAllowCustomValues: true,
			expected: searchableSelectInputSelectionPresentation{
				Display:           SearchableSelectInputSelectionDisplayTags,
				AllowCustomValues: true,
			},
		},
		{
			name:             "multi select tags without custom values",
			isMultiSelect:    true,
			selectionDisplay: SearchableSelectInputSelectionDisplayTags,
			expected: searchableSelectInputSelectionPresentation{
				Display: SearchableSelectInputSelectionDisplayTags,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualPresentation := searchableSelectInputSelectionPresentationResolver(
				testCase.isMultiSelect, testCase.selectionDisplay,
				testCase.shouldAllowCustomValues,
			)
			if actualPresentation != testCase.expected {
				t.Errorf(
					"SelectionPresentationMismatch: got %+v, want %+v",
					actualPresentation, testCase.expected,
				)
			}
		})
	}
}

func TestSearchableSelectInputRendersTagsAndCustomValues(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SearchableSelectInput(SearchableSelectInputSettings{
		InputName:               "countries",
		Label:                   "Countries",
		FlatOptions:             []string{"Argentina", "Brazil"},
		TwoWayStatePath:         "countries",
		IsMultiSelect:           true,
		SelectionDisplay:        SearchableSelectInputSelectionDisplayTags,
		ShouldAllowCustomValues: true,
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SearchableSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{
		"optionLabelResolver(selectedValue)", "Remove ", "userInput.trim()",
		`name="countries"`,
	} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

func TestSearchableSelectInputRendersDropdownMaxHeightAndEmptyState(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := SearchableSelectInput(SearchableSelectInputSettings{
		InputName:              "country",
		Label:                  "Country",
		FlatOptions:            []string{"Argentina", "Brazil"},
		TwoWayStatePath:        "country",
		DropdownMaxHeightClass: "max-h-32",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SearchableSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{
		"max-h-32", "No matches", "searchableSelectInput",
		"role=\"combobox\"",
	} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

func TestSearchableSelectInputRendersOpaqueDropdownAboveStickyHeaders(t *testing.T) {
	settings := SearchableSelectInputSettings{
		InputName:       "country",
		Label:           "Country",
		FlatOptions:     []string{"Argentina", "Brazil"},
		TwoWayStatePath: "country",
	}
	var buffer bytes.Buffer
	renderErr := SearchableSelectInput(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("SearchableSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"z-20", "bg-neutral-800/95"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}

	settings.DropdownBackgroundColor = "emerald-900"
	coloredBuffer := bytes.Buffer{}
	renderErr = SearchableSelectInput(settings).Render(context.Background(), &coloredBuffer)
	if renderErr != nil {
		t.Fatalf("SearchableSelectInputRenderFailed: %v", renderErr)
	}
	coloredHtml := coloredBuffer.String()
	if !strings.Contains(coloredHtml, "bg-emerald-900") {
		t.Errorf("RenderedHtmlMissingCustomDropdownBackground")
	}
	if strings.Contains(coloredHtml, "bg-neutral-800/95") {
		t.Errorf("RenderedHtmlKeepsDefaultDropdownBackground")
	}
}

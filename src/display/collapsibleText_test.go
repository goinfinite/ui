package uiDisplay

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestCollapsibleTextMaxLinesResolver(t *testing.T) {
	tests := []struct {
		maxVisibleLines int
		expected        int
	}{
		{0, 3},
		{-5, 3},
		{1, 1},
		{3, 3},
		{12, 12},
	}
	for _, test := range tests {
		result := collapsibleTextMaxLinesResolver(test.maxVisibleLines)
		if result != test.expected {
			t.Errorf("maxLinesResolver(%d) = %d, expected %d",
				test.maxVisibleLines, result, test.expected)
		}
	}
}

func TestCollapsibleTextSizeClassResolver(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{CollapsibleTextSizeXs, "text-xs"},
		{CollapsibleTextSizeSm, "text-sm"},
		{CollapsibleTextSizeMd, "text-base"},
		{CollapsibleTextSizeLg, "text-lg"},
		{"", "text-sm"},
		{"unknown", "text-sm"},
	}
	for _, test := range tests {
		result := collapsibleTextSizeClassResolver(test.size)
		if result != test.expected {
			t.Errorf("sizeClassResolver(%q) = %q, expected %q",
				test.size, result, test.expected)
		}
	}
}

func TestCollapsibleTextLabelsResolver(t *testing.T) {
	tests := []struct {
		expandLabel           string
		collapseLabel         string
		expectedExpandLabel   string
		expectedCollapseLabel string
	}{
		{"", "", "show more", "show less"},
		{"Ver mais", "", "Ver mais", "show less"},
		{"", "Ver menos", "show more", "Ver menos"},
		{"Ver mais", "Ver menos", "Ver mais", "Ver menos"},
	}
	for _, test := range tests {
		expandLabel, collapseLabel := collapsibleTextLabelsResolver(
			test.expandLabel, test.collapseLabel,
		)
		if expandLabel != test.expectedExpandLabel {
			t.Errorf("labelsResolver(%q, %q) expand = %q, expected %q",
				test.expandLabel, test.collapseLabel,
				expandLabel, test.expectedExpandLabel)
		}
		if collapseLabel != test.expectedCollapseLabel {
			t.Errorf("labelsResolver(%q, %q) collapse = %q, expected %q",
				test.expandLabel, test.collapseLabel,
				collapseLabel, test.expectedCollapseLabel)
		}
	}
}

func TestCollapsibleTextRender(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CollapsibleText(CollapsibleTextSettings{
		Text:            "A long log line",
		MaxVisibleLines: 5,
		ExpandLabel:     "Open",
		CollapseLabel:   "Close",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CollapsibleTextRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	textIdPattern := regexp.MustCompile(`id="(ui-collapsible-text-\d+)"`)
	textIdMatches := textIdPattern.FindStringSubmatch(renderedHtml)
	if textIdMatches == nil {
		t.Fatal("RenderedHtmlMissingTextId")
	}
	expectedFragments := []string{
		">A long log line</span>",
		`aria-controls="` + textIdMatches[1] + `"`,
		`x-show="hasToggle"`,
		`collapsibleText(&#34;A long log line&#34;,5,&#34;Open&#34;,&#34;Close&#34;)`,
		`x-text="displayText"`,
		`x-text="toggleLabel"`,
		`whitespace-pre-wrap`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCollapsibleTextRenderDefaults(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CollapsibleText(CollapsibleTextSettings{
		Text: "Short text",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CollapsibleTextRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`collapsibleText(&#34;Short text&#34;,3,&#34;show more&#34;,&#34;show less&#34;)`,
		"text-sm",
		"text-neutral-300",
		"text-secondary-400",
		`ph-caret-down`,
		`ph-caret-up`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "shrink-0") {
		t.Error("RenderedHtmlContainsIconWithoutIconSetting")
	}
}

func TestCollapsibleTextRenderToggleIcons(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CollapsibleText(CollapsibleTextSettings{
		Text:         "A long log line",
		ExpandIcon:   "ph-plus",
		CollapseIcon: "ph-minus",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CollapsibleTextRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`ph-plus`,
		`ph-minus`,
		`x-show="!isExpanded"`,
		`x-show="isExpanded"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCollapsibleTextRenderIcons(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CollapsibleText(CollapsibleTextSettings{
		Text:           "A long log line",
		FlankIconLeft:  "ph-file-text",
		FlankIconRight: "ph-caret-up-down",
		IconColor:      "sky-400",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CollapsibleTextRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`ph-file-text`,
		`ph-caret-up-down`,
		"text-sky-400",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "text-neutral-400") {
		t.Error("RenderedHtmlContainsDefaultIconColorWithIconColorSet")
	}
}

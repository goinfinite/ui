package uiStructural

import "testing"

func TestTabOrientationResolver(t *testing.T) {
	tests := []struct {
		name        string
		orientation string
		want        string
	}{
		{name: "empty defaults to horizontal", orientation: "", want: TabOrientationHorizontal},
		{name: "horizontal", orientation: TabOrientationHorizontal, want: TabOrientationHorizontal},
		{name: "vertical", orientation: TabOrientationVertical, want: TabOrientationVertical},
		{name: "unknown defaults to horizontal", orientation: "diagonal", want: TabOrientationHorizontal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabOrientationResolver(tt.orientation); got != tt.want {
				t.Errorf("tabOrientationResolver(%q) = %q, want %q", tt.orientation, got, tt.want)
			}
		})
	}
}

func TestTabAriaLabelResolver(t *testing.T) {
	tests := []struct {
		name      string
		ariaLabel string
		want      string
	}{
		{name: "empty defaults to Tabs", ariaLabel: "", want: "Tabs"},
		{name: "custom label", ariaLabel: "Settings sections", want: "Settings sections"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabAriaLabelResolver(tt.ariaLabel); got != tt.want {
				t.Errorf("tabAriaLabelResolver(%q) = %q, want %q", tt.ariaLabel, got, tt.want)
			}
		})
	}
}

func TestTabSelectedExpressionBuilder(t *testing.T) {
	tests := []struct {
		name              string
		selectedValuePath string
		itemValue         string
		want              string
	}{
		{
			name:              "compares the bound path to the item value",
			selectedValuePath: "selectedTab",
			itemValue:         "security",
			want:              `selectedTab === "security"`,
		},
		{
			name:              "supports nested paths",
			selectedValuePath: "settings.activeSection",
			itemValue:         "general",
			want:              `settings.activeSection === "general"`,
		},
		{
			name:              "escapes a value that contains a quote",
			selectedValuePath: "selectedTab",
			itemValue:         `it's`,
			want:              `selectedTab === "it's"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabSelectedExpressionBuilder(tt.selectedValuePath, tt.itemValue); got != tt.want {
				t.Errorf("tabSelectedExpressionBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabSideResolver(t *testing.T) {
	tests := []struct {
		name string
		side string
		want string
	}{
		{name: "empty defaults to left", side: "", want: TabSideLeft},
		{name: "left", side: TabSideLeft, want: TabSideLeft},
		{name: "right", side: TabSideRight, want: TabSideRight},
		{name: "unknown defaults to left", side: "diagonal", want: TabSideLeft},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabSideResolver(tt.side); got != tt.want {
				t.Errorf("tabSideResolver(%q) = %q, want %q", tt.side, got, tt.want)
			}
		})
	}
}

func TestTabBorderRadiusTokenResolver(t *testing.T) {
	tests := []struct {
		name         string
		borderRadius string
		want         string
	}{
		{name: "empty defaults to md", borderRadius: "", want: TabBorderRadiusMd},
		{name: "none", borderRadius: TabBorderRadiusNone, want: TabBorderRadiusNone},
		{name: "xs", borderRadius: TabBorderRadiusXs, want: TabBorderRadiusXs},
		{name: "sm", borderRadius: TabBorderRadiusSm, want: TabBorderRadiusSm},
		{name: "md", borderRadius: TabBorderRadiusMd, want: TabBorderRadiusMd},
		{name: "lg", borderRadius: TabBorderRadiusLg, want: TabBorderRadiusLg},
		{name: "xl", borderRadius: TabBorderRadiusXl, want: TabBorderRadiusXl},
		{name: "unknown defaults to md", borderRadius: "huge", want: TabBorderRadiusMd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabBorderRadiusTokenResolver(tt.borderRadius); got != tt.want {
				t.Errorf("tabBorderRadiusTokenResolver(%q) = %q, want %q", tt.borderRadius, got, tt.want)
			}
		})
	}
}

func TestTabEdgeClassResolver(t *testing.T) {
	tests := []struct {
		name         string
		borderRadius string
		isVertical   bool
		isRight      bool
		want         string
	}{
		{
			name:         "horizontal defaults to rounded top md",
			borderRadius: "",
			isVertical:   false,
			want:         "rounded-t-md border-b-2 border-transparent -mb-px",
		},
		{
			name:         "vertical left rounds the left edge",
			borderRadius: "",
			isVertical:   true,
			isRight:      false,
			want:         "rounded-l-md border-r-2 border-transparent -mr-px",
		},
		{
			name:         "vertical right rounds the right edge",
			borderRadius: "",
			isVertical:   true,
			isRight:      true,
			want:         "rounded-r-md border-l-2 border-transparent -ml-px",
		},
		{
			name:         "custom radius replaces the md token",
			borderRadius: TabBorderRadiusNone,
			isVertical:   false,
			want:         "rounded-t-none border-b-2 border-transparent -mb-px",
		},
		{
			name:         "unknown radius falls back to md",
			borderRadius: "huge",
			isVertical:   true,
			isRight:      true,
			want:         "rounded-r-md border-l-2 border-transparent -ml-px",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabEdgeClassResolver(tt.borderRadius, tt.isVertical, tt.isRight); got != tt.want {
				t.Errorf("tabEdgeClassResolver() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabPercentClassResolver(t *testing.T) {
	tests := []struct {
		name        string
		classPrefix string
		percent     uint
		want        string
	}{
		{name: "zero leaves the size unset", classPrefix: "max-w", percent: 0, want: ""},
		{name: "percent builds the arbitrary class", classPrefix: "max-w", percent: 80, want: "max-w-[80%]"},
		{name: "min prefix builds the floor class", classPrefix: "min-h", percent: 25, want: "min-h-[25%]"},
		{name: "hundred renders the full class", classPrefix: "max-h", percent: 100, want: "max-h-full"},
		{name: "over a hundred renders the full class", classPrefix: "min-w", percent: 150, want: "min-w-full"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabPercentClassResolver(tt.classPrefix, tt.percent); got != tt.want {
				t.Errorf("tabPercentClassResolver() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabConstraintClassesResolver(t *testing.T) {
	tests := []struct {
		name              string
		componentSettings TabSettings
		want              string
	}{
		{
			name:              "empty settings render no constraints",
			componentSettings: TabSettings{},
			want:              "",
		},
		{
			name: "percent constraints join in axis order",
			componentSettings: TabSettings{
				MinWidthPercent:  25,
				MaxWidthPercent:  80,
				MinHeightPercent: 10,
				MaxHeightPercent: 60,
			},
			want: "min-w-[25%] max-w-[80%] min-h-[10%] max-h-[60%]",
		},
		{
			name: "class constraints join after the percent constraints",
			componentSettings: TabSettings{
				MaxWidthPercent: 80,
				MinHeightClass:  "min-h-48",
				MaxHeightClass:  "max-h-64",
			},
			want: "max-w-[80%] min-h-48 max-h-64",
		},
		{
			name:              "hundred percent renders the full class",
			componentSettings: TabSettings{MaxWidthPercent: 100},
			want:              "max-w-full",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabConstraintClassesResolver(tt.componentSettings); got != tt.want {
				t.Errorf("tabConstraintClassesResolver() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabAlignmentClassResolver(t *testing.T) {
	tests := []struct {
		name      string
		alignment string
		want      string
	}{
		{name: "empty defaults to the top alignment", alignment: "", want: "self-start"},
		{name: "top", alignment: TabAlignmentTop, want: "self-start"},
		{name: "center", alignment: TabAlignmentCenter, want: "self-center"},
		{name: "bottom", alignment: TabAlignmentBottom, want: "self-end"},
		{name: "unknown defaults to the top alignment", alignment: "middle", want: "self-start"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabAlignmentClassResolver(tt.alignment); got != tt.want {
				t.Errorf("tabAlignmentClassResolver() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabBadgeClassesResolver(t *testing.T) {
	tests := []struct {
		name string
		item TabItemSettings
		want string
	}{
		{
			name: "defaults render the neutral badge",
			item: TabItemSettings{Label: "Security", Value: "security", BadgeCount: 3},
			want: "px-1.5 py-0.5 text-xs font-bold rounded-md bg-neutral-50/10 text-neutral-50/80",
		},
		{
			name: "custom colors and radius replace the defaults",
			item: TabItemSettings{
				BadgeBackgroundColor: "red-500/20",
				BadgeBorderRadius:    TabBorderRadiusSm,
				BadgeTextColor:       "red-50",
			},
			want: "px-1.5 py-0.5 text-xs font-bold rounded-sm bg-red-500/20 text-red-50",
		},
		{
			name: "custom ring frames the badge",
			item: TabItemSettings{
				BadgeRingColor:     "red-500/50",
				BadgeRingThickness: TabRingThicknessSm,
			},
			want: "px-1.5 py-0.5 text-xs font-bold rounded-md bg-neutral-50/10 text-neutral-50/80 ring-1.5 ring-red-500/50",
		},
		{
			name: "unknown radius falls back to md",
			item: TabItemSettings{BadgeBorderRadius: "huge"},
			want: "px-1.5 py-0.5 text-xs font-bold rounded-md bg-neutral-50/10 text-neutral-50/80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabBadgeClassesResolver(tt.item); got != tt.want {
				t.Errorf("tabBadgeClassesResolver() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabClickExpressionBuilder(t *testing.T) {
	tests := []struct {
		name              string
		selectedValuePath string
		itemValue         string
		isUrlHashSynced   bool
		want              string
	}{
		{
			name:              "sets the bound path",
			selectedValuePath: "selectedTab",
			itemValue:         "general",
			isUrlHashSynced:   false,
			want:              `selectedTab = "general"`,
		},
		{
			name:              "url hash synced click also writes the hash",
			selectedValuePath: "selectedTab",
			itemValue:         "general",
			isUrlHashSynced:   true,
			want:              `selectedTab = "general"; location.hash = "general"`,
		},
		{
			name:              "escapes a value that contains a quote",
			selectedValuePath: "selectedTab",
			itemValue:         `it's`,
			isUrlHashSynced:   true,
			want:              `selectedTab = "it's"; location.hash = "it's"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabClickExpressionBuilder(tt.selectedValuePath, tt.itemValue, tt.isUrlHashSynced); got != tt.want {
				t.Errorf("tabClickExpressionBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabUrlHashSyncExpressionBuilder(t *testing.T) {
	items := []TabItemSettings{
		{Label: "General", Value: "general"},
		{Label: "Security", Value: "security"},
	}
	want := `selectedTab = ["general","security"].includes(location.hash.slice(1)) ? location.hash.slice(1) : selectedTab`
	if got := tabUrlHashSyncExpressionBuilder("selectedTab", items); got != want {
		t.Errorf("tabUrlHashSyncExpressionBuilder() = %q, want %q", got, want)
	}
}

func TestTabArrowKeyExpressionBuilder(t *testing.T) {
	tests := []struct {
		name              string
		selectedValuePath string
		siblingDirection  string
		want              string
	}{
		{
			name:              "previous sibling focus and select",
			selectedValuePath: "selectedTab",
			siblingDirection:  "previous",
			want:              "if ($el.previousElementSibling) { $el.previousElementSibling.focus(); selectedTab = $el.previousElementSibling.dataset.tabValue }",
		},
		{
			name:              "next sibling focus and select",
			selectedValuePath: "selectedTab",
			siblingDirection:  "next",
			want:              "if ($el.nextElementSibling) { $el.nextElementSibling.focus(); selectedTab = $el.nextElementSibling.dataset.tabValue }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabArrowKeyExpressionBuilder(tt.selectedValuePath, tt.siblingDirection); got != tt.want {
				t.Errorf("tabArrowKeyExpressionBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabEdgeKeyExpressionBuilder(t *testing.T) {
	tests := []struct {
		name              string
		selectedValuePath string
		edgeProperty      string
		want              string
	}{
		{
			name:              "home focuses and selects the first tab",
			selectedValuePath: "selectedTab",
			edgeProperty:      "firstElementChild",
			want:              "$el.parentElement.firstElementChild.focus(); selectedTab = $el.parentElement.firstElementChild.dataset.tabValue",
		},
		{
			name:              "end focuses and selects the last tab",
			selectedValuePath: "selectedTab",
			edgeProperty:      "lastElementChild",
			want:              "$el.parentElement.lastElementChild.focus(); selectedTab = $el.parentElement.lastElementChild.dataset.tabValue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabEdgeKeyExpressionBuilder(tt.selectedValuePath, tt.edgeProperty); got != tt.want {
				t.Errorf("tabEdgeKeyExpressionBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabIdTextBuilder(t *testing.T) {
	tests := []struct {
		name     string
		idPrefix string
		index    int
		want     string
	}{
		{name: "first tab id", idPrefix: "ui-tab-1", index: 0, want: "ui-tab-1-0"},
		{name: "later tab id", idPrefix: "ui-tab-1", index: 2, want: "ui-tab-1-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabIdTextBuilder(tt.idPrefix, tt.index); got != tt.want {
				t.Errorf("tabIdTextBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabPanelIdTextBuilder(t *testing.T) {
	tests := []struct {
		name     string
		idPrefix string
		index    int
		want     string
	}{
		{name: "first panel id", idPrefix: "ui-tab-1", index: 0, want: "ui-tab-1-panel-0"},
		{name: "later panel id", idPrefix: "ui-tab-1", index: 2, want: "ui-tab-1-panel-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabPanelIdTextBuilder(tt.idPrefix, tt.index); got != tt.want {
				t.Errorf("tabPanelIdTextBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabValuesArrayLiteralBuilder(t *testing.T) {
	items := []TabItemSettings{
		{Label: "General", Value: "general"},
		{Label: "Security", Value: "security"},
	}
	if got := tabValuesArrayLiteralBuilder(items); got != `["general","security"]` {
		t.Errorf("tabValuesArrayLiteralBuilder() = %q, want %q", got, `["general","security"]`)
	}
}

func TestTabAnySelectedExpressionBuilder(t *testing.T) {
	items := []TabItemSettings{
		{Label: "General", Value: "general"},
		{Label: "Security", Value: "security"},
	}
	want := `["general","security"].includes(selectedTab)`
	if got := tabAnySelectedExpressionBuilder("selectedTab", items); got != want {
		t.Errorf("tabAnySelectedExpressionBuilder() = %q, want %q", got, want)
	}
}

func TestTabTabindexExpressionBuilder(t *testing.T) {
	tests := []struct {
		name                  string
		selectedExpression    string
		anySelectedExpression string
		isFirst               bool
		want                  string
	}{
		{
			name:                  "non-first tab is focusable only when selected",
			selectedExpression:    "selectedTab === 'security'",
			anySelectedExpression: "['general','security'].includes(selectedTab)",
			isFirst:               false,
			want:                  "selectedTab === 'security' ? '0' : '-1'",
		},
		{
			name:                  "first tab stays focusable when nothing matches",
			selectedExpression:    "selectedTab === 'general'",
			anySelectedExpression: "['general','security'].includes(selectedTab)",
			isFirst:               true,
			want:                  "selectedTab === 'general' || !(['general','security'].includes(selectedTab)) ? '0' : '-1'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabTabindexExpressionBuilder(tt.selectedExpression, tt.anySelectedExpression, tt.isFirst); got != tt.want {
				t.Errorf("tabTabindexExpressionBuilder() = %q, want %q", got, tt.want)
			}
		})
	}
}

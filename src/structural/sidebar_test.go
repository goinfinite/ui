package uiStructural

import "testing"

func TestSidebarIsInlineAttachmentExpressionBuilder(t *testing.T) {
	tests := []struct {
		name                          string
		attachmentModeTwoWayStatePath string
		expected                      string
	}{
		{
			name:                          "EmptyPathFallsBackToTrue",
			attachmentModeTwoWayStatePath: "",
			expected:                      "true",
		},
		{
			name:                          "BoundPathComparesAgainstInline",
			attachmentModeTwoWayStatePath: "attachmentMode",
			expected:                      "attachmentMode === 'inline'",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualExpression := sidebarIsInlineAttachmentExpressionBuilder(
				testCase.attachmentModeTwoWayStatePath,
			)

			if actualExpression != testCase.expected {
				t.Errorf(
					"actualExpression = %q, want %q",
					actualExpression, testCase.expected,
				)
			}
		})
	}
}

func TestSidebarDynamicClassesBuilder(t *testing.T) {
	tests := []struct {
		name              string
		componentSettings SidebarSettings
		attachmentMode    string
		expected          string
	}{
		{
			name:              "NoStatePathsBuildsEmptyObject",
			componentSettings: SidebarSettings{},
			attachmentMode:    SidebarAttachmentModeInline,
			expected:          "{}",
		},
		{
			name: "VisibleAndCollapsedPathsAddTheirToggles",
			componentSettings: SidebarSettings{
				IsVisibleTwoWayStatePath:   "isVisible",
				IsCollapsedTwoWayStatePath: "isCollapsed",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected:       "{'invisible': !isVisible,'!w-16': isCollapsed,}",
		},
		{
			name: "AttachmentPathAddsTheFixedToggle",
			componentSettings: SidebarSettings{
				AttachmentModeTwoWayStatePath: "attachmentMode",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected:       "{'!fixed top-0': attachmentMode === 'fixed',}",
		},
		{
			name: "OffCanvasPathAddsTheAbsoluteToggle",
			componentSettings: SidebarSettings{
				IsOffCanvasTwoWayStatePath: "isOffCanvas",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected:       "{'!absolute top-0 z-50': isOffCanvas,}",
		},
		{
			name: "SidePathAddsOnlyPositionAnchors",
			componentSettings: SidebarSettings{
				SideTwoWayStatePath: "side",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected:       "{'left-0': side === 'left','right-0': side === 'right',}",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualExpression := sidebarDynamicClassesBuilder(
				testCase.componentSettings, testCase.attachmentMode,
			)

			if actualExpression != testCase.expected {
				t.Errorf(
					"actualExpression = %q, want %q",
					actualExpression, testCase.expected,
				)
			}
		})
	}
}

func TestSidebarWrapperDynamicClassesBuilder(t *testing.T) {
	tests := []struct {
		name              string
		componentSettings SidebarSettings
		attachmentMode    string
		expected          string
	}{
		{
			name:              "NoStatePathsBuildsEmptyObject",
			componentSettings: SidebarSettings{},
			attachmentMode:    SidebarAttachmentModeInline,
			expected:          "{}",
		},
		{
			name: "VisibleOffCanvasAndCollapsedPathsAddTheirToggles",
			componentSettings: SidebarSettings{
				IsVisibleTwoWayStatePath:   "isVisible",
				IsOffCanvasTwoWayStatePath: "isOffCanvas",
				IsCollapsedTwoWayStatePath: "isCollapsed",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected: "{'invisible': !isVisible,'!w-0': isOffCanvas," +
				"'!w-16': isCollapsed,}",
		},
		{
			name: "SidePathWithInlineAttachmentAddsFloatClasses",
			componentSettings: SidebarSettings{
				SideTwoWayStatePath: "side",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected: "{'float-left': side === 'left' && true," +
				"'float-right': side === 'right' && true,}",
		},
		{
			name: "SidePathWithBoundAttachmentGatesFloatsOnInline",
			componentSettings: SidebarSettings{
				SideTwoWayStatePath:           "side",
				AttachmentModeTwoWayStatePath: "attachmentMode",
			},
			attachmentMode: SidebarAttachmentModeInline,
			expected: "{'float-left': side === 'left' && attachmentMode === 'inline'," +
				"'float-right': side === 'right' && attachmentMode === 'inline',}",
		},
		{
			name: "SidePathWithFixedAttachmentOmitsFloatClasses",
			componentSettings: SidebarSettings{
				SideTwoWayStatePath: "side",
			},
			attachmentMode: SidebarAttachmentModeFixed,
			expected:       "{}",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actualExpression := sidebarWrapperDynamicClassesBuilder(
				testCase.componentSettings, testCase.attachmentMode,
			)

			if actualExpression != testCase.expected {
				t.Errorf(
					"actualExpression = %q, want %q",
					actualExpression, testCase.expected,
				)
			}
		})
	}
}

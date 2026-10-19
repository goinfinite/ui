package uiControl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestButtonRendersPaddingAndSizeConstraints(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Button(ButtonSettings{
		PaddingSize:              uiToolset.PaddingSizeNone,
		MinWidthClass:            "min-w-7",
		MinHeightClass:           "min-h-7",
		Size:                     ButtonSizeMd,
		IconLeft:                 "ph-binary",
		AriaLabelOneWayStatePath: "isVisible ? 'Hide' : 'Show'",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ButtonRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		"p-0", "min-w-7", "min-h-7",
		`:aria-label="isVisible ? &#39;Hide&#39; : &#39;Show&#39;"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestButtonRendersAriaExpanded(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Button(ButtonSettings{
		AriaExpandedOneWayStatePath: "isOpen",
		AriaLabel:                   "Toggle date picker",
		IconLeft:                    "ph-calendar-blank",
		OnClickFunc:                 "togglePopover()",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("ButtonRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`:aria-expanded="isOpen"`,
		`aria-label="Toggle date picker"`,
		`@click="togglePopover()"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestButtonRendersActionDisplay(t *testing.T) {
	testCases := []struct {
		name               string
		settings           ButtonSettings
		expectedFragments  []string
		forbiddenFragments []string
	}{
		{
			name:     "PlainButtonHasNoActionScope",
			settings: ButtonSettings{Label: "Save", IsSubmit: true, OnClickFunc: "save()"},
			expectedFragments: []string{
				`@click="save()"`,
			},
			forbiddenFragments: []string{
				"buttonAction(", "actionLabelResolver", ":aria-busy", "trackActionPromise",
			},
		},
		{
			name: "TrackedSubmitButtonFollowsItsForm",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
			},
			expectedFragments: []string{
				`class="contents"`,
				`buttonAction({&#34;idle&#34;:&#34;Save&#34;,&#34;pending&#34;:&#34;Save&#34;,&#34;succeeded&#34;:&#34;Save&#34;,&#34;failed&#34;:&#34;Save&#34;},{&#34;idle&#34;:&#34;&#34;,&#34;pending&#34;:&#34;ph-spinner animate-spin&#34;,&#34;succeeded&#34;:&#34;ph-check-circle&#34;,&#34;failed&#34;:&#34;ph-warning-circle&#34;},true)`,
				`:disabled="isActionPending(false)"`,
				`:aria-busy="isActionPending(false)"`,
				`x-text="actionLabelResolver(false)"`,
			},
		},
		{
			name: "TrackedActionButtonFollowsItsClickPromise",
			settings: ButtonSettings{
				Label: "Sync", OnClickFunc: "syncAccounts()", ShouldTrackAction: true,
			},
			expectedFragments: []string{
				`},false)`,
				`@click="trackActionPromise(syncAccounts())"`,
			},
		},
		{
			name: "CustomActionIconsReplaceTheDefaults",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
				PendingIcon: "ph-hourglass", SucceededIcon: "ph-thumbs-up",
				FailedIcon: "ph-skull",
			},
			expectedFragments: []string{
				`&#34;pending&#34;:&#34;ph-hourglass animate-spin&#34;`,
				`&#34;succeeded&#34;:&#34;ph-thumbs-up&#34;`,
				`&#34;failed&#34;:&#34;ph-skull&#34;`,
			},
		},
		{
			name: "PendingIconAnimationConstantReplacesTheSpin",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
				PendingIconAnimation: uiToolset.AnimationPulse,
			},
			expectedFragments: []string{
				`&#34;pending&#34;:&#34;ph-spinner animate-pulse&#34;`,
			},
			forbiddenFragments: []string{
				"animate-spin",
			},
		},
		{
			name: "PendingIconAnimationCustomClassPassesThrough",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
				PendingIconAnimation: "animate-wiggle",
			},
			expectedFragments: []string{
				`&#34;pending&#34;:&#34;ph-spinner animate-wiggle&#34;`,
			},
		},
		{
			name: "PendingIconAnimationNoneDropsTheAnimation",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
				PendingIconAnimation: uiToolset.AnimationNone,
			},
			expectedFragments: []string{
				`&#34;pending&#34;:&#34;ph-spinner&#34;`,
			},
			forbiddenFragments: []string{
				"animate-spin",
			},
		},
		{
			name: "ExternalPendingPathMergesWithDisabledPath",
			settings: ButtonSettings{
				Label:                     "Save",
				PendingLabel:              "Saving",
				IsPendingOneWayStatePath:  "isSaving",
				IsDisabledOneWayStatePath: "!isValid",
			},
			expectedFragments: []string{
				`&#34;pending&#34;:&#34;Saving&#34;`,
				`},false)`,
				`:disabled="(!isValid) || isActionPending((isSaving))"`,
				`x-text="actionLabelResolver((isSaving))"`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := Button(testCase.settings).Render(
				context.Background(), &buffer,
			)
			if renderErr != nil {
				t.Fatalf("ButtonRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			for _, expectedFragment := range testCase.expectedFragments {
				if !strings.Contains(renderedHtml, expectedFragment) {
					t.Errorf("RenderedHtmlMissingFragment: %q in %s", expectedFragment, renderedHtml)
				}
			}
			for _, forbiddenFragment := range testCase.forbiddenFragments {
				if strings.Contains(renderedHtml, forbiddenFragment) {
					t.Errorf("RenderedHtmlHasForbiddenFragment: %q", forbiddenFragment)
				}
			}
		})
	}
}

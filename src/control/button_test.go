package uiControl

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

type buttonRenderCase struct {
	name               string
	settings           ButtonSettings
	expectedFragments  []string
	forbiddenFragments []string
}

func renderButtonHtml(test *testing.T, settings ButtonSettings) string {
	test.Helper()
	var buffer bytes.Buffer
	renderErr := Button(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		test.Fatalf("ButtonRenderFailed: %v", renderErr)
	}
	return buffer.String()
}

func assertButtonHtml(
	test *testing.T,
	renderedHtml string,
	expectedFragments []string,
	forbiddenFragments []string,
) {
	test.Helper()
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			test.Errorf("RenderedHtmlMissingFragment: %q in %s", expectedFragment, renderedHtml)
		}
	}
	for _, forbiddenFragment := range forbiddenFragments {
		if strings.Contains(renderedHtml, forbiddenFragment) {
			test.Errorf("RenderedHtmlHasForbiddenFragment: %q", forbiddenFragment)
		}
	}
}

func runButtonRenderCases(test *testing.T, testCases []buttonRenderCase) {
	test.Helper()
	for _, testCase := range testCases {
		test.Run(testCase.name, func(subtest *testing.T) {
			renderedHtml := renderButtonHtml(subtest, testCase.settings)
			assertButtonHtml(
				subtest, renderedHtml,
				testCase.expectedFragments, testCase.forbiddenFragments,
			)
		})
	}
}

func TestButtonRendersPaddingAndSizeConstraints(test *testing.T) {
	renderedHtml := renderButtonHtml(test, ButtonSettings{
		PaddingSize:              uiToolset.PaddingSizeNone,
		MinWidthClass:            "min-w-7",
		MinHeightClass:           "min-h-7",
		Size:                     ButtonSizeMd,
		IconLeft:                 "ph-binary",
		AriaLabelOneWayStatePath: "isVisible ? 'Hide' : 'Show'",
	})
	assertButtonHtml(test, renderedHtml, []string{
		"p-0", "min-w-7", "min-h-7",
		`:aria-label="isVisible ? &#39;Hide&#39; : &#39;Show&#39;"`,
	}, nil)
}

func TestButtonRendersAriaExpanded(test *testing.T) {
	renderedHtml := renderButtonHtml(test, ButtonSettings{
		AriaExpandedOneWayStatePath: "isOpen",
		AriaLabel:                   "Toggle date picker",
		IconLeft:                    "ph-calendar-blank",
		OnClickFunc:                 "togglePopover()",
	})
	assertButtonHtml(test, renderedHtml, []string{
		`:aria-expanded="isOpen"`,
		`aria-label="Toggle date picker"`,
		`@click="togglePopover()"`,
	}, nil)
}

func TestButtonRendersActionDisplay(test *testing.T) {
	testCases := []buttonRenderCase{
		{
			name:     "PlainButtonHasNoButtonState",
			settings: ButtonSettings{Label: "Save", IsSubmit: true, OnClickFunc: "save()"},
			expectedFragments: []string{
				`@click="save()"`,
			},
			forbiddenFragments: []string{
				"buttonState(", "actionLabelResolver", ":aria-busy", "trackActionPromise",
			},
		},
		{
			name: "TrackedSubmitButtonFollowsItsForm",
			settings: ButtonSettings{
				Label: "Save", IsSubmit: true, ShouldTrackAction: true,
			},
			expectedFragments: []string{
				`class="contents"`,
				`buttonState({&#34;idle&#34;:&#34;&#34;,&#34;pending&#34;:&#34;ph-spinner animate-spin&#34;,&#34;succeeded&#34;:&#34;ph-check-circle&#34;,&#34;failed&#34;:&#34;ph-warning-circle&#34;},true,0)`,
				`:disabled="isActionPending(false)"`,
				`:aria-busy="isActionPending(false)"`,
				`x-text="actionLabelResolver(&#34;Save&#34;, &#34;Save&#34;, &#34;Save&#34;, &#34;Save&#34;, false)"`,
			},
		},
		{
			name: "TrackedActionButtonFollowsItsClickPromise",
			settings: ButtonSettings{
				Label: "Sync", OnClickFunc: "syncAccounts()", ShouldTrackAction: true,
			},
			expectedFragments: []string{
				`},false,0)`,
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
				`x-text="actionLabelResolver(&#34;Save&#34;, &#34;Saving&#34;, &#34;Save&#34;, &#34;Save&#34;, (isSaving))"`,
				`},false,0)`,
				`:disabled="(!isValid) || isActionPending((isSaving))"`,
			},
		},
	}

	runButtonRenderCases(test, testCases)
}

func TestButtonRendersFormAssociation(test *testing.T) {
	testCases := []buttonRenderCase{
		{
			name:     "FormIdRendersTheFormAttribute",
			settings: ButtonSettings{Label: "Save", IsSubmit: true, FormId: "account-form"},
			expectedFragments: []string{
				`type="submit"`,
				`form="account-form"`,
			},
		},
		{
			name:     "EmptyFormIdOmitsTheFormAttribute",
			settings: ButtonSettings{Label: "Save", IsSubmit: true},
			expectedFragments: []string{
				`type="submit"`,
			},
			forbiddenFragments: []string{
				"form=",
			},
		},
	}

	runButtonRenderCases(test, testCases)
}

func TestButtonRendersIconTop(test *testing.T) {
	testCases := []buttonRenderCase{
		{
			name: "IconTopCombinesWithSideIconsAndLabel",
			settings: ButtonSettings{
				Label:     "Upload",
				IconTop:   "ph-cloud-arrow-up",
				IconLeft:  "ph-folder",
				IconRight: "ph-arrow-right",
			},
			expectedFragments: []string{
				"flex flex-col",
				"ph-cloud-arrow-up",
				"ph-folder",
				"ph-arrow-right",
				"flex flex-row items-center justify-center",
				"Upload",
			},
		},
		{
			name:     "IconTopWithoutALabelCentersTheIcon",
			settings: ButtonSettings{IconTop: "ph-plus", AriaLabel: "Add"},
			expectedFragments: []string{
				"flex flex-col",
				"ph-plus",
				`aria-label="Add"`,
			},
			forbiddenFragments: []string{
				"flex flex-row items-center justify-center",
			},
		},
		{
			name: "IconTopStatePathBindsTheIconClass",
			settings: ButtonSettings{
				Label:                  "Status",
				IconTopOneWayStatePath: "statusIconClass",
			},
			expectedFragments: []string{
				`:class="statusIconClass"`,
			},
		},
	}

	runButtonRenderCases(test, testCases)
}

func TestButtonRendersDynamicLabels(test *testing.T) {
	testCases := []buttonRenderCase{
		{
			name:     "StatePathAloneCreatesTheLabel",
			settings: ButtonSettings{LabelOneWayStatePath: "buttonLabel"},
			expectedFragments: []string{
				`<span x-text="buttonLabel">`,
			},
		},
		{
			name: "ActionLabelsBindToOneWayStatePaths",
			settings: ButtonSettings{
				Label:                         "Save",
				LabelOneWayStatePath:          "idleLabel",
				ShouldTrackAction:             true,
				PendingLabel:                  "Saving",
				PendingLabelOneWayStatePath:   "pendingLabel",
				SucceededLabel:                "Saved",
				SucceededLabelOneWayStatePath: "succeededLabel",
				FailedLabel:                   "Save failed",
				FailedLabelOneWayStatePath:    "failedLabel",
			},
			expectedFragments: []string{
				`x-text="actionLabelResolver(idleLabel, pendingLabel, succeededLabel, failedLabel, false)"`,
			},
		},
		{
			name: "CountdownUsesTheDynamicLabel",
			settings: ButtonSettings{
				Label:                "Delete",
				LabelOneWayStatePath: "deleteLabel",
				CountdownSeconds:     15,
			},
			expectedFragments: []string{
				`x-text="isCounting ? (deleteLabel) + &#34; (&#34; + countdownSecondsRemaining + &#34;s)&#34; : deleteLabel"`,
			},
		},
	}

	runButtonRenderCases(test, testCases)
}

func TestButtonMergesActionAndCountdownIntoOneState(test *testing.T) {
	renderedHtml := renderButtonHtml(test, ButtonSettings{
		Label:             "Delete",
		IsSubmit:          true,
		ShouldTrackAction: true,
		CountdownSeconds:  15,
	})
	stateScopeCount := strings.Count(renderedHtml, `x-data="buttonState(`)
	if stateScopeCount != 1 {
		test.Errorf("ButtonStateScopeCountMismatch: got %d", stateScopeCount)
	}
	assertButtonHtml(test, renderedHtml, []string{
		`:disabled="isActionPending(false) || isCounting"`,
		`x-init="startCountdown()"`,
		"actionLabelResolver(",
		"countdownSecondsRemaining",
	}, nil)
}

func TestButtonRendersCountdown(test *testing.T) {
	testCases := []buttonRenderCase{
		{
			name:     "CountdownStartsOnMountAndShowsTheSeconds",
			settings: ButtonSettings{Label: "Delete", CountdownSeconds: 15},
			expectedFragments: []string{
				"buttonState(",
				",15)",
				`x-init="startCountdown()"`,
				`:disabled="isCounting"`,
				"countdownSecondsRemaining",
				"x-text=",
			},
		},
		{
			name: "CountdownRestartsFromTheTriggerPath",
			settings: ButtonSettings{
				Label:                           "Delete",
				CountdownSeconds:                15,
				CountdownTriggerOneWayStatePath: "deleteAttemptCount",
			},
			expectedFragments: []string{
				`x-effect="restartCountdownOnTrigger(deleteAttemptCount)"`,
			},
			forbiddenFragments: []string{
				`x-init="startCountdown()"`,
			},
		},
		{
			name:     "NoCountdownKeepsThePlainButton",
			settings: ButtonSettings{Label: "Delete"},
			forbiddenFragments: []string{
				"buttonState(", "isCounting", "countdownSecondsRemaining",
			},
		},
	}

	runButtonRenderCases(test, testCases)
}

package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestPasswordInputRulesResolver(t *testing.T) {
	testCases := []struct {
		name     string
		rules    PasswordInputRules
		expected passwordInputRulesConfig
	}{
		{
			name:  "empty rules fall back to the defaults",
			rules: PasswordInputRules{},
			expected: passwordInputRulesConfig{
				MinLength:           6,
				MaxLength:           64,
				Length:              16,
				IncludeLowercase:    true,
				IncludeUppercase:    true,
				IncludeNumbers:      true,
				IncludeSpecialChars: true,
			},
		},
		{
			name: "length only keeps every character class",
			rules: PasswordInputRules{
				MinLength:        12,
				MaxLength:        24,
				GenerationLength: 18,
			},
			expected: passwordInputRulesConfig{
				MinLength:           12,
				MaxLength:           24,
				Length:              18,
				IncludeLowercase:    true,
				IncludeUppercase:    true,
				IncludeNumbers:      true,
				IncludeSpecialChars: true,
			},
		},
		{
			name: "explicit character classes are kept",
			rules: PasswordInputRules{
				ShouldIncludeUppercaseChars: true,
				ShouldIncludeNumbers:        true,
			},
			expected: passwordInputRulesConfig{
				MinLength:        6,
				MaxLength:        64,
				Length:           16,
				IncludeUppercase: true,
				IncludeNumbers:   true,
			},
		},
		{
			name: "max length below min length clamps up",
			rules: PasswordInputRules{
				MinLength: 10,
				MaxLength: 4,
			},
			expected: passwordInputRulesConfig{
				MinLength:           10,
				MaxLength:           10,
				Length:              16,
				IncludeLowercase:    true,
				IncludeUppercase:    true,
				IncludeNumbers:      true,
				IncludeSpecialChars: true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualConfig := passwordInputRulesResolver(testCase.rules)
			if actualConfig != testCase.expected {
				t.Errorf(
					"RulesMismatch: got %+v, want %+v",
					actualConfig, testCase.expected,
				)
			}
		})
	}
}

func TestPasswordStrengthCriteriaItemsResolver(t *testing.T) {
	config := passwordInputRulesResolver(PasswordInputRules{})
	criteriaItems := passwordStrengthCriteriaItemsResolver(config)
	if len(criteriaItems) != 5 {
		t.Fatalf("CriteriaLengthMismatch: got %d, want 5", len(criteriaItems))
	}
	expectedLabels := []string{
		"Between 6 and 64 characters",
		"At least 1 number",
		"At least 1 uppercase letter",
		"At least 1 lowercase letter",
		"At least 1 special character",
	}
	expectedStateKeys := []string{
		"isLongEnough",
		"hasNumbers",
		"hasUppercaseChars",
		"hasLowercaseChars",
		"hasSpecialChars",
	}
	for index, criteriaItem := range criteriaItems {
		if criteriaItem.Label != expectedLabels[index] {
			t.Errorf(
				"CriteriaLabelMismatch(%d): got %q, want %q",
				index, criteriaItem.Label, expectedLabels[index],
			)
		}
		if criteriaItem.StateKey != expectedStateKeys[index] {
			t.Errorf(
				"CriteriaStateKeyMismatch(%d): got %q, want %q",
				index, criteriaItem.StateKey, expectedStateKeys[index],
			)
		}
	}

	trimmedConfig := passwordInputRulesResolver(PasswordInputRules{
		ShouldIncludeUppercaseChars: true,
		ShouldIncludeNumbers:        true,
	})
	trimmedItems := passwordStrengthCriteriaItemsResolver(trimmedConfig)
	if len(trimmedItems) != 3 {
		t.Errorf("TrimmedCriteriaLengthMismatch: got %d, want 3", len(trimmedItems))
	}
}

func TestPasswordInputRendersActionsAndStrengthMeter(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := PasswordInput(PasswordInputSettings{
		InputName:                "password",
		Label:                    "Password",
		LabelOneWayStatePath:     "passwordLabel",
		TwoWayStatePath:          "password",
		ShouldShowGenerateButton: true,
		ShouldShowStrengthMeter:  true,
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("PasswordInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{
		"Show password", "Generate random password",
		"passwordStrengthPercentage", "Between 6 and 64 characters",
		"show/hide the password", "generate random password",
		`role="tooltip"`, "Caps Lock is on", `@keydown="updateCapsLockState($event)"`,
		`x-text="passwordLabel"`, `placeholder="Password"`,
	} {
		if !strings.Contains(renderedHtml, expectedText) {
			testingHandle.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

func TestPasswordInputRendersCapsLockWarningSettings(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := PasswordInput(PasswordInputSettings{
		InputName:                          "password",
		Label:                              "Password",
		TwoWayStatePath:                    "password",
		CapsLockWarningText:                "Watch the modifier",
		CapsLockWarningTextCase:            uiToolset.TextCaseUpper,
		CapsLockWarningTextColor:           "red-400",
		CapsLockWarningTextOneWayStatePath: "capsLockWarningText",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("PasswordInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{
		`role="status"`,
		`x-text="capsLockWarningText"`,
		"Watch the modifier",
		"text-red-400",
		"uppercase",
		`@keydown="updateCapsLockState($event)"`,
	} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}
}

func TestPasswordInputHidesCapsLockWarning(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := PasswordInput(PasswordInputSettings{
		InputName:                 "password",
		Label:                     "Password",
		TwoWayStatePath:           "password",
		ShouldHideCapsLockWarning: true,
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("PasswordInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, absentText := range []string{
		`role="status"`,
		"Caps Lock is on",
		`@keydown="updateCapsLockState($event)"`,
	} {
		if strings.Contains(renderedHtml, absentText) {
			t.Errorf("RenderedHtmlHasUnexpected(%q)", absentText)
		}
	}
}

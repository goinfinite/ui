package uiControl

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestButtonIconOnlyWidthRemResolver(t *testing.T) {
	testCases := []struct {
		size     string
		expected float64
	}{
		{size: ButtonSizeXs, expected: 1.5},
		{size: ButtonSizeSm, expected: 1.75},
		{size: ButtonSizeMd, expected: 2.125},
		{size: ButtonSizeLg, expected: 2.5},
		{size: ButtonSizeXl, expected: 3},
		{size: "", expected: 2},
	}
	for _, testCase := range testCases {
		t.Run(testCase.size, func(t *testing.T) {
			actualWidth := ButtonIconOnlyWidthRemResolver(testCase.size)
			if actualWidth != testCase.expected {
				t.Errorf(
					"IconOnlyWidthMismatch(%q): got %v, want %v",
					testCase.size, actualWidth, testCase.expected,
				)
			}
		})
	}
}

func firstKnownRemResolver(
	classAttributeMatch []string, remByClass map[string]float64,
) (float64, bool) {
	if classAttributeMatch == nil {
		return 0, false
	}
	for classToken := range strings.FieldsSeq(classAttributeMatch[1]) {
		remValue, isKnownClass := remByClass[classToken]
		if isKnownClass {
			return remValue, true
		}
	}
	return 0, false
}

func TestButtonIconOnlyWidthMatchesRenderedClasses(t *testing.T) {
	paddingRemByClass := map[string]float64{
		"p-1": 0.25, "p-1.5": 0.375, "p-2": 0.5, "p-2.5": 0.625, "p-3": 0.75,
	}
	iconRemByClass := map[string]float64{
		"text-base": 1, "text-lg": 1.125, "text-xl": 1.25, "text-2xl": 1.5,
	}
	buttonClassPattern := regexp.MustCompile(`<button[^>]*class="([^"]*)"`)
	iconClassPattern := regexp.MustCompile(`<i class="([^"]*)"`)
	sizes := []string{ButtonSizeXs, ButtonSizeSm, ButtonSizeMd, ButtonSizeLg, ButtonSizeXl, ""}
	for _, size := range sizes {
		t.Run(size, func(t *testing.T) {
			var buffer bytes.Buffer
			renderErr := Button(ButtonSettings{
				Size:     size,
				IconLeft: "ph-x",
			}).Render(context.Background(), &buffer)
			if renderErr != nil {
				t.Fatalf("ButtonRenderFailed: %v", renderErr)
			}
			renderedHtml := buffer.String()
			paddingRem, hasPadding := firstKnownRemResolver(
				buttonClassPattern.FindStringSubmatch(renderedHtml),
				paddingRemByClass,
			)
			iconRem, hasIcon := firstKnownRemResolver(
				iconClassPattern.FindStringSubmatch(renderedHtml), iconRemByClass,
			)
			if !hasPadding || !hasIcon {
				t.Fatalf(
					"RenderedButtonMissingSizeClasses: padding=%v icon=%v",
					hasPadding, hasIcon,
				)
			}
			renderedWidthRem := paddingRem*2 + iconRem
			contractWidthRem := ButtonIconOnlyWidthRemResolver(size)
			if contractWidthRem != renderedWidthRem {
				t.Errorf(
					"ContractWidthMismatch(%q): contract %v, rendered %v",
					size, contractWidthRem, renderedWidthRem,
				)
			}
		})
	}
}

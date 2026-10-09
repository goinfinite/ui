package uiToolset

import (
	_ "embed"

	"github.com/a-h/templ"
)

const (
	TooltipPositionTop    string = "top"
	TooltipPositionBottom string = "bottom"
	TooltipPositionLeft   string = "left"
	TooltipPositionRight  string = "right"
)

//go:embed tooltipState.js
var tooltipAlpineState string

var TooltipAlpineStateOnce = templ.NewOnceHandle(
	templ.WithComponent(MinifierTemplateJs(&tooltipAlpineState)),
)

type TooltipSurfaceSettings struct {
	BackgroundColor string
	TextColor       string
	RingColor       string
	RingThickness   string

	// OptionalFields
	MinWidthClass  string
	MaxWidthClass  string
	MinHeightClass string
	MaxHeightClass string
}

func TooltipSurfaceClassesResolver(settings TooltipSurfaceSettings) string {
	surfaceClasses := "fixed z-100 invisible w-fit p-1.5 text-xs rounded-md shadow-md"
	surfaceClasses += " transition-opacity duration-200 ease-out pointer-events-none"
	surfaceClasses += " " + BackgroundColorClassResolver(
		settings.BackgroundColor, "bg-neutral-800/95",
	)
	surfaceClasses += " " + TextColorClassResolver(settings.TextColor, "text-neutral-50")
	surfaceClasses += " " + CompactRingClassResolver(settings.RingColor, settings.RingThickness)
	surfaceClasses += " " + tooltipSizeClassResolver(settings.MaxWidthClass, "max-w-96")
	if sizeConstraintClasses := NonEmptyClassJoiner([]string{
		settings.MinWidthClass, settings.MinHeightClass, settings.MaxHeightClass,
	}); sizeConstraintClasses != "" {
		surfaceClasses += " " + sizeConstraintClasses
	}
	return surfaceClasses
}

func tooltipSizeClassResolver(sizeClass, fallbackClass string) string {
	if sizeClass != "" {
		return sizeClass
	}
	return fallbackClass
}

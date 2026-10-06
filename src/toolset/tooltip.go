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

func TooltipSurfaceClassesResolver(
	backgroundColor, textColor, ringColor, ringThickness string,
) string {
	surfaceClasses := "fixed z-100 invisible w-fit p-1.5 text-xs rounded-md shadow-md"
	surfaceClasses += " transition-opacity duration-150 pointer-events-none"
	surfaceClasses += " " + BackgroundColorClassResolver(
		backgroundColor, "bg-neutral-800/95",
	)
	surfaceClasses += " " + TextColorClassResolver(textColor, "text-neutral-50")
	return surfaceClasses + " " + CompactRingClassResolver(ringColor, ringThickness)
}

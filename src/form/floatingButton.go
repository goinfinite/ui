package uiForm

import uiControl "github.com/goinfinite/ui/src/control"

func floatingButtonSizeResolver(style FloatingButtonStyle) string {
	if style.Size != "" {
		return style.Size
	}
	return uiControl.ButtonSizeXs
}

func floatingActionWidthRemResolver(style FloatingButtonStyle) float64 {
	return uiControl.ButtonIconOnlyWidthRemResolver(
		floatingButtonSizeResolver(style),
	)
}

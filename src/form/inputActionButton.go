package uiForm

import (
	uiControl "github.com/goinfinite/ui/src/control"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

const (
	InputActionButtonStyleBoxed string = "boxed"
	InputActionButtonStylePlain string = "plain"
)

type inputActionButtonStyle struct {
	BackgroundColor      string
	BackgroundColorHover string
	RingColor            string
	RingColorHover       string
}

type inputActionButtonSize struct {
	ButtonSize     string
	MinHeightClass string
	MinWidthClass  string
}

func inputActionButtonStyleResolver(actionButtonStyle string) inputActionButtonStyle {
	if actionButtonStyle == InputActionButtonStylePlain {
		return inputActionButtonStyle{
			BackgroundColor:      "transparent",
			BackgroundColorHover: "neutral-50/5",
		}
	}
	return inputActionButtonStyle{
		BackgroundColor:      "neutral-50/5",
		BackgroundColorHover: "neutral-50/20",
		RingColor:            "neutral-50/5",
		RingColorHover:       "neutral-50/30",
	}
}

func inputActionButtonSizeResolver(size string) inputActionButtonSize {
	switch size {
	case InputFieldSizeXs:
		return inputActionButtonSize{ButtonSize: uiControl.ButtonSizeXs, MinHeightClass: "min-h-5", MinWidthClass: "min-w-5"}
	case InputFieldSizeSm:
		return inputActionButtonSize{ButtonSize: uiControl.ButtonSizeSm, MinHeightClass: "min-h-6", MinWidthClass: "min-w-6"}
	case InputFieldSizeLg:
		return inputActionButtonSize{ButtonSize: uiControl.ButtonSizeLg, MinHeightClass: "min-h-8", MinWidthClass: "min-w-8"}
	case InputFieldSizeXl:
		return inputActionButtonSize{ButtonSize: uiControl.ButtonSizeXl, MinHeightClass: "min-h-9", MinWidthClass: "min-w-9"}
	}
	return inputActionButtonSize{ButtonSize: uiControl.ButtonSizeMd, MinHeightClass: "min-h-7", MinWidthClass: "min-w-7"}
}

func inputActionButtonSettingsResolver(
	size, actionButtonStyle, ariaLabel, icon, onClickFunc, tooltipContent string,
) uiControl.ButtonSettings {
	actionStyle := inputActionButtonStyleResolver(actionButtonStyle)
	actionSize := inputActionButtonSizeResolver(size)
	return uiControl.ButtonSettings{
		AriaLabel:            ariaLabel,
		BackgroundColor:      actionStyle.BackgroundColor,
		BackgroundColorHover: actionStyle.BackgroundColorHover,
		IconLeft:             icon,
		MinHeightClass:       actionSize.MinHeightClass,
		MinWidthClass:        actionSize.MinWidthClass,
		OnClickFunc:          onClickFunc,
		PaddingSize:          uiToolset.PaddingSizeNone,
		RingColor:            actionStyle.RingColor,
		RingColorHover:       actionStyle.RingColorHover,
		RingThickness:        uiToolset.RingThicknessXs,
		Size:                 actionSize.ButtonSize,
		TextColor:            "neutral-100",
		TextColorHover:       "neutral-100",
		TooltipContent:       tooltipContent,
	}
}

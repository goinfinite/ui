package uiControl

type buttonSizeTokens struct {
	Gap        string
	IconRem    float64
	IconSize   string
	Padding    string
	PaddingRem float64
	TextSize   string
}

func buttonSizeTokensResolver(size string) buttonSizeTokens {
	switch size {
	case ButtonSizeXs:
		return buttonSizeTokens{
			Gap:        "gap-0.5",
			IconRem:    1,
			IconSize:   "text-base",
			Padding:    "p-1",
			PaddingRem: 0.25,
			TextSize:   "text-xs",
		}
	case ButtonSizeSm:
		return buttonSizeTokens{
			Gap:        "gap-0.75",
			IconRem:    1,
			IconSize:   "text-base",
			Padding:    "p-1.5",
			PaddingRem: 0.375,
			TextSize:   "text-sm",
		}
	case ButtonSizeMd:
		return buttonSizeTokens{
			Gap:        "gap-1",
			IconRem:    1.125,
			IconSize:   "text-lg",
			Padding:    "p-2",
			PaddingRem: 0.5,
			TextSize:   "text-base",
		}
	case ButtonSizeLg:
		return buttonSizeTokens{
			Gap:        "gap-1.25",
			IconRem:    1.25,
			IconSize:   "text-xl",
			Padding:    "p-2.5",
			PaddingRem: 0.625,
			TextSize:   "text-lg",
		}
	case ButtonSizeXl:
		return buttonSizeTokens{
			Gap:        "gap-1.5",
			IconRem:    1.5,
			IconSize:   "text-2xl",
			Padding:    "p-3",
			PaddingRem: 0.75,
			TextSize:   "text-xl",
		}
	}
	return buttonSizeTokens{
		Gap:        "gap-1",
		IconRem:    1,
		IconSize:   "text-base",
		Padding:    "p-2",
		PaddingRem: 0.5,
		TextSize:   "text-base",
	}
}

func ButtonIconOnlyWidthRemResolver(size string) float64 {
	sizeTokens := buttonSizeTokensResolver(size)
	return sizeTokens.PaddingRem*2 + sizeTokens.IconRem
}

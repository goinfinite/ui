package uiToolset

const (
	BorderRadiusNone string = "none"
	BorderRadiusXs   string = "xs"
	BorderRadiusSm   string = "sm"
	BorderRadiusMd   string = "md"
	BorderRadiusLg   string = "lg"
	BorderRadiusXl   string = "xl"
	BorderRadius2xl  string = "2xl"
	BorderRadius3xl  string = "3xl"
	BorderRadiusFull string = "full"

	ShapeCircular string = "circular"
	ShapeRounded  string = "rounded"
	ShapeSquare   string = "square"

	ShadowSizeNone string = "none"
	ShadowSizeXs   string = "xs"
	ShadowSizeSm   string = "sm"
	ShadowSizeMd   string = "md"
	ShadowSizeLg   string = "lg"
	ShadowSizeXl   string = "xl"

	RingThicknessXs string = "xs"
	RingThicknessSm string = "sm"
	RingThicknessMd string = "md"
	RingThicknessLg string = "lg"
	RingThicknessXl string = "xl"

	PaddingSizeNone string = "none"
	PaddingSizeXs   string = "xs"
	PaddingSizeSm   string = "sm"
	PaddingSizeMd   string = "md"
	PaddingSizeLg   string = "lg"
	PaddingSizeXl   string = "xl"

	GapSizeNone string = "none"
	GapSizeXs   string = "xs"
	GapSizeSm   string = "sm"
	GapSizeMd   string = "md"
	GapSizeLg   string = "lg"
	GapSizeXl   string = "xl"
)

func BackgroundColorClassResolver(backgroundColor, fallbackClass string) string {
	if backgroundColor != "" {
		return "bg-" + backgroundColor
	}
	return fallbackClass
}

func TextColorClassResolver(textColor, fallbackClass string) string {
	if textColor != "" {
		return "text-" + textColor
	}
	return fallbackClass
}

func BorderColorClassResolver(borderColor, fallbackClass string) string {
	if borderColor != "" {
		return "border-" + borderColor
	}
	return fallbackClass
}

func BorderRadiusTokenResolver(borderRadius, fallbackToken string) string {
	switch borderRadius {
	case BorderRadiusNone, BorderRadiusXs, BorderRadiusSm, BorderRadiusMd,
		BorderRadiusLg, BorderRadiusXl, BorderRadius2xl, BorderRadius3xl,
		BorderRadiusFull:
		return borderRadius
	}
	return fallbackToken
}

func BorderRadiusClassResolver(borderRadius, fallbackClass string) string {
	resolvedToken := BorderRadiusTokenResolver(borderRadius, "")
	if resolvedToken == "" {
		return fallbackClass
	}
	return "rounded-" + resolvedToken
}

func ShapeClassResolver(shape, fallbackClass string) string {
	switch shape {
	case ShapeCircular:
		return "rounded-full"
	case ShapeRounded:
		return "rounded"
	case ShapeSquare:
		return "rounded-none"
	}
	return fallbackClass
}

func ShadowClassResolver(shadowSize, fallbackClass string) string {
	switch shadowSize {
	case ShadowSizeNone:
		return "shadow-none"
	case ShadowSizeXs:
		return "shadow-xs"
	case ShadowSizeSm:
		return "shadow-sm"
	case ShadowSizeMd:
		return "shadow-md"
	case ShadowSizeLg:
		return "shadow-lg"
	case ShadowSizeXl:
		return "shadow-xl"
	}
	return fallbackClass
}

func RingThicknessClassResolver(ringThickness, fallbackClass string) string {
	switch ringThickness {
	case RingThicknessXs:
		return "ring-1"
	case RingThicknessSm:
		return "ring-1.5"
	case RingThicknessMd:
		return "ring-2"
	case RingThicknessLg:
		return "ring-2.5"
	case RingThicknessXl:
		return "ring-3"
	}
	return fallbackClass
}

func RingClassResolver(ringColor, ringThickness string) string {
	if ringColor == "" || ringThickness == "" {
		return ""
	}
	return RingThicknessClassResolver(ringThickness, "ring-1") + " ring-" + ringColor
}

func CompactRingClassResolver(ringColor, ringThickness string) string {
	if ringColor == "" {
		return ""
	}
	ringThicknessClass := "ring-1"
	switch ringThickness {
	case RingThicknessXs:
		ringThicknessClass = "ring-0.5"
	case RingThicknessSm:
		ringThicknessClass = "ring-1"
	case RingThicknessMd:
		ringThicknessClass = "ring-1.5"
	case RingThicknessLg:
		ringThicknessClass = "ring-2"
	case RingThicknessXl:
		ringThicknessClass = "ring-2.5"
	}
	return ringThicknessClass + " ring-" + ringColor
}

func PaddingClassResolver(paddingSize, fallbackClass string) string {
	switch paddingSize {
	case PaddingSizeNone:
		return "p-0"
	case PaddingSizeXs:
		return "p-3"
	case PaddingSizeSm:
		return "p-4"
	case PaddingSizeMd:
		return "p-5"
	case PaddingSizeLg:
		return "p-6"
	case PaddingSizeXl:
		return "p-8"
	}
	return fallbackClass
}

func CompactPaddingClassResolver(paddingSize, fallbackClass string) string {
	switch paddingSize {
	case PaddingSizeNone:
		return "p-0"
	case PaddingSizeXs:
		return "p-1"
	case PaddingSizeSm:
		return "p-1.5"
	case PaddingSizeMd:
		return "p-2"
	case PaddingSizeLg:
		return "p-2.5"
	case PaddingSizeXl:
		return "p-3"
	}
	return fallbackClass
}

func GapClassResolver(gapSize, fallbackClass string) string {
	switch gapSize {
	case GapSizeNone:
		return "gap-0"
	case GapSizeXs:
		return "gap-1"
	case GapSizeSm:
		return "gap-2"
	case GapSizeMd:
		return "gap-3"
	case GapSizeLg:
		return "gap-6"
	case GapSizeXl:
		return "gap-8"
	}
	return fallbackClass
}

package uiToolset

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

func BorderRadiusClassResolver(borderRadius, fallbackClass string) string {
	switch borderRadius {
	case "none":
		return "rounded-none"
	case "xs":
		return "rounded-xs"
	case "sm":
		return "rounded-sm"
	case "md":
		return "rounded"
	case "lg":
		return "rounded-lg"
	case "xl":
		return "rounded-xl"
	}
	return fallbackClass
}

func ShadowClassResolver(shadowSize, fallbackClass string) string {
	switch shadowSize {
	case "none":
		return "shadow-none"
	case "xs":
		return "shadow-xs"
	case "sm":
		return "shadow-sm"
	case "md":
		return "shadow-md"
	case "lg":
		return "shadow-lg"
	case "xl":
		return "shadow-xl"
	}
	return fallbackClass
}

func RingClassResolver(ringColor, ringThickness string) string {
	if ringColor == "" || ringThickness == "" {
		return ""
	}
	ringThicknessClass := "ring-1"
	switch ringThickness {
	case "sm":
		ringThicknessClass = "ring-1.5"
	case "md":
		ringThicknessClass = "ring-2"
	case "lg":
		ringThicknessClass = "ring-2.5"
	case "xl":
		ringThicknessClass = "ring-3"
	}
	return ringThicknessClass + " ring-" + ringColor
}

func PaddingClassResolver(paddingSize, fallbackClass string) string {
	switch paddingSize {
	case "none":
		return "p-0"
	case "xs":
		return "p-3"
	case "sm":
		return "p-4"
	case "md":
		return "p-5"
	case "lg":
		return "p-6"
	case "xl":
		return "p-8"
	}
	return fallbackClass
}

func CompactPaddingClassResolver(paddingSize, fallbackClass string) string {
	switch paddingSize {
	case "none":
		return "p-0"
	case "xs":
		return "p-1"
	case "sm":
		return "p-1.5"
	case "md":
		return "p-2"
	case "lg":
		return "p-2.5"
	case "xl":
		return "p-3"
	}
	return fallbackClass
}

func GapClassResolver(gapSize, fallbackClass string) string {
	switch gapSize {
	case "none":
		return "gap-0"
	case "xs":
		return "gap-1"
	case "sm":
		return "gap-2"
	case "md":
		return "gap-3"
	case "lg":
		return "gap-6"
	case "xl":
		return "gap-8"
	}
	return fallbackClass
}

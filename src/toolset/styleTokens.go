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
		return "rounded-md"
	case "lg":
		return "rounded-lg"
	case "xl":
		return "rounded-xl"
	case "full":
		return "rounded-full"
	}
	return fallbackClass
}

func ShapeClassResolver(shape, fallbackClass string) string {
	switch shape {
	case "circular":
		return "rounded-full"
	case "rounded":
		return "rounded"
	case "square":
		return "rounded-none"
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

func RingThicknessClassResolver(ringThickness, fallbackClass string) string {
	switch ringThickness {
	case "xs":
		return "ring-1"
	case "sm":
		return "ring-1.5"
	case "md":
		return "ring-2"
	case "lg":
		return "ring-2.5"
	case "xl":
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
	case "xs":
		ringThicknessClass = "ring-0.5"
	case "sm":
		ringThicknessClass = "ring-1"
	case "md":
		ringThicknessClass = "ring-1.5"
	case "lg":
		ringThicknessClass = "ring-2"
	case "xl":
		ringThicknessClass = "ring-2.5"
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

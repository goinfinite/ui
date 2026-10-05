package uiStructural

type HorizontalAlignment string

const (
	HorizontalAlignmentLeft   HorizontalAlignment = "left"
	HorizontalAlignmentCenter HorizontalAlignment = "center"
	HorizontalAlignmentRight  HorizontalAlignment = "right"
)

func (alignment HorizontalAlignment) justifyClassResolver() string {
	switch alignment {
	case HorizontalAlignmentCenter:
		return "justify-center"
	case HorizontalAlignmentRight:
		return "justify-end"
	}
	return "justify-start"
}

type TextAlignment string

const (
	TextAlignmentLeft   TextAlignment = "left"
	TextAlignmentCenter TextAlignment = "center"
	TextAlignmentRight  TextAlignment = "right"
)

func (alignment TextAlignment) alignmentClassResolver() string {
	switch alignment {
	case TextAlignmentCenter:
		return "text-center"
	case TextAlignmentRight:
		return "text-right"
	}
	return "text-left"
}

func (alignment TextAlignment) justifyClassResolver() string {
	return HorizontalAlignment(alignment).justifyClassResolver()
}

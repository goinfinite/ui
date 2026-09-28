package uiToolset

const (
	TextCaseNone       string = "none"
	TextCaseLower      string = "lower"
	TextCaseUpper      string = "upper"
	TextCaseCapitalize string = "capitalize"
)

func TextCaseClassResolver(textCase string) string {
	switch textCase {
	case TextCaseLower:
		return "lowercase"
	case TextCaseUpper:
		return "uppercase"
	case TextCaseCapitalize:
		return "capitalize"
	}
	return ""
}

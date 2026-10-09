package uiForm

import (
	_ "embed"
	"fmt"
	"sync/atomic"

	"github.com/a-h/templ"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

//go:embed searchableSelectInputState.js
var searchableSelectInputAlpineState string

var searchableSelectInputAlpineStateOnce = templ.NewOnceHandle(
	templ.WithComponent(
		uiToolset.MinifierTemplateJs(&searchableSelectInputAlpineState),
	),
)

const (
	SearchableSelectInputSelectionDisplayText string = "text"
	SearchableSelectInputSelectionDisplayTags string = "tags"
)

type searchableSelectInputSelectionPresentation struct {
	Display           string
	AllowCustomValues bool
}

func searchableSelectInputSelectionPresentationResolver(
	isMultiSelect bool, selectionDisplay string, shouldAllowCustomValues bool,
) searchableSelectInputSelectionPresentation {
	presentation := searchableSelectInputSelectionPresentation{
		Display: SearchableSelectInputSelectionDisplayText,
	}
	if !isMultiSelect {
		return presentation
	}
	if selectionDisplay != SearchableSelectInputSelectionDisplayTags {
		return presentation
	}
	presentation.Display = SearchableSelectInputSelectionDisplayTags
	presentation.AllowCustomValues = shouldAllowCustomValues
	return presentation
}

var searchableSelectInputConfigScriptIdCounter atomic.Uint64

func generateSearchableSelectInputConfigScriptId() string {
	return fmt.Sprintf(
		"searchable-select-input-config-%d",
		searchableSelectInputConfigScriptIdCounter.Add(1),
	)
}

type searchableSelectInputItemConfig struct {
	Label          string `json:"label"`
	Value          string `json:"value"`
	SearchableText string `json:"searchableText"`
}

type searchableSelectInputRemoteConfig struct {
	Url            string `json:"url"`
	QueryParam     string `json:"queryParam"`
	MinQueryLength uint   `json:"minQueryLength"`
	DebounceMs     uint   `json:"debounceMs"`
}

type searchableSelectInputClientConfig struct {
	Items          []searchableSelectInputItemConfig  `json:"items"`
	InitialOptions []searchableSelectInputItemConfig  `json:"initialOptions"`
	IsMultiSelect  bool                               `json:"isMultiSelect"`
	Remote         *searchableSelectInputRemoteConfig `json:"remote,omitempty"`
}

func searchableSelectInputItemsResolver(
	flatOptions []string, labelValueOptions []SelectLabelValueOption,
) []searchableSelectInputItemConfig {
	items := make(
		[]searchableSelectInputItemConfig, 0,
		len(flatOptions)+len(labelValueOptions),
	)
	for _, option := range flatOptions {
		items = append(items, searchableSelectInputItemConfig{
			Label:          option,
			Value:          option,
			SearchableText: option,
		})
	}
	for _, option := range labelValueOptions {
		items = append(items, searchableSelectInputItemConfig{
			Label:          option.Label,
			Value:          option.Value,
			SearchableText: option.Label,
		})
	}
	return items
}

func searchableSelectInputClientConfigResolver(
	isMultiSelect bool,
	flatOptions []string, labelValueOptions, initialOptions []SelectLabelValueOption,
	remoteConfig *searchableSelectInputRemoteConfig,
) searchableSelectInputClientConfig {
	return searchableSelectInputClientConfig{
		Items:          searchableSelectInputItemsResolver(flatOptions, labelValueOptions),
		InitialOptions: searchableSelectInputItemsResolver(nil, initialOptions),
		IsMultiSelect:  isMultiSelect,
		Remote:         remoteConfig,
	}
}

func searchableSelectInputDropdownClassesResolver(
	settings searchableSelectInputShellSettings,
) string {
	dropdownBackgroundClass := "bg-neutral-800/95"
	if settings.DropdownBackgroundColor != "" {
		dropdownBackgroundClass = "bg-" + settings.DropdownBackgroundColor
	}
	dropdownHeightClasses := selectDropdownHeightClassesResolver(
		settings.DropdownMinHeightClass, settings.DropdownMaxHeightClass,
	)
	return "absolute left-0 z-20 w-full overflow-auto rounded-md " +
		"border border-neutral-50/5 shadow-lg hover:border-neutral-50/30 " +
		dropdownHeightClasses + " " + dropdownBackgroundClass
}

func searchableSelectInputHasSelectionExpressionBuilder(statePath string) string {
	return statePath + ".length > 0"
}

func searchableSelectInputSelectedExpressionBuilder(
	statePath, valueExpression string, isMultiSelect bool,
) string {
	if isMultiSelect {
		return statePath + ".includes(" + valueExpression + ")"
	}
	return "String(" + statePath + ") === " + valueExpression
}

func searchableSelectInputSingleSelectExpressionBuilder(
	statePath, valueExpression, labelExpression, onChangeSuffix string,
) string {
	return "selectedLabelCache.set(" + valueExpression + ", " + labelExpression + ")" +
		"; " + statePath + " = " + valueExpression +
		"; userInput = " + labelExpression +
		"; isUserSearching = false" +
		"; closeDropdown()" + onChangeSuffix
}

func searchableSelectInputMultiToggleExpressionBuilder(
	statePath, valueExpression, labelExpression, onChangeSuffix string,
) string {
	return "selectedLabelCache.set(" + valueExpression + ", " + labelExpression + ")" +
		"; " + statePath + " = " + statePath + ".includes(" + valueExpression + ")" +
		" ? " + statePath + ".filter((item) => item !== " + valueExpression + ")" +
		" : [..." + statePath + ", " + valueExpression + "]" +
		"; userInput = ''; isUserSearching = false" + onChangeSuffix
}

func searchableSelectInputClearExpressionBuilder(
	statePath string, isMultiSelect bool, onChangeSuffix string,
) string {
	emptyValue := "''"
	if isMultiSelect {
		emptyValue = "[]"
	}
	return statePath + " = " + emptyValue + "; userInput = ''; isUserSearching = false" + onChangeSuffix
}

func searchableSelectInputCustomValueAddExpressionBuilder(
	statePath, onChangeSuffix string,
) string {
	return "if (userInput.trim() !== '') { " + statePath + " = " +
		statePath + ".includes(userInput.trim()) ? " + statePath +
		" : [..." + statePath + ", userInput.trim()] } userInput = ''" +
		onChangeSuffix
}

func searchableSelectInputTagRemoveExpressionBuilder(
	statePath, valueExpression, onChangeSuffix string,
) string {
	return statePath + " = " + statePath + ".filter((item) => item !== " +
		valueExpression + ")" + onChangeSuffix
}

func searchableSelectInputTagBackspaceExpressionBuilder(
	statePath, onChangeSuffix string,
) string {
	return "if (userInput === '' && " + statePath + ".length > 0) { " +
		statePath + " = " + statePath + ".slice(0, -1) }" + onChangeSuffix
}

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
	Items             []searchableSelectInputItemConfig  `json:"items"`
	IsMultiSelect     bool                               `json:"isMultiSelect"`
	AllowCustomValues bool                               `json:"allowCustomValues"`
	Remote            *searchableSelectInputRemoteConfig `json:"remote,omitempty"`
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
	isMultiSelect, allowCustomValues bool,
	flatOptions []string, labelValueOptions []SelectLabelValueOption,
	remoteConfig *searchableSelectInputRemoteConfig,
) searchableSelectInputClientConfig {
	return searchableSelectInputClientConfig{
		Items:             searchableSelectInputItemsResolver(flatOptions, labelValueOptions),
		IsMultiSelect:     isMultiSelect,
		AllowCustomValues: allowCustomValues,
		Remote:            remoteConfig,
	}
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
	return "selectedLabelCache[" + valueExpression + "] = " + labelExpression +
		"; " + statePath + " = " + valueExpression +
		"; userInput = " + labelExpression +
		"; closeDropdown()" + onChangeSuffix
}

func searchableSelectInputMultiToggleExpressionBuilder(
	statePath, valueExpression, labelExpression, onChangeSuffix string,
) string {
	return "selectedLabelCache[" + valueExpression + "] = " + labelExpression +
		"; " + statePath + " = " + statePath + ".includes(" + valueExpression + ")" +
		" ? " + statePath + ".filter((item) => item !== " + valueExpression + ")" +
		" : [..." + statePath + ", " + valueExpression + "]" +
		"; userInput = ''" + onChangeSuffix
}

func searchableSelectInputClearExpressionBuilder(
	statePath string, isMultiSelect bool, onChangeSuffix string,
) string {
	emptyValue := "''"
	if isMultiSelect {
		emptyValue = "[]"
	}
	return statePath + " = " + emptyValue + "; userInput = ''" + onChangeSuffix
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

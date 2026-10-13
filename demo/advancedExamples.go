package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	uiDisplay "github.com/goinfinite/ui/src/display"
)

const advancedExamplesAssetDirectory = "assets/advancedExamples"

type advancedExamplePanel struct {
	sectionId string
	panelKey  string
	content   templ.Component
}

type advancedExampleRegistry struct {
	panels []advancedExamplePanel
}

type advancedExampleRegistryContextKey struct{}

func (registry *advancedExampleRegistry) contextAttacher(
	parentContext context.Context,
) context.Context {
	return context.WithValue(
		parentContext, advancedExampleRegistryContextKey{}, registry,
	)
}

func advancedExampleRegistryFinder(
	renderContext context.Context,
) (*advancedExampleRegistry, bool) {
	registry, hasRegistry := renderContext.
		Value(advancedExampleRegistryContextKey{}).(*advancedExampleRegistry)
	return registry, hasRegistry
}

func advancedExampleAssetPathBuilder(sectionId string) string {
	return advancedExamplesAssetDirectory + "/" + sectionId + ".json"
}

func advancedExamplePanelKeyBuilder(panelTitle string) string {
	isKeySeparator := func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsDigit(character)
	}
	titleWords := strings.FieldsFunc(strings.ToLower(panelTitle), isKeySeparator)
	return strings.Join(titleWords, "-")
}

func (registry *advancedExampleRegistry) registerLazyPanels(
	sectionId string, items []uiDisplay.AccordionItemSettings,
) []uiDisplay.AccordionItemSettings {
	assetPath := advancedExampleAssetPathBuilder(sectionId)
	placeholderItems := make([]uiDisplay.AccordionItemSettings, 0, len(items))
	for _, item := range items {
		panelKey := advancedExamplePanelKeyBuilder(item.Title)
		registry.panels = append(registry.panels, advancedExamplePanel{
			sectionId: sectionId,
			panelKey:  panelKey,
			content:   item.Content,
		})
		item.Content = DemoAdvancedExamplePlaceholder(assetPath, panelKey)
		placeholderItems = append(placeholderItems, item)
	}
	return placeholderItems
}

func (registry *advancedExampleRegistry) renderPanelsBySection() (
	map[string]map[string]string, error,
) {
	panelsBySection := map[string]map[string]string{}
	for _, panel := range registry.panels {
		sectionPanels, hasSection := panelsBySection[panel.sectionId]
		if !hasSection {
			sectionPanels = map[string]string{}
			panelsBySection[panel.sectionId] = sectionPanels
		}
		_, isDuplicateKey := sectionPanels[panel.panelKey]
		if isDuplicateKey {
			return nil, errors.New(
				"DuplicateAdvancedExamplePanelKey (" + panel.sectionId + "): " +
					panel.panelKey,
			)
		}

		panelBody := &bytes.Buffer{}
		onceHandleFreshContext := context.Background()
		renderErr := panel.content.Render(onceHandleFreshContext, panelBody)
		if renderErr != nil {
			return nil, errors.New(
				"RenderAdvancedExamplePanelFailed (" + panel.sectionId + "/" +
					panel.panelKey + "): " + renderErr.Error(),
			)
		}
		sectionPanels[panel.panelKey] = panelBody.String()
	}
	return panelsBySection, nil
}

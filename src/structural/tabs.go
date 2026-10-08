package uiStructural

import (
	"strconv"
	"strings"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

var tabIdPrefixGenerator = uiToolset.NewComponentIdPrefixGenerator("ui-tab")

func tabOrientationResolver(orientation string) string {
	if orientation == TabOrientationVertical {
		return TabOrientationVertical
	}
	return TabOrientationHorizontal
}

func tabSideResolver(side string) string {
	if side == TabSideRight {
		return TabSideRight
	}
	return TabSideLeft
}

func tabEdgeClassResolver(
	borderRadius string, isVertical, isRight bool,
) string {
	radius := uiToolset.BorderRadiusTokenResolver(
		borderRadius, uiToolset.BorderRadiusMd,
	)
	if isVertical {
		if isRight {
			return "rounded-r-" + radius +
				" border-l-2 border-transparent -ml-px"
		}
		return "rounded-l-" + radius +
			" border-r-2 border-transparent -mr-px"
	}
	return "rounded-t-" + radius +
		" border-b-2 border-transparent -mb-px"
}

func tabPercentClassResolver(classPrefix string, percent uint) string {
	if percent == 0 {
		return ""
	}
	if percent >= 100 {
		return classPrefix + "-full"
	}
	return classPrefix + "-[" + strconv.Itoa(int(percent)) + "%]"
}

func tabConstraintClassesResolver(componentSettings TabSettings) string {
	constraintClasses := []string{
		tabPercentClassResolver("min-w", componentSettings.MinWidthPercent),
		tabPercentClassResolver("max-w", componentSettings.MaxWidthPercent),
		tabPercentClassResolver("min-h", componentSettings.MinHeightPercent),
		tabPercentClassResolver("max-h", componentSettings.MaxHeightPercent),
		componentSettings.MinWidthClass,
		componentSettings.MaxWidthClass,
		componentSettings.MinHeightClass,
		componentSettings.MaxHeightClass,
	}
	nonEmptyClasses := []string{}
	for _, constraintClass := range constraintClasses {
		if constraintClass != "" {
			nonEmptyClasses = append(nonEmptyClasses, constraintClass)
		}
	}
	return strings.Join(nonEmptyClasses, " ")
}

func tabAlignmentClassResolver(alignment string) string {
	switch alignment {
	case TabAlignmentCenter:
		return "self-center"
	case TabAlignmentBottom:
		return "self-end"
	}
	return "self-start"
}

func tabBadgeClassesResolver(item TabItemSettings) string {
	badgeClasses := "px-1.5 py-0.5 text-xs font-bold"
	badgeClasses += " rounded-" + uiToolset.BorderRadiusTokenResolver(
		item.BadgeBorderRadius, uiToolset.BorderRadiusMd,
	)
	badgeClasses += " " + uiToolset.BackgroundColorClassResolver(item.BadgeBackgroundColor, "bg-neutral-50/10")
	badgeClasses += " " + uiToolset.TextColorClassResolver(item.BadgeTextColor, "text-neutral-50/80")
	badgeRingClasses := uiToolset.RingClassResolver(item.BadgeRingColor, item.BadgeRingThickness)
	if badgeRingClasses != "" {
		badgeClasses += " " + badgeRingClasses
	}
	return badgeClasses
}

func tabAriaLabelResolver(ariaLabel string) string {
	if ariaLabel != "" {
		return ariaLabel
	}
	return "Tabs"
}

func tabSelectedExpressionBuilder(selectedValuePath, itemValue string) string {
	return selectedValuePath + " === " + strconv.Quote(itemValue)
}

func tabValuesArrayLiteralBuilder(items []TabItemSettings) string {
	itemValues := make([]string, 0, len(items))
	for _, item := range items {
		itemValues = append(itemValues, strconv.Quote(item.Value))
	}
	return "[" + strings.Join(itemValues, ",") + "]"
}

func tabAnySelectedExpressionBuilder(
	selectedValuePath string, items []TabItemSettings,
) string {
	return tabValuesArrayLiteralBuilder(items) + ".includes(" + selectedValuePath + ")"
}

func tabTabindexExpressionBuilder(
	selectedExpression, anySelectedExpression string, isFirst bool,
) string {
	if isFirst {
		return selectedExpression + " || !(" + anySelectedExpression +
			") ? '0' : '-1'"
	}
	return selectedExpression + " ? '0' : '-1'"
}

func tabClickExpressionBuilder(
	selectedValuePath, itemValue string, isUrlHashSynced bool,
) string {
	quotedValue := strconv.Quote(itemValue)
	clickExpression := selectedValuePath + " = " + quotedValue
	if isUrlHashSynced {
		clickExpression += "; location.hash = " + quotedValue
	}
	return clickExpression
}

func tabUrlHashSyncExpressionBuilder(
	selectedValuePath string, items []TabItemSettings,
) string {
	return selectedValuePath + " = " + tabValuesArrayLiteralBuilder(items) +
		".includes(location.hash.slice(1)) ? location.hash.slice(1) : " +
		selectedValuePath
}

func tabArrowKeyExpressionBuilder(
	selectedValuePath, siblingDirection string,
) string {
	sibling := "$el." + siblingDirection + "ElementSibling"
	return "if (" + sibling + ") { " + sibling + ".focus(); " +
		selectedValuePath + " = " + sibling + ".dataset.tabValue }"
}

func tabEdgeKeyExpressionBuilder(
	selectedValuePath, edgeProperty string,
) string {
	edge := "$el.parentElement." + edgeProperty
	return edge + ".focus(); " + selectedValuePath + " = " + edge +
		".dataset.tabValue"
}

func tabIdTextBuilder(idPrefix string, index int) string {
	return idPrefix + "-" + strconv.Itoa(index)
}

func tabPanelIdTextBuilder(idPrefix string, index int) string {
	return idPrefix + "-panel-" + strconv.Itoa(index)
}

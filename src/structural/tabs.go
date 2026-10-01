package uiStructural

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

var tabIdCounter atomic.Uint64

func generateTabIdPrefix() string {
	return fmt.Sprintf("ui-tab-%d", tabIdCounter.Add(1))
}

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

func tabBorderRadiusTokenResolver(borderRadius string) string {
	switch borderRadius {
	case TabBorderRadiusNone, TabBorderRadiusXs, TabBorderRadiusSm,
		TabBorderRadiusMd, TabBorderRadiusLg, TabBorderRadiusXl:
		return borderRadius
	}
	return TabBorderRadiusMd
}

func tabEdgeClassResolver(
	borderRadius string, isVertical, isRight bool,
) string {
	radius := tabBorderRadiusTokenResolver(borderRadius)
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

func tabTextColorClassResolver(textColor string) string {
	if textColor != "" {
		return "text-" + textColor
	}
	return "text-neutral-50/80"
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

func tabPaddingClassResolver(paddingSize string) string {
	switch paddingSize {
	case TabPaddingSizeNone:
		return "p-0"
	case TabPaddingSizeXs:
		return "p-3"
	case TabPaddingSizeSm:
		return "p-4"
	case TabPaddingSizeMd:
		return "p-5"
	case TabPaddingSizeLg:
		return "p-6"
	case TabPaddingSizeXl:
		return "p-8"
	}
	return ""
}

func tabGapClassResolver(gapSize string) string {
	switch gapSize {
	case TabGapSizeNone:
		return "gap-0"
	case TabGapSizeXs:
		return "gap-1"
	case TabGapSizeSm:
		return "gap-2"
	case TabGapSizeMd:
		return "gap-3"
	case TabGapSizeLg:
		return "gap-6"
	case TabGapSizeXl:
		return "gap-8"
	}
	return "gap-1"
}

func tabRingClassResolver(ringColor, ringThickness string) string {
	if ringColor == "" || ringThickness == "" {
		return ""
	}
	ringThicknessClass := "ring-1"
	switch ringThickness {
	case TabRingThicknessSm:
		ringThicknessClass = "ring-1.5"
	case TabRingThicknessMd:
		ringThicknessClass = "ring-2"
	case TabRingThicknessLg:
		ringThicknessClass = "ring-2.5"
	case TabRingThicknessXl:
		ringThicknessClass = "ring-3"
	}
	return ringThicknessClass + " ring-" + ringColor
}

func tabShadowClassResolver(shadowSize string) string {
	switch shadowSize {
	case TabShadowSizeNone:
		return "shadow-none"
	case TabShadowSizeXs:
		return "shadow-xs"
	case TabShadowSizeSm:
		return "shadow-sm"
	case TabShadowSizeMd:
		return "shadow-md"
	case TabShadowSizeLg:
		return "shadow-lg"
	case TabShadowSizeXl:
		return "shadow-xl"
	}
	return ""
}

func tabBadgeBackgroundColorClassResolver(backgroundColor string) string {
	if backgroundColor != "" {
		return "bg-" + backgroundColor
	}
	return "bg-neutral-50/10"
}

func tabBadgeTextColorClassResolver(textColor string) string {
	if textColor != "" {
		return "text-" + textColor
	}
	return "text-neutral-50/80"
}

func tabBadgeClassesResolver(item TabItemSettings) string {
	badgeClasses := "px-1.5 py-0.5 text-xs font-bold"
	badgeClasses += " rounded-" + tabBorderRadiusTokenResolver(item.BadgeBorderRadius)
	badgeClasses += " " + tabBadgeBackgroundColorClassResolver(item.BadgeBackgroundColor)
	badgeClasses += " " + tabBadgeTextColorClassResolver(item.BadgeTextColor)
	badgeRingClasses := tabRingClassResolver(item.BadgeRingColor, item.BadgeRingThickness)
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

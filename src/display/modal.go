package uiDisplay

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	modalContentHeightClass          string = "h-auto"
	modalContentHeightCeilingPercent int    = 85
)

type modalSizeClassSettings struct {
	Width   string
	Height  string
	Padding string
}

var modalSizeClassesBySize = map[string]modalSizeClassSettings{
	ModalSizeXs:   {Width: "w-[40%]", Height: "h-[40%]", Padding: "p-3"},
	ModalSizeSm:   {Width: "w-[50%]", Height: "h-[50%]", Padding: "p-3.5"},
	ModalSizeMd:   {Width: "w-[60%]", Height: "h-[60%]", Padding: "p-4"},
	ModalSizeLg:   {Width: "w-[70%]", Height: "h-[70%]", Padding: "p-4.5"},
	ModalSizeXl:   {Width: "w-[80%]", Height: "h-[80%]", Padding: "p-5"},
	ModalSizeXxl:  {Width: "w-[90%]", Height: "h-[90%]", Padding: "p-5.5"},
	ModalSizeFull: {Width: "w-full", Height: "h-full", Padding: "p-5.5"},
}

var modalSizeOrder = []string{
	ModalSizeXs, ModalSizeSm, ModalSizeMd,
	ModalSizeLg, ModalSizeXl, ModalSizeXxl, ModalSizeFull,
}

var modalDefaultPossibleSizes = []string{
	ModalSizeMd, ModalSizeLg, ModalSizeXl, ModalSizeXxl,
}

func modalPossibleSizesResolver(possibleSizes []string) []string {
	if len(possibleSizes) == 0 {
		return modalDefaultPossibleSizes
	}
	resolvedSizes := []string{}
	for _, possibleSize := range possibleSizes {
		if _, isKnownSize := modalSizeClassesBySize[possibleSize]; !isKnownSize {
			continue
		}
		resolvedSizes = append(resolvedSizes, possibleSize)
	}
	if len(resolvedSizes) == 0 {
		return modalDefaultPossibleSizes
	}
	return resolvedSizes
}

func modalInitialSizeResolver(initialSize string, possibleSizes []string) string {
	if initialSize == "" {
		initialSize = ModalSizeMd
	}
	if slices.Contains(possibleSizes, initialSize) {
		return initialSize
	}
	return possibleSizes[0]
}

func modalReachableSizesResolver(initialSize string, possibleSizes []string) []string {
	for index, possibleSize := range possibleSizes {
		if possibleSize == initialSize {
			return possibleSizes[index:]
		}
	}
	return possibleSizes
}

func modalPercentageClassBuilder(classPrefix string, percent int) string {
	if percent <= 0 {
		return ""
	}
	if percent >= 100 {
		return classPrefix + "-full"
	}
	return classPrefix + "-[" + strconv.Itoa(percent) + "%]"
}

func modalSizeClassesResolver(
	initialSize string, widthPercent, heightPercent int, shouldHeightMatchContent bool,
) string {
	sizeClasses, isKnownSize := modalSizeClassesBySize[initialSize]
	if !isKnownSize {
		sizeClasses = modalSizeClassesBySize[ModalSizeMd]
	}
	widthClasses := sizeClasses.Width
	if widthPercent > 0 {
		widthClasses = modalPercentageClassBuilder("w", widthPercent)
	}
	heightClasses := sizeClasses.Height
	if heightPercent > 0 {
		heightClasses = modalPercentageClassBuilder("h", heightPercent)
	}
	if shouldHeightMatchContent {
		heightClasses = modalContentHeightClass
	}
	return widthClasses + " " + heightClasses + " " + sizeClasses.Padding
}

func modalSizeConstraintClassesResolver(
	minWidthPercent, maxWidthPercent, minHeightPercent, maxHeightPercent int,
	shouldHeightMatchContent bool,
) string {
	heightCeilingPercent := maxHeightPercent
	if shouldHeightMatchContent && heightCeilingPercent == 0 {
		heightCeilingPercent = modalContentHeightCeilingPercent
	}
	constraintClasses := []string{
		modalPercentageClassBuilder("min-w", minWidthPercent),
		modalPercentageClassBuilder("max-w", maxWidthPercent),
		modalPercentageClassBuilder("min-h", minHeightPercent),
		modalPercentageClassBuilder("max-h", heightCeilingPercent),
	}
	nonEmptyClasses := []string{}
	for _, constraintClass := range constraintClasses {
		if constraintClass != "" {
			nonEmptyClasses = append(nonEmptyClasses, constraintClass)
		}
	}
	if len(nonEmptyClasses) == 0 {
		return ""
	}
	return strings.Join(nonEmptyClasses, " ")
}

func modalSizeClassMapExpressionBuilder(
	sizePath string, widthPercent, heightPercent int, shouldHeightMatchContent bool,
) string {
	var mapExpression strings.Builder
	mapExpression.WriteString("{")
	for index, size := range modalSizeOrder {
		if index > 0 {
			mapExpression.WriteString(",")
		}
		fmt.Fprintf(
			&mapExpression, " '%s': %s === '%s'",
			modalSizeClassesResolver(size, widthPercent, heightPercent, shouldHeightMatchContent),
			sizePath, size,
		)
	}
	mapExpression.WriteString(" }")
	return mapExpression.String()
}

func modalIsResizableResolver(
	isUnresizable bool, widthPercent, heightPercent int, reachableSizes []string,
) bool {
	if isUnresizable {
		return false
	}
	if widthPercent > 0 && heightPercent > 0 {
		return false
	}
	return len(reachableSizes) > 1
}

func modalCanEnlargeExpressionBuilder(sizePath string, reachableSizes []string) string {
	if len(reachableSizes) < 2 {
		return ""
	}
	largestSize := reachableSizes[len(reachableSizes)-1]
	return sizePath + " !== '" + largestSize + "'"
}

func modalCanReduceExpressionBuilder(sizePath string, reachableSizes []string) string {
	if len(reachableSizes) < 2 {
		return ""
	}
	smallestSize := reachableSizes[0]
	return sizePath + " !== '" + smallestSize + "'"
}

func modalEnlargeExpressionBuilder(sizePath string, reachableSizes []string) string {
	largestSize := reachableSizes[len(reachableSizes)-1]
	var expression strings.Builder
	fmt.Fprintf(&expression, "%s = ", sizePath)
	for index := 0; index < len(reachableSizes)-1; index++ {
		fmt.Fprintf(
			&expression, "%s === '%s' ? '%s' : ",
			sizePath, reachableSizes[index], reachableSizes[index+1],
		)
	}
	fmt.Fprintf(&expression, "'%s'", largestSize)
	return expression.String()
}

func modalReduceExpressionBuilder(sizePath string, reachableSizes []string) string {
	smallestSize := reachableSizes[0]
	var expression strings.Builder
	fmt.Fprintf(&expression, "%s = ", sizePath)
	for index := len(reachableSizes) - 1; index > 0; index-- {
		fmt.Fprintf(
			&expression, "%s === '%s' ? '%s' : ",
			sizePath, reachableSizes[index], reachableSizes[index-1],
		)
	}
	fmt.Fprintf(&expression, "'%s'", smallestSize)
	return expression.String()
}

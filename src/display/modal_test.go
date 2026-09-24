package uiDisplay

import (
	"reflect"
	"testing"
)

func TestModalPossibleSizesResolver(t *testing.T) {
	tests := []struct {
		name          string
		possibleSizes []string
		expected      []string
	}{
		{"EmptyInputFallsBackToDefaultRange", nil, modalDefaultPossibleSizes},
		{"EmptySliceFallsBackToDefaultRange", []string{}, modalDefaultPossibleSizes},
		{"UnknownOnlyFallsBackToDefaultRange", []string{"bogus"}, modalDefaultPossibleSizes},
		{
			"KnownSubsetIsKeptInOrder",
			[]string{ModalSizeXl, ModalSizeXxl, ModalSizeFull},
			[]string{ModalSizeXl, ModalSizeXxl, ModalSizeFull},
		},
		{
			"UnknownEntriesAreDropped",
			[]string{ModalSizeSm, "bogus", ModalSizeMd},
			[]string{ModalSizeSm, ModalSizeMd},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalPossibleSizesResolver(testCase.possibleSizes)

			if !reflect.DeepEqual(resolved, testCase.expected) {
				t.Errorf("resolved = %v, want %v", resolved, testCase.expected)
			}
		})
	}
}

func TestModalInitialSizeResolver(t *testing.T) {
	tests := []struct {
		name          string
		initialSize   string
		possibleSizes []string
		expected      string
	}{
		{"EmptyFallsBackToMedium", "", modalSizeOrder, ModalSizeMd},
		{"KnownSizeIsKept", ModalSizeLg, modalSizeOrder, ModalSizeLg},
		{"UnknownFallsBackToFirst", "bogus", modalSizeOrder, ModalSizeXs},
		{"EmptyUsesFirstOfSubset", "", []string{ModalSizeXl, ModalSizeXxl}, ModalSizeXl},
		{"UnknownUsesFirstOfSubset", ModalSizeMd, []string{ModalSizeXl, ModalSizeXxl}, ModalSizeXl},
		{"KnownSizeInSubsetIsKept", ModalSizeXxl, []string{ModalSizeXl, ModalSizeXxl}, ModalSizeXxl},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalInitialSizeResolver(testCase.initialSize, testCase.possibleSizes)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalReachableSizesResolver(t *testing.T) {
	tests := []struct {
		name          string
		initialSize   string
		possibleSizes []string
		expected      []string
	}{
		{
			"EntriesBelowInitialAreDropped",
			ModalSizeMd,
			modalSizeOrder,
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl, ModalSizeXxl, ModalSizeFull},
		},
		{"InitialAtFirstKeepsEveryEntry", ModalSizeXs, modalSizeOrder, modalSizeOrder},
		{"InitialAtLastKeepsOnlyLast", ModalSizeFull, modalSizeOrder, []string{ModalSizeFull}},
		{
			"SubsetStartsAtInitial",
			ModalSizeMd,
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
		},
		{
			"SubsetDropsEntriesBeforeInitial",
			ModalSizeLg,
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			[]string{ModalSizeLg, ModalSizeXl},
		},
		{
			"UnknownInitialKeepsEveryEntry",
			"bogus",
			[]string{ModalSizeMd, ModalSizeLg},
			[]string{ModalSizeMd, ModalSizeLg},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalReachableSizesResolver(testCase.initialSize, testCase.possibleSizes)

			if !reflect.DeepEqual(resolved, testCase.expected) {
				t.Errorf("resolved = %v, want %v", resolved, testCase.expected)
			}
		})
	}
}

func TestModalPercentageClassBuilder(t *testing.T) {
	tests := []struct {
		name        string
		classPrefix string
		percent     int
		expected    string
	}{
		{"ZeroPercentBuildsNothing", "w", 0, ""},
		{"NegativePercentBuildsNothing", "w", -5, ""},
		{"HalfWidthBuildsArbitraryClass", "w", 50, "w-[50%]"},
		{"HeightPercentBuildsArbitraryClass", "h", 45, "h-[45%]"},
		{"HundredPercentBuildsFullClass", "w", 100, "w-full"},
		{"OverHundredPercentClampsToFull", "h", 120, "h-full"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalPercentageClassBuilder(testCase.classPrefix, testCase.percent)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalSizeClassesResolver(t *testing.T) {
	tests := []struct {
		name                 string
		initialSize          string
		widthPercent         int
		heightPercent        int
		isHeightContentSized bool
		expected             string
	}{
		{"MediumUsesItsOwnBox", ModalSizeMd, 0, 0, false, "w-[60%] h-[60%] p-4"},
		{"ExtraSmallUsesItsOwnBox", ModalSizeXs, 0, 0, false, "w-[40%] h-[40%] p-3"},
		{"UnknownSizeFallsBackToMedium", "bogus", 0, 0, false, "w-[60%] h-[60%] p-4"},
		{"WidthPercentReplacesOnlyWidth", ModalSizeMd, 70, 0, false, "w-[70%] h-[60%] p-4"},
		{"HeightPercentReplacesOnlyHeight", ModalSizeMd, 0, 45, false, "w-[60%] h-[45%] p-4"},
		{
			"BothPercentsReplaceBothAxes",
			ModalSizeXxl, 95, 85, false, "w-[95%] h-[85%] p-5.5",
		},
		{"ContentHeightReplacesHeight", ModalSizeMd, 0, 0, true, "w-[60%] h-auto max-h-[85%] p-4"},
		{
			"ContentHeightWinsOverHeightPercent",
			ModalSizeMd, 50, 45, true, "w-[50%] h-auto max-h-[85%] p-4",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalSizeClassesResolver(
				testCase.initialSize,
				testCase.widthPercent,
				testCase.heightPercent,
				testCase.isHeightContentSized,
			)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalSizeClassMapExpressionBuilder(t *testing.T) {
	resolved := modalSizeClassMapExpressionBuilder("modalSize", 70, 45, false)
	expectedPrefix := "{ 'w-[70%] h-[45%] p-3': modalSize === 'xs',"
	if len(resolved) < len(expectedPrefix) || resolved[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("resolved = %q, want prefix %q", resolved, expectedPrefix)
	}
	expectedSuffix := " 'w-[70%] h-[45%] p-5.5': modalSize === 'full' }"
	if len(resolved) < len(expectedSuffix) ||
		resolved[len(resolved)-len(expectedSuffix):] != expectedSuffix {
		t.Errorf("resolved = %q, want suffix %q", resolved, expectedSuffix)
	}
}

func TestModalIsResizableResolver(t *testing.T) {
	tests := []struct {
		name           string
		isUnresizable  bool
		widthPercent   int
		heightPercent  int
		reachableSizes []string
		expected       bool
	}{
		{"DefaultIsResizable", false, 0, 0, modalSizeOrder, true},
		{"UnresizableWins", true, 0, 0, modalSizeOrder, false},
		{"BothPinnedRetiresResize", false, 95, 85, modalSizeOrder, false},
		{"WidthOnlyStillResizes", false, 95, 0, modalSizeOrder, true},
		{"HeightOnlyStillResizes", false, 0, 85, modalSizeOrder, true},
		{"SingleReachableSizeRetiresResize", false, 0, 0, []string{ModalSizeXl}, false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalIsResizableResolver(
				testCase.isUnresizable,
				testCase.widthPercent,
				testCase.heightPercent,
				testCase.reachableSizes,
			)

			if resolved != testCase.expected {
				t.Errorf("resolved = %v, want %v", resolved, testCase.expected)
			}
		})
	}
}

func TestModalCanEnlargeExpressionBuilder(t *testing.T) {
	tests := []struct {
		name           string
		sizePath       string
		reachableSizes []string
		expected       string
	}{
		{
			"NotAtLargestCanEnlarge",
			"modalSize",
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			"modalSize !== 'xl'",
		},
		{"SingleSizeCannotEnlarge", "modalSize", []string{ModalSizeXl}, ""},
		{"EmptyCannotEnlarge", "modalSize", nil, ""},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalCanEnlargeExpressionBuilder(testCase.sizePath, testCase.reachableSizes)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalCanReduceExpressionBuilder(t *testing.T) {
	tests := []struct {
		name           string
		sizePath       string
		reachableSizes []string
		expected       string
	}{
		{
			"NotAtSmallestCanReduce",
			"modalSize",
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			"modalSize !== 'md'",
		},
		{"SingleSizeCannotReduce", "modalSize", []string{ModalSizeXl}, ""},
		{"EmptyCannotReduce", "modalSize", nil, ""},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalCanReduceExpressionBuilder(testCase.sizePath, testCase.reachableSizes)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalEnlargeExpressionBuilder(t *testing.T) {
	tests := []struct {
		name           string
		sizePath       string
		reachableSizes []string
		expected       string
	}{
		{
			"StepsUpToLargest",
			"modalSize",
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			"modalSize = modalSize === 'md' ? 'lg' : modalSize === 'lg' ? 'xl' : 'xl'",
		},
		{
			"SingleSizeStaysPut",
			"modalSize",
			[]string{ModalSizeXl},
			"modalSize = 'xl'",
		},
		{
			"ThreeSizesClampAtLargest",
			"modalSize",
			[]string{ModalSizeXl, ModalSizeXxl, ModalSizeFull},
			"modalSize = modalSize === 'xl' ? 'xxl' : modalSize === 'xxl' ? 'full' : 'full'",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalEnlargeExpressionBuilder(testCase.sizePath, testCase.reachableSizes)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

func TestModalReduceExpressionBuilder(t *testing.T) {
	tests := []struct {
		name           string
		sizePath       string
		reachableSizes []string
		expected       string
	}{
		{
			"StepsDownToSmallest",
			"modalSize",
			[]string{ModalSizeMd, ModalSizeLg, ModalSizeXl},
			"modalSize = modalSize === 'xl' ? 'lg' : modalSize === 'lg' ? 'md' : 'md'",
		},
		{
			"SingleSizeStaysPut",
			"modalSize",
			[]string{ModalSizeXl},
			"modalSize = 'xl'",
		},
		{
			"ThreeSizesClampAtSmallest",
			"modalSize",
			[]string{ModalSizeXl, ModalSizeXxl, ModalSizeFull},
			"modalSize = modalSize === 'full' ? 'xxl' : modalSize === 'xxl' ? 'xl' : 'xl'",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			resolved := modalReduceExpressionBuilder(testCase.sizePath, testCase.reachableSizes)

			if resolved != testCase.expected {
				t.Errorf("resolved = %q, want %q", resolved, testCase.expected)
			}
		})
	}
}

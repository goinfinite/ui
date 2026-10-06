package uiToolset

import "testing"

func TestBackgroundColorClassResolver(t *testing.T) {
	testCases := []struct {
		name            string
		backgroundColor string
		fallbackClass   string
		want            string
	}{
		{name: "empty takes the fallback", backgroundColor: "", fallbackClass: "bg-neutral-950/20", want: "bg-neutral-950/20"},
		{name: "token is prefixed", backgroundColor: "neutral-800/50", fallbackClass: "bg-neutral-950/20", want: "bg-neutral-800/50"},
		{name: "empty without a fallback stays unset", backgroundColor: "", fallbackClass: "", want: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := BackgroundColorClassResolver(testCase.backgroundColor, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("BackgroundColorClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestTextColorClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		textColor     string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", textColor: "", fallbackClass: "text-neutral-50", want: "text-neutral-50"},
		{name: "token is prefixed", textColor: "secondary-500", fallbackClass: "text-neutral-50", want: "text-secondary-500"},
		{name: "empty without a fallback stays unset", textColor: "", fallbackClass: "", want: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := TextColorClassResolver(testCase.textColor, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("TextColorClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestBorderColorClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		borderColor   string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", borderColor: "", fallbackClass: "border-neutral-50/5", want: "border-neutral-50/5"},
		{name: "token is prefixed", borderColor: "amber-500", fallbackClass: "border-neutral-50/5", want: "border-amber-500"},
		{name: "empty without a fallback stays unset", borderColor: "", fallbackClass: "", want: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := BorderColorClassResolver(testCase.borderColor, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("BorderColorClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestBorderRadiusClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		borderRadius  string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", borderRadius: "", fallbackClass: "rounded-lg", want: "rounded-lg"},
		{name: "none", borderRadius: "none", fallbackClass: "rounded-lg", want: "rounded-none"},
		{name: "xs", borderRadius: "xs", fallbackClass: "rounded-lg", want: "rounded-xs"},
		{name: "sm", borderRadius: "sm", fallbackClass: "rounded-lg", want: "rounded-sm"},
		{name: "md", borderRadius: "md", fallbackClass: "rounded-lg", want: "rounded-md"},
		{name: "lg", borderRadius: "lg", fallbackClass: "rounded-lg", want: "rounded-lg"},
		{name: "xl", borderRadius: "xl", fallbackClass: "rounded-lg", want: "rounded-xl"},
		{name: "full", borderRadius: "full", fallbackClass: "rounded-lg", want: "rounded-full"},
		{name: "unknown takes the fallback", borderRadius: "huge", fallbackClass: "rounded-lg", want: "rounded-lg"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := BorderRadiusClassResolver(testCase.borderRadius, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("BorderRadiusClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestShapeClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		shape         string
		fallbackClass string
		want          string
	}{
		{name: "circular", shape: "circular", fallbackClass: "rounded", want: "rounded-full"},
		{name: "rounded", shape: "rounded", fallbackClass: "rounded-full", want: "rounded"},
		{name: "square", shape: "square", fallbackClass: "rounded", want: "rounded-none"},
		{name: "empty takes the fallback", shape: "", fallbackClass: "rounded-full", want: "rounded-full"},
		{name: "unknown takes the fallback", shape: "hexagon", fallbackClass: "rounded-full", want: "rounded-full"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ShapeClassResolver(testCase.shape, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("ShapeClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestShadowClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		shadowSize    string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", shadowSize: "", fallbackClass: "", want: ""},
		{name: "none", shadowSize: "none", fallbackClass: "", want: "shadow-none"},
		{name: "xs", shadowSize: "xs", fallbackClass: "", want: "shadow-xs"},
		{name: "sm", shadowSize: "sm", fallbackClass: "", want: "shadow-sm"},
		{name: "md", shadowSize: "md", fallbackClass: "", want: "shadow-md"},
		{name: "lg", shadowSize: "lg", fallbackClass: "", want: "shadow-lg"},
		{name: "xl", shadowSize: "xl", fallbackClass: "", want: "shadow-xl"},
		{name: "unknown takes the fallback", shadowSize: "huge", fallbackClass: "", want: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ShadowClassResolver(testCase.shadowSize, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("ShadowClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestRingThicknessClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		ringThickness string
		fallbackClass string
		want          string
	}{
		{name: "xs", ringThickness: "xs", fallbackClass: "ring-0", want: "ring-1"},
		{name: "sm", ringThickness: "sm", fallbackClass: "ring-0", want: "ring-1.5"},
		{name: "md", ringThickness: "md", fallbackClass: "ring-0", want: "ring-2"},
		{name: "lg", ringThickness: "lg", fallbackClass: "ring-0", want: "ring-2.5"},
		{name: "xl", ringThickness: "xl", fallbackClass: "ring-0", want: "ring-3"},
		{name: "empty takes the fallback", ringThickness: "", fallbackClass: "ring-0", want: "ring-0"},
		{name: "unknown takes the fallback", ringThickness: "huge", fallbackClass: "ring-0", want: "ring-0"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := RingThicknessClassResolver(testCase.ringThickness, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("RingThicknessClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestRingClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		ringColor     string
		ringThickness string
		want          string
	}{
		{name: "empty color and thickness render no ring", ringColor: "", ringThickness: "", want: ""},
		{name: "color without thickness renders no ring", ringColor: "red-500/50", ringThickness: "", want: ""},
		{name: "thickness without color renders no ring", ringColor: "", ringThickness: "md", want: ""},
		{name: "xs thickness", ringColor: "red-500/50", ringThickness: "xs", want: "ring-1 ring-red-500/50"},
		{name: "sm thickness", ringColor: "red-500/50", ringThickness: "sm", want: "ring-1.5 ring-red-500/50"},
		{name: "md thickness", ringColor: "red-500/50", ringThickness: "md", want: "ring-2 ring-red-500/50"},
		{name: "lg thickness", ringColor: "red-500/50", ringThickness: "lg", want: "ring-2.5 ring-red-500/50"},
		{name: "xl thickness", ringColor: "red-500/50", ringThickness: "xl", want: "ring-3 ring-red-500/50"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := RingClassResolver(testCase.ringColor, testCase.ringThickness)
			if got != testCase.want {
				t.Errorf("RingClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestCompactRingClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		ringColor     string
		ringThickness string
		want          string
	}{
		{name: "empty color renders no ring", ringColor: "", ringThickness: "md", want: ""},
		{name: "empty thickness defaults to the thin ring", ringColor: "red-500/50", ringThickness: "", want: "ring-1 ring-red-500/50"},
		{name: "xs thickness", ringColor: "red-500/50", ringThickness: "xs", want: "ring-0.5 ring-red-500/50"},
		{name: "sm thickness", ringColor: "red-500/50", ringThickness: "sm", want: "ring-1 ring-red-500/50"},
		{name: "md thickness", ringColor: "red-500/50", ringThickness: "md", want: "ring-1.5 ring-red-500/50"},
		{name: "lg thickness", ringColor: "red-500/50", ringThickness: "lg", want: "ring-2 ring-red-500/50"},
		{name: "xl thickness", ringColor: "red-500/50", ringThickness: "xl", want: "ring-2.5 ring-red-500/50"},
		{name: "unknown thickness defaults to the thin ring", ringColor: "red-500/50", ringThickness: "huge", want: "ring-1 ring-red-500/50"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := CompactRingClassResolver(testCase.ringColor, testCase.ringThickness)
			if got != testCase.want {
				t.Errorf("CompactRingClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestPaddingClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		paddingSize   string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", paddingSize: "", fallbackClass: "", want: ""},
		{name: "none", paddingSize: "none", fallbackClass: "", want: "p-0"},
		{name: "xs", paddingSize: "xs", fallbackClass: "", want: "p-3"},
		{name: "sm", paddingSize: "sm", fallbackClass: "", want: "p-4"},
		{name: "md", paddingSize: "md", fallbackClass: "", want: "p-5"},
		{name: "lg", paddingSize: "lg", fallbackClass: "", want: "p-6"},
		{name: "xl", paddingSize: "xl", fallbackClass: "", want: "p-8"},
		{name: "unknown without a fallback stays unset", paddingSize: "huge", fallbackClass: "", want: ""},
		{name: "unknown takes the fallback", paddingSize: "huge", fallbackClass: "p-5", want: "p-5"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := PaddingClassResolver(testCase.paddingSize, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("PaddingClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestCompactPaddingClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		paddingSize   string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", paddingSize: "", fallbackClass: "p-3", want: "p-3"},
		{name: "none", paddingSize: "none", fallbackClass: "p-3", want: "p-0"},
		{name: "xs", paddingSize: "xs", fallbackClass: "p-3", want: "p-1"},
		{name: "sm", paddingSize: "sm", fallbackClass: "p-3", want: "p-1.5"},
		{name: "md", paddingSize: "md", fallbackClass: "p-3", want: "p-2"},
		{name: "lg", paddingSize: "lg", fallbackClass: "p-3", want: "p-2.5"},
		{name: "xl", paddingSize: "xl", fallbackClass: "p-3", want: "p-3"},
		{name: "unknown takes the fallback", paddingSize: "huge", fallbackClass: "p-2", want: "p-2"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := CompactPaddingClassResolver(testCase.paddingSize, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("CompactPaddingClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestGapClassResolver(t *testing.T) {
	testCases := []struct {
		name          string
		gapSize       string
		fallbackClass string
		want          string
	}{
		{name: "empty takes the fallback", gapSize: "", fallbackClass: "gap-1", want: "gap-1"},
		{name: "none", gapSize: "none", fallbackClass: "gap-1", want: "gap-0"},
		{name: "xs", gapSize: "xs", fallbackClass: "gap-1", want: "gap-1"},
		{name: "sm", gapSize: "sm", fallbackClass: "gap-1", want: "gap-2"},
		{name: "md", gapSize: "md", fallbackClass: "gap-1", want: "gap-3"},
		{name: "lg", gapSize: "lg", fallbackClass: "gap-1", want: "gap-6"},
		{name: "xl", gapSize: "xl", fallbackClass: "gap-1", want: "gap-8"},
		{name: "unknown takes the fallback", gapSize: "huge", fallbackClass: "gap-1", want: "gap-1"},
		{name: "unknown without a fallback stays unset", gapSize: "huge", fallbackClass: "", want: ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := GapClassResolver(testCase.gapSize, testCase.fallbackClass)
			if got != testCase.want {
				t.Errorf("GapClassResolver() = %q, want %q", got, testCase.want)
			}
		})
	}
}

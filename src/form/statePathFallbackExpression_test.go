package uiForm

import "testing"

func TestStatePathFallbackExpressionBuilder(testingHandle *testing.T) {
	testCases := []struct {
		name          string
		statePath     string
		fallback      string
		expectedValue string
	}{
		{
			name:          "empty path yields no expression",
			statePath:     "",
			fallback:      "Label",
			expectedValue: "",
		},
		{
			name:          "path falls back to a quoted label",
			statePath:     "fieldLabel",
			fallback:      "Role",
			expectedValue: `(fieldLabel) || "Role"`,
		},
		{
			name:          "fallback quotes are escaped",
			statePath:     "fieldLabel",
			fallback:      `Say "hi"`,
			expectedValue: `(fieldLabel) || "Say \"hi\""`,
		},
	}
	for _, testCase := range testCases {
		testingHandle.Run(testCase.name, func(subtestHandle *testing.T) {
			actualValue := statePathFallbackExpressionBuilder(
				testCase.statePath, testCase.fallback,
			)
			if actualValue != testCase.expectedValue {
				subtestHandle.Errorf(
					"FallbackExpressionMismatch: got %q, want %q",
					actualValue, testCase.expectedValue,
				)
			}
		})
	}
}

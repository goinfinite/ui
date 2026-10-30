package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestToggleSwitchRendersLabelStatePath(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := ToggleSwitch(ToggleSwitchSettings{
		Label:                "Notifications",
		LabelOneWayStatePath: "toggleLabel",
		TwoWayStatePath:      "notificationsEnabled",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("ToggleSwitchRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="toggleLabel"`,
		"Notifications",
		`x-model="notificationsEnabled"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestFileUploadInputRendersLabelStatePaths(testingHandle *testing.T) {
	var buffer bytes.Buffer
	renderErr := FileUploadInput(FileUploadInputSettings{
		InputName:                  "attachments",
		Label:                      "Attachments",
		LabelOneWayStatePath:       "attachmentsLabel",
		Description:                "Upload up to three files.",
		DescriptionOneWayStatePath: "attachmentsDescription",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		testingHandle.Fatalf("FileUploadInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedFragment := range []string{
		`x-text="attachmentsLabel"`,
		`x-text="attachmentsDescription"`,
		"Attachments",
		"Upload up to three files.",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

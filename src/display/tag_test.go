package uiDisplay

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestTagRendersTinySizeClasses(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Tag(TagSettings{
		Size:              TagSizeTiny,
		OuterLeftLabel:    "cpu",
		InnerLabel:        "cores",
		OnRemoveFunc:      "removeCpu()",
		RemoveButtonLabel: "Remove cpu",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("TagRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"text-[0.625rem] p-0.25 gap-0.25",
		"text-[0.625rem] p-0.5 gap-0.5",
		"-m-0.5 p-0.5",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestTagRendersBoundValueWithoutLowercase(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := Tag(TagSettings{
		OuterLeftLabel:            "CPU cores",
		InnerValueOneWayStatePath: "cpuCores",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("TagRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if !strings.Contains(renderedHtml, `x-text="cpuCores"`) {
		t.Error("RenderedValueMissingBoundAttribute")
	}
	if strings.Contains(renderedHtml, `lowercase" x-text="cpuCores"`) {
		t.Error("RenderedValueMustNotForceLowercase")
	}
}

package uiToolset

import (
	"strings"
	"testing"
)

func TestComponentIdPrefixGenerator(t *testing.T) {
	generator := NewComponentIdPrefixGenerator("record")

	firstId := generator.GenerateNext()
	secondId := generator.GenerateNext()

	if !strings.HasPrefix(firstId, "record-") {
		t.Errorf("IdPrefixMissingComponentName: got %q, want prefix %q", firstId, "record-")
	}
	if firstId == secondId {
		t.Errorf("IdPrefixNotUnique: %q == %q", firstId, secondId)
	}

	otherGenerator := NewComponentIdPrefixGenerator("other")
	if otherGenerator.GenerateNext() == firstId {
		t.Errorf("IdPrefixSharedAcrossGenerators: got %q", firstId)
	}
}

package uiForm

import (
	"testing"

	uiControl "github.com/goinfinite/ui/src/control"
)

func TestFloatingButtonSizeResolver(t *testing.T) {
	if actual := floatingButtonSizeResolver(FloatingButtonStyle{}); actual != uiControl.ButtonSizeXs {
		t.Errorf("DefaultSizeMismatch: got %q", actual)
	}

	providedStyle := FloatingButtonStyle{Size: uiControl.ButtonSizeLg}
	if actual := floatingButtonSizeResolver(providedStyle); actual != uiControl.ButtonSizeLg {
		t.Errorf("ProvidedSizeMismatch: got %q", actual)
	}
}
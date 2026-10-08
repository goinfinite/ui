package uiStructural

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestCarouselItemTooltip(t *testing.T) {
	renderTooltip := func(
		t *testing.T, settings CarouselItemTooltipSettings,
	) string {
		t.Helper()
		var buffer bytes.Buffer
		itemContent := templ.Raw("<span>item</span>")
		renderContext := templ.WithChildren(context.Background(), itemContent)
		renderErr := CarouselItemTooltip(settings).Render(renderContext, &buffer)
		if renderErr != nil {
			t.Fatalf("CarouselItemTooltipRenderFailed: %v", renderErr)
		}
		return buffer.String()
	}

	t.Run("content teleports a tooltip to the body", func(t *testing.T) {
		renderedHtml := renderTooltip(t, CarouselItemTooltipSettings{
			Content: "Full description",
		})
		for _, expectedFragment := range []string{
			`class="h-full"`,
			`x-data="tooltip(`,
			`x-ref="trigger"`,
			`x-teleport="body"`,
			`role="tooltip"`,
			`<span>item</span>`,
			"Full description",
		} {
			if !strings.Contains(renderedHtml, expectedFragment) {
				t.Errorf("RenderedHtmlMissing: %s", expectedFragment)
			}
		}
	})

	t.Run("the trigger describes the tooltip", func(t *testing.T) {
		renderedHtml := renderTooltip(t, CarouselItemTooltipSettings{
			Content: "Full description",
		})
		describedById := regexp.MustCompile(`aria-describedby="([^"]{1,128})"`).
			FindStringSubmatch(renderedHtml)
		tooltipId := regexp.MustCompile(`id="(carousel-item-tooltip-\d{1,18})"`).
			FindStringSubmatch(renderedHtml)
		if describedById == nil || tooltipId == nil {
			t.Fatalf("RenderedHtmlMissingIdWiring: %q", renderedHtml)
		}
		if describedById[1] != tooltipId[1] {
			t.Errorf(
				"TooltipIdMismatch: aria-describedby=%q id=%q",
				describedById[1], tooltipId[1],
			)
		}
	})

	t.Run("no content renders a plain wrapper", func(t *testing.T) {
		renderedHtml := renderTooltip(t, CarouselItemTooltipSettings{})
		if !strings.Contains(renderedHtml, `class="h-full"`) {
			t.Errorf("RenderedHtmlMissingWrapper: %q", renderedHtml)
		}
		for _, unexpectedFragment := range []string{
			"x-teleport", "aria-describedby", "x-data",
		} {
			if strings.Contains(renderedHtml, unexpectedFragment) {
				t.Errorf("RenderedHtmlHasTooltip: %q", unexpectedFragment)
			}
		}
	})

	t.Run("customization tokens reach the surface", func(t *testing.T) {
		renderedHtml := renderTooltip(t, CarouselItemTooltipSettings{
			Content:         "Full description",
			BackgroundColor: "primary-900/95",
			TextColor:       "amber-50",
			RingColor:       "primary-500/40",
			RingThickness:   uiToolset.RingThicknessXs,
		})
		for _, expectedFragment := range []string{
			"bg-primary-900/95", "text-amber-50", "ring-0.5 ring-primary-500/40",
		} {
			if !strings.Contains(renderedHtml, expectedFragment) {
				t.Errorf("RenderedHtmlMissingToken: %s", expectedFragment)
			}
		}
	})
}

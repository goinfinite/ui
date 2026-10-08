package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRemoteSearchableSelectInputRemoteConfigResolver(t *testing.T) {
	testCases := []struct {
		name           string
		optionsUrl     string
		queryParam     string
		minQueryLength uint
		debounceMs     uint
		expected       searchableSelectInputRemoteConfig
	}{
		{
			name:       "defaults",
			optionsUrl: "/api/countries",
			expected: searchableSelectInputRemoteConfig{
				Url:            "/api/countries",
				QueryParam:     "q",
				MinQueryLength: 3,
				DebounceMs:     300,
			},
		},
		{
			name:           "explicit settings win over the defaults",
			optionsUrl:     "/api/countries",
			queryParam:     "countryName",
			minQueryLength: 1,
			debounceMs:     800,
			expected: searchableSelectInputRemoteConfig{
				Url:            "/api/countries",
				QueryParam:     "countryName",
				MinQueryLength: 1,
				DebounceMs:     800,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualConfig := remoteSearchableSelectInputRemoteConfigResolver(
				testCase.optionsUrl, testCase.queryParam,
				testCase.minQueryLength, testCase.debounceMs,
			)
			if *actualConfig != testCase.expected {
				t.Errorf(
					"RemoteConfigMismatch: got %+v, want %+v",
					*actualConfig, testCase.expected,
				)
			}
		})
	}
}

func TestRemoteSearchableSelectInputRendersOpaqueDropdownAboveStickyHeaders(t *testing.T) {
	settings := RemoteSearchableSelectInputSettings{
		InputName:       "country",
		Label:           "Country",
		OptionsUrl:      "/api/countries",
		TwoWayStatePath: "country",
	}
	var buffer bytes.Buffer
	renderErr := RemoteSearchableSelectInput(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("RemoteSearchableSelectInputRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	for _, expectedText := range []string{"z-20", "bg-neutral-800/95"} {
		if !strings.Contains(renderedHtml, expectedText) {
			t.Errorf("RenderedHtmlMissing(%q)", expectedText)
		}
	}

	settings.DropdownBackgroundColor = "emerald-900"
	coloredBuffer := bytes.Buffer{}
	renderErr = RemoteSearchableSelectInput(settings).Render(context.Background(), &coloredBuffer)
	if renderErr != nil {
		t.Fatalf("RemoteSearchableSelectInputRenderFailed: %v", renderErr)
	}
	coloredHtml := coloredBuffer.String()
	if !strings.Contains(coloredHtml, "bg-emerald-900") {
		t.Errorf("RenderedHtmlMissingCustomDropdownBackground")
	}
	if strings.Contains(coloredHtml, "bg-neutral-800/95") {
		t.Errorf("RenderedHtmlKeepsDefaultDropdownBackground")
	}
}

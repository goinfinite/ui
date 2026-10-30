package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiControl "github.com/goinfinite/ui/src/control"
	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestDateTimePickerModeResolver(t *testing.T) {
	tests := []struct {
		mode     string
		expected string
	}{
		{"", DateTimePickerModeDate},
		{DateTimePickerModeDate, DateTimePickerModeDate},
		{DateTimePickerModeDateTime, DateTimePickerModeDateTime},
		{DateTimePickerModeTime, DateTimePickerModeTime},
		{"invalid", DateTimePickerModeDate},
	}
	for _, test := range tests {
		settings := DateTimePickerSettings{Mode: test.mode}
		if resolved := dateTimePickerModeResolver(settings); resolved != test.expected {
			t.Errorf("modeResolver(%s) = %s, expected %s", test.mode, resolved, test.expected)
		}
	}
}

func TestDateTimePickerLocaleResolver(t *testing.T) {
	settings := DateTimePickerSettings{}
	if resolved := dateTimePickerLocaleResolver(settings); resolved != "en-US" {
		t.Errorf("localeResolver default = %s, expected en-US", resolved)
	}
	settings.Locale = "pt-BR"
	if resolved := dateTimePickerLocaleResolver(settings); resolved != "pt-BR" {
		t.Errorf("localeResolver custom = %s, expected pt-BR", resolved)
	}
}

func TestDateTimePickerRangeEndInputNameResolver(t *testing.T) {
	settings := DateTimePickerSettings{InputName: "booking"}
	if resolved := dateTimePickerRangeEndInputNameResolver(settings); resolved != "booking-end" {
		t.Errorf("rangeEndInputName default = %s, expected booking-end", resolved)
	}
	settings.RangeEndInputName = "bookingEnd"
	if resolved := dateTimePickerRangeEndInputNameResolver(settings); resolved != "bookingEnd" {
		t.Errorf("rangeEndInputName custom = %s, expected bookingEnd", resolved)
	}
}

func TestDateTimePickerInputPlaceholderResolver(t *testing.T) {
	tests := []struct {
		mode     string
		expected string
	}{
		{DateTimePickerModeDate, "yyyy-mm-dd"},
		{DateTimePickerModeTime, "HH:MM"},
		{DateTimePickerModeDateTime, "yyyy-mm-ddTHH:MM"},
	}
	for _, test := range tests {
		resolved := dateTimePickerInputPlaceholderResolver(test.mode)
		if resolved != test.expected {
			t.Errorf("placeholderResolver(%s) = %s, expected %s", test.mode, resolved, test.expected)
		}
	}
}

func TestDateTimePickerStyleConfigResolver(t *testing.T) {
	defaultStyle := dateTimePickerStyleConfigResolver(DateTimePickerStyle{})
	expectedDefaults := map[string]string{
		"daySelectedBackgroundColor": defaultStyle.DaySelectedBackgroundColor,
		"dayTextColor":               defaultStyle.DayTextColor,
		"dayTodayRingColor":          defaultStyle.DayTodayRingColor,
	}
	expectedDefaultValues := map[string]string{
		"daySelectedBackgroundColor": "secondary-500",
		"dayTextColor":               "neutral-100",
		"dayTodayRingColor":          "neutral-50/30",
	}
	for key, value := range expectedDefaults {
		if value != expectedDefaultValues[key] {
			t.Errorf("%s default = %s, expected %s", key, value, expectedDefaultValues[key])
		}
	}
	customStyle := dateTimePickerStyleConfigResolver(DateTimePickerStyle{
		DaySelectedBackgroundColor: "emerald-500",
		DayTextColor:               "neutral-50",
	})
	if customStyle.DaySelectedBackgroundColor != "emerald-500" {
		t.Errorf("daySelectedBackgroundColor custom = %s", customStyle.DaySelectedBackgroundColor)
	}
	if customStyle.DayTextColor != "neutral-50" {
		t.Errorf("dayTextColor custom = %s", customStyle.DayTextColor)
	}
}

func TestDateTimePickerDoneButtonSizeResolver(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"", "xs"},
		{"xs", "xs"},
		{"sm", "sm"},
		{"md", "md"},
		{"lg", "lg"},
		{"xl", "xl"},
		{"invalid", "xs"},
	}
	for _, test := range tests {
		resolved := dateTimePickerDoneButtonSizeResolver(test.size)
		if resolved != test.expected {
			t.Errorf("doneButtonSize(%s) = %s, expected %s", test.size, resolved, test.expected)
		}
	}
}

func TestDateTimePickerTimeSelectorSizeResolver(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"", SelectInputSizeMd},
		{SelectInputSizeXs, SelectInputSizeXs},
		{SelectInputSizeSm, SelectInputSizeSm},
		{SelectInputSizeLg, SelectInputSizeLg},
		{SelectInputSizeXl, SelectInputSizeXl},
		{"invalid", SelectInputSizeMd},
	}
	for _, test := range tests {
		resolved := dateTimePickerTimeSelectorSizeResolver(test.size)
		if resolved != test.expected {
			t.Errorf("timeSelectorSize(%s) = %s, expected %s", test.size, resolved, test.expected)
		}
	}
}

func TestDateTimePickerDoneButtonTextCaseResolver(t *testing.T) {
	if resolved := dateTimePickerDoneButtonTextCaseResolver(""); resolved != uiToolset.TextCaseLower {
		t.Errorf("doneButtonTextCase default = %s, expected lower", resolved)
	}
	if resolved := dateTimePickerDoneButtonTextCaseResolver(uiToolset.TextCaseUpper); resolved != uiToolset.TextCaseUpper {
		t.Errorf("doneButtonTextCase custom = %s, expected upper", resolved)
	}
}

func TestDateTimePickerDoneButtonIconSizeClassResolver(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"", "text-xs"},
		{uiControl.ButtonSizeXs, "text-xs"},
		{uiControl.ButtonSizeSm, "text-sm"},
		{uiControl.ButtonSizeMd, "text-base"},
		{uiControl.ButtonSizeLg, "text-lg"},
		{uiControl.ButtonSizeXl, "text-xl"},
	}
	for _, test := range tests {
		resolved := dateTimePickerDoneButtonIconSizeClassResolver(test.size)
		if resolved != test.expected {
			t.Errorf("doneButtonIconSizeClass(%s) = %s, expected %s", test.size, resolved, test.expected)
		}
	}
}

func TestDateTimePickerTimezoneOptionsResolver(t *testing.T) {
	defaultOptions := dateTimePickerTimezoneOptionsResolver("", nil)
	if len(defaultOptions) == 0 {
		t.Fatal("timezoneOptions default is empty")
	}
	if defaultOptions[0].Value != "UTC" {
		t.Errorf("timezoneOptions default first = %s, expected UTC", defaultOptions[0].Value)
	}
	selectedOptions := dateTimePickerTimezoneOptionsResolver("America/Manaus", nil)
	if selectedOptions[0].Value != "America/Manaus" {
		t.Errorf("timezoneOptions selected first = %s, expected America/Manaus", selectedOptions[0].Value)
	}
	if selectedOptions[1].Value != "UTC" {
		t.Errorf("timezoneOptions selected second = %s, expected UTC", selectedOptions[1].Value)
	}
	customOptions := dateTimePickerTimezoneOptionsResolver("", []string{"America/Sao_Paulo"})
	if len(customOptions) != 1 || customOptions[0].Value != "America/Sao_Paulo" {
		t.Errorf("timezoneOptions custom = %+v", customOptions)
	}
	if customOptions[0].Label != "America/Sao Paulo" {
		t.Errorf("timezoneOptions label = %s, expected America/Sao Paulo", customOptions[0].Label)
	}
}

func TestDateTimePickerConfigJsonBuilder(t *testing.T) {
	settings := DateTimePickerSettings{
		Label:           "Booking date",
		FirstDayOfWeek:  "sunday",
		IsRange:         true,
		Locale:          "pt-BR",
		MaxValue:        "2026-12-31",
		MinValue:        "2026-01-01",
		Mode:            DateTimePickerModeDateTime,
		RangeEndValue:   "2026-10-20T18:00",
		Timezone:        "America/Sao_Paulo",
		TwoWayStatePath: "bookingDate",
		Value:           "2026-10-16T09:00",
	}
	configJson := dateTimePickerConfigJsonBuilder(settings)
	expectedFragments := []string{
		`"mode":"datetime"`,
		`"locale":"pt-BR"`,
		`"minValue":"2026-01-01"`,
		`"maxValue":"2026-12-31"`,
		`"firstDayOfWeek":"sunday"`,
		`"isRange":true`,
		`"timezone":"America/Sao_Paulo"`,
		`"daySelectedBackgroundColor":"secondary-500"`,
		`"dayTextColor":"neutral-100"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(configJson, expectedFragment) {
			t.Errorf("configJson missing fragment %s: %s", expectedFragment, configJson)
		}
	}
}

func TestDateTimePickerDataExpressionBuilder(t *testing.T) {
	staticSettings := DateTimePickerSettings{Value: "2026-10-16"}
	if resolved := dateTimePickerDataExpressionBuilder(staticSettings, "{}"); resolved != `dateTimePicker({}, "2026-10-16", "", "")` {
		t.Errorf("dataExpression static = %s", resolved)
	}
	pathSettings := DateTimePickerSettings{
		TwoWayStatePath:         "bookingStart",
		RangeEndTwoWayStatePath: "bookingEnd",
	}
	if resolved := dateTimePickerDataExpressionBuilder(pathSettings, "{}"); resolved != `dateTimePicker({}, bookingStart, bookingEnd, "")` {
		t.Errorf("dataExpression paths = %s", resolved)
	}
	timezonePathSettings := DateTimePickerSettings{TimezoneTwoWayStatePath: "bookingTimezone"}
	if resolved := dateTimePickerDataExpressionBuilder(timezonePathSettings, "{}"); resolved != `dateTimePicker({}, "", "", bookingTimezone)` {
		t.Errorf("dataExpression timezone path = %s", resolved)
	}
}

func TestDateTimePickerWritebackExpressionBuilder(t *testing.T) {
	settings := DateTimePickerSettings{}
	if resolved := dateTimePickerWritebackExpressionBuilder(settings); resolved != "" {
		t.Errorf("writebackExpression empty = %s", resolved)
	}
	settings.TwoWayStatePath = "bookingDate"
	if resolved := dateTimePickerWritebackExpressionBuilder(settings); resolved != "initialized && (bookingDate = startDraft)" {
		t.Errorf("writebackExpression single = %s", resolved)
	}
	settings.RangeEndTwoWayStatePath = "bookingEnd"
	if resolved := dateTimePickerWritebackExpressionBuilder(settings); resolved != "initialized && (bookingDate = startDraft, bookingEnd = endDraft)" {
		t.Errorf("writebackExpression range = %s", resolved)
	}
	settings.TimezoneTwoWayStatePath = "bookingTimezone"
	if resolved := dateTimePickerWritebackExpressionBuilder(settings); resolved != "initialized && (bookingDate = startDraft, bookingEnd = endDraft, bookingTimezone = timezone)" {
		t.Errorf("writebackExpression timezone = %s", resolved)
	}
}

func renderDateTimePicker(t *testing.T, settings DateTimePickerSettings) string {
	t.Helper()
	var buffer bytes.Buffer
	renderErr := DateTimePicker(settings).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("Render failed: %v", renderErr)
	}
	return buffer.String()
}

func TestDateTimePickerRenderDateMode(t *testing.T) {
	renderedHtml := renderDateTimePicker(t, DateTimePickerSettings{
		Label:           "Booking date",
		InputName:       "booking",
		TwoWayStatePath: "bookingDate",
		Value:           "2026-10-16",
	})
	expectedFragments := []string{
		`x-effect="initialized &amp;&amp; (bookingDate = startDraft)"`,
		`, bookingDate, &#34;&#34;, &#34;&#34;)`,
		`x-data="dateTimePicker({&#34;mode&#34;:&#34;date&#34;`,
		`placeholder="yyyy-mm-dd"`,
		`aria-label="Toggle date picker"`,
		`:aria-expanded="isOpen"`,
		`ph-check text-xs`,
		`lowercase`,
		`bg-neutral-50/5 hover:bg-neutral-50/20`,
		`ring-neutral-50/5 hover:ring-neutral-50/30`,
		`x-model="startTypedValue"`,
		`commitTypedValueForSide('start')`,
		`x-show="isOpen"`,
		`overflow-visible`,
		`x-text="monthYearLabel"`,
		`aria-label="Previous year"`,
		`aria-label="Next month"`,
		`x-for="weekdayLabel in weekdayLabels"`,
		`:data-day-index="dayIndex"`,
		`:class="dayClassesBuilder(day)"`,
		`:class="dayNumberClassesBuilder(day)"`,
		`:aria-current="day.isSelected ? 'date' : false"`,
		`name="booking"`,
		`:value="startDraft"`,
		`@click.outside="closePopover()"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	unexpectedFragments := []string{
		`>Hour</legend>`,
		`>Minute</legend>`,
		`aria-label="Booking window end"`,
		`name="booking-end"`,
	}
	for _, unexpectedFragment := range unexpectedFragments {
		if strings.Contains(renderedHtml, unexpectedFragment) {
			t.Errorf("rendered html unexpectedly contains fragment %s", unexpectedFragment)
		}
	}
}

func TestDateTimePickerRenderRangeDateMode(t *testing.T) {
	renderedHtml := renderDateTimePicker(t, DateTimePickerSettings{
		Label:                   "Booking window",
		InputName:               "booking",
		IsRange:                 true,
		RangeEndTwoWayStatePath: "bookingEnd",
		TwoWayStatePath:         "bookingStart",
	})
	expectedFragments := []string{
		`x-effect="initialized &amp;&amp; (bookingStart = startDraft, bookingEnd = endDraft)"`,
		`, bookingStart, bookingEnd, &#34;&#34;)`,
		`placeholder="yyyy-mm-dd"`,
		`x-model="startTypedValue"`,
		`x-model="endTypedValue"`,
		`activateEnd(&#39;start&#39;)`,
		`activateEnd(&#39;end&#39;)`,
		`commitTypedValueForSide(&#39;end&#39;)`,
		`aria-label="Booking window start"`,
		`aria-label="Booking window end"`,
		`>Start</span>`,
		`>End</span>`,
		`my-1.5 w-px shrink-0 bg-neutral-50/10`,
		`name="booking"`,
		`name="booking-end"`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, "→") {
		t.Error("rendered html still contains the arrow range divider")
	}
}

func TestDateTimePickerRenderTimeMode(t *testing.T) {
	renderedHtml := renderDateTimePicker(t, DateTimePickerSettings{
		Label:    "Start time",
		Mode:     DateTimePickerModeTime,
		Timezone: "America/Sao_Paulo",
	})
	expectedFragments := []string{
		`, &#34;&#34;, &#34;&#34;, &#34;America/Sao_Paulo&#34;)`,
		`placeholder="HH:MM"`,
		`>Hour</legend>`,
		`>Minute</legend>`,
		`x-ref="selectTrigger"`,
		`String(timeHour) === $el.value`,
		`String(timeMinute) === $el.value`,
		`commitTime()`,
		`max-h-40`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
	if strings.Contains(renderedHtml, `x-text="monthYearLabel"`) {
		t.Error("rendered html unexpectedly contains the calendar")
	}
	if strings.Contains(renderedHtml, "timePartIsDisabled") {
		t.Error("rendered html still contains the native select disabled expression")
	}
	if strings.Contains(renderedHtml, ">Timezone</legend>") {
		t.Error("rendered html unexpectedly contains the timezone select in time mode")
	}
}

func TestDateTimePickerRenderDateTimeMode(t *testing.T) {
	renderedHtml := renderDateTimePicker(t, DateTimePickerSettings{
		Label: "Publish at",
		Mode:  DateTimePickerModeDateTime,
	})
	if !strings.Contains(renderedHtml, `x-text="monthYearLabel"`) {
		t.Error("rendered html missing the calendar")
	}
	if !strings.Contains(renderedHtml, `>Hour</legend>`) {
		t.Error("rendered html missing the time selects")
	}
	if !strings.Contains(renderedHtml, `placeholder="yyyy-mm-ddTHH:MM"`) {
		t.Error("rendered html missing the datetime placeholder")
	}
}

func TestDateTimePickerRenderCustomizations(t *testing.T) {
	renderedHtml := renderDateTimePicker(t, DateTimePickerSettings{
		Label:                   "Visit",
		Mode:                    DateTimePickerModeDateTime,
		TextCase:                uiToolset.TextCaseUpper,
		DoneButtonIconLeft:      "ph-calendar-check",
		DoneButtonIconRight:     "ph-seal-check",
		DoneButtonLabel:         "Apply",
		DoneButtonSize:          uiControl.ButtonSizeSm,
		DoneButtonStyle:         DateTimePickerDoneButtonStyle{BackgroundColor: "secondary-500", IconSizeClass: "text-sm", TextColor: "secondary-950", TextSizeClass: "text-base"},
		TimeSelectorSize:        SelectInputSizeSm,
		TimeSelector:            DateTimePickerTimeSelectorStyle{BackgroundColor: "secondary-950/40", DropdownBackgroundColor: "secondary-950"},
		Timezone:                "America/Sao_Paulo",
		TimezoneInputName:       "visitTimezone",
		TimezoneTwoWayStatePath: "visitTimezone",
		Style:                   DateTimePickerStyle{CalendarBackgroundColor: "secondary-900", CalendarBorderColor: "secondary-500/30", CalendarBorderRadius: uiToolset.BorderRadiusLg, CalendarButtonIconColor: "secondary-300"},
	})
	expectedFragments := []string{
		`&#34;timezone&#34;:&#34;America/Sao_Paulo&#34;`,
		`, visitTimezone)`,
		`>Timezone</legend>`,
		`ph-calendar-check text-sm`,
		`ph-seal-check`,
		` Apply <`,
		`text-secondary-300`,
		`bg-secondary-500`,
		`text-secondary-950`,
		`bg-secondary-950/40`,
		`bg-secondary-950`,
		`bg-secondary-900`,
		`border-secondary-500/30`,
		`rounded-lg`,
		`text-neutral-50/80 text-xs uppercase`,
		`name="visitTimezone"`,
		`:value="timezone"`,
		`visitTimezone = timezone`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
}

func TestDateTimePickerRendersLabelStatePath(testingHandle *testing.T) {
	renderedHtml := renderDateTimePicker(testingHandle, DateTimePickerSettings{
		Label:                "Booking date",
		LabelOneWayStatePath: "bookingLabel",
		InputName:            "booking",
		TwoWayStatePath:      "bookingDate",
	})
	for _, expectedFragment := range []string{
		`x-text="bookingLabel"`,
		`:aria-label="(bookingLabel) || &#34;Booking date&#34;"`,
		"Booking date",
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
}

func TestDateTimePickerRangeAriaLabelExpressionBuilder(testingHandle *testing.T) {
	testCases := []struct {
		name          string
		settings      DateTimePickerSettings
		sideLabel     string
		staticLabel   string
		expectedValue string
	}{
		{
			name:          "empty path yields no expression",
			settings:      DateTimePickerSettings{Label: "Booking date"},
			sideLabel:     "start",
			staticLabel:   "Booking date start",
			expectedValue: "",
		},
		{
			name:          "path falls back to the static aria label",
			settings:      DateTimePickerSettings{LabelOneWayStatePath: "bookingLabel"},
			sideLabel:     "end",
			staticLabel:   "Booking date end",
			expectedValue: `(bookingLabel ? (bookingLabel) + " end" : "Booking date end")`,
		},
	}
	for _, testCase := range testCases {
		testingHandle.Run(testCase.name, func(subtestHandle *testing.T) {
			actualValue := testCase.settings.rangeAriaLabelExpressionBuilder(
				testCase.sideLabel, testCase.staticLabel,
			)
			if actualValue != testCase.expectedValue {
				subtestHandle.Errorf(
					"RangeAriaLabelExpressionMismatch: got %q, want %q",
					actualValue, testCase.expectedValue,
				)
			}
		})
	}
}

func TestDateTimePickerRendersRangeLabelStatePath(testingHandle *testing.T) {
	renderedHtml := renderDateTimePicker(testingHandle, DateTimePickerSettings{
		Label:                   "Booking date",
		LabelOneWayStatePath:    "bookingLabel",
		IsRange:                 true,
		InputName:               "booking",
		TwoWayStatePath:         "bookingStartDate",
		RangeEndTwoWayStatePath: "bookingEndDate",
	})
	for _, expectedFragment := range []string{
		`aria-label="Booking date start"`,
		`aria-label="Booking date end"`,
		`:aria-label="(bookingLabel ? (bookingLabel) + &#34; start&#34; : &#34;Booking date start&#34;)"`,
		`:aria-label="(bookingLabel ? (bookingLabel) + &#34; end&#34; : &#34;Booking date end&#34;)"`,
	} {
		if !strings.Contains(renderedHtml, expectedFragment) {
			testingHandle.Errorf("rendered html missing fragment %s", expectedFragment)
		}
	}
}

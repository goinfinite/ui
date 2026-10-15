package uiForm

import (
	"bytes"
	"context"
	"strings"
	"testing"

	uiToolset "github.com/goinfinite/ui/src/toolset"
)

func TestCronScheduleRendersFieldAndPreviewStyle(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CronSchedule(CronScheduleSettings{
		Label: "Cron Schedule",
		Dropdown: CronScheduleDropdownStyle{
			BackgroundColor: "violet-900",
			BorderColor:     "violet-500/40",
			BorderRadius:    uiToolset.BorderRadiusLg,
		},
		Field: CronScheduleFieldStyle{
			BackgroundColor: "violet-950",
			BorderColor:     "violet-500/40",
			BorderRadius:    uiToolset.BorderRadiusLg,
		},
		Preview: CronSchedulePreviewStyle{
			CodeBackgroundColor: "violet-500/20",
			CodeTextColor:       "violet-200",
			TextColor:           "violet-100",
		},
		RadioStyle: RadioInputStyle{
			ActiveColor:    "violet-400",
			FocusRingColor: "violet-300",
			TextColor:      "violet-100",
		},
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CronScheduleRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		"bg-violet-950", "border-violet-500/40", "rounded-lg",
		"bg-violet-900",
		"bg-violet-500/20", "text-violet-200", "text-violet-100",
		"hover:text-violet-400", "text-violet-400 font-bold",
		"peer-focus-visible:ring-violet-300",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCronScheduleRendersSizeAndTextCaseOnPreview(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CronSchedule(CronScheduleSettings{
		Label:    "Cron Schedule",
		Size:     InputFieldSizeXl,
		TextCase: uiToolset.TextCaseUpper,
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CronScheduleRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`class="text-lg uppercase"`,
		"ph-bold ph-clock text-lg",
		"font-mono text-lg",
		"text-neutral-400 uppercase",
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCronSchedulePredefinedOptionsBuilder(t *testing.T) {
	options := cronSchedulePredefinedOptionsBuilder()
	if len(options) == 0 {
		t.Fatal("PredefinedOptionsEmpty")
	}
	seenValues := map[string]bool{}
	hasFiveDays := false
	for _, option := range options {
		if option.Label == "" {
			t.Errorf("PredefinedOptionMissingLabel: %q", option.Value)
		}
		if option.Value == "" {
			t.Errorf("PredefinedOptionMissingValue: %q", option.Label)
		}
		if seenValues[option.Value] {
			t.Errorf("PredefinedOptionValueRepeated: %q", option.Value)
		}
		seenValues[option.Value] = true
		if option.Label != "Every 5 days" {
			continue
		}
		hasFiveDays = true
		if option.Value != "0 0 */5 * *" {
			t.Errorf("EveryFiveDaysValueMismatch: %q", option.Value)
		}
	}
	if !hasFiveDays {
		t.Error("EveryFiveDaysOptionMissing")
	}
}

func TestCronScheduleFieldDescriptorsBuilder(t *testing.T) {
	expectedFieldNames := []string{"minute", "hour", "day", "month", "weekday"}
	descriptors := cronScheduleFieldDescriptorsBuilder()
	if len(descriptors) != len(expectedFieldNames) {
		t.Fatalf(
			"FieldDescriptorsLengthMismatch: got %d, want %d",
			len(descriptors), len(expectedFieldNames),
		)
	}
	for fieldIndex, descriptor := range descriptors {
		fieldName := expectedFieldNames[fieldIndex]
		expectedStatePath := "cronSchedule." + fieldName
		if descriptor.StatePath != expectedStatePath {
			t.Errorf(
				"FieldStatePathMismatch(%s): got %q, want %q",
				fieldName, descriptor.StatePath, expectedStatePath,
			)
		}
		if len(descriptor.Options) == 0 {
			t.Errorf("FieldOptionsEmpty(%s)", fieldName)
		}
		if !strings.Contains(descriptor.InvalidStateExpression, fieldName) {
			t.Errorf(
				"FieldInvalidExpressionMissingFieldName(%s): %q",
				fieldName, descriptor.InvalidStateExpression,
			)
		}
	}
	weekdayDescriptor := descriptors[4]
	hasWeekendsOption := false
	for _, option := range weekdayDescriptor.Options {
		if option.Label != "Only weekends" {
			continue
		}
		hasWeekendsOption = true
		if option.Value != "6,0" {
			t.Errorf("WeekendsValueMismatch: %q", option.Value)
		}
	}
	if !hasWeekendsOption {
		t.Error("WeekendsOptionMissing")
	}
}

func TestCronScheduleMonthAndWeekdayOrdinals(t *testing.T) {
	expectedMonths := map[string]string{
		"January": "1", "February": "2", "March": "3", "April": "4",
		"May": "5", "June": "6", "July": "7", "August": "8",
		"September": "9", "October": "10", "November": "11", "December": "12",
	}
	expectedWeekdays := map[string]string{
		"Every Monday": "1", "Every Tuesday": "2", "Every Wednesday": "3",
		"Every Thursday": "4", "Every Friday": "5", "Every Saturday": "6",
		"Every Sunday": "7",
	}

	assertOrdinalOptions(t, cronScheduleMonthOptionsBuilder(), expectedMonths)
	assertOrdinalOptions(t, cronScheduleWeekdayOptionsBuilder(), expectedWeekdays)
}

func assertOrdinalOptions(
	t *testing.T, options []SelectLabelValueOption, expected map[string]string,
) {
	t.Helper()
	matched := map[string]bool{}
	for _, option := range options {
		expectedValue, isOrdinal := expected[option.Label]
		if !isOrdinal {
			continue
		}
		matched[option.Label] = true
		if option.Value != expectedValue {
			t.Errorf("OrdinalValueMismatch(%s): got %q, want %q", option.Label, option.Value, expectedValue)
		}
	}
	for label := range expected {
		if !matched[label] {
			t.Errorf("OrdinalOptionMissing(%s)", label)
		}
	}
}

func TestCronScheduleRendersFieldsAndBoundState(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CronSchedule(CronScheduleSettings{
		Label:           "Cron Schedule",
		TwoWayStatePath: "schedule",
		InputName:       "schedule",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CronScheduleRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	expectedFragments := []string{
		`x-data="cronScheduleInput(`,
		`x-init="cronScheduleInit(schedule)"`,
		`x-effect="schedule = cronExpression"`,
		`x-effect="cronScheduleMode === 'custom' && (cronExpression = cronCustomExpressionBuilder())"`,
		`name="schedule"`,
		`x-text="cronPreviewTextBuilder()"`,
		`Cron Schedule`,
		`Predefined Schedule`,
		`cronSchedule.minute`,
		`cronSchedule.hour`,
		`cronSchedule.day`,
		`cronSchedule.month`,
		`cronSchedule.weekday`,
		`0 0 */5 * *`,
	}
	for _, expectedFragment := range expectedFragments {
		if !strings.Contains(renderedHtml, expectedFragment) {
			t.Errorf("RenderedHtmlMissingFragment: %q", expectedFragment)
		}
	}
}

func TestCronScheduleOmitsPathBindingWhenUnbound(t *testing.T) {
	var buffer bytes.Buffer
	renderErr := CronSchedule(CronScheduleSettings{
		Label: "Cron Schedule",
	}).Render(context.Background(), &buffer)
	if renderErr != nil {
		t.Fatalf("CronScheduleRenderFailed: %v", renderErr)
	}
	renderedHtml := buffer.String()
	if strings.Contains(renderedHtml, `x-init="cronScheduleInit(`) {
		t.Error("UnboundCronScheduleRendersPathInit")
	}
	if strings.Contains(renderedHtml, `type="hidden"`) {
		t.Error("UnboundCronScheduleRendersHiddenInput")
	}
}

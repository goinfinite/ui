package uiForm

import "strconv"

type cronScheduleFieldDescriptor struct {
	Label                  string
	StatePath              string
	Options                []SelectLabelValueOption
	InvalidStateExpression string
}

func cronSchedulePredefinedOptionsBuilder() []SelectLabelValueOption {
	return []SelectLabelValueOption{
		{Label: "Hourly", Value: "@hourly"},
		{Label: "Every 3 hours", Value: "0 */3 * * *"},
		{Label: "Every 6 hours", Value: "0 */6 * * *"},
		{Label: "Every 12 hours", Value: "0 */12 * * *"},
		{Label: "Daily at midnight", Value: "0 0 * * *"},
		{Label: "Daily at 1 AM", Value: "0 1 * * *"},
		{Label: "Daily at 2 AM", Value: "0 2 * * *"},
		{Label: "Every 2 days", Value: "0 0 */2 * *"},
		{Label: "Every 3 days", Value: "0 0 */3 * *"},
		{Label: "Every 5 days", Value: "0 0 */5 * *"},
		{Label: "Weekly", Value: "@weekly"},
		{Label: "Monthly", Value: "@monthly"},
		{Label: "Annually", Value: "@annually"},
		{Label: "On Reboot", Value: "@reboot"},
	}
}

func cronScheduleNumberedOptionsBuilder(
	labelPrefix string, values []int,
) []SelectLabelValueOption {
	options := make([]SelectLabelValueOption, 0, len(values))
	for _, value := range values {
		valueText := strconv.Itoa(value)
		options = append(options, SelectLabelValueOption{
			Label: labelPrefix + valueText,
			Value: valueText,
		})
	}
	return options
}

func cronScheduleSteppedOptionsBuilder(
	labelPrefix, labelSuffix string, steps []int,
) []SelectLabelValueOption {
	options := make([]SelectLabelValueOption, 0, len(steps))
	for _, step := range steps {
		stepText := strconv.Itoa(step)
		options = append(options, SelectLabelValueOption{
			Label: labelPrefix + stepText + labelSuffix,
			Value: "*/" + stepText,
		})
	}
	return options
}

func cronScheduleOrdinalOptionsBuilder(
	labelPrefix string, names []string,
) []SelectLabelValueOption {
	options := make([]SelectLabelValueOption, 0, len(names))
	for nameIndex, name := range names {
		options = append(options, SelectLabelValueOption{
			Label: labelPrefix + name,
			Value: strconv.Itoa(nameIndex + 1),
		})
	}
	return options
}

func cronScheduleMinuteOptionsBuilder() []SelectLabelValueOption {
	options := cronScheduleNumberedOptionsBuilder("At minute ", []int{0, 5, 10, 15, 30, 45})
	options = append(options, SelectLabelValueOption{Label: "Every minute", Value: "*"})
	return append(options, cronScheduleSteppedOptionsBuilder(
		"Every ", " minutes", []int{5, 10, 15, 30, 45},
	)...)
}

func cronScheduleHourOptionsBuilder() []SelectLabelValueOption {
	options := []SelectLabelValueOption{{Label: "Every hour", Value: "*"}}
	options = append(options, cronScheduleSteppedOptionsBuilder(
		"Every ", " hours", []int{2, 6, 12},
	)...)
	return append(options,
		SelectLabelValueOption{Label: "At midnight", Value: "0"},
		SelectLabelValueOption{Label: "At 2h00 AM", Value: "2"},
		SelectLabelValueOption{Label: "At 6h00 AM", Value: "6"},
		SelectLabelValueOption{Label: "At midday", Value: "12"},
		SelectLabelValueOption{Label: "At 14h00 (2h00 PM)", Value: "14"},
		SelectLabelValueOption{Label: "At 18h00 (6h00 PM)", Value: "18"},
	)
}

func cronScheduleDayOptionsBuilder() []SelectLabelValueOption {
	options := []SelectLabelValueOption{
		{Label: "Every day", Value: "*"},
		{Label: "Weekly", Value: "*/7"},
		{Label: "Every 15 days", Value: "*/15"},
		{Label: "Every 30 days", Value: "*/30"},
	}
	return append(options, cronScheduleNumberedOptionsBuilder(
		"At day ", []int{1, 5, 10, 15, 20, 25, 30},
	)...)
}

func cronScheduleMonthOptionsBuilder() []SelectLabelValueOption {
	options := []SelectLabelValueOption{
		{Label: "Every month", Value: "*"},
		{Label: "Every 2 months", Value: "*/2"},
		{Label: "Every 6 months", Value: "*/6"},
		{Label: "Every year", Value: "*/12"},
	}
	return append(options, cronScheduleOrdinalOptionsBuilder("", []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	})...)
}

func cronScheduleWeekdayOptionsBuilder() []SelectLabelValueOption {
	options := []SelectLabelValueOption{
		{Label: "Every day", Value: "*"},
		{Label: "Monday to Friday", Value: "1-5"},
		{Label: "Only weekends", Value: "6,0"},
	}
	return append(options, cronScheduleOrdinalOptionsBuilder("Every ", []string{
		"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
	})...)
}

func cronScheduleFieldDescriptorBuilder(
	label, fieldName string, options []SelectLabelValueOption,
) cronScheduleFieldDescriptor {
	return cronScheduleFieldDescriptor{
		Label:                  label,
		StatePath:              "cronSchedule." + fieldName,
		Options:                options,
		InvalidStateExpression: "!cronFieldIsValid('" + fieldName + "', cronSchedule." + fieldName + ")",
	}
}

func cronScheduleFieldDescriptorsBuilder() []cronScheduleFieldDescriptor {
	return []cronScheduleFieldDescriptor{
		cronScheduleFieldDescriptorBuilder("Minute", "minute", cronScheduleMinuteOptionsBuilder()),
		cronScheduleFieldDescriptorBuilder("Hour", "hour", cronScheduleHourOptionsBuilder()),
		cronScheduleFieldDescriptorBuilder("Day", "day", cronScheduleDayOptionsBuilder()),
		cronScheduleFieldDescriptorBuilder("Month", "month", cronScheduleMonthOptionsBuilder()),
		cronScheduleFieldDescriptorBuilder("Weekday", "weekday", cronScheduleWeekdayOptionsBuilder()),
	}
}

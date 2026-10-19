package uiForm

import (
	"fmt"
	"slices"
	"strings"
)

var dateTimePickerHourOptions = dateTimePickerZeroPaddedOptionsResolver(24)
var dateTimePickerMinuteOptions = dateTimePickerZeroPaddedOptionsResolver(60)

var dateTimePickerDefaultTimezones = []string{
	"UTC",
	"Africa/Cairo",
	"Africa/Johannesburg",
	"Africa/Lagos",
	"Africa/Nairobi",
	"America/Argentina/Buenos_Aires",
	"America/Bogota",
	"America/Chicago",
	"America/Denver",
	"America/Lima",
	"America/Los_Angeles",
	"America/Mexico_City",
	"America/New_York",
	"America/Santiago",
	"America/Sao_Paulo",
	"Asia/Bangkok",
	"Asia/Dubai",
	"Asia/Hong_Kong",
	"Asia/Kolkata",
	"Asia/Seoul",
	"Asia/Shanghai",
	"Asia/Singapore",
	"Asia/Tokyo",
	"Australia/Brisbane",
	"Australia/Melbourne",
	"Australia/Perth",
	"Australia/Sydney",
	"Europe/Amsterdam",
	"Europe/Athens",
	"Europe/Berlin",
	"Europe/Brussels",
	"Europe/Copenhagen",
	"Europe/Helsinki",
	"Europe/Istanbul",
	"Europe/Lisbon",
	"Europe/London",
	"Europe/Madrid",
	"Europe/Moscow",
	"Europe/Oslo",
	"Europe/Paris",
	"Europe/Rome",
	"Europe/Stockholm",
	"Europe/Vienna",
	"Europe/Warsaw",
	"Europe/Zurich",
	"Pacific/Auckland",
}

func dateTimePickerZeroPaddedOptionsResolver(limit int) []string {
	options := make([]string, limit)
	for index := range options {
		options[index] = fmt.Sprintf("%02d", index)
	}
	return options
}

func dateTimePickerTimezoneLabelFormatter(timezone string) string {
	return strings.ReplaceAll(timezone, "_", " ")
}

func dateTimePickerTimezoneOptionsResolver(
	selectedTimezone string, configuredOptions []string,
) []SelectLabelValueOption {
	timezones := configuredOptions
	if len(timezones) == 0 {
		timezones = dateTimePickerDefaultTimezones
	}
	if selectedTimezone != "" && !slices.Contains(timezones, selectedTimezone) {
		timezones = append([]string{selectedTimezone}, timezones...)
	}
	options := make([]SelectLabelValueOption, 0, len(timezones))
	for _, timezone := range timezones {
		options = append(options, SelectLabelValueOption{
			Label: dateTimePickerTimezoneLabelFormatter(timezone),
			Value: timezone,
		})
	}
	return options
}

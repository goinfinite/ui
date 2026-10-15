package uiForm

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

func cronScheduleFieldDescriptorsBuilder() []cronScheduleFieldDescriptor {
	return []cronScheduleFieldDescriptor{
		{
			Label:     "Minute",
			StatePath: "cronSchedule.minute",
			Options: []SelectLabelValueOption{
				{Label: "At minute 0", Value: "0"},
				{Label: "At minute 5", Value: "5"},
				{Label: "At minute 10", Value: "10"},
				{Label: "At minute 15", Value: "15"},
				{Label: "At minute 30", Value: "30"},
				{Label: "At minute 45", Value: "45"},
				{Label: "Every minute", Value: "*"},
				{Label: "Every 5 minutes", Value: "*/5"},
				{Label: "Every 10 minutes", Value: "*/10"},
				{Label: "Every 15 minutes", Value: "*/15"},
				{Label: "Every 30 minutes", Value: "*/30"},
				{Label: "Every 45 minutes", Value: "*/45"},
			},
			InvalidStateExpression: "!cronFieldIsValid('minute', cronSchedule.minute)",
		},
		{
			Label:     "Hour",
			StatePath: "cronSchedule.hour",
			Options: []SelectLabelValueOption{
				{Label: "Every hour", Value: "*"},
				{Label: "Every 2 hours", Value: "*/2"},
				{Label: "Every 6 hours", Value: "*/6"},
				{Label: "Every 12 hours", Value: "*/12"},
				{Label: "At midnight", Value: "0"},
				{Label: "At 2h00 AM", Value: "2"},
				{Label: "At 6h00 AM", Value: "6"},
				{Label: "At midday", Value: "12"},
				{Label: "At 14h00 (2h00 PM)", Value: "14"},
				{Label: "At 18h00 (6h00 PM)", Value: "18"},
			},
			InvalidStateExpression: "!cronFieldIsValid('hour', cronSchedule.hour)",
		},
		{
			Label:     "Day",
			StatePath: "cronSchedule.day",
			Options: []SelectLabelValueOption{
				{Label: "Every day", Value: "*"},
				{Label: "Weekly", Value: "*/7"},
				{Label: "Every 15 days", Value: "*/15"},
				{Label: "Every 30 days", Value: "*/30"},
				{Label: "At day 1", Value: "1"},
				{Label: "At day 5", Value: "5"},
				{Label: "At day 10", Value: "10"},
				{Label: "At day 15", Value: "15"},
				{Label: "At day 20", Value: "20"},
				{Label: "At day 25", Value: "25"},
				{Label: "At day 30", Value: "30"},
			},
			InvalidStateExpression: "!cronFieldIsValid('day', cronSchedule.day)",
		},
		{
			Label:     "Month",
			StatePath: "cronSchedule.month",
			Options: []SelectLabelValueOption{
				{Label: "Every month", Value: "*"},
				{Label: "Every 2 months", Value: "*/2"},
				{Label: "Every 6 months", Value: "*/6"},
				{Label: "Every year", Value: "*/12"},
				{Label: "January", Value: "1"},
				{Label: "February", Value: "2"},
				{Label: "March", Value: "3"},
				{Label: "April", Value: "4"},
				{Label: "May", Value: "5"},
				{Label: "June", Value: "6"},
				{Label: "July", Value: "7"},
				{Label: "August", Value: "8"},
				{Label: "September", Value: "9"},
				{Label: "October", Value: "10"},
				{Label: "November", Value: "11"},
				{Label: "December", Value: "12"},
			},
			InvalidStateExpression: "!cronFieldIsValid('month', cronSchedule.month)",
		},
		{
			Label:     "Weekday",
			StatePath: "cronSchedule.weekday",
			Options: []SelectLabelValueOption{
				{Label: "Every day", Value: "*"},
				{Label: "Monday to Friday", Value: "1-5"},
				{Label: "Only weekends", Value: "6,0"},
				{Label: "Every Monday", Value: "1"},
				{Label: "Every Tuesday", Value: "2"},
				{Label: "Every Wednesday", Value: "3"},
				{Label: "Every Thursday", Value: "4"},
				{Label: "Every Friday", Value: "5"},
				{Label: "Every Saturday", Value: "6"},
				{Label: "Every Sunday", Value: "7"},
			},
			InvalidStateExpression: "!cronFieldIsValid('weekday', cronSchedule.weekday)",
		},
	}
}

UiToolset.RegisterAlpineState(() => {
  const fieldNames = ["minute", "hour", "day", "month", "weekday"];
  const fieldBounds = {
    minute: [0, 59],
    hour: [0, 23],
    day: [1, 31],
    month: [1, 12],
    weekday: [0, 7],
  };
  const weekdayNames = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
  ];
  const monthNames = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];
  const macroPreviews = {
    "@hourly": "Every hour on the hour",
    "@daily": "Daily at midnight",
    "@weekly": "Every Sunday at midnight",
    "@monthly": "On day 1 of every month at midnight",
    "@yearly": "On January 1 at midnight",
    "@annually": "On January 1 at midnight",
    "@reboot": "When the system reboots",
  };

  Alpine.data("cronScheduleInput", (predefinedValues) => ({
    cronScheduleMode: "predefined",
    cronExpression: "0 0 * * *",
    cronSchedule: {
      minute: "0",
      hour: "0",
      day: "*",
      month: "*",
      weekday: "*",
    },
    predefinedValues: predefinedValues || [],

    cronScheduleInit(initialExpression) {
      const expression = String(initialExpression ?? "").trim();
      if (expression === "") {
        return;
      }
      this.cronExpression = expression;
      if (
        expression.startsWith("@") ||
        this.predefinedValues.includes(expression)
      ) {
        this.cronScheduleMode = "predefined";
        return;
      }
      const parts = expression.split(/\s+/);
      if (parts.length !== 5) {
        return;
      }
      this.cronSchedule = {
        minute: parts[0],
        hour: parts[1],
        day: parts[2],
        month: parts[3],
        weekday: parts[4],
      };
      this.cronScheduleMode = "custom";
    },

    cronFieldIsValid(fieldName, fieldValue) {
      const bounds = fieldBounds[fieldName];
      if (bounds === undefined) {
        return false;
      }
      const fieldText = String(fieldValue ?? "").trim();
      if (fieldText === "") {
        return false;
      }
      return fieldText
        .split(",")
        .every((part) =>
          this.cronFieldPartIsValid(part.trim(), bounds[0], bounds[1]),
        );
    },

    cronFieldPartIsValid(part, minValue, maxValue) {
      if (part === "*") {
        return true;
      }
      const stepMatch = part.match(/^(\*|\d{1,2}|\d{1,2}-\d{1,2})\/(\d{1,3})$/);
      if (stepMatch !== null) {
        if (Number(stepMatch[2]) < 1) {
          return false;
        }
        return this.cronFieldPartIsValid(stepMatch[1], minValue, maxValue);
      }
      const rangeMatch = part.match(/^(\d{1,2})-(\d{1,2})$/);
      if (rangeMatch !== null) {
        const rangeStartsAboveEnd =
          Number(rangeMatch[1]) > Number(rangeMatch[2]);
        if (rangeStartsAboveEnd) {
          return false;
        }
        return (
          this.cronFieldNumberInBounds(rangeMatch[1], minValue, maxValue) &&
          this.cronFieldNumberInBounds(rangeMatch[2], minValue, maxValue)
        );
      }
      return this.cronFieldNumberInBounds(part, minValue, maxValue);
    },

    cronFieldNumberInBounds(numberText, minValue, maxValue) {
      if (!/^\d{1,2}$/.test(numberText)) {
        return false;
      }
      const number = Number(numberText);
      return number >= minValue && number <= maxValue;
    },

    cronScheduleIsValid() {
      const expression = String(this.cronExpression ?? "").trim();
      if (Object.hasOwn(macroPreviews, expression)) {
        return true;
      }
      const parts = expression.split(/\s+/);
      if (parts.length !== 5) {
        return false;
      }
      return fieldNames.every((fieldName, fieldIndex) =>
        this.cronFieldIsValid(fieldName, parts[fieldIndex]),
      );
    },

    cronCustomExpressionBuilder() {
      const parts = fieldNames.map((fieldName) => this.cronSchedule[fieldName]);
      const allFieldsValid = fieldNames.every((fieldName, fieldIndex) =>
        this.cronFieldIsValid(fieldName, parts[fieldIndex]),
      );
      if (!allFieldsValid) {
        return "";
      }
      return parts.join(" ");
    },

    cronPreviewTextBuilder() {
      const expression = String(this.cronExpression ?? "").trim();
      if (expression === "") {
        return "Invalid schedule: fix the highlighted fields";
      }
      const macroPreview = macroPreviews[expression];
      if (macroPreview !== undefined) {
        return macroPreview;
      }
      const parts = expression.split(/\s+/);
      if (parts.length !== 5) {
        return "Invalid schedule: expected five fields";
      }
      const allFieldsValid = fieldNames.every((fieldName, fieldIndex) =>
        this.cronFieldIsValid(fieldName, parts[fieldIndex]),
      );
      if (!allFieldsValid) {
        return "Invalid schedule: fix the highlighted fields";
      }
      const timePhrase = this.cronTimePhraseBuilder(parts[0], parts[1]);
      const dayPhrase = this.cronDayPhraseBuilder(parts[2], parts[3], parts[4]);
      return `${timePhrase}, ${dayPhrase}`;
    },

    cronTimePhraseBuilder(minute, hour) {
      if (minute === "*" && hour === "*") {
        return "Every minute";
      }
      if (minute === "*") {
        return `Every minute during hour ${hour}`;
      }
      const minuteStepMatch = minute.match(/^\*\/(\d{1,3})$/);
      if (minuteStepMatch !== null && hour === "*") {
        return `Every ${minuteStepMatch[1]} minutes`;
      }
      const hourStepMatch = hour.match(/^\*\/(\d{1,3})$/);
      if (/^\d{1,2}$/.test(minute) && hourStepMatch !== null) {
        return `Every ${hourStepMatch[1]} hours at minute ${minute}`;
      }
      if (/^\d{1,2}$/.test(minute) && /^\d{1,2}$/.test(hour)) {
        return `At ${this.cronNumberPaddedFormatter(hour)}:${this.cronNumberPaddedFormatter(minute)}`;
      }
      if (hour === "*") {
        return `At minute ${minute} of every hour`;
      }
      return `At minute ${minute} of hour ${hour}`;
    },

    cronDayPhraseBuilder(day, month, weekday) {
      if (day === "*" && month === "*" && weekday === "*") {
        return "every day";
      }
      if (day === "*" && month === "*") {
        if (weekday === "1-5") {
          return "on weekdays";
        }
        if (weekday === "6,0" || weekday === "0,6") {
          return "on weekends";
        }
        if (/^\d{1,2}$/.test(weekday)) {
          return `on ${this.cronWeekdayNameResolver(weekday)}`;
        }
        return `on weekday ${weekday}`;
      }
      const dayStepMatch = day.match(/^\*\/(\d{1,3})$/);
      if (dayStepMatch !== null && month === "*" && weekday === "*") {
        return `every ${dayStepMatch[1]} days`;
      }
      if (/^\d{1,2}$/.test(day) && /^\d{1,2}$/.test(month) && weekday === "*") {
        return `on ${this.cronMonthNameResolver(month)} ${day}`;
      }
      if (/^\d{1,2}$/.test(day) && month === "*" && weekday === "*") {
        return `on day ${day} of every month`;
      }
      if (day === "*" && /^\d{1,2}$/.test(month) && weekday === "*") {
        return `in ${this.cronMonthNameResolver(month)}`;
      }
      return `on day ${day}, month ${month}, weekday ${weekday}`;
    },

    cronNumberPaddedFormatter(numberText) {
      return String(numberText).padStart(2, "0");
    },

    cronWeekdayNameResolver(weekday) {
      return weekdayNames[Number(weekday) % 7];
    },

    cronMonthNameResolver(month) {
      return monthNames[Number(month) - 1];
    },
  }));
});

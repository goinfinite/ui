UiToolset.RegisterAlpineState(() => {
  const weekStartByWeekdayName = {
    monday: 1,
    sunday: 0,
  };

  const referenceSunday = new Date(2024, 0, 7);

  const dateTimePickerStateFactory = (
    config,
    startValue,
    endValue,
    timezoneValue,
  ) => {
    const options = config || {};
    const styleOptions = options.style || {};
    const dayStyle = {
      disabledTextColor: styleOptions.dayDisabledTextColor || "neutral-600",
      hoverBackgroundColor:
        styleOptions.dayHoverBackgroundColor || "neutral-50/10",
      mutedTextColor: styleOptions.dayMutedTextColor || "neutral-500",
      rangeBackgroundColor:
        styleOptions.dayRangeBackgroundColor || "secondary-500/20",
      rangeTextColor: styleOptions.dayRangeTextColor || "neutral-50",
      selectedBackgroundColor:
        styleOptions.daySelectedBackgroundColor || "secondary-500",
      selectedTextColor: styleOptions.daySelectedTextColor || "neutral-950",
      textColor: styleOptions.dayTextColor || "neutral-100",
      todayRingColor: styleOptions.dayTodayRingColor || "neutral-50/30",
    };
    const startText = String(startValue ?? "").trim();
    const endText = String(endValue ?? "").trim();
    const configuredTimezone = String(options.timezone ?? "").trim();
    const boundTimezone = String(timezoneValue ?? "").trim();
    return {
      mode: options.mode || "date",
      locale: options.locale || "en-US",
      minValue: options.minValue || "",
      maxValue: options.maxValue || "",
      firstDayOfWeekOverride: options.firstDayOfWeek || "",
      isRange: Boolean(options.isRange),
      timezone: boundTimezone || configuredTimezone,
      startDraft: startText,
      endDraft: endText,
      startTypedValue: "",
      endTypedValue: "",
      activeEnd: "start",
      viewYear: 2000,
      viewMonth: 0,
      isOpen: false,
      initialized: false,
      timeHour: "00",
      timeMinute: "00",

      init() {
        const seedDate = this.dateObjectFromIso(this.startDraft) || new Date();
        this.viewYear = seedDate.getFullYear();
        this.viewMonth = seedDate.getMonth();
        this.syncTimeSelectsFromDraft();
        this.syncTypedValues();
        this.initialized = true;
      },

      get hasValue() {
        return this.startDraft !== "" || this.endDraft !== "";
      },

      get activeDraft() {
        return this.activeEnd === "end" ? this.endDraft : this.startDraft;
      },

      get firstDayOfWeek() {
        if (this.firstDayOfWeekOverride !== "") {
          const override =
            weekStartByWeekdayName[this.firstDayOfWeekOverride.toLowerCase()];
          if (override !== undefined) {
            return override;
          }
        }
        try {
          const locale = new Intl.Locale(this.locale);
          const weekInfo =
            typeof locale.getWeekInfo === "function"
              ? locale.getWeekInfo()
              : locale.weekInfo;
          if (weekInfo && typeof weekInfo.firstDay === "number") {
            return weekInfo.firstDay % 7;
          }
        } catch (localeError) {
          console.warn("DateTimePickerLocaleWeekInfoUnavailable", localeError);
        }
        return 1;
      },

      get weekdayLabels() {
        const formatter = new Intl.DateTimeFormat(this.locale, {
          weekday: "short",
        });
        const labels = [];
        for (let offset = 0; offset < 7; offset++) {
          const weekdayIndex = (this.firstDayOfWeek + offset) % 7;
          const weekdayDate = new Date(
            referenceSunday.getFullYear(),
            referenceSunday.getMonth(),
            referenceSunday.getDate() + weekdayIndex,
          );
          labels.push(formatter.format(weekdayDate));
        }
        return labels;
      },

      get monthYearLabel() {
        const viewDate = new Date(this.viewYear, this.viewMonth, 1);
        return new Intl.DateTimeFormat(this.locale, {
          month: "long",
          year: "numeric",
        }).format(viewDate);
      },

      get calendarDays() {
        const firstOfMonth = new Date(this.viewYear, this.viewMonth, 1);
        const leadingOffset =
          (firstOfMonth.getDay() - this.firstDayOfWeek + 7) % 7;
        const gridStart = new Date(
          this.viewYear,
          this.viewMonth,
          1 - leadingOffset,
        );
        const startIso = this.startDraft
          ? this.datePartExtractor(this.startDraft)
          : "";
        const endIso =
          this.isRange && this.endDraft
            ? this.datePartExtractor(this.endDraft)
            : "";
        const todayIso = this.isoFromDateObject(new Date());
        const days = [];
        for (let gridIndex = 0; gridIndex < 42; gridIndex++) {
          const date = new Date(
            gridStart.getFullYear(),
            gridStart.getMonth(),
            gridStart.getDate() + gridIndex,
          );
          const iso = this.isoFromDateObject(date);
          const isRangeStart =
            this.isRange && startIso !== "" && iso === startIso;
          const isRangeEnd = this.isRange && endIso !== "" && iso === endIso;
          const isInRange =
            this.isRange &&
            startIso !== "" &&
            endIso !== "" &&
            iso > startIso &&
            iso < endIso;
          days.push({
            iso,
            dayNumber: date.getDate(),
            isCurrentMonth: date.getMonth() === this.viewMonth,
            isToday: iso === todayIso,
            isRangeStart,
            isRangeEnd,
            isInRange,
            isSelected:
              isRangeStart || isRangeEnd || (!this.isRange && iso === startIso),
            isDisabled: this.dayIsOutOfBounds(iso),
          });
        }
        return days;
      },

      selectDay(day) {
        if (day.isDisabled) {
          return;
        }
        const candidateDraft = this.composeDateDraft(day.iso);
        if (this.valueIsOutOfBounds(candidateDraft)) {
          return;
        }
        this.assignDraft(candidateDraft);
        if (this.isRange) {
          this.normalizeRangeOrder();
          this.activeEnd = this.activeEnd === "start" ? "end" : "start";
          this.syncTimeSelectsFromDraft();
        }
        this.syncTypedValues();
      },

      activateEnd(end) {
        this.activeEnd = end;
        this.syncTimeSelectsFromDraft();
        this.jumpViewToActiveDraft();
        this.isOpen = true;
      },

      commitTime() {
        const timeText = `${this.timeHour}:${this.timeMinute}`;
        if (this.mode === "time") {
          if (this.valueIsOutOfBounds(timeText)) {
            this.syncTimeSelectsFromDraft();
            return;
          }
          this.assignDraft(timeText);
          if (this.isRange) {
            this.normalizeRangeOrder();
          }
          this.syncTypedValues();
          return;
        }
        const dateText =
          this.datePartExtractor(this.activeDraft) ||
          this.isoFromDateObject(new Date());
        const candidateDraft = this.composeDateDraft(dateText);
        if (this.valueIsOutOfBounds(candidateDraft)) {
          this.syncTimeSelectsFromDraft();
          return;
        }
        this.assignDraft(candidateDraft);
        if (this.isRange) {
          this.normalizeRangeOrder();
        }
        this.syncTypedValues();
      },

      commitTimezone() {
        if (this.mode !== "datetime") {
          return;
        }
        this.startDraft = this.dateTimezoneComposer(this.startDraft);
        this.endDraft = this.dateTimezoneComposer(this.endDraft);
        this.syncTimeSelectsFromDraft();
        this.syncTypedValues();
      },

      commitTypedValueForSide(side) {
        const typedText = String(
          side === "end" ? this.endTypedValue : this.startTypedValue,
        ).trim();
        if (
          typedText === this.draftDisplayTextFormatter(this.sideDraft(side))
        ) {
          return;
        }
        if (typedText === "") {
          this.assignSideDraft(side, "");
          this.syncTypedValues();
          return;
        }
        const parsedValue = this.parseIsoInput(typedText);
        if (parsedValue === null) {
          this.syncTypedValues();
          return;
        }
        this.assignSideDraft(side, parsedValue);
        if (this.isRange) {
          this.normalizeRangeOrder();
        }
        this.syncTimeSelectsFromDraft();
        this.syncTypedValues();
      },

      togglePopover() {
        if (this.isOpen) {
          this.closePopover();
          return;
        }
        this.openPopover();
      },

      openPopover() {
        this.isOpen = true;
        this.jumpViewToActiveDraft();
      },

      closePopover() {
        this.isOpen = false;
      },

      shiftMonth(offset) {
        const target = new Date(this.viewYear, this.viewMonth + offset, 1);
        this.viewYear = target.getFullYear();
        this.viewMonth = target.getMonth();
      },

      shiftYear(offset) {
        this.viewYear += offset;
      },

      dayFocusTargetResolver(event) {
        const dayButtons = event.currentTarget.querySelectorAll(
          "button[data-day-index]",
        );
        if (dayButtons.length === 0) {
          return null;
        }
        return {
          dayButtons,
          currentIndex: Number(event.target.dataset.dayIndex ?? 0),
        };
      },

      focusDayOffset(event, offset) {
        const dayFocusTarget = this.dayFocusTargetResolver(event);
        if (dayFocusTarget === null) {
          return;
        }
        const { dayButtons, currentIndex } = dayFocusTarget;
        const targetIndex = Math.min(
          Math.max(currentIndex + offset, 0),
          dayButtons.length - 1,
        );
        dayButtons[targetIndex].focus();
      },

      focusDayEdge(event, edge) {
        const dayFocusTarget = this.dayFocusTargetResolver(event);
        if (dayFocusTarget === null) {
          return;
        }
        const { dayButtons, currentIndex } = dayFocusTarget;
        const weekStartIndex = currentIndex - (currentIndex % 7);
        const weekEndIndex = Math.min(
          weekStartIndex + 6,
          dayButtons.length - 1,
        );
        const step = edge === "start" ? 1 : -1;
        let targetIndex = edge === "start" ? weekStartIndex : weekEndIndex;
        while (
          targetIndex >= weekStartIndex &&
          targetIndex <= weekEndIndex &&
          dayButtons[targetIndex].disabled
        ) {
          targetIndex += step;
        }
        if (targetIndex < weekStartIndex || targetIndex > weekEndIndex) {
          return;
        }
        dayButtons[targetIndex].focus();
      },

      dayClassesBuilder(day) {
        const classes = [
          "group",
          "relative",
          "flex",
          "h-9",
          "w-full",
          "items-center",
          "justify-center",
          "text-sm",
          "transition-colors",
        ];
        if (day.isDisabled) {
          classes.push(
            "cursor-not-allowed",
            "bg-transparent",
            `text-${dayStyle.disabledTextColor}`,
          );
          return classes.join(" ");
        }
        classes.push("cursor-pointer");
        if (day.isInRange) {
          classes.push(
            `bg-${dayStyle.rangeBackgroundColor}`,
            `text-${dayStyle.rangeTextColor}`,
          );
        }
        if (day.isRangeStart) {
          classes.push(
            "rounded-l-md",
            `bg-${dayStyle.selectedBackgroundColor}`,
            "font-bold",
            `text-${dayStyle.selectedTextColor}`,
          );
        }
        if (day.isRangeEnd) {
          classes.push(
            "rounded-r-md",
            `bg-${dayStyle.selectedBackgroundColor}`,
            "font-bold",
            `text-${dayStyle.selectedTextColor}`,
          );
        }
        if (day.isSelected && !this.isRange) {
          classes.push(
            "rounded-md",
            `bg-${dayStyle.selectedBackgroundColor}`,
            "font-bold",
            `text-${dayStyle.selectedTextColor}`,
          );
        }
        if (!day.isInRange && !day.isSelected) {
          classes.push("bg-transparent");
          classes.push(
            day.isCurrentMonth
              ? `text-${dayStyle.textColor}`
              : `text-${dayStyle.mutedTextColor}`,
          );
        }
        return classes.join(" ");
      },

      dayNumberClassesBuilder(day) {
        const classes = [
          "flex",
          "h-7",
          "w-7",
          "items-center",
          "justify-center",
          "rounded-full",
          "transition-colors",
        ];
        if (!day.isInRange && !day.isSelected && !day.isDisabled) {
          classes.push(`group-hover:bg-${dayStyle.hoverBackgroundColor}`);
        }
        if (day.isToday) {
          classes.push("ring-1", `ring-${dayStyle.todayRingColor}`);
        }
        return classes.join(" ");
      },

      dayAriaLabelFormatter(day) {
        const date = this.dateObjectFromIso(day.iso);
        if (date === null) {
          return day.iso;
        }
        return new Intl.DateTimeFormat(this.locale, {
          weekday: "long",
          year: "numeric",
          month: "long",
          day: "numeric",
        }).format(date);
      },

      timezoneOffsetSuffix(wallClockIso) {
        if (this.timezone === "" || this.mode !== "datetime") {
          return "";
        }
        return this.timezoneOffsetResolver(wallClockIso);
      },

      timezoneOffsetMinutes(date) {
        try {
          const offsetParts = new Intl.DateTimeFormat("en-US", {
            timeZone: this.timezone,
            timeZoneName: "longOffset",
          }).formatToParts(date);
          const timezoneName =
            offsetParts.find((part) => part.type === "timeZoneName")?.value ??
            "";
          if (timezoneName === "GMT" || timezoneName === "UTC") {
            return 0;
          }
          const offsetMatch = timezoneName.match(
            /^GMT([+-])(\d{1,2})(?::(\d{2}))?$/,
          );
          if (offsetMatch === null) {
            return null;
          }
          const offsetSign = offsetMatch[1] === "-" ? -1 : 1;
          const offsetHours = Number(offsetMatch[2]);
          const offsetMinutes = Number(offsetMatch[3] ?? "00");
          return offsetSign * (offsetHours * 60 + offsetMinutes);
        } catch (timezoneError) {
          console.warn(
            "DateTimePickerTimezoneOffsetUnavailable",
            timezoneError,
          );
          return null;
        }
      },

      timezoneOffsetResolver(wallClockIso) {
        const wallClockDate =
          this.dateObjectFromIso(wallClockIso) || new Date();
        const wallClockUtcMs = Date.UTC(
          wallClockDate.getFullYear(),
          wallClockDate.getMonth(),
          wallClockDate.getDate(),
          wallClockDate.getHours(),
          wallClockDate.getMinutes(),
        );
        const initialOffsetMinutes = this.timezoneOffsetMinutes(
          new Date(wallClockUtcMs),
        );
        if (initialOffsetMinutes === null) {
          return "";
        }
        const correctedOffsetMinutes = this.timezoneOffsetMinutes(
          new Date(wallClockUtcMs - initialOffsetMinutes * 60000),
        );
        if (correctedOffsetMinutes === null) {
          return "";
        }
        const offsetSign = correctedOffsetMinutes < 0 ? "-" : "+";
        const absoluteOffsetMinutes = Math.abs(correctedOffsetMinutes);
        const offsetHourText = String(
          Math.floor(absoluteOffsetMinutes / 60),
        ).padStart(2, "0");
        const offsetMinuteText = String(absoluteOffsetMinutes % 60).padStart(
          2,
          "0",
        );
        return `${offsetSign}${offsetHourText}:${offsetMinuteText}`;
      },

      timezoneWallClockNowParts() {
        const now = new Date();
        if (this.timezone !== "") {
          try {
            const clockParts = new Intl.DateTimeFormat("en-US", {
              timeZone: this.timezone,
              hour: "2-digit",
              minute: "2-digit",
              hourCycle: "h23",
            }).formatToParts(now);
            const hour = clockParts.find((part) => part.type === "hour")?.value;
            const minute = clockParts.find(
              (part) => part.type === "minute",
            )?.value;
            if (hour && minute) {
              return { hour, minute };
            }
          } catch (timezoneError) {
            console.warn(
              "DateTimePickerTimezoneClockUnavailable",
              timezoneError,
            );
          }
        }
        return {
          hour: String(now.getHours()).padStart(2, "0"),
          minute: String(now.getMinutes()).padStart(2, "0"),
        };
      },

      dateTimezoneComposer(isoValue) {
        if (isoValue === "" || this.mode !== "datetime") {
          return isoValue;
        }
        const dateText = this.datePartExtractor(isoValue);
        const timeText =
          this.timePartExtractor(isoValue) ||
          `${this.timeHour}:${this.timeMinute}`;
        const wallClockText = `${dateText}T${timeText}`;
        return `${wallClockText}${this.timezoneOffsetSuffix(wallClockText)}`;
      },

      assignDraft(isoValue) {
        if (!this.isRange) {
          this.startDraft = isoValue;
          return;
        }
        this.assignSideDraft(this.activeEnd, isoValue);
      },

      assignSideDraft(side, isoValue) {
        if (side === "end") {
          this.endDraft = isoValue;
          return;
        }
        this.startDraft = isoValue;
      },

      sideDraft(side) {
        return side === "end" ? this.endDraft : this.startDraft;
      },

      composeDateDraft(isoDate) {
        if (this.mode === "datetime") {
          const wallClockText = `${isoDate}T${this.timeHour}:${this.timeMinute}`;
          return `${wallClockText}${this.timezoneOffsetSuffix(wallClockText)}`;
        }
        return isoDate;
      },

      normalizeRangeOrder() {
        if (!this.startDraft || !this.endDraft) {
          return;
        }
        const start = this.boundComparisonTextResolver(this.startDraft);
        const end = this.boundComparisonTextResolver(this.endDraft);
        if (start > end) {
          const previousStart = this.startDraft;
          this.startDraft = this.endDraft;
          this.endDraft = previousStart;
        }
      },

      syncTimeSelectsFromDraft() {
        if (this.mode === "date") {
          return;
        }
        const timeText = this.timePartExtractor(this.activeDraft);
        if (timeText === "") {
          const nowParts = this.timezoneWallClockNowParts();
          this.timeHour = nowParts.hour;
          this.timeMinute = nowParts.minute;
          return;
        }
        this.timeHour = timeText.slice(0, 2);
        this.timeMinute = timeText.slice(3, 5);
      },

      syncTypedValues() {
        this.startTypedValue = this.draftDisplayTextFormatter(this.startDraft);
        this.endTypedValue = this.draftDisplayTextFormatter(this.endDraft);
      },

      jumpViewToActiveDraft() {
        if (this.mode === "time") {
          return;
        }
        const date = this.dateObjectFromIso(this.activeDraft);
        if (date === null) {
          return;
        }
        this.viewYear = date.getFullYear();
        this.viewMonth = date.getMonth();
      },

      calendarDateIsValid(isoDate) {
        const dateMatch = String(isoDate ?? "").match(
          /^(\d{4})-(\d{2})-(\d{2})$/,
        );
        if (dateMatch === null) {
          return false;
        }
        const dateYear = Number(dateMatch[1]);
        const dateMonth = Number(dateMatch[2]) - 1;
        const dateDay = Number(dateMatch[3]);
        const date = new Date(dateYear, dateMonth, dateDay);
        return (
          date.getFullYear() === dateYear &&
          date.getMonth() === dateMonth &&
          date.getDate() === dateDay
        );
      },

      parseIsoInput(isoText) {
        const trimmed = String(isoText ?? "").trim();
        if (trimmed === "") {
          return "";
        }
        if (this.mode === "time") {
          if (!/^([01]\d|2[0-3]):[0-5]\d$/.test(trimmed)) {
            return null;
          }
          return this.valueIsOutOfBounds(trimmed) ? null : trimmed;
        }
        if (this.mode === "date") {
          if (
            !/^\d{4}-\d{2}-\d{2}$/.test(trimmed) ||
            !this.calendarDateIsValid(trimmed)
          ) {
            return null;
          }
          return this.valueIsOutOfBounds(trimmed) ? null : trimmed;
        }
        if (
          !/^\d{4}-\d{2}-\d{2}T([01]\d|2[0-3]):[0-5]\d(?:[+-]\d{2}:\d{2}|Z)?$/.test(
            trimmed,
          )
        ) {
          return null;
        }
        const wallClockText = trimmed.slice(0, 16);
        if (
          !this.calendarDateIsValid(wallClockText.slice(0, 10)) ||
          this.valueIsOutOfBounds(wallClockText)
        ) {
          return null;
        }
        return wallClockText + this.timezoneOffsetSuffix(wallClockText);
      },

      draftDisplayTextFormatter(iso) {
        const date = this.dateObjectFromIso(iso);
        if (date === null) {
          return "";
        }
        if (this.mode === "time") {
          return new Intl.DateTimeFormat(this.locale, {
            hour: "2-digit",
            minute: "2-digit",
          }).format(date);
        }
        const formatterOptions = {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
        };
        if (this.mode === "datetime") {
          formatterOptions.hour = "2-digit";
          formatterOptions.minute = "2-digit";
        }
        return new Intl.DateTimeFormat(this.locale, formatterOptions).format(
          date,
        );
      },

      datePartExtractor(iso) {
        if (!iso) {
          return "";
        }
        if (this.mode === "time") {
          return iso.slice(0, 5);
        }
        return iso.slice(0, 10);
      },

      timePartExtractor(iso) {
        if (!iso) {
          return "";
        }
        if (this.mode === "time") {
          return iso.slice(0, 5);
        }
        const separatorIndex = iso.indexOf("T");
        if (separatorIndex < 0) {
          return "";
        }
        return iso.slice(separatorIndex + 1, separatorIndex + 6);
      },

      dayIsOutOfBounds(isoDate) {
        if (this.mode === "time") {
          return false;
        }
        const minDate = this.minValue ? this.minValue.slice(0, 10) : "";
        const maxDate = this.maxValue ? this.maxValue.slice(0, 10) : "";
        if (minDate && isoDate < minDate) {
          return true;
        }
        if (maxDate && isoDate > maxDate) {
          return true;
        }
        return false;
      },

      valueIsOutOfBounds(iso) {
        if (!iso) {
          return false;
        }
        const comparableValue = this.boundComparisonTextResolver(iso);
        const comparableMin = this.fullDayBoundComparisonTextResolver(
          this.minValue,
          "00:00",
        );
        const comparableMax = this.fullDayBoundComparisonTextResolver(
          this.maxValue,
          "23:59",
        );
        if (comparableMin && comparableValue < comparableMin) {
          return true;
        }
        if (comparableMax && comparableValue > comparableMax) {
          return true;
        }
        return false;
      },

      boundComparisonTextResolver(iso) {
        if (!iso) {
          return "";
        }
        if (this.mode === "time") {
          return iso.slice(-5);
        }
        if (this.mode === "date") {
          return iso.slice(0, 10);
        }
        return iso.slice(0, 16);
      },

      fullDayBoundComparisonTextResolver(bound, dayEdgeTimeText) {
        const comparableBound = this.boundComparisonTextResolver(bound);
        const boundIsDateOnly = comparableBound.length === 10;
        if (this.mode !== "datetime" || !boundIsDateOnly) {
          return comparableBound;
        }
        return `${comparableBound}T${dayEdgeTimeText}`;
      },

      dateObjectFromIso(iso) {
        const isoText = String(iso ?? "").trim();
        if (isoText === "") {
          return null;
        }
        const dateTimeMatch = isoText.match(
          /^(\d{4})-(\d{2})-(\d{2})T([01]\d|2[0-3]):([0-5]\d)(?:[+-]\d{2}:\d{2}|Z)?$/,
        );
        if (dateTimeMatch) {
          return new Date(
            Number(dateTimeMatch[1]),
            Number(dateTimeMatch[2]) - 1,
            Number(dateTimeMatch[3]),
            Number(dateTimeMatch[4]),
            Number(dateTimeMatch[5]),
          );
        }
        const dateMatch = isoText.match(/^(\d{4})-(\d{2})-(\d{2})$/);
        if (dateMatch) {
          return new Date(
            Number(dateMatch[1]),
            Number(dateMatch[2]) - 1,
            Number(dateMatch[3]),
          );
        }
        const timeMatch = isoText.match(/^([01]\d|2[0-3]):([0-5]\d)$/);
        if (timeMatch) {
          return new Date(
            referenceSunday.getFullYear(),
            referenceSunday.getMonth(),
            referenceSunday.getDate(),
            Number(timeMatch[1]),
            Number(timeMatch[2]),
          );
        }
        return null;
      },

      isoFromDateObject(date) {
        const year = String(date.getFullYear()).padStart(4, "0");
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
      },
    };
  };

  Alpine.data("dateTimePicker", dateTimePickerStateFactory);
});

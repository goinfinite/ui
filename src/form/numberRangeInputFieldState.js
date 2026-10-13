UiToolset.RegisterAlpineState(() => {
  function decimalPlacesCounter(numberText) {
    const [mantissaText, exponentText = "0"] = String(numberText)
      .toLowerCase()
      .split("e");
    const decimalSeparatorIndex = mantissaText.indexOf(".");
    let mantissaDecimalPlaces = 0;
    if (decimalSeparatorIndex !== -1) {
      mantissaDecimalPlaces = mantissaText.length - decimalSeparatorIndex - 1;
    }
    return Math.max(mantissaDecimalPlaces - Number(exponentText), 0);
  }

  Alpine.data(
    "numberRangeInputField",
    (minValueText, maxValueText, stepText) => ({
      numberRangeTypedValue: "",

      syncNumberRangeTypedValue(boundValue) {
        this.numberRangeTypedValue = String(boundValue ?? "");
      },

      boundedNumberRangeValueResolver(typedValue, fallbackValue) {
        const minValue = Number(minValueText);
        const maxValue = Number(maxValueText);
        const stepSize = Number(stepText);

        let candidateValue = Number.parseFloat(typedValue);
        if (Number.isNaN(candidateValue)) {
          candidateValue = Number.parseFloat(fallbackValue);
        }
        if (Number.isNaN(candidateValue)) {
          candidateValue = minValue;
        }

        const clampedValue = Math.min(
          Math.max(candidateValue, minValue),
          maxValue,
        );
        if (!(stepSize > 0)) {
          return clampedValue;
        }
        const floatingPointTolerance = 1e-9;
        const lastStepCountBelowMax = Math.floor(
          (maxValue - minValue) / stepSize + floatingPointTolerance,
        );
        const nearestStepCount = Math.round(
          (clampedValue - minValue) / stepSize,
        );
        const stepCountFromMin = Math.min(
          nearestStepCount,
          lastStepCountBelowMax,
        );
        const steppedValue = minValue + stepCountFromMin * stepSize;
        const valuePrecision = Math.max(
          decimalPlacesCounter(minValueText),
          decimalPlacesCounter(stepText),
        );
        return Number(steppedValue.toFixed(valuePrecision));
      },

      commitNumberRangeTypedValue(fallbackValue) {
        const boundedValue = this.boundedNumberRangeValueResolver(
          this.numberRangeTypedValue,
          fallbackValue,
        );
        this.numberRangeTypedValue = String(boundedValue);
        return boundedValue;
      },
    }),
  );
});

UiToolset.RegisterAlpineState(() => {
  Alpine.data("selectInput", (labelValueOptionsScriptId) => ({
    ...UiToolset.SelectDropdown.controller,
    labelValueOptionsScriptId: labelValueOptionsScriptId || "",
    isOpen: false,
    openUpward: false,
    selectedValue: "",
    selectedItems: [],

    dynamicOptionLabelResolver(options, optionValue) {
      if (!Array.isArray(options) || optionValue === undefined) {
        return "";
      }

      const matchedOption = options.find(
        (option) => String(option?.value ?? option) === String(optionValue),
      );
      if (matchedOption === undefined) {
        return "";
      }
      if (typeof matchedOption === "object" && matchedOption !== null) {
        return matchedOption.label ?? matchedOption.value ?? "";
      }
      return matchedOption;
    },

    dynamicOptionsDisplayResolver(options, selectedValue) {
      if (!Array.isArray(selectedValue)) {
        return this.dynamicOptionLabelResolver(options, selectedValue);
      }
      if (selectedValue.length === 0) {
        return null;
      }

      const selectedLabels = selectedValue
        .map((optionValue) =>
          this.dynamicOptionLabelResolver(options, optionValue),
        )
        .filter((optionLabel) => optionLabel !== "");
      if (selectedLabels.length === 0) {
        return null;
      }
      return selectedLabels.join(", ");
    },

    selectedItemsDisplayFormatter(items) {
      if (!items || items.length === 0) {
        return null;
      }

      let labelValueOptions = [];
      if (this.labelValueOptionsScriptId) {
        try {
          labelValueOptions = JSON.parse(
            document.getElementById(this.labelValueOptionsScriptId)
              ?.textContent || "[]",
          );
        } catch (parseError) {
          console.error(
            `SelectInputInvalidLabelValueOptionsJson: ${parseError.message}`,
          );
        }
      }

      if (labelValueOptions.length === 0) {
        return items.join(", ");
      }

      const valueToLabelLookup = labelValueOptions.reduce((lookup, option) => {
        lookup[option.value] = option.label;
        return lookup;
      }, {});

      return items
        .map((value) => valueToLabelLookup[value] ?? value)
        .join(", ");
    },
  }));
});

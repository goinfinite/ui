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

    selectedOptionValueIsPresent(selectedItems, optionValue) {
      if (!Array.isArray(selectedItems)) {
        return false;
      }
      return selectedItems.some(
        (selectedValue) => String(selectedValue) === String(optionValue),
      );
    },

    toggleSelectedOptionValue(selectedItems, optionValue, isChecked) {
      const currentItems = Array.isArray(selectedItems) ? selectedItems : [];
      const itemsWithoutOption = currentItems.filter(
        (selectedValue) => String(selectedValue) !== String(optionValue),
      );
      if (!isChecked) {
        return itemsWithoutOption;
      }
      return [...itemsWithoutOption, optionValue];
    },

    dynamicOptionDisplayLabelResolver(options, optionValue) {
      const optionLabel = this.dynamicOptionLabelResolver(options, optionValue);
      if (optionLabel !== "") {
        return optionLabel;
      }
      if (optionValue === undefined || optionValue === null) {
        return "";
      }
      return String(optionValue);
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
          this.dynamicOptionDisplayLabelResolver(options, optionValue),
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

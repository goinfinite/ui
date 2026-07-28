UiToolset.RegisterAlpineState(() => {
  Alpine.data("multiSelectInput", (labelValueOptionsScriptId) => ({
    labelValueOptionsScriptId: labelValueOptionsScriptId || "",
    selectedItems: [],
    shouldExpandOptions: false,
    selectedItemsDisplayFormatter(items) {
      if (!items || items.length === 0) {
        return null;
      }

      let labelValueOptions = [];
      if (this.labelValueOptionsScriptId) {
        try {
          labelValueOptions = JSON.parse(
            document.getElementById(this.labelValueOptionsScriptId)?.textContent || "[]"
          );
        } catch (parseError) {
          console.error(
            `MultiSelectInputInvalidLabelValueOptionsJson: ${parseError.message}`
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
    closeDropdown() {
      this.shouldExpandOptions = false;
    },
    toggleDropdownDisplay() {
      this.shouldExpandOptions = !this.shouldExpandOptions;
    },
  }));
});
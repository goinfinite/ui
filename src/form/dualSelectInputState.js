UiToolset.RegisterAlpineState(() => {
  Alpine.data("dualSelectInput", (optionsScriptId) => ({
    ...UiToolset.SelectDropdown.controller,
    optionsScriptId: optionsScriptId || "",
    options: [],
    selectedValue: "",
    isOpen: false,
    openUpward: false,

    init() {
      if (!this.optionsScriptId) {
        return;
      }
      try {
        this.options = JSON.parse(
          document.getElementById(this.optionsScriptId)?.textContent || "[]",
        );
      } catch (parseError) {
        console.error(
          `DualSelectInputInvalidOptionsJson: ${parseError.message}`,
        );
      }
    },

    hasUnmatchedSelectValue(selectedValue) {
      if (
        selectedValue === "" ||
        selectedValue === null ||
        selectedValue === undefined
      ) {
        return false;
      }
      return !this.options.some(
        (option) => String(option.value) === String(selectedValue),
      );
    },

    selectedLabelResolver(selectedValue) {
      const matchedOption = this.options.find(
        (option) => String(option.value) === String(selectedValue),
      );
      if (matchedOption !== undefined) {
        return matchedOption.label;
      }
      if (
        selectedValue === "" ||
        selectedValue === null ||
        selectedValue === undefined
      ) {
        return "";
      }
      return "Custom";
    },
  }));
});

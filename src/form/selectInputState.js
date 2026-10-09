UiToolset.RegisterAlpineState(() => {
  Alpine.data("selectInput", (labelValueOptionsScriptId) => ({
    labelValueOptionsScriptId: labelValueOptionsScriptId || "",
    isOpen: false,
    openUpward: false,
    selectedValue: "",
    selectedItems: [],

    toggleDropdown() {
      this.isOpen = !this.isOpen;
      if (!this.isOpen) {
        return;
      }
      this.$nextTick(() => {
        this.openUpward = UiToolset.SelectDropdown.openUpwardResolver(
          this.$refs.selectTrigger,
          this.$refs.selectDropdown,
        );
      });
    },

    closeDropdown() {
      this.isOpen = false;
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

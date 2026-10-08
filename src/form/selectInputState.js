UiToolset.RegisterAlpineState(() => {
  Alpine.data("selectInput", () => ({
    isOpen: false,
    openUpward: false,
    selectedValue: "",

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
  }));
});

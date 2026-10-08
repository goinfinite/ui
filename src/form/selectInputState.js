UiToolset.RegisterAlpineState(() => {
  const selectOpenUpwardSlackPx = 8;

  function selectClipperBottomResolver(triggerElement) {
    let availableBottom = window.innerHeight;
    for (
      let ancestor = triggerElement.parentElement;
      ancestor !== null;
      ancestor = ancestor.parentElement
    ) {
      const ancestorStyle = getComputedStyle(ancestor);
      if (
        ancestorStyle.overflowX === "visible" &&
        ancestorStyle.overflowY === "visible"
      ) {
        continue;
      }
      availableBottom = Math.min(
        availableBottom,
        ancestor.getBoundingClientRect().bottom,
      );
    }
    return availableBottom;
  }

  function selectOpenUpwardResolver(trigger, dropdown) {
    if (!trigger || !dropdown) {
      return false;
    }
    const triggerBottom = trigger.getBoundingClientRect().bottom;
    const availableBottom = selectClipperBottomResolver(trigger);
    return (
      availableBottom - triggerBottom <
      dropdown.offsetHeight + selectOpenUpwardSlackPx
    );
  }

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
        this.openUpward = selectOpenUpwardResolver(
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

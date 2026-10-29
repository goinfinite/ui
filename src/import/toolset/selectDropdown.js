const selectDropdownOpenUpwardSlackPx = 8;

function selectDropdownClipperBoundsResolver(triggerElement) {
  let availableTop = 0;
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
    const ancestorRect = ancestor.getBoundingClientRect();
    availableTop = Math.max(availableTop, ancestorRect.top);
    availableBottom = Math.min(availableBottom, ancestorRect.bottom);
  }
  return { availableTop, availableBottom };
}

function selectDropdownOpenUpwardResolver(trigger, dropdown) {
  if (!trigger || !dropdown) {
    return false;
  }
  const triggerRect = trigger.getBoundingClientRect();
  const { availableTop, availableBottom } =
    selectDropdownClipperBoundsResolver(trigger);
  const requiredHeight =
    dropdown.offsetHeight + selectDropdownOpenUpwardSlackPx;
  const spaceBelow = availableBottom - triggerRect.bottom;
  if (spaceBelow >= requiredHeight) {
    return false;
  }
  const spaceAbove = triggerRect.top - availableTop;
  if (spaceAbove >= requiredHeight) {
    return true;
  }
  return spaceAbove > spaceBelow;
}

const selectDropdownController = {
  toggleDropdown() {
    if (this.isOpen) {
      this.closeDropdown();
      return;
    }
    this.openWithOptionFocus();
  },

  openWithOptionFocus() {
    const trigger = this.$refs.selectTrigger;
    const dropdown = this.$refs.selectDropdown;
    this.isOpen = true;
    this.$nextTick(() => {
      this.openUpward = selectDropdownOpenUpwardResolver(trigger, dropdown);
      requestAnimationFrame(() => {
        if (!this.isOpen) {
          return;
        }
        const checkedOption =
          dropdown.querySelector("input:checked") ??
          dropdown.querySelector("input");
        if (checkedOption !== null) {
          checkedOption.focus();
        }
      });
    });
  },

  closeDropdown() {
    this.isOpen = false;
  },

  focusAdjacentOption(direction) {
    const dropdown = this.$refs.selectDropdown;
    const options = [
      ...dropdown.querySelectorAll(
        'input[type="radio"], input[type="checkbox"]',
      ),
    ];
    if (options.length === 0) {
      return;
    }
    const focusedIndex = options.indexOf(document.activeElement);
    const firstIndex = direction > 0 ? 0 : options.length - 1;
    const steppedIndex = Math.min(
      Math.max(focusedIndex + direction, 0),
      options.length - 1,
    );
    const nextOption = options[focusedIndex === -1 ? firstIndex : steppedIndex];
    nextOption.focus();
    if (nextOption.type === "radio") {
      nextOption.click();
    }
  },

  closeWithFocusReturn() {
    const wasOpen = this.isOpen;
    this.isOpen = false;
    if (!wasOpen) {
      return;
    }
    this.$nextTick(() => this.$refs.selectTrigger.focus());
  },

  commitFocusedOption() {
    const focusedOption = document.activeElement;
    const isOptionFocused =
      focusedOption !== null &&
      this.$refs.selectDropdown.contains(focusedOption) &&
      focusedOption.matches('input[type="radio"], input[type="checkbox"]');
    if (!isOptionFocused) {
      this.closeWithFocusReturn();
      return;
    }
    focusedOption.click();
    this.closeWithFocusReturn();
  },
};

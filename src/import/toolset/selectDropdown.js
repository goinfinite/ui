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

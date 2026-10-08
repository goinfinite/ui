const selectDropdownOpenUpwardSlackPx = 8;

function selectDropdownClipperBottomResolver(triggerElement) {
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

function selectDropdownOpenUpwardResolver(trigger, dropdown) {
  if (!trigger || !dropdown) {
    return false;
  }
  const triggerBottom = trigger.getBoundingClientRect().bottom;
  const availableBottom = selectDropdownClipperBottomResolver(trigger);
  return (
    availableBottom - triggerBottom <
    dropdown.offsetHeight + selectDropdownOpenUpwardSlackPx
  );
}

const buttonTooltipOffsetPx = 6;
const buttonTooltipViewportPaddingPx = 4;

function clampNumber(value, minimum, maximum) {
  return Math.min(Math.max(value, minimum), maximum);
}

function resolveFlippedAxisStart(
  triggerStart,
  triggerEnd,
  tooltipSize,
  viewportSize,
  prefersAfter,
) {
  const beforeStart = triggerStart - tooltipSize - buttonTooltipOffsetPx;
  const afterStart = triggerEnd + buttonTooltipOffsetPx;
  let start = prefersAfter ? afterStart : beforeStart;
  if (start < buttonTooltipViewportPaddingPx) {
    start = afterStart;
  }
  if (start + tooltipSize > viewportSize - buttonTooltipViewportPaddingPx) {
    start = beforeStart;
  }
  return start;
}

function resolveCenteredAxisStart(
  triggerStart,
  triggerSize,
  tooltipSize,
  viewportSize,
) {
  const centeredStart = triggerStart + triggerSize / 2 - tooltipSize / 2;
  return clampNumber(
    centeredStart,
    buttonTooltipViewportPaddingPx,
    viewportSize - tooltipSize - buttonTooltipViewportPaddingPx,
  );
}

function resolveTooltipCoordinates(triggerRect, tooltipRect, position) {
  const isVerticalAxisPrimary = position === "top" || position === "bottom";
  if (isVerticalAxisPrimary) {
    return {
      top: resolveFlippedAxisStart(
        triggerRect.top,
        triggerRect.bottom,
        tooltipRect.height,
        window.innerHeight,
        position === "bottom",
      ),
      left: resolveCenteredAxisStart(
        triggerRect.left,
        triggerRect.width,
        tooltipRect.width,
        window.innerWidth,
      ),
    };
  }
  return {
    top: resolveCenteredAxisStart(
      triggerRect.top,
      triggerRect.height,
      tooltipRect.height,
      window.innerHeight,
    ),
    left: resolveFlippedAxisStart(
      triggerRect.left,
      triggerRect.right,
      tooltipRect.width,
      window.innerWidth,
      position === "right",
    ),
  };
}

UiToolset.RegisterAlpineState(() => {
  Alpine.data("buttonTooltip", (position) => ({
    position: position || "top",
    isVisible: false,
    tooltipCoordinates: { top: 0, left: 0 },
    viewportChangeHandler: null,

    get resolvedTooltipStyle() {
      return {
        visibility: this.isVisible ? "visible" : "hidden",
        opacity: this.isVisible ? "1" : "0",
        top: `${this.tooltipCoordinates.top}px`,
        left: `${this.tooltipCoordinates.left}px`,
      };
    },

    init() {
      this.viewportChangeHandler = () => {
        if (this.isVisible) {
          this.show();
        }
      };
      window.addEventListener("scroll", this.viewportChangeHandler, true);
      window.addEventListener("resize", this.viewportChangeHandler);
    },

    destroy() {
      window.removeEventListener("scroll", this.viewportChangeHandler, true);
      window.removeEventListener("resize", this.viewportChangeHandler);
    },

    show() {
      const trigger = this.$refs.trigger;
      const tooltip = this.$refs.tooltip;
      if (!trigger || !tooltip) {
        return;
      }
      this.tooltipCoordinates = resolveTooltipCoordinates(
        trigger.getBoundingClientRect(),
        tooltip.getBoundingClientRect(),
        this.position,
      );
      this.isVisible = true;
    },

    hide() {
      this.isVisible = false;
    },
  }));
});

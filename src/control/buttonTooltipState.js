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
  const maximumStart = Math.max(
    buttonTooltipViewportPaddingPx,
    viewportSize - tooltipSize - buttonTooltipViewportPaddingPx,
  );
  return clampNumber(start, buttonTooltipViewportPaddingPx, maximumStart);
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
    tooltipPosition: position || "top",
    isTooltipVisible: false,
    tooltipCoordinates: { top: 0, left: 0 },
    tooltipViewportChangeHandler: null,

    get resolvedTooltipStyle() {
      return {
        visibility: this.isTooltipVisible ? "visible" : "hidden",
        opacity: this.isTooltipVisible ? "1" : "0",
        top: `${this.tooltipCoordinates.top}px`,
        left: `${this.tooltipCoordinates.left}px`,
      };
    },

    init() {
      this.tooltipViewportChangeHandler = () => {
        if (this.isTooltipVisible) {
          this.showTooltip();
        }
      };
      window.addEventListener(
        "scroll",
        this.tooltipViewportChangeHandler,
        true,
      );
      window.addEventListener("resize", this.tooltipViewportChangeHandler);
    },

    destroy() {
      window.removeEventListener(
        "scroll",
        this.tooltipViewportChangeHandler,
        true,
      );
      window.removeEventListener("resize", this.tooltipViewportChangeHandler);
    },

    showTooltip() {
      const trigger = this.$refs.trigger;
      const tooltip = this.$refs.tooltip;
      if (!trigger || !tooltip) {
        return;
      }
      this.tooltipCoordinates = resolveTooltipCoordinates(
        trigger.getBoundingClientRect(),
        tooltip.getBoundingClientRect(),
        this.tooltipPosition,
      );
      this.isTooltipVisible = true;
    },

    hideTooltip() {
      this.isTooltipVisible = false;
    },
  }));
});

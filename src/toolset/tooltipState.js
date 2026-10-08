UiToolset.RegisterAlpineState(() => {
  const tooltipOffsetPx = 6;
  const tooltipViewportPaddingPx = 4;
  const tooltipShowDelayMs = 200;

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
    const beforeStart = triggerStart - tooltipSize - tooltipOffsetPx;
    const afterStart = triggerEnd + tooltipOffsetPx;
    let start = prefersAfter ? afterStart : beforeStart;
    if (start < tooltipViewportPaddingPx) {
      start = afterStart;
    }
    if (start + tooltipSize > viewportSize - tooltipViewportPaddingPx) {
      start = beforeStart;
    }
    const maximumStart = Math.max(
      tooltipViewportPaddingPx,
      viewportSize - tooltipSize - tooltipViewportPaddingPx,
    );
    return clampNumber(start, tooltipViewportPaddingPx, maximumStart);
  }

  function resolveCenteredAxisStart(
    triggerStart,
    triggerSize,
    tooltipSize,
    viewportSize,
  ) {
    const centeredStart = triggerStart + triggerSize / 2 - tooltipSize / 2;
    const maximumStart = Math.max(
      tooltipViewportPaddingPx,
      viewportSize - tooltipSize - tooltipViewportPaddingPx,
    );
    return clampNumber(centeredStart, tooltipViewportPaddingPx, maximumStart);
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

  Alpine.data("tooltip", (position) => ({
    tooltipPosition: position || "top",
    isTooltipVisible: false,
    tooltipCoordinates: { top: 0, left: 0 },
    tooltipViewportChangeHandler: null,
    tooltipShowTimeoutId: null,

    get resolvedTooltipStyle() {
      return {
        visibility: this.isTooltipVisible ? "visible" : "hidden",
        opacity: this.isTooltipVisible ? "1" : "0",
        top: `${this.tooltipCoordinates.top}px`,
        left: `${this.tooltipCoordinates.left}px`,
      };
    },

    updateTooltipCoordinates() {
      const trigger = this.$refs.trigger;
      const tooltip = this.$refs.tooltip;
      if (!trigger || !tooltip) {
        return false;
      }
      this.tooltipCoordinates = resolveTooltipCoordinates(
        trigger.getBoundingClientRect(),
        tooltip.getBoundingClientRect(),
        this.tooltipPosition,
      );
      return true;
    },

    attachViewportListeners() {
      if (this.tooltipViewportChangeHandler) {
        return;
      }
      this.tooltipViewportChangeHandler = () => {
        if (this.isTooltipVisible) {
          this.updateTooltipCoordinates();
        }
      };
      window.addEventListener(
        "scroll",
        this.tooltipViewportChangeHandler,
        true,
      );
      window.addEventListener("resize", this.tooltipViewportChangeHandler);
    },

    detachViewportListeners() {
      if (!this.tooltipViewportChangeHandler) {
        return;
      }
      window.removeEventListener(
        "scroll",
        this.tooltipViewportChangeHandler,
        true,
      );
      window.removeEventListener("resize", this.tooltipViewportChangeHandler);
      this.tooltipViewportChangeHandler = null;
    },

    showTooltip() {
      clearTimeout(this.tooltipShowTimeoutId);
      this.tooltipShowTimeoutId = setTimeout(
        () => this.revealTooltip(),
        tooltipShowDelayMs,
      );
    },

    revealTooltip() {
      if (!this.updateTooltipCoordinates()) {
        return;
      }
      this.isTooltipVisible = true;
      this.attachViewportListeners();
    },

    hideTooltip() {
      clearTimeout(this.tooltipShowTimeoutId);
      this.isTooltipVisible = false;
      this.detachViewportListeners();
    },

    destroy() {
      clearTimeout(this.tooltipShowTimeoutId);
      this.detachViewportListeners();
    },
  }));
});

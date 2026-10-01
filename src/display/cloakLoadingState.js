UiToolset.RegisterAlpineState(() => {
  let isAlpineInitialized = false;
  document.addEventListener(
    "alpine:initialized",
    () => {
      isAlpineInitialized = true;
    },
    { once: true },
  );

  Alpine.data("cloakLoading", (hideDelayMilliseconds) => ({
    isHideScheduled: false,

    init() {
      if (isAlpineInitialized || document.readyState === "complete") {
        this.scheduleHide();
        return;
      }
      document.addEventListener(
        "alpine:initialized",
        () => this.scheduleHide(),
        { once: true },
      );
      window.addEventListener("load", () => this.scheduleHide(), {
        once: true,
      });
    },

    scheduleHide() {
      if (this.isHideScheduled) {
        return;
      }
      this.isHideScheduled = true;
      setTimeout(() => {
        this.$el.style.display = "none";
      }, hideDelayMilliseconds);
    },
  }));
});

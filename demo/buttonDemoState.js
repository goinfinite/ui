UiToolset.RegisterAlpineState(() => {
  const demoActionDurationMs = 1500;

  function demoActionPromiseBuilder(shouldSucceed) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        if (shouldSucceed) {
          resolve();
          return;
        }
        reject(new Error("DemoMirrorUnreachable"));
      }, demoActionDurationMs);
    });
  }

  Alpine.data("buttonDemoActions", () => ({
    isSyncing: false,
    dynamicButtonLabel: "Dynamic idle",
    dynamicPendingLabel: "Dynamic pending",
    dynamicSucceededLabel: "Dynamic succeeded",
    dynamicFailedLabel: "Dynamic failed",

    syncAccounts() {
      return demoActionPromiseBuilder(true);
    },

    syncBrokenMirror() {
      return demoActionPromiseBuilder(false);
    },

    startSharedSync() {
      this.isSyncing = true;
      setTimeout(() => {
        this.isSyncing = false;
      }, demoActionDurationMs);
    },
  }));
});

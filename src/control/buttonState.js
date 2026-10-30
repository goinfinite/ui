UiToolset.RegisterAlpineState(() => {
  const actionResultDisplayMs = 2000;

  function afterSubmitDispatchScheduler(pendingCallback) {
    const remainingListenersDelayMs = 0;
    setTimeout(pendingCallback, remainingListenersDelayMs);
  }

  Alpine.data("buttonState", (actionIcons, isFormAction, countdownSeconds) => ({
    actionStatus: "idle",
    actionResultTimeoutId: null,
    actionListenersController: new AbortController(),
    actionButton: null,
    submittedForm: null,
    formObserver: null,
    countdownSecondsRemaining: Math.max(
      0,
      Math.floor(Number(countdownSeconds) || 0),
    ),
    countdownIntervalId: null,
    hasReceivedCountdownTrigger: false,
    lastCountdownTriggerValue: undefined,

    init() {
      if (!isFormAction) {
        return;
      }
      this.actionButton = this.$el.querySelector("button");
      const formId = this.actionButton?.getAttribute("form") ?? "";
      const submittedForm = this.resolveSubmittedForm();
      if (submittedForm instanceof HTMLFormElement) {
        this.submittedForm = submittedForm;
        this.attachFormListeners();
        return;
      }
      if (formId === "") {
        console.error("ButtonActionFormNotFound");
        return;
      }
      if (document.getElementById(formId) !== null) {
        console.error("ButtonActionFormNotFound");
        return;
      }
      this.observeSubmittedForm(formId);
    },

    destroy() {
      this.actionListenersController.abort();
      this.formObserver?.disconnect();
      clearTimeout(this.actionResultTimeoutId);
      this.stopCountdown();
    },

    resolveSubmittedForm() {
      const formId = this.actionButton?.getAttribute("form") ?? "";
      if (formId === "") {
        return this.actionButton?.closest("form") ?? null;
      }
      const namedForm = document.getElementById(formId);
      if (!(namedForm instanceof HTMLFormElement)) {
        return null;
      }
      return namedForm;
    },

    observeSubmittedForm(formId) {
      this.formObserver = new MutationObserver(() => {
        const submittedForm = document.getElementById(formId);
        if (!(submittedForm instanceof HTMLFormElement)) {
          return;
        }
        this.formObserver.disconnect();
        this.formObserver = null;
        this.submittedForm = submittedForm;
        this.attachFormListeners();
      });
      this.formObserver.observe(document.body, {
        childList: true,
        subtree: true,
      });
    },

    attachFormListeners() {
      const listenerOptions = {
        signal: this.actionListenersController.signal,
      };
      this.submittedForm.addEventListener(
        "submit",
        (submitEvent) => this.blockRepeatedSubmission(submitEvent),
        { ...listenerOptions, capture: true },
      );
      this.submittedForm.addEventListener(
        "submit",
        (submitEvent) => this.markNativeSubmissionPending(submitEvent),
        listenerOptions,
      );
      this.submittedForm.addEventListener(
        "htmx:beforeSend",
        (requestEvent) => this.markRequestPending(requestEvent),
        listenerOptions,
      );
      this.submittedForm.addEventListener(
        "htmx:afterRequest",
        (requestEvent) => this.showRequestResult(requestEvent),
        listenerOptions,
      );
      window.addEventListener(
        "pageshow",
        (pageShowEvent) => this.resetRestoredPage(pageShowEvent),
        listenerOptions,
      );
    },

    markActionPending() {
      clearTimeout(this.actionResultTimeoutId);
      this.actionStatus = "pending";
    },

    showActionResult(resultStatus) {
      this.actionStatus = resultStatus;
      this.actionResultTimeoutId = setTimeout(() => {
        this.actionStatus = "idle";
      }, actionResultDisplayMs);
    },

    trackActionPromise(actionResult) {
      const isPromise = typeof actionResult?.then === "function";
      if (!isPromise) {
        return;
      }
      this.markActionPending();
      actionResult.then(
        () => this.showActionResult("succeeded"),
        (actionError) => {
          console.error(`ButtonActionFailed: ${actionError?.message}`);
          this.showActionResult("failed");
        },
      );
    },

    blockRepeatedSubmission(submitEvent) {
      if (this.actionStatus !== "pending") {
        return;
      }
      submitEvent.preventDefault();
      submitEvent.stopImmediatePropagation();
    },

    markNativeSubmissionPending(submitEvent) {
      const formTarget = this.submittedForm.target;
      const isPageNavigating = formTarget === "" || formTarget === "_self";
      if (!isPageNavigating) {
        return;
      }
      afterSubmitDispatchScheduler(() => {
        if (submitEvent.defaultPrevented) {
          return;
        }
        this.markActionPending();
      });
    },

    markRequestPending(requestEvent) {
      const isOwnFormRequest =
        requestEvent.detail.elt === this.submittedForm ||
        requestEvent.detail.elt === this.actionButton;
      if (!isOwnFormRequest) {
        return;
      }
      this.markActionPending();
    },

    showRequestResult(requestEvent) {
      const isOwnFormRequest =
        requestEvent.detail.elt === this.submittedForm ||
        requestEvent.detail.elt === this.actionButton;
      if (!isOwnFormRequest) {
        return;
      }
      let resultStatus = "failed";
      if (requestEvent.detail.successful) {
        resultStatus = "succeeded";
      }
      this.showActionResult(resultStatus);
    },

    resetRestoredPage(pageShowEvent) {
      if (!pageShowEvent.persisted) {
        return;
      }
      this.actionStatus = "idle";
    },

    actionDisplayStatusResolver(isPendingExternally) {
      if (isPendingExternally) {
        return "pending";
      }
      return this.actionStatus;
    },

    isActionPending(isPendingExternally) {
      const displayStatus =
        this.actionDisplayStatusResolver(isPendingExternally);
      return displayStatus === "pending";
    },

    actionLabelResolver(
      idleLabel,
      pendingLabel,
      succeededLabel,
      failedLabel,
      isPendingExternally,
    ) {
      const displayStatus =
        this.actionDisplayStatusResolver(isPendingExternally);
      if (displayStatus === "pending") {
        return pendingLabel;
      }
      if (displayStatus === "succeeded") {
        return succeededLabel;
      }
      if (displayStatus === "failed") {
        return failedLabel;
      }
      return idleLabel;
    },

    actionIconResolver(isPendingExternally) {
      const displayStatus =
        this.actionDisplayStatusResolver(isPendingExternally);
      return actionIcons[displayStatus];
    },

    startCountdown() {
      this.stopCountdown();
      const initialSeconds = Math.max(
        0,
        Math.floor(Number(countdownSeconds) || 0),
      );
      this.countdownSecondsRemaining = initialSeconds;
      if (initialSeconds === 0) {
        return;
      }
      this.countdownIntervalId = setInterval(() => {
        this.countdownSecondsRemaining = Math.max(
          0,
          this.countdownSecondsRemaining - 1,
        );
        if (this.countdownSecondsRemaining === 0) {
          this.stopCountdown();
        }
      }, 1000);
    },

    stopCountdown() {
      if (this.countdownIntervalId === null) {
        return;
      }
      clearInterval(this.countdownIntervalId);
      this.countdownIntervalId = null;
    },

    restartCountdownOnTrigger(triggerValue) {
      const isSameTriggerValue =
        this.hasReceivedCountdownTrigger &&
        triggerValue === this.lastCountdownTriggerValue;
      if (isSameTriggerValue) {
        return;
      }
      this.hasReceivedCountdownTrigger = true;
      this.lastCountdownTriggerValue = triggerValue;
      this.startCountdown();
    },

    get isCounting() {
      return this.countdownSecondsRemaining > 0;
    },
  }));
});

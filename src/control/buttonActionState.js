UiToolset.RegisterAlpineState(() => {
  const actionResultDisplayMs = 2000;

  const actionIconByStatus = {
    idle: "",
    pending: "ph-spinner animate-spin",
    succeeded: "ph-check-circle",
    failed: "ph-warning-circle",
  };

  function afterSubmitDispatchScheduler(pendingCallback) {
    const remainingListenersDelayMs = 0;
    setTimeout(pendingCallback, remainingListenersDelayMs);
  }

  Alpine.data("buttonAction", (actionLabelByStatus, isFormAction) => ({
    actionStatus: "idle",
    actionResultTimeoutId: null,
    actionListenersController: new AbortController(),
    submittedForm: null,

    init() {
      if (!isFormAction) {
        return;
      }
      const submittedForm = this.$el.closest("form");
      if (!submittedForm) {
        console.error("ButtonActionFormNotFound");
        return;
      }
      this.submittedForm = submittedForm;

      const listenerOptions = {
        signal: this.actionListenersController.signal,
      };
      submittedForm.addEventListener(
        "submit",
        (submitEvent) => this.blockRepeatedSubmission(submitEvent),
        { ...listenerOptions, capture: true },
      );
      submittedForm.addEventListener(
        "submit",
        (submitEvent) => this.markNativeSubmissionPending(submitEvent),
        listenerOptions,
      );
      submittedForm.addEventListener(
        "htmx:beforeSend",
        (requestEvent) => this.markRequestPending(requestEvent),
        listenerOptions,
      );
      submittedForm.addEventListener(
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

    destroy() {
      this.actionListenersController.abort();
      clearTimeout(this.actionResultTimeoutId);
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
      const isOwnFormRequest = requestEvent.detail.elt === this.submittedForm;
      if (!isOwnFormRequest) {
        return;
      }
      this.markActionPending();
    },

    showRequestResult(requestEvent) {
      const isOwnFormRequest = requestEvent.detail.elt === this.submittedForm;
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

    actionLabelResolver(isPendingExternally) {
      const displayStatus =
        this.actionDisplayStatusResolver(isPendingExternally);
      return actionLabelByStatus[displayStatus];
    },

    actionIconResolver(isPendingExternally) {
      const displayStatus =
        this.actionDisplayStatusResolver(isPendingExternally);
      return actionIconByStatus[displayStatus];
    },
  }));
});

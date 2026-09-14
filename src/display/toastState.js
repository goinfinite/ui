UiToolset.RegisterAlpineState(() => {
  Alpine.store("toast", {
    toastVisible: false,
    toastMessage: "",
    toastType: "danger",
    displayDurationMs: 10000,
    dismissTimeoutId: null,

    displayToast(message, toastType) {
      clearTimeout(this.dismissTimeoutId);
      this.toastVisible = true;
      this.toastMessage = message;
      this.toastType = toastType;
      this.dismissTimeoutId = setTimeout(() => {
        this.clearToast();
      }, this.displayDurationMs);
    },

    displayToastWithApiResponse(apiResponse, httpStatusCode) {
      const responseDisplay = UiToolset.ResolveApiResponseDisplay(
        apiResponse,
        httpStatusCode,
      );
      if (responseDisplay.message === "") {
        return;
      }

      let toastType = responseDisplay.outcome;
      if (toastType === "error") {
        toastType = "danger";
      }

      this.displayToast(responseDisplay.message, toastType);
    },

    clearToast() {
      clearTimeout(this.dismissTimeoutId);
      this.dismissTimeoutId = null;
      this.toastVisible = false;
      this.toastMessage = "";
    },
  });
});

document.addEventListener("htmx:afterRequest", (event) => {
  const httpResponseObject = event.detail.xhr;
  if (
    !httpResponseObject?.getResponseHeader("Content-Type")
      ?.includes("application/json")
  ) {
    return;
  }

  const responseData = httpResponseObject.responseText;
  if (responseData === "") {
    return;
  }

  let parsedResponse;
  try {
    parsedResponse = JSON.parse(responseData);
  } catch (parseError) {
    console.error("ToastApiResponseParseFailed", parseError);
    Alpine.store("toast").displayToast("UnexpectedResponse", "danger");
    return;
  }

  Alpine.store("toast").displayToastWithApiResponse(
    parsedResponse,
    httpResponseObject.status,
  );
});

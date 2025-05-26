document.addEventListener("alpine:init", () => {
  Alpine.store("toast", {
    toastVisible: false,
    toastMessage: "",
    toastType: "danger",

    displayToast(message, toastType) {
      this.toastVisible = true;
      this.toastMessage = message;
      this.toastType = toastType;
      setTimeout(() => {
        this.clearToast();
      }, 4000);
    },

    clearToast() {
      this.toastVisible = false;
      this.toastMessage = "";
    },
  });
});

document.addEventListener("htmx:afterRequest", (event) => {
  const httpErrorStatusCodeWithMessage = {
    400: "BadRequest",
    401: "Unauthorized",
    403: "Forbidden",
    404: "NotFound",
    500: "InternalServerError",
    502: "BadGateway",
    503: "ServiceUnavailable",
    504: "GatewayTimeout",
  };
  const httpResponseObject = event.detail.xhr;
  if (
    httpResponseObject.getResponseHeader("Content-Type") !== "application/json"
  ) {
    return;
  }

  const responseData = httpResponseObject.responseText;
  if (responseData === "") {
    return;
  }

  const parsedResponse = JSON.parse(responseData);
  if (parsedResponse.body === undefined || parsedResponse.body === "") {
    return;
  }

  httpResponseStatusCode = httpResponseObject.status;

  let toastType = "success";
  let toastMessage = "Success";
  if (httpResponseStatusCode == 207) {
    toastType = "partialSuccess";
    toastMessage = "PartialSuccess";
  }

  if (httpResponseStatusCode >= 400) {
    toastType = "danger";
    toastMessage = "Error";
    if (httpErrorStatusCodeWithMessage[httpResponseStatusCode]) {
      toastMessage = httpErrorStatusCodeWithMessage[httpResponseStatusCode];
    }
  }

  if (typeof parsedResponse.body === "string") {
    toastMessage = parsedResponse.body;
  }

  Alpine.store("toast").displayToast(toastMessage, toastType);
});

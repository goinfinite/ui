function displayApiErrorToast(error) {
  const toastStore = Alpine.store("toast");
  const hasApiResponse = error?.apiResponse !== undefined;
  if (!hasApiResponse) {
    toastStore.displayToast(error.message, "danger");
    return;
  }

  toastStore.displayToastWithApiResponse(
    error.apiResponse,
    error.httpStatusCode,
  );
}

let activeJsonAjaxRequestCount = 0;

async function jsonAjax(
  method = "POST",
  url,
  payload = {},
  shouldDisplayToast = true,
) {
  if (
    typeof method !== "string" ||
    !["GET", "POST", "PUT", "DELETE"].includes(method)
  ) {
    throw new Error(`InvalidHttpMethod`);
  }

  if (typeof url !== "string" || url.trim() === "") {
    throw new Error("InvalidUrl");
  }

  if (typeof payload !== "object" || payload === null) {
    throw new Error("InvalidPayload");
  }

  const payloadIsEmpty = Object.keys(payload).length === 0;

  if (method === "GET" && !payloadIsEmpty) {
    throw new Error("GetRequestPayloadNotAllowed");
  }

  if (typeof shouldDisplayToast !== "boolean") {
    throw new Error("InvalidShouldDisplayToast");
  }

  activeJsonAjaxRequestCount++;
  toggleLoadingOverlay(true);

  try {
    const response = await fetch(url, {
      method,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: payloadIsEmpty ? undefined : JSON.stringify(payload),
    });

    const contentType = response.headers.get("Content-Type");
    if (!contentType?.includes("application/json")) {
      throw new Error("UnexpectedResponseContentType");
    }

    const responseData = await response.json();

    if (!response.ok) {
      const responseDisplay = resolveApiResponseDisplay(
        responseData,
        response.status,
      );
      const responseErrorMessage = responseDisplay.message || "UnknownError";
      const responseError = new Error(responseErrorMessage);
      responseError.apiResponse = responseData;
      responseError.httpStatusCode = response.status;
      throw responseError;
    }

    if (shouldDisplayToast) {
      Alpine.store("toast").displayToastWithApiResponse(
        responseData,
        response.status,
      );
    }

    return responseData?.body;
  } catch (error) {
    if (shouldDisplayToast) {
      displayApiErrorToast(error);
    }

    throw error;
  } finally {
    activeJsonAjaxRequestCount--;
    if (activeJsonAjaxRequestCount === 0) {
      toggleLoadingOverlay(false);
    }
  }
}

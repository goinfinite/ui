function toggleLoadingOverlay(status = true) {
  const loadingOverlayElement = document.getElementById("loading-overlay");
  const loadingOverlayElementExists = loadingOverlayElement !== null;
  if (!loadingOverlayElementExists) {
    return;
  }

  const loadingOverlayElementIsLoading =
    loadingOverlayElement.classList.contains("htmx-request");

  if (loadingOverlayElementIsLoading && status) {
    return;
  }

  if (!loadingOverlayElementIsLoading && !status) {
    return;
  }

  if (!status) {
    loadingOverlayElement.classList.remove("htmx-request");
    return;
  }

  loadingOverlayElement.classList.add("htmx-request");
}

async function jsonAjax(
  method = "POST",
  url,
  payload = {},
  shouldDisplayToast = true
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

  if (typeof shouldDisplayToast !== "boolean") {
    shouldDisplayToast = true;
  }

  toggleLoadingOverlay(true);

  return fetch(url, {
    method,
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: Object.keys(payload).length > 0 ? JSON.stringify(payload) : undefined,
  })
    .then((response) => {
      toggleLoadingOverlay(false);

      const contentType = response.headers.get("Content-Type");
      if (!(contentType && contentType.includes("application/json"))) {
        throw new Error("UnexpectedResponseContentType");
      }

      return response.json().then((responseData) => {
        if (!response.ok) {
          throw new Error(
            responseData?.body ? responseData.body : "UnknownError"
          );
        }

        if (
          shouldDisplayToast &&
          responseData?.body &&
          typeof responseData.body === "string"
        ) {
          Alpine.store("toast").displayToast(responseData.body, "success");
        }

        return responseData?.body;
      });
    })
    .catch((error) => {
      toggleLoadingOverlay(false);

      if (shouldDisplayToast) {
        Alpine.store("toast").displayToast(error.message, "danger");
      }

      throw error;
    });
}

function createRandomPassword() {
  const passwordLength = 16;
  const chars =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+";

  let passwordContent = "";
  let passwordIterationCount = 0;
  while (passwordIterationCount < passwordLength) {
    const randomIndex = Math.floor(Math.random() * chars.length);
    const indexAfterRandomIndex = randomIndex + 1;
    passwordContent += chars.substring(randomIndex, indexAfterRandomIndex);

    passwordIterationCount++;
  }

  return passwordContent;
}

function registerAlpineState(stateFunction) {
  if (window.Alpine) {
    stateFunction();
    return;
  }

  document.addEventListener("alpine:init", stateFunction);
}

window.UiToolset = {
  CreateRandomPassword: createRandomPassword,
  ToggleLoadingOverlay: toggleLoadingOverlay,
  RegisterAlpineState: registerAlpineState,
};

registerAlpineState(() => {
  window.UiToolset.JsonAjax = jsonAjax;
});

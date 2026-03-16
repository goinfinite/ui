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
  const letterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
  const numberChars = "0123456789";
  const specialChars = "!@#$%^&*()_+";
  const allChars = letterChars + numberChars + specialChars;

  const pickRandomChar = (charSet) => {
    const randomIndex = Math.floor(Math.random() * charSet.length);
    return charSet[randomIndex];
  };

  let passwordChars = [];
  let passwordIterationCount = 0;
  while (passwordIterationCount < passwordLength) {
    passwordChars.push(pickRandomChar(allChars));
    passwordIterationCount++;
  }

  const hasLetter = passwordChars.some((char) => /[a-zA-Z]/.test(char));
  const hasNumber = passwordChars.some((char) => /[0-9]/.test(char));
  const hasSpecial = passwordChars.some((char) => /[^a-zA-Z0-9]/.test(char));

  const missingClasses = [];
  if (!hasLetter) missingClasses.push(letterChars);
  if (!hasNumber) missingClasses.push(numberChars);
  if (!hasSpecial) missingClasses.push(specialChars);

  let replacementPositionIndex = 0;
  for (const missingCharSet of missingClasses) {
    passwordChars[replacementPositionIndex] = pickRandomChar(missingCharSet);
    replacementPositionIndex++;
  }

  return passwordChars.join("");
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

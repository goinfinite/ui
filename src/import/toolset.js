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

function randomNumberGenerator(rangeSize) {
  const randomValues = crypto.getRandomValues(new Uint32Array(1));
  const rawRandomInteger = randomValues[0];
  const numberWithinRange = rawRandomInteger % rangeSize;
  return numberWithinRange;
}

function createRandomPassword() {
  const passwordLength = 16;
  const letterChars =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
  const numberChars = "0123456789";
  const specialChars = "!@#$%^&*()_+";
  const allChars = letterChars + numberChars + specialChars;

  let passwordChars = [];
  for (
    let charIndex = 0;
    charIndex < passwordLength;
    charIndex++
  ) {
    const randomPosition = randomNumberGenerator(allChars.length);
    passwordChars.push(allChars[randomPosition]);
  }

  const letterPosition = randomNumberGenerator(passwordLength);

  let numberPosition = randomNumberGenerator(passwordLength);
  while (numberPosition === letterPosition) {
    numberPosition = randomNumberGenerator(passwordLength);
  }

  let specialPosition = randomNumberGenerator(passwordLength);
  while (
    specialPosition === letterPosition ||
    specialPosition === numberPosition
  ) {
    specialPosition = randomNumberGenerator(passwordLength);
  }

  passwordChars[letterPosition] =
    letterChars[randomNumberGenerator(letterChars.length)];
  passwordChars[numberPosition] =
    numberChars[randomNumberGenerator(numberChars.length)];
  passwordChars[specialPosition] =
    specialChars[randomNumberGenerator(specialChars.length)];

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

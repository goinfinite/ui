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

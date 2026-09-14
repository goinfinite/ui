function registerAlpineState(stateFunction) {
  if (window.Alpine) {
    stateFunction();
    return;
  }

  document.addEventListener("alpine:init", stateFunction);
}

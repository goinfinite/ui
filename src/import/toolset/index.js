window.UiToolset = {
  CreateRandomPassword: createRandomPassword,
  ResolveApiResponseDisplay: resolveApiResponseDisplay,
  ToggleLoadingOverlay: toggleLoadingOverlay,
  RegisterAlpineState: registerAlpineState,
};

registerAlpineState(() => {
  window.UiToolset.JsonAjax = jsonAjax;
});

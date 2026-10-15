window.UiToolset = {
  CreateRandomPassword: createRandomPassword,
  ResolveApiResponseDisplay: resolveApiResponseDisplay,
  SelectDropdown: {
    controller: selectDropdownController,
    openUpwardResolver: selectDropdownOpenUpwardResolver,
  },
  ToggleLoadingOverlay: toggleLoadingOverlay,
  RegisterAlpineState: registerAlpineState,
};

registerAlpineState(() => {
  window.UiToolset.JsonAjax = jsonAjax;
});

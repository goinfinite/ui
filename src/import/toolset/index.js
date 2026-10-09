window.UiToolset = {
  CreateRandomPassword: createRandomPassword,
  ResolveApiResponseDisplay: resolveApiResponseDisplay,
  SelectDropdown: {
    openUpwardResolver: selectDropdownOpenUpwardResolver,
  },
  ToggleLoadingOverlay: toggleLoadingOverlay,
  RegisterAlpineState: registerAlpineState,
};

registerAlpineState(() => {
  window.UiToolset.JsonAjax = jsonAjax;
});

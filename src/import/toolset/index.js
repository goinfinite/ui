window.UiToolset = {
  CreateRandomPassword: createRandomPassword,
  ResolveApiResponseDisplay: resolveApiResponseDisplay,
  SelectDropdown: {
    controller: selectDropdownController,
    openUpwardResolver: selectDropdownOpenUpwardResolver,
  },
  ToggleLoadingOverlay: toggleLoadingOverlay,
  UnixTimestamp: {
    ToIsoDate: toIsoDate,
    ToRfc1123: toRfc1123,
  },
  RegisterAlpineState: registerAlpineState,
};

registerAlpineState(() => {
  window.UiToolset.JsonAjax = jsonAjax;
});

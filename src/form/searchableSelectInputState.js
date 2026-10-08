UiToolset.RegisterAlpineState(() => {
  const searchableSelectInputDefaultConfig = {
    items: [],
    isMultiSelect: false,
    remote: null,
  };

  const searchableSelectInputDefaultRemoteConfig = {
    url: "",
    queryParam: "q",
    minQueryLength: 3,
    debounceMs: 300,
  };

  Alpine.data("searchableSelectInput", (configScriptId) => ({
    config: { ...searchableSelectInputDefaultConfig },
    isOpen: false,
    openUpward: false,
    userInput: "",
    selectedLabelCache: new Map(),
    remoteOptions: [],
    isLoading: false,
    hasError: false,
    isAwaitingRemoteOptions: false,
    remoteRequestId: 0,
    remoteDebounceTimer: null,

    init() {
      if (!configScriptId) {
        return;
      }
      try {
        const configElement = document.getElementById(configScriptId);
        const parsedConfig = JSON.parse(configElement?.textContent || "{}");
        this.config = { ...this.config, ...parsedConfig };
        if (this.config.remote) {
          this.config.remote = {
            ...searchableSelectInputDefaultRemoteConfig,
            ...this.config.remote,
          };
        }
      } catch (parseError) {
        console.error(
          `SearchableSelectInputInvalidConfigJson: ${parseError.message}`,
        );
      }
    },

    destroy() {
      clearTimeout(this.remoteDebounceTimer);
    },

    openDropdown() {
      this.isOpen = true;
      this.$nextTick(() => {
        this.openUpward = UiToolset.SelectDropdown.openUpwardResolver(
          this.$refs.searchableSelectTrigger,
          this.$refs.searchableSelectDropdown,
        );
      });
    },

    closeDropdown() {
      this.isOpen = false;
    },

    get isRemote() {
      return this.config.remote !== null && this.config.remote !== undefined;
    },

    get filteredItems() {
      const filterText = this.userInput.trim().toLowerCase();
      if (filterText === "") {
        return this.config.items;
      }
      return this.config.items.filter((item) => {
        const searchableText = item.searchableText || item.label || "";
        return searchableText.toLowerCase().includes(filterText);
      });
    },

    isItemVisible(value) {
      return this.filteredItems.some((item) => item.value === value);
    },

    get hasNoMatches() {
      if (this.isRemote) {
        return (
          this.userInput.length >= this.config.remote.minQueryLength &&
          !this.isLoading &&
          !this.isAwaitingRemoteOptions &&
          !this.hasError &&
          this.remoteOptions.length === 0
        );
      }
      return this.filteredItems.length === 0;
    },

    get shouldPromptMinQueryLength() {
      return (
        this.isRemote &&
        this.userInput.length > 0 &&
        this.userInput.length < this.config.remote.minQueryLength
      );
    },

    get minQueryLengthPromptResolver() {
      return `Type at least ${this.config.remote.minQueryLength} characters`;
    },

    optionLabelResolver(value) {
      const cachedLabel = this.selectedLabelCache.get(value);
      if (cachedLabel !== undefined) {
        return cachedLabel;
      }
      const option =
        this.config.items.find((item) => item.value === value) ||
        this.remoteOptions.find((item) => item.value === value);
      if (option !== undefined) {
        return option.label;
      }
      return value;
    },

    optionLabelsTextResolver(values) {
      if (!Array.isArray(values) || values.length === 0) {
        return "";
      }
      return values.map((value) => this.optionLabelResolver(value)).join(", ");
    },

    onInputChanged() {
      if (!this.isRemote) {
        return;
      }
      clearTimeout(this.remoteDebounceTimer);
      const remoteConfig = this.config.remote;
      if (this.userInput.length < remoteConfig.minQueryLength) {
        this.remoteRequestId++;
        this.remoteOptions = [];
        this.isLoading = false;
        this.hasError = false;
        this.isAwaitingRemoteOptions = false;
        return;
      }
      this.isAwaitingRemoteOptions = true;
      this.remoteDebounceTimer = setTimeout(
        () => this.fetchRemoteOptions(),
        remoteConfig.debounceMs,
      );
    },

    async fetchRemoteOptions() {
      const remoteConfig = this.config.remote;
      const requestId = ++this.remoteRequestId;
      this.isLoading = true;
      this.hasError = false;
      this.isAwaitingRemoteOptions = false;
      try {
        const requestUrl = new URL(remoteConfig.url, window.location.href);
        requestUrl.searchParams.set(remoteConfig.queryParam, this.userInput);
        const response = await fetch(requestUrl.toString(), {
          headers: { Accept: "application/json" },
        });
        if (!response.ok) {
          throw new Error(`BadHttpResponseCode: ${response.status}`);
        }
        const responseData = await response.json();
        if (requestId !== this.remoteRequestId) {
          return;
        }
        const responseItems = Array.isArray(responseData)
          ? responseData
          : (responseData?.body ?? []);
        this.remoteOptions = responseItems.map((item) =>
          typeof item === "string"
            ? { value: item, label: item }
            : {
                value: String(item.value),
                label: String(item.label ?? item.value),
              },
        );
      } catch (error) {
        if (requestId !== this.remoteRequestId) {
          return;
        }
        this.remoteOptions = [];
        this.hasError = true;
        console.error(`SearchableSelectInputRemoteError: ${error.message}`);
      } finally {
        if (requestId === this.remoteRequestId) {
          this.isLoading = false;
        }
      }
    },
  }));
});

(() => {
  const unloadedPlaceholderSelector =
    "[data-demo-advanced-example-status='pending'], [data-demo-advanced-example-status='failed']";
  const assetRequests = new Map();

  function fetchPanelsAsset(assetPath) {
    if (!assetRequests.has(assetPath)) {
      const assetRequest = fetch(assetPath).then((response) => {
        if (!response.ok) {
          throw new Error(
            `AdvancedExampleAssetFetchFailed (${assetPath}): ${response.status}`,
          );
        }
        return response.json();
      });
      assetRequests.set(assetPath, assetRequest);
    }
    return assetRequests.get(assetPath);
  }

  function isExecutableScript(script) {
    return script.type === "" || script.type === "text/javascript";
  }

  function alreadyExecutedScriptRemover(panelFragment) {
    const executedScriptTexts = new Set();
    for (const pageScript of document.scripts) {
      executedScriptTexts.add(pageScript.textContent);
    }
    for (const script of panelFragment.querySelectorAll("script")) {
      if (!isExecutableScript(script)) {
        continue;
      }
      if (executedScriptTexts.has(script.textContent)) {
        script.remove();
        continue;
      }
      executedScriptTexts.add(script.textContent);
    }
  }

  async function loadPanel(placeholder) {
    placeholder.dataset.demoAdvancedExampleStatus = "loading";
    const assetPath = placeholder.dataset.demoAdvancedExampleAsset;
    const panelKey = placeholder.dataset.demoAdvancedExampleKey;
    try {
      const panelsByKey = await fetchPanelsAsset(assetPath);
      const panelHtml = panelsByKey[panelKey];
      if (panelHtml === undefined) {
        throw new Error(
          `AdvancedExamplePanelMissing (${assetPath}): ${panelKey}`,
        );
      }
      const panelFragment = document
        .createRange()
        .createContextualFragment(panelHtml);
      alreadyExecutedScriptRemover(panelFragment);
      placeholder.replaceChildren(panelFragment);
      window.htmx?.process(placeholder);
      placeholder.dataset.demoAdvancedExampleStatus = "loaded";
    } catch (loadError) {
      assetRequests.delete(assetPath);
      console.error(loadError.message);
      placeholder.dataset.demoAdvancedExampleStatus = "failed";
      placeholder.replaceChildren(
        Object.assign(document.createElement("p"), {
          className: "text-sm text-red-400",
          textContent:
            "The example failed to load. Close and open the panel to retry.",
        }),
      );
    }
  }

  document.addEventListener(
    "toggle",
    (event) => {
      const details = event.target;
      if (!(details instanceof HTMLDetailsElement) || !details.open) {
        return;
      }
      const placeholder = details.querySelector(unloadedPlaceholderSelector);
      if (!placeholder) {
        return;
      }
      void loadPanel(placeholder);
    },
    true,
  );
})();

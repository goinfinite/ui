UiToolset.RegisterAlpineState(() => {
  const remoteSearchableSelectDemoUrlMarker = "remote-searchable-select-demo";
  const remoteSearchableSelectErrorUrlMarker =
    "remote-searchable-select-demo-error";
  const remoteSearchableSelectDemoLatencyMs = 400;

  const demoCountryNames = [
    "Argentina",
    "Australia",
    "Austria",
    "Belgium",
    "Brazil",
    "Cameroon",
    "Canada",
    "Chile",
    "China",
    "Colombia",
    "Denmark",
    "Egypt",
    "Finland",
    "France",
    "Germany",
    "Ghana",
    "Greece",
    "India",
    "Indonesia",
    "Ireland",
    "Italy",
    "Japan",
    "Kenya",
    "Mexico",
    "Morocco",
    "Netherlands",
    "New Zealand",
    "Nigeria",
    "Norway",
    "Peru",
    "Poland",
    "Portugal",
    "Romania",
    "Senegal",
    "South Africa",
    "South Korea",
    "Spain",
    "Sweden",
    "Switzerland",
    "Thailand",
    "Turkey",
    "United Kingdom",
    "United States",
    "Uruguay",
    "Vietnam",
  ];

  const originalFetch = window.fetch.bind(window);
  window.fetch = async (input, init) => {
    const requestUrl =
      typeof input === "string" ? input : (input?.url ?? undefined);
    if (!requestUrl?.includes(remoteSearchableSelectDemoUrlMarker)) {
      return originalFetch(input, init);
    }

    await new Promise((resolve) =>
      setTimeout(resolve, remoteSearchableSelectDemoLatencyMs),
    );

    if (requestUrl.includes(remoteSearchableSelectErrorUrlMarker)) {
      return new Response(JSON.stringify({ body: [] }), {
        status: 500,
        headers: { "content-type": "application/json" },
      });
    }

    const request = new URL(requestUrl, window.location.href);
    const query = (
      request.searchParams.get("countryName") ??
      request.searchParams.get("q") ??
      ""
    ).toLowerCase();
    const matchingNames = demoCountryNames.filter((countryName) =>
      countryName.toLowerCase().includes(query),
    );
    const responseBody = JSON.stringify({
      body: matchingNames.map((countryName) => ({
        label: countryName,
        value: countryName,
      })),
    });
    return new Response(responseBody, {
      status: 200,
      headers: { "content-type": "application/json" },
    });
  };
});

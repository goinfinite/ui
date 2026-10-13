function examplePanel(page, section, title) {
  return page.locator(`${section} details`).filter({
    has: page
      .locator("summary span.font-bold")
      .getByText(title, { exact: true }),
  });
}

// The demo loads a panel on the details toggle event, and Lightpanda never
// fires it. The dispatch covers Lightpanda; on a real engine the loader
// ignores the second event because the panel already left the pending state.
async function openPanelDetails(details) {
  await details.locator(":scope > summary").click();
  await details.dispatchEvent("toggle");
}

async function waitForLazyPanelContent(target) {
  const lazyPlaceholder = target
    .locator("[data-demo-advanced-example-status]")
    .first();
  if ((await lazyPlaceholder.count()) === 0) {
    return;
  }
  await target
    .locator("[data-demo-advanced-example-status='loaded']")
    .first()
    .waitFor({ state: "attached" });
}

export async function openExamplePanel(page, section, title) {
  const target = examplePanel(page, section, title);
  const isPanelOpen = await target.evaluate((element) => element.open);
  if (!isPanelOpen) {
    await openPanelDetails(target);
  }
  await waitForLazyPanelContent(target);
  return target;
}

export async function openEveryExamplePanel(page, section) {
  const lazyPanels = page.locator(
    `${section} details:has(> div [data-demo-advanced-example-status])`,
  );
  for (const lazyPanel of await lazyPanels.all()) {
    await openPanelDetails(lazyPanel);
    await waitForLazyPanelContent(lazyPanel);
  }
}

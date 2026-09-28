function examplePanel(page, section, title) {
  return page.locator(`${section} details`).filter({
    has: page
      .locator("summary span.font-bold")
      .getByText(title, { exact: true }),
  });
}

export async function openExamplePanel(page, section, title) {
  const target = examplePanel(page, section, title);
  const isPanelOpen = await target.evaluate((element) => element.open);
  if (!isPanelOpen) {
    await target.locator("summary").click();
  }
  return target;
}

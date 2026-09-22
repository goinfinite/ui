import { expect, test } from "@playwright/test";

const textAreaSection = "#text-area-demo";
const collapsedHeightPx = 144;
const expandedHeightPx = 432;

async function waitForAlpineHydration(page) {
  await page.waitForFunction((section) => {
    const fieldset = document.querySelector(`${section} fieldset`);
    return window.Alpine && fieldset.classList.contains("mt-0");
  }, textAreaSection);
}

async function firstLineCenter(page, index) {
  return page
    .locator(`${textAreaSection} textarea`)
    .nth(index)
    .evaluate((ta) => {
      const rect = ta.getBoundingClientRect();
      const style = getComputedStyle(ta);
      return (
        rect.top +
        parseFloat(style.paddingTop) +
        parseFloat(style.lineHeight) / 2
      );
    });
}

async function firstIconCenter(page, index) {
  const fieldset = page.locator(`${textAreaSection} fieldset`).nth(index);
  await fieldset.locator("textarea").hover();
  return fieldset
    .locator(".ph-arrows-out")
    .first()
    .evaluate((icon) => {
      const rect = icon.getBoundingClientRect();
      return rect.top + rect.height / 2;
    });
}

async function heightOf(page, index) {
  return page
    .locator(`${textAreaSection} textarea`)
    .nth(index)
    .evaluate((ta) => ta.getBoundingClientRect().height);
}

async function toggleExpand(page, index) {
  await page.locator(`${textAreaSection} textarea`).nth(index).hover();
  await page
    .locator(`${textAreaSection} fieldset`)
    .nth(index)
    .locator(".ph-arrows-out, .ph-arrows-in")
    .first()
    .click();
}

async function expectStableHeight(page, index, expected) {
  await expect
    .poll(async () => heightOf(page, index), { timeout: 5000 })
    .toBeCloseTo(expected, 1);
}

test.describe("TextArea", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${textAreaSection.slice(1)}`);
    await waitForAlpineHydration(page);
  });

  test("@smoke default size renders at the md scale height", async ({
    page,
  }) => {
    expect(await heightOf(page, 0)).toBeCloseTo(collapsedHeightPx, 1);
  });

  test("@smoke expand toggle grows the box and collapse restores it", async ({
    page,
  }) => {
    await toggleExpand(page, 0);
    await expectStableHeight(page, 0, expandedHeightPx);
    await toggleExpand(page, 0);
    await expectStableHeight(page, 0, collapsedHeightPx);
  });

  test("read-only instance expands without a state path", async ({ page }) => {
    await toggleExpand(page, 1);
    await expectStableHeight(page, 1, expandedHeightPx);
  });

  test("action buttons reveal on focus and expand with the keyboard", async ({
    page,
  }) => {
    const fieldset = page.locator(`${textAreaSection} fieldset`).first();
    const actions = fieldset.locator("div.hidden").first();

    await expect(actions).toBeHidden();
    await fieldset.locator("textarea").focus();
    await expect(actions).toBeVisible();

    await fieldset
      .getByRole("button", { name: "Toggle text area height" })
      .press("Enter");
    await expectStableHeight(page, 0, expandedHeightPx);
  });

  test("floating icons sit on the first text line when empty", async ({
    page,
  }) => {
    await expect
      .poll(
        async () =>
          Math.abs(
            (await firstIconCenter(page, 0)) - (await firstLineCenter(page, 0)),
          ),
        { timeout: 5000 },
      )
      .toBeLessThanOrEqual(2);
  });

  test("floating icons sit on the first text line when filled", async ({
    page,
  }) => {
    await page
      .locator(`${textAreaSection} textarea`)
      .first()
      .fill("Lorem ipsum dolor sit amet");
    await expect
      .poll(
        async () =>
          Math.abs(
            (await firstIconCenter(page, 0)) - (await firstLineCenter(page, 0)),
          ),
        { timeout: 5000 },
      )
      .toBeLessThanOrEqual(2);
  });
});

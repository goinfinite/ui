import { test, expect } from "@playwright/test";

async function openToolsetPage(page) {
  await page.goto("/index.html");
  await page.waitForFunction(() => window.UiToolset?.JsonAjax);
}

test.describe("JsonAjax", () => {
  test.beforeEach(async ({ page }) => {
    await openToolsetPage(page);
  });

  test("@toolset GET with a nonempty payload rejects without calling fetch", async ({ page }) => {
    const rejection = await page.evaluate(async () => {
      let fetchCallCount = 0;
      window.fetch = () => {
        fetchCallCount++;
        return Promise.resolve();
      };
      try {
        await window.UiToolset.JsonAjax("GET", "/resource", { filter: "x" }, false);
      } catch (error) {
        return { message: error.message, fetchCallCount };
      }
      return { message: "", fetchCallCount };
    });

    expect(rejection).toEqual({
      message: "GetRequestPayloadNotAllowed",
      fetchCallCount: 0,
    });
  });

  test("@toolset a non-boolean shouldDisplayToast rejects without calling fetch", async ({ page }) => {
    const rejection = await page.evaluate(async () => {
      let fetchCallCount = 0;
      window.fetch = () => {
        fetchCallCount++;
        return Promise.resolve();
      };
      try {
        await window.UiToolset.JsonAjax("POST", "/resource", {}, "false");
      } catch (error) {
        return { message: error.message, fetchCallCount };
      }
      return { message: "", fetchCallCount };
    });

    expect(rejection).toEqual({
      message: "InvalidShouldDisplayToast",
      fetchCallCount: 0,
    });
  });

  test("@toolset a finished request hides the loading overlay", async ({ page }) => {
    const overlayStates = await page.evaluate(async () => {
      const overlay = document.getElementById("loading-overlay");
      const overlayIsShown = () => overlay.classList.contains("htmx-request");
      window.fetch = async () => ({
        headers: { get: () => "application/json" },
        ok: true,
        status: 200,
        json: async () => ({ body: {} }),
      });

      const request = window.UiToolset.JsonAjax("POST", "/resource", {}, false);
      const shownDuringRequest = overlayIsShown();
      await request;
      return { shownDuringRequest, shownAfterRequest: overlayIsShown() };
    });

    expect(overlayStates).toEqual({
      shownDuringRequest: true,
      shownAfterRequest: false,
    });
  });

  test("@toolset the loading overlay stays until the last concurrent request finishes", async ({ page }) => {
    const overlayStates = await page.evaluate(async () => {
      const overlay = document.getElementById("loading-overlay");
      const overlayIsShown = () => overlay.classList.contains("htmx-request");
      let resolveFirst;
      let resolveSecond;
      window.fetch = () =>
        new Promise((resolve) => {
          if (!resolveFirst) {
            resolveFirst = resolve;
            return;
          }
          resolveSecond = resolve;
        });
      const response = () => ({
        headers: { get: () => "application/json" },
        ok: true,
        status: 200,
        json: async () => ({ body: {} }),
      });

      const firstRequest = window.UiToolset.JsonAjax("POST", "/first", {}, false);
      const secondRequest = window.UiToolset.JsonAjax("POST", "/second", {}, false);
      const shownDuringBoth = overlayIsShown();
      resolveFirst(response());
      await firstRequest;
      const shownAfterFirst = overlayIsShown();
      resolveSecond(response());
      await secondRequest;
      return {
        shownDuringBoth,
        shownAfterFirst,
        shownAfterSecond: overlayIsShown(),
      };
    });

    expect(overlayStates).toEqual({
      shownDuringBoth: true,
      shownAfterFirst: true,
      shownAfterSecond: false,
    });
  });
});

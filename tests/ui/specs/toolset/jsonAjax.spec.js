import { expect, test } from "@playwright/test";

async function openToolsetPage(page) {
  await page.goto("/index.html");
  await page.waitForFunction(() => window.UiToolset?.JsonAjax);
}

test.describe("JsonAjax", () => {
  test.beforeEach(async ({ page }) => {
    await openToolsetPage(page);
  });

  test("@toolset GET with a nonempty payload rejects without calling fetch", async ({
    page,
  }) => {
    const rejection = await page.evaluate(async () => {
      let fetchCallCount = 0;
      window.fetch = () => {
        fetchCallCount++;
        return Promise.resolve();
      };
      try {
        await window.UiToolset.JsonAjax(
          "GET",
          "/resource",
          { filter: "x" },
          false,
        );
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

  test("@toolset a non-boolean shouldDisplayToast rejects without calling fetch", async ({
    page,
  }) => {
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

  test("@toolset a finished request hides the loading overlay", async ({
    page,
  }) => {
    const overlay = page.locator("#loading-overlay").first();
    await page.evaluate(() => {
      window.pendingFetch = new Promise((resolve) => {
        window.releaseFetch = () =>
          resolve({
            headers: { get: () => "application/json" },
            ok: true,
            status: 200,
            json: async () => ({ body: {} }),
          });
      });
      window.fetch = () => window.pendingFetch;
      window.pendingRequest = window.UiToolset.JsonAjax(
        "POST",
        "/resource",
        {},
        false,
      );
    });

    await expect(overlay).toBeVisible();

    await page.evaluate(async () => {
      window.releaseFetch();
      await window.pendingRequest;
    });
    await expect(overlay).toBeHidden();
  });

  test("@toolset the loading overlay stays until the last concurrent request finishes", async ({
    page,
  }) => {
    const overlay = page.locator("#loading-overlay").first();
    await page.evaluate(() => {
      const pendingFetches = [];
      window.releaseNextFetch = () => pendingFetches.shift()();
      window.fetch = () =>
        new Promise((resolve) => {
          pendingFetches.push(() =>
            resolve({
              headers: { get: () => "application/json" },
              ok: true,
              status: 200,
              json: async () => ({ body: {} }),
            }),
          );
        });
      window.pendingRequests = [
        window.UiToolset.JsonAjax("POST", "/first", {}, false),
        window.UiToolset.JsonAjax("POST", "/second", {}, false),
      ];
    });

    await expect(overlay).toBeVisible();

    await page.evaluate(async () => {
      window.releaseNextFetch();
      await window.pendingRequests[0];
    });
    await expect(overlay).toBeVisible();

    await page.evaluate(async () => {
      window.releaseNextFetch();
      await window.pendingRequests[1];
    });
    await expect(overlay).toBeHidden();
  });
});

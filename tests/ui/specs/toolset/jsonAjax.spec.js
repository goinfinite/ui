import { expect, test } from "../../fixture.js";

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

  test("@toolset a non-ok response rejects with the API message and status", async ({
    page,
  }) => {
    const rejection = await page.evaluate(async () => {
      window.fetch = () =>
        Promise.resolve({
          headers: { get: () => "application/json" },
          ok: false,
          status: 422,
          json: async () => ({ readableMessage: "Name is taken" }),
        });
      try {
        await window.UiToolset.JsonAjax("POST", "/resource", {}, false);
      } catch (error) {
        return {
          message: error.message,
          httpStatusCode: error.httpStatusCode,
          apiResponse: error.apiResponse,
        };
      }
      return { message: "", httpStatusCode: 0, apiResponse: null };
    });

    expect(rejection).toEqual({
      message: "Name is taken",
      httpStatusCode: 422,
      apiResponse: { readableMessage: "Name is taken" },
    });
  });

  test("@toolset a rejected fetch rejects and hides the loading overlay", async ({
    page,
  }) => {
    const overlay = page.locator("#loading-overlay").first();
    const rejection = await page.evaluate(async () => {
      window.fetch = () => Promise.reject(new Error("Network down"));
      try {
        await window.UiToolset.JsonAjax("POST", "/resource", {}, false);
      } catch (error) {
        return error.message;
      }
      return "";
    });

    expect(rejection).toBe("Network down");
    await expect(overlay).toBeHidden();
  });

  test("@toolset a non-JSON response rejects without parsing the body", async ({
    page,
  }) => {
    const rejection = await page.evaluate(async () => {
      let jsonCallCount = 0;
      window.fetch = () =>
        Promise.resolve({
          headers: { get: () => "text/html" },
          ok: true,
          status: 200,
          json: async () => {
            jsonCallCount++;
            return {};
          },
        });
      try {
        await window.UiToolset.JsonAjax("POST", "/resource", {}, false);
      } catch (error) {
        return { message: error.message, jsonCallCount };
      }
      return { message: "", jsonCallCount };
    });

    expect(rejection).toEqual({
      message: "UnexpectedResponseContentType",
      jsonCallCount: 0,
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

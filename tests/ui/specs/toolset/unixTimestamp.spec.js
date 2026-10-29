import { expect, test } from "../../fixture.js";

async function openToolsetPage(page) {
  await page.goto("/index.html");
  await page.waitForFunction(() => window.UiToolset?.UnixTimestamp);
}

test.describe("UnixTimestamp", () => {
  test.beforeEach(async ({ page }) => {
    await openToolsetPage(page);
  });

  test("@toolset seconds resolve to an ISO date and an RFC 1123 string", async ({
    page,
  }) => {
    const formatting = await page.evaluate(() => ({
      isoDate: window.UiToolset.UnixTimestamp.ToIsoDate(1609459200),
      rfc1123: window.UiToolset.UnixTimestamp.ToRfc1123(1609459200),
    }));

    expect(formatting).toEqual({
      isoDate: "2021-01-01",
      rfc1123: "Fri, 01 Jan 2021 00:00:00 GMT",
    });
  });

  test("@toolset a timestamp inside the day keeps the UTC day", async ({
    page,
  }) => {
    const formatting = await page.evaluate(() => ({
      isoDate: window.UiToolset.UnixTimestamp.ToIsoDate(1609462861),
      rfc1123: window.UiToolset.UnixTimestamp.ToRfc1123(1609462861),
    }));

    expect(formatting).toEqual({
      isoDate: "2021-01-01",
      rfc1123: "Fri, 01 Jan 2021 01:01:01 GMT",
    });
  });

  test("@toolset a non-number timestamp rejects", async ({ page }) => {
    const rejection = await page.evaluate(() => {
      try {
        window.UiToolset.UnixTimestamp.ToIsoDate("1609459200");
      } catch (error) {
        return error.message;
      }
      return "";
    });

    expect(rejection).toBe("InvalidUnixTimestamp");
  });

  test("@toolset an unrepresentable timestamp rejects", async ({ page }) => {
    const rejection = await page.evaluate(() => {
      try {
        window.UiToolset.UnixTimestamp.ToRfc1123(1e20);
      } catch (error) {
        return error.message;
      }
      return "";
    });

    expect(rejection).toBe("InvalidUnixTimestamp");
  });

  test("@toolset year zero formats and a negative year rejects", async ({
    page,
  }) => {
    const boundary = await page.evaluate(() => {
      const isoDate = window.UiToolset.UnixTimestamp.ToIsoDate(-62167219200);
      const rfc1123 = window.UiToolset.UnixTimestamp.ToRfc1123(-62167219200);
      let rejection = "";
      try {
        window.UiToolset.UnixTimestamp.ToIsoDate(-62167219201);
      } catch (error) {
        rejection = error.message;
      }
      return { isoDate, rfc1123, rejection };
    });

    expect(boundary).toEqual({
      isoDate: "0000-01-01",
      rfc1123: "Sat, 01 Jan 0000 00:00:00 GMT",
      rejection: "InvalidUnixTimestamp",
    });
  });

  test("@toolset a year above four digits rejects", async ({ page }) => {
    const rejections = await page.evaluate(() => {
      const messages = [];
      for (const formatter of [
        window.UiToolset.UnixTimestamp.ToIsoDate,
        window.UiToolset.UnixTimestamp.ToRfc1123,
      ]) {
        try {
          formatter(253402300800);
        } catch (error) {
          messages.push(error.message);
        }
      }
      return messages;
    });

    expect(rejections).toEqual([
      "InvalidUnixTimestamp",
      "InvalidUnixTimestamp",
    ]);
  });
});

import { expect, test } from "@playwright/test";
import { openExamplePanel } from "../../examplePanel.js";

const cardSection = "#card-demo";

test.describe("Card @structural", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html#${cardSection.slice(1)}`);
  });

  test("@smoke card renders the heading block with icon and description", async ({
    page,
  }) => {
    const defaultCard = page
      .locator(cardSection)
      .locator("div.rounded-lg")
      .filter({ hasText: "Default Card" })
      .first();

    await expect(defaultCard.locator("h2")).toHaveText("Default Card");
    await expect(defaultCard.locator("i.ph-cube")).toHaveCount(1);
    await expect(defaultCard).toContainText(
      "Rounded edges and an icon heading",
    );
    await expect(defaultCard).toContainText(
      "Cards group related content on a surface.",
    );
  });

  test("@smoke square card drops the rounded edges and keeps the action", async ({
    page,
  }) => {
    await openExamplePanel(page, cardSection, "Surface & Edges");
    const squareCard = page
      .locator(cardSection)
      .locator("div")
      .filter({ hasText: "Square Card" })
      .filter({ hasText: "Cards group related content" })
      .last();

    await expect(squareCard).toHaveCSS("border-radius", "0px");
    await expect(squareCard.locator("h2")).toHaveText("Square Card");
    await expect(
      squareCard.getByRole("button", { name: "Action" }),
    ).toBeVisible();
  });

  test("content-only card renders middle content without a heading", async ({
    page,
  }) => {
    await openExamplePanel(page, cardSection, "Content-Only Card");
    const contentOnlyCard = page
      .locator(cardSection)
      .locator("div")
      .filter({ hasText: "Content-only card body without a heading block" })
      .last();

    await expect(contentOnlyCard.locator("h2")).toHaveCount(0);
    await expect(contentOnlyCard).toContainText(
      "Content-only card body without a heading block",
    );
  });

  test("@smoke content gap example sets the card gap scale", async ({
    page,
  }) => {
    await openExamplePanel(page, cardSection, "Content Gap");
    const tightCard = page
      .locator(cardSection)
      .locator("div.grid > div")
      .filter({ hasText: "Tight Gap" })
      .first();
    const looseCard = page
      .locator(cardSection)
      .locator("div.grid > div")
      .filter({ hasText: "Loose Gap" })
      .first();

    await expect(tightCard).toHaveCSS("row-gap", "8px");
    await expect(looseCard).toHaveCSS("row-gap", "32px");
  });
});

import { openExamplePanel } from "../../examplePanel.js";
import { expect, test } from "../../fixture.js";

const fileUploadSection = "#file-upload-input-demo";
const archiveForm = `${fileUploadSection} form[hx-post="file-upload-demo"]`;
const archiveInput = `${archiveForm} input[type=file][name=archiveFiles]`;

function archiveFile(fileName, sizeBytes) {
  return {
    name: fileName,
    mimeType: "application/octet-stream",
    buffer: Buffer.alloc(sizeBytes),
  };
}

function submittedFileNames(page, fileInputSelector) {
  return page
    .locator(fileInputSelector)
    .evaluate((fileInput) =>
      Array.from(fileInput.files).map((pickedFile) => pickedFile.name),
    );
}

test.describe("FileUploadInput", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`/index.html${fileUploadSection}`);
  });

  test("@smoke picked files appear in the list with their sizes", async ({
    page,
  }) => {
    await page
      .locator(archiveInput)
      .setInputFiles([
        archiveFile("backup.tar", 2048),
        archiveFile("site.zip", 10),
      ]);

    const fileRows = page.locator(`${archiveForm} li`);
    await expect(fileRows).toHaveCount(2);
    await expect(fileRows.nth(0)).toContainText("backup.tar");
    await expect(fileRows.nth(0)).toContainText("2.0 KiB");
    await expect(fileRows.nth(1)).toContainText("10 B");
    expect(await submittedFileNames(page, archiveInput)).toEqual([
      "backup.tar",
      "site.zip",
    ]);
  });

  test("files that break a limit are rejected with a reason", async ({
    page,
  }) => {
    await page
      .locator(archiveInput)
      .setInputFiles([
        archiveFile("huge.tar", 6 * 1024 * 1024),
        archiveFile("notes.txt", 10),
        archiveFile("one.tar", 10),
        archiveFile("two.tar", 10),
        archiveFile("three.tar", 10),
        archiveFile("four.tar", 10),
      ]);

    const rejectionRows = page.locator(`${archiveForm} ul[aria-live] li`);
    await expect(rejectionRows).toHaveCount(3);
    await expect(rejectionRows.nth(0)).toContainText("Larger than 5.0 MiB");
    await expect(rejectionRows.nth(1)).toContainText("File type not accepted");
    await expect(rejectionRows.nth(2)).toContainText("File limit reached (3)");
    expect(await submittedFileNames(page, archiveInput)).toEqual([
      "one.tar",
      "two.tar",
      "three.tar",
    ]);
  });

  test("removing a file takes it out of the form submission", async ({
    page,
  }) => {
    await page
      .locator(archiveInput)
      .setInputFiles([
        archiveFile("keep.tar", 10),
        archiveFile("drop.tar", 10),
      ]);

    await page
      .locator(`${archiveForm} button[aria-label="Remove drop.tar"]`)
      .click();

    await expect(page.locator(`${archiveForm} li`)).toHaveCount(1);
    expect(await submittedFileNames(page, archiveInput)).toEqual(["keep.tar"]);
  });

  test("a successful upload marks the files uploaded and empties the input", async ({
    page,
  }) => {
    await page
      .locator(archiveInput)
      .setInputFiles([archiveFile("backup.tar", 4096)]);

    await page.locator(`${archiveForm} button[type=submit]`).click();

    await expect(
      page.locator(`${archiveForm} [role=progressbar]`).first(),
    ).toBeVisible();
    await expect(
      page.locator(`${archiveForm} li [role=img][aria-label=Uploaded]`),
    ).toBeVisible({ timeout: 5000 });
    expect(await submittedFileNames(page, archiveInput)).toEqual([]);
  });

  test("a single-file input replaces the pending file", async ({ page }) => {
    const singleFilePanel = await openExamplePanel(
      page,
      fileUploadSection,
      "Single File",
    );
    const avatarInput = singleFilePanel.locator("input[type=file]");
    const avatarFile = (fileName) => ({
      name: fileName,
      mimeType: "image/png",
      buffer: Buffer.alloc(10),
    });

    await avatarInput.setInputFiles([avatarFile("first.png")]);
    await avatarInput.setInputFiles([avatarFile("second.png")]);

    await expect(singleFilePanel.locator("li")).toHaveCount(1);
    await expect(singleFilePanel.locator("li")).toContainText("second.png");
  });

  test("a failed upload keeps the files for a retry", async ({ page }) => {
    const failurePanel = await openExamplePanel(
      page,
      fileUploadSection,
      "Upload Failure",
    );
    const failingInput = failurePanel.locator("input[type=file]");
    await failingInput.setInputFiles([archiveFile("report.tar", 2048)]);

    await failurePanel.locator("button[type=submit]").click();

    await expect(
      failurePanel.locator("li [role=img][aria-label='Upload failed']"),
    ).toBeVisible({ timeout: 5000 });
    const retryFileNames = await failingInput.evaluate((fileInput) =>
      Array.from(fileInput.files).map((pickedFile) => pickedFile.name),
    );
    expect(retryFileNames).toEqual(["report.tar"]);
  });
});

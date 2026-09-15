import { createServer } from "node:http";
import { readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";

const docsRoot = fileURLToPath(new URL("../../docs/", import.meta.url));
const listenPort = Number(process.argv[2] ?? 8377);

const contentTypes = {
  ".css": "text/css",
  ".html": "text/html",
  ".ico": "image/x-icon",
  ".jpg": "image/jpeg",
  ".js": "text/javascript",
  ".json": "application/json",
  ".png": "image/png",
  ".svg": "image/svg+xml",
  ".webp": "image/webp",
};

async function serveFile(request, response) {
  let fileName = "";
  try {
    const requestPath = decodeURIComponent(new URL(request.url, "http://localhost").pathname);
    fileName = requestPath.endsWith("/") ? `${requestPath}index.html` : requestPath;
  } catch {
    response.writeHead(400).end("bad request");
    return;
  }
  if (fileName.includes("..")) {
    response.writeHead(400).end("bad request");
    return;
  }
  try {
    const body = await readFile(`${docsRoot}${fileName}`);
    const extension = fileName.slice(fileName.lastIndexOf("."));
    response.writeHead(200, { "content-type": contentTypes[extension] ?? "application/octet-stream" });
    response.end(body);
  } catch {
    response.writeHead(404).end("not found");
  }
}

createServer(serveFile).listen(listenPort, "127.0.0.1");

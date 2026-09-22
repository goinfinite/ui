import { createHash, randomUUID } from "node:crypto";
import { mkdir, readFile, rename, writeFile } from "node:fs/promises";
import { createServer } from "node:http";
import { resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";

const docsRoot = resolve(
  fileURLToPath(new URL("../../docs/", import.meta.url)),
);
const assetCacheRoot = fileURLToPath(
  new URL("../.cache/demo-assets/", import.meta.url),
);
const listenPort = Number(process.argv[2] ?? 8377);
const proxyPrefix = "/proxy/";
const approvedAssetHosts = new Set([
  "cdn.jsdelivr.net",
  "fonts.googleapis.com",
  "fonts.gstatic.com",
  "img.shields.io",
  "sonarcloud.io",
  "ui.demo.goinfinite.net",
]);
const cacheFillAttempts = 2;
const cacheFillTimeoutMs = 8000;

const contentTypes = {
  ".css": "text/css",
  ".html": "text/html",
  ".ico": "image/x-icon",
  ".jpg": "image/jpeg",
  ".js": "text/javascript",
  ".json": "application/json",
  ".png": "image/png",
  ".svg": "image/svg+xml",
  ".ttf": "font/ttf",
  ".webp": "image/webp",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
};

function contentTypeOf(fileName) {
  const extension = fileName.slice(fileName.lastIndexOf("."));
  return contentTypes[extension] ?? "application/octet-stream";
}

function encodeProxyUrl(externalUrl) {
  return `${proxyPrefix}${externalUrl.replace("://", "/")}`;
}

function decodeProxyUrl(proxyPath) {
  const [scheme, ...pathParts] = proxyPath.split("/");
  return `${scheme}://${pathParts.join("/")}`;
}

function rewriteTagUrl(tag) {
  const urlMatch = tag.match(/\b(src|href)="(https?:\/\/[^"]{1,4096})"/);
  if (urlMatch === null) {
    return tag;
  }
  return tag
    .replace(urlMatch[2], () => encodeProxyUrl(urlMatch[2]))
    .replace(/ integrity="[^"]{0,4096}"/, "");
}

function rewriteExternalAssetUrls(html) {
  return html
    .replace(/<script\b[^>]{0,4096}>/g, rewriteTagUrl)
    .replace(/<link\b[^>]{0,4096}>/g, rewriteTagUrl)
    .replace(/<img\b[^>]{0,4096}>/g, rewriteTagUrl);
}

function rewriteAbsoluteCssUrls(css) {
  return css.replace(
    /url\((https?:\/\/[^)]{1,4096})\)/g,
    (_match, externalUrl) => `url(${encodeProxyUrl(externalUrl)})`,
  );
}

async function serveFile(request, response) {
  let fileName = "";
  try {
    const requestPath = decodeURIComponent(
      new URL(request.url, "http://localhost").pathname,
    );
    fileName = requestPath.endsWith("/")
      ? `${requestPath}index.html`
      : requestPath;
  } catch {
    response.writeHead(400).end("bad request");
    return;
  }
  const resolvedPath = resolve(docsRoot, `.${fileName}`);
  if (
    resolvedPath !== docsRoot &&
    !resolvedPath.startsWith(`${docsRoot}${sep}`)
  ) {
    response.writeHead(400).end("bad request");
    return;
  }
  try {
    const fileBody = await readFile(resolvedPath);
    const body = fileName.endsWith(".html")
      ? rewriteExternalAssetUrls(fileBody.toString("utf8"))
      : fileBody;
    response.writeHead(200, { "content-type": contentTypeOf(fileName) });
    response.end(body);
  } catch {
    response.writeHead(404).end("not found");
  }
}

async function readCachedAsset(cachePath) {
  try {
    const body = await readFile(cachePath);
    const contentType = (await readFile(`${cachePath}.type`, "utf8")).trim();
    return { body, contentType };
  } catch {
    return null;
  }
}

async function fetchExternalAsset(externalUrl) {
  for (let attempt = 1; attempt <= cacheFillAttempts; attempt++) {
    try {
      const upstream = await fetch(externalUrl, {
        redirect: "manual",
        signal: AbortSignal.timeout(cacheFillTimeoutMs),
      });
      if (upstream.ok) {
        return upstream;
      }
    } catch {}
  }
  return null;
}

async function fillAssetCache(externalUrl, cachePath) {
  const upstream = await fetchExternalAsset(externalUrl);
  if (upstream === null) {
    return null;
  }
  const contentType =
    upstream.headers.get("content-type") ?? "application/octet-stream";
  const rawBody = Buffer.from(await upstream.arrayBuffer());
  const body = contentType.includes("text/css")
    ? Buffer.from(rewriteAbsoluteCssUrls(rawBody.toString("utf8")))
    : rawBody;
  const temporaryPath = `${cachePath}.${randomUUID()}.tmp`;
  await writeFile(temporaryPath, body);
  await rename(temporaryPath, cachePath);
  await writeFile(`${cachePath}.type`, contentType);
  return { body, contentType };
}

function parseApprovedAssetUrl(externalUrl) {
  let assetUrl;
  try {
    assetUrl = new URL(externalUrl);
  } catch {
    return null;
  }
  const hasCredentials = assetUrl.username !== "" || assetUrl.password !== "";
  if (
    assetUrl.protocol !== "https:" ||
    hasCredentials ||
    !approvedAssetHosts.has(assetUrl.host)
  ) {
    return null;
  }
  return assetUrl;
}

async function serveProxiedAsset(request, response) {
  const proxyRequest = new URL(request.url, "http://localhost");
  const externalUrl = decodeProxyUrl(
    `${proxyRequest.pathname}${proxyRequest.search}`.slice(proxyPrefix.length),
  );
  if (parseApprovedAssetUrl(externalUrl) === null) {
    response.writeHead(403).end("asset not approved");
    return;
  }
  const cacheKey = createHash("sha256").update(externalUrl).digest("hex");
  const cachePath = `${assetCacheRoot}${cacheKey}`;
  const asset =
    (await readCachedAsset(cachePath)) ??
    (await fillAssetCache(externalUrl, cachePath));
  if (asset === null) {
    response.writeHead(502).end("asset unavailable");
    return;
  }
  response.writeHead(200, { "content-type": asset.contentType });
  response.end(asset.body);
}

await mkdir(assetCacheRoot, { recursive: true });

createServer((request, response) => {
  if (request.url?.startsWith(proxyPrefix)) {
    serveProxiedAsset(request, response).catch(() => {
      if (!response.headersSent) {
        response.writeHead(502);
      }
      response.end("asset unavailable");
    });
    return;
  }
  serveFile(request, response);
}).listen(listenPort, "127.0.0.1");

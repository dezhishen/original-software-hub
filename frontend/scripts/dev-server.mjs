import { createServer } from "node:http";
import { readFile, stat } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const rootDir = path.resolve(scriptDir, "..");
const distDir = path.join(rootDir, "dist");
const port = Number(process.env.PORT || 5173);

const MIME_TYPES = {
  ".html": "text/html; charset=utf-8",
  ".js": "application/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".gif": "image/gif",
  ".webp": "image/webp",
  ".ico": "image/x-icon"
};

function contentType(filePath) {
  const ext = path.extname(filePath).toLowerCase();
  return MIME_TYPES[ext] || "application/octet-stream";
}

function cacheHeaders(filePath) {
  const ext = path.extname(filePath).toLowerCase();
  if (ext === ".html" || ext === ".js" || ext === ".json" || ext === ".css") {
    return {
      "Cache-Control": "no-store, no-cache, must-revalidate, proxy-revalidate",
      Pragma: "no-cache",
      Expires: "0"
    };
  }

  return {
    "Cache-Control": "public, max-age=3600"
  };
}

function safePathFromUrl(urlPath) {
  const cleanPath = decodeURIComponent((urlPath || "/").split("?")[0]);
  const normalized = path.normalize(cleanPath).replace(/^([.][.][/\\])+/, "");
  const relative = normalized === "/" ? "index.html" : normalized.replace(/^[/\\]+/, "");
  return path.join(distDir, relative);
}

const server = createServer(async (req, res) => {
  try {
    const filePath = safePathFromUrl(req.url || "/");

    if (!filePath.startsWith(distDir)) {
      res.writeHead(403, { "Content-Type": "text/plain; charset=utf-8" });
      res.end("Forbidden");
      return;
    }

    let finalPath = filePath;
    let fileStat;

    try {
      fileStat = await stat(finalPath);
      if (fileStat.isDirectory()) {
        finalPath = path.join(finalPath, "index.html");
      }
    } catch {
      // fall through to 404 handling
    }

    const data = await readFile(finalPath);
    res.writeHead(200, {
      "Content-Type": contentType(finalPath),
      ...cacheHeaders(finalPath)
    });
    res.end(data);
  } catch {
    res.writeHead(404, { "Content-Type": "text/plain; charset=utf-8" });
    res.end("Not Found");
  }
});

server.listen(port, "127.0.0.1", () => {
  console.log(`[frontend] dev server running at http://127.0.0.1:${port}`);
});

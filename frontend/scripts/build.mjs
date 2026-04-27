import { mkdir, readdir, readFile, rm, stat, writeFile, copyFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { minify as minifyHtml } from "html-minifier-terser";
import { minify as minifyJs } from "terser";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const sourceDir = path.resolve(scriptDir, "..");
const outputDir = path.join(sourceDir, "dist");
const ignoredNames = new Set(["dist", "node_modules", "package.json", "package-lock.json", ".gitignore", "scripts"]);

async function collectFiles(dir, files = []) {
  const entries = await readdir(dir, { withFileTypes: true });
  for (const entry of entries) {
    if (ignoredNames.has(entry.name)) {
      continue;
    }

    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      await collectFiles(fullPath, files);
      continue;
    }

    files.push(fullPath);
  }

  return files;
}

async function ensureParentDir(filePath) {
  await mkdir(path.dirname(filePath), { recursive: true });
}

function getOutputPath(filePath) {
  return path.join(outputDir, path.relative(sourceDir, filePath));
}

async function buildHtml(sourcePath, outputPath) {
  const input = await readFile(sourcePath, "utf8");
  const output = await minifyHtml(input, {
    collapseBooleanAttributes: true,
    collapseInlineTagWhitespace: true,
    collapseWhitespace: true,
    keepClosingSlash: true,
    minifyCSS: true,
    minifyJS: true,
    processConditionalComments: true,
    removeAttributeQuotes: false,
    removeComments: true,
    removeEmptyAttributes: false,
    removeOptionalTags: false,
    removeRedundantAttributes: true,
    removeScriptTypeAttributes: true,
    removeStyleLinkTypeAttributes: true,
    sortAttributes: true,
    sortClassName: true,
    useShortDoctype: true
  });

  await ensureParentDir(outputPath);
  await writeFile(outputPath, output, "utf8");
}

async function buildJs(sourcePath, outputPath) {
  const input = await readFile(sourcePath, "utf8");
  const result = await minifyJs(input, {
    compress: {
      passes: 2,
      pure_getters: true
    },
    format: {
      comments: false
    },
    mangle: true,
    sourceMap: false
  });

  if (!result.code) {
    throw new Error(`JS 压缩失败: ${sourcePath}`);
  }

  await ensureParentDir(outputPath);
  await writeFile(outputPath, result.code, "utf8");
}

async function buildJson(sourcePath, outputPath) {
  const input = await readFile(sourcePath, "utf8");
  const parsed = JSON.parse(input);
  await ensureParentDir(outputPath);
  await writeFile(outputPath, JSON.stringify(parsed), "utf8");
}

async function copyAsset(sourcePath, outputPath) {
  await ensureParentDir(outputPath);
  await copyFile(sourcePath, outputPath);
}

async function buildFile(sourcePath) {
  const outputPath = getOutputPath(sourcePath);
  if (sourcePath.endsWith(".html")) {
    await buildHtml(sourcePath, outputPath);
    return;
  }

  if (sourcePath.endsWith(".js")) {
    await buildJs(sourcePath, outputPath);
    return;
  }

  if (sourcePath.endsWith(".json")) {
    await buildJson(sourcePath, outputPath);
    return;
  }

  await copyAsset(sourcePath, outputPath);
}

async function main() {
  await rm(outputDir, { recursive: true, force: true });
  await mkdir(outputDir, { recursive: true });

  const files = await collectFiles(sourceDir);
  await Promise.all(files.map(buildFile));

  const distStats = await stat(outputDir);
  if (!distStats.isDirectory()) {
    throw new Error("dist 目录生成失败");
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
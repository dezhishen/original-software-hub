// bundle-detail.js — software-detail.html 入口，所有模块合并（无构建工具版）
(function () {
  "use strict";

  // ── config ─────────────────────────────────────────────────────────────────
  const APP_DATA_SOURCE_CONFIG = window.APP_DATA_SOURCE_CONFIG;
  if (!APP_DATA_SOURCE_CONFIG) {
    throw new Error("缺少 APP_DATA_SOURCE_CONFIG，请先加载 config.js");
  }

  // ── shared/constants ───────────────────────────────────────────────────────
  const DEFAULT_JSONP_TIMEOUT_MS = 8000;
  const DEFAULT_JSONP_CALLBACK_PARAM = "callback";

  // ── shared/utils ───────────────────────────────────────────────────────────
  function escapeHtml(value) {
    return String(value)
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#39;");
  }

  function escapeAttr(value) {
    return escapeHtml(value).replaceAll("`", "&#96;");
  }

  function isObject(value) {
    return value !== null && typeof value === "object";
  }

  function toPositiveNumber(value, fallback) {
    const num = Number(value);
    return num > 0 ? num : fallback;
  }

  // ── infra/http/json-client ─────────────────────────────────────────────────
  async function fetchJson(url) {
    const response = await fetch(url, { cache: "no-cache" });
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${url}`);
    }
    return response.json();
  }

  // ── infra/http/jsonp-client ────────────────────────────────────────────────
  function fetchJsonp(source) {
    const url = source.url;
    const callbackParam = source.callbackParam || DEFAULT_JSONP_CALLBACK_PARAM;
    const timeoutMs = source.timeoutMs || DEFAULT_JSONP_TIMEOUT_MS;
    const callbackName = `__jsonp_cb_${Date.now()}_${Math.random().toString(36).slice(2)}`;

    return new Promise((resolve, reject) => {
      const script = document.createElement("script");
      let settled = false;
      const previousFixedCallback = window.callback;

      const cleanup = () => {
        settled = true;
        delete window[callbackName];
        if (previousFixedCallback === undefined) {
          delete window.callback;
        } else {
          window.callback = previousFixedCallback;
        }
        script.remove();
      };

      const timer = setTimeout(() => {
        if (!settled) {
          cleanup();
          reject(new Error(`JSONP timeout: ${url}`));
        }
      }, timeoutMs);

      window[callbackName] = (data) => {
        clearTimeout(timer);
        cleanup();
        resolve(data);
      };

      // 兼容固定 callback(...) 的静态 JSONP 文件。
      window.callback = (data) => {
        clearTimeout(timer);
        cleanup();
        resolve(data);
      };

      script.onerror = () => {
        clearTimeout(timer);
        cleanup();
        reject(new Error(`JSONP load error: ${url}`));
      };

      const separator = url.includes("?") ? "&" : "?";
      script.src = `${url}${separator}${encodeURIComponent(callbackParam)}=${callbackName}`;
      document.head.appendChild(script);
    });
  }

  // ── infra/http/source-client ───────────────────────────────────────────────
  async function fetchBySource(source) {
    if (source.mode === "jsonp") {
      return fetchJsonp(source);
    }
    return fetchJson(source.url);
  }

  // ── infra/config/config-normalizer ─────────────────────────────────────────
  function normalizeDataSourceConfig(config) {
    if (!isObject(config)) {
      throw new Error("数据源配置格式错误");
    }

    const endpoint = config.endpoint;
    if (!isObject(endpoint)) {
      throw new Error("数据源配置缺少 endpoint");
    }

    const type = String(endpoint.type || "json").trim().toLowerCase();
    const url = String(endpoint.url || "").trim();
    const indexPath = String(endpoint.indexPath || "index.json").trim() || "index.json";
    const callbackParam =
      String(endpoint.callbackParam || DEFAULT_JSONP_CALLBACK_PARAM).trim() ||
      DEFAULT_JSONP_CALLBACK_PARAM;
    const timeoutMs = toPositiveNumber(endpoint.timeoutMs, DEFAULT_JSONP_TIMEOUT_MS);

    if (!url) throw new Error("数据源配置缺少 endpoint.url");
    if (type !== "json" && type !== "jsonp") throw new Error("endpoint.type 仅支持 json 或 jsonp");

    return { endpoint: { type, url, indexPath, callbackParam, timeoutMs } };
  }

  // ── infra/config/config-repository ─────────────────────────────────────────
  function resolveUrl(baseUrl, path) {
    const base = String(baseUrl).endsWith("/") ? baseUrl : `${baseUrl}/`;
    return new URL(path, new URL(base, window.location.href)).toString();
  }

  async function loadDataSourceConfig() {
    const normalized = normalizeDataSourceConfig(APP_DATA_SOURCE_CONFIG);
    const { type, url, indexPath, callbackParam, timeoutMs } = normalized.endpoint;

    const indexSource = {
      mode: type,
      url: resolveUrl(url, indexPath),
      callbackParam,
      timeoutMs
    };

    const indexPayload = await fetchBySource(indexSource);

    const softwareListEntry = indexPayload?.softwareList;
    let listPath, listMode;

    if (isObject(softwareListEntry)) {
      listPath = String(softwareListEntry.path || "").trim();
      listMode = String(softwareListEntry.mode || type).trim().toLowerCase() || type;
    } else {
      listPath = String(softwareListEntry || "").trim();
      listMode = type;
    }

    if (!listPath) throw new Error("index 缺少 softwareList.path");

    return {
      softwareList: {
        mode: listMode,
        url: resolveUrl(url, listPath),
        callbackParam,
        timeoutMs
      },
      softwareSourceDefaults: {
        mode: type,
        baseUrl: url,
        callbackParam,
        timeoutMs
      }
    };
  }

  function resolveSoftwareVersionSource(source, defaults) {
    if (!isObject(defaults)) throw new Error("版本源默认配置无效");

    let mode = String(defaults.mode || "json").trim().toLowerCase() || "json";
    let sourcePath = "";
    let callbackParam =
      String(defaults.callbackParam || DEFAULT_JSONP_CALLBACK_PARAM).trim() ||
      DEFAULT_JSONP_CALLBACK_PARAM;
    let timeoutMs = toPositiveNumber(defaults.timeoutMs, DEFAULT_JSONP_TIMEOUT_MS);

    if (typeof source === "string") {
      sourcePath = source.trim();
    } else if (isObject(source)) {
      sourcePath = String(source.path || source.url || "").trim();
      const m = String(source.mode || source.type || mode).trim().toLowerCase();
      mode = m || mode;
      callbackParam = String(source.callbackParam || callbackParam).trim() || callbackParam;
      timeoutMs = toPositiveNumber(source.timeoutMs, timeoutMs);
    }

    if (!sourcePath) throw new Error("软件缺少 source 配置");
    if (mode !== "json" && mode !== "jsonp") throw new Error("source.mode 仅支持 json 或 jsonp");

    return { mode, url: resolveUrl(defaults.baseUrl, sourcePath), callbackParam, timeoutMs };
  }

  // ── domain/validators/payload-validator ────────────────────────────────────
  function normalizeSoftwareSource(source) {
    if (typeof source === "string") {
      return source.trim() || null;
    }
    if (!isObject(source)) return null;

    const path = String(source.path || source.url || "").trim();
    if (!path) return null;

    const mode = String(source.mode || source.type || "").trim().toLowerCase();
    if (mode && mode !== "json" && mode !== "jsonp") return null;

    const result = { path };
    if (mode) result.mode = mode;

    const callbackParam = String(source.callbackParam || "").trim();
    if (callbackParam) result.callbackParam = callbackParam;

    const timeoutMs = Number(source.timeoutMs);
    if (timeoutMs > 0) result.timeoutMs = timeoutMs;

    return result;
  }

  function normalizeSoftwareItem(item) {
    if (!isObject(item)) return null;

    const id = String(item.id || "").trim();
    const name = String(item.name || "").trim();
    const organization = String(item.organization || "").trim();
    const officialWebsite = String(item.officialWebsite || "").trim();
    const source = normalizeSoftwareSource(item.source);

    if (!id || !name || !organization || !officialWebsite || !source) return null;

    return {
      id,
      name,
      description: String(item.description || "").trim(),
      organization,
      officialWebsite,
      source
    };
  }

  function normalizeSoftwareListPayload(payload) {
    const items = Array.isArray(payload?.items) ? payload.items : [];
    return { items: items.map(normalizeSoftwareItem).filter(Boolean) };
  }

  function normalizeLink(item) {
    if (!isObject(item)) return null;
    const type = String(item.type || "").trim().toLowerCase();
    const label = String(item.label || "").trim();
    const url = String(item.url || "").trim();
    if (!label || !url || type !== "direct") return null;
    return { type, label, url };
  }

  function normalizeVariant(item) {
    if (!isObject(item)) return null;
    const architecture = String(item.architecture || "").trim();
    const platform = String(item.platform || "").trim();
    const links = Array.isArray(item.links)
      ? item.links.map(normalizeLink).filter(Boolean)
      : [];
    return { architecture, platform, links };
  }

  function normalizeVersion(item) {
    if (!isObject(item)) return null;
    const version = String(item.version || "").trim();
    const releaseDate = String(item.releaseDate || "").trim();
    const officialUrl = String(item.officialUrl || "").trim();
    const variants = Array.isArray(item.variants)
      ? item.variants.map(normalizeVariant).filter(Boolean)
      : [];
    return { version, releaseDate, officialUrl, variants };
  }

  function normalizeSoftwareVersionPayload(payload) {
    const raw = Array.isArray(payload?.versions)
      ? payload.versions
      : Array.isArray(payload)
        ? payload
        : [];
    return { versions: raw.map(normalizeVersion).filter(Boolean) };
  }

  // ── ui/renderers/software-detail-renderer ──────────────────────────────────
  function renderDetailEmpty(container, title, description) {
    if (!container) return;
    container.className = "min-h-[280px] grid place-items-center text-center";
    container.innerHTML = `<div><h2 class="text-xl font-semibold text-slate-700" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(title)}</h2><p class="mt-2 text-sm text-slate-500">${escapeHtml(description || "")}</p></div>`;
  }

  function renderSoftwareDetail({ container, software, versions }) {
    if (!container) return;
    if (!software) {
      renderDetailEmpty(container, "请选择一个软件", "");
      return;
    }

    container.className = "text-left";
    container.innerHTML = `
      <div class="mb-5 grid gap-2 border-b border-slate-200 pb-5">
        <h2 class="text-2xl font-bold text-slate-900" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(software.name)}</h2>
        <p class="text-sm leading-6 text-slate-600">${escapeHtml(software.description)}</p>
        <p class="text-sm text-slate-500">所属机构：${escapeHtml(software.organization)}</p>
        <a class="inline-flex w-fit items-center rounded-lg border border-brand-500/40 bg-brand-50 px-3 py-1.5 text-sm font-semibold text-brand-700 hover:bg-brand-100" target="_blank" rel="noopener noreferrer"
           href="${escapeAttr(software.officialWebsite)}">访问官网</a>
      </div>
      <div id="versionsContainer" class="grid gap-4"></div>
    `;

    const versionsContainer = container.querySelector("#versionsContainer");
    if (!versionsContainer) return;

    if (!Array.isArray(versions) || versions.length === 0) {
      versionsContainer.innerHTML =
        '<p class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-6 text-sm text-slate-600">暂无版本信息，请访问官网获取最新版本。</p>';
      return;
    }

    versions.forEach((v) => {
      const card = document.createElement("div");
      card.className = "overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:border-brand-500/50 hover:shadow-[0_10px_24px_rgba(15,157,132,0.15)]";

      const officialBtn = v.officialUrl
        ? `<a class="inline-flex items-center rounded-md border border-amber-300 bg-amber-50 px-2.5 py-1 text-xs font-semibold text-amber-700 hover:bg-amber-100" target="_blank" rel="noopener noreferrer"
              href="${escapeAttr(v.officialUrl)}">前往官网发布页</a>`
        : "";

      const variantRows = (v.variants || [])
        .map((variant) => {
          const directLinks = (variant.links || [])
            .map(
              (link) =>
                `<a class="inline-flex items-center rounded-md border border-brand-500/30 bg-brand-50 px-2.5 py-1 text-xs font-semibold text-brand-700 hover:bg-brand-100" target="_blank" rel="noopener noreferrer"
                    href="${escapeAttr(link.url)}">${escapeHtml(link.label)}</a>`
            )
            .join("");

          const directLinksHtml = directLinks
            ? `<div class="flex flex-wrap gap-2">${directLinks}</div>`
            : "暂无直链";

          return `
            <tr>
              <td class="whitespace-nowrap px-3 py-2 text-sm text-slate-700">${escapeHtml(variant.architecture || "-")}</td>
              <td class="whitespace-nowrap px-3 py-2 text-sm text-slate-700">${escapeHtml(variant.platform || "-")}</td>
              <td class="px-3 py-2 text-sm text-slate-700">${directLinksHtml}</td>
            </tr>`;
        })
        .join("");

      card.innerHTML = `
        <div class="flex flex-wrap items-center gap-2 border-b border-slate-200 bg-slate-50 px-3 py-3">
          <span class="rounded-full bg-brand-50 px-2.5 py-1 text-xs font-semibold text-brand-700" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(v.version || "-")}</span>
          <span class="text-xs text-slate-500">${escapeHtml(v.releaseDate || "")}</span>
          ${officialBtn}
        </div>
        ${
          variantRows
            ? `<div class="overflow-x-auto">
                <table class="min-w-full border-collapse">
                  <thead class="bg-slate-100">
                    <tr><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600">架构</th><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600">平台</th><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600">直接下载</th></tr>
                  </thead>
                  <tbody class="divide-y divide-slate-200">${variantRows}</tbody>
                </table>
              </div>`
            : '<p class="px-4 py-5 text-sm text-slate-600">该版本暂无构建信息。</p>'
        }
      `;

      versionsContainer.appendChild(card);
    });
  }

  // ── app/detail-bootstrap ───────────────────────────────────────────────────
  function hideOverlay(overlay) {
    if (!overlay) return;
    overlay.classList.add("hidden");
    setTimeout(() => { overlay.style.display = "none"; }, 400);
  }

  function setupSmartBackNavigation() {
    const backButton = document.querySelector("#smartBackButton");
    if (!backButton) return;

    backButton.addEventListener("click", () => {
      const fallbackUrl = new URL("./index.html", window.location.href).toString();
      const referrer = document.referrer || "";
      const hasHistory = window.history.length > 1;
      const isSameOriginReferrer = referrer.startsWith(window.location.origin);

      if (hasHistory && isSameOriginReferrer) {
        window.history.back();
        return;
      }

      window.location.href = fallbackUrl;
    });
  }

  async function bootstrapSoftwareDetailApp() {
    const container = document.querySelector("#softwareDetail");
    const overlay = document.querySelector("#loadingOverlay");
    const softwareId = new URLSearchParams(window.location.search).get("id") || "";

    setupSmartBackNavigation();

    if (!container) { hideOverlay(overlay); return; }

    if (!softwareId.trim()) {
      renderDetailEmpty(container, "缺少软件 ID", "请从软件目录页重新进入。");
      hideOverlay(overlay);
      return;
    }

    try {
      const dataSource = await loadDataSourceConfig();
      const rawList = await fetchBySource(dataSource.softwareList);
      const software = normalizeSoftwareListPayload(rawList).items.find(
        (item) => item.id === softwareId
      );

      if (!software) {
        renderDetailEmpty(container, "未找到软件", "请返回目录页选择有效的软件。");
        hideOverlay(overlay);
        return;
      }

      document.title = `${software.name} - Original Software Hub`;
      const breadcrumbCurrent = document.querySelector("#breadcrumbCurrent");
      if (breadcrumbCurrent) {
        breadcrumbCurrent.textContent = software.name;
      }

      const source = resolveSoftwareVersionSource(software.source, dataSource.softwareSourceDefaults);
      const rawVersions = await fetchBySource(source);
      const { versions } = normalizeSoftwareVersionPayload(rawVersions);

      renderSoftwareDetail({ container, software, versions });
      hideOverlay(overlay);
    } catch (error) {
      const message = error instanceof Error ? error.message : "未知错误";
      renderDetailEmpty(container, "加载失败", `详情加载失败：${message}`);
      hideOverlay(overlay);
    }
  }

  // ── entry ──────────────────────────────────────────────────────────────────
  bootstrapSoftwareDetailApp();
})();

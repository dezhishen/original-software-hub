// bundle-app.js — index.html 入口，所有模块合并（无构建工具版）
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
    const icon = String(item.icon || "").trim();
    const organization = String(item.organization || "").trim();
    const officialWebsite = String(item.officialWebsite || "").trim();
    const source = normalizeSoftwareSource(item.source);

    if (!id || !name || !organization || !officialWebsite || !source) return null;

    return {
      id,
      name,
      icon,
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

  // ── ui/pages/home-page ─────────────────────────────────────────────────────
  function getHomeDom() {
    return {
      list: document.querySelector("#softwareList"),
      searchInput: document.querySelector("#searchInput")
    };
  }

  function bindHomeEvents(dom, handlers) {
    dom.searchInput?.addEventListener("input", (event) => {
      handlers.onKeywordChange(event.target.value);
    });
  }

  function showLoadError(dom, error) {
    if (!dom.list) return;
    const message = error instanceof Error ? error.message : "未知错误";
    dom.list.innerHTML = `<p class="rounded-xl border border-rose-200 bg-rose-50 px-4 py-6 text-sm text-rose-700">数据加载失败：${escapeHtml(message)}</p>`;
  }

  // ── ui/renderers/software-list-renderer ────────────────────────────────────
  function renderSoftwareList({ container, softwares, keyword, onSelect }) {
    if (!container) return { filtered: [], firstId: "" };

    const kw = String(keyword || "").trim().toLowerCase();
    const filtered = softwares.filter((s) => {
      if (!kw) return true;
      return `${s.name} ${s.organization}`.toLowerCase().includes(kw);
    });

    container.innerHTML = "";

    if (filtered.length === 0) {
      container.innerHTML = '<p class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-6 text-sm text-slate-600">没有匹配的软件，请尝试其他关键词。</p>';
      return { filtered, firstId: "" };
    }

    filtered.forEach((software) => {
      const card = document.createElement("article");
      card.className = "cursor-pointer rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:-translate-y-0.5 hover:border-brand-500/50 hover:shadow-[0_10px_24px_rgba(15,157,132,0.15)]";
      const iconMarkup = renderSoftwareIcon(software);
      card.innerHTML = `
        <div class="mb-3 flex items-center gap-3">
          ${iconMarkup}
          <h3 class="text-lg font-semibold text-slate-900" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(software.name)}</h3>
        </div>
        <p class="mt-2 text-sm leading-6 text-slate-600">${escapeHtml(software.description)}</p>
        <p class="mt-2 text-xs text-slate-500">机构：${escapeHtml(software.organization)}</p>
      `;
      card.addEventListener("click", () => onSelect(software.id));
      container.appendChild(card);
    });

    return { filtered, firstId: filtered[0]?.id || "" };
  }

  function renderSoftwareIcon(software) {
    const icon = String(software?.icon || "").trim();
    if (/^https?:\/\//i.test(icon) || icon.startsWith("/") || icon.startsWith("./")) {
      return `<img class="h-9 w-9 rounded-lg border border-slate-200 bg-white p-1 object-contain" src="${escapeAttr(icon)}" alt="${escapeAttr(software.name)} 图标" loading="lazy" />`;
    }

    const fallback = escapeHtml(icon || software.name.slice(0, 1).toUpperCase() || "?");
    return `<span class="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-sm font-semibold text-slate-600">${fallback}</span>`;
  }

  // ── app/bootstrap ──────────────────────────────────────────────────────────
  function hideOverlay(overlay) {
    if (!overlay) return;
    overlay.classList.add("hidden");
    setTimeout(() => { overlay.style.display = "none"; }, 400);
  }

  async function bootstrapHomeApp() {
    const dom = getHomeDom();
    const loadingOverlay = document.querySelector("#loadingOverlay");
    const state = { softwares: [], keyword: "" };

    bindHomeEvents(dom, {
      onKeywordChange(keyword) {
        state.keyword = keyword.trim().toLowerCase();
        renderAll();
      }
    });

    try {
      const dataSource = await loadDataSourceConfig();
      const rawList = await fetchBySource(dataSource.softwareList);
      state.softwares = normalizeSoftwareListPayload(rawList).items;
      renderAll();
      hideOverlay(loadingOverlay);
    } catch (error) {
      hideOverlay(loadingOverlay);
      showLoadError(dom, error);
    }

    function renderAll() {
      renderSoftwareList({
        container: dom.list,
        softwares: state.softwares,
        keyword: state.keyword,
        onSelect(softwareId) {
          const nextUrl = new URL("./software-detail.html", window.location.href);
          nextUrl.searchParams.set("id", softwareId);
          window.location.href = nextUrl.toString();
        }
      });
    }
  }

  // ── entry ──────────────────────────────────────────────────────────────────
  bootstrapHomeApp();
})();

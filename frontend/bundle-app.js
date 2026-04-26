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

  function relativeTime(dateStr) {
    if (!dateStr) return "";
    const raw = String(dateStr).trim();
    const isDateOnly = /^\d{4}-\d{2}-\d{2}$/.test(raw);
    const date = new Date(
      isDateOnly
        ? `${raw}T12:00:00+08:00`
        : raw
    );
    if (isNaN(date)) return escapeHtml(String(dateStr));

    const diffMs = Date.now() - date.getTime();
    const diffMin = Math.floor(diffMs / 60000);
    const diffHr  = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHr  / 24);
    let rel;
    if (diffMs < 0)         rel = "刚刚";
    else if (isDateOnly && diffDay < 1) rel = "今天";
    else if (diffMin < 1)   rel = "刚刚";
    else if (diffMin < 60)  rel = `${diffMin} 分钟前`;
    else if (diffHr  < 24)  rel = `${diffHr} 小时前`;
    else if (diffDay < 30)  rel = `${diffDay} 天前`;
    else if (diffDay < 365) rel = `${Math.floor(diffDay / 30)} 个月前`;
    else                    rel = `${Math.floor(diffDay / 365)} 年前`;

    const realTime = date.toLocaleString("zh-CN", {
      timeZone: "Asia/Shanghai",
      year: "numeric", month: "2-digit", day: "2-digit",
      hour: "2-digit", minute: "2-digit"
    });
    return `<time datetime="${escapeAttr(date.toISOString())}" title="${escapeAttr(realTime + "（北京时间）")}" class="cursor-help underline decoration-dotted decoration-slate-400">${escapeHtml(rel)}</time>`;
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
      generatedAt: String(indexPayload?.meta?.generatedAt || "").trim(),
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
      pinyin: String(item.pinyin || "").trim().toLowerCase(),
      icon,
      description: String(item.description || "").trim(),
      organization,
      officialWebsite,
      tags: Array.isArray(item.tags)
        ? item.tags.map((tag) => String(tag || "").trim()).filter(Boolean)
        : [],
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
  function renderSoftwareList({ container, softwares, keyword, onSelect, onTagSelect }) {
    if (!container) return { filtered: [], firstId: "" };

    const kw = String(keyword || "").trim().toLowerCase();
    const tagKw = kw.startsWith("#") ? kw.slice(1).trim() : "";
    const filtered = softwares.filter((s) => {
      if (!kw) return true;
      const tags = Array.isArray(s.tags) ? s.tags : [];
      if (tagKw) {
        return tags.some((tag) => String(tag || "").toLowerCase().includes(tagKw));
      }
      return `${s.name} ${s.organization} ${tags.join(" ")} ${s.pinyin || ""}`.toLowerCase().includes(kw);
    });

    container.innerHTML = "";

    if (filtered.length === 0) {
      container.innerHTML = '<p class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-6 text-sm text-slate-600">没有匹配的软件，请尝试其他关键词。</p>';
      return { filtered, firstId: "" };
    }

    filtered.forEach((software) => {
      const card = document.createElement("article");
      card.className = "relative overflow-hidden cursor-pointer rounded-xl border border-slate-200/85 bg-white/60 p-4 shadow-[0_14px_30px_rgba(15,70,56,0.14)] backdrop-blur-md transition hover:-translate-y-1 hover:border-brand-500/55 hover:shadow-[0_20px_40px_rgba(15,157,132,0.22)] dark:border-slate-700/80 dark:bg-slate-800/60 dark:shadow-[0_14px_30px_rgba(2,6,23,0.45)] dark:hover:shadow-[0_20px_40px_rgba(15,157,132,0.25)]";
      const iconMarkup = renderSoftwareIcon(software);
      const rawIcon = String(software?.icon || "").trim();
      const bgWatermark = rawIcon
        ? `<div class="pointer-events-none absolute -bottom-3 -right-3 h-28 w-28 select-none opacity-[0.08] dark:opacity-[0.05]" style="background-image:url('${escapeAttr(rawIcon)}');background-size:contain;background-repeat:no-repeat;background-position:center;"></div>`
        : "";
      const tagsMarkup = (software.tags || [])
        .map(tag => `<button type="button" data-tag="${escapeAttr(tag)}" class="inline-block rounded-full bg-brand-50/80 px-2 py-0.5 text-xs font-medium text-brand-700 transition hover:bg-brand-100 dark:bg-slate-700/50 dark:text-brand-400 dark:hover:bg-slate-700">#${escapeHtml(tag)}</button>`)
        .join(" ");
      card.innerHTML = `
        ${bgWatermark}
        <div class="relative mb-3 flex items-center gap-3">
          ${iconMarkup}
          <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(software.name)}</h3>
        </div>
        <p class="relative mt-2 text-sm leading-6 text-slate-600 dark:text-slate-400">${escapeHtml(software.description)}</p>
        ${tagsMarkup ? `<div class="relative mt-3 flex flex-wrap gap-1.5">${tagsMarkup}</div>` : ""}
        <p class="relative mt-2 text-xs text-slate-500 dark:text-slate-500">机构：${escapeHtml(software.organization)}</p>
      `;
      card.addEventListener("click", () => onSelect(software.id));
      card.querySelectorAll("[data-tag]").forEach((tagBtn) => {
        tagBtn.addEventListener("click", (event) => {
          event.stopPropagation();
          onTagSelect?.(tagBtn.getAttribute("data-tag") || "");
        });
      });
      container.appendChild(card);
    });

    return { filtered, firstId: filtered[0]?.id || "" };
  }

  function renderAppFooter(updatedAt) {
    const footer = document.querySelector("#appFooter");
    if (!footer) return;
    const timeHtml = updatedAt ? `数据更新时间：${relativeTime(updatedAt)}` : "";
    footer.innerHTML = [
      timeHtml,
      `<p>本站数据来源于各软件官方渠道，所有下载链接均指向官方或官方镜像地址，仅供参考。本站不对链接可用性、文件安全性及版本准确性作任何保证，请自行核实后使用。</p>`
    ].filter(Boolean).join("\n");
  }

  function initDarkMode() {
    const html = document.documentElement;
    const btn = document.querySelector("#darkModeToggle");
    const saved = localStorage.getItem("osh-dark-mode");
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    const isDark = saved !== null ? saved === "1" : prefersDark;
    const sunIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>`;
    const moonIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>`;
    function applyDark(dark) {
      html.classList.toggle("dark", dark);
      if (!btn) return;
      btn.innerHTML = dark ? sunIcon : moonIcon;
      btn.setAttribute("aria-label", dark ? "切换白天模式" : "切换夜间模式");
    }
    applyDark(isDark);
    btn?.addEventListener("click", () => {
      const nowDark = !html.classList.contains("dark");
      localStorage.setItem("osh-dark-mode", nowDark ? "1" : "0");
      applyDark(nowDark);
    });
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
    initDarkMode();
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
      renderAppFooter(dataSource.generatedAt);
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
        onTagSelect(tag) {
          const cleanTag = String(tag || "").trim();
          if (!cleanTag) return;
          state.keyword = `#${cleanTag.toLowerCase()}`;
          if (dom.searchInput) dom.searchInput.value = `#${cleanTag}`;
          renderAll();
        },
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

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
  function renderAppFooter(updatedAt) {
    const footer = document.querySelector("#appFooter");
    if (!footer) return;
    const timeHtml = updatedAt ? `数据更新时间：${relativeTime(updatedAt)}` : "";
    footer.innerHTML = [
      timeHtml,
      `<p>本站数据来源于各软件官方渠道，所有下载链接均指向官方或官方镜像地址，仅供参考。本站不对链接可用性、文件安全性及版本准确性作任何保证，请自行核实后使用，建议优先前往软件官网下载安装。</p>`
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

  function renderDetailEmpty(container, title, description) {
    if (!container) return;
    container.className = "min-h-[280px] grid place-items-center text-center";
    container.innerHTML = `<div><h2 class="text-xl font-semibold text-slate-700" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(title)}</h2><p class="mt-2 text-sm text-slate-500">${escapeHtml(description || "")}</p></div>`;
  }

  function detectCurrentPlatform() {
    const ua = String(navigator.userAgent || "").toLowerCase();
    const platform = String((navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || "").toLowerCase();
    const source = `${ua} ${platform}`;

    if (/iphone|ipad|ipod/.test(source)) return { id: "ios", label: "iOS" };
    if (/android/.test(source)) return { id: "android", label: "Android" };
    if (/mac|darwin/.test(source)) return { id: "macos", label: "macOS" };
    if (/win/.test(source)) return { id: "windows", label: "Windows" };
    if (/linux|x11/.test(source)) return { id: "linux", label: "Linux" };
    return { id: "web", label: "Web" };
  }

  function detectCurrentArchitecture() {
    const ua = String(navigator.userAgent || "").toLowerCase();
    const uaArch = String((navigator.userAgentData && navigator.userAgentData.architecture) || "").toLowerCase();
    const source = `${uaArch} ${ua}`;

    if (/arm64|aarch64|armv8/.test(source)) return { id: "arm64", label: "ARM64" };
    if (/x86_64|win64|wow64|amd64|x64/.test(source)) return { id: "x64", label: "x64" };
    if (/i[3-6]86|x86/.test(source)) return { id: "x86", label: "x86" };
    return { id: "universal", label: "通用" };
  }

  function platformMatchesCurrent(variantPlatform, currentPlatformId) {
    const p = String(variantPlatform || "").toLowerCase();
    if (!p) return false;

    switch (currentPlatformId) {
      case "windows":
        return p.includes("windows");
      case "macos":
        return p.includes("mac");
      case "linux":
        return p.includes("linux");
      case "android":
        return p.includes("android");
      case "ios":
        return p.includes("ios") || p.includes("iphone") || p.includes("ipad");
      case "web":
        return p.includes("web");
      default:
        return false;
    }
  }

  function architectureScore(variantArchitecture, currentArchId) {
    const arch = String(variantArchitecture || "").toLowerCase();
    const has = (keyword) => arch.includes(keyword);

    if (has("universal") || has("通用")) return 85;

    switch (currentArchId) {
      case "arm64":
        if (has("arm64") || has("arm")) return 100;
        if (has("x64") || has("amd64")) return 55;
        if (has("x86") || has("32")) return 35;
        break;
      case "x64":
        if (has("x64") || has("amd64")) return 100;
        if (has("x86/x64")) return 100;
        if (has("x86") || has("32")) return 70;
        if (has("arm64") || has("arm")) return 40;
        break;
      case "x86":
        if (has("x86") || has("32")) return 100;
        if (has("x64") || has("amd64")) return 60;
        if (has("arm64") || has("arm")) return 30;
        break;
      default:
        return 50;
    }

    return 50;
  }

  function renderSoftwareDetail({ container, software, versions }) {
    if (!container) return;
    if (!software) {
      renderDetailEmpty(container, "请选择一个软件", "");
      return;
    }

    const currentPlatform = detectCurrentPlatform();
    const currentArchitecture = detectCurrentArchitecture();

    container.className = "text-left";
    container.innerHTML = `
      <div class="mb-5 grid gap-2 border-b border-slate-200 pb-5">
        <h2 class="text-2xl font-semibold text-slate-900" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(software.name)}</h2>
        <p class="text-sm leading-6 text-slate-600">${escapeHtml(software.description)}</p>
        <p class="text-sm text-slate-500">所属机构：${escapeHtml(software.organization)}</p>
        <p class="text-xs text-slate-500 dark:text-slate-400">当前检测环境：<span class="rounded-full bg-brand-50 px-2 py-0.5 font-medium text-brand-700 dark:bg-slate-700/60 dark:text-brand-300">${escapeHtml(currentPlatform.label)} / ${escapeHtml(currentArchitecture.label)}</span></p>
        <a class="inline-flex w-fit items-center rounded-lg border border-brand-500/35 bg-brand-50 px-3 py-1.5 text-sm font-medium text-brand-700 hover:bg-brand-100" target="_blank" rel="noopener noreferrer"
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
      card.className = "overflow-hidden rounded-xl border border-slate-200/90 bg-white shadow-[0_4px_12px_rgba(15,70,56,0.08)] transition hover:-translate-y-0.5 hover:border-brand-500/40 hover:shadow-[0_8px_16px_rgba(15,157,132,0.12)] dark:border-slate-700/80 dark:bg-slate-800/90 dark:shadow-[0_6px_16px_rgba(2,6,23,0.35)] dark:hover:shadow-[0_10px_20px_rgba(15,157,132,0.18)]";

      const officialBtn = v.officialUrl
        ? `<a class="inline-flex items-center rounded-md border border-amber-300 bg-amber-50 px-2.5 py-1 text-xs font-semibold text-amber-700 hover:bg-amber-100 dark:border-amber-700/60 dark:bg-amber-900/25 dark:text-amber-300 dark:hover:bg-amber-900/40" target="_blank" rel="noopener noreferrer"
              href="${escapeAttr(v.officialUrl)}">前往官网发布页</a>`
        : "";

      const sortedVariants = [...(v.variants || [])].sort((a, b) => {
        const aPlatformScore = platformMatchesCurrent(a.platform, currentPlatform.id) ? 1 : 0;
        const bPlatformScore = platformMatchesCurrent(b.platform, currentPlatform.id) ? 1 : 0;
        if (aPlatformScore !== bPlatformScore) return bPlatformScore - aPlatformScore;

        if (aPlatformScore === 1) {
          const aArchScore = architectureScore(a.architecture, currentArchitecture.id);
          const bArchScore = architectureScore(b.architecture, currentArchitecture.id);
          if (aArchScore !== bArchScore) return bArchScore - aArchScore;
        }

        return 0;
      });

      const first = sortedVariants[0];
      const hasCurrentDeviceRow = !!first && platformMatchesCurrent(first.platform, currentPlatform.id);

      const variantRows = sortedVariants
        .map((variant, index) => {
          const isCurrentDevice = hasCurrentDeviceRow && index === 0;
          const directLinks = (variant.links || [])
            .map(
              (link) =>
                `<a class="inline-flex items-center rounded-md border border-brand-500/30 bg-brand-50 px-2.5 py-1 text-xs font-semibold text-brand-700 hover:bg-brand-100 dark:border-brand-500/40 dark:bg-slate-700/50 dark:text-brand-300 dark:hover:bg-slate-700" target="_blank" rel="noopener noreferrer"
                    href="${escapeAttr(link.url)}">${escapeHtml(link.label)}</a>`
            )
            .join("");

          const directLinksHtml = directLinks
            ? `<div class="flex flex-wrap gap-2">${directLinks}</div>`
            : "暂无直链";

          return `
            <tr class="bg-white even:bg-slate-50 hover:bg-slate-100/70 dark:bg-slate-800 dark:even:bg-slate-800/75 dark:hover:bg-slate-700/60 ${isCurrentDevice ? "font-semibold" : ""}">
              <td class="whitespace-nowrap px-3 py-2 text-sm text-slate-700 dark:text-slate-200">${escapeHtml(variant.architecture || "-")}</td>
              <td class="whitespace-nowrap px-3 py-2 text-sm text-slate-700 dark:text-slate-200">${escapeHtml(variant.platform || "-")}${isCurrentDevice ? ' <span class="ml-1 text-[11px] font-semibold text-brand-700 dark:text-brand-300">当前设备</span>' : ""}</td>
              <td class="px-3 py-2 text-sm text-slate-700 dark:text-slate-200">${directLinksHtml}</td>
            </tr>`;
        })
        .join("");

      card.innerHTML = `
        <div class="flex flex-wrap items-center gap-2 border-b border-slate-200 bg-slate-50 px-3 py-3 dark:border-slate-700 dark:bg-slate-900/45">
          <span class="rounded-full bg-brand-50 px-2.5 py-1 text-xs font-medium text-brand-700 dark:bg-slate-700/50 dark:text-brand-300" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(v.version || "-")}</span>
          <span class="text-xs text-slate-500 dark:text-slate-400">${escapeHtml(v.releaseDate || "")}</span>
          ${officialBtn}
        </div>
        ${
          variantRows
            ? `<div class="overflow-x-auto">
                <table class="min-w-full border-collapse">
                  <thead class="bg-slate-100 dark:bg-slate-900/55">
                    <tr><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600 dark:text-slate-300">架构</th><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600 dark:text-slate-300">平台</th><th class="px-3 py-2 text-left text-xs font-semibold tracking-wide text-slate-600 dark:text-slate-300">直接下载</th></tr>
                  </thead>
                  <tbody class="divide-y divide-slate-200 dark:divide-slate-700 dark:bg-slate-800">${variantRows}</tbody>
                </table>
              </div>`
            : '<p class="px-4 py-5 text-sm text-slate-600 dark:text-slate-400">该版本暂无构建信息。</p>'
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
    initDarkMode();
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
      renderAppFooter(dataSource.generatedAt);
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

// bundle-app.js — index.html 入口，所有模块合并（无构建工具版）
(function () {
  "use strict";

  // ── config ─────────────────────────────────────────────────────────────────
  const APP_DATA_SOURCE_CONFIG = window.APP_DATA_SOURCE_CONFIG;
  if (!APP_DATA_SOURCE_CONFIG) {
    throw new Error("缺少 APP_DATA_SOURCE_CONFIG，请先加载 config.js");
  }
  const dataRepository = window.OSH_DATA_REPOSITORY;
  if (!dataRepository) {
    throw new Error("缺少 OSH_DATA_REPOSITORY，请先加载 bundle-shared.js");
  }

  // ── shared/constants ───────────────────────────────────────────────────────
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

  // ── ui/pages/home-page ─────────────────────────────────────────────────────
  function getAppDom() {
    return {
      homeHero: document.querySelector("#homeHero"),
      homeHeroMeta: document.querySelector("#homeHeroMeta"),
      detailHero: document.querySelector("#detailHero"),
      detailHeroTitle: document.querySelector("#detailHeroTitle"),
      detailHeroDescription: document.querySelector("#detailHeroDescription"),
      detailHeroMeta: document.querySelector("#detailHeroMeta"),
      homeSection: document.querySelector("#homeSection"),
      detailSection: document.querySelector("#detailSection"),
      detailContainer: document.querySelector("#softwareDetail"),
      list: document.querySelector("#softwareList"),
      searchInput: document.querySelector("#searchInput"),
      loadingOverlay: document.querySelector("#loadingOverlay"),
      loadingMessage: document.querySelector("#loadingMessage")
    };
  }

  function bindAppEvents(dom, handlers) {
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
      card.className = "relative overflow-hidden cursor-pointer rounded-xl border border-slate-200/90 bg-white/92 p-4 shadow-[0_6px_16px_rgba(15,70,56,0.08)] transition hover:-translate-y-0.5 hover:border-brand-500/45 hover:shadow-[0_10px_20px_rgba(15,157,132,0.14)] dark:border-slate-700/80 dark:bg-slate-800/88 dark:shadow-[0_6px_16px_rgba(2,6,23,0.35)] dark:hover:shadow-[0_10px_20px_rgba(15,157,132,0.18)]";
      const iconMarkup = renderSoftwareIcon(software);
      const rawIcon = String(software?.icon || "").trim();
      const bgWatermark = rawIcon
        ? `<div class="pointer-events-none absolute -bottom-3 -right-3 h-24 w-24 select-none opacity-[0.05] dark:opacity-[0.035]" style="background-image:url('${escapeAttr(rawIcon)}');background-size:contain;background-repeat:no-repeat;background-position:center;"></div>`
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

  function renderHomeHeroMeta(dom, updatedAt) {
    if (!dom?.homeHeroMeta) return;
    if (!updatedAt) {
      dom.homeHeroMeta.textContent = "";
      return;
    }

    dom.homeHeroMeta.innerHTML = `数据更新于 ${relativeTime(updatedAt)}`;
  }

  function renderDetailHeroMeta(dom, updatedAt) {
    if (!dom?.detailHeroMeta) return;
    if (!updatedAt) {
      dom.detailHeroMeta.textContent = "";
      return;
    }

    dom.detailHeroMeta.innerHTML = `数据更新于 ${relativeTime(updatedAt)}`;
  }

  function renderAppFooter() {
    const footer = document.querySelector("#appFooter");
    if (!footer) return;
    footer.innerHTML = [
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

  function renderSoftwareIcon(software) {
    const icon = String(software?.icon || "").trim();
    if (/^https?:\/\//i.test(icon) || icon.startsWith("/") || icon.startsWith("./")) {
      return `<img class="h-9 w-9 rounded-lg border border-slate-200 bg-white p-1 object-contain" src="${escapeAttr(icon)}" alt="${escapeAttr(software.name)} 图标" loading="lazy" />`;
    }

    const fallback = escapeHtml(icon || software.name.slice(0, 1).toUpperCase() || "?");
    return `<span class="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-sm font-semibold text-slate-600">${fallback}</span>`;
  }

  function renderDetailHeroIconBackground(software) {
    const icon = String(software?.icon || "").trim();
    if (/^https?:\/\//i.test(icon) || icon.startsWith("/") || icon.startsWith("./")) {
      return `
        <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-xl">
          <div class="absolute -bottom-6 -right-6 h-28 w-28 opacity-[0.08] blur-xl dark:opacity-[0.06]" style="background-image:url('${escapeAttr(icon)}');background-size:contain;background-repeat:no-repeat;background-position:center;"></div>
          <div class="absolute inset-0 bg-gradient-to-br from-white/0 via-white/0 to-white/35 dark:to-slate-800/25"></div>
        </div>`;
    }

    return "";
  }

  function normalizeLink(item) {
    if (!isObject(item)) return null;
    const type = String(item.type || "").trim().toLowerCase();
    const label = String(item.label || "").trim();
    const url = String(item.url || "").trim();
    if (!label || !url || !["direct", "webpage"].includes(type)) return null;
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
    return {
      updatedAt: String(payload?.updatedAt || "").trim(),
      versions: raw.map(normalizeVersion).filter(Boolean)
    };
  }

  function renderDetailEmpty(container, title, description) {
    if (!container) return;
    container.className = "min-h-[280px] grid place-items-center text-center";
    container.innerHTML = `<div><h2 class="text-xl font-semibold text-slate-700 dark:text-slate-200" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(title)}</h2><p class="mt-2 text-sm text-slate-500 dark:text-slate-400">${escapeHtml(description || "")}</p></div>`;
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
    const platform = String(variantPlatform || "").toLowerCase();
    if (!platform) return false;

    switch (currentPlatformId) {
      case "windows":
        return platform.includes("windows");
      case "macos":
        return platform.includes("mac");
      case "linux":
        return platform.includes("linux");
      case "android":
        return platform.includes("android");
      case "ios":
        return platform.includes("ios") || platform.includes("iphone") || platform.includes("ipad");
      case "web":
        return platform.includes("web");
      default:
        return false;
    }
  }

  function architectureScore(variantArchitecture, currentArchId) {
    const architecture = String(variantArchitecture || "").toLowerCase();
    const has = (keyword) => architecture.includes(keyword);

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

  function renderSoftwareDetail({ container, software, versions, onBack, onNavigateHome }) {
    if (!container) return;
    if (!software) {
      renderDetailEmpty(container, "请选择一个软件", "");
      return;
    }

    const currentPlatform = detectCurrentPlatform();
    const currentArchitecture = detectCurrentArchitecture();
    const detailIconMarkup = renderSoftwareIcon(software);
    const detailIconBackground = renderDetailHeroIconBackground(software);
    const tagsMarkup = (software.tags || [])
      .map(tag => `<span class="inline-block rounded-full bg-brand-50/80 px-2 py-0.5 text-xs font-medium text-brand-700 dark:bg-slate-700/50 dark:text-brand-400">#${escapeHtml(tag)}</span>`)
      .join(" ");

    container.className = "text-left";
    container.innerHTML = `
      <article class="relative overflow-hidden rounded-xl border border-slate-200/90 bg-white/92 p-4 shadow-[0_6px_16px_rgba(15,70,56,0.08)] dark:border-slate-700/80 dark:bg-slate-800/88 dark:shadow-[0_6px_16px_rgba(2,6,23,0.35)]">
        ${detailIconBackground}
        <div class="relative flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
          <div class="min-w-0 flex-1">
            <div class="mb-2 flex flex-wrap items-center gap-2">
              <button type="button" data-detail-back class="inline-flex items-center rounded-lg bg-brand-500 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-brand-700 focus:outline-none focus:ring-4 focus:ring-brand-500/20">
                ← 返回上一页
              </button>
              <button type="button" data-detail-home class="inline-flex items-center rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 transition hover:border-brand-500/40 hover:text-brand-700 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-300">
                回到软件目录
              </button>
            </div>
            <div class="mb-3 flex items-center gap-3">
              <div class="shrink-0">
                <span class="block [&>img]:h-12 [&>img]:w-12 [&>img]:rounded-lg [&>img]:border [&>img]:border-slate-200 [&>img]:bg-white [&>img]:p-1 [&>span]:inline-flex [&>span]:h-12 [&>span]:w-12 [&>span]:items-center [&>span]:justify-center [&>span]:rounded-lg [&>span]:border [&>span]:border-slate-200 [&>span]:bg-slate-50 dark:[&>img]:border-slate-700 dark:[&>img]:bg-slate-800 dark:[&>span]:border-slate-700 dark:[&>span]:bg-slate-800/90">${detailIconMarkup}</span>
              </div>
              <div class="min-w-0 flex flex-1 flex-wrap items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                <span class="rounded-full bg-white/80 px-2.5 py-1 shadow-sm dark:bg-slate-800/80">所属机构：${escapeHtml(software.organization)}</span>
                <span class="rounded-full bg-brand-50 px-2.5 py-1 font-medium text-brand-700 dark:bg-slate-700/60 dark:text-brand-300">当前检测环境：${escapeHtml(currentPlatform.label)} / ${escapeHtml(currentArchitecture.label)}</span>
              </div>
            </div>
            ${tagsMarkup ? `<div class="mt-3 flex flex-wrap gap-1.5">${tagsMarkup}</div>` : ""}
          </div>
          <div class="relative flex shrink-0 items-center">
            <a class="inline-flex w-fit items-center rounded-lg border border-brand-500/35 bg-brand-50 px-3 py-1.5 text-xs font-medium text-brand-700 hover:bg-brand-100 dark:border-brand-500/40 dark:bg-slate-700/50 dark:text-brand-300 dark:hover:bg-slate-700" target="_blank" rel="noopener noreferrer"
              href="${escapeAttr(software.officialWebsite)}">访问官网</a>
          </div>
        </div>
        <section class="relative mt-4 border-t border-slate-200/90 pt-3 dark:border-slate-700/80">
          <div class="mb-2 flex items-center justify-between gap-3">
            <div>
              <h3 class="text-sm font-semibold text-slate-800 dark:text-slate-100" style="font-family: 'Space Grotesk', sans-serif;">版本信息</h3>
              <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">以下版本与下载入口均归属于 ${escapeHtml(software.name)}</p>
            </div>
          </div>
          <div id="versionsContainer" class="grid gap-2.5"></div>
        </section>
      </article>
    `;

    container.querySelector("[data-detail-back]")?.addEventListener("click", () => {
      onBack?.();
    });
    container.querySelector("[data-detail-home]")?.addEventListener("click", () => {
      onNavigateHome?.();
    });

    const versionsContainer = container.querySelector("#versionsContainer");
    if (!versionsContainer) return;

    if (!Array.isArray(versions) || versions.length === 0) {
      versionsContainer.innerHTML =
        '<p class="rounded-lg border border-slate-200 bg-slate-50/90 px-3 py-4 text-sm text-slate-600 dark:border-slate-700 dark:bg-slate-800/70 dark:text-slate-400">暂无版本信息，请访问官网获取最新版本。</p>';
      return;
    }

    versions.forEach((versionItem) => {
      const card = document.createElement("div");
      card.className = "overflow-hidden rounded-lg border border-slate-200/90 bg-white/95 dark:border-slate-700/80 dark:bg-slate-800/90";

      const officialBtn = versionItem.officialUrl
          ? `<a class="inline-flex items-center rounded-md border border-amber-300 bg-amber-50 px-2 py-1 text-[11px] font-semibold text-amber-700 hover:bg-amber-100 dark:border-amber-700/60 dark:bg-amber-900/25 dark:text-amber-300 dark:hover:bg-amber-900/40" target="_blank" rel="noopener noreferrer"
              href="${escapeAttr(versionItem.officialUrl)}">前往官网发布页</a>`
        : "";

      const sortedVariants = [...(versionItem.variants || [])].sort((left, right) => {
        const leftPlatformScore = platformMatchesCurrent(left.platform, currentPlatform.id) ? 1 : 0;
        const rightPlatformScore = platformMatchesCurrent(right.platform, currentPlatform.id) ? 1 : 0;
        if (leftPlatformScore !== rightPlatformScore) return rightPlatformScore - leftPlatformScore;

        if (leftPlatformScore === 1) {
          const leftArchScore = architectureScore(left.architecture, currentArchitecture.id);
          const rightArchScore = architectureScore(right.architecture, currentArchitecture.id);
          if (leftArchScore !== rightArchScore) return rightArchScore - leftArchScore;
        }

        return 0;
      });

      const firstVariant = sortedVariants[0];
      const hasCurrentDeviceRow = !!firstVariant && platformMatchesCurrent(firstVariant.platform, currentPlatform.id);

      const variantRows = sortedVariants
        .map((variant, index) => {
          const isCurrentDevice = hasCurrentDeviceRow && index === 0;
          const directLinks = (variant.links || [])
            .map(
              (link) => {
                const tone = link.type === "webpage"
                  ? "border-slate-300 bg-white text-slate-700 hover:border-brand-500/40 hover:text-brand-700 dark:border-slate-600 dark:bg-slate-800/80 dark:text-slate-300 dark:hover:bg-slate-700"
                  : "border-brand-500/30 bg-brand-50 text-brand-700 hover:bg-brand-100 dark:border-brand-500/40 dark:bg-slate-700/50 dark:text-brand-300 dark:hover:bg-slate-700";
                const suffix = link.type === "webpage" ? "<span class=\"ml-1 text-[10px] font-medium opacity-70\">页面</span>" : "";
                return `<a class="inline-flex items-center rounded-md border px-2 py-1 text-[11px] font-semibold ${tone}" target="_blank" rel="noopener noreferrer"
                    href="${escapeAttr(link.url)}">${escapeHtml(link.label)}${suffix}</a>`;
              }
            )
            .join("");

          const directLinksHtml = directLinks
            ? `<div class="flex flex-wrap gap-1.5">${directLinks}</div>`
            : "暂无直链";

          return `
            <tr class="bg-white even:bg-slate-50 hover:bg-slate-100/70 dark:bg-slate-800 dark:even:bg-slate-800/75 dark:hover:bg-slate-700/60 ${isCurrentDevice ? "font-semibold" : ""}">
              <td class="whitespace-nowrap px-2.5 py-2 text-[13px] text-slate-700 dark:text-slate-200">${escapeHtml(variant.platform || "-")}${isCurrentDevice ? ' <span class="ml-1 text-[10px] font-semibold text-brand-700 dark:text-brand-300">当前设备</span>' : ""}</td>
              <td class="whitespace-nowrap px-2.5 py-2 text-[13px] text-slate-700 dark:text-slate-200">${escapeHtml(variant.architecture || "-")}</td>
              <td class="px-2.5 py-2 text-[13px] text-slate-700 dark:text-slate-200">${directLinksHtml}</td>
            </tr>`;
        })
        .join("");

      card.innerHTML = `
        <div class="flex flex-wrap items-center gap-2 border-b border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-700 dark:bg-slate-900/45">
          <span class="rounded-full bg-brand-50 px-2.5 py-1 text-[11px] font-medium text-brand-700 dark:bg-slate-700/50 dark:text-brand-300" style="font-family: 'Space Grotesk', sans-serif;">${escapeHtml(versionItem.version || "-")}</span>
          <span class="text-[11px] text-slate-500 dark:text-slate-400">${escapeHtml(versionItem.releaseDate || "")}</span>
          ${officialBtn}
        </div>
        ${
          variantRows
            ? `<div class="overflow-x-auto">
                <table class="min-w-full border-collapse">
                  <thead class="bg-slate-100 dark:bg-slate-900/55">
                    <tr><th class="px-2.5 py-2 text-left text-[11px] font-semibold tracking-wide text-slate-600 dark:text-slate-300">平台</th><th class="px-2.5 py-2 text-left text-[11px] font-semibold tracking-wide text-slate-600 dark:text-slate-300">架构</th><th class="px-2.5 py-2 text-left text-[11px] font-semibold tracking-wide text-slate-600 dark:text-slate-300">下载入口</th></tr>
                  </thead>
                  <tbody class="divide-y divide-slate-200 dark:divide-slate-700 dark:bg-slate-800">${variantRows}</tbody>
                </table>
              </div>`
            : '<p class="px-3 py-4 text-sm text-slate-600 dark:text-slate-400">该版本暂无构建信息。</p>'
        }
      `;

      versionsContainer.appendChild(card);
    });
  }

  // ── app/bootstrap ──────────────────────────────────────────────────────────
  function showOverlay(overlay, messageNode, message) {
    if (!overlay) return;
    if (messageNode && message) {
      messageNode.textContent = message;
    }
    overlay.style.display = "";
    overlay.classList.remove("hidden");
  }

  function hideOverlay(overlay) {
    if (!overlay) return;
    overlay.classList.add("hidden");
    setTimeout(() => { overlay.style.display = "none"; }, 400);
  }

  function getRouteSoftwareId() {
    return new URLSearchParams(window.location.search).get("id") || "";
  }

  function buildRouteUrl(softwareId) {
    const nextUrl = new URL("./index.html", window.location.href);
    if (softwareId) {
      nextUrl.searchParams.set("id", softwareId);
    }
    return nextUrl.toString();
  }

  function renderHomeLayout(dom) {
    dom.homeHero?.classList.remove("hidden");
    dom.detailHero?.classList.add("hidden");
    dom.homeSection?.classList.remove("hidden");
    dom.detailSection?.classList.add("hidden");
    document.title = "Original Software Hub";
  }

  function renderDetailLayout(dom, software, updatedAt) {
    dom.homeHero?.classList.add("hidden");
    dom.detailHero?.classList.remove("hidden");
    dom.homeSection?.classList.add("hidden");
    dom.detailSection?.classList.remove("hidden");
    const softwareName = String(software?.name || "").trim();
    const softwareDescription = String(software?.description || "").trim();
    if (dom.detailHeroTitle) {
      dom.detailHeroTitle.textContent = softwareName || "软件详情";
    }
    if (dom.detailHeroDescription) {
      dom.detailHeroDescription.textContent = softwareDescription || "查看软件介绍、版本发布时间和下载入口。";
    }
    renderDetailHeroMeta(dom, updatedAt);
    document.title = softwareName ? `${softwareName} - Original Software Hub` : "Original Software Hub";
  }

  async function bootstrapHomeApp() {
    initDarkMode();
    const dom = getAppDom();
    const state = {
      softwares: [],
      keyword: "",
      latestRenderToken: 0
    };

    bindAppEvents(dom, {
      onKeywordChange(keyword) {
        state.keyword = keyword.trim().toLowerCase();
        if (!getRouteSoftwareId()) {
          renderHomeList();
        }
      },
      onBack() {
        const referrer = document.referrer || "";
        const hasHistory = window.history.length > 1;
        const isSameOriginReferrer = referrer.startsWith(window.location.origin);

        if (hasHistory && isSameOriginReferrer) {
          window.history.back();
          return;
        }

        navigateToHome({ replace: true });
      },
      onNavigateHome() {
        navigateToHome();
      }
    });

    window.addEventListener("popstate", () => {
      renderCurrentRoute();
    });

    try {
      showOverlay(dom.loadingOverlay, dom.loadingMessage, "正在加载软件列表...");
      const catalog = await dataRepository.loadSoftwareCatalog();
      state.softwares = catalog.softwares;
      renderHomeHeroMeta(dom, catalog.generatedAt);
      renderAppFooter();
      await renderCurrentRoute();
    } catch (error) {
      hideOverlay(dom.loadingOverlay);
      renderHomeLayout(dom);
      showLoadError(dom, error);
    }

    function navigateToHome(options) {
      const replace = !!options?.replace;
      const targetUrl = buildRouteUrl("");
      const method = replace ? "replaceState" : "pushState";
      window.history[method]({}, "", targetUrl);
      renderCurrentRoute();
    }

    function navigateToSoftware(softwareId, options) {
      const replace = !!options?.replace;
      const targetUrl = buildRouteUrl(softwareId);
      const method = replace ? "replaceState" : "pushState";
      window.history[method]({ softwareId }, "", targetUrl);
      renderCurrentRoute();
    }

    function renderHomeList() {
      renderHomeLayout(dom);
      if (dom.searchInput && dom.searchInput.value !== state.keyword) {
        dom.searchInput.value = state.keyword;
      }
      renderSoftwareList({
        container: dom.list,
        softwares: state.softwares,
        keyword: state.keyword,
        onTagSelect(tag) {
          const cleanTag = String(tag || "").trim();
          if (!cleanTag) return;
          state.keyword = `#${cleanTag.toLowerCase()}`;
          if (dom.searchInput) dom.searchInput.value = `#${cleanTag}`;
          renderHomeList();
        },
        onSelect(softwareId) {
          navigateToSoftware(softwareId);
        }
      });
    }

    async function renderCurrentRoute() {
      const routeSoftwareId = getRouteSoftwareId().trim();
      const renderToken = ++state.latestRenderToken;

      if (!routeSoftwareId) {
        renderHomeList();
        hideOverlay(dom.loadingOverlay);
        return;
      }

      renderDetailLayout(dom, "详情", "");
      renderDetailEmpty(dom.detailContainer, "正在准备详情", "请稍候...");
      showOverlay(dom.loadingOverlay, dom.loadingMessage, "正在加载详情...");

      try {
        const software = await dataRepository.getSoftwareById(routeSoftwareId);
        if (renderToken !== state.latestRenderToken) return;

        if (!software) {
          renderDetailLayout(dom, { name: "未找到软件", description: "请返回目录页选择有效的软件。" }, "");
          renderDetailEmpty(dom.detailContainer, "未找到软件", "请返回目录页选择有效的软件。");
          hideOverlay(dom.loadingOverlay);
          return;
        }

        renderDetailLayout(dom, software, "");
        const rawVersions = await dataRepository.loadSoftwareVersions(software);
        if (renderToken !== state.latestRenderToken) return;

        const { versions, updatedAt } = normalizeSoftwareVersionPayload(rawVersions);
        renderDetailLayout(dom, software, updatedAt);
        renderSoftwareDetail({
          container: dom.detailContainer,
          software,
          versions,
          onBack() {
            const referrer = document.referrer || "";
            const hasHistory = window.history.length > 1;
            const isSameOriginReferrer = referrer.startsWith(window.location.origin);

            if (hasHistory && isSameOriginReferrer) {
              window.history.back();
              return;
            }

            navigateToHome({ replace: true });
          },
          onNavigateHome() {
            navigateToHome();
          }
        });
        hideOverlay(dom.loadingOverlay);
      } catch (error) {
        if (renderToken !== state.latestRenderToken) return;
        const message = error instanceof Error ? error.message : "未知错误";
        renderDetailLayout(dom, { name: "加载失败", description: `详情加载失败：${message}` }, "");
        renderDetailEmpty(dom.detailContainer, "加载失败", `详情加载失败：${message}`);
        hideOverlay(dom.loadingOverlay);
      }
    }
  }

  // ── entry ──────────────────────────────────────────────────────────────────
  bootstrapHomeApp();
})();

// 共享数据入口：统一配置解析、请求加载与内存级 Promise 缓存。
(function () {
  "use strict";

  const APP_DATA_SOURCE_CONFIG = window.APP_DATA_SOURCE_CONFIG;
  if (!APP_DATA_SOURCE_CONFIG) {
    throw new Error("缺少 APP_DATA_SOURCE_CONFIG，请先加载 config.js");
  }

  if (window.OSH_DATA_REPOSITORY) {
    return;
  }

  const DEFAULT_JSONP_TIMEOUT_MS = 8000;
  const DEFAULT_JSONP_CALLBACK_PARAM = "callback";

  function isObject(value) {
    return value !== null && typeof value === "object";
  }

  function toPositiveNumber(value, fallback) {
    const num = Number(value);
    return num > 0 ? num : fallback;
  }

  async function fetchJson(url) {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${url}`);
    }
    return response.json();
  }

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

  function fetchBySource(source) {
    if (source.mode === "jsonp") {
      return fetchJsonp(source);
    }
    return fetchJson(source.url);
  }

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

  function resolveUrl(baseUrl, path) {
    const base = String(baseUrl).endsWith("/") ? baseUrl : `${baseUrl}/`;
    return new URL(path, new URL(base, window.location.href)).toString();
  }

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
      pinyin: String(item.pinyin || "").trim().toLowerCase(),
      icon: String(item.icon || "").trim(),
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
      const nextMode = String(source.mode || source.type || mode).trim().toLowerCase();
      mode = nextMode || mode;
      callbackParam = String(source.callbackParam || callbackParam).trim() || callbackParam;
      timeoutMs = toPositiveNumber(source.timeoutMs, timeoutMs);
    }

    if (!sourcePath) throw new Error("软件缺少 source 配置");
    if (mode !== "json" && mode !== "jsonp") throw new Error("source.mode 仅支持 json 或 jsonp");

    return { mode, url: resolveUrl(defaults.baseUrl, sourcePath), callbackParam, timeoutMs };
  }

  function createCachedPromiseStore() {
    const store = new Map();
    return {
      get(key, loader) {
        if (store.has(key)) {
          return store.get(key);
        }

        const promise = Promise.resolve()
          .then(loader)
          .catch((error) => {
            store.delete(key);
            throw error;
          });

        store.set(key, promise);
        return promise;
      }
    };
  }

  function createDataRepository(config) {
    const requestCache = createCachedPromiseStore();
    const versionCache = createCachedPromiseStore();
    let dataSourcePromise = null;
    let softwareCatalogPromise = null;
    let softwareMapPromise = null;

    function loadRequest(source) {
      const cacheKey = [
        source.mode,
        source.url,
        source.callbackParam || "",
        String(source.timeoutMs || "")
      ].join("::");

      return requestCache.get(cacheKey, () => fetchBySource(source));
    }

    function loadDataSourceConfig() {
      if (dataSourcePromise) {
        return dataSourcePromise;
      }

      dataSourcePromise = Promise.resolve()
        .then(() => {
          const normalized = normalizeDataSourceConfig(config);
          const { type, url, indexPath, callbackParam, timeoutMs } = normalized.endpoint;

          const indexSource = {
            mode: type,
            url: resolveUrl(url, indexPath),
            callbackParam,
            timeoutMs
          };

          return loadRequest(indexSource).then((indexPayload) => {
            const softwareListEntry = indexPayload?.softwareList;
            let listPath;
            let listMode;

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
          });
        })
        .catch((error) => {
          dataSourcePromise = null;
          throw error;
        });

      return dataSourcePromise;
    }

    function loadSoftwareCatalog() {
      if (softwareCatalogPromise) {
        return softwareCatalogPromise;
      }

      softwareCatalogPromise = loadDataSourceConfig()
        .then(async (dataSource) => {
          const rawList = await loadRequest(dataSource.softwareList);
          return {
            generatedAt: dataSource.generatedAt,
            softwareSourceDefaults: dataSource.softwareSourceDefaults,
            softwares: normalizeSoftwareListPayload(rawList).items
          };
        })
        .catch((error) => {
          softwareCatalogPromise = null;
          throw error;
        });

      return softwareCatalogPromise;
    }

    function loadSoftwareMap() {
      if (softwareMapPromise) {
        return softwareMapPromise;
      }

      softwareMapPromise = loadSoftwareCatalog()
        .then((catalog) => {
          return new Map(catalog.softwares.map((software) => [software.id, software]));
        })
        .catch((error) => {
          softwareMapPromise = null;
          throw error;
        });

      return softwareMapPromise;
    }

    function getSoftwareById(softwareId) {
      const normalizedId = String(softwareId || "").trim();
      if (!normalizedId) {
        return Promise.resolve(null);
      }

      return loadSoftwareMap().then((softwareMap) => softwareMap.get(normalizedId) || null);
    }

    function loadSoftwareVersions(software) {
      const softwareId = String(software?.id || "").trim();
      if (!softwareId) {
        return Promise.reject(new Error("缺少软件 ID"));
      }

      return versionCache.get(softwareId, async () => {
        const catalog = await loadSoftwareCatalog();
        const source = resolveSoftwareVersionSource(software.source, catalog.softwareSourceDefaults);
        return loadRequest(source);
      });
    }

    return {
      loadDataSourceConfig,
      loadSoftwareCatalog,
      getSoftwareById,
      loadSoftwareVersions
    };
  }

  window.OSH_DATA_REPOSITORY = createDataRepository(APP_DATA_SOURCE_CONFIG);
})();
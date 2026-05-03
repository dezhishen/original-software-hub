// 数据获取服务 — bundle-shared.js 的 ES 模块版本（单例）

const DEFAULT_TIMEOUT_MS = 8000
const DEFAULT_CALLBACK_PARAM = 'callback'

function isObject(v) {
  return v !== null && typeof v === 'object'
}

function toPositiveNumber(v, fallback) {
  const n = Number(v)
  return n > 0 ? n : fallback
}

async function fetchJson(url) {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${url}`)
  return res.json()
}

function fetchJsonp(source) {
  const { url, callbackParam = DEFAULT_CALLBACK_PARAM, timeoutMs = DEFAULT_TIMEOUT_MS } = source
  const callbackName = `__jsonp_cb_${Date.now()}_${Math.random().toString(36).slice(2)}`
  return new Promise((resolve, reject) => {
    const script = document.createElement('script')
    let settled = false
    const prevFixed = window.callback
    const cleanup = () => {
      settled = true
      delete window[callbackName]
      if (prevFixed === undefined) delete window.callback
      else window.callback = prevFixed
      script.remove()
    }
    const timer = setTimeout(() => {
      if (!settled) { cleanup(); reject(new Error(`JSONP timeout: ${url}`)) }
    }, timeoutMs)
    window[callbackName] = window.callback = (data) => {
      clearTimeout(timer); cleanup(); resolve(data)
    }
    script.onerror = () => { clearTimeout(timer); cleanup(); reject(new Error(`JSONP load error: ${url}`)) }
    script.src = `${url}${url.includes('?') ? '&' : '?'}${encodeURIComponent(callbackParam)}=${callbackName}`
    document.head.appendChild(script)
  })
}

function fetchBySource(source) {
  return source.mode === 'jsonp' ? fetchJsonp(source) : fetchJson(source.url)
}

function normalizeDataSourceConfig(config) {
  if (!isObject(config)) throw new Error('数据源配置格式错误')
  const ep = config.endpoint
  if (!isObject(ep)) throw new Error('数据源配置缺少 endpoint')
  const type = String(ep.type || 'json').trim().toLowerCase()
  const url = String(ep.url || '').trim()
  const indexPath = String(ep.indexPath || 'index.json').trim() || 'index.json'
  const callbackParam = String(ep.callbackParam || DEFAULT_CALLBACK_PARAM).trim() || DEFAULT_CALLBACK_PARAM
  const timeoutMs = toPositiveNumber(ep.timeoutMs, DEFAULT_TIMEOUT_MS)
  if (!url) throw new Error('数据源配置缺少 endpoint.url')
  if (type !== 'json' && type !== 'jsonp') throw new Error('endpoint.type 仅支持 json 或 jsonp')
  return { endpoint: { type, url, indexPath, callbackParam, timeoutMs } }
}

function resolveUrl(baseUrl, path) {
  const base = String(baseUrl).endsWith('/') ? baseUrl : `${baseUrl}/`
  return new URL(path, new URL(base, window.location.href)).toString()
}

function normalizeSoftwareSource(source) {
  if (typeof source === 'string') return source.trim() || null
  if (!isObject(source)) return null
  const path = String(source.path || source.url || '').trim()
  if (!path) return null
  const mode = String(source.mode || source.type || '').trim().toLowerCase()
  if (mode && mode !== 'json' && mode !== 'jsonp') return null
  const result = { path }
  if (mode) result.mode = mode
  const callbackParam = String(source.callbackParam || '').trim()
  if (callbackParam) result.callbackParam = callbackParam
  const timeoutMs = Number(source.timeoutMs)
  if (timeoutMs > 0) result.timeoutMs = timeoutMs
  return result
}

function normalizeSoftwareItem(item) {
  if (!isObject(item)) return null
  const id = String(item.id || '').trim()
  const name = String(item.name || '').trim()
  const organization = String(item.organization || '').trim()
  const officialWebsite = String(item.officialWebsite || '').trim()
  const source = normalizeSoftwareSource(item.source)
  if (!id || !name || !organization || !officialWebsite || !source) return null
  return {
    id,
    name,
    pinyin: String(item.pinyin || '').trim().toLowerCase(),
    icon: String(item.icon || '').trim(),
    description: String(item.description || '').trim(),
    organization,
    officialWebsite,
    tags: Array.isArray(item.tags)
      ? item.tags.map((t) => String(t || '').trim()).filter(Boolean).slice(0, 3)
      : [],
    source
  }
}

function normalizeSoftwareListPayload(payload) {
  const items = Array.isArray(payload?.items) ? payload.items : []
  return { items: items.map(normalizeSoftwareItem).filter(Boolean) }
}

function resolveSoftwareVersionSource(source, defaults) {
  if (!isObject(defaults)) throw new Error('版本源默认配置无效')
  let mode = String(defaults.mode || 'json').trim().toLowerCase() || 'json'
  let sourcePath = ''
  let callbackParam = String(defaults.callbackParam || DEFAULT_CALLBACK_PARAM).trim() || DEFAULT_CALLBACK_PARAM
  let timeoutMs = toPositiveNumber(defaults.timeoutMs, DEFAULT_TIMEOUT_MS)
  if (typeof source === 'string') {
    sourcePath = source.trim()
  } else if (isObject(source)) {
    sourcePath = String(source.path || source.url || '').trim()
    const nextMode = String(source.mode || source.type || mode).trim().toLowerCase()
    mode = nextMode || mode
    callbackParam = String(source.callbackParam || callbackParam).trim() || callbackParam
    timeoutMs = toPositiveNumber(source.timeoutMs, timeoutMs)
  }
  if (!sourcePath) throw new Error('软件缺少 source 配置')
  if (mode !== 'json' && mode !== 'jsonp') throw new Error('source.mode 仅支持 json 或 jsonp')
  return { mode, url: resolveUrl(defaults.baseUrl, sourcePath), callbackParam, timeoutMs }
}

function createCachedStore() {
  const store = new Map()
  return {
    get(key, loader) {
      if (store.has(key)) return store.get(key)
      const p = Promise.resolve().then(loader).catch((e) => { store.delete(key); throw e })
      store.set(key, p)
      return p
    }
  }
}

function createDataRepository(config) {
  const requestCache = createCachedStore()
  const versionCache = createCachedStore()
  let dataSourcePromise = null
  let softwareCatalogPromise = null
  let softwareMapPromise = null

  function loadRequest(source) {
    const key = [source.mode, source.url, source.callbackParam || '', String(source.timeoutMs || '')].join('::')
    return requestCache.get(key, () => fetchBySource(source))
  }

  function loadDataSourceConfig() {
    if (dataSourcePromise) return dataSourcePromise
    dataSourcePromise = Promise.resolve()
      .then(() => {
        const { type, url, indexPath, callbackParam, timeoutMs } = normalizeDataSourceConfig(config).endpoint
        return loadRequest({ mode: type, url: resolveUrl(url, indexPath), callbackParam, timeoutMs })
          .then((indexPayload) => {
            const listEntry = indexPayload?.softwareList
            let listPath, listMode
            if (isObject(listEntry)) {
              listPath = String(listEntry.path || '').trim()
              listMode = String(listEntry.mode || type).trim().toLowerCase() || type
            } else {
              listPath = String(listEntry || '').trim()
              listMode = type
            }
            if (!listPath) throw new Error('index 缺少 softwareList.path')
            return {
              generatedAt: String(indexPayload?.meta?.generatedAt || '').trim(),
              softwareList: { mode: listMode, url: resolveUrl(url, listPath), callbackParam, timeoutMs },
              softwareSourceDefaults: { mode: type, baseUrl: url, callbackParam, timeoutMs }
            }
          })
      })
      .catch((e) => { dataSourcePromise = null; throw e })
    return dataSourcePromise
  }

  function loadSoftwareCatalog() {
    if (softwareCatalogPromise) return softwareCatalogPromise
    softwareCatalogPromise = loadDataSourceConfig()
      .then(async (ds) => {
        const rawList = await loadRequest(ds.softwareList)
        return {
          generatedAt: ds.generatedAt,
          softwareSourceDefaults: ds.softwareSourceDefaults,
          softwares: normalizeSoftwareListPayload(rawList).items
        }
      })
      .catch((e) => { softwareCatalogPromise = null; throw e })
    return softwareCatalogPromise
  }

  function loadSoftwareMap() {
    if (softwareMapPromise) return softwareMapPromise
    softwareMapPromise = loadSoftwareCatalog()
      .then((catalog) => new Map(catalog.softwares.map((s) => [s.id, s])))
      .catch((e) => { softwareMapPromise = null; throw e })
    return softwareMapPromise
  }

  function getSoftwareById(softwareId) {
    const id = String(softwareId || '').trim()
    if (!id) return Promise.resolve(null)
    return loadSoftwareMap().then((m) => m.get(id) || null)
  }

  function loadSoftwareVersions(software) {
    const id = String(software?.id || '').trim()
    if (!id) return Promise.reject(new Error('缺少软件 ID'))
    return versionCache.get(id, async () => {
      const catalog = await loadSoftwareCatalog()
      const source = resolveSoftwareVersionSource(software.source, catalog.softwareSourceDefaults)
      return loadRequest(source)
    })
  }

  return { loadSoftwareCatalog, getSoftwareById, loadSoftwareVersions }
}

// 读取运行时 config（由 public/config.js 注入）或使用默认值
const runtimeConfig = window.APP_DATA_SOURCE_CONFIG ?? {
  endpoint: { type: 'json', url: './data/json', indexPath: 'index.json', timeoutMs: 8000 }
}

export const dataRepository = createDataRepository(runtimeConfig)

function isObject(v) {
  return v !== null && typeof v === 'object'
}

function normalizeLink(item) {
  if (!isObject(item)) return null
  const type = String(item.type || '').trim().toLowerCase()
  const label = String(item.label || '').trim()
  const url = String(item.url || '').trim()
  if (!label || !url || !['direct', 'webpage', 'store'].includes(type)) return null
  return { type, label, url }
}

function normalizePlatformPackage(item) {
  if (!isObject(item)) return null
  return {
    architecture: String(item.architecture || '').trim(),
    links: Array.isArray(item.links) ? item.links.map(normalizeLink).filter(Boolean) : []
  }
}

function normalizePlatformRelease(item, fallback) {
  if (!isObject(item)) return null
  const platform = String(item.platform || '').trim()
  if (!platform) return null
  return {
    platform,
    version: String(item.version || fallback?.version || '').trim(),
    releaseDate: String(item.releaseDate || fallback?.releaseDate || '').trim(),
    officialUrl: String(item.officialUrl || fallback?.officialUrl || '').trim(),
    packages: Array.isArray(item.packages)
      ? item.packages.map(normalizePlatformPackage).filter(Boolean)
      : []
  }
}

function normalizeVariant(item) {
  if (!isObject(item)) return null
  return {
    architecture: String(item.architecture || '').trim(),
    platform: String(item.platform || '').trim(),
    links: Array.isArray(item.links) ? item.links.map(normalizeLink).filter(Boolean) : []
  }
}

function buildPlatformReleasesFromVariants(variants, fallback) {
  const groups = new Map()
  ;(variants || []).forEach((variant) => {
    const platformLabel = String(variant?.platform || '').trim() || '其他'
    if (!groups.has(platformLabel)) {
      groups.set(platformLabel, {
        platform: platformLabel,
        version: String(fallback?.version || '').trim(),
        releaseDate: String(fallback?.releaseDate || '').trim(),
        officialUrl: String(fallback?.officialUrl || '').trim(),
        packages: []
      })
    }
    groups.get(platformLabel).packages.push({
      architecture: String(variant?.architecture || '').trim(),
      links: Array.isArray(variant?.links) ? variant.links : []
    })
  })
  return [...groups.values()]
}

function normalizeVersion(item) {
  if (!isObject(item)) return null
  const version = String(item.version || '').trim()
  const releaseDate = String(item.releaseDate || '').trim()
  const officialUrl = String(item.officialUrl || '').trim()
  const variants = Array.isArray(item.variants)
    ? item.variants.map(normalizeVariant).filter(Boolean)
    : []
  const platforms = Array.isArray(item.platforms)
    ? item.platforms
        .map((e) => normalizePlatformRelease(e, { version, releaseDate, officialUrl }))
        .filter(Boolean)
    : []
  return {
    version,
    releaseDate,
    officialUrl,
    variants,
    platforms: platforms.length
      ? platforms
      : buildPlatformReleasesFromVariants(variants, { version, releaseDate, officialUrl })
  }
}

function mergePlatformReleases(entries) {
  const groups = new Map()
  const order = []
  ;(entries || []).forEach((entry) => {
    const platform = String(entry?.platform || '').trim() || '其他'
    if (!groups.has(platform)) {
      groups.set(platform, { platform, version: '', releaseDate: '', officialUrl: '', packages: [] })
      order.push(platform)
    }
    const current = groups.get(platform)
    const nextDate = String(entry?.releaseDate || '').trim()
    const currentDate = String(current.releaseDate || '').trim()
    const shouldReplaceMeta = !currentDate || (nextDate && nextDate >= currentDate)
    if (shouldReplaceMeta) {
      current.version = String(entry?.version || current.version || '').trim()
      current.releaseDate = nextDate || currentDate
      current.officialUrl = String(entry?.officialUrl || current.officialUrl || '').trim()
    }
    const seen = new Set(
      (current.packages || []).flatMap((pkg) =>
        (pkg.links || []).map((link) => `${String(pkg.architecture || '').trim()}::${String(link.url || '').trim()}`)
      )
    )
    ;(entry?.packages || []).forEach((pkg) => {
      const arch = String(pkg?.architecture || '').trim()
      const links = []
      ;(pkg?.links || []).forEach((link) => {
        const url = String(link?.url || '').trim()
        if (!url) return
        const key = `${arch}::${url}`
        if (seen.has(key)) return
        seen.add(key)
        links.push(link)
      })
      if (links.length) current.packages.push({ architecture: pkg?.architecture || '', links })
    })
  })
  return order.map((p) => groups.get(p)).filter(Boolean)
}

function flattenVersionEntriesToPlatforms(versionEntries) {
  const all = []
  ;(versionEntries || []).forEach((v) => {
    ;(v?.platforms || []).forEach((p) => all.push(p))
  })
  return mergePlatformReleases(all)
}

/**
 * 将原始版本 payload 规范化为 { updatedAt, platforms[] }
 * @param {unknown} payload
 * @returns {{ updatedAt: string, platforms: PlatformRelease[] }}
 */
export function normalizeSoftwareVersionPayload(payload) {
  const directPlatforms = Array.isArray(payload?.platforms)
    ? payload.platforms.map((e) => normalizePlatformRelease(e, {})).filter(Boolean)
    : []
  const rawVersions = Array.isArray(payload?.versions)
    ? payload.versions
    : Array.isArray(payload)
      ? payload
      : []
  const normalizedVersions = rawVersions.map(normalizeVersion).filter(Boolean)
  return {
    updatedAt: String(payload?.updatedAt || '').trim(),
    platforms: directPlatforms.length
      ? mergePlatformReleases(directPlatforms)
      : flattenVersionEntriesToPlatforms(normalizedVersions)
  }
}

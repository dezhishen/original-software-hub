function asObject(value) {
  return value && typeof value === 'object' ? value : null
}

function asArray(value) {
  return Array.isArray(value) ? value : []
}

function asPositiveNumber(value, fallback) {
  const n = Number(value)
  return Number.isFinite(n) && n > 0 ? n : fallback
}

function escapeHtml(value) {
  return String(value || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function mapScriptItem(item) {
  const source = asObject(item)
  if (!source) return null
  const src = String(source.src || '').trim()
  const inline = String(source.inline || '').trim()
  if (!src && !inline) return null
  return {
    src: src || undefined,
    inline: inline || undefined,
    type: String(source.type || '').trim() || undefined,
    async: source.async,
    defer: source.defer,
    attrs: asObject(source.attrs) || undefined
  }
}

function platformToEmbed(provider, slotName) {
  const source = asObject(provider)
  if (!source) return null
  const platform = String(source.platform || source.type || '').trim().toLowerCase()
  const id = String(source.id || `${slotName}-${Math.random().toString(36).slice(2)}`).trim()
  const title = String(source.title || '').trim()

  if (platform === 'iframe') {
    const src = String(source.src || '').trim()
    const srcdoc = String(source.srcdoc || '').trim()
    if (!src && !srcdoc) return null
    return { id, type: 'iframe', src: src || undefined, srcdoc: srcdoc || undefined, title }
  }

  if (platform === 'custom-html' || platform === 'html') {
    const html = String(source.html || '').trim()
    if (!html && !title) return null
    return { id, type: 'html', html: html || undefined, title }
  }

  if (platform === 'adsense') {
    const client = String(source.client || '').trim()
    const slot = String(source.slot || '').trim()
    const fallbackHtml = String(source.fallbackHtml || '').trim()
    if ((!client || !slot) && (fallbackHtml || title)) {
      return {
        id,
        type: 'html',
        html: fallbackHtml || undefined,
        title: title || '公益合作位'
      }
    }
    if (!client || !slot) return null
    const format = String(source.format || 'auto').trim() || 'auto'
    const responsive = source.responsive === undefined ? 'true' : String(Boolean(source.responsive))
    return {
      id,
      type: 'script',
      title: title || 'Google AdSense',
      slotHtml:
        `<ins class="adsbygoogle" style="display:block" data-ad-client="${escapeHtml(client)}" data-ad-slot="${escapeHtml(slot)}" data-ad-format="${escapeHtml(format)}" data-full-width-responsive="${escapeHtml(responsive)}"></ins>`,
      scripts: [
        {
          src: 'https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js',
          async: true,
          attrs: { crossorigin: 'anonymous' }
        },
        {
          inline: '(adsbygoogle = window.adsbygoogle || []).push({});'
        }
      ]
    }
  }

  if (platform === 'script' || platform === 'script-snippet') {
    const slotHtml = String(source.slotHtml || '').trim()
    const scripts = asArray(source.scripts).map(mapScriptItem).filter(Boolean)
    if (!slotHtml && !scripts.length && !title) return null
    return {
      id,
      type: 'script',
      title,
      slotHtml: slotHtml || undefined,
      scripts
    }
  }

  return null
}

function normalizeSlot(slotConfig, slotName, fallbackRotateMs) {
  const slot = asObject(slotConfig) || {}
  const rotateMs = asPositiveNumber(slot.rotateMs, fallbackRotateMs)
  const providers = asArray(slot.providers)
  const embeds = providers.map((provider) => platformToEmbed(provider, slotName)).filter(Boolean)
  return { rotateMs, embeds }
}

function createRuntimeAdConfig(rawConfig) {
  const config = asObject(rawConfig) || {}
  const slots = asObject(config.slots) || {}
  const defaultRotateMs = asPositiveNumber(config.defaultRotateMs, 6000)

  return {
    getSlot(slotName) {
      const key = String(slotName || '').trim()
      const slot = slots[key]
      return normalizeSlot(slot, key || 'slot', defaultRotateMs)
    }
  }
}

const runtimeConfig = createRuntimeAdConfig(window.APP_AD_CONFIG)

export function getAdSlotConfig(slotName) {
  return runtimeConfig.getSlot(slotName)
}

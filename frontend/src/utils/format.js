/**
 * 将日期字符串解析为 Date + 是否仅日期 的元组
 * @param {string} dateStr
 * @returns {{ date: Date, isDateOnly: boolean } | null}
 */
export function parseDate(dateStr) {
  if (!dateStr) return null
  const raw = String(dateStr).trim()
  const isDateOnly = /^\d{4}-\d{2}-\d{2}$/.test(raw)
  const date = new Date(isDateOnly ? `${raw}T12:00:00+08:00` : raw)
  return isNaN(date.getTime()) ? null : { date, isDateOnly }
}

/**
 * 返回相对时间文本（如"3 天前"）
 * @param {string} dateStr
 * @returns {string}
 */
export function relativeTimeText(dateStr) {
  const parsed = parseDate(dateStr)
  if (!parsed) return String(dateStr || '')
  const { date, isDateOnly } = parsed
  const diffMs = Date.now() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  const diffHr = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHr / 24)
  if (diffMs < 0) return '刚刚'
  if (isDateOnly && diffDay < 1) return '今天'
  if (diffMin < 1) return '刚刚'
  if (diffMin < 60) return `${diffMin} 分钟前`
  if (diffHr < 24) return `${diffHr} 小时前`
  if (diffDay < 30) return `${diffDay} 天前`
  if (diffDay < 365) return `${Math.floor(diffDay / 30)} 个月前`
  return `${Math.floor(diffDay / 365)} 年前`
}

/**
 * 返回本地时间文本（北京时间）
 * @param {string} dateStr
 * @returns {string}
 */
export function absoluteTimeText(dateStr) {
  const parsed = parseDate(dateStr)
  if (!parsed) return ''
  return parsed.date.toLocaleString('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * 返回 ISO 8601 字符串
 * @param {string} dateStr
 * @returns {string}
 */
export function isoTimeText(dateStr) {
  const parsed = parseDate(dateStr)
  return parsed?.date.toISOString() ?? ''
}

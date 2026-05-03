/**
 * useDetailBackground
 *
 * 管理详情页的 body 背景特效：
 *   - 将软件图标 URL 注入 CSS 变量 --detail-icon-url
 *   - 添加/移除 detail-icon-bg / detail-icon-bg-enter 类
 *   - 组件卸载时自动清理
 */
import { onUnmounted } from 'vue'

function isIconUrl(icon) {
  return icon && (/^https?:\/\//i.test(icon) || icon.startsWith('/') || icon.startsWith('./'))
}

function escapeCssUrl(src) {
  return src.replace(/\\/g, '\\\\').replace(/"/g, '\\"')
}

export function useDetailBackground() {
  let enterTimer = 0

  function applyTo(target, icon) {
    target.style.setProperty('--detail-icon-url', `url("${escapeCssUrl(icon)}")`)
    target.classList.add('detail-icon-bg')
    target.classList.remove('detail-icon-bg-enter')
  }

  function clearFrom(target) {
    target.classList.remove('detail-icon-bg', 'detail-icon-bg-enter')
    target.style.removeProperty('--detail-icon-url')
  }

  function apply(software) {
    const icon = String(software?.icon || '').trim()
    if (!isIconUrl(icon)) {
      clear()
      return
    }

    const body = document.body
    const root = document.documentElement
    applyTo(body, icon)
    applyTo(root, icon)

    // 触发入场动画
    if (enterTimer) clearTimeout(enterTimer)
    body.classList.remove('detail-icon-bg-enter')
    root.classList.remove('detail-icon-bg-enter')
    void body.offsetWidth  // force reflow
    body.classList.add('detail-icon-bg-enter')
    root.classList.add('detail-icon-bg-enter')
    enterTimer = window.setTimeout(() => {
      body.classList.remove('detail-icon-bg-enter')
      root.classList.remove('detail-icon-bg-enter')
      enterTimer = 0
    }, 620)
  }

  function clear() {
    if (enterTimer) { clearTimeout(enterTimer); enterTimer = 0 }
    clearFrom(document.body)
    clearFrom(document.documentElement)
  }

  onUnmounted(clear)

  return { applyBackground: apply, clearBackground: clear }
}

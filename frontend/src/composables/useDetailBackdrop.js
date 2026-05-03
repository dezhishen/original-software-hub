import { computed } from 'vue'

function isIconUrl(icon) {
  return icon && (/^https?:\/\//i.test(icon) || icon.startsWith('/') || icon.startsWith('./'))
}

function toAbsoluteUrl(icon) {
  try {
    return new URL(icon, window.location.href).toString()
  } catch {
    return icon
  }
}

function escapeCssUrl(src) {
  return src.replace(/\\/g, '\\\\').replace(/"/g, '\\"')
}

export function useDetailBackdrop(pageState, route) {
  const iconUrl = computed(() => String(pageState.detailSoftware?.icon || '').trim())
  const isDetailRoute = computed(() => route.path.startsWith('/software/'))

  const showDetailBackdrop = computed(() => {
    return isDetailRoute.value && isIconUrl(iconUrl.value)
  })

  const detailBackdropStyle = computed(() => {
    return {
      '--detail-icon-url': `url("${escapeCssUrl(toAbsoluteUrl(iconUrl.value))}")`
    }
  })

  return { showDetailBackdrop, detailBackdropStyle }
}

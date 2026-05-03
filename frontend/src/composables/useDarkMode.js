import { ref, watch } from 'vue'

// 模块级单例，所有组件共享同一状态
const isDark = ref(
  (() => {
    try {
      const saved = localStorage.getItem('osh-dark-mode')
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      return saved !== null ? saved === '1' : prefersDark
    } catch {
      return false
    }
  })()
)

watch(
  isDark,
  (val) => {
    document.documentElement.classList.toggle('dark', val)
    try { localStorage.setItem('osh-dark-mode', val ? '1' : '0') } catch {}
  },
  { immediate: true }
)

export function useDarkMode() {
  return {
    isDark,
    toggle: () => { isDark.value = !isDark.value }
  }
}

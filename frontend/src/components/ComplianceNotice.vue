<template>
  <section class="relative w-full rounded-xl border border-slate-200/80 bg-white/78 px-3 py-2.5 text-[11px] text-slate-500 shadow-[0_4px_14px_rgba(15,70,56,0.08)] backdrop-blur-sm dark:border-slate-700/85 dark:bg-slate-900/88 dark:text-slate-300 md:ml-auto md:max-w-[340px]">
    <p class="leading-5">下载入口来自官方渠道或官方镜像，仅作导航参考。</p>

    <button
      ref="triggerRef"
      type="button"
      class="mt-1.5 inline-flex items-center gap-1 text-[11px] font-medium text-brand-700 transition hover:text-brand-800 dark:text-brand-400 dark:hover:text-brand-300"
      @mouseenter="openPopover"
      @mouseleave="scheduleClose"
      @focus="openPopover"
      @blur="scheduleClose"
      @click.prevent
    >
      合规与隐私说明
      <span class="text-[10px]">ⓘ</span>
    </button>
  </section>

  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed z-[2147483647] w-[min(92vw,390px)] rounded-xl border border-slate-200/90 bg-white p-3 text-[11px] shadow-[0_16px_36px_rgba(2,6,23,0.28)] dark:border-slate-700/90 dark:bg-slate-900"
      :style="popoverStyle"
      @mouseenter="cancelClose"
      @mouseleave="scheduleClose"
    >
      <p class="mb-1.5 text-[10px] font-medium tracking-[0.04em] text-slate-400 dark:text-slate-500">合规与隐私说明</p>
      <div class="space-y-1.5 leading-5 text-slate-500 dark:text-slate-400">
        <p>本站不托管安装包，不替代官方发布页；链接可用性、文件安全性和版本准确性请自行核验。</p>
        <p>遵循数据最小化实践（参考 GDPR 原则），仅请求软件目录与版本数据，不要求登录，不采集身份信息。</p>
        <p>仅在浏览器本地保存主题偏好（localStorage）；可随时清除站点数据。</p>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'

const triggerRef = ref(null)
const visible = ref(false)
const top = ref(0)
const left = ref(0)
let closeTimer = 0

const popoverStyle = computed(() => ({
  top: `${top.value}px`,
  left: `${left.value}px`
}))

function recalcPosition() {
  const trigger = triggerRef.value
  if (!trigger) return
  const rect = trigger.getBoundingClientRect()
  const panelWidth = Math.min(window.innerWidth * 0.92, 390)
  const gap = 8
  top.value = rect.bottom + gap
  left.value = Math.min(Math.max(8, rect.right - panelWidth), window.innerWidth - panelWidth - 8)
}

function cancelClose() {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = 0
  }
}

function openPopover() {
  cancelClose()
  recalcPosition()
  visible.value = true
}

function scheduleClose() {
  cancelClose()
  closeTimer = window.setTimeout(() => {
    visible.value = false
    closeTimer = 0
  }, 120)
}

function onViewportChange() {
  if (visible.value) recalcPosition()
}

onMounted(() => {
  window.addEventListener('scroll', onViewportChange, true)
  window.addEventListener('resize', onViewportChange)
})

onUnmounted(() => {
  if (closeTimer) clearTimeout(closeTimer)
  window.removeEventListener('scroll', onViewportChange, true)
  window.removeEventListener('resize', onViewportChange)
})
</script>

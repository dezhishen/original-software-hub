<template>
  <aside
    class="relative overflow-hidden rounded-2xl border border-dashed border-brand-500/28 bg-[linear-gradient(180deg,rgba(255,255,255,0.94),rgba(236,253,245,0.96))] p-4 shadow-[0_12px_32px_rgba(15,157,132,0.08)] dark:border-brand-500/22 dark:bg-[linear-gradient(180deg,rgba(15,23,42,0.88),rgba(12,74,55,0.24))]"
    :class="wrapperClass"
    :aria-label="ariaLabel"
  >
    <div class="pointer-events-none absolute inset-x-0 top-0 h-24 bg-[radial-gradient(circle_at_top,rgba(16,185,129,0.18),transparent_70%)] dark:bg-[radial-gradient(circle_at_top,rgba(16,185,129,0.12),transparent_70%)]"></div>

    <div class="relative flex h-full flex-col gap-4">
      <div v-if="hasMultipleEmbeds" class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            class="inline-flex h-8 w-8 items-center justify-center rounded-full border border-brand-500/22 bg-white/80 text-brand-700 transition hover:border-brand-500/40 hover:bg-white focus:outline-none focus:ring-4 focus:ring-brand-500/15 dark:border-brand-400/20 dark:bg-slate-950/42 dark:text-brand-300 dark:hover:bg-slate-900"
            aria-label="上一条广告"
            @click="goPrev"
          >
            <span class="text-sm leading-none">‹</span>
          </button>
          <button
            type="button"
            class="inline-flex h-8 w-8 items-center justify-center rounded-full border border-brand-500/22 bg-white/80 text-brand-700 transition hover:border-brand-500/40 hover:bg-white focus:outline-none focus:ring-4 focus:ring-brand-500/15 dark:border-brand-400/20 dark:bg-slate-950/42 dark:text-brand-300 dark:hover:bg-slate-900"
            aria-label="下一条广告"
            @click="goNext"
          >
            <span class="text-sm leading-none">›</span>
          </button>
        </div>

        <div class="flex items-center gap-1.5">
          <button
            v-for="(item, index) in normalizedEmbeds"
            :key="item.id || index"
            type="button"
            class="h-2.5 rounded-full transition"
            :class="index === activeIndex ? 'w-6 bg-brand-600 dark:bg-brand-400' : 'w-2.5 bg-brand-500/24 hover:bg-brand-500/40 dark:bg-brand-400/24 dark:hover:bg-brand-400/40'"
            :aria-label="`切换到第 ${index + 1} 条广告`"
            :aria-current="index === activeIndex ? 'true' : undefined"
            @click="setActive(index)"
          ></button>
        </div>
      </div>

      <div class="flex min-h-[220px] flex-1 items-center justify-center rounded-2xl border border-white/80 bg-white/60 px-4 dark:border-slate-700/70 dark:bg-slate-950/24">
        <AdEmbedRenderer
          :key="activeEmbed?.id || activeIndex"
          :embed="activeEmbed"
          :aria-label="ariaLabel"
        />
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import AdEmbedRenderer from './AdEmbedRenderer.vue'

const props = defineProps({
  sticky: { type: Boolean, default: true },
  ariaLabel: { type: String, default: '广告位占位区域' },
  embeds: { type: Array, default: () => [] },
  rotateMs: { type: Number, default: 6000 }
})

const activeIndex = ref(0)
let rotateTimer = null

const wrapperClass = computed(() => (
  props.sticky ? 'xl:sticky xl:top-24' : ''
))

const normalizedEmbeds = computed(() => props.embeds.filter((item) => {
  if (!item || typeof item !== 'object') return false
  const type = String(item.type || '').toLowerCase()
  if (type === 'script') return Boolean(item.slotHtml || item.scripts || item.title)
  if (type === 'iframe') return Boolean(item.src || item.srcdoc)
  if (type === 'html') return Boolean(item.html || item.title)
  return Boolean(item.src || item.srcdoc || item.html || item.slotHtml || item.scripts || item.title)
}))
const activeEmbed = computed(() => normalizedEmbeds.value[activeIndex.value] ?? null)
const hasMultipleEmbeds = computed(() => normalizedEmbeds.value.length > 1)

function setActive(index) {
  if (!normalizedEmbeds.value.length) return
  activeIndex.value = ((index % normalizedEmbeds.value.length) + normalizedEmbeds.value.length) % normalizedEmbeds.value.length
  restartRotation()
}

function goPrev() {
  setActive(activeIndex.value - 1)
}

function goNext() {
  setActive(activeIndex.value + 1)
}

function stopRotation() {
  if (rotateTimer) {
    window.clearInterval(rotateTimer)
    rotateTimer = null
  }
}

function startRotation() {
  stopRotation()
  if (!hasMultipleEmbeds.value || props.rotateMs <= 0) return
  rotateTimer = window.setInterval(() => {
    activeIndex.value = (activeIndex.value + 1) % normalizedEmbeds.value.length
  }, props.rotateMs)
}

function restartRotation() {
  startRotation()
}

watch(normalizedEmbeds, (embeds) => {
  if (!embeds.length) {
    activeIndex.value = 0
    stopRotation()
    return
  }
  if (activeIndex.value >= embeds.length) activeIndex.value = 0
  startRotation()
}, { immediate: true })

watch(() => props.rotateMs, () => {
  startRotation()
})

onMounted(() => {
  startRotation()
})

onBeforeUnmount(() => {
  stopRotation()
})
</script>
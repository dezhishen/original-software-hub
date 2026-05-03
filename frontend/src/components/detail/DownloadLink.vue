<!--
  DownloadLink — 单个下载入口徽章
  Props:
    link — { type: 'direct' | 'webpage' | 'store', label: string, url: string }
-->
<template>
  <a
    :href="link.url"
    target="_blank"
    rel="noopener noreferrer"
    class="inline-flex items-center rounded-md border px-2 py-1 text-[11px] font-semibold"
    :class="styleClass"
  >
    {{ link.label }}
    <span v-if="link.type === 'webpage'" class="ml-1 text-[10px] font-medium opacity-70">页面</span>
    <span v-else-if="link.type === 'store'" class="ml-1 text-[10px] font-medium opacity-70">商店</span>
  </a>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  link: { type: Object, required: true }
})

// 不同链接类型对应的配色
const STYLE_MAP = {
  webpage: 'border-slate-300 bg-white text-slate-700 hover:border-brand-500/40 hover:text-brand-700 dark:border-slate-600 dark:bg-slate-800/80 dark:text-slate-300 dark:hover:bg-slate-700',
  store:   'border-emerald-300/60 bg-emerald-50 text-emerald-700 hover:bg-emerald-100 dark:border-emerald-700/50 dark:bg-emerald-900/25 dark:text-emerald-300 dark:hover:bg-emerald-900/35',
  direct:  'border-brand-500/30 bg-brand-50 text-brand-700 hover:bg-brand-100 dark:border-brand-500/40 dark:bg-slate-700/50 dark:text-brand-300 dark:hover:bg-slate-700'
}

const styleClass = computed(() => STYLE_MAP[props.link.type] ?? STYLE_MAP.direct)
</script>

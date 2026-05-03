<!--
  SoftwareIcon — 统一渲染软件图标
  - 有效图片 URL → <img>
  - 其他 → 文字占位 <span>
  Props:
    software  — SoftwareItem 对象（需要 icon、name 字段）
    size      — 'sm' (h-9 w-9) | 'md' (h-12 w-12)，默认 'sm'
-->
<template>
  <img
    v-if="isImageUrl"
    :src="software.icon"
    :alt="`${software.name} 图标`"
    :class="imgClass"
    loading="lazy"
  />
  <span v-else :class="spanClass">{{ fallbackChar }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  software: { type: Object, required: true },
  size: { type: String, default: 'sm' }
})

const icon = computed(() => String(props.software?.icon || '').trim())

const isImageUrl = computed(() =>
  icon.value &&
  (/^https?:\/\//i.test(icon.value) ||
   icon.value.startsWith('/') ||
   icon.value.startsWith('./'))
)

const fallbackChar = computed(
  () => (!isImageUrl.value && icon.value)
    || props.software?.name?.slice(0, 1)?.toUpperCase()
    || '?'
)

const sizeClass = computed(() => props.size === 'md' ? 'h-12 w-12' : 'h-9 w-9')

const imgClass = computed(() =>
  `${sizeClass.value} rounded-lg border border-slate-200 bg-white p-1 object-contain dark:border-slate-700 dark:bg-slate-800`
)

const spanClass = computed(() =>
  `inline-flex ${sizeClass.value} items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-sm font-semibold text-slate-600 dark:border-slate-700 dark:bg-slate-800/90`
)
</script>

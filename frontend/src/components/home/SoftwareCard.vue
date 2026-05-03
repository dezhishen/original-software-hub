<template>
  <article
    class="software-card relative overflow-hidden cursor-pointer rounded-xl border border-slate-200/90 bg-white/92 p-3.5 shadow-[0_6px_16px_rgba(15,70,56,0.08)] transition hover:-translate-y-0.5 hover:border-brand-500/45 hover:shadow-[0_10px_20px_rgba(15,157,132,0.14)] dark:border-slate-700/80 dark:bg-slate-800/88 dark:shadow-[0_6px_16px_rgba(2,6,23,0.35)] dark:hover:shadow-[0_10px_20px_rgba(15,157,132,0.18)]"
    @click="$emit('select', software.id)"
  >
    <!-- 图标背景虚化层（仅对图片图标显示） -->
    <div v-if="isImageIcon" class="software-card-bg pointer-events-none absolute inset-0">
      <div class="absolute -right-6 -top-6 h-24 w-24 rounded-full bg-brand-500/10 blur-2xl dark:bg-brand-500/14"></div>
      <div class="absolute -left-5 -bottom-5 h-20 w-20 rounded-full bg-amber-200/35 blur-2xl dark:bg-slate-600/18"></div>
      <div
        class="software-card-bg-icon absolute inset-0 opacity-[0.26] blur-[2px] dark:opacity-[0.2]"
        :style="{ backgroundImage: `url('${software.icon}')` }"
      ></div>
      <div class="absolute inset-0 bg-gradient-to-r from-white/32 via-white/20 to-white/10 dark:from-slate-800/45 dark:via-slate-800/30 dark:to-slate-800/18"></div>
    </div>

    <!-- 图标 + 名称 -->
    <div class="relative mb-2.5 flex items-center gap-3">
      <SoftwareIcon :software="software" size="sm" />
      <h3
        class="text-[17px] font-semibold leading-6 text-slate-900 dark:text-slate-100"
        style="font-family: 'Space Grotesk', sans-serif;"
      >{{ software.name }}</h3>
    </div>

    <!-- 描述 -->
    <p class="software-card-description relative mt-1.5 text-[13px] leading-5 text-slate-600 dark:text-slate-400">
      {{ software.description }}
    </p>

    <!-- 标签 -->
    <div v-if="software.tags?.length" class="relative mt-2.5 flex flex-wrap gap-1.5">
      <button
        v-for="tag in software.tags"
        :key="tag"
        type="button"
        class="inline-block rounded-full bg-brand-50/80 px-2 py-0.5 text-xs font-medium text-brand-700 transition hover:bg-brand-100 dark:bg-slate-700/50 dark:text-brand-400 dark:hover:bg-slate-700"
        @click.stop="$emit('tag-select', tag)"
      >#{{ tag }}</button>
    </div>

    <!-- 机构 -->
    <p class="relative mt-2 text-xs text-slate-500 dark:text-slate-500">
      机构：{{ software.organization }}
    </p>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import SoftwareIcon from '@/components/SoftwareIcon.vue'

const props = defineProps({
  software: { type: Object, required: true }
})

defineEmits(['select', 'tag-select'])

const isImageIcon = computed(() => {
  const icon = String(props.software?.icon || '').trim()
  return icon && (/^https?:\/\//i.test(icon) || icon.startsWith('/') || icon.startsWith('./'))
})
</script>


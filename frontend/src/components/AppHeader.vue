<template>
  <header class="app-shell mx-auto pt-3 pb-2 text-left md:pt-4 md:pb-2">
    <div class="rounded-2xl border border-slate-200/80 bg-white/72 px-3 py-2.5 shadow-[0_6px_18px_rgba(15,70,56,0.07)] backdrop-blur-sm md:flex md:items-center md:justify-between md:px-4 md:py-3 dark:border-slate-700/75 dark:bg-slate-800/58 dark:shadow-[0_6px_16px_rgba(2,6,23,0.28)]">
      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-2">
          <span class="inline-block h-2 w-2 rounded-full bg-brand-500/80 shadow-[0_0_0_5px_rgba(15,157,132,0.12)]"></span>
        <h1
          class="text-lg font-semibold tracking-[0.01em] text-slate-900 md:text-xl dark:text-slate-100"
          style="font-family: 'Space Grotesk', sans-serif;"
        >{{ title }}</h1>
        <span
          class="inline-flex rounded-full border border-brand-500/22 bg-brand-50/80 px-2.5 py-0.5 text-[10px] font-medium tracking-[0.08em] text-brand-700 dark:border-brand-500/35 dark:bg-slate-800/80 dark:text-brand-400"
          style="font-family: 'Space Grotesk', sans-serif;"
        >{{ badge }}</span>
      </div>
      <p
        v-if="description"
        class="mt-1 hidden truncate pr-2 text-[13px] leading-5 text-slate-600 md:block dark:text-slate-400"
      >{{ description }}</p>
      <div class="mt-1 flex flex-wrap items-center gap-2">
        <p v-if="updatedAt" class="text-xs text-slate-500 dark:text-slate-400">
          数据更新于 <RelativeTime :dateStr="updatedAt" />
        </p>
        <span class="inline-flex rounded-md border border-slate-200/85 bg-white/75 px-2 py-0.5 text-[11px] text-slate-500 dark:border-slate-700/80 dark:bg-slate-800/70 dark:text-slate-400">
          官方直链优先
        </span>
      </div>
    </div>
      <ComplianceNotice class="mt-2 md:mt-0 md:ml-4 md:shrink-0" />
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { pageState } from '@/stores/pageState'
import ComplianceNotice from './ComplianceNotice.vue'
import RelativeTime from './RelativeTime.vue'

const isDetail = computed(() => pageState.mode === 'detail')

const title = computed(() =>
  isDetail.value ? (pageState.detailSoftware?.name || '软件详情') : '常用软件下载导航'
)
const badge = computed(() => (isDetail.value ? 'SOFTWARE DETAIL' : 'ORIGINAL SOFTWARE HUB'))
const description = computed(() =>
  isDetail.value
    ? (pageState.detailSoftware?.description || '查看软件介绍、版本发布时间和下载入口。')
    : '帮你快速找到常见软件的官网、适合自己设备的版本，以及对应的官方下载方式。'
)
const updatedAt = computed(() =>
  isDetail.value ? pageState.detailUpdatedAt : pageState.homeUpdatedAt
)
</script>

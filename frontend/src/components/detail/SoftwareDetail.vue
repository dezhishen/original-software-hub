<template>
  <article
    class="relative overflow-hidden rounded-2xl border border-slate-200/90 bg-white/20 backdrop-blur-sm p-4 shadow-[0_6px_16px_rgba(15,70,56,0.08)] md:p-5 xl:p-6 dark:border-slate-700/80 dark:bg-slate-800/40 dark:shadow-[0_6px_16px_rgba(2,6,23,0.35)]"
  >
    <!-- Detail hero background -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-xl">
      <div class="absolute -right-12 -top-12 h-48 w-48 rounded-full bg-brand-500/8 blur-3xl dark:bg-brand-500/12"></div>
      <div class="absolute -left-12 -bottom-8 h-40 w-40 rounded-full bg-amber-200/25 blur-3xl dark:bg-slate-600/15"></div>
      <div class="absolute inset-0 bg-gradient-to-r from-white/30 via-white/20 to-white/10 dark:from-slate-800/40 dark:via-slate-800/30 dark:to-slate-800/20"></div>
    </div>

    <div class="relative flex flex-col gap-3 lg:grid lg:grid-cols-[minmax(0,1fr)_auto] lg:items-start lg:gap-4">
      <div class="min-w-0 flex-1">
        <!-- Back / Home buttons -->
        <div class="mb-2 flex flex-wrap items-center gap-2">
          <button
            type="button"
            class="inline-flex items-center rounded-lg bg-brand-500 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-brand-700 focus:outline-none focus:ring-4 focus:ring-brand-500/20"
            @click="$emit('back')"
          >← 返回上一页</button>
          <button
            type="button"
            class="inline-flex items-center rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 transition hover:border-brand-500/40 hover:text-brand-700 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-300"
            @click="$emit('home')"
          >回到软件目录</button>
        </div>

        <!-- 图标 + 环境 meta -->
        <div class="mb-2.5 flex items-center gap-3">
          <div class="shrink-0">
            <SoftwareIcon :software="software" size="md" />
          </div>
          <div class="min-w-0 flex flex-1 flex-wrap items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
            <span class="rounded-full bg-white/80 px-2.5 py-1 shadow-sm dark:bg-slate-800/80">
              所属机构：{{ software.organization }}
            </span>
            <span class="rounded-full bg-brand-50 px-2.5 py-1 font-medium text-brand-700 dark:bg-slate-700/60 dark:text-brand-300">
              当前检测环境：{{ currentPlatform.label }} / {{ currentArchitecture.label }}
            </span>
          </div>
        </div>

        <!-- Tags -->
        <div v-if="software.tags?.length" class="mt-2 flex flex-wrap gap-1.5">
          <span
            v-for="tag in software.tags"
            :key="tag"
            class="inline-block rounded-full bg-brand-50/80 px-2 py-0.5 text-xs font-medium text-brand-700 dark:bg-slate-700/50 dark:text-brand-400"
          >#{{ tag }}</span>
        </div>
      </div>

      <!-- Official website link -->
      <div class="relative flex shrink-0 items-start lg:pt-0.5">
        <a
          :href="software.officialWebsite"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex w-fit items-center rounded-lg border border-brand-500/35 bg-brand-50 px-3 py-1.5 text-xs font-medium text-brand-700 hover:bg-brand-100 dark:border-brand-500/40 dark:bg-slate-700/50 dark:text-brand-300 dark:hover:bg-slate-700"
        >访问官网</a>
      </div>
    </div>

    <!-- Versions section -->
    <section class="relative mt-3 border-t border-slate-200/90 pt-2.5 dark:border-slate-700/80">
      <div class="mb-2">
        <h3
          class="text-sm font-semibold text-slate-700 dark:text-slate-200"
          style="font-family: 'Space Grotesk', sans-serif;"
        >版本信息</h3>
        <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">
          以下版本与下载入口均归属于 {{ software.name }}
        </p>
      </div>

      <p
        v-if="!platforms?.length"
        class="rounded-lg border border-slate-200 bg-slate-50/90 px-3 py-4 text-sm text-slate-600 dark:border-slate-700 dark:bg-slate-800/70 dark:text-slate-400"
      >暂无版本信息，请访问官网获取最新版本。</p>

      <PlatformTabs
        v-else
        :platforms="platforms"
        :current-platform="currentPlatform"
        :current-architecture="currentArchitecture"
      />
    </section>
  </article>
</template>

<script setup>
import PlatformTabs from './PlatformTabs.vue'
import SoftwareIcon from '@/components/SoftwareIcon.vue'
import { usePlatform } from '@/composables/usePlatform'

const props = defineProps({
  software: { type: Object, required: true },
  platforms: { type: Array, default: () => [] }
})

defineEmits(['back', 'home'])

const { currentPlatform, currentArchitecture } = usePlatform()
</script>

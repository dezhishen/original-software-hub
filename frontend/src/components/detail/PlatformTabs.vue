<template>
  <!-- Platform tabs -->
  <div class="overflow-hidden rounded-lg border border-slate-200/90 bg-white/20 dark:border-slate-700/80 dark:bg-slate-800/20">
    <!-- Tab buttons -->
    <div
      v-if="sortedPlatforms.length > 1"
      role="tablist"
      class="flex flex-wrap items-end gap-1 border-b border-slate-200/80 bg-slate-50/60 px-3 pt-1.5 dark:border-slate-700/70 dark:bg-slate-900/35"
    >
      <button
        v-for="platform in sortedPlatforms"
        :key="platform.platform"
        type="button"
        role="tab"
        :aria-selected="platform.platform === activePlatform"
        class="-mb-px border-b-2 px-3 py-2 text-[12px] font-semibold transition"
        :class="platform.platform === activePlatform
          ? 'border-brand-500 text-brand-700 dark:border-brand-400 dark:text-brand-300'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:border-slate-600'"
        @click="activePlatform = platform.platform"
      >
        {{ platform.platform }}
        <span class="ml-1 rounded bg-slate-200/80 px-1.5 py-0.5 text-[10px] font-medium text-slate-600 dark:bg-slate-700/70 dark:text-slate-300">
          {{ platform.packages.length }}
        </span>
        <span v-if="isPlatformCurrent(platform.platform)" class="ml-1 text-[10px] opacity-80">当前</span>
      </button>
    </div>

    <!-- Tab panels -->
    <div class="platform-panels">
      <div
        v-for="platform in sortedPlatforms"
        :key="platform.platform"
        role="tabpanel"
        :class="{ hidden: platform.platform !== activePlatform }"
      >
        <!-- Platform meta row -->
        <div class="mb-2 border-b border-slate-200/70 bg-gradient-to-r from-sky-50/80 to-teal-50/70 px-3 py-2 dark:border-slate-700/70 dark:from-slate-800/50 dark:to-slate-800/20">
          <div class="flex flex-wrap items-center gap-2">
            <span class="inline-flex items-center rounded-md border border-sky-300/70 bg-sky-100/80 px-2 py-1 text-[11px] font-semibold text-sky-700 dark:border-sky-700/60 dark:bg-sky-900/30 dark:text-sky-300">
              平台：{{ platform.platform }}
            </span>
            <span
              class="rounded-full bg-brand-50 px-2.5 py-1 text-[11px] font-medium text-brand-700 dark:bg-slate-700/50 dark:text-brand-300"
              style="font-family: 'Space Grotesk', sans-serif;"
            >{{ platform.version || '-' }}</span>
            <span class="text-[11px] text-slate-600 dark:text-slate-300">{{ platform.releaseDate || '' }}</span>
            <a
              v-if="platform.officialUrl"
              :href="platform.officialUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center rounded-md border border-amber-300/50 bg-amber-100/20 px-2 py-1 text-[11px] font-semibold text-amber-700 hover:bg-amber-100/30 dark:border-amber-700/40 dark:bg-amber-900/15 dark:text-amber-300 dark:hover:bg-amber-900/25"
            >前往官网发布页</a>
            <span class="text-[11px] text-slate-500 dark:text-slate-400">共 {{ platform.packages.length }} 个安装包</span>
          </div>
        </div>

        <!-- Packages table -->
        <div v-if="sortedPackages(platform).length" class="overflow-x-auto">
          <table class="min-w-full border-collapse">
            <thead class="bg-transparent">
              <tr>
                <th class="px-2.5 py-2 text-left text-[11px] font-semibold tracking-wide text-slate-600 dark:text-slate-300">架构</th>
                <th class="px-2.5 py-2 text-left text-[11px] font-semibold tracking-wide text-slate-600 dark:text-slate-300">下载入口</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-slate-700">
              <tr
                v-for="(pkg, idx) in sortedPackages(platform)"
                :key="idx"
                class="bg-white/10 even:bg-slate-50/10 hover:bg-slate-100/20 dark:bg-slate-800/10 dark:even:bg-slate-800/15 dark:hover:bg-slate-700/20"
                :class="{ 'font-semibold': idx === 0 && isPlatformCurrent(platform.platform) }"
              >
                <td class="whitespace-nowrap px-2.5 py-2 text-[13px] text-slate-700 dark:text-slate-200">
                  {{ pkg.architecture || '-' }}
                </td>
                <td class="px-2.5 py-2 text-[13px] text-slate-700 dark:text-slate-200">
                  <div v-if="pkg.links?.length" class="flex flex-wrap gap-1.5">
                    <DownloadLink
                      v-for="link in pkg.links"
                      :key="link.url"
                      :link="link"
                    />
                  </div>
                  <span v-else class="text-slate-400">暂无直链</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-else class="px-3 py-4 text-sm text-slate-600 dark:text-slate-400">该平台暂无构建信息。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import DownloadLink from './DownloadLink.vue'
import { normalizePlatformId, architectureScore } from '@/composables/usePlatform'

const props = defineProps({
  platforms: { type: Array, required: true },
  currentPlatform: { type: Object, required: true },
  currentArchitecture: { type: Object, required: true }
})

function isPlatformCurrent(label) {
  return normalizePlatformId(label) === props.currentPlatform.id
}

const sortedPlatforms = computed(() => {
  return [...props.platforms].sort((a, b) => {
    const aCur = isPlatformCurrent(a.platform) ? 1 : 0
    const bCur = isPlatformCurrent(b.platform) ? 1 : 0
    if (aCur !== bCur) return bCur - aCur
    return String(a.platform).localeCompare(String(b.platform), 'zh-CN')
  })
})

const activePlatform = ref(sortedPlatforms.value[0]?.platform || '')

function sortedPackages(platform) {
  return [...(platform.packages || [])].sort((a, b) => {
    const sa = architectureScore(a.architecture, props.currentArchitecture.id)
    const sb = architectureScore(b.architecture, props.currentArchitecture.id)
    if (sa !== sb) return sb - sa
    return String(a.architecture || '').localeCompare(String(b.architecture || ''), 'zh-CN')
  })
}
</script>

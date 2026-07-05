<template>
  <div class="flex min-h-0 flex-1 flex-col">
    <!-- Search bar -->
    <div class="mb-2 border-b border-slate-200 pb-2 dark:border-slate-700">
      <input
        v-model="keyword"
        type="search"
        placeholder="搜索软件名称 / 机构 / 标签（#社交）/ 拼音（weixin）"
        class="w-full rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm outline-none ring-brand-500/20 transition focus:border-brand-500 focus:ring-4 md:max-w-[34rem] dark:border-slate-700 dark:bg-slate-900/90 dark:text-slate-100 dark:placeholder-slate-500"
      />
    </div>

    <!-- Category bar -->
    <CategoryBar
      :categories="categoryStats"
      :active="activeCategory"
      :total-count="softwares.length"
      @update:active="activeCategory = $event; keyword = ''"
    />

    <!-- Software grid -->
    <div
      class="home-list-panel grid auto-rows-max content-start gap-3 overflow-y-auto overscroll-contain pt-1 pr-1 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-[repeat(auto-fill,minmax(240px,1fr))]"
      aria-live="polite"
    >
      <p
        v-if="filtered.length === 0"
        class="rounded-xl border border-slate-200 bg-slate-50 px-4 py-6 text-sm text-slate-600 dark:border-slate-700 dark:bg-slate-800/50"
      >没有匹配的软件，请尝试其他关键词。</p>

      <SoftwareCard
        v-for="software in filtered"
        :key="software.id"
        :software="software"
        @select="$emit('select', $event)"
        @tag-select="onTagSelect"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import SoftwareCard from './SoftwareCard.vue'
import CategoryBar from './CategoryBar.vue'

const props = defineProps({
  softwares: { type: Array, default: () => [] },
  initialKeyword: { type: String, default: '' }
})

const emit = defineEmits(['select'])

const keyword = ref(props.initialKeyword)
const activeCategory = ref('')

/** 统计每个分类的软件数量 */
const categoryStats = computed(() => {
  const map = {}
  for (const s of props.softwares) {
    const cats = Array.isArray(s.categories) ? s.categories : []
    for (const c of cats) {
      map[c] = (map[c] || 0) + 1
    }
  }
  return Object.entries(map).map(([key, count]) => ({ key, count }))
})

/** 按分类过滤 */
const byCategory = computed(() => {
  if (!activeCategory.value) return props.softwares
  return props.softwares.filter((s) => {
    const cats = Array.isArray(s.categories) ? s.categories : []
    return cats.includes(activeCategory.value)
  })
})

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  const source = byCategory.value
  if (!kw) return source
  const tagKw = kw.startsWith('#') ? kw.slice(1).trim() : ''
  return source.filter((s) => {
    const tags = Array.isArray(s.tags) ? s.tags : []
    if (tagKw) return tags.some((t) => String(t || '').toLowerCase().includes(tagKw))
    return `${s.name} ${s.organization} ${tags.join(' ')} ${s.pinyin || ''}`.toLowerCase().includes(kw)
  })
})

function onTagSelect(tag) {
  keyword.value = `#${tag}`
}
</script>

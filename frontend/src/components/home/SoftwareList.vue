<template>
  <div>
    <!-- Search bar -->
    <div class="mb-2 border-b border-slate-200 pb-2 dark:border-slate-700">
      <input
        v-model="keyword"
        type="search"
        placeholder="搜索软件名称 / 机构 / 标签（#社交）/ 拼音（weixin）"
        class="w-full rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm outline-none ring-brand-500/20 transition focus:border-brand-500 focus:ring-4 dark:border-slate-700 dark:bg-slate-900/90 dark:text-slate-100 dark:placeholder-slate-500"
      />
    </div>

    <!-- Software grid -->
    <div
      class="home-list-panel grid auto-rows-max content-start grid-cols-1 gap-3.5 overflow-y-auto overscroll-contain pt-1 pr-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
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
import { ref, computed, watch } from 'vue'
import SoftwareCard from './SoftwareCard.vue'

const props = defineProps({
  softwares: { type: Array, default: () => [] },
  initialKeyword: { type: String, default: '' }
})

const emit = defineEmits(['select'])

const keyword = ref(props.initialKeyword)

watch(() => props.initialKeyword, (v) => { keyword.value = v })

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return props.softwares
  const tagKw = kw.startsWith('#') ? kw.slice(1).trim() : ''
  return props.softwares.filter((s) => {
    const tags = Array.isArray(s.tags) ? s.tags : []
    if (tagKw) return tags.some((t) => String(t || '').toLowerCase().includes(tagKw))
    return `${s.name} ${s.organization} ${tags.join(' ')} ${s.pinyin || ''}`.toLowerCase().includes(kw)
  })
})

function onTagSelect(tag) {
  keyword.value = `#${tag}`
}
</script>

<template>
  <div
    class="category-bar -mx-3 mb-2 overflow-x-auto overscroll-x-contain border-b border-slate-200 px-3 pb-2 scrollbar-none md:-mx-4 md:px-4 dark:border-slate-700"
  >
    <div class="flex gap-1.5">
      <button
        type="button"
        :class="activeClass('')"
        @click="$emit('update:active', '')"
      >
        全部分类
        <span class="ml-1 text-[11px] opacity-60">{{ totalCount }}</span>
      </button>

      <button
        v-for="cat in visibleCategories"
        :key="cat.key"
        type="button"
        :class="activeClass(cat.key)"
        @click="$emit('update:active', cat.key)"
      >
        {{ cat.label }}
        <span class="ml-1 text-[11px] opacity-60">{{ cat.count }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  categories: { type: Array, default: () => [] },
  active: { type: String, default: '' },
  totalCount: { type: Number, default: 0 },
})

defineEmits(['update:active'])

const labelMap = {
  browser: '浏览器',
  security: '安全',
  office: '办公',
  media: '影音',
  design: '设计',
  development: '开发',
  communication: '通讯',
  remote: '远程',
  'cloud-storage': '云盘',
  'input-method': '输入法',
  compression: '压缩',
  virtualization: '虚拟机',
  translation: '翻译',
  'note-taking': '笔记',
  gaming: '游戏',
}

const visibleCategories = computed(() =>
  props.categories
    .filter((c) => c.count > 0)
    .sort((a, b) => b.count - a.count || a.key.localeCompare(b.key))
    .map((c) => ({
      key: c.key,
      label: labelMap[c.key] || c.key,
      count: c.count,
    }))
)

function activeClass(key) {
  const isActive = props.active === key
  return [
    'shrink-0 whitespace-nowrap rounded-full px-3.5 py-1.5 text-sm font-medium transition',
    isActive
      ? 'bg-brand-500 text-white shadow-sm dark:bg-brand-700'
      : 'bg-slate-100 text-slate-700 hover:bg-slate-200 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700',
  ].join(' ')
}
</script>

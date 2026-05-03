<!--
  ErrorDisplay — 居中展示错误信息 + 返回按钮
  Props:
    title   — 标题文本
    message — 说明文本
    action  — 按钮文字，默认"← 返回目录"
  Emits:
    action  — 按钮被点击
-->
<template>
  <div class="min-h-[280px] grid place-items-center text-center px-4">
    <div>
      <h2
        class="text-xl font-semibold text-slate-700 dark:text-slate-200"
        style="font-family: 'Space Grotesk', sans-serif;"
      >{{ title }}</h2>
      <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">{{ message }}</p>
      <button
        type="button"
        class="mt-4 inline-flex items-center rounded-lg bg-brand-500 px-3 py-1.5 text-xs font-medium text-white hover:bg-brand-700"
        @click="$emit('action')"
      >{{ action }}</button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title:   { type: String, required: true },
  message: { type: String, default: '' },
  action:  { type: String, default: '← 返回目录' }
})

defineEmits(['action'])
</script>

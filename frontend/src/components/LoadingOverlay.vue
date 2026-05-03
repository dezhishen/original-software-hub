<template>
  <Transition name="overlay">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex flex-col items-center justify-center gap-4 bg-white/90 backdrop-blur-sm dark:bg-slate-900/90"
    >
      <div class="h-11 w-11 animate-spin rounded-full border-4 border-brand-500/20 border-t-brand-500"></div>
      <p class="text-sm text-slate-600 dark:text-slate-400">{{ message }}</p>
    </div>
  </Transition>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  message: { type: String, default: '正在加载...' }
})
</script>

<style scoped>
.overlay-enter-active,
.overlay-leave-active {
  transition: opacity 0.3s ease;
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}
</style>

<template>
  <section
    class="rounded-2xl border border-slate-200/90 bg-white/92 p-3 shadow-[0_8px_20px_rgba(15,70,56,0.07)] md:p-4 dark:border-slate-700 dark:bg-slate-800/88 dark:shadow-[0_8px_20px_rgba(2,6,23,0.35)]"
  >
    <LoadingOverlay :visible="loading" message="正在加载软件列表..." />
    <div
      v-if="error"
      class="rounded-xl border border-rose-200 bg-rose-50 px-4 py-6 text-sm text-rose-700 dark:border-rose-800 dark:bg-rose-900/20 dark:text-rose-300"
    >数据加载失败：{{ error }}</div>
    <SoftwareList
      v-else
      :softwares="softwares"
      @select="navigateToDetail"
    />
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import SoftwareList from '@/components/home/SoftwareList.vue'
import LoadingOverlay from '@/components/LoadingOverlay.vue'
import { dataRepository } from '@/services/dataRepository'
import { pageState } from '@/stores/pageState'

const router = useRouter()
const softwares = ref([])
const loading = ref(true)
const error = ref('')

async function loadHomeCatalog() {
  pageState.mode = 'home'
  pageState.detailSoftware = null
  pageState.detailUpdatedAt = ''
  pageState.transitionLoading = false
  pageState.transitionMessage = '正在加载...'
  try {
    const catalog = await dataRepository.loadSoftwareCatalog()
    softwares.value = catalog.softwares
    pageState.homeUpdatedAt = catalog.generatedAt
  } catch (e) {
    error.value = e instanceof Error ? e.message : '未知错误'
  } finally {
    loading.value = false
  }
}

onMounted(loadHomeCatalog)

function navigateToDetail(softwareId) {
  const to = `/software/${softwareId}`
  pageState.transitionMessage = '正在加载详情...'
  pageState.transitionLoading = true

  if (typeof document.startViewTransition === 'function') {
    try {
      document.startViewTransition(() => router.push(to))
      return
    } catch {
      // Fallback for browsers/environments where the API exists but cannot run now.
    }
  }
  router.push(to)
}
</script>

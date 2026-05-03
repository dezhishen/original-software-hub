<template>
  <div class="pt-1">
    <LoadingOverlay :visible="loading" message="正在加载详情..." />

    <ErrorDisplay
      v-if="errorKind"
      :title="errorKind === 'not-found' ? '未找到软件' : '加载失败'"
      :message="errorKind === 'not-found' ? '请返回目录页选择有效的软件。' : errorMessage"
      @action="router.push('/')"
    />

    <SoftwareDetail
      v-else-if="software && !loading"
      :software="software"
      :platforms="platforms"
      @back="goBack"
      @home="router.push('/')"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import SoftwareDetail from '@/components/detail/SoftwareDetail.vue'
import LoadingOverlay from '@/components/LoadingOverlay.vue'
import ErrorDisplay from '@/components/ErrorDisplay.vue'
import { dataRepository } from '@/services/dataRepository'
import { normalizeSoftwareVersionPayload } from '@/utils/normalize'
import { pageState } from '@/stores/pageState'

const props = defineProps({ id: { type: String, required: true } })
const router = useRouter()

const software = ref(null)
const platforms = ref([])
const loading = ref(true)
const errorKind = ref('')      // '' | 'not-found' | 'error'
const errorMessage = ref('')

async function loadDetail(id) {
  loading.value = true
  errorKind.value = ''
  software.value = null
  platforms.value = []
  pageState.mode = 'detail'
  pageState.detailSoftware = null
  pageState.detailUpdatedAt = ''

  try {
    const sw = await dataRepository.getSoftwareById(id)
    if (!sw) {
      errorKind.value = 'not-found'
      return
    }

    software.value = sw
    pageState.detailSoftware = sw

    const rawVersions = await dataRepository.loadSoftwareVersions(sw)
    const { platforms: pf, updatedAt } = normalizeSoftwareVersionPayload(rawVersions)
    platforms.value = pf
    pageState.detailUpdatedAt = updatedAt
  } catch (e) {
    errorKind.value = 'error'
    errorMessage.value = e instanceof Error ? e.message : '未知错误'
  } finally {
    loading.value = false
  }
}

watch(() => props.id, loadDetail, { immediate: true })

function goBack() {
  const sameOrigin = document.referrer.startsWith(window.location.origin)
  if (window.history.length > 1 && sameOrigin) router.back()
  else router.replace('/')
}
</script>

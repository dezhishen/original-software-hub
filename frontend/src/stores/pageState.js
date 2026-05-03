import { reactive } from 'vue'

// 全局页面状态 — 由 HomeView/DetailView 写入，由 AppHeader 读取
export const pageState = reactive({
  /** @type {'home' | 'detail'} */
  mode: 'home',
  homeUpdatedAt: '',
  /** @type {import('../services/dataRepository').SoftwareItem | null} */
  detailSoftware: null,
  detailUpdatedAt: ''
})

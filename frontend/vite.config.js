import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { resolve } from 'path'
import { fileURLToPath } from 'url'

const __dirname = fileURLToPath(new URL('.', import.meta.url))

export default defineConfig({
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      // 只预缓存构建产物，运行时数据 JSON 不缓存
      workbox: {
        globPatterns: ['**/*.{js,css,html,svg,png,ico}'],
        navigateFallback: null,
        runtimeCaching: [
          {
            // 软件目录与版本 JSON — 网络优先，离线时返回缓存
            urlPattern: /\/data\/json\/.+\.json(\?.*)?$/,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'osh-data-cache',
              networkTimeoutSeconds: 5,
              cacheableResponse: { statuses: [0, 200] }
            }
          }
        ]
      },
      manifest: false,          // 使用 public/manifest.webmanifest
      includeManifestIcons: false
    })
  ],
  base: './',
  build: {
    outDir: 'dist',
    emptyOutDir: true
  },
  resolve: {
    alias: { '@': resolve(__dirname, './src') }
  }
})

<template>
  <div class="relative min-h-screen">
    <Transition name="detail-bg-fade" mode="out-in">
      <div
        v-if="showDetailBackdrop"
        class="global-detail-backdrop"
        :style="detailBackdropStyle"
      ></div>
    </Transition>

    <div class="relative z-10 flex min-h-screen flex-col">
      <AppHeader />
      <main class="app-shell mx-auto flex-1 min-h-0 pb-5 md:pb-6">
        <RouterView v-slot="{ Component, route }">
          <Transition name="route-fade-slide" mode="out-in">
            <component :is="Component" :key="route.fullPath" />
          </Transition>
        </RouterView>
      </main>
      <DarkModeToggle />
    </div>

    <LoadingOverlay :visible="pageState.transitionLoading" :message="pageState.transitionMessage" />
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import LoadingOverlay from '@/components/LoadingOverlay.vue'
import { pageState } from '@/stores/pageState'
import { useDetailBackdrop } from '@/composables/useDetailBackdrop'

const route = useRoute()
const { showDetailBackdrop, detailBackdropStyle } = useDetailBackdrop(pageState, route)
</script>

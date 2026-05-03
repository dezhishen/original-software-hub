<template>
  <div class="relative min-h-screen">
    <Transition name="detail-bg-fade" mode="out-in">
      <div
        v-if="showDetailBackdrop"
        class="global-detail-backdrop"
        :style="detailBackdropStyle"
      ></div>
    </Transition>

    <div class="relative z-10">
      <AppHeader />
      <main class="app-shell mx-auto pb-8 md:pb-9">
        <RouterView v-slot="{ Component, route }">
          <Transition name="route-fade-slide" mode="out-in">
            <component :is="Component" :key="route.fullPath" />
          </Transition>
        </RouterView>
      </main>
      <DarkModeToggle />
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import { pageState } from '@/stores/pageState'
import { useDetailBackdrop } from '@/composables/useDetailBackdrop'

const route = useRoute()
const { showDetailBackdrop, detailBackdropStyle } = useDetailBackdrop(pageState, route)
</script>

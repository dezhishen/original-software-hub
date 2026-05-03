<template>
  <AppHeader />
  <main class="app-shell mx-auto pb-8 md:pb-9">
    <RouterView v-slot="{ Component, route }">
      <Transition name="route-fade-slide" mode="out-in">
        <component :is="Component" :key="route.fullPath" />
      </Transition>
    </RouterView>
  </main>
  <DarkModeToggle />
</template>

<script setup>
import AppHeader from '@/components/AppHeader.vue'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
</script>

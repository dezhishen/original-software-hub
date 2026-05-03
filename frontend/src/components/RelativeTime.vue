<template>
  <time
    v-if="dateStr"
    :datetime="iso"
    :title="`${absolute}（北京时间）`"
    class="cursor-help underline decoration-dotted decoration-slate-400"
  >{{ relative }}</time>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { relativeTimeText, absoluteTimeText, isoTimeText } from '@/utils/format'

const props = defineProps({ dateStr: String })

const tick = ref(0)
let timer

onMounted(() => { timer = setInterval(() => tick.value++, 60000) })
onUnmounted(() => clearInterval(timer))

const relative = computed(() => { tick.value; return relativeTimeText(props.dateStr) })
const absolute = computed(() => absoluteTimeText(props.dateStr))
const iso = computed(() => isoTimeText(props.dateStr))
</script>

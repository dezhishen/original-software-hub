<template>
  <iframe
    v-if="embed?.type === 'iframe' && (embed.src || embed.srcdoc)"
    :src="embed.src || undefined"
    :srcdoc="embed.srcdoc || undefined"
    :title="embed.title || ariaLabel"
    class="h-full min-h-[320px] w-full rounded-[1.1rem] border-0 bg-white dark:bg-slate-950"
    loading="lazy"
    referrerpolicy="strict-origin-when-cross-origin"
    allowfullscreen
  ></iframe>

  <div
    v-else-if="embed?.type === 'html' && embed.html"
    class="ad-slot-html h-full min-h-[320px] w-full overflow-hidden rounded-[1.1rem]"
    v-html="embed.html"
  ></div>

  <div
    v-else-if="embed?.type === 'script'"
    ref="scriptHost"
    class="ad-slot-script h-full min-h-[320px] w-full overflow-hidden rounded-[1.1rem]"
  ></div>

  <div
    v-else-if="embed"
    class="flex h-full min-h-[320px] w-full items-center justify-center rounded-[1.1rem] border border-dashed border-brand-500/18 bg-[linear-gradient(180deg,rgba(255,255,255,0.3),rgba(255,255,255,0.12))] px-4 text-center text-sm text-slate-400 dark:border-brand-400/14 dark:bg-[linear-gradient(180deg,rgba(15,23,42,0.16),rgba(15,23,42,0.04))] dark:text-slate-500"
  >
    {{ embed.title || '广告内容待接入' }}
  </div>

  <div
    v-else
    class="flex h-full min-h-[320px] w-full items-center justify-center rounded-[1.1rem] border border-dashed border-brand-500/18 bg-[linear-gradient(180deg,rgba(255,255,255,0.3),rgba(255,255,255,0.12))] dark:border-brand-400/14 dark:bg-[linear-gradient(180deg,rgba(15,23,42,0.16),rgba(15,23,42,0.04))]"
  ></div>
</template>

<script setup>
import { nextTick, onMounted, ref, watch } from 'vue'

const props = defineProps({
  embed: { type: Object, default: null },
  ariaLabel: { type: String, default: '广告位占位区域' }
})

const scriptHost = ref(null)

function applyScriptAttrs(scriptEl, attrs) {
  if (!attrs || typeof attrs !== 'object') return
  Object.entries(attrs).forEach(([key, value]) => {
    if (value === undefined || value === null) return
    scriptEl.setAttribute(key, String(value))
  })
}

async function renderScriptEmbed() {
  const host = scriptHost.value
  if (!host) return
  host.innerHTML = ''

  if (!props.embed || props.embed.type !== 'script') return

  if (props.embed.slotHtml) {
    host.insertAdjacentHTML('beforeend', props.embed.slotHtml)
  }

  const scripts = Array.isArray(props.embed.scripts)
    ? props.embed.scripts
    : (props.embed.scripts ? [props.embed.scripts] : [])

  for (const item of scripts) {
    if (!item || typeof item !== 'object') continue
    const scriptEl = document.createElement('script')
    if (item.src) scriptEl.src = String(item.src)
    if (item.type) scriptEl.type = String(item.type)
    if (item.async !== undefined) scriptEl.async = Boolean(item.async)
    if (item.defer !== undefined) scriptEl.defer = Boolean(item.defer)
    applyScriptAttrs(scriptEl, item.attrs)
    if (item.inline) scriptEl.textContent = String(item.inline)
    host.appendChild(scriptEl)
  }
}

watch(
  () => props.embed,
  async (embed) => {
    if (embed?.type !== 'script') return
    await nextTick()
    renderScriptEmbed()
  },
  { immediate: true, deep: true }
)

onMounted(() => {
  if (props.embed?.type === 'script') renderScriptEmbed()
})
</script>
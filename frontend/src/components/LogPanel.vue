<template>
  <div v-if="logs && logs.length" class="mt-auto bg-slate-900/80 border border-slate-700/60 rounded-xl overflow-hidden shadow-lg backdrop-blur-md">
    <div class="flex items-center justify-between px-3.5 py-2 bg-slate-950/70 cursor-pointer select-none border-b border-slate-800" @click="isOpen = !isOpen">
      <div class="flex items-center gap-1.5 text-xs font-bold text-slate-300 uppercase tracking-wide">
        <span>{{ t('battle_log_title') }}</span>
        <span class="text-slate-500 font-mono">({{ logs.length }})</span>
      </div>
      <button class="btn btn-ghost btn-xs h-5 min-h-0 text-slate-400 p-0">
        {{ isOpen ? '▲' : '▼' }}
      </button>
    </div>

    <div v-show="isOpen" ref="logBodyRef" class="p-3 max-h-28 overflow-y-auto font-mono text-[11.5px] leading-relaxed text-slate-400 flex flex-col gap-1">
      <div v-for="(log, i) in logs" :key="i" class="text-slate-300/90 hover:text-slate-100 transition-colors">
        {{ formatLog(log) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { useI18n } from '../composables/useI18n';

const props = defineProps({
  logs: { type: Array, default: () => [] },
});

const { t, lang } = useI18n();
const isOpen = ref(true);
const logBodyRef = ref(null);

function formatLog(log) {
  if (!log) return '';
  if (log.includes(' / ')) {
    const parts = log.split(' / ');
    return lang.value === 'zh' ? parts[0] : (parts[1] || parts[0]);
  }
  return log;
}

watch(
  () => props.logs,
  () => {
    nextTick(() => {
      if (logBodyRef.value) {
        logBodyRef.value.scrollTop = logBodyRef.value.scrollHeight;
      }
    });
  },
  { deep: true }
);
</script>

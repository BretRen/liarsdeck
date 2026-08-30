<template>
  <div v-if="logs && logs.length" class="log-panel glass-panel">
    <div class="log-header" @click="isOpen = !isOpen">
      <div class="log-title">
        <span>{{ t('battle_log_title') }}</span>
        <span class="log-count">({{ logs.length }})</span>
      </div>
      <button class="btn-icon collapse-btn">
        {{ isOpen ? '▲' : '▼' }}
      </button>
    </div>

    <div v-show="isOpen" ref="logBodyRef" class="log-body">
      <div v-for="(log, i) in logs" :key="i" class="log-entry">
        {{ log }}
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

const { t } = useI18n();
const isOpen = ref(true);
const logBodyRef = ref(null);

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

<style scoped>
.log-panel {
  background: #10131b;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  overflow: hidden;
  margin-top: auto;
}

.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background: #141822;
  cursor: pointer;
  user-select: none;
}

.log-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.log-count {
  color: var(--text-muted);
}

.collapse-btn {
  font-size: 9px;
  padding: 2px 6px;
  background: transparent;
  border: none;
}

.log-body {
  padding: 6px 12px;
  max-height: 100px;
  overflow-y: auto;
  font-family: 'SF Mono', 'Fira Code', Consolas, monospace;
  font-size: 11.5px;
  line-height: 1.5;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-top: 1px solid var(--border-subtle);
}

.log-entry {
  word-break: break-all;
}

.log-entry:last-child {
  color: var(--text-primary);
  font-weight: 500;
}
</style>

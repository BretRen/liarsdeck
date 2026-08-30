<template>
  <div v-if="isDisconnected" class="dc-modal-overlay">
    <div class="dc-modal glass-panel">
      <!-- Modal Header with Language Toggle -->
      <div class="dc-header">
        <h2 class="dc-title font-serif">{{ t('dc_title') }}</h2>
        <button class="btn-icon lang-btn" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Modal Body -->
      <div class="dc-body">
        <div class="spinner-container">
          <div class="pulsing-ring"></div>
          <span class="dc-icon">⚡</span>
        </div>

        <p class="dc-desc">{{ t('dc_desc') }}</p>

        <!-- 30s Grace Countdown when match is in progress -->
        <div v-if="state.status === 'paused' && pauseRemainingSeconds > 0" class="grace-box">
          <span class="grace-label">{{ t('pause_countdown_label') }}</span>
          <span class="grace-seconds">{{ pauseRemainingSeconds }}s</span>
          <p class="grace-tip">{{ t('dc_grace_tip') }}</p>
        </div>

        <div class="dc-status">
          <span class="pulse-dot"></span>
          <span>{{ t('dc_reconnecting') }}</span>
        </div>
      </div>

      <!-- Modal Actions -->
      <div class="dc-actions">
        <button class="btn-primary flex-1" :disabled="isReconnecting" @click="tryReconnect">
          {{ t('dc_retry_btn') }}
        </button>
        <button class="btn-secondary flex-1" @click="exitToLobby">
          {{ t('dc_exit_btn') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const { t, lang, toggleLang } = useI18n();
const { state, isDisconnected, isReconnecting, pauseRemainingSeconds, tryReconnect, exitToLobby } =
  useGameStore();
</script>

<style scoped>
.dc-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(10px);
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: fadeIn 0.25s ease;
}

.dc-modal {
  width: 100%;
  max-width: 420px;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.8), 0 0 30px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.dc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-brass);
  background: #10131a;
}

.dc-title {
  font-size: 1.15rem;
  color: #ff8787;
  margin: 0;
}

.lang-btn {
  font-weight: 700;
}

.dc-body {
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
}

.spinner-container {
  position: relative;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulsing-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid var(--border-brass);
  border-top-color: var(--gold-accent);
  animation: spin 1.2s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.dc-icon {
  font-size: 22px;
  color: var(--gold-accent);
}

.dc-desc {
  font-size: 0.95rem;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0;
}

.grace-box {
  background: #1c1313;
  border: 1px solid #742a2a;
  border-radius: var(--radius-md);
  padding: 10px 14px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.grace-label {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
}

.grace-seconds {
  font-family: var(--font-serif);
  font-size: 1.8rem;
  font-weight: 900;
  color: #ff8787;
  line-height: 1;
}

.grace-tip {
  font-size: 12px;
  color: #ffa8a8;
  margin: 0;
}

.dc-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--gold-accent);
  font-weight: 600;
}

.pulse-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--gold-accent);
  box-shadow: 0 0 8px var(--gold-accent);
  animation: pulseGold 1.5s infinite;
}

.dc-actions {
  display: flex;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-brass);
  background: #10131a;
}

.flex-1 {
  flex: 1;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
}
</style>

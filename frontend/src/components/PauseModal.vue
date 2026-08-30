<template>
  <div v-if="isPaused && !isDisconnected" class="pause-modal-overlay">
    <div class="pause-modal glass-panel">
      <!-- Header -->
      <div class="pause-header">
        <div class="header-title-group">
          <span class="pause-icon-badge">⏸</span>
          <h2 class="font-serif">{{ t('pause_modal_title') }}</h2>
        </div>
        <button class="btn-icon lang-btn" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Body -->
      <div class="pause-body">
        <div class="pause-desc">
          {{ t('pause_modal_desc', { name: state.paused_player }) }}
        </div>

        <!-- 30s Countdown Display -->
        <div class="countdown-card" :class="{ urgent: pauseRemainingSeconds <= 10 }">
          <span class="countdown-label">{{ t('pause_countdown_label') }}</span>
          <div class="countdown-clock">
            <span class="countdown-num">{{ pauseRemainingSeconds }}</span>
            <span class="countdown-unit">s</span>
          </div>

          <!-- Progress Bar -->
          <div class="progress-track">
            <div
              class="progress-fill"
              :style="{ width: `${(Math.max(0, pauseRemainingSeconds) / 30) * 100}%` }"
            ></div>
          </div>
        </div>

        <p class="pause-tip-text">
          {{ t('pause_tip') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const { t, lang, toggleLang } = useI18n();
const { state, isDisconnected, pauseRemainingSeconds } = useGameStore();

const isPaused = computed(() => state.value.status === 'paused');
</script>

<style scoped>
.pause-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.82);
  backdrop-filter: blur(8px);
  z-index: 2500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  animation: fadeIn 0.25s ease;
}

.pause-modal {
  width: 100%;
  max-width: 440px;
  background: #141720;
  border: 1.5px solid #d4af37;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.85), 0 0 25px rgba(212, 175, 55, 0.2);
  overflow: hidden;
}

.pause-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #10131a;
  border-bottom: 1px solid var(--border-brass);
}

.header-title-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pause-icon-badge {
  font-size: 16px;
  color: var(--gold-accent);
}

.pause-header h2 {
  font-size: 1.15rem;
  color: var(--gold-accent);
  margin: 0;
}

.lang-btn {
  font-weight: 700;
}

.pause-body {
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 18px;
}

.pause-desc {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.5;
}

.countdown-card {
  width: 100%;
  background: #0f121a;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  padding: 14px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.countdown-label {
  font-size: 11.5px;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.countdown-clock {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.countdown-num {
  font-family: var(--font-serif);
  font-size: 2.4rem;
  font-weight: 900;
  color: var(--gold-accent);
  line-height: 1;
}

.countdown-unit {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-secondary);
}

.countdown-card.urgent .countdown-num {
  color: #ff8787;
  text-shadow: 0 0 12px rgba(255, 135, 135, 0.6);
}

.progress-track {
  width: 100%;
  height: 6px;
  background: #242938;
  border-radius: 3px;
  overflow: hidden;
  margin-top: 4px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #c59b27, #f59e0b);
  border-radius: 3px;
  transition: width 0.25s linear;
}

.countdown-card.urgent .progress-fill {
  background: linear-gradient(90deg, #c92a2a, #ff8787);
}

.pause-tip-text {
  font-size: 0.85rem;
  color: var(--text-muted);
  line-height: 1.5;
  margin: 0;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
}
</style>

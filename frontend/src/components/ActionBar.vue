<template>
  <div class="action-bar glass-panel">
    <!-- Waiting Phase -->
    <template v-if="status === 'waiting'">
      <div v-if="isPlayer" class="action-group">
        <button class="btn-primary" @click="$emit('toggle-ready')">
          {{ isReady ? t('unready_btn') : t('ready_btn') }}
        </button>

        <button
          class="btn-success"
          :disabled="!canStart"
          @click="$emit('start-game')"
          v-if="amHost"
        >
          {{ t('start_game_btn') }}
          <span v-if="!canStart" class="btn-subtip">
            ({{ allReady ? t('need_more_players') : t('all_ready_needed') }})
          </span>
        </button>
      </div>

      <div v-else class="status-hint">
        <span>{{ t('spectator_banner') }}</span>
      </div>
    </template>

    <!-- Playing Phase -->
    <template v-else-if="status === 'playing'">
      <div v-if="isPlayer && isMyTurn" class="action-group">
        <!-- Turn Timer Badge inside Action Bar -->
        <div v-if="remainingSeconds > 0" class="action-timer" :class="{ urgent: remainingSeconds <= 8 }">
          <span class="timer-label">{{ t('timeout_warn') }}:</span>
          <span class="timer-val">{{ remainingSeconds }}s</span>
        </div>

        <button
          class="btn-primary btn-play"
          :disabled="!canPlay"
          @click="$emit('play-cards')"
        >
          {{ t('play_cards_btn') }} ({{ selectedCount }}/3)
        </button>

        <button
          class="btn-danger btn-liar"
          :class="{ 'pulse-liar': canCallLiar }"
          :disabled="!canCallLiar"
          @click="$emit('call-liar')"
        >
          {{ t('call_liar_btn') }}
        </button>
      </div>

      <div v-else class="status-hint">
        <span class="pulse-dot"></span>
        <span v-if="isPlayer">{{ t('status_playing') }} ({{ t('status_waiting') }})</span>
        <span v-else>{{ t('spectator_banner') }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

defineProps({
  status: { type: String, default: 'waiting' },
  isPlayer: { type: Boolean, default: true },
  isMyTurn: { type: Boolean, default: false },
  isReady: { type: Boolean, default: false },
  amHost: { type: Boolean, default: false },
  canPlay: { type: Boolean, default: false },
  canCallLiar: { type: Boolean, default: false },
  canStart: { type: Boolean, default: false },
  allReady: { type: Boolean, default: false },
  selectedCount: { type: Number, default: 0 },
});

defineEmits(['toggle-ready', 'start-game', 'play-cards', 'call-liar']);
const { t } = useI18n();
const store = useGameStore();

const remainingSeconds = computed(() => store.remainingSeconds.value);
</script>

<style scoped>
.action-bar {
  padding: 12px 18px;
  margin-bottom: 14px;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 64px;
}

.action-group {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
  width: 100%;
}

.action-group button {
  flex: 1;
  min-width: 140px;
  max-width: 260px;
  padding: 11px 18px;
  font-size: 14px;
}

.action-timer {
  display: flex;
  align-items: baseline;
  gap: 4px;
  background: #0f121a;
  border: 1px solid var(--border-brass);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  color: var(--gold-accent);
}

.action-timer.urgent {
  background: #2b1414;
  border-color: #821c1c;
  color: #ff8787;
  animation: pulseRed 1s infinite;
}

.timer-val {
  font-family: var(--font-serif);
  font-weight: 700;
  font-size: 15px;
}

.btn-subtip {
  font-size: 11px;
  font-weight: 500;
  opacity: 0.85;
}

.btn-liar {
  font-family: var(--font-serif);
  letter-spacing: 0.5px;
  font-weight: 700;
}

.status-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.pulse-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--gold-accent);
  box-shadow: 0 0 8px var(--gold-accent);
  animation: pulseGold 1.5s infinite;
}

@keyframes pulseGold {
  0%, 100% { opacity: 0.4; transform: scale(0.9); }
  50% { opacity: 1; transform: scale(1.2); }
}
</style>

<template>
  <div class="action-bar glass-panel">
    <!-- Waiting Phase Controls -->
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
          🎮 {{ t('start_game_btn') }}
          <span v-if="!canStart" class="btn-subtip">
            ({{ allReady ? t('need_more_players') : t('all_ready_needed') }})
          </span>
        </button>
      </div>

      <div v-else class="spectator-action-hint">
        <span>👀 {{ t('spectator_banner') }}</span>
      </div>
    </template>

    <!-- Playing Phase Controls -->
    <template v-else-if="status === 'playing'">
      <div v-if="isPlayer && isMyTurn" class="action-group">
        <button
          class="btn-primary btn-play"
          :disabled="!canPlay"
          @click="$emit('play-cards')"
        >
          🃏 {{ t('play_cards_btn') }} ({{ selectedCount }}/3)
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

      <div v-else class="waiting-turn-hint">
        <span class="pulse-dot"></span>
        <span v-if="isPlayer">{{ t('status_playing') }} ({{ t('status_waiting') }})</span>
        <span v-else>{{ t('spectator_banner') }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';

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
</script>

<style scoped>
.action-bar {
  padding: 14px 20px;
  margin-bottom: 16px;
  background: rgba(18, 21, 28, 0.95);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 68px;
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
  max-width: 280px;
  padding: 12px 20px;
  font-size: 15px;
}

.btn-subtip {
  font-size: 11px;
  font-weight: 500;
  opacity: 0.85;
}

.btn-liar {
  font-family: var(--font-heading);
  letter-spacing: 0.5px;
}

.waiting-turn-hint, .spectator-action-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 500;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent-gold);
  box-shadow: 0 0 10px var(--accent-gold);
  animation: pulseGold 1.5s infinite;
}

@keyframes pulseGold {
  0%, 100% { opacity: 0.4; transform: scale(0.9); }
  50% { opacity: 1; transform: scale(1.2); }
}
</style>

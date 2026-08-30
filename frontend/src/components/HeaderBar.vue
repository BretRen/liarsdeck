<template>
  <header class="header-bar glass-panel">
    <!-- Left: Room Info -->
    <div class="header-left">
      <div class="room-code-badge">
        <span class="room-code-label">ROOM</span>
        <span class="room-code-val">{{ roomCode }}</span>
      </div>
      <button class="btn-secondary btn-sm" @click="copyInvite" :title="t('invite_btn')">
        📋 {{ t('invite_btn') }}
      </button>
    </div>

    <!-- Center: Status & Timer -->
    <div class="header-center">
      <div class="status-pill" :class="status">
        <span class="status-dot"></span>
        {{ statusLabel }}
      </div>
      <div v-if="status === 'playing' && deadline" class="timer-badge" :class="{ warning: remainingTime <= 10 }">
        <span class="timer-icon">⏱</span>
        <span class="timer-text">{{ Math.max(0, remainingTime) }}s</span>
      </div>
    </div>

    <!-- Right: Utility Controls -->
    <div class="header-right">
      <button class="btn-icon" @click="audio.toggleMute" :title="audio.isMuted.value ? t('audio_off') : t('audio_on')">
        {{ audio.isMuted.value ? '🔇' : '🔊' }}
      </button>
      <button class="btn-icon" @click="$emit('open-rules')" :title="t('rules_btn')">
        📖
      </button>
      <button class="btn-icon lang-btn" @click="toggleLang">
        {{ lang.toUpperCase() }}
      </button>
      <button class="btn-icon exit-btn" @click="$emit('leave')" title="离开房间">
        🚪
      </button>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  roomCode: { type: String, default: '' },
  status: { type: String, default: 'waiting' },
  deadline: { type: Number, default: 0 },
  currentUnix: { type: Number, default: 0 },
});

defineEmits(['open-rules', 'leave']);

const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const store = useGameStore();

const remainingTime = computed(() => {
  return props.deadline - props.currentUnix;
});

const statusLabel = computed(() => {
  if (props.status === 'playing') return t('status_playing');
  if (props.status === 'game_over') return t('status_game_over');
  return t('status_waiting');
});

function copyInvite() {
  store.copyInvite();
}
</script>

<style scoped>
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 18px;
  margin-bottom: 16px;
  background: rgba(18, 21, 28, 0.9);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  flex-wrap: wrap;
  gap: 12px;
}

.header-left, .header-center, .header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.room-code-badge {
  display: flex;
  align-items: center;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid var(--border-gold);
  border-radius: var(--radius-md);
  padding: 4px 10px;
  gap: 6px;
}

.room-code-label {
  font-size: 10px;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.room-code-val {
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 1.1rem;
  letter-spacing: 2px;
  color: var(--accent-gold);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
  border-radius: var(--radius-sm);
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-pill.waiting {
  background: rgba(156, 163, 175, 0.12);
  color: var(--text-secondary);
  border: 1px solid rgba(156, 163, 175, 0.25);
}
.status-pill.waiting .status-dot { background: var(--text-secondary); }

.status-pill.playing {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}
.status-pill.playing .status-dot { background: #34d399; box-shadow: 0 0 8px #34d399; }

.status-pill.game_over {
  background: rgba(245, 158, 11, 0.15);
  color: var(--accent-gold);
  border: 1px solid var(--border-gold);
}
.status-pill.game_over .status-dot { background: var(--accent-gold); }

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.timer-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
  padding: 4px 10px;
  border-radius: 20px;
  font-variant-numeric: tabular-nums;
  font-weight: 700;
  font-size: 13px;
}

.timer-badge.warning {
  background: rgba(239, 68, 68, 0.3);
  border-color: #ef4444;
  color: #ffffff;
  animation: pulseRed 1s infinite;
}

.lang-btn {
  font-weight: 700;
  font-size: 12px;
}

.exit-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5);
  color: #f87171;
}
</style>

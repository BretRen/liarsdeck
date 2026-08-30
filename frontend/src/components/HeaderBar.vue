<template>
  <header class="header-bar glass-panel">
    <!-- Left: Room Code & Invite -->
    <div class="header-left">
      <div class="room-tag">
        <span class="room-label">ROOM</span>
        <span class="room-code">{{ roomCode }}</span>
      </div>
      <button class="btn-secondary btn-sm" @click="copyInvite" :title="t('invite_btn')">
        {{ t('invite_btn') }}
      </button>
    </div>

    <!-- Center: Match Status & Turn Countdown Timer -->
    <div class="header-center">
      <div class="status-indicator" :class="status">
        <span class="status-dot"></span>
        <span>{{ statusLabel }}</span>
      </div>

      <!-- Prominent Countdown Timer -->
      <div
        v-if="status === 'playing' && deadline"
        class="turn-timer"
        :class="{ urgent: remainingSeconds <= 8 }"
      >
        <span class="timer-num">{{ remainingSeconds }}</span>
        <span class="timer-sec">s</span>
      </div>
    </div>

    <!-- Right: Utility Controls -->
    <div class="header-right">
      <button class="btn-icon" @click="audio.toggleMute" :title="audio.isMuted.value ? t('audio_off') : t('audio_on')">
        {{ audio.isMuted.value ? 'Muted' : 'Sound' }}
      </button>
      <button class="btn-icon" @click="$emit('open-rules')">
        {{ t('rules_btn') }}
      </button>
      <button class="btn-icon lang-btn" @click="toggleLang">
        {{ lang.toUpperCase() }}
      </button>
      <button class="btn-icon exit-btn" @click="$emit('leave')" title="离开房间">
        退出
      </button>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useAudio } from '../composables/useAudio';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  roomCode: { type: String, default: '' },
  status: { type: String, default: 'waiting' },
  deadline: { type: Number, default: 0 },
});

defineEmits(['open-rules', 'leave']);

const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const store = useGameStore();

const remainingSeconds = computed(() => store.remainingSeconds.value);

const statusLabel = computed(() => {
  if (props.status === 'playing') return t('status_playing');
  if (props.status === 'paused') return t('status_paused');
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
  padding: 10px 16px;
  margin-bottom: 14px;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  flex-wrap: wrap;
  gap: 10px;
}

.header-left, .header-center, .header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.room-tag {
  display: flex;
  align-items: baseline;
  background: #0f121a;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-sm);
  padding: 4px 10px;
  gap: 6px;
}

.room-label {
  font-size: 10px;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.room-code {
  font-family: var(--font-serif);
  font-weight: 700;
  font-size: 1.1rem;
  letter-spacing: 2px;
  color: var(--gold-accent);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  background: #0f1219;
  border: 1px solid var(--border-subtle);
}

.status-indicator.waiting { color: var(--text-secondary); }
.status-indicator.waiting .status-dot { background: var(--text-muted); }

.status-indicator.playing { color: #51cf66; border-color: rgba(81, 207, 102, 0.3); }
.status-indicator.playing .status-dot { background: #51cf66; box-shadow: 0 0 6px #51cf66; }

.status-indicator.game_over { color: var(--gold-accent); border-color: var(--border-brass-bright); }
.status-indicator.game_over .status-dot { background: var(--gold-accent); }

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

/* Prominent Turn Timer */
.turn-timer {
  display: flex;
  align-items: baseline;
  gap: 2px;
  background: #1c1313;
  border: 1px solid #742a2a;
  color: #ff8787;
  padding: 3px 10px;
  border-radius: var(--radius-sm);
  font-variant-numeric: tabular-nums;
  font-weight: 700;
}

.timer-num {
  font-size: 16px;
  font-family: var(--font-serif);
}

.timer-sec {
  font-size: 11px;
  color: var(--text-muted);
}

.turn-timer.urgent {
  background: #3b1212;
  border-color: #e03131;
  color: #ffffff;
  animation: pulseRed 1s infinite ease-in-out;
}

.lang-btn {
  font-weight: 700;
}

.exit-btn:hover {
  background: #2b1414;
  border-color: #821c1c;
  color: #ff8787;
}
</style>

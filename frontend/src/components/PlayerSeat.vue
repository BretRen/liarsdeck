<template>
  <div
    class="player-seat"
    :class="{
      active: isActiveTurn && !player.is_spectator,
      dead: !player.is_alive && !player.is_spectator,
      spectator: player.is_spectator,
      'is-me': isMe,
    }"
  >
    <!-- Turn Indicator & Countdown Ring -->
    <div v-if="isActiveTurn && !player.is_spectator" class="turn-bar">
      <span class="turn-label">CURRENT TURN</span>
      <span v-if="remainingSeconds > 0" class="turn-seconds">{{ remainingSeconds }}s</span>
    </div>

    <!-- Header: Name & Role Tags -->
    <div class="player-header">
      <div class="player-name-wrapper">
        <span class="player-name" :title="player.nickname">{{ player.nickname }}</span>
        <span v-if="player.is_host" class="tag host-tag">{{ t('host_tag') }}</span>
        <span v-if="player.is_spectator" class="tag spec-tag">{{ t('spec_tag') }}</span>
        <span v-if="player.is_alive && !player.is_spectator && player.is_connected === false" class="tag offline-tag">{{ t('offline_tag') }}</span>
        <span v-if="isMe" class="tag me-tag">[{{ t('me_tag') }}]</span>
      </div>

      <!-- Kick Button for Host -->
      <button
        v-if="amHost && !player.is_host && !isMe"
        class="kick-btn"
        :title="t('kick_btn')"
        @click="$emit('kick', player.id)"
      >
        ✕
      </button>
    </div>

    <!-- Body for Non-Spectators -->
    <div v-if="!player.is_spectator" class="player-body">
      <!-- 6-Chamber Revolver Visual Indicator -->
      <div class="cylinder-row" :title="`${t('bullets_label')}: ${player.bullets}/6`">
        <div
          v-for="i in 6"
          :key="i"
          class="chamber-bullet"
          :class="{
            loaded: player.is_alive && i <= player.bullets,
            empty: player.is_alive && i > player.bullets,
            fatal: !player.is_alive,
          }"
        ></div>
      </div>

      <div class="player-meta">
        <div v-if="player.is_alive" class="cards-count">
          <span class="cards-label">{{ t('hand_count_label') }}:</span>
          <span class="cards-num">{{ player.hand ? player.hand.length : 0 }}</span>
        </div>
        <div v-else class="dead-notice">
          {{ t('dead_tag') }}
        </div>
      </div>

      <!-- Ready Status during Waiting Phase -->
      <div v-if="gameStatus === 'waiting'" class="ready-status-box" :class="{ ready: player.is_ready }">
        <span>{{ player.is_ready ? t('ready_status') : t('unready_status') }}</span>
      </div>
    </div>

    <!-- Spectator Body -->
    <div v-else class="spectator-body">
      <span>{{ t('watching_tag') }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  player: { type: Object, required: true },
  isActiveTurn: { type: Boolean, default: false },
  isMe: { type: Boolean, default: false },
  amHost: { type: Boolean, default: false },
  gameStatus: { type: String, default: 'waiting' },
});

defineEmits(['kick']);
const { t } = useI18n();
const store = useGameStore();

const remainingSeconds = computed(() => store.remainingSeconds.value);
</script>

<style scoped>
.player-seat {
  position: relative;
  flex: 1;
  min-width: 125px;
  background: #151821;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  transition: all 0.2s ease;
}

.player-seat.is-me {
  border-color: #3b5bdb;
  background: #161a26;
}

.player-seat.active {
  border-color: var(--gold-accent);
  background: #1c1c1a;
  box-shadow: 0 0 16px var(--gold-glow);
  transform: translateY(-2px);
}

.player-seat.dead {
  opacity: 0.4;
  filter: grayscale(0.9);
  border-color: #491212;
}

.player-seat.spectator {
  border-style: dashed;
  border-color: var(--border-subtle);
  background: #101217;
}

.turn-bar {
  position: absolute;
  top: -8px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--gold-accent);
  color: #0d0f14;
  font-family: var(--font-sans);
  font-weight: 800;
  font-size: 9px;
  letter-spacing: 0.5px;
  padding: 1px 8px;
  border-radius: 2px;
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.turn-seconds {
  background: #0d0f14;
  color: #ffffff;
  padding: 0 4px;
  border-radius: 2px;
  font-size: 8.5px;
}

.player-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
}

.player-name-wrapper {
  display: flex;
  align-items: center;
  gap: 4px;
  overflow: hidden;
}

.player-name {
  font-weight: 700;
  font-size: 0.88rem;
  color: var(--text-primary);
  max-width: 85px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag {
  font-size: 9px;
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 2px;
}

.host-tag {
  background: rgba(212, 175, 55, 0.2);
  color: var(--gold-accent);
  border: 1px solid var(--border-brass);
}

.spec-tag {
  background: #202430;
  color: var(--text-muted);
}

.offline-tag {
  background: #3b1414;
  border: 1px solid #822222;
  color: #ff8787;
  animation: pulseRed 1.2s infinite ease-in-out;
}

.me-tag {
  color: #748ffc;
}

.kick-btn {
  padding: 1px 5px;
  font-size: 9px;
  background: #2b1414;
  border: 1px solid #742a2a;
  color: #ff8787;
  border-radius: 2px;
}
.kick-btn:hover {
  background: #c92a2a;
  color: #fff;
}

.player-body {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

/* 6-Chamber Visual */
.cylinder-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  padding: 3px;
  background: #0d0f14;
  border-radius: 4px;
  border: 1px solid var(--border-subtle);
}

.chamber-bullet {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  background: #242938;
  border: 1px solid #363d50;
  transition: all 0.2s ease;
}

.chamber-bullet.loaded {
  background: #c59b27;
  border-color: #e0b438;
  box-shadow: 0 0 3px rgba(212, 175, 55, 0.5);
}

.chamber-bullet.empty {
  background: #141720;
  border-color: #242835;
}

.chamber-bullet.fatal {
  background: #c92a2a;
  border-color: #ff8787;
}

.player-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-secondary);
}

.cards-count {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.cards-num {
  font-family: var(--font-serif);
  font-weight: 700;
  color: var(--gold-accent);
}

.dead-notice {
  font-size: 11px;
  font-weight: 700;
  color: #ff8787;
}

.ready-status-box {
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 4px;
  border-radius: 2px;
  background: #1a1e28;
  color: var(--text-muted);
}

.ready-status-box.ready {
  background: rgba(43, 138, 62, 0.2);
  color: #51cf66;
  border: 1px solid rgba(43, 138, 62, 0.4);
}

.spectator-body {
  padding: 4px 0;
  text-align: center;
  font-size: 11px;
  color: var(--text-muted);
}
</style>

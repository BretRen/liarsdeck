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
    <!-- Turn Indicator Glow Banner -->
    <div v-if="isActiveTurn && !player.is_spectator" class="turn-ribbon">
      <span>TURN</span>
    </div>

    <!-- Player Header Info -->
    <div class="player-header">
      <div class="player-name-wrapper">
        <span class="player-name">{{ player.nickname }}</span>
        <span v-if="player.is_host" class="tag host-tag" :title="t('host_tag')">👑 {{ t('host_tag') }}</span>
        <span v-if="player.is_spectator" class="tag spec-tag">👀 {{ t('spec_tag') }}</span>
        <span v-if="isMe" class="tag me-tag">({{ t('me_tag') }})</span>
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

    <!-- Non-Spectator Info -->
    <div v-if="!player.is_spectator" class="player-body">
      <!-- Revolver Cylinder Slots Indicator (6 Chambers) -->
      <div class="revolver-cylinder" :title="`${t('bullets_label')}: ${player.bullets}/6`">
        <div
          v-for="i in 6"
          :key="i"
          class="chamber-slot"
          :class="{
            loaded: player.is_alive && i <= player.bullets,
            fired: player.is_alive && i > player.bullets,
            fatal: !player.is_alive,
          }"
        ></div>
      </div>

      <div class="player-meta">
        <div v-if="player.is_alive" class="meta-item">
          <span class="meta-icon">🃏</span>
          <span class="meta-val">{{ player.hand ? player.hand.length : 0 }}</span>
        </div>
        <div v-else class="dead-badge">
          {{ t('dead_tag') }}
        </div>
      </div>

      <!-- Ready Status in Waiting Stage -->
      <div v-if="gameStatus === 'waiting'" class="ready-badge" :class="{ ready: player.is_ready }">
        <span class="ready-icon">{{ player.is_ready ? '✓' : '…' }}</span>
        <span>{{ player.is_ready ? t('ready_status') : t('unready_status') }}</span>
      </div>
    </div>

    <!-- Spectator Info -->
    <div v-else class="spectator-body">
      <span class="spec-watching-text">👀 {{ t('watching_tag') }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';

const props = defineProps({
  player: { type: Object, required: true },
  isActiveTurn: { type: Boolean, default: false },
  isMe: { type: Boolean, default: false },
  amHost: { type: Boolean, default: false },
  gameStatus: { type: String, default: 'waiting' },
});

defineEmits(['kick']);
const { t } = useI18n();
</script>

<style scoped>
.player-seat {
  position: relative;
  flex: 1;
  min-width: 130px;
  background: var(--bg-card);
  border: 1.5px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 12px 14px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.35);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.player-seat.is-me {
  border-color: rgba(59, 130, 246, 0.4);
}

.player-seat.active {
  border-color: var(--accent-gold);
  background: radial-gradient(circle at 50% 0%, rgba(245, 158, 11, 0.15) 0%, var(--bg-card) 80%);
  box-shadow: 0 0 20px rgba(245, 158, 11, 0.3), 0 8px 24px rgba(0, 0, 0, 0.6);
  transform: translateY(-4px);
}

.player-seat.dead {
  opacity: 0.45;
  filter: grayscale(0.8);
  border-color: rgba(239, 68, 68, 0.2);
}

.player-seat.spectator {
  border-style: dashed;
  border-color: var(--border-subtle);
  background: rgba(18, 21, 28, 0.4);
}

.turn-ribbon {
  position: absolute;
  top: -9px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #0a0b0e;
  font-family: var(--font-heading);
  font-weight: 900;
  font-size: 10px;
  letter-spacing: 1px;
  padding: 2px 10px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.5);
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
  gap: 6px;
  flex-wrap: wrap;
}

.player-name {
  font-weight: 700;
  font-size: 0.92rem;
  color: var(--text-primary);
  max-width: 90px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
}

.host-tag {
  background: rgba(245, 158, 11, 0.18);
  color: var(--accent-gold);
  border: 1px solid var(--border-gold);
}

.spec-tag {
  background: rgba(156, 163, 175, 0.15);
  color: var(--text-secondary);
}

.me-tag {
  color: var(--accent-blue);
  font-weight: 700;
}

.kick-btn {
  padding: 2px 6px;
  font-size: 10px;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
  border-radius: 4px;
}
.kick-btn:hover {
  background: #ef4444;
  color: #fff;
}

.player-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* Revolver Chamber Slots */
.revolver-cylinder {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 4px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
}

.chamber-slot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  transition: all 0.25s ease;
}

.chamber-slot.loaded {
  background: linear-gradient(135deg, #fbbf24, #d97706);
  box-shadow: 0 0 5px rgba(245, 158, 11, 0.6);
}

.chamber-slot.fired {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--border-subtle);
}

.chamber-slot.fatal {
  background: #ef4444;
  box-shadow: 0 0 6px #ef4444;
}

.player-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  color: var(--text-secondary);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
}

.dead-badge {
  font-size: 12px;
  font-weight: 700;
  color: #f87171;
}

.ready-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 6px;
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
}

.ready-badge.ready {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.spectator-body {
  padding: 8px 0;
  text-align: center;
}

.spec-watching-text {
  font-size: 12px;
  color: var(--text-muted);
}
</style>

<template>
  <div class="app-root">
    <!-- Global Broadcast Banner -->
    <transition name="broadcast-slide">
      <div v-if="globalBroadcast.visible" class="global-broadcast-banner">
        <div class="broadcast-inner glass-panel">
          <div class="broadcast-icon">📢</div>
          <div class="broadcast-content">
            <span class="broadcast-tag font-serif">{{ t('global_broadcast_title') }}</span>
            <span class="broadcast-msg">{{ globalBroadcast.message }}</span>
          </div>
          <button class="broadcast-close" @click="dismissBroadcast" title="Close">✕</button>
        </div>
      </div>
    </transition>

    <!-- View Switcher -->
    <LobbyView v-if="!connected" @open-rules="showRules = true" />
    <GameView v-else @open-rules="showRules = true" />

    <!-- Rules Modal -->
    <RulesModal v-if="showRules" @close="showRules = false" />

    <!-- Admin Console (Ctrl + X) -->
    <AdminModal :is-open="showAdmin" @close="showAdmin = false" />

    <!-- Pause Modal for all players in room when game is paused -->
    <PauseModal />

    <!-- Disconnect / Reconnect Modal for disconnected player -->
    <DisconnectModal />

    <!-- Global Toast Notification -->
    <transition name="toast-fade">
      <div v-if="toast" class="global-toast">
        {{ toast }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from './composables/useI18n';
import { useGameStore } from './composables/useGameStore';
import LobbyView from './views/LobbyView.vue';
import GameView from './views/GameView.vue';
import RulesModal from './components/RulesModal.vue';
import AdminModal from './components/AdminModal.vue';
import PauseModal from './components/PauseModal.vue';
import DisconnectModal from './components/DisconnectModal.vue';

const { t } = useI18n();
const showRules = ref(false);
const showAdmin = ref(false);
const { connected, toast, globalBroadcast, dismissBroadcast } = useGameStore();

function onKeyDown(e) {
  if ((e.ctrlKey || e.metaKey) && (e.key === 'x' || e.key === 'X')) {
    e.preventDefault();
    showAdmin.value = !showAdmin.value;
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeyDown);
});

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown);
});
</script>

<style scoped>
.app-root {
  width: 100%;
  position: relative;
}

/* Global Broadcast Banner */
.global-broadcast-banner {
  position: fixed;
  top: 14px;
  left: 50%;
  transform: translateX(-50%);
  width: 92%;
  max-width: 820px;
  z-index: 10000;
}

.broadcast-inner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 18px;
  background: rgba(26, 20, 10, 0.95);
  border: 1.5px solid var(--accent-gold);
  border-radius: var(--radius-lg);
  box-shadow: 0 12px 36px rgba(217, 119, 6, 0.25), 0 0 20px rgba(0, 0, 0, 0.8);
  animation: pulse-glow 3s infinite alternate;
}

@keyframes pulse-glow {
  0% {
    box-shadow: 0 8px 24px rgba(217, 119, 6, 0.2), 0 0 15px rgba(0, 0, 0, 0.8);
  }
  100% {
    box-shadow: 0 12px 36px rgba(217, 119, 6, 0.45), 0 0 25px rgba(217, 119, 6, 0.2);
  }
}

.broadcast-icon {
  font-size: 20px;
  animation: horn-shake 1.5s ease infinite;
}

@keyframes horn-shake {
  0%, 100% { transform: rotate(0deg); }
  20% { transform: rotate(-12deg); }
  40% { transform: rotate(12deg); }
  60% { transform: rotate(-8deg); }
  80% { transform: rotate(8deg); }
}

.broadcast-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  text-align: left;
}

.broadcast-tag {
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--accent-gold);
  font-weight: 700;
  text-transform: uppercase;
}

.broadcast-msg {
  font-size: 14px;
  color: #fff;
  font-weight: 600;
  line-height: 1.4;
  word-break: break-word;
}

.broadcast-close {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: var(--text-muted);
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  transition: all 0.2s ease;
}

.broadcast-close:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: var(--accent-crimson);
  color: #fff;
}

.broadcast-slide-enter-active,
.broadcast-slide-leave-active {
  transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.broadcast-slide-enter-from,
.broadcast-slide-leave-to {
  opacity: 0;
  transform: translate(-50%, -30px);
}

/* Toast */
.global-toast {
  position: fixed;
  bottom: 28px;
  left: 50%;
  transform: translateX(-50%);
  background: #181d28;
  border: 1px solid var(--border-brass);
  color: var(--text-primary);
  padding: 9px 20px;
  border-radius: var(--radius-md);
  font-size: 13.5px;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.7);
  z-index: 9999;
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translate(-50%, 16px);
}
</style>

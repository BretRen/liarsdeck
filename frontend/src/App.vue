<template>
  <div class="app-root">
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
import { useGameStore } from './composables/useGameStore';
import LobbyView from './views/LobbyView.vue';
import GameView from './views/GameView.vue';
import RulesModal from './components/RulesModal.vue';
import AdminModal from './components/AdminModal.vue';
import PauseModal from './components/PauseModal.vue';
import DisconnectModal from './components/DisconnectModal.vue';

const showRules = ref(false);
const showAdmin = ref(false);
const { connected, toast } = useGameStore();

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

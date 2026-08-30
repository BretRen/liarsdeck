<template>
  <div class="app-root">
    <!-- View Switcher -->
    <LobbyView v-if="!connected" @open-rules="showRules = true" />
    <GameView v-else @open-rules="showRules = true" />

    <!-- Rules Modal -->
    <RulesModal v-if="showRules" @close="showRules = false" />

    <!-- Global Toast Notification -->
    <transition name="toast-fade">
      <div v-if="toast" class="global-toast">
        {{ toast }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useGameStore } from './composables/useGameStore';
import LobbyView from './views/LobbyView.vue';
import GameView from './views/GameView.vue';
import RulesModal from './components/RulesModal.vue';

const showRules = ref(false);
const { connected, toast } = useGameStore();
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
  background: rgba(24, 29, 40, 0.95);
  border: 1px solid var(--border-gold);
  color: #ffffff;
  padding: 10px 24px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.8), 0 0 15px rgba(245, 158, 11, 0.2);
  backdrop-filter: blur(10px);
  z-index: 9999;
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}
</style>

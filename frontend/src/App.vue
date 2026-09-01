<template>
  <div class="app-root p-3 md:p-5 min-h-screen">
    <!-- Fullscreen Fatal Gunshot Flash & Blood Vignette (Teleported to Body for True 100vw/100vh) -->
    <Teleport to="body">
      <transition name="toast-fade">
        <div v-if="isFatalShotActive" class="fatal-blood-overlay fixed inset-0 w-screen h-screen pointer-events-none z-[999999]">
          <!-- 1. Muzzle Flash -->
          <div class="fatal-flash fixed inset-0 w-screen h-screen"></div>

          <!-- 2. Fullscreen Deep Radial Blood Vignette -->
          <div class="fatal-vignette fixed inset-0 w-screen h-screen"></div>

          <!-- 3. Screen Edge Blood Splatter Accents -->
          <svg class="fixed inset-0 w-screen h-screen pointer-events-none opacity-85" xmlns="http://www.w3.org/2000/svg">
            <!-- Top left splatter -->
            <circle cx="3%" cy="4%" r="45" fill="rgba(159, 18, 57, 0.85)" filter="blur(8px)" />
            <circle cx="6%" cy="8%" r="22" fill="rgba(190, 18, 60, 0.75)" filter="blur(4px)" />
            <circle cx="2%" cy="14%" r="28" fill="rgba(136, 19, 55, 0.8)" filter="blur(6px)" />
            <!-- Top right splatter -->
            <circle cx="96%" cy="3%" r="50" fill="rgba(159, 18, 57, 0.85)" filter="blur(8px)" />
            <circle cx="93%" cy="7%" r="24" fill="rgba(190, 18, 60, 0.75)" filter="blur(4px)" />
            <circle cx="97%" cy="13%" r="30" fill="rgba(136, 19, 55, 0.8)" filter="blur(6px)" />
            <!-- Bottom left splatter -->
            <circle cx="3%" cy="95%" r="55" fill="rgba(136, 19, 55, 0.9)" filter="blur(10px)" />
            <circle cx="8%" cy="91%" r="28" fill="rgba(190, 18, 60, 0.8)" filter="blur(5px)" />
            <!-- Bottom right splatter -->
            <circle cx="97%" cy="96%" r="60" fill="rgba(136, 19, 55, 0.9)" filter="blur(10px)" />
            <circle cx="92%" cy="92%" r="32" fill="rgba(190, 18, 60, 0.8)" filter="blur(5px)" />
          </svg>
        </div>
      </transition>
    </Teleport>

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
    <LobbyView v-if="!connected" @open-rules="showRules = true" @open-changelog="showChangelog = true" />
    <GameView v-else @open-rules="showRules = true" @open-changelog="showChangelog = true" />

    <!-- Rules Modal -->
    <RulesModal v-if="showRules" @close="showRules = false" />

    <!-- Changelog Modal (Markdown Rendering) -->
    <ChangelogModal v-if="showChangelog" @close="showChangelog = false" />

    <!-- Admin Console (Ctrl + X) -->
    <AdminModal :is-open="showAdmin" @close="showAdmin = false" />

    <!-- Pause Modal for all players in room when game is paused -->
    <PauseModal />

    <!-- Disconnect / Reconnect Modal for disconnected player -->
    <DisconnectModal />

    <!-- Mandatory OAuth2 PKCE Login Modal -->
    <LoginModal />

    <!-- Global Toast Notification -->
    <transition name="toast-fade">
      <div v-if="toast" class="global-toast">
        {{ toast }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { useI18n } from './composables/useI18n';
import { useAuth } from './composables/useAuth';
import { useGameStore } from './composables/useGameStore';
import LobbyView from './views/LobbyView.vue';
import GameView from './views/GameView.vue';
import RulesModal from './components/RulesModal.vue';
import AdminModal from './components/AdminModal.vue';
import ChangelogModal from './components/ChangelogModal.vue';
import PauseModal from './components/PauseModal.vue';
import DisconnectModal from './components/DisconnectModal.vue';
import LoginModal from './components/LoginModal.vue';

const { t } = useI18n();
const { handleCallback } = useAuth();
const showRules = ref(false);
const showAdmin = ref(false);
const showChangelog = ref(false);
const { connected, toast, globalBroadcast, dismissBroadcast, isScreenShaking, isFatalShotActive } = useGameStore();

watch(isScreenShaking, (val) => {
  if (val) {
    document.body.classList.add('animate-screen-shake');
  } else {
    document.body.classList.remove('animate-screen-shake');
  }
});

function onKeyDown(e) {
  if ((e.ctrlKey || e.metaKey) && (e.key === 'x' || e.key === 'X')) {
    e.preventDefault();
    showAdmin.value = !showAdmin.value;
  }
}

onMounted(async () => {
  window.addEventListener('keydown', onKeyDown);
  // Check if we are handling an OAuth2 callback
  await handleCallback();
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
  top: 16px;
  left: 50%;
  transform: translateX(-50%);
  width: 92%;
  max-width: 820px;
  z-index: 10000;
}

.broadcast-inner {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 20px;
  background: rgba(15, 23, 42, 0.95);
  border: 1px solid rgba(99, 102, 241, 0.4);
  border-radius: var(--radius-lg);
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.8), 0 0 20px rgba(99, 102, 241, 0.2);
  animation: pulse-glow 3s infinite alternate;
}

@keyframes pulse-glow {
  0% {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.8), 0 0 15px rgba(99, 102, 241, 0.15);
  }
  100% {
    box-shadow: 0 12px 36px rgba(0, 0, 0, 0.9), 0 0 25px rgba(99, 102, 241, 0.35);
  }
}

.broadcast-icon {
  font-size: 22px;
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
  color: #818cf8;
  font-weight: 800;
  text-transform: uppercase;
}

.broadcast-msg {
  font-size: 14px;
  color: #f8fafc;
  font-weight: 600;
  line-height: 1.4;
  word-break: break-word;
}

.broadcast-close {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #94a3b8;
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
  background: rgba(239, 68, 68, 0.25);
  border-color: #ef4444;
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
</style>

<template>
  <div class="lobby-view">
    <!-- Top Bar Controls -->
    <div class="lobby-topbar">
      <button class="btn-icon" @click="audio.toggleMute" :title="audio.isMuted.value ? t('audio_off') : t('audio_on')">
        {{ audio.isMuted.value ? 'Muted' : 'Sound' }}
      </button>
      <button class="btn-icon lang-btn" @click="toggleLang">
        {{ lang.toUpperCase() }}
      </button>
    </div>

    <div class="lobby-hero">
      <h1 class="hero-title font-serif">Liar's Deck</h1>
      <p class="hero-subtitle">{{ t('app_subtitle') }}</p>
    </div>

    <!-- Main Card Box -->
    <div class="lobby-card glass-panel">
      <!-- Nickname Input -->
      <div class="form-group">
        <label>{{ t('nickname') }}</label>
        <input
          v-model="nicknameInput"
          type="text"
          :placeholder="t('nickname_ph')"
          maxlength="12"
          @keyup.enter="handleEnter"
        />
      </div>

      <!-- Mode 1: Create -->
      <template v-if="mode === 'create'">
        <div class="mode-actions">
          <button class="btn-primary full-btn" @click="onCreateRoom">
            {{ t('lobby_create_title') }}
          </button>
          <div class="sub-actions">
            <button class="btn-secondary flex-1" @click="mode = 'join'">
              {{ t('join_btn') }}
            </button>
            <button class="btn-secondary flex-1" @click="mode = 'spectate'">
              {{ t('spectate_btn') }}
            </button>
          </div>
        </div>
      </template>

      <!-- Mode 2: Join with Code -->
      <template v-else-if="mode === 'join'">
        <div class="form-group">
          <label>{{ t('room_code') }}</label>
          <input
            v-model="roomCodeInput"
            type="text"
            :placeholder="t('room_code_ph')"
            maxlength="6"
            class="room-code-input"
            @input="roomCodeInput = roomCodeInput.toUpperCase()"
            @keyup.enter="onJoinRoom"
          />
        </div>
        <div class="mode-actions">
          <button class="btn-primary full-btn" @click="onJoinRoom">
            {{ t('join_btn') }}
          </button>
          <button class="btn-secondary full-btn" @click="mode = 'create'">
            ← {{ t('back') }}
          </button>
        </div>
      </template>

      <!-- Mode 3: Spectate with Code -->
      <template v-else-if="mode === 'spectate'">
        <div class="form-group">
          <label>{{ t('room_code') }}</label>
          <input
            v-model="roomCodeInput"
            type="text"
            :placeholder="t('room_code_ph')"
            maxlength="6"
            class="room-code-input"
            @input="roomCodeInput = roomCodeInput.toUpperCase()"
            @keyup.enter="onSpectateRoom"
          />
        </div>
        <div class="mode-actions">
          <button class="btn-primary full-btn" @click="onSpectateRoom">
            {{ t('spectate_btn') }}
          </button>
          <button class="btn-secondary full-btn" @click="mode = 'create'">
            ← {{ t('back') }}
          </button>
        </div>
      </template>
    </div>

    <!-- Rulebook Link -->
    <button class="rules-link" @click="$emit('open-rules')">
      {{ t('rules_btn') }}
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useAudio } from '../composables/useAudio';
import { useGameStore } from '../composables/useGameStore';

const emit = defineEmits(['open-rules']);
const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const store = useGameStore();

const nicknameInput = ref('Player' + Math.floor(Math.random() * 900 + 100));
const roomCodeInput = ref('');
const mode = ref('create'); // 'create' | 'join' | 'spectate'

onMounted(() => {
  const params = new URLSearchParams(window.location.search);
  const room = params.get('room');
  if (room) {
    roomCodeInput.value = room.toUpperCase();
    mode.value = 'join';
  }
});

function handleEnter() {
  if (mode.value === 'create') onCreateRoom();
  else if (mode.value === 'join') onJoinRoom();
  else onSpectateRoom();
}

function onCreateRoom() {
  if (!nicknameInput.value.trim()) {
    store.showToast(t('err_enter_nickname'));
    return;
  }
  store.connect('create', '', nicknameInput.value.trim());
}

function onJoinRoom() {
  if (!nicknameInput.value.trim()) {
    store.showToast(t('err_enter_nickname'));
    return;
  }
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  store.connect('join', roomCodeInput.value.trim().toUpperCase(), nicknameInput.value.trim());
}

function onSpectateRoom() {
  const nick = nicknameInput.value.trim() || 'Spectator' + Math.floor(Math.random() * 100);
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  store.connect('spectate', roomCodeInput.value.trim().toUpperCase(), nick);
}
</script>

<style scoped>
.lobby-view {
  min-height: 85vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px 10px;
  position: relative;
}

.lobby-topbar {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  gap: 8px;
}

.lang-btn {
  font-weight: 700;
}

.lobby-hero {
  text-align: center;
  margin-bottom: 24px;
}

.hero-title {
  font-size: 3rem;
  font-weight: 900;
  color: var(--text-primary);
  letter-spacing: 2px;
  text-shadow: 0 0 24px var(--gold-glow);
  margin-bottom: 4px;
}

.hero-subtitle {
  font-size: 0.95rem;
  color: var(--text-secondary);
}

.lobby-card {
  width: 100%;
  max-width: 390px;
  padding: 28px 24px;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.6);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 18px;
  text-align: left;
}

.form-group label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  letter-spacing: 0.5px;
}

.room-code-input {
  text-transform: uppercase;
}

.mode-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sub-actions {
  display: flex;
  gap: 10px;
}

.full-btn {
  width: 100%;
  padding: 12px;
  font-size: 14px;
}

.flex-1 {
  flex: 1;
}

.rules-link {
  margin-top: 20px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 13px;
  border: none;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 4px;
}
.rules-link:hover {
  color: var(--gold-accent);
}
</style>

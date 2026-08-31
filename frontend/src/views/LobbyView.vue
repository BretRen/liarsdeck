<template>
  <div class="lobby-view">
    <!-- Top Bar Controls -->
    <div class="lobby-topbar">
      <!-- User Profile Pill -->
      <div v-if="isAuthenticated && user" class="user-pill glass-panel">
        <img v-if="avatar" :src="avatar" alt="Avatar" class="user-avatar" />
        <div v-else class="user-avatar-placeholder">👤</div>
        <span class="user-name">{{ nickname }}</span>
        <button class="btn-logout" :title="t('logout_btn')" @click="onLogout">
          {{ t('logout_btn') }}
        </button>
      </div>

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
      <!-- Verified Nickname Badge -->
      <div class="form-group">
        <label>{{ t('nickname') }}</label>
        <div class="verified-nickname-box">
          <span class="verified-icon">🔒</span>
          <span class="verified-name">{{ nickname || 'Player' }}</span>
          <span class="verified-tag">AUTHENTICATED</span>
        </div>
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
import { useAuth } from '../composables/useAuth';
import { useGameStore } from '../composables/useGameStore';

const emit = defineEmits(['open-rules']);
const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const auth = useAuth();
const { user, nickname, avatar, isAuthenticated, logout } = auth;
const store = useGameStore();

const roomCodeInput = ref('');
const mode = ref('create'); // 'create' | 'join' | 'spectate'

onMounted(() => {
  const params = new URLSearchParams(window.location.search);
  const room = params.get('room');
  const saved = store.getSavedSession();

  if (saved) {
    if (!room || saved.roomCode === room.toUpperCase()) {
      roomCodeInput.value = saved.roomCode;
      const playerName = nickname.value || saved.nickname;
      store.connect('reconnect', saved.roomCode, playerName, saved.token);
      return;
    }
  }

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
  const playerName = nickname.value.trim() || 'Player';
  store.clearSession();
  store.connect('create', '', playerName);
}

function onJoinRoom() {
  const playerName = nickname.value.trim() || 'Player';
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  const code = roomCodeInput.value.trim().toUpperCase();
  const saved = store.getSavedSession();
  if (saved && saved.roomCode === code) {
    store.connect('reconnect', code, playerName, saved.token);
  } else {
    store.connect('join', code, playerName);
  }
}

function onSpectateRoom() {
  const playerName = nickname.value.trim() || 'Spectator';
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  store.connect('spectate', roomCodeInput.value.trim().toUpperCase(), playerName);
}

function onLogout() {
  if (confirm(t('logout_confirm'))) {
    store.clearSession();
    logout();
  }
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
  align-items: center;
  gap: 8px;
  z-index: 10;
}

/* User Profile Pill */
.user-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px 4px 6px;
  background: rgba(18, 22, 32, 0.85);
  border: 1px solid var(--border-brass);
  border-radius: 20px;
}

.user-avatar {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid var(--accent-gold);
}

.user-avatar-placeholder {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
}

.user-name {
  font-size: 12.5px;
  font-weight: 700;
  color: var(--text-primary);
  max-width: 110px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-logout {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.btn-logout:hover {
  color: #fca5a5;
  background: rgba(239, 68, 68, 0.15);
}

.lobby-hero {
  text-align: center;
  margin-bottom: 24px;
}

.hero-title {
  font-size: 2.8rem;
  color: var(--text-primary);
  letter-spacing: 2px;
  text-shadow: 0 4px 12px rgba(0, 0, 0, 0.8), 0 0 20px rgba(212, 175, 55, 0.3);
  margin-bottom: 6px;
}

.hero-subtitle {
  font-size: 12.5px;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.lobby-card {
  width: 100%;
  max-width: 360px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  text-align: left;
}

.form-group label {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* Verified Nickname Box */
.verified-nickname-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #090b10;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-md);
  padding: 10px 14px;
  height: 42px;
  box-sizing: border-box;
}

.verified-icon {
  font-size: 13px;
  color: var(--accent-gold);
}

.verified-name {
  flex: 1;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.verified-tag {
  font-size: 9px;
  font-weight: 900;
  background: rgba(212, 175, 55, 0.15);
  border: 1px solid rgba(212, 175, 55, 0.4);
  color: var(--accent-gold);
  padding: 2px 6px;
  border-radius: 4px;
  letter-spacing: 0.5px;
}

.room-code-input {
  text-align: center;
  letter-spacing: 2px;
  font-weight: 700;
  font-family: inherit;
  width: 100%;
  box-sizing: border-box;
  padding: 10px 14px;
  font-size: 14px;
  height: 42px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-brass);
  background: #090b10;
  color: #fff;
}

.room-code-input:focus {
  outline: none;
  border-color: var(--gold-accent);
  box-shadow: 0 0 10px rgba(212, 175, 55, 0.3);
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
}

.flex-1 {
  flex: 1;
}

.rules-link {
  margin-top: 18px;
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 13px;
  cursor: pointer;
  text-decoration: underline;
  transition: color 0.2s ease;
}

.rules-link:hover {
  color: var(--gold-accent);
}
</style>

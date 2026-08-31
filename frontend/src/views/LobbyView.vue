<template>
  <div class="min-h-[85vh] flex flex-col items-center justify-center p-5 relative">
    <!-- Top Bar Controls -->
    <div class="absolute top-3 right-3 flex items-center gap-2 z-10">
      <!-- User Profile Pill -->
      <div v-if="isAuthenticated && user" class="flex items-center gap-2 px-3 py-1.5 bg-slate-900/80 border border-slate-700/60 rounded-full shadow-md backdrop-blur-md">
        <img v-if="avatar" :src="avatar" alt="Avatar" class="w-5 h-5 rounded-full object-cover ring-1 ring-indigo-500/50" />
        <div v-else class="w-5 h-5 rounded-full bg-slate-800 flex items-center justify-center text-[10px]">👤</div>
        <span class="text-xs font-bold text-slate-200 max-w-[120px] truncate">{{ nickname }}</span>
        <button class="btn btn-ghost btn-xs text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 px-1.5 h-6 min-h-0" :title="t('logout_btn')" @click="onLogout">
          {{ t('logout_btn') }}
        </button>
      </div>

      <button class="btn btn-sm btn-ghost bg-slate-900/70 border border-slate-700/60 text-slate-300 hover:bg-slate-800" @click="audio.toggleMute" :title="audio.isMuted.value ? t('audio_off') : t('audio_on')">
        {{ audio.isMuted.value ? 'Muted' : 'Sound' }}
      </button>
      <button class="btn btn-sm btn-ghost bg-slate-900/70 border border-slate-700/60 text-slate-300 hover:bg-slate-800" @click="toggleLang">
        {{ lang.toUpperCase() }}
      </button>
    </div>

    <!-- Hero Title -->
    <div class="text-center mb-8">
      <h1 class="text-5xl md:text-6xl font-black font-serif text-slate-100 tracking-wider drop-shadow-[0_4px_16px_rgba(0,0,0,0.8)] mb-2">
        Liar's Deck
      </h1>
      <p class="text-xs md:text-sm text-slate-400 tracking-wide">{{ t('app_subtitle') }}</p>
    </div>

    <!-- Main Card Box (DaisyUI Card) -->
    <div class="card w-full max-w-sm bg-slate-900/80 border border-slate-800/80 shadow-2xl shadow-black/80 backdrop-blur-xl rounded-2xl p-6 flex flex-col gap-4">
      <!-- Verified Nickname Badge -->
      <div class="flex flex-col gap-1.5 text-left">
        <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">{{ t('nickname') }}</label>
        <div class="flex items-center gap-2 bg-slate-950/80 border border-slate-800 rounded-lg px-3.5 h-11">
          <span class="text-xs text-indigo-400">🔒</span>
          <span class="flex-1 text-sm font-bold text-slate-100 truncate">{{ nickname || 'Player' }}</span>
          <span class="badge badge-sm badge-neutral border-indigo-500/30 text-indigo-300 text-[9px] font-extrabold tracking-wider">
            AUTHENTICATED
          </span>
        </div>
      </div>

      <!-- Mode 1: Create -->
      <template v-if="mode === 'create'">
        <div class="flex flex-col gap-2.5 mt-1">
          <button class="btn btn-primary w-full h-11 font-bold tracking-wide shadow-lg shadow-indigo-600/25 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white transition-all active:scale-[0.98]" @click="onCreateRoom">
            {{ t('lobby_create_title') }}
          </button>
          <div class="flex gap-2">
            <button class="btn btn-neutral flex-1 h-10 min-h-0 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-200 font-semibold" @click="mode = 'join'">
              {{ t('join_btn') }}
            </button>
            <button class="btn btn-neutral flex-1 h-10 min-h-0 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-200 font-semibold" @click="mode = 'spectate'">
              {{ t('spectate_btn') }}
            </button>
          </div>
        </div>
      </template>

      <!-- Mode 2: Join with Code -->
      <template v-else-if="mode === 'join'">
        <div class="flex flex-col gap-1.5 text-left">
          <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">{{ t('room_code') }}</label>
          <input
            v-model="roomCodeInput"
            type="text"
            :placeholder="t('room_code_ph')"
            maxlength="6"
            class="input input-bordered w-full h-11 text-center font-bold tracking-widest text-base bg-slate-950/80 border-slate-800 text-slate-100 focus:border-indigo-500 uppercase"
            @input="roomCodeInput = roomCodeInput.toUpperCase()"
            @keyup.enter="onJoinRoom"
          />
        </div>
        <div class="flex flex-col gap-2.5 mt-1">
          <button class="btn btn-primary w-full h-11 font-bold tracking-wide shadow-lg shadow-indigo-600/25 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" @click="onJoinRoom">
            {{ t('join_btn') }}
          </button>
          <button class="btn btn-neutral w-full h-10 min-h-0 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300 font-semibold" @click="mode = 'create'">
            ← {{ t('back') }}
          </button>
        </div>
      </template>

      <!-- Mode 3: Spectate with Code -->
      <template v-else-if="mode === 'spectate'">
        <div class="flex flex-col gap-1.5 text-left">
          <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">{{ t('room_code') }}</label>
          <input
            v-model="roomCodeInput"
            type="text"
            :placeholder="t('room_code_ph')"
            maxlength="6"
            class="input input-bordered w-full h-11 text-center font-bold tracking-widest text-base bg-slate-950/80 border-slate-800 text-slate-100 focus:border-indigo-500 uppercase"
            @input="roomCodeInput = roomCodeInput.toUpperCase()"
            @keyup.enter="onSpectateRoom"
          />
        </div>
        <div class="flex flex-col gap-2.5 mt-1">
          <button class="btn btn-primary w-full h-11 font-bold tracking-wide shadow-lg shadow-indigo-600/25 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" @click="onSpectateRoom">
            {{ t('spectate_btn') }}
          </button>
          <button class="btn btn-neutral w-full h-10 min-h-0 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300 font-semibold" @click="mode = 'create'">
            ← {{ t('back') }}
          </button>
        </div>
      </template>
    </div>

    <!-- Rulebook Link -->
    <button class="mt-6 text-xs text-slate-400 hover:text-indigo-400 underline underline-offset-4 transition-colors" @click="$emit('open-rules')">
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
  const sub = user.value?.sub || '';

  if (saved && (!sub || saved.token === sub)) {
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
  const sub = user.value?.sub || '';
  store.clearSession();
  store.connect('create', '', playerName, sub);
}

function onJoinRoom() {
  const playerName = nickname.value.trim() || 'Player';
  const sub = user.value?.sub || '';
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  const code = roomCodeInput.value.trim().toUpperCase();
  const saved = store.getSavedSession();
  if (saved && saved.roomCode === code && saved.token === sub) {
    store.connect('reconnect', code, playerName, saved.token);
  } else {
    store.connect('join', code, playerName, sub);
  }
}

function onSpectateRoom() {
  const playerName = nickname.value.trim() || 'Spectator';
  const sub = user.value?.sub || '';
  if (!roomCodeInput.value.trim()) {
    store.showToast(t('err_enter_code'));
    return;
  }
  store.connect('spectate', roomCodeInput.value.trim().toUpperCase(), playerName, sub);
}

function onLogout() {
  if (confirm(t('logout_confirm'))) {
    store.clearSession();
    logout();
  }
}
</script>

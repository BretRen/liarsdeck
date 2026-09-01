<template>
  <div class="min-h-[88vh] flex flex-col items-center justify-center p-4 md:p-6 relative max-w-4xl mx-auto w-full">
    <!-- Top Bar Controls -->
    <div class="w-full flex items-center justify-between gap-2 mb-4 z-10">
      <!-- Changelog & Rules Buttons on Left -->
      <div class="flex items-center gap-2">
        <button
          class="btn btn-sm btn-ghost bg-slate-900/80 border border-slate-700/60 text-slate-300 hover:bg-slate-800 hover:text-indigo-300 font-semibold gap-1.5"
          @click="$emit('open-changelog')"
        >
          <span>📜</span>
          <span>{{ t('changelog_btn') }}</span>
        </button>
        <button
          class="btn btn-sm btn-ghost bg-slate-900/80 border border-slate-700/60 text-slate-300 hover:bg-slate-800 hover:text-indigo-300 font-semibold gap-1.5"
          @click="$emit('open-rules')"
        >
          <span>📖</span>
          <span>{{ t('rules_btn') }}</span>
        </button>
      </div>

      <!-- User Profile & Audio & Lang on Right -->
      <div class="flex items-center gap-2">
        <!-- User Profile Pill -->
        <div v-if="isAuthenticated && user" class="hidden sm:flex items-center gap-2 px-3 py-1 bg-slate-900/80 border border-slate-700/60 rounded-full shadow-md backdrop-blur-md">
          <img v-if="avatar" :src="avatar" alt="Avatar" class="w-5 h-5 rounded-full object-cover ring-1 ring-indigo-500/50" />
          <div v-else class="w-5 h-5 rounded-full bg-slate-800 flex items-center justify-center text-[10px]">👤</div>
          <span class="text-xs font-bold text-slate-200 max-w-[110px] truncate">{{ nickname }}</span>
          <button class="btn btn-ghost btn-xs text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 px-1.5 h-6 min-h-0" :title="t('logout_btn')" @click="onLogout">
            {{ t('logout_btn') }}
          </button>
        </div>

        <!-- Volume Dropdown -->
        <div class="dropdown dropdown-end">
          <button tabindex="0" class="btn btn-sm btn-ghost bg-slate-900/70 border border-slate-700/60 text-slate-300 hover:bg-slate-800 gap-1" :title="t('volume_label')">
            <span>{{ audio.isMuted.value ? '🔇' : '🔊' }}</span>
            <span class="text-[11px] font-mono font-bold">{{ audio.isMuted.value ? '0%' : `${Math.round(audio.masterVolume.value * 100)}%` }}</span>
          </button>
          <div tabindex="0" class="dropdown-content z-[100] menu p-3.5 shadow-2xl bg-slate-900/95 border border-slate-700/80 rounded-2xl w-56 backdrop-blur-xl mt-2 flex flex-col gap-2.5">
            <div class="flex items-center justify-between text-xs font-bold text-slate-300">
              <span>{{ t('volume_label') }}</span>
              <button class="btn btn-xs btn-ghost text-slate-400 hover:text-white" @click="audio.toggleMute">
                {{ audio.isMuted.value ? 'Unmute' : 'Mute' }}
              </button>
            </div>
            <input
              type="range"
              min="0"
              max="1"
              step="0.05"
              :value="audio.isMuted.value ? 0 : audio.masterVolume.value"
              class="range range-xs range-primary"
              @input="onVolumeInput"
            />
          </div>
        </div>

        <button class="btn btn-sm btn-ghost bg-slate-900/70 border border-slate-700/60 text-slate-300 hover:bg-slate-800 font-mono font-bold" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>
    </div>

    <!-- Hero Title -->
    <div class="text-center mb-6">
      <h1 class="text-4xl md:text-5xl font-black font-serif text-slate-100 tracking-wider drop-shadow-[0_4px_16px_rgba(0,0,0,0.8)] mb-1">
        Liar's Deck
      </h1>
      <p class="text-xs md:text-sm text-slate-400 tracking-wide">{{ t('app_subtitle') }}</p>
    </div>

    <!-- Grid: Left Play Box + Right Public Tables -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-5 w-full items-start">
      <!-- Main Action Card Box (5 Cols) -->
      <div class="md:col-span-5 card w-full bg-slate-900/85 border border-slate-800/80 shadow-2xl shadow-black/80 backdrop-blur-xl rounded-2xl p-5 flex flex-col gap-4">
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
              class="input input-bordered w-full h-11 px-3.5 bg-slate-950/80 border-slate-800 text-slate-100 focus:border-indigo-500 text-sm font-mono uppercase tracking-wider placeholder:normal-case placeholder:tracking-normal placeholder:font-sans placeholder:text-slate-500"
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
              class="input input-bordered w-full h-11 px-3.5 bg-slate-950/80 border-slate-800 text-slate-100 focus:border-indigo-500 text-sm font-mono uppercase tracking-wider placeholder:normal-case placeholder:tracking-normal placeholder:font-sans placeholder:text-slate-500"
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

      <!-- Public Tables Browser (7 Cols) -->
      <div class="md:col-span-7 card w-full bg-slate-900/85 border border-slate-800/80 shadow-2xl shadow-black/80 backdrop-blur-xl rounded-2xl p-5 flex flex-col gap-3 min-h-[300px]">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-800 pb-2.5">
          <div class="flex items-center gap-2">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <h2 class="text-sm font-bold font-serif text-slate-100 tracking-wide">{{ t('public_rooms_title') }}</h2>
            <span class="badge badge-sm bg-slate-800 text-slate-400 border-slate-700 font-mono">{{ publicRooms.length }}</span>
          </div>
          <button
            class="btn btn-ghost btn-xs text-slate-400 hover:text-indigo-400 gap-1 font-semibold"
            :disabled="fetchingRooms"
            @click="fetchPublicRooms"
          >
            <span :class="{ 'animate-spin': fetchingRooms }">🔄</span>
            <span>{{ t('refresh_btn') }}</span>
          </button>
        </div>

        <!-- Room List Cards -->
        <div class="flex flex-col gap-2 overflow-y-auto max-h-[340px] pr-1">
          <div
            v-for="rm in publicRooms"
            :key="rm.room_code"
            class="flex items-center justify-between p-3 rounded-xl bg-slate-950/70 border border-slate-800 hover:border-indigo-500/50 hover:bg-slate-900 transition-all gap-2"
          >
            <div class="flex items-center gap-2.5 overflow-hidden">
              <div class="w-9 h-9 rounded-lg bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center font-mono font-bold text-indigo-400 text-xs shrink-0">
                {{ rm.room_code }}
              </div>
              <div class="flex flex-col overflow-hidden text-left">
                <div class="flex items-center gap-1.5">
                  <span class="font-bold text-xs text-slate-100 truncate">{{ rm.host_name }}</span>
                  <span class="badge badge-xs bg-slate-800 text-slate-400 border-slate-700 text-[9px]">{{ t('room_host') }}</span>
                  <span v-if="rm.game_mode === 'items'" class="badge badge-xs bg-slate-800 text-amber-300 border-slate-700 text-[9px]">
                    {{ t('mode_items_title') }}
                  </span>
                  <span v-else class="badge badge-xs bg-slate-800 text-slate-300 border-slate-700 text-[9px]">
                    {{ t('mode_classic_title') }}
                  </span>
                </div>
                <div class="flex items-center gap-2 text-[10px] text-slate-500 font-mono">
                  <span :class="rm.status === 'waiting' ? 'text-emerald-400' : rm.status === 'game_over' ? 'text-slate-500' : 'text-amber-400'">
                    ● {{ rm.status === 'waiting' ? t('table_status_waiting') : rm.status === 'game_over' ? t('status_game_over') : t('table_status_playing') }}
                  </span>
                  <span>{{ rm.player_count }}/{{ rm.max_players }} {{ t('table_players') }}</span>
                </div>
              </div>
            </div>

            <!-- Quick Join -->
            <button
              class="btn btn-xs btn-primary bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white font-bold px-3 shrink-0"
              @click="quickJoinRoom(rm.room_code)"
            >
              {{ t('quick_join_btn') }}
            </button>
          </div>

          <!-- Empty State -->
          <div v-if="publicRooms.length === 0 && !fetchingRooms" class="flex flex-col items-center justify-center py-12 gap-2 text-center text-slate-500">
            <span class="text-2xl opacity-60">🃏</span>
            <p class="text-xs">{{ t('public_rooms_empty') }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Room Modal Dialog -->
    <CreateRoomModal
      v-if="showCreateModal"
      @confirm="handleConfirmCreate"
      @close="showCreateModal = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useAudio } from '../composables/useAudio';
import { useAuth } from '../composables/useAuth';
import { useGameStore } from '../composables/useGameStore';
import CreateRoomModal from '../components/CreateRoomModal.vue';

const emit = defineEmits(['open-rules', 'open-changelog']);
const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const auth = useAuth();
const { user, nickname, avatar, isAuthenticated, logout } = auth;
const store = useGameStore();

const roomCodeInput = ref('');
const mode = ref('create'); // 'create' | 'join' | 'spectate'
const showCreateModal = ref(false);
const publicRooms = ref([]);
const fetchingRooms = ref(false);
let pollInterval = null;

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

  fetchPublicRooms();
  pollInterval = setInterval(fetchPublicRooms, 5000);
});

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval);
});

function onVolumeInput(e) {
  audio.setVolume(e.target.value);
}

async function fetchPublicRooms() {
  fetchingRooms.value = true;
  try {
    const res = await fetch('/api/rooms');
    if (res.ok) {
      const data = await res.json();
      if (data.success && Array.isArray(data.rooms)) {
        publicRooms.value = data.rooms;
      }
    }
  } catch (_) {}
  finally {
    fetchingRooms.value = false;
  }
}

function quickJoinRoom(code) {
  roomCodeInput.value = code;
  onJoinRoom();
}

function onCreateRoom() {
  showCreateModal.value = true;
}

function handleConfirmCreate({ mode: roomMode, maxPlayers }) {
  showCreateModal.value = false;
  const playerName = nickname.value.trim() || 'Player';
  const sub = user.value?.sub || '';
  store.clearSession();
  store.connect('create', '', playerName, sub, { mode: roomMode, maxPlayers });
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

<template>
  <header class="flex items-center justify-between p-3.5 mb-3.5 bg-slate-900/80 border border-slate-700/60 rounded-2xl shadow-xl shadow-black/50 backdrop-blur-xl flex-wrap gap-2.5">
    <!-- Left: Room Code & Invite -->
    <div class="flex items-center gap-2">
      <div class="flex items-center bg-slate-950/80 border border-slate-750 rounded-lg px-2.5 py-1 gap-1.5 shadow-inner">
        <span class="text-[10px] font-extrabold tracking-widest text-slate-400">ROOM</span>
        <span class="text-xs font-mono font-bold tracking-wider text-indigo-300">{{ roomCode }}</span>
      </div>
      <button class="btn btn-xs btn-neutral bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300 font-semibold" @click="copyInvite" :title="t('invite_btn')">
        {{ t('invite_btn') }}
      </button>
    </div>

    <!-- Center: Match Status & Turn Countdown Timer -->
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2 px-3 py-1 bg-slate-950/70 border border-slate-800 rounded-full text-xs font-semibold" :class="statusBadgeClass">
        <span class="w-2 h-2 rounded-full" :class="statusDotClass"></span>
        <span>{{ statusLabel }}</span>
      </div>

      <!-- Prominent Countdown Timer -->
      <div
        v-if="status === 'playing' && deadline"
        class="flex items-baseline gap-0.5 px-3 py-1 rounded-lg border font-mono font-bold transition-all shadow-md"
        :class="remainingSeconds <= 8 ? 'bg-rose-950/70 border-rose-500/80 text-rose-400 animate-pulse' : 'bg-slate-950/80 border-indigo-500/40 text-indigo-300'"
      >
        <span class="text-base">{{ remainingSeconds }}</span>
        <span class="text-xs text-slate-400">s</span>
      </div>
    </div>

    <!-- Right: Utility Controls -->
    <div class="flex items-center gap-1.5">
      <button class="btn btn-xs btn-ghost bg-slate-800/80 hover:bg-slate-700 text-slate-300" @click="audio.toggleMute" :title="audio.isMuted.value ? t('audio_off') : t('audio_on')">
        {{ audio.isMuted.value ? 'Muted' : 'Sound' }}
      </button>
      <button class="btn btn-xs btn-ghost bg-slate-800/80 hover:bg-slate-700 text-slate-300" @click="$emit('open-rules')">
        {{ t('rules_btn') }}
      </button>
      <button class="btn btn-xs btn-ghost bg-slate-800/80 hover:bg-slate-700 text-slate-300" @click="toggleLang">
        {{ lang.toUpperCase() }}
      </button>
      <button class="btn btn-xs btn-error bg-rose-600/80 hover:bg-rose-600 border-none text-white font-semibold" @click="$emit('leave')" title="离开房间">
        退出
      </button>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useAudio } from '../composables/useAudio';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  roomCode: { type: String, default: '' },
  status: { type: String, default: 'waiting' },
  deadline: { type: Number, default: 0 },
});

defineEmits(['open-rules', 'leave']);

const { t, lang, toggleLang } = useI18n();
const audio = useAudio();
const store = useGameStore();

const remainingSeconds = computed(() => store.remainingSeconds.value);

const statusLabel = computed(() => {
  if (props.status === 'playing') return t('status_playing');
  if (props.status === 'paused') return t('status_paused');
  if (props.status === 'game_over') return t('status_game_over');
  return t('status_waiting');
});

const statusBadgeClass = computed(() => {
  if (props.status === 'playing') return 'text-emerald-400 border-emerald-500/30';
  if (props.status === 'paused') return 'text-amber-400 border-amber-500/30';
  if (props.status === 'game_over') return 'text-indigo-400 border-indigo-500/30';
  return 'text-slate-300 border-slate-700';
});

const statusDotClass = computed(() => {
  if (props.status === 'playing') return 'bg-emerald-400 animate-pulse';
  if (props.status === 'paused') return 'bg-amber-400 animate-ping';
  if (props.status === 'game_over') return 'bg-indigo-400';
  return 'bg-slate-400';
});

function copyInvite() {
  store.copyInvite();
}
</script>

<template>
  <div v-if="isDisconnected" class="fixed inset-0 z-[3000] flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-in fade-in duration-200">
    <div class="card w-full max-w-md bg-slate-900/95 border border-rose-500/40 rounded-2xl shadow-2xl shadow-black/90 overflow-hidden">
      <!-- Modal Header with Language Toggle -->
      <div class="flex items-center justify-between px-6 py-4 bg-slate-950/80 border-b border-slate-800">
        <h2 class="text-base font-bold font-serif text-rose-400">{{ t('dc_title') }}</h2>
        <button class="btn btn-ghost btn-xs text-slate-400 hover:text-white" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 flex flex-col items-center text-center gap-4">
        <div class="relative w-16 h-16 rounded-full bg-rose-950/40 border border-rose-500/30 flex items-center justify-center">
          <div class="absolute inset-0 rounded-full border border-rose-400/40 animate-ping"></div>
          <span class="text-2xl">⚡</span>
        </div>

        <p class="text-xs text-slate-300 leading-relaxed">{{ t('dc_desc') }}</p>

        <!-- 30s Grace Countdown when match is in progress -->
        <div v-if="state.status === 'paused' && pauseRemainingSeconds > 0" class="w-full p-3.5 bg-slate-950/80 border border-slate-800 rounded-xl flex flex-col items-center gap-1">
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">{{ t('pause_countdown_label') }}</span>
          <span class="text-2xl font-mono font-black text-rose-400">{{ pauseRemainingSeconds }}s</span>
          <p class="text-[11px] text-slate-400">{{ t('dc_grace_tip') }}</p>
        </div>

        <div class="flex items-center gap-2 text-xs font-semibold text-sky-400">
          <span class="w-2 h-2 rounded-full bg-sky-400 animate-ping"></span>
          <span>{{ t('dc_reconnecting') }}</span>
        </div>
      </div>

      <!-- Modal Actions -->
      <div class="flex gap-2 p-4 bg-slate-950/70 border-t border-slate-800">
        <button class="btn btn-primary flex-1 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" :disabled="isReconnecting" @click="tryReconnect">
          {{ t('dc_retry_btn') }}
        </button>
        <button class="btn btn-neutral flex-1 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300 font-semibold" @click="exitToLobby">
          {{ t('dc_exit_btn') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const { t, lang, toggleLang } = useI18n();
const { state, isDisconnected, isReconnecting, pauseRemainingSeconds, tryReconnect, exitToLobby } =
  useGameStore();
</script>

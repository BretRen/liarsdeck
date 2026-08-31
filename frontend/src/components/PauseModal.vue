<template>
  <div v-if="isPaused && !isDisconnected" class="fixed inset-0 z-[2500] flex items-center justify-center p-4 bg-slate-950/85 backdrop-blur-md animate-in fade-in duration-200">
    <div class="card w-full max-w-md bg-slate-900/95 border border-amber-500/40 rounded-2xl shadow-2xl shadow-black/90 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 bg-slate-950/80 border-b border-slate-800">
        <div class="flex items-center gap-2">
          <span class="text-amber-400 text-lg">⏸</span>
          <h2 class="text-lg font-bold font-serif text-slate-100">{{ t('pause_modal_title') }}</h2>
        </div>
        <button class="btn btn-ghost btn-xs text-slate-400 hover:text-white" @click="toggleLang">
          {{ lang.toUpperCase() }}
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 flex flex-col items-center text-center gap-5">
        <div class="text-sm text-slate-300">
          {{ t('pause_modal_desc', { name: state.paused_player }) }}
        </div>

        <!-- 30s Countdown Display -->
        <div class="w-full bg-slate-950/80 border border-slate-800 rounded-xl p-5 flex flex-col items-center gap-2" :class="{ 'border-rose-500/60 shadow-lg shadow-rose-950/40': pauseRemainingSeconds <= 10 }">
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">{{ t('pause_countdown_label') }}</span>
          <div class="flex items-baseline gap-1 font-mono font-black" :class="pauseRemainingSeconds <= 10 ? 'text-rose-400' : 'text-amber-400'">
            <span class="text-4xl leading-none">{{ pauseRemainingSeconds }}</span>
            <span class="text-sm text-slate-500">s</span>
          </div>

          <!-- Progress Bar -->
          <div class="w-full h-2 bg-slate-900 rounded-full overflow-hidden mt-1 border border-slate-800">
            <div
              class="h-full rounded-full transition-all duration-300"
              :class="pauseRemainingSeconds <= 10 ? 'bg-rose-500' : 'bg-amber-400'"
              :style="{ width: `${(Math.max(0, pauseRemainingSeconds) / 30) * 100}%` }"
            ></div>
          </div>
        </div>

        <p class="text-xs text-slate-400 leading-relaxed max-w-xs">
          {{ t('pause_tip') }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const { t, lang, toggleLang } = useI18n();
const { state, isDisconnected, pauseRemainingSeconds } = useGameStore();

const isPaused = computed(() => state.value.status === 'paused');
</script>

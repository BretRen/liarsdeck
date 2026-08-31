<template>
  <div
    class="relative flex-1 min-w-[125px] bg-slate-900/85 border rounded-xl p-2.5 flex flex-col gap-1.5 transition-all duration-200 backdrop-blur-md"
    :class="[
      isActiveTurn && !player.is_spectator ? 'border-indigo-500 bg-slate-900 shadow-lg shadow-indigo-500/25 -translate-y-0.5' : '',
      isMe ? 'ring-1 ring-sky-500/50 border-sky-500/40' : 'border-slate-800',
      !player.is_alive && !player.is_spectator ? 'opacity-40 grayscale border-rose-950/60 bg-slate-950' : '',
      player.is_spectator ? 'border-dashed border-slate-800 bg-slate-950/60' : ''
    ]"
  >
    <!-- Turn Indicator & Countdown Ring -->
    <div v-if="isActiveTurn && !player.is_spectator" class="absolute -top-2.5 left-1/2 -translate-x-1/2 bg-indigo-600 text-white font-extrabold text-[9px] tracking-wider px-2 py-0.5 rounded shadow-md shadow-indigo-600/40 flex items-center gap-1.5 whitespace-nowrap z-10">
      <span>TURN</span>
      <span v-if="remainingSeconds > 0" class="bg-slate-950 text-indigo-200 px-1 rounded text-[8.5px]">{{ remainingSeconds }}s</span>
    </div>

    <!-- Header: Name & Role Tags -->
    <div class="flex items-center justify-between gap-1">
      <div class="flex items-center gap-1 overflow-hidden">
        <span class="font-bold text-xs md:text-sm text-slate-100 max-w-[85px] truncate" :title="player.nickname">{{ player.nickname }}</span>
        <span v-if="player.is_host" class="badge badge-xs bg-indigo-500/20 text-indigo-300 border-indigo-500/40 text-[9px] font-bold">{{ t('host_tag') }}</span>
        <span v-if="player.is_spectator" class="badge badge-xs bg-slate-800 text-slate-400 border-slate-700 text-[9px]">{{ t('spec_tag') }}</span>
        <span v-if="player.is_alive && !player.is_spectator && player.is_connected === false" class="badge badge-xs bg-rose-950 text-rose-300 border-rose-700 text-[9px] animate-pulse">{{ t('offline_tag') }}</span>
        <span v-if="isMe" class="text-[10px] font-bold text-sky-400">[{{ t('me_tag') }}]</span>
      </div>

      <!-- Kick Button for Host -->
      <button
        v-if="amHost && !player.is_host && !isMe"
        class="btn btn-ghost btn-xs text-slate-500 hover:text-rose-400 hover:bg-rose-950/40 h-5 w-5 min-h-0 p-0"
        :title="t('kick_btn')"
        @click="$emit('kick', player.id)"
      >
        ✕
      </button>
    </div>

    <!-- Body for Non-Spectators -->
    <div v-if="!player.is_spectator" class="flex flex-col gap-1.5">
      <!-- 6-Chamber Revolver Visual Indicator -->
      <div class="flex gap-1 justify-center py-1 bg-slate-950/80 rounded-md border border-slate-800/80" :title="`${t('bullets_label')}: ${player.bullets}/6`">
        <div
          v-for="i in 6"
          :key="i"
          class="w-2.5 h-3.5 rounded-xs transition-all border"
          :class="[
            player.is_alive && i <= player.bullets ? 'bg-sky-500/80 border-sky-400 shadow-[0_0_6px_rgba(56,189,248,0.5)]' : '',
            player.is_alive && i > player.bullets ? 'bg-slate-900 border-slate-800' : '',
            !player.is_alive ? 'bg-rose-950 border-rose-900' : ''
          ]"
        ></div>
      </div>

      <div class="flex items-center justify-between text-xs font-semibold px-0.5">
        <div v-if="player.is_alive" class="flex items-center gap-1 text-slate-400">
          <span class="text-[11px]">{{ t('hand_count_label') }}:</span>
          <span class="text-slate-200 font-bold font-mono">{{ player.hand ? player.hand.length : 0 }}</span>
        </div>
        <div v-else class="text-rose-400 font-bold text-[11px]">
          {{ t('dead_tag') }}
        </div>
      </div>

      <!-- Ready Status during Waiting Phase -->
      <div v-if="gameStatus === 'waiting'" class="text-center py-0.5 rounded text-[10.5px] font-bold transition-colors" :class="player.is_ready ? 'bg-emerald-950/60 border border-emerald-500/40 text-emerald-400' : 'bg-slate-950/60 border border-slate-800 text-slate-500'">
        <span>{{ player.is_ready ? t('ready_status') : t('unready_status') }}</span>
      </div>
    </div>

    <!-- Spectator Body -->
    <div v-else class="py-2 text-center text-xs text-slate-500 font-medium">
      <span>{{ t('watching_tag') }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  player: { type: Object, required: true },
  isActiveTurn: { type: Boolean, default: false },
  isMe: { type: Boolean, default: false },
  amHost: { type: Boolean, default: false },
  gameStatus: { type: String, default: 'waiting' },
});

defineEmits(['kick']);
const { t } = useI18n();
const store = useGameStore();

const remainingSeconds = computed(() => store.remainingSeconds.value);
</script>

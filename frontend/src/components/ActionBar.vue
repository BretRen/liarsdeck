<template>
  <div class="p-3 mb-3.5 bg-slate-900/80 border border-slate-700/60 rounded-2xl shadow-xl shadow-black/50 backdrop-blur-xl flex items-center justify-center min-h-[64px]">
    <!-- Waiting Phase -->
    <template v-if="computedStatus === 'waiting'">
      <div v-if="isPlayer" class="flex items-center gap-3 flex-wrap justify-center">
        <button class="btn btn-sm md:btn-md btn-primary px-6 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" @click="$emit('toggle-ready')">
          {{ isReady ? t('unready_btn') : t('ready_btn') }}
        </button>

        <button
          class="btn btn-sm md:btn-md btn-success px-6 font-bold text-white shadow-lg shadow-emerald-600/30 bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-500 hover:to-emerald-600 border-none disabled:opacity-40"
          :disabled="!canStart"
          @click="$emit('start-game')"
          v-if="amHost"
        >
          {{ t('start_game_btn') }}
          <span v-if="!canStart" class="text-[10px] font-normal opacity-80">
            ({{ allReady ? t('need_more_players') : t('all_ready_needed') }})
          </span>
        </button>
      </div>

      <div v-else class="text-xs text-slate-400 font-medium">
        <span>{{ t('spectator_banner') }}</span>
      </div>
    </template>

    <!-- Playing Phase -->
    <template v-else-if="computedStatus === 'playing'">
      <div v-if="isPlayer && isAlive && isMyTurn" class="flex items-center gap-3 flex-wrap justify-center">
        <!-- Turn Timer Badge inside Action Bar -->
        <div v-if="remainingSeconds > 0" class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border font-mono font-bold text-xs" :class="remainingSeconds <= 8 ? 'bg-rose-950/80 border-rose-500 text-rose-300 animate-pulse' : 'bg-slate-950/80 border-indigo-500/40 text-indigo-300'">
          <span class="text-[11px] uppercase tracking-wider text-slate-400">{{ t('timeout_warn') }}:</span>
          <span class="text-sm font-extrabold">{{ remainingSeconds }}s</span>
        </div>

        <button
          class="btn btn-sm md:btn-md btn-primary px-6 font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white disabled:opacity-40"
          :disabled="!canPlay"
          @click="$emit('play-cards')"
        >
          {{ t('play_cards_btn') }} ({{ selectedCount }}/3)
        </button>

        <button
          class="btn btn-sm md:btn-md btn-error px-6 font-bold shadow-lg shadow-rose-600/30 bg-gradient-to-r from-rose-600 to-rose-700 hover:from-rose-500 hover:to-rose-600 border-none text-white disabled:opacity-40"
          :class="{ 'animate-bounce': canCallLiar }"
          :disabled="!canCallLiar"
          @click="$emit('call-liar')"
        >
          {{ t('call_liar_btn') }}
        </button>
      </div>

      <!-- waiting for turn OR dead player watching -->
      <div v-else class="flex items-center gap-2 text-xs text-slate-400 font-medium">
        <span class="w-2 h-2 rounded-full bg-indigo-400 animate-pulse"></span>
        <span v-if="isPlayer && !isAlive">已淘汰 — 观战中</span>
        <span v-else-if="isPlayer">{{ t('status_playing') }} ({{ t('status_waiting') }})</span>
        <span v-else>{{ t('spectator_banner') }}</span>
      </div>
    </template>

    <!-- Game Over Phase -->
    <template v-else-if="computedStatus === 'game_over'">
      <div class="text-xs text-slate-400 font-medium">对局已结束</div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  status: { type: String, default: '' },
  gameStatus: { type: String, default: '' },
  isPlayer: { type: Boolean, default: true },
  isAlive: { type: Boolean, default: true },
  isMyTurn: { type: Boolean, default: false },
  isReady: { type: Boolean, default: false },
  amHost: { type: Boolean, default: false },
  canPlay: { type: Boolean, default: false },
  canCallLiar: { type: Boolean, default: false },
  canStart: { type: Boolean, default: false },
  allReady: { type: Boolean, default: false },
  selectedCount: { type: Number, default: 0 },
});

defineEmits(['toggle-ready', 'start-game', 'play-cards', 'call-liar']);
const { t } = useI18n();
const store = useGameStore();

const computedStatus = computed(() => props.status || props.gameStatus || 'waiting');
const remainingSeconds = computed(() => store.remainingSeconds.value);
</script>

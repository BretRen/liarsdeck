<template>
  <div class="flex flex-col items-center justify-center p-5 mb-3.5 bg-slate-900/80 border border-slate-700/60 rounded-2xl shadow-2xl shadow-black/60 backdrop-blur-xl min-h-[190px] gap-3.5">
    <!-- Playing State: Table True Card & Face-down Stack -->
    <template v-if="status === 'playing'">
      <div class="flex flex-col items-center gap-1.5">
        <div class="font-serif text-[11px] font-bold tracking-widest text-indigo-400 uppercase">
          {{ t('table_card_label') }}
        </div>
        
        <div class="p-1">
          <div class="w-20 h-28 bg-white border-2 border-indigo-400 rounded-lg shadow-xl shadow-indigo-500/15 flex flex-col justify-between p-2 font-black relative select-none" :class="isRedCard ? 'text-rose-600' : 'text-slate-900'">
            <span class="text-xs text-left leading-none font-serif">{{ tableCard }}</span>
            <span class="text-4xl text-center leading-none font-serif">{{ tableCard }}</span>
            <span class="text-xs text-right leading-none font-serif">{{ tableCard }}</span>
            <div v-if="tableCard === '2'" class="absolute inset-x-0 bottom-1 mx-auto w-max px-1.5 py-0.5 bg-indigo-600 text-white rounded text-[8px] font-extrabold tracking-widest">
              WILD
            </div>
          </div>
        </div>

        <div class="text-xs text-slate-400 font-medium">{{ t('wild_card_tip') }}</div>
      </div>

      <!-- Face-down Played Pile -->
      <div v-if="lastPlayedCnt > 0" class="flex flex-col items-center gap-1 mt-1">
        <div class="flex items-center justify-center h-12">
          <div
            v-for="i in lastPlayedCnt"
            :key="i"
            class="w-9 h-13 bg-indigo-950 border border-indigo-500/40 rounded-sm flex items-center justify-center shadow-lg shadow-black/60 -mx-2 transition-transform p-0.5"
            :style="{ transform: `rotate(${(i - (lastPlayedCnt + 1) / 2) * 6}deg)` }"
          >
            <div class="w-full h-full border border-dashed border-indigo-400/30 rounded-xs bg-[repeating-linear-gradient(45deg,#1e1b4b,#1e1b4b_2px,#312e81_2px,#312e81_4px)]"></div>
          </div>
        </div>
        <div class="text-xs font-semibold text-slate-300">
          {{ lastPlayedCnt }} {{ t('cards_on_table') }}
        </div>
      </div>
    </template>

    <!-- Game Over State: Champion Display -->
    <template v-else-if="status === 'game_over'">
      <div class="flex flex-col items-center text-center gap-2 py-3">
        <div class="text-xs font-extrabold text-slate-400 tracking-wider uppercase">{{ t('winner_label') }}</div>
        <div class="text-4xl font-black font-serif text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-sky-300 to-indigo-400 drop-shadow-[0_0_25px_rgba(99,102,241,0.5)]">
          {{ winner || '—' }}
        </div>

        <div class="mt-3">
          <button v-if="amHost" class="btn btn-primary px-8 h-11 text-sm font-bold shadow-lg shadow-indigo-600/30 bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 border-none text-white" @click="$emit('reset')">
            {{ t('play_again_btn') }}
          </button>
          <p v-else class="text-xs text-slate-400">
            {{ t('wait_host_reset') }}
          </p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';

const props = defineProps({
  status: { type: String, default: 'waiting' },
  tableCard: { type: String, default: '' },
  lastPlayedCnt: { type: Number, default: 0 },
  winner: { type: String, default: '' },
  amHost: { type: Boolean, default: false },
});

defineEmits(['reset']);
const { t } = useI18n();

const isRedCard = computed(() => props.tableCard === 'Q' || props.tableCard === '2');
</script>

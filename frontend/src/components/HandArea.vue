<template>
  <div class="p-3.5 mb-3.5 bg-slate-900/80 border border-slate-700/60 rounded-2xl shadow-xl shadow-black/50 backdrop-blur-xl flex flex-col gap-2.5">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-baseline gap-2">
        <h3 class="text-sm font-bold text-slate-100 font-serif">{{ t('my_hand_title') }}</h3>
        <span class="text-xs font-bold text-indigo-400 font-mono">({{ hand.length }})</span>
      </div>

      <div class="flex items-center gap-1.5 text-xs">
        <span class="font-semibold px-2 py-0.5 rounded bg-slate-950/80 border border-slate-800 text-slate-300">
          {{ t('selected_count') }} {{ selectedIndexes.length }}/3
        </span>
        <span v-if="!isMyTurn" class="text-slate-500 font-medium">({{ t('status_waiting') }})</span>
      </div>
    </div>

    <!-- Hand Cards Row -->
    <div class="flex items-center justify-center gap-2.5 min-h-[115px] overflow-x-auto py-2 px-1" :class="{ 'opacity-70': !isMyTurn }">
      <div
        v-for="(card, index) in hand"
        :key="index"
        class="playing-card w-[72px] h-[104px] bg-white border border-slate-300 rounded-lg shadow-lg flex flex-col justify-between p-1.5 font-serif font-black select-none transition-all duration-200 cursor-pointer relative shrink-0"
        :class="[
          card === 'Q' || card === '2' ? 'text-rose-600' : 'text-slate-900',
          selectedIndexes.includes(index) ? '-translate-y-3 shadow-xl shadow-indigo-500/30 ring-2 ring-indigo-500 z-10 scale-105' : 'hover:-translate-y-1',
          !isMyTurn ? 'cursor-not-allowed opacity-80' : ''
        ]"
        @click="onCardClick(index)"
      >
        <span class="text-[11px] text-left leading-none">{{ card }}</span>
        <span class="text-3xl text-center leading-none">{{ card }}</span>
        <span class="text-[11px] text-right leading-none">{{ card }}</span>
        <div v-if="card === '2'" class="absolute inset-x-0 bottom-1 mx-auto w-max px-1 py-0.2 bg-indigo-600 text-white rounded text-[7px] font-extrabold tracking-widest">
          WILD
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';

const props = defineProps({
  hand: { type: Array, default: () => [] },
  selectedIndexes: { type: Array, default: () => [] },
  isMyTurn: { type: Boolean, default: false },
});

const emit = defineEmits(['toggle-select']);
const { t } = useI18n();

function onCardClick(index) {
  if (props.isMyTurn) {
    emit('toggle-select', index);
  }
}
</script>

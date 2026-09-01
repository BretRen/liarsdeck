<template>
  <div v-if="currentStep" class="fixed inset-0 z-[1500] flex items-center justify-center p-5 backdrop-blur-md animate-in fade-in duration-200" :class="overlayBg">
    <!-- Muzzle Flash & Blood Vignette on fatal shot -->
    <div
      v-if="currentStep === 'shot' && stepData.fatal"
      class="fixed inset-0 pointer-events-none z-[1501] bg-[radial-gradient(ellipse_at_center,_transparent_30%,_rgba(225,29,72,0.65)_100%)] animate-pulse"
    ></div>

    <div
      class="card w-full max-w-md p-7 bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col items-center text-center gap-4 transition-all"
      :class="{ 'animate-screen-shake': currentStep === 'shot' && stepData.fatal }"
    >
      <!-- 1. Liar Call Step -->
      <template v-if="currentStep === 'liar_call'">
        <div class="badge badge-error badge-sm font-extrabold tracking-widest uppercase py-1 px-3 shadow-lg shadow-rose-600/30 animate-pulse">
          CHALLENGE
        </div>
        <h2 class="text-2xl font-black font-serif text-slate-100 tracking-wide">{{ t('event_liar_alert') }}</h2>
        <div class="text-sm md:text-base text-slate-300 flex items-center justify-center gap-1.5 flex-wrap">
          <span class="font-bold text-sky-400">{{ stepData.caller }}</span>
          <span>{{ t('event_calls_out') }}</span>
          <span class="font-bold text-rose-400">{{ stepData.accused }}</span>
          <span>{{ t('event_liar_claim') }}</span>
        </div>
      </template>

      <!-- 2. Reveal Cards Step (3D Flip Animation) -->
      <template v-if="currentStep === 'reveal'">
        <div class="badge badge-info badge-sm font-extrabold tracking-widest uppercase py-1 px-3 shadow-lg shadow-sky-600/30">
          VERIFICATION
        </div>
        <h2 class="text-xl font-bold font-serif text-slate-100 tracking-wide">
          {{ stepData.accused }} - {{ t('event_cards_revealed') }}
        </h2>
        <div class="flex gap-3 justify-center mt-2 flex-wrap perspective-1000">
          <div
            v-for="(card, i) in stepData.cards"
            :key="i"
            class="playing-card w-[72px] h-[102px] bg-white border border-slate-300 rounded-xl shadow-2xl flex flex-col justify-between p-2 font-serif font-black select-none preserve-3d animate-flip-in"
            :class="[
              card === 'Q' || card === '2' ? 'text-rose-600' : 'text-slate-900',
              isCardHonest(card) ? 'shadow-sky-500/40 ring-2 ring-sky-400/80' : 'shadow-rose-600/50 ring-2 ring-rose-500/90'
            ]"
            :style="{ animationDelay: `${i * 0.16}s` }"
          >
            <span class="text-xs text-left leading-none font-sans font-extrabold">{{ card }}</span>
            <span class="text-3xl text-center leading-none">{{ card }}</span>
            <span class="text-xs text-right leading-none font-sans font-extrabold">{{ card }}</span>
            <div v-if="card === '2'" class="absolute inset-x-0 bottom-1.5 mx-auto w-max px-1.5 bg-indigo-600 text-white rounded text-[7.5px] font-extrabold tracking-widest shadow">
              WILD
            </div>
          </div>
        </div>
      </template>

      <!-- 3. Shot Step -->
      <template v-if="currentStep === 'shot'">
        <div class="badge badge-sm font-extrabold tracking-widest uppercase py-1 px-3" :class="stepData.fatal ? 'badge-error text-white shadow-lg shadow-rose-600/40' : 'badge-success text-white shadow-lg shadow-emerald-600/40'">
          {{ stepData.fatal ? 'FATAL ROUND' : 'DRY FIRE' }}
        </div>
        <h2
          class="text-2xl font-black font-serif tracking-wide"
          :class="stepData.fatal ? 'text-rose-400 drop-shadow-[0_0_20px_rgba(244,63,94,0.8)] animate-pulse' : 'text-emerald-400 drop-shadow-[0_0_15px_rgba(16,185,129,0.5)]'"
        >
          {{ stepData.fatal ? t('event_bang_title') : t('event_click_title') }}
        </h2>
        <p class="text-sm text-slate-300 leading-relaxed">
          <span class="font-bold text-slate-100">{{ stepData.target }}</span>
          {{ stepData.fatal ? t('event_bang_sub') : t('event_click_sub') }}
        </p>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  currentStep: { type: String, default: '' },
  stepData: { type: Object, default: () => ({}) },
});

const { t } = useI18n();
const store = useGameStore();

function isCardHonest(card) {
  const tableCard = store.state.value.table_card;
  if (!tableCard) return true;
  return card === tableCard || card === '2';
}

const overlayBg = computed(() => {
  if (props.currentStep === 'liar_call') return 'bg-rose-950/70';
  if (props.currentStep === 'reveal') return 'bg-slate-950/80';
  if (props.currentStep === 'shot') return props.stepData.fatal ? 'bg-rose-950/85' : 'bg-slate-950/85';
  return 'bg-slate-950/75';
});
</script>

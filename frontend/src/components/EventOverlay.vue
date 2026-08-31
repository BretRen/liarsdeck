<template>
  <div v-if="currentStep" class="fixed inset-0 z-[1500] flex items-center justify-center p-5 backdrop-blur-md animate-in fade-in duration-200" :class="overlayBg">
    <!-- Muzzle Flash on fatal shot -->
    <div v-if="currentStep === 'shot' && stepData.fatal" class="fixed inset-0 bg-rose-600/30 z-[1501] pointer-events-none animate-ping"></div>

    <div class="card w-full max-w-md p-7 bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col items-center text-center gap-3.5" :class="{ 'animate-bounce': currentStep === 'shot' && stepData.fatal }">
      <!-- 1. Liar Call Step -->
      <template v-if="currentStep === 'liar_call'">
        <div class="badge badge-error badge-sm font-extrabold tracking-widest uppercase py-1 px-3">CHALLENGE</div>
        <h2 class="text-2xl font-black font-serif text-slate-100 tracking-wide">{{ t('event_liar_alert') }}</h2>
        <div class="text-sm md:text-base text-slate-300 flex items-center justify-center gap-1.5 flex-wrap">
          <span class="font-bold text-sky-400">{{ stepData.caller }}</span>
          <span>{{ t('event_calls_out') }}</span>
          <span class="font-bold text-rose-400">{{ stepData.accused }}</span>
          <span>{{ t('event_liar_claim') }}</span>
        </div>
      </template>

      <!-- 2. Reveal Cards Step -->
      <template v-if="currentStep === 'reveal'">
        <div class="badge badge-info badge-sm font-extrabold tracking-widest uppercase py-1 px-3">VERIFICATION</div>
        <h2 class="text-xl font-bold font-serif text-slate-100 tracking-wide">{{ stepData.accused }} - {{ t('event_cards_revealed') }}</h2>
        <div class="flex gap-2.5 justify-center mt-2 flex-wrap">
          <div
            v-for="(card, i) in stepData.cards"
            :key="i"
            class="playing-card w-[68px] h-[98px] bg-white border border-slate-300 rounded-lg shadow-xl flex flex-col justify-between p-1.5 font-serif font-black select-none animate-in zoom-in spin-in-12 duration-300"
            :class="card === 'Q' || card === '2' ? 'text-rose-600' : 'text-slate-900'"
            :style="{ animationDelay: `${i * 0.12}s` }"
          >
            <span class="text-[10px] text-left leading-none">{{ card }}</span>
            <span class="text-3xl text-center leading-none">{{ card }}</span>
            <span class="text-[10px] text-right leading-none">{{ card }}</span>
            <div v-if="card === '2'" class="absolute inset-x-0 bottom-1 mx-auto w-max px-1 bg-indigo-600 text-white rounded text-[7px] font-extrabold tracking-widest">
              WILD
            </div>
          </div>
        </div>
      </template>

      <!-- 3. Shot Step -->
      <template v-if="currentStep === 'shot'">
        <div class="badge badge-sm font-extrabold tracking-widest uppercase py-1 px-3" :class="stepData.fatal ? 'badge-error text-white' : 'badge-success text-white'">
          {{ stepData.fatal ? 'FATAL ROUND' : 'DRY FIRE' }}
        </div>
        <h2 class="text-2xl font-black font-serif tracking-wide" :class="stepData.fatal ? 'text-rose-400 drop-shadow-[0_0_15px_rgba(244,63,94,0.6)]' : 'text-emerald-400 drop-shadow-[0_0_15px_rgba(16,185,129,0.5)]'">
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

const props = defineProps({
  currentStep: { type: String, default: '' },
  stepData: { type: Object, default: () => ({}) },
});

const { t } = useI18n();

const overlayBg = computed(() => {
  if (props.currentStep === 'liar_call') return 'bg-rose-950/70';
  if (props.currentStep === 'reveal') return 'bg-slate-950/80';
  if (props.currentStep === 'shot') return props.stepData.fatal ? 'bg-rose-950/85' : 'bg-slate-950/85';
  return 'bg-slate-950/75';
});
</script>

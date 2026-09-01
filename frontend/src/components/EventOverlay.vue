<template>
  <div v-if="currentStep" class="fixed inset-0 z-[1500] flex items-center justify-center p-5 backdrop-blur-md animate-in fade-in duration-200" :class="overlayBg">
    <div
      class="card w-full max-w-md p-7 bg-slate-900/95 border border-slate-700/80 rounded-2xl shadow-2xl shadow-black/80 flex flex-col items-center text-center gap-4 transition-all"
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

      <!-- 3. Shot Step (Russian Roulette Revolver / Shotgun Double Fire) -->
      <template v-if="currentStep === 'shot'">
        <!-- Badge -->
        <div class="flex items-center gap-2 flex-wrap justify-center">
          <span
            v-if="stepData.double_shot"
            class="badge badge-warning badge-sm font-extrabold tracking-wider uppercase py-1 px-3 shadow-lg shadow-amber-600/40 animate-pulse"
          >
            ⚡ {{ shotStage === 1 ? (lang === 'zh' ? '第 1 枪 / 2' : 'SHOT 1 / 2') : (lang === 'zh' ? '第 2 枪 / 2' : 'SHOT 2 / 2') }}
          </span>
          <span
            v-if="stepData.armor_blocked"
            class="badge badge-info badge-sm font-extrabold tracking-wider uppercase py-1 px-3 shadow-lg shadow-sky-600/40"
          >
            {{ t('event_shot_badge_armor') }}
          </span>
          <span
            v-else
            class="badge badge-sm font-extrabold tracking-widest uppercase py-1 px-3"
            :class="stepData.fatal ? 'badge-error text-white shadow-lg shadow-rose-600/40' : 'badge-success text-white shadow-lg shadow-emerald-600/40'"
          >
            {{ stepData.fatal ? t('event_shot_badge_fatal') : t('event_shot_badge_blank') }}
          </span>
        </div>

        <!-- 6-Chamber Revolver Graphic with Double Fire / Armor Animation -->
        <div
          class="relative w-24 h-24 rounded-full border-2 border-slate-700 bg-slate-950 flex items-center justify-center shadow-inner my-2"
          :class="[
            stepData.double_shot ? (shotStage === 1 ? 'animate-recoil-1' : 'animate-recoil-2') : '',
            stepData.armor_blocked ? 'animate-armor-shield' : ''
          ]"
        >
          <!-- Central Hub -->
          <div class="absolute w-7 h-7 rounded-full bg-slate-800 border border-slate-600 flex items-center justify-center">
            <span v-if="stepData.armor_blocked" class="text-xs">🛡️</span>
            <span v-else-if="stepData.double_shot" class="text-xs text-amber-400 font-black">
              {{ shotStage === 1 ? '①' : '②' }}
            </span>
          </div>

          <!-- Double Shot Muzzle Flashes -->
          <div
            v-if="stepData.double_shot && shotStage === 1"
            class="absolute -top-3 text-amber-400 text-base animate-ping select-none"
          >
            💥
          </div>
          <div
            v-if="stepData.double_shot && shotStage === 2"
            class="absolute -top-3 text-rose-500 text-lg animate-ping select-none"
          >
            🔥
          </div>

          <!-- 6 chambers arranged in a circle -->
          <div
            v-for="idx in 6"
            :key="idx"
            class="absolute w-4 h-4 rounded-full border transition-all duration-300"
            :class="getChamberClass(idx)"
            :style="{
              transform: `rotate(${(idx - 1) * 60}deg) translate(0, -32px)`
            }"
          ></div>
        </div>

        <!-- Stage Title -->
        <h2
          class="text-2xl font-black font-serif tracking-wide transition-all"
          :class="stepData.fatal && (shotStage === 2 || !stepData.double_shot) ? 'text-rose-400 drop-shadow-[0_0_20px_rgba(244,63,94,0.8)]' : (stepData.armor_blocked ? 'text-sky-300 drop-shadow-[0_0_15px_rgba(56,189,248,0.6)]' : 'text-emerald-400 drop-shadow-[0_0_15px_rgba(16,185,129,0.5)]')"
        >
          <template v-if="stepData.double_shot">
            <span v-if="shotStage === 1" class="text-amber-300">
              ⚡ {{ lang === 'zh' ? '第 1 发扣击：空枪！' : '1st Shot: CLICK!' }}
            </span>
            <span v-else-if="stepData.fatal" class="text-rose-400">
              💥 {{ lang === 'zh' ? '第 2 发扣击：中弹！' : '2nd Shot: BANG!' }}
            </span>
            <span v-else-if="stepData.armor_blocked" class="text-sky-300">
              🛡️ {{ t('event_shot_badge_armor') }}
            </span>
            <span v-else class="text-emerald-400">
              🛡️ {{ lang === 'zh' ? '第 2 发扣击：幸存！' : '2nd Shot: SURVIVED!' }}
            </span>
          </template>
          <template v-else-if="stepData.fatal">
            {{ t('event_bang_title') }}
          </template>
          <template v-else-if="stepData.armor_blocked">
            {{ t('event_shot_badge_armor') }}
          </template>
          <template v-else>
            {{ t('event_click_title') }}
          </template>
        </h2>

        <p class="text-sm text-slate-300 leading-relaxed max-w-xs">
          <span class="font-bold text-slate-100 mr-1">{{ stepData.target }}</span>
          <template v-if="stepData.fatal && (shotStage === 2 || !stepData.double_shot)">
            {{ t('event_bang_sub') }}
          </template>
          <template v-else-if="stepData.armor_blocked">
            {{ t('event_shot_armor_sub') }}
          </template>
          <template v-else-if="stepData.double_shot">
            {{ t('event_shot_double_sub') }}
          </template>
          <template v-else>
            {{ t('event_click_sub') }}
          </template>
        </p>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';

const props = defineProps({
  currentStep: { type: String, default: '' },
  stepData: { type: Object, default: () => ({}) },
});

const { t, lang } = useI18n();
const store = useGameStore();

const shotStage = ref(1);
let shotStageTimer = null;

watch(
  () => props.currentStep,
  (val) => {
    if (val === 'shot') {
      shotStage.value = 1;
      if (props.stepData?.double_shot) {
        if (shotStageTimer) clearTimeout(shotStageTimer);
        shotStageTimer = setTimeout(() => {
          shotStage.value = 2;
        }, 1100);
      }
    } else {
      shotStage.value = 1;
      if (shotStageTimer) clearTimeout(shotStageTimer);
    }
  },
  { immediate: true }
);

function isCardHonest(card) {
  const tableCard = store.state.value.table_card;
  if (!tableCard) return true;
  return card === tableCard || card === '2';
}

function getChamberClass(idx) {
  const isFatal = props.stepData.fatal;
  const isDouble = props.stepData.double_shot;
  const isArmor = props.stepData.armor_blocked;

  if (idx === 1) {
    if (shotStage.value === 1 && isDouble) {
      return 'bg-amber-400 border-amber-200 shadow-[0_0_16px_rgba(251,191,36,1)] scale-125';
    }
    if (isArmor) {
      return 'bg-sky-500 border-sky-300 shadow-[0_0_14px_rgba(56,189,248,0.95)]';
    }
    if (isFatal && (!isDouble || shotStage.value === 1)) {
      return 'bg-rose-500 border-rose-300 shadow-[0_0_16px_rgba(244,63,94,0.95)] animate-ping';
    }
    return 'bg-emerald-500 border-emerald-300 shadow-[0_0_8px_rgba(16,185,129,0.7)]';
  }

  if (idx === 2 && isDouble) {
    if (shotStage.value === 2) {
      if (isFatal) {
        return 'bg-rose-500 border-rose-300 shadow-[0_0_18px_rgba(244,63,94,1)] scale-125 animate-ping';
      }
      return 'bg-amber-400 border-amber-200 shadow-[0_0_16px_rgba(251,191,36,1)] scale-125';
    }
    return 'bg-slate-800 border-slate-600';
  }

  return 'bg-slate-900 border-slate-700';
}

const overlayBg = computed(() => {
  if (props.currentStep === 'liar_call') return 'bg-rose-950/70';
  if (props.currentStep === 'reveal') return 'bg-slate-950/80';
  if (props.currentStep === 'shot') return props.stepData.fatal ? 'bg-rose-950/60' : 'bg-slate-950/80';
  return 'bg-slate-950/75';
});
</script>

<style scoped>
@keyframes recoil1 {
  0% { transform: scale(1) translateY(0) rotate(0deg); }
  25% { transform: scale(1.22) translateY(-10px) rotate(-4deg); }
  60% { transform: scale(0.95) translateY(3px) rotate(2deg); }
  100% { transform: scale(1) translateY(0) rotate(0deg); }
}
.animate-recoil-1 {
  animation: recoil1 0.45s cubic-bezier(0.175, 0.885, 0.32, 1.275) both;
}

@keyframes recoil2 {
  0% { transform: scale(1) translateY(0) rotate(0deg); }
  25% { transform: scale(1.28) translateY(-14px) rotate(6deg); }
  60% { transform: scale(0.92) translateY(4px) rotate(-3deg); }
  100% { transform: scale(1) translateY(0) rotate(0deg); }
}
.animate-recoil-2 {
  animation: recoil2 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275) both;
}

@keyframes armorShield {
  0% { transform: scale(1); filter: drop-shadow(0 0 0 rgba(56, 189, 248, 0)); }
  40% { transform: scale(1.25); filter: drop-shadow(0 0 20px rgba(56, 189, 248, 0.9)); }
  100% { transform: scale(1); filter: drop-shadow(0 0 6px rgba(56, 189, 248, 0.4)); }
}
.animate-armor-shield {
  animation: armorShield 0.8s ease-out both;
}
</style>

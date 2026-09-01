<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md">
    <div class="card w-full max-w-lg bg-slate-900 border border-slate-700/80 shadow-2xl shadow-black/90 rounded-2xl p-6 flex flex-col gap-5 text-left animate-in fade-in zoom-in duration-200">
      
      <!-- 弹窗顶栏 -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <div>
          <h2 class="text-base font-bold text-slate-100 tracking-wide">
            {{ t('create_room_modal_title') }}
          </h2>
          <p class="text-xs text-slate-400">
            {{ t('create_room_modal_subtitle') }}
          </p>
        </div>
        <button
          class="btn btn-ghost btn-xs btn-circle text-slate-400 hover:text-white"
          @click="$emit('close')"
        >
          ✕
        </button>
      </div>

      <!-- 1. 人数选择 (2 - 4 人) -->
      <div class="flex flex-col gap-2">
        <label class="text-xs font-bold text-slate-300 tracking-wider">
          {{ t('select_players_label') }}
        </label>
        <div class="grid grid-cols-3 gap-2.5">
          <button
            v-for="num in [2, 3, 4]"
            :key="num"
            type="button"
            class="btn btn-sm h-10 flex items-center justify-center rounded-xl border transition-all"
            :class="selectedPlayers === num 
              ? 'bg-indigo-600 border-indigo-400 text-white shadow-md shadow-indigo-600/30' 
              : 'bg-slate-950/70 border-slate-800 text-slate-300 hover:bg-slate-800 hover:border-slate-700'"
            @click="selectedPlayers = num"
          >
            <span class="text-xs font-semibold">
              {{ num === 2 ? t('players_2') : num === 3 ? t('players_3') : t('players_4') }}
            </span>
          </button>
        </div>
      </div>

      <!-- 2. 游戏模式选择 (普通 / 道具) -->
      <div class="flex flex-col gap-2">
        <label class="text-xs font-bold text-slate-300 tracking-wider">
          {{ t('select_mode_label') }}
        </label>
        <div class="grid grid-cols-2 gap-3">
          <!-- 普通模式 -->
          <div
            class="p-3.5 rounded-xl border cursor-pointer transition-all flex flex-col gap-1 select-none"
            :class="selectedMode === 'classic'
              ? 'bg-indigo-950/40 border-indigo-500 shadow-sm shadow-indigo-600/20'
              : 'bg-slate-950/60 border-slate-800 hover:border-slate-700 hover:bg-slate-900'"
            @click="onSelectMode('classic')"
          >
            <div class="flex items-center justify-between">
              <span class="font-bold text-sm text-slate-100">{{ t('mode_classic_title') }}</span>
              <span
                v-if="selectedMode === 'classic'"
                class="w-2.5 h-2.5 rounded-full bg-indigo-400"
              ></span>
            </div>
            <p class="text-[11px] text-slate-400 leading-snug">
              {{ t('mode_classic_desc') }}
            </p>
          </div>

          <!-- 道具模式 -->
          <div
            class="p-3.5 rounded-xl border cursor-pointer transition-all flex flex-col gap-1 select-none"
            :class="selectedMode === 'items'
              ? 'bg-rose-950/40 border-rose-500 shadow-sm shadow-rose-600/30'
              : 'bg-slate-950/60 border-slate-800 hover:border-slate-700 hover:bg-slate-900'"
            @click="onSelectMode('items')"
          >
            <div class="flex items-center justify-between">
              <span class="font-bold text-sm text-slate-100">{{ t('mode_items_title') }}</span>
              <span
                v-if="selectedMode === 'items'"
                class="w-2.5 h-2.5 rounded-full bg-rose-500"
              ></span>
            </div>
            <p class="text-[11px] text-slate-400 leading-snug">
              {{ t('mode_items_desc') }}
            </p>
          </div>
        </div>
      </div>

      <!-- 3. 提示警告 (道具模式激活) -->
      <transition name="fade">
        <div
          v-if="selectedMode === 'items'"
          class="flex flex-col gap-1 p-3 rounded-xl bg-amber-950/30 border border-amber-600/50 text-amber-200"
        >
          <span class="text-xs font-bold text-amber-300">
            {{ t('items_mode_warning_title') }}
          </span>
          <p class="text-xs text-amber-200/90 leading-relaxed">
            {{ t('items_mode_warning_desc') }}
          </p>
        </div>
      </transition>

      <!-- 底部操作按钮 -->
      <div class="flex items-center justify-end gap-3 border-t border-slate-800 pt-4 mt-1">
        <button
          type="button"
          class="btn btn-neutral btn-sm h-10 px-4 bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300"
          @click="$emit('close')"
        >
          {{ t('cancel') }}
        </button>

        <button
          type="button"
          class="btn btn-sm h-10 px-6 font-bold transition-all"
          :class="confirmButtonClass"
          :disabled="isButtonLocked"
          @click="onConfirm"
        >
          <span v-if="selectedMode === 'items' && countdown > 0">
            {{ t('create_room_confirm_lock', { s: countdown }) }}
          </span>
          <span v-else>
            {{ t('create_room_confirm_classic') }}
          </span>
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue';
import { useI18n } from '../composables/useI18n';

const emit = defineEmits(['confirm', 'close']);
const { t } = useI18n();

const selectedPlayers = ref(4);
const selectedMode = ref('classic');
const countdown = ref(0);
let timer = null;

function onSelectMode(mode) {
  selectedMode.value = mode;
  if (timer) {
    clearInterval(timer);
    timer = null;
  }

  if (mode === 'items') {
    countdown.value = 3;
    timer = setInterval(() => {
      countdown.value--;
      if (countdown.value <= 0) {
        clearInterval(timer);
        timer = null;
      }
    }, 1000);
  } else {
    countdown.value = 0;
  }
}

const isButtonLocked = computed(() => {
  if (selectedMode.value === 'items') {
    return countdown.value > 0;
  }
  return false;
});

const confirmButtonClass = computed(() => {
  if (selectedMode.value === 'items') {
    if (countdown.value > 0) {
      return 'btn-disabled bg-slate-800 border-slate-700 text-slate-500 cursor-not-allowed opacity-60';
    }
    return 'btn-primary bg-indigo-600 hover:bg-indigo-500 border-none text-white shadow-md shadow-indigo-600/30';
  }
  return 'btn-primary bg-indigo-600 hover:bg-indigo-500 border-none text-white shadow-md shadow-indigo-600/25';
});

function onConfirm() {
  if (isButtonLocked.value) return;
  emit('confirm', {
    mode: selectedMode.value,
    maxPlayers: selectedPlayers.value,
  });
}

onUnmounted(() => {
  if (timer) clearInterval(timer);
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>

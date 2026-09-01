<template>
  <div class="flex flex-col min-h-[90vh] max-w-5xl mx-auto w-full">
    <!-- Header Bar -->
    <HeaderBar
      :room-code="state.room_code || myRoomCode"
      :status="state.status"
      :deadline="state.deadline"
      @open-rules="$emit('open-rules')"
      @open-changelog="$emit('open-changelog')"
      @leave="onLeave"
    />

    <!-- Spectator Banner -->
    <div v-if="isSpectator" class="py-2 px-4 mb-3.5 bg-indigo-950/40 border border-dashed border-indigo-500/40 rounded-xl text-center text-xs text-indigo-300 backdrop-blur-md">
      <span>👀 {{ t('spectator_banner') }}</span>
    </div>

    <!-- Players Grid -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-4">
      <PlayerSeat
        v-for="(p, idx) in state.players"
        :key="p.id"
        :player="p"
        :is-active-turn="state.status === 'playing' && state.current_turn === idx"
        :is-me="p.id === myId || p.nickname === myNickname"
        :am-host="amHost"
        :game-status="state.status"
        @kick="kickPlayer"
      />
    </div>

    <!-- Center Table -->
    <TableArea
      :status="state.status"
      :table-card="state.table_card"
      :last-played-cnt="state.last_played_cnt"
      :winner="state.winner"
      :am-host="amHost"
      @reset="resetGame"
    />

    <!-- Double Damage Active Alert Banner -->
    <div
      v-if="state.double_damage && state.status === 'playing'"
      class="flex items-center justify-center p-2 my-1 mx-auto max-w-lg bg-slate-900 border border-indigo-500/70 rounded-xl text-indigo-200 text-xs select-none shadow-md"
    >
      <span>{{ t('double_damage_banner') }}</span>
    </div>

    <!-- Devil Items Bar (Only in Devil's Items Mode) -->
    <ItemBar
      v-if="isPlayer && state.status === 'playing' && state.game_mode === 'items' && myPlayer && myPlayer.is_alive"
      :items="myPlayer.items || []"
      :is-my-turn="isMyTurn"
      :table-has-cards="Boolean(state.last_played_cnt > 0)"
      @use-item="useItem"
    />

    <!-- Action Controls -->
    <ActionBar
      :status="state.status"
      :game-status="state.status"
      :is-my-turn="isMyTurn"
      :is-spectator="isSpectator"
      :is-alive="myPlayer ? myPlayer.is_alive : false"
      :is-player="isPlayer"
      :is-ready="myPlayer ? myPlayer.is_ready : false"
      :am-host="amHost"
      :can-play="canPlay"
      :can-call-liar="canCallLiar"
      :can-start="canStart"
      :all-ready="allReady"
      :selected-count="selectedIndexes.length"
      @toggle-ready="toggleReady"
      @start-game="startGame"
      @play-cards="playCards"
      @call-liar="callLiar"
    />

    <!-- Player Hand Area -->
    <HandArea
      v-if="isPlayer && state.status === 'playing' && myPlayer && myPlayer.is_alive"
      :hand="myHand"
      :selected-indexes="selectedIndexes"
      :is-my-turn="isMyTurn"
      @toggle-select="toggleCardSelect"
    />

    <!-- Battle Logs Panel -->
    <LogPanel :logs="state.logs" />

    <!-- Event Overlays (Liar call, Reveal, Gunshot) -->
    <EventOverlay
      :current-step="currentStep"
      :step-data="currentStepData"
    />

    <!-- Item Used Notification -->
    <Transition name="item-toast">
      <div
        v-if="itemUsedEvent"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 z-[1600] flex items-center gap-3 px-5 py-3 rounded-2xl bg-slate-900/95 border border-amber-500/60 shadow-2xl shadow-black/60 backdrop-blur-xl select-none"
      >
        <div class="item-used-icon w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/40 flex items-center justify-center text-amber-300 font-bold text-xs shrink-0">
          {{ getItemIcon(itemUsedEvent.item) }}
        </div>
        <div class="flex flex-col">
          <span class="text-[11px] font-bold text-amber-300 leading-tight">{{ itemUsedEvent.nickname }} 使用了道具</span>
          <span class="text-xs font-bold text-slate-100">{{ getItemName(itemUsedEvent.item) }}</span>
        </div>
      </div>
    </Transition>

    <!-- 放大镜查看结果弹窗 -->
    <div
      v-if="showEagleEyeModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md animate-in fade-in zoom-in duration-200"
    >
      <div class="card w-full max-w-sm bg-slate-900 border border-slate-700 shadow-2xl rounded-2xl p-6 flex flex-col items-center gap-4 text-center">
        <h3 class="text-sm font-bold text-slate-100">{{ t('eagle_eye_modal_title') }}</h3>
        <p class="text-xs text-slate-300">
          {{ t('eagle_eye_modal_desc') }}
        </p>

        <!-- Card Graphic -->
        <div class="w-20 h-28 rounded-xl border border-slate-700 bg-slate-950 shadow-md flex flex-col items-center justify-center gap-1">
          <span class="text-2xl font-bold text-amber-400">
            {{ inspectedCard }}
          </span>
          <span class="text-[10px] text-slate-400 font-mono">
            {{ inspectedCard === '2' ? 'WILD' : 'REAL' }}
          </span>
        </div>

        <button
          type="button"
          class="btn btn-sm btn-primary w-full bg-indigo-600 hover:bg-indigo-500 text-white font-bold"
          @click="closeEagleEyeModal"
        >
          {{ t('eagle_eye_close') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';
import { useGameStore } from '../composables/useGameStore';
import HeaderBar from '../components/HeaderBar.vue';
import PlayerSeat from '../components/PlayerSeat.vue';
import TableArea from '../components/TableArea.vue';
import HandArea from '../components/HandArea.vue';
import ActionBar from '../components/ActionBar.vue';
import ItemBar from '../components/ItemBar.vue';
import EventOverlay from '../components/EventOverlay.vue';
import LogPanel from '../components/LogPanel.vue';

defineEmits(['open-rules', 'open-changelog']);
const { t } = useI18n();
const {
  state,
  currentUnix,
  selectedIndexes,
  currentStep,
  currentStepData,
  myNickname,
  myId,
  myRoomCode,
  myPlayer,
  myHand,
  isPlayer,
  isSpectator,
  amHost,
  isMyTurn,
  canPlay,
  canCallLiar,
  canStart,
  allReady,
  isScreenShaking,
  inspectedCard,
  showEagleEyeModal,
  closeEagleEyeModal,
  useItem,
  itemUsedEvent,
  toggleCardSelect,
  playCards,
  callLiar,
  toggleReady,
  startGame,
  resetGame,
  kickPlayer,
  disconnect,
} = useGameStore();

function getItemIcon(item) {
  const icons = { eagle_eye: '🔍', sawed_off: '⚡', hard_liquor: '🍺', kevlar_armor: '🛡', fate_shift: '🎲' };
  return icons[item] || '?';
}

function getItemName(item) {
  const map = { eagle_eye: '放大镜', sawed_off: '猎枪', hard_liquor: '啤酒', kevlar_armor: '防弹衣', fate_shift: '骰子' };
  return map[item] || item;
}

function onLeave() {
  if (confirm('确定要离开房间吗？')) {
    disconnect();
  }
}
</script>

<style scoped>
/* Item Used Toast Transition */
.item-toast-enter-active {
  animation: itemToastIn 0.36s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}
.item-toast-leave-active {
  animation: itemToastOut 0.3s ease both;
}
@keyframes itemToastIn {
  from { opacity: 0; transform: translate(-50%, 24px) scale(0.88); }
  to   { opacity: 1; transform: translate(-50%, 0) scale(1); }
}
@keyframes itemToastOut {
  from { opacity: 1; transform: translate(-50%, 0) scale(1); }
  to   { opacity: 0; transform: translate(-50%, -12px) scale(0.92); }
}

/* Icon pulse when toast appears */
.item-used-icon {
  animation: iconPulse 0.5s ease 0.1s both;
}
@keyframes iconPulse {
  0%   { transform: scale(1); }
  40%  { transform: scale(1.35); }
  70%  { transform: scale(0.9); }
  100% { transform: scale(1); }
}
</style>

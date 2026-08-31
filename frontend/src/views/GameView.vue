<template>
  <div class="game-view">
    <!-- Header Bar -->
    <HeaderBar
      :room-code="state.room_code || myRoomCode"
      :status="state.status"
      :deadline="state.deadline"
      :current-unix="currentUnix"
      @open-rules="$emit('open-rules')"
      @leave="onLeave"
    />

    <!-- Spectator Banner -->
    <div v-if="isSpectator" class="spectator-banner glass-panel">
      <span>👀 {{ t('spectator_banner') }}</span>
    </div>

    <!-- Players Grid -->
    <div class="players-container">
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

    <!-- Action Controls -->
    <ActionBar
      :status="state.status"
      :is-player="isPlayer"
      :is-my-turn="isMyTurn"
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
import EventOverlay from '../components/EventOverlay.vue';
import LogPanel from '../components/LogPanel.vue';

defineEmits(['open-rules']);
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
  toggleCardSelect,
  playCards,
  callLiar,
  toggleReady,
  startGame,
  resetGame,
  kickPlayer,
  disconnect,
} = useGameStore();

function onLeave() {
  if (confirm('确定要离开房间吗？')) {
    disconnect();
  }
}
</script>

<style scoped>
.game-view {
  display: flex;
  flex-direction: column;
  min-height: 90vh;
}

.spectator-banner {
  padding: 8px 16px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px dashed var(--accent-blue);
  border-radius: var(--radius-md);
  text-align: center;
  font-size: 13px;
  color: #93c5fd;
  margin-bottom: 14px;
}

.players-container {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

@media (max-width: 600px) {
  .players-container {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>

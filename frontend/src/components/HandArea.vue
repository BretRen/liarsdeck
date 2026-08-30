<template>
  <div class="hand-area glass-panel">
    <div class="hand-header">
      <div class="hand-title-group">
        <span class="hand-icon">🎴</span>
        <h3>{{ t('my_hand_title') }}</h3>
        <span class="hand-count-pill">({{ hand.length }})</span>
      </div>

      <div class="selection-status">
        <span class="selection-count">
          {{ t('selected_count') }} {{ selectedIndexes.length }}/3
        </span>
        <span v-if="!isMyTurn" class="not-turn-tip">({{ t('status_waiting') }})</span>
      </div>
    </div>

    <!-- Hand Cards Carousel / Grid -->
    <div class="hand-cards-container" :class="{ 'not-my-turn': !isMyTurn }">
      <div
        v-for="(card, index) in hand"
        :key="index"
        class="playing-card hand-card"
        :class="[
          `rank-${card}`,
          {
            selected: selectedIndexes.includes(index),
            disabled: !isMyTurn,
          },
        ]"
        @click="onCardClick(index)"
      >
        <span class="card-corner">{{ card }}</span>
        <span class="rank-main">{{ card }}</span>
        <span class="card-corner-bottom">{{ card }}</span>
        <div v-if="card === '2'" class="card-wild-badge">WILD</div>
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

<style scoped>
.hand-area {
  padding: 16px 20px;
  margin-bottom: 16px;
  background: rgba(18, 21, 28, 0.9);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hand-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: 8px;
}

.hand-title-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hand-title-group h3 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.hand-count-pill {
  font-size: 12px;
  font-weight: 600;
  color: var(--accent-gold);
}

.selection-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selection-count {
  font-size: 13px;
  font-weight: 600;
  color: var(--accent-gold);
}

.not-turn-tip {
  font-size: 12px;
  color: var(--text-muted);
}

.hand-cards-container {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
  min-height: 120px;
  padding-top: 18px;
}

.hand-cards-container.not-my-turn .hand-card {
  cursor: default;
  opacity: 0.85;
}
.hand-cards-container.not-my-turn .hand-card:hover {
  transform: none;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.6);
  border-color: #544438;
}
</style>

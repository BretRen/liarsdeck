<template>
  <div class="table-area glass-panel">
    <!-- Playing State: Table True Card & Face-down Pile -->
    <template v-if="status === 'playing'">
      <div class="table-card-section">
        <div class="table-label">{{ t('table_card_label') }}</div>
        
        <div class="true-card-container">
          <div class="playing-card true-card" :class="`rank-${tableCard}`">
            <span class="card-corner">{{ tableCard }}</span>
            <span class="rank-main">{{ tableCard }}</span>
            <span class="card-corner-bottom">{{ tableCard }}</span>
            <div v-if="tableCard === '2'" class="card-wild-badge">WILD</div>
          </div>
        </div>

        <div class="table-tip">{{ t('wild_card_tip') }}</div>
      </div>

      <!-- Last Played Face-down Pile -->
      <div v-if="lastPlayedCnt > 0" class="played-pile-section">
        <div class="pile-cards">
          <div
            v-for="i in lastPlayedCnt"
            :key="i"
            class="facedown-card"
            :style="{ transform: `rotate(${(i - (lastPlayedCnt + 1) / 2) * 8}deg)` }"
          >
            <span class="facedown-pattern">🂠</span>
          </div>
        </div>
        <div class="pile-info">
          {{ lastPlayedCnt }} {{ t('cards_on_table') }}
        </div>
      </div>
    </template>

    <!-- Game Over State: Champion Display -->
    <template v-else-if="status === 'game_over'">
      <div class="game-over-section">
        <div class="trophy-icon">🏆</div>
        <div class="winner-title">{{ t('winner_label') }}</div>
        <div class="winner-name">{{ winner || '—' }}</div>

        <div class="game-over-actions">
          <button v-if="amHost" class="btn-success btn-lg" @click="$emit('reset')">
            🔄 {{ t('play_again_btn') }}
          </button>
          <p v-else class="wait-host-text">
            {{ t('wait_host_reset') }}
          </p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';

defineProps({
  status: { type: String, default: 'waiting' },
  tableCard: { type: String, default: '' },
  lastPlayedCnt: { type: Number, default: 0 },
  winner: { type: String, default: '' },
  amHost: { type: Boolean, default: false },
});

defineEmits(['reset']);
const { t } = useI18n();
</script>

<style scoped>
.table-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  margin-bottom: 16px;
  background: radial-gradient(circle at 50% 50%, rgba(24, 29, 40, 0.95) 0%, rgba(14, 17, 24, 0.95) 100%);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  min-height: 190px;
  gap: 16px;
}

.table-card-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.table-label {
  font-family: var(--font-heading);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
  color: var(--accent-gold);
}

.true-card-container {
  padding: 4px;
}

.true-card {
  width: 80px;
  height: 116px;
  font-size: 32px;
  border-width: 2.5px;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.7), 0 0 20px rgba(245, 158, 11, 0.3);
  cursor: default;
  transform: none !important;
}

.table-tip {
  font-size: 12px;
  color: var(--text-muted);
}

/* Face-down pile */
.played-pile-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
}

.pile-cards {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 48px;
}

.facedown-card {
  width: 38px;
  height: 52px;
  background: linear-gradient(135deg, #374151 0%, #1f2937 100%);
  border: 1.5px solid #4b5563;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-gold);
  font-size: 18px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
  margin: 0 -8px;
  transition: transform 0.2s ease;
}

.pile-info {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
}

/* Game Over */
.game-over-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 8px;
  animation: fadeIn 0.4s ease;
}

.trophy-icon {
  font-size: 40px;
}

.winner-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.winner-name {
  font-family: var(--font-heading);
  font-size: 1.8rem;
  font-weight: 900;
  color: var(--accent-gold);
  text-shadow: 0 0 20px rgba(245, 158, 11, 0.5);
}

.game-over-actions {
  margin-top: 10px;
}

.btn-lg {
  padding: 12px 28px;
  font-size: 15px;
}

.wait-host-text {
  font-size: 13px;
  color: var(--text-muted);
}
</style>

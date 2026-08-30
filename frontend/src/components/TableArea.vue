<template>
  <div class="table-area glass-panel">
    <!-- Playing State: Table True Card & Face-down Stack -->
    <template v-if="status === 'playing'">
      <div class="table-card-section">
        <div class="table-title">{{ t('table_card_label') }}</div>
        
        <div class="true-card-placard">
          <div class="playing-card true-card" :class="`rank-${tableCard}`">
            <span class="card-corner">{{ tableCard }}</span>
            <span class="rank-main">{{ tableCard }}</span>
            <span class="card-corner-bottom">{{ tableCard }}</span>
            <div v-if="tableCard === '2'" class="card-wild-badge">WILD</div>
          </div>
        </div>

        <div class="table-tip">{{ t('wild_card_tip') }}</div>
      </div>

      <!-- Face-down Played Pile -->
      <div v-if="lastPlayedCnt > 0" class="played-stack-section">
        <div class="pile-cards">
          <div
            v-for="i in lastPlayedCnt"
            :key="i"
            class="facedown-card"
            :style="{ transform: `rotate(${(i - (lastPlayedCnt + 1) / 2) * 6}deg)` }"
          >
            <div class="card-back-pattern"></div>
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
        <div class="winner-title">{{ t('winner_label') }}</div>
        <div class="winner-name">{{ winner || '—' }}</div>

        <div class="game-over-actions">
          <button v-if="amHost" class="btn-primary btn-lg" @click="$emit('reset')">
            {{ t('play_again_btn') }}
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
  padding: 20px;
  margin-bottom: 14px;
  background: #11141c;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  min-height: 180px;
  gap: 14px;
}

.table-card-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.table-title {
  font-family: var(--font-serif);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
  color: var(--gold-accent);
  text-transform: uppercase;
}

.true-card-placard {
  padding: 4px;
}

.true-card {
  width: 76px;
  height: 110px;
  font-size: 32px;
  border-width: 2px;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.6), 0 0 12px rgba(212, 175, 55, 0.3);
  cursor: default;
  transform: none !important;
}

.table-tip {
  font-size: 12px;
  color: var(--text-muted);
}

/* Face-down stack */
.played-stack-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.pile-cards {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 44px;
}

.facedown-card {
  width: 36px;
  height: 50px;
  background: #2b1810;
  border: 1.5px solid #543324;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.5);
  margin: 0 -8px;
  padding: 3px;
}

.card-back-pattern {
  width: 100%;
  height: 100%;
  border: 1px dashed rgba(212, 175, 55, 0.4);
  border-radius: 2px;
  background: repeating-linear-gradient(
    45deg,
    rgba(0,0,0,0.1),
    rgba(0,0,0,0.1) 2px,
    transparent 2px,
    transparent 4px
  );
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
  gap: 6px;
}

.winner-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.winner-name {
  font-family: var(--font-serif);
  font-size: 2rem;
  font-weight: 900;
  color: var(--gold-accent);
  text-shadow: 0 0 20px var(--gold-glow);
}

.game-over-actions {
  margin-top: 8px;
}

.btn-lg {
  padding: 11px 26px;
  font-size: 14px;
}

.wait-host-text {
  font-size: 13px;
  color: var(--text-muted);
}
</style>

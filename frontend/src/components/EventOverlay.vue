<template>
  <div v-if="currentStep" class="event-overlay" :class="currentStep">
    <!-- Muzzle Flash on fatal shot -->
    <div v-if="currentStep === 'shot' && stepData.fatal" class="muzzle-flash"></div>

    <div class="event-card glass-panel" :class="{ shake: currentStep === 'shot' && stepData.fatal }">
      <!-- 1. Liar Call Step -->
      <template v-if="currentStep === 'liar_call'">
        <div class="event-badge pulse-liar">CHALLENGE</div>
        <h2 class="event-title">{{ t('event_liar_alert') }}</h2>
        <div class="call-details">
          <span class="highlight caller">{{ stepData.caller }}</span>
          <span class="claim-text">{{ t('event_calls_out') }}</span>
          <span class="highlight accused">{{ stepData.accused }}</span>
          <span class="claim-text">{{ t('event_liar_claim') }}</span>
        </div>
      </template>

      <!-- 2. Reveal Cards Step -->
      <template v-if="currentStep === 'reveal'">
        <div class="event-badge">VERIFICATION</div>
        <h2 class="event-title">{{ stepData.accused }} - {{ t('event_cards_revealed') }}</h2>
        <div class="revealed-cards-row">
          <div
            v-for="(card, i) in stepData.cards"
            :key="i"
            class="playing-card reveal-card"
            :class="`rank-${card}`"
            :style="{ animationDelay: `${i * 0.12}s` }"
          >
            <span class="card-corner">{{ card }}</span>
            <span class="rank-main">{{ card }}</span>
            <span class="card-corner-bottom">{{ card }}</span>
            <div v-if="card === '2'" class="card-wild-badge">WILD</div>
          </div>
        </div>
      </template>

      <!-- 3. Shot Step -->
      <template v-if="currentStep === 'shot'">
        <div class="event-badge" :class="{ 'badge-fatal': stepData.fatal, 'badge-blank': !stepData.fatal }">
          {{ stepData.fatal ? 'FATAL ROUND' : 'DRY FIRE' }}
        </div>
        <h2 class="event-title" :class="{ fatal: stepData.fatal, blank: !stepData.fatal }">
          {{ stepData.fatal ? t('event_bang_title') : t('event_click_title') }}
        </h2>
        <p class="shot-sub">
          <span class="highlight">{{ stepData.target }}</span>
          {{ stepData.fatal ? t('event_bang_sub') : t('event_click_sub') }}
        </p>
      </template>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../composables/useI18n';

defineProps({
  currentStep: { type: String, default: '' },
  stepData: { type: Object, default: () => ({}) },
});

const { t } = useI18n();
</script>

<style scoped>
.event-overlay {
  position: fixed;
  inset: 0;
  z-index: 1500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  backdrop-filter: blur(8px);
  animation: fadeIn 0.2s ease;
}

.event-overlay.liar_call {
  background: rgba(43, 10, 10, 0.7);
}

.event-overlay.reveal {
  background: rgba(13, 15, 20, 0.8);
}

.event-overlay.shot {
  background: rgba(5, 6, 8, 0.85);
}

.event-card {
  width: 100%;
  max-width: 440px;
  padding: 28px 24px;
  background: #141720;
  border: 1px solid var(--border-brass);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 14px;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.8);
}

.event-badge {
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 1.5px;
  padding: 3px 10px;
  border-radius: var(--radius-sm);
  background: #2b1414;
  border: 1px solid #742a2a;
  color: #ff8787;
}

.badge-fatal {
  background: #3b1010;
  border-color: #e03131;
  color: #ffffff;
}

.badge-blank {
  background: rgba(43, 138, 62, 0.25);
  border-color: #2b8a3e;
  color: #51cf66;
}

.event-title {
  font-family: var(--font-serif);
  font-size: 1.4rem;
  font-weight: 900;
  color: var(--text-primary);
  margin: 0;
}

.event-title.fatal {
  color: #ff8787;
  text-shadow: 0 0 16px rgba(239, 68, 68, 0.5);
}

.event-title.blank {
  color: #51cf66;
  text-shadow: 0 0 16px rgba(81, 207, 102, 0.5);
}

.call-details {
  font-size: 1rem;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  flex-wrap: wrap;
}

.highlight {
  font-weight: 700;
  color: var(--text-primary);
}

.highlight.caller {
  color: #748ffc;
}

.highlight.accused {
  color: #ff8787;
}

/* Revealed cards */
.revealed-cards-row {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 6px;
}

.reveal-card {
  width: 70px;
  height: 102px;
  font-size: 28px;
  cursor: default;
  transform: none;
  animation: cardFlip 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
}

@keyframes cardFlip {
  0% { transform: scale(0.6) rotateY(90deg); opacity: 0; }
  100% { transform: scale(1) rotateY(0deg); opacity: 1; }
}

.shot-sub {
  font-size: 0.95rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>

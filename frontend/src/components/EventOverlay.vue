<template>
  <div v-if="currentStep" class="event-overlay" :class="currentStep">
    <!-- Muzzle Flash overlay on fatal shot -->
    <div v-if="currentStep === 'shot' && stepData.fatal" class="muzzle-flash"></div>

    <div class="event-card glass-panel" :class="{ shake: currentStep === 'shot' && stepData.fatal }">
      <!-- 1. Liar Call Step -->
      <template v-if="currentStep === 'liar_call'">
        <div class="siren-icon pulse-liar">🚨</div>
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
        <h2 class="event-title">{{ stepData.accused }} - {{ t('event_cards_revealed') }}</h2>
        <div class="revealed-cards-row">
          <div
            v-for="(card, i) in stepData.cards"
            :key="i"
            class="playing-card reveal-card"
            :class="`rank-${card}`"
            :style="{ animationDelay: `${i * 0.15}s` }"
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
        <div class="shot-icon">
          {{ stepData.fatal ? '💥' : '💨' }}
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
  backdrop-filter: blur(12px);
  animation: fadeIn 0.25s ease;
}

.event-overlay.liar_call {
  background: rgba(185, 28, 28, 0.25);
}

.event-overlay.reveal {
  background: rgba(10, 11, 14, 0.75);
}

.event-overlay.shot {
  background: rgba(0, 0, 0, 0.8);
}

.event-card {
  width: 100%;
  max-width: 480px;
  padding: 32px 24px;
  background: #141722;
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-xl);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.8), 0 0 30px rgba(0, 0, 0, 0.5);
}

.siren-icon {
  font-size: 54px;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.15);
}

.event-title {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 900;
  color: var(--text-primary);
}

.event-title.fatal {
  color: #ef4444;
  text-shadow: 0 0 20px rgba(239, 68, 68, 0.6);
}

.event-title.blank {
  color: #34d399;
  text-shadow: 0 0 20px rgba(52, 211, 153, 0.6);
}

.call-details {
  font-size: 1.05rem;
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
  color: #60a5fa;
}

.highlight.accused {
  color: #f87171;
}

/* Revealed cards */
.revealed-cards-row {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 10px;
}

.reveal-card {
  width: 74px;
  height: 106px;
  font-size: 30px;
  cursor: default;
  transform: none;
  animation: cardFlip 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) backwards;
}

@keyframes cardFlip {
  0% { transform: scale(0.5) rotateY(90deg); opacity: 0; }
  100% { transform: scale(1) rotateY(0deg); opacity: 1; }
}

.shot-icon {
  font-size: 60px;
}

.shot-sub {
  font-size: 1rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
